import { Router } from 'express';
import { z } from 'zod';
import { authMiddleware } from '../middleware/auth';
import { prisma } from '../db/client';
import { getContracts } from '../contracts';
import { genRequestId } from '../utils/ids';
import { env } from '../config/env';
import { AddressLike, Contract, Interface, JsonRpcProvider, parseEther, zeroPadValue } from 'ethers';

const router = Router();

const initiateSchema = z.object({
  verifier_id: z.string(),
  currency: z.string().default('tDUST'),
  auto_release_hours: z.number().int().min(1).max(168).optional(),
  wallet_address: z.string().optional(), // for non-custodial; server may already have mapping
});

type InitiateBody = z.infer<typeof initiateSchema>;

router.post('/initiate', authMiddleware, async (req, res) => {
  const parse = initiateSchema.safeParse(req.body);
  if (!parse.success) return res.status(400).json({ error: parse.error.flatten() });
  const body = parse.data as InitiateBody;

  // Ensure user exists
  const user = await prisma.user.upsert({
    where: { clerkUserId: req.user!.id },
    update: {},
    create: { clerkUserId: req.user!.id, walletAddress: body.wallet_address },
  });
  if (body.wallet_address && user.walletAddress !== body.wallet_address) {
    await prisma.user.update({ where: { id: user.id }, data: { walletAddress: body.wallet_address } });
  }

  // Resolve verifier by ID or onchain address
  let verifier = await prisma.verifier.findUnique({ where: { id: body.verifier_id } });
  if (!verifier && body.verifier_id.startsWith('0x') && body.verifier_id.length === 42) {
    verifier = await prisma.verifier.findUnique({ where: { onchainAddress: body.verifier_id } });
    if (!verifier) {
      // Create placeholder verifier record
      verifier = await prisma.verifier.create({ data: { name: 'Verifier', onchainAddress: body.verifier_id, currency: body.currency } });
    }
  }
  if (!verifier) return res.status(404).json({ error: 'Verifier not found' });

  const { provider, marketplace, escrow, iface, addresses } = getContracts();

  // Calculate on-chain verification fee and prepare transaction
  let verificationFee: bigint;
  try {
    verificationFee = await marketplace.calculateVerificationFee(verifier.onchainAddress);
  } catch (e) {
    return res.status(400).json({ error: 'Failed to calculate verification fee' });
  }

  const walletAddress = user.walletAddress || body.wallet_address;
  if (env.ESCROW_MODE === 'noncustodial') {
    if (!walletAddress) return res.status(400).json({ error: 'Missing user wallet_address for non-custodial flow' });

    const now = BigInt(Math.floor(Date.now() / 1000));
    const nonce = BigInt(Date.now());
    const requestId = genRequestId(walletAddress, verifier.onchainAddress, nonce, now);

    const data = iface.escrow.encodeFunctionData('createEscrow', [requestId, verifier.onchainAddress]);

    // Gas estimate
    let gasLimit: bigint | undefined;
    try {
      gasLimit = await provider.estimateGas({
        from: walletAddress,
        to: addresses.escrow,
        data,
        value: verificationFee,
      } as any);
    } catch {
      // provide a safe default if estimation fails
      gasLimit = 200000n;
    }

    // Persist escrow record
    await prisma.escrow.create({
      data: {
        id: requestId,
        requestId,
        userId: user.id,
        verifierId: verifier.id,
        amount: verificationFee,
        currency: body.currency,
        autoReleaseAt: body.auto_release_hours ? new Date(Date.now() + body.auto_release_hours * 3600 * 1000) : null,
        status: 'submitted',
      }
    });

    return res.json({
      escrow_id: requestId,
      chain: { chainId: env.CHAIN_ID, rpcUrl: env.RPC_URL },
      tx: {
        to: addresses.escrow,
        data,
        value: verificationFee.toString(),
        gasLimit: gasLimit.toString(),
      }
    });
  } else {
    // Custodial: server submits the tx using signer
    const signer = (escrow.runner as any);
    if (!signer || !('provider' in signer)) {
      return res.status(500).json({ error: 'Server signer not configured' });
    }

    const now = BigInt(Math.floor(Date.now() / 1000));
    const nonce = BigInt(Date.now());
    const requestId = genRequestId(signer.address, verifier.onchainAddress, nonce, now);

    try {
      const tx = await escrow.createEscrow(requestId, verifier.onchainAddress, { value: verificationFee });
      const receipt = await tx.wait();

      await prisma.escrow.create({
        data: {
          id: requestId,
          requestId,
          userId: user.id,
          verifierId: verifier.id,
          amount: verificationFee,
          currency: body.currency,
          txHash: receipt?.hash,
          status: 'submitted',
        }
      });

      return res.json({ escrow_id: requestId, status: 'submitted', tx_hash: receipt?.hash });
    } catch (e: any) {
      return res.status(500).json({ error: 'Escrow submission failed', details: e?.message });
    }
  }
});

router.get('/status/:escrowId', authMiddleware, async (req, res) => {
  const escrow = await prisma.escrow.findUnique({ where: { id: req.params.escrowId }, include: { verification: true, credential: true } });
  if (!escrow) return res.status(404).json({ error: 'Escrow not found' });

  const steps = [
    { key: 'created', status: ['submitted','in_progress','completed','refunded','cancelled'].includes(escrow.status) ? 'done' : 'pending' },
    { key: 'funds_locked', status: ['in_progress','completed','refunded'].includes(escrow.status) ? 'done' : 'pending' },
    { key: 'verification', status: escrow.verification ? (escrow.verification.status === 'completed' ? 'done' : 'in_progress') : 'pending' },
    { key: 'fraud_check', status: ['completed','refunded'].includes(escrow.status) ? 'done' : 'pending' },
    { key: 'settlement', status: ['completed','refunded'].includes(escrow.status) ? 'done' : 'pending' },
  ];

  res.json({ escrowId: escrow.id, status: escrow.status, steps });
});

export default router;
