import { PrismaClient } from '@prisma/client';
import * as dotenv from 'dotenv';
import path from 'path';

dotenv.config({ path: path.resolve(process.cwd(), '.env') });

const prisma = new PrismaClient();

async function main() {
  const defaults = {
    currency: 'HBAR' as const,
  };

  const defaultVerifierAddr = process.env.DEFAULT_VERIFIER_ADDRESS || '0x000000000000000000000000000000000000dEaD';

  const verifiers = [
    {
      name: 'Alpha Verifier',
      onchainAddress: defaultVerifierAddr,
      fee: BigInt(0),
      currency: defaults.currency,
      rating: 4.8,
      metadata: {
        specialties: ['KYC', 'Document Verification'],
        languages: ['en'],
        uri: 'ipfs://bafy...alpha',
      },
    },
    {
      name: 'Beta Verifier',
      onchainAddress: '0x000000000000000000000000000000000000bEEF',
      fee: BigInt(0),
      currency: defaults.currency,
      rating: 4.5,
      metadata: {
        specialties: ['Address Verification'],
        languages: ['en', 'es'],
        uri: 'ipfs://bafy...beta',
      },
    },
    {
      name: 'Gamma Verifier',
      onchainAddress: '0x000000000000000000000000000000000000CAFE',
      fee: BigInt(0),
      currency: defaults.currency,
      rating: 4.9,
      metadata: {
        specialties: ['Enhanced Due Diligence'],
        languages: ['en', 'fr'],
        uri: 'ipfs://bafy...gamma',
      },
    },
  ];

  for (const v of verifiers) {
    await prisma.verifier.upsert({
      where: { onchainAddress: v.onchainAddress },
      update: {
        name: v.name,
        fee: v.fee,
        currency: v.currency,
        rating: v.rating ?? undefined,
        metadata: v.metadata as any,
        status: 'active',
      },
      create: v as any,
    });
  }

  const count = await prisma.verifier.count();
  console.log(`Seeded verifiers. Total count: ${count}`);
}

main()
  .then(async () => {
    await prisma.$disconnect();
    process.exit(0);
  })
  .catch(async (e) => {
    console.error(e);
    await prisma.$disconnect();
    process.exit(1);
  });