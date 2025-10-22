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

async function fetchOnchainMetadata(verifierAddress: string) {
  const { marketplace } = getContracts();
  try {
    // Use the new getVerifierDetails function for comprehensive data
    const [name, metadataURI, baseFee, isActive, reputationScore, totalVerifications, successfulVerifications, stakedAmount] = 
      await marketplace.getVerifierDetails(verifierAddress);
    
    return {
      name,
      metadataURI,
      baseFee: baseFee.toString(),
      isActive,
      reputationScore: Number(reputationScore),
      totalVerifications: Number(totalVerifications),
      successfulVerifications: Number(successfulVerifications),
      stakedAmount: stakedAmount.toString(),
    };
  } catch (error) {
    // Fallback to original getVerifier if new function not available
    try {
      const [, stakedAmount, reputationScore, totalVerifications, successfulVerifications, lastActivityTimestamp, isActive, metadata] = 
        await marketplace.getVerifier(verifierAddress);
      const fee = await marketplace.calculateVerificationFee(verifierAddress);
      
      return {
        name: null, // Will be extracted from metadata or use fallback
        metadataURI: metadata,
        baseFee: fee.toString(),
        isActive,
        reputationScore: Number(reputationScore),
        totalVerifications: Number(totalVerifications),
        successfulVerifications: Number(successfulVerifications),
        stakedAmount: stakedAmount.toString(),
        lastActivityTimestamp: Number(lastActivityTimestamp),
      };
    } catch {
      return null;
    }
  }
}

router.get('/', authMiddleware, async (_req, res) => {
  const dbVerifiers = await prisma.verifier.findMany({ orderBy: { createdAt: 'desc' } });
  
  const withOnchain = await Promise.all(dbVerifiers.map(async (v) => {
    const onchainData = await fetchOnchainMetadata(v.onchainAddress);
    
    if (!onchainData) {
      return { ...serializeVerifier(v), onchain: null };
    }
    
    return {
      ...serializeVerifier(v),
      onchain: onchainData,
      // Add resolved metadata if name is available from contract
      ...(onchainData.name && { onchainResolved: { name: onchainData.name } })
    };
  }));
  
  res.json({ verifiers: withOnchain });
});

router.get('/:id', authMiddleware, async (req, res) => {
  const v = await prisma.verifier.findUnique({ where: { id: req.params.id } });
  if (!v) return res.status(404).json({ error: 'Verifier not found' });
  
  const onchainData = await fetchOnchainMetadata(v.onchainAddress);
  
  if (!onchainData) {
    return res.json({ ...serializeVerifier(v), onchain: null });
  }
  
  return res.json({
    ...serializeVerifier(v),
    onchain: onchainData,
    // Add resolved metadata if name is available from contract
    ...(onchainData.name && { onchainResolved: { name: onchainData.name } })
  });
});

export default router;