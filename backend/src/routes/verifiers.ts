import { Router } from 'express';
import { prisma } from '../db/client';
import { getContracts } from '../contracts';
import { authMiddleware } from '../middleware/auth';

const router = Router();

function serializeVerifier(v: any) {
  return {
    ...v,
    fee: v.fee ? v.fee.toString() : '0',
  };
}

router.get('/', authMiddleware, async (_req, res) => {
  const dbVerifiers = await prisma.verifier.findMany({ orderBy: { createdAt: 'desc' } });
  const { marketplace } = getContracts();
  const withOnchain = await Promise.all(dbVerifiers.map(async (v) => {
    try {
      const [, stakedAmount, reputationScore, totalVerifications, successfulVerifications, lastActivityTimestamp, isActive, metadata] = await marketplace.getVerifier(v.onchainAddress);
      const fee = await marketplace.calculateVerificationFee(v.onchainAddress);
      return {
        ...serializeVerifier(v),
        onchain: {
          stakedAmount: stakedAmount.toString(),
          reputationScore: Number(reputationScore),
          totalVerifications: Number(totalVerifications),
          successfulVerifications: Number(successfulVerifications),
          lastActivityTimestamp: Number(lastActivityTimestamp),
          isActive,
          metadata,
          fee: fee.toString(),
        }
      };
    } catch {
      return { ...serializeVerifier(v), onchain: null };
    }
  }));
  res.json({ verifiers: withOnchain });
});

router.get('/:id', authMiddleware, async (req, res) => {
  const v = await prisma.verifier.findUnique({ where: { id: req.params.id } });
  if (!v) return res.status(404).json({ error: 'Verifier not found' });
  const { marketplace } = getContracts();
  try {
    const [, stakedAmount, reputationScore, totalVerifications, successfulVerifications, lastActivityTimestamp, isActive, metadata] = await marketplace.getVerifier(v.onchainAddress);
    const fee = await marketplace.calculateVerificationFee(v.onchainAddress);
    return res.json({
      ...serializeVerifier(v),
      onchain: {
        stakedAmount: stakedAmount.toString(),
        reputationScore: Number(reputationScore),
        totalVerifications: Number(totalVerifications),
        successfulVerifications: Number(successfulVerifications),
        lastActivityTimestamp: Number(lastActivityTimestamp),
        isActive,
        metadata,
        fee: fee.toString(),
      }
    });
  } catch {
    return res.json({ ...serializeVerifier(v), onchain: null });
  }
});

export default router;