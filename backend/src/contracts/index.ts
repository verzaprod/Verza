import { Contract, JsonRpcProvider, Wallet, Interface, getAddress, formatEther } from 'ethers';
import fs from 'fs';
import path from 'path';
import { env, contractsConfigPath } from '../config/env';
import { logger } from '../logger';

export type Contracts = {
  provider: JsonRpcProvider;
  signer?: Wallet;
  addresses: {
    escrow: string;
    registry: string;
    marketplace: string;
  };
  escrow: Contract;
  registry: Contract;
  marketplace: Contract;
  iface: {
    escrow: Interface;
    registry: Interface;
    marketplace: Interface;
  };
};

function loadArtifact(relPath: string): any {
  const abs = path.resolve(process.cwd(), relPath);
  const raw = fs.readFileSync(abs, 'utf-8');
  return JSON.parse(raw);
}

function readAddresses() {
  const configPath = contractsConfigPath();
  const raw = fs.readFileSync(configPath, 'utf-8');
  const cfg = JSON.parse(raw);
  const key = env.NETWORK;
  const net = cfg.networks?.[key];
  if (!net) throw new Error(`${key} not found in contract-config.json`);
  return {
    escrow: env.ESCROW_ADDRESS ?? net.escrowContract,
    registry: env.VC_REGISTRY_ADDRESS ?? net.vcRegistry,
    marketplace: env.VERIFIER_MARKETPLACE_ADDRESS ?? net.verifierMarketplace,
    chainId: Number(net.chainId ?? env.CHAIN_ID),
  };
}

export function getContracts(): Contracts {
  const provider = new JsonRpcProvider(env.RPC_URL, env.CHAIN_ID, { batchMaxCount: 1 });
  const signer = env.SERVER_PRIVATE_KEY ? new Wallet(env.SERVER_PRIVATE_KEY, provider) : undefined;

  const addresses = readAddresses();

  const escrowArtifact = loadArtifact(path.join('..','contracts','artifacts','contracts','EscrowContract.sol','EscrowContract.json'));
  const registryArtifact = loadArtifact(path.join('..','contracts','artifacts','contracts','VCRegistry.sol','VCRegistry.json'));
  const marketplaceArtifact = loadArtifact(path.join('..','contracts','artifacts','contracts','VerifierMarketplace.sol','VerifierMarketplace.json'));

  const escrowIface = new Interface(escrowArtifact.abi);
  const registryIface = new Interface(registryArtifact.abi);
  const marketplaceIface = new Interface(marketplaceArtifact.abi);

  const escrow = new Contract(addresses.escrow, escrowArtifact.abi, signer ?? provider);
  const registry = new Contract(addresses.registry, registryArtifact.abi, signer ?? provider);
  const marketplace = new Contract(addresses.marketplace, marketplaceArtifact.abi, signer ?? provider);

  logger.info({ network: env.NETWORK, chainId: env.CHAIN_ID, escrow: addresses.escrow, registry: addresses.registry, marketplace: addresses.marketplace }, 'Loaded contracts');

  return {
    provider,
    signer,
    addresses: { escrow: addresses.escrow, registry: addresses.registry, marketplace: addresses.marketplace },
    escrow,
    registry,
    marketplace,
    iface: { escrow: escrowIface, registry: registryIface, marketplace: marketplaceIface },
  };
}