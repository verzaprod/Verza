import path from 'path';

function arg(name: string): string | undefined {
  const i = process.argv.indexOf(name);
  return i !== -1 ? process.argv[i + 1] : undefined;
}

async function main() {
  const key = arg('--key');
  if (!key) throw new Error('Missing --key');
  try {
    const cjsPath = path.join(process.cwd(), 'interfaces', 'managed', 'registry', 'contract', 'index.cjs');
    const mod = await import(cjsPath);
    const rt = await import('@midnight-ntwrk/compact-runtime');
    const hexToBytes = (h: string, len: number) => {
      const s = h.startsWith('0x') ? h.slice(2) : h;
      const b = Buffer.from(s, 'hex');
      if (b.length !== len) throw new Error(`Invalid length for hex: expected ${len}, got ${b.length}`);
      return new Uint8Array(b);
    };
    const contract = new (mod as any).Contract({});
    const init = contract.initialState({ initialPrivateState: {}, initialZswapLocalState: {} });
    const ctx = {
      originalState: init.currentContractState,
      currentPrivateState: init.currentPrivateState,
      currentZswapLocalState: init.currentZswapLocalState,
      transactionContext: new (rt as any).QueryContext(init.currentContractState.data, (rt as any).dummyContractAddress()),
    };
    const k = hexToBytes(key, 32);
    const out = contract.circuits.getRecord(ctx, k);
    const value = out.result?.[0] as Uint8Array;
    const hex = '0x' + Buffer.from(value).toString('hex');
    console.log(JSON.stringify({ status: 'ok', value: hex, via: 'compact' }));
    return;
  } catch {}
  const file = path.join(process.cwd(), 'interfaces', 'managed', 'local-registry.json');
  const store = await import('fs');
  let obj: Record<string, string> = {};
  try { obj = JSON.parse(store.readFileSync(file, 'utf-8')); } catch {}
  console.log(JSON.stringify({ status: 'ok', value: obj[key] ?? null, via: 'local' }));
}

main().catch((e) => {
  console.error(JSON.stringify({ error: e.message }));
  process.exit(1);
});
