import { Router } from 'express';
import path from 'path';
import { spawn } from 'child_process';
import { env } from '../config/env';
import { authMiddleware } from '../middleware/auth';
import { z } from 'zod';

const router = Router();

router.get('/health', async (_req, res) => {
  const proc = spawn('wsl', [env.MIDNIGHT_COMPACT_PATH, '--version']);
  let out = '';
  proc.stdout.on('data', (d) => (out += d.toString()));
  proc.on('exit', () => {
    res.json({ compactVersion: out.trim() || null });
  });
  proc.on('error', () => res.status(500).json({ error: 'compact not available' }));
});

router.post('/compile', authMiddleware, async (_req, res) => {
  const cwd = path.resolve(process.cwd(), 'midnightcontract');
  const npmCmd = process.platform === 'win32' ? 'npm.cmd' : 'npm';
  const proc = spawn(npmCmd, ['run', 'compile'], { cwd });
  proc.on('exit', (code) => {
    if (code === 0) res.json({ status: 'ok' });
    else res.status(500).json({ status: 'failed' });
  });
  proc.on('error', (e) => res.status(500).json({ error: e.message }));
});

const setSchema = z.object({ key: z.string(), value: z.string() });

router.post('/registry/set', authMiddleware, async (req, res) => {
  const parse = setSchema.safeParse(req.body);
  if (!parse.success) return res.status(400).json({ error: 'Missing key/value' });
  const { key, value } = parse.data;
  const cwd = path.resolve(process.cwd(), 'midnightcontract');
  const args = ['--loader','ts-node/esm','./scripts/registry-set.ts','--key', key, '--value', value];
  const proc = spawn('node', args, { cwd });
  let out = '';
  let err = '';
  proc.stdout.on('data', (d) => (out += d.toString()));
  proc.stderr.on('data', (d) => (err += d.toString()));
  proc.on('exit', (code) => {
    if (code === 0) {
      try { return res.json(JSON.parse(out)); } catch { return res.json({ status: 'ok' }); }
    }
    try { return res.status(500).json(JSON.parse(err)); } catch { return res.status(500).json({ error: 'registry set failed' }); }
  });
});

router.get('/registry/get', authMiddleware, async (req, res) => {
  const key = String(req.query.key || '');
  if (!key) return res.status(400).json({ error: 'Missing key' });
  const cwd = path.resolve(process.cwd(), 'midnightcontract');
  const args = ['--loader','ts-node/esm','./scripts/registry-get.ts','--key', key];
  const proc = spawn('node', args, { cwd });
  let out = '';
  let err = '';
  proc.stdout.on('data', (d) => (out += d.toString()));
  proc.stderr.on('data', (d) => (err += d.toString()));
  proc.on('exit', (code) => {
    if (code === 0) {
      try { return res.json(JSON.parse(out)); } catch { return res.json({ status: 'ok' }); }
    }
    try { return res.status(500).json(JSON.parse(err)); } catch { return res.status(500).json({ error: 'registry get failed' }); }
  });
});

export default router;

