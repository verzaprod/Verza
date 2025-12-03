import { getContracts } from '../contracts';
import { logger } from '../logger';
import { prisma } from '../db/client';
import { keccak256, toUtf8Bytes } from 'ethers';

export async function startChainWorker() {
  const { escrow, provider, registry, signer, iface } = getContracts();
  logger.info('Chain worker starting: subscribing to Escrow and VCRegistry events');

  // Escrow lifecycle events
  escrow.on('EscrowCreated', async (requestId: string, user: string, verifier: string, amount: bigint) => {
    try {
      await prisma.escrow.upsert({
        where: { id: requestId },
        update: { status: 'submitted', amount },
        create: {
          id: requestId,
          requestId,
          amount,
          currency: 'tDUST',
          status: 'submitted',
          user: { connectOrCreate: { where: { clerkUserId: user.toLowerCase() }, create: { clerkUserId: user.toLowerCase(), walletAddress: user } } },
          verifier: { connectOrCreate: { where: { onchainAddress: verifier }, create: { name: 'Verifier', onchainAddress: verifier, currency: 'HBAR' } } },
        },
      });
      logger.info({ requestId }, 'EscrowCreated processed');
    } catch (e) {
      logger.error({ e, requestId }, 'Failed to process EscrowCreated');
    }
  });

  escrow.on('FundsLocked', async (requestId: string, amount: bigint, expiresAt: bigint) => {
    try {
      await prisma.escrow.update({
        where: { id: requestId },
        data: { status: 'in_progress', autoReleaseAt: new Date(Number(expiresAt) * 1000) },
      });
      logger.info({ requestId }, 'FundsLocked processed');
    } catch (e) {
      logger.error({ e, requestId }, 'Failed to process FundsLocked');
    }
  });

  escrow.on('FundsReleased', async (requestId: string) => {
    try {
      // Mark escrow completed
      const escrowRecord = await prisma.escrow.update({ where: { id: requestId }, data: { status: 'completed' }, include: { user: true, credential: true } });
      logger.info({ requestId }, 'FundsReleased processed');

      // Attempt VC issuance when settlement completes
      if (!signer) {
        logger.warn({ requestId }, 'Skipping VC issuance: server signer not configured');
        return;
      }
      if (escrowRecord.credential) {
        logger.info({ requestId }, 'Credential already exists, skipping issuance');
        return;
      }

      const holder = escrowRecord.user.walletAddress;
      const hederaDID = escrowRecord.user.did;
      if (!holder) {
        logger.warn({ requestId }, 'Skipping VC issuance: user walletAddress missing');
        return;
      }
      if (!hederaDID) {
        logger.warn({ requestId }, 'Skipping VC issuance: user DID missing');
        return;
      }

      // Derive a deterministic VC hash from escrow context
      const vcPayload = JSON.stringify({ escrowId: requestId, userId: escrowRecord.userId, verifierId: escrowRecord.verifierId });
      const vcHash = keccak256(toUtf8Bytes(vcPayload));

      // Minimal metadata URI inline (data URI)
      const meta = {
        name: 'Verza Identity Credential',
        description: 'Issued upon successful verification and settlement.',
        escrowId: requestId,
        did: hederaDID,
        holder,
        network: 'midnightTestnet',
        issuedAt: new Date().toISOString(),
      };
      const metadataURI = `data:application/json;base64,${Buffer.from(JSON.stringify(meta)).toString('base64')}`;

      try {
        const tx = await registry.issueCredential(
          holder,
          hederaDID,
          vcHash,
          0, // CredentialType.Identity
          metadataURI,
          '', // schemaURI empty to avoid approval requirement
          [], // claims
          0   // expirationPeriod (0 => default or none)
        );
        const receipt = await tx.wait();

        // Parse logs to extract tokenId
        let tokenId: bigint | null = null;
        for (const log of receipt.logs) {
          try {
            const parsed = iface.registry.parseLog({ topics: Array.from(log.topics), data: log.data });
            if (parsed?.name === 'CredentialIssued') {
              tokenId = parsed.args[0] as bigint;
              break;
            }
          } catch {}
        }

        if (tokenId === null) {
          logger.warn({ requestId }, 'Issued VC but tokenId not parsed from logs');
        }

        await prisma.credential.create({
          data: {
            escrowId: requestId,
            userId: escrowRecord.userId,
            tokenId: tokenId ?? 0n,
            tokenUri: metadataURI,
            type: 'identity',
          }
        });
        logger.info({ requestId, tokenId: tokenId?.toString() }, 'VC issuance persisted');
      } catch (e: any) {
        logger.error({ requestId, err: e?.message }, 'VC issuance failed');
      }
    } catch (e) {
      logger.error({ e, requestId }, 'Failed to finalize FundsReleased');
    }
  });

  escrow.on('RefundIssued', async (requestId: string) => {
    try {
      await prisma.escrow.update({ where: { id: requestId }, data: { status: 'refunded' } });
      logger.info({ requestId }, 'RefundIssued processed');
    } catch (e) {
      logger.error({ e, requestId }, 'Failed to process RefundIssued');
    }
  });

  escrow.on('EscrowCancelled', async (requestId: string) => {
    try {
      await prisma.escrow.update({ where: { id: requestId }, data: { status: 'cancelled' } });
      logger.info({ requestId }, 'EscrowCancelled processed');
    } catch (e) {
      logger.error({ e, requestId }, 'Failed to process EscrowCancelled');
    }
  });

  // VCRegistry events
  registry.on('CredentialIssued', async (tokenId: bigint, vcHash: string, issuer: string, holder: string, hederaDID: string) => {
    logger.info({ tokenId: tokenId.toString(), holder, issuer, hederaDID }, 'CredentialIssued observed');
    // Optional: we could reconcile with DB here if issuance occurred out-of-band
  });

  provider.on('error', (err) => {
    logger.error({ err }, 'Provider error in chain worker');
  });
}
