import path from 'path';

function arg(name: string): string | undefined {
  const i = process.argv.indexOf(name);
  return i !== -1 ? process.argv[i + 1] : undefined;
}

async function main() {
  const key = arg('--key');
  const value = arg('--value');
  if (!key || !value) throw new Error('Missing --key/--value');
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
    const v = hexToBytes(value, 64);
    contract.circuits.setRecord(ctx, k, v);
    console.log(JSON.stringify({ status: 'ok', via: 'compact' }));
    return;
  } catch {}
  // Fallback: local JSON registry
  const file = path.join(process.cwd(), 'interfaces', 'managed', 'local-registry.json');
  const store = await import('fs');
  let obj: Record<string, string> = {};
  try { obj = JSON.parse(store.readFileSync(file, 'utf-8')); } catch {}
  obj[key] = value;
  store.writeFileSync(file, JSON.stringify(obj, null, 2));
  console.log(JSON.stringify({ status: 'ok', via: 'local' }));
}

main().catch((e) => {
  console.error(JSON.stringify({ error: e.message }));
  process.exit(1);
});
