import { spawnSync } from 'child_process';

export function checkCompactVersion(): string | null {
  const res = spawnSync('wsl', ['/home/ekko/.local/bin/compact', '--version'], { encoding: 'utf-8' });
  if (res.status === 0) return res.stdout.trim();
  return null;
}
