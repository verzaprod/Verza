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
const seed = process.env.MIDNIGHT_WALLET_SEED || '0'.repeat(64);
const log = arg('log', 'error');

const strip0x = (s) => (s && s.startsWith('0x')) ? s.slice(2) : s;

try {
  const walletMod = await import('@midnight-ntwrk/wallet');
  const zswap = await import('@midnight-ntwrk/zswap');
  const fetchMod = await import('node-fetch');
  const fetchFn = fetchMod.default || fetchMod;
  const networkId = zswap.NetworkId?.[network] ?? zswap.NetworkId?.TestNet;

  const wallet = await walletMod.WalletBuilder.build(indexer, ws, proof, rpc, seed, networkId, log);

  if (transcriptB64) {
    const jsonText = Buffer.from(transcriptB64, 'base64').toString('utf-8');
    const proofData = JSON.parse(jsonText);

    const unprovenOffer = new zswap.UnprovenOffer();
    const unproven = new zswap.UnprovenTransaction(unprovenOffer);

    const proveFn = async () => {
      const resp = await fetchFn(`${proof}/prove`, {
        method: 'POST',
        headers: { 'content-type': 'application/json' },
        body: JSON.stringify({ proofData, network }),
      });
      if (!resp.ok) throw new Error(`Prove failed: ${resp.status}`);
      const out = await resp.json();
      const txHex = out.txHex || out.transaction || out.tx || '';
      if (!txHex) throw new Error('Missing transaction in proof response');
      const raw = Buffer.from(strip0x(txHex), 'hex');
      return zswap.Transaction.deserialize(raw, networkId);
    };

    let tx;
    try {
      tx = await zswap.Transaction.fromUnproven(proveFn, unproven);
    } catch (e) {
      console.error(JSON.stringify({ error: e?.message ?? 'prove failed' }));
      process.exit(1);
    }

    let recipe;
    try {
      recipe = await wallet.balanceTransaction(tx, []);
    } catch (e) {
      console.error(JSON.stringify({ error: e?.message ?? 'balance failed' }));
      process.exit(1);
    }

    let finalTx = null;
    try {
      if (recipe?.type === 'NothingToProve') {
        finalTx = recipe.transaction;
      } else {
        finalTx = await wallet.proveTransaction(recipe);
      }
    } catch (e) {
      console.error(JSON.stringify({ error: e?.message ?? 'prove recipe failed' }));
      process.exit(1);
    }

    try {
      const identifier = await wallet.submitTransaction(finalTx);
      console.log(JSON.stringify({ status: 'ok', network, identifier }));
      process.exit(0);
    } catch (e) {
      console.error(JSON.stringify({ error: e?.message ?? 'submit failed' }));
      process.exit(1);
    }
  }

  console.log(JSON.stringify({ status: 'ok', network, walletCreated: !!wallet }));
  process.exit(0);
} catch (e) {
  console.error(JSON.stringify({ error: e?.message ?? 'submit failed' }));
  process.exit(1);
}
