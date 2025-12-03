const arg = (name, def) => {
  const i = process.argv.indexOf(`--${name}`);
  return i >= 0 ? String(process.argv[i + 1] ?? '') : (def ?? '');
};

const indexer = arg('indexer');
const ws = arg('ws');
const rpc = arg('rpc');
const proof = arg('proof');
const network = arg('network', 'TestNet');
const transcriptB64 = arg('transcript');

try {
  const walletMod = await import('@midnight-ntwrk/wallet');
  let networkId = network;
  try {
    const zswapMod = await import('@midnight-ntwrk/zswap');
    networkId = zswapMod.NetworkId?.[network] ?? zswapMod.NetworkId?.TestNet;
  } catch {}
  const wallet = await walletMod.WalletBuilder.build(indexer, ws, proof, rpc, networkId, 'error');
  if (transcriptB64) {
    const transcriptJson = Buffer.from(transcriptB64, 'base64').toString('utf-8');
    const transcript = JSON.parse(transcriptJson);
    let txId = null;
    let receipt = null;
    try {
      if (wallet.submitTranscript) {
        txId = await wallet.submitTranscript(transcript);
      } else if (wallet.submit) {
        receipt = await wallet.submit(transcript);
      }
      if (wallet.awaitReceipt && txId) {
        receipt = await wallet.awaitReceipt(txId);
      }
    } catch (e) {
      console.error(JSON.stringify({ error: e?.message ?? 'submit failed' }));
      process.exit(1);
    }
    console.log(JSON.stringify({ status: 'ok', network, txId, receipt }));
    process.exit(0);
  }
  console.log(JSON.stringify({ status: 'ok', network, walletCreated: !!wallet }));
  process.exit(0);
} catch (e) {
  console.error(JSON.stringify({ error: e?.message ?? 'submit failed' }));
  process.exit(1);
}
