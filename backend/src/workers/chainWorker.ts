import { getContracts } from '../contracts';
import { logger } from '../logger';
import { prisma } from '../db/client';

export async function startChainWorker() {
  const { escrow, provider } = getContracts();
  logger.info('Chain worker starting: subscribing to Escrow events');

  escrow.on('EscrowCreated', async (requestId: string, user: string, verifier: string, amount: bigint) => {
    try {
      await prisma.escrow.upsert({
        where: { id: requestId },
        update: { status: 'submitted', amount },
        create: {
          id: requestId,
          requestId,
          amount,
          currency: 'HBAR',
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
      await prisma.escrow.update({ where: { id: requestId }, data: { status: 'completed' } });
      logger.info({ requestId }, 'FundsReleased processed');
    } catch (e) {
      logger.error({ e, requestId }, 'Failed to process FundsReleased');
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

  provider.on('error', (err) => {
    logger.error({ err }, 'Provider error in chain worker');
  });
}