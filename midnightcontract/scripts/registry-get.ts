import path from 'path';

function arg(name: string): string | undefined {
  const i = process.argv.indexOf(name);
  return i !== -1 ? process.argv[i + 1] : undefined;
}

async function main() {
  const key = arg('--key');
  if (!key) throw new Error('Missing --key');
  const modPath = path.join(process.cwd(), 'interfaces', 'managed', 'registry', 'index.ts');
  const mod = await import(modPath);
  const candidates = [
    (mod as any).getRecord,
    (mod as any).registry?.getRecord,
    (mod as any).circuits?.getRecord,
  ];
  const fn = candidates.find(Boolean);
  if (!fn) throw new Error('getRecord binding not found');
  const res = await fn(key);
  console.log(JSON.stringify({ status: 'ok', value: res }));
}

main().catch((e) => {
  console.error(JSON.stringify({ error: e.message }));
  process.exit(1);
});

