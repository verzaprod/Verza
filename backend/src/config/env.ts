import * as dotenv from 'dotenv';
import path from 'path';

dotenv.config({ path: path.resolve(process.cwd(), '.env') });

type EscrowMode = 'noncustodial' | 'custodial';

function toBool(val: any, def = false): boolean {
  if (val === undefined) return def;
  return ['1', 'true', 'yes', 'on'].includes(String(val).toLowerCase());
}

export const env = {
  NODE_ENV: process.env.NODE_ENV ?? 'development',
  PORT: Number(process.env.PORT ?? 3001),
  RPC_URL: process.env.RPC_URL ?? 'http://localhost',
  CHAIN_ID: Number(process.env.CHAIN_ID ?? 0),
  NETWORK: process.env.NETWORK ?? 'midnightTestnet',
  ESCROW_MODE: (process.env.ESCROW_MODE ?? 'noncustodial') as EscrowMode,
  AUTH_BYPASS: toBool(process.env.AUTH_BYPASS ?? 'true'),
  CLERK_JWKS_URL: process.env.CLERK_JWKS_URL ?? '',
  SERVER_PRIVATE_KEY: process.env.SERVER_PRIVATE_KEY ?? '',
  ENABLE_WORKER: toBool(process.env.ENABLE_WORKER ?? 'false'),
  STORAGE_PROVIDER: process.env.STORAGE_PROVIDER ?? 'local',
  CONTRACTS_CONFIG_PATH: process.env.CONTRACTS_CONFIG_PATH ?? path.join('..','contracts','contract-config.json'),
  DEFAULT_VERIFIER_ADDRESS: process.env.DEFAULT_VERIFIER_ADDRESS ?? '',
  // Optional contract address overrides
  ESCROW_ADDRESS: process.env.ESCROW_ADDRESS,
  VC_REGISTRY_ADDRESS: process.env.VC_REGISTRY_ADDRESS,
  VERIFIER_MARKETPLACE_ADDRESS: process.env.VERIFIER_MARKETPLACE_ADDRESS,
  MIDNIGHT_COMPACT_PATH: process.env.MIDNIGHT_COMPACT_PATH ?? '/home/ekko/.local/bin/compact',
  ENABLE_MIDNIGHT_WORKER: toBool(process.env.ENABLE_MIDNIGHT_WORKER ?? 'false'),
};

export function projectRoot(): string {
  return process.cwd();
}

export function contractsConfigPath(): string {
  return path.resolve(projectRoot(), env.CONTRACTS_CONFIG_PATH);
}
