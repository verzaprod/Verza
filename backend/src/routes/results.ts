import { Router } from 'express';
import { authMiddleware } from '../middleware/auth';
import { prisma } from '../db/client';
import { env } from '../config/env';
import { getContracts } from '../contracts';

const router = Router();

router.get('/results/:escrowId', authMiddleware, async (req, res) => {
  const escrow = await prisma.escrow.findUnique({ where: { id: req.params.escrowId }, include: { credential: true, user: true, verifier: true } });
  if (!escrow) return res.status(404).json({ error: 'Escrow not found' });
  const verified = !!escrow.credential;

  const credential = verified ? {
    type: escrow.credential!.type,
    issuer: { name: 'Verza', did: escrow.user.did ?? null },
    subject: { did: escrow.user.did ?? null },
    chain: { chainId: env.CHAIN_ID, registry: getContracts().addresses.registry, tokenId: escrow.credential!.tokenId.toString(), tokenURI: escrow.credential!.tokenUri },
    issuedAt: escrow.credential!.issuedAt.toISOString(),
    attributes: {}
  } : undefined;

  res.json({
    escrowId: escrow.id,
    verified,
    credential,
    verifier: { id: escrow.verifierId, name: escrow.verifier.name, rating: escrow.verifier.rating }
  });
});

export default router;