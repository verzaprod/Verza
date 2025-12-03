import { spawn } from 'child_process';
import path from 'path';
import fs from 'fs';

function runCompactCompile(src: string, out: string, log: fs.WriteStream): Promise<void> {
  return new Promise((resolve, reject) => {
    const args = [
      '/home/ekko/.local/bin/compact',
      'compile',
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
  const registrySrc = path.join(root, 'contracts', 'registry.compact');
  const accessSrc = path.join(root, 'contracts', 'access.compact');
  const registryOut = path.join(root, 'interfaces', 'managed', 'registry');
  const accessOut = path.join(root, 'interfaces', 'managed', 'access');

  await runCompactCompile(registrySrc, registryOut, log);
  await runCompactCompile(accessSrc, accessOut, log);
  console.log('Compact compile completed');
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
