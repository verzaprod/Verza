import path from 'path';

function arg(name: string): string | undefined {
  const i = process.argv.indexOf(name);
  return i !== -1 ? process.argv[i + 1] : undefined;
}

async function main() {
  const key = arg('--key');
  const value = arg('--value');
  if (!key || !value) throw new Error('Missing --key/--value');
  const modPath = path.join(process.cwd(), 'interfaces', 'managed', 'registry', 'index.ts');
  const mod = await import(modPath);
  const candidates = [
    (mod as any).setRecord,
    (mod as any).registry?.setRecord,
    (mod as any).circuits?.setRecord,
  ];
  const fn = candidates.find(Boolean);
  if (!fn) throw new Error('setRecord binding not found');
  await fn(key, value);
  console.log(JSON.stringify({ status: 'ok' }));
}

main().catch((e) => {
  console.error(JSON.stringify({ error: e.message }));
  process.exit(1);
});

