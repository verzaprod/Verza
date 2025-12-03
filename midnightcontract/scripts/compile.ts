import { spawn } from 'child_process';
import path from 'path';
import fs from 'fs';

function runCompactCompile(src: string, out: string, log: fs.WriteStream): Promise<void> {
  return new Promise((resolve, reject) => {
    const version = process.env.COMPACT_VERSION ?? '0.26.0';
    const args = [
      '/home/ekko/.local/bin/compact',
      'compile',
      `+${version}`,
      src,
      out,
    ];
    const proc = spawn('wsl', args);
    proc.stdout.on('data', (d) => log.write(d));
    proc.stderr.on('data', (d) => log.write(d));
    proc.on('exit', (code) => {
      if (code === 0) resolve();
      else reject(new Error(`compact compile failed (${code})`));
    });
    proc.on('error', (err) => reject(err));
  });
}

async function main() {
  const root = process.cwd();
  const logPath = path.join(root, 'compile.log');
  const log = fs.createWriteStream(logPath, { flags: 'w' });
  const counterSrcWin = path.join(root, 'contracts', 'counter.compact');
  const counterOutWin = path.join(root, 'interfaces', 'managed', 'counter');
  const registrySrcWin = path.join(root, 'contracts', 'registry.compact');
  const registryOutWin = path.join(root, 'interfaces', 'managed', 'registry');
  const escrowSrcWin = path.join(root, 'contracts', 'escrow.compact');
  const escrowOutWin = path.join(root, 'interfaces', 'managed', 'escrow');

  const toWslPath = (p: string) => {
    const m = p.match(/^([A-Za-z]):\\(.*)$/);
    if (!m) return p.replace(/\\/g, '/');
    const drive = m[1].toLowerCase();
    const rest = m[2].replace(/\\/g, '/');
    return `/mnt/${drive}/${rest}`;
  };

  [counterOutWin, registryOutWin, escrowOutWin].forEach((d) => {
    fs.mkdirSync(d, { recursive: true });
  });

  const counterSrc = toWslPath(counterSrcWin);
  const counterOut = toWslPath(counterOutWin);
  const registrySrc = toWslPath(registrySrcWin);
  const registryOut = toWslPath(registryOutWin);
  const escrowSrc = toWslPath(escrowSrcWin);
  const escrowOut = toWslPath(escrowOutWin);

  await runCompactCompile(counterSrc, counterOut, log);
  await runCompactCompile(registrySrc, registryOut, log);
  await runCompactCompile(escrowSrc, escrowOut, log);
  console.log('Compact compile completed');
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
