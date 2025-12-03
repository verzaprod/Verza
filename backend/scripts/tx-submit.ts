import * as process from 'process';

const arg = (name: string, def?: string) => {
  const i = process.argv.indexOf(`--${name}`);
  return i >= 0 ? String(process.argv[i + 1] ?? '') : (def ?? '');
};

const indexer = arg('indexer');
const ws = arg('ws');
const rpc = arg('rpc');
const proof = arg('proof');
const network = arg('network', 'TestNet');

try {
  const walletMod = await import('@midnight-ntwrk/wallet');
  let networkId: any = network;
  try {
    const zswapMod = await import('@midnight-ntwrk/zswap');
    networkId = (zswapMod as any).NetworkId?.[network] ?? (zswapMod as any).NetworkId?.TestNet;
  } catch {}
  const wallet = await (walletMod as any).WalletBuilder.build(indexer, ws, proof, rpc, networkId, 'error');
  console.log(JSON.stringify({ status: 'ok', network, walletCreated: !!wallet }));
  process.exit(0);
} catch (e: any) {
  console.error(JSON.stringify({ error: e?.message ?? 'submit failed' }));
  process.exit(1);
}
