import { Router } from 'express';
import path from 'path';
import { spawn } from 'child_process';
import { createRequire } from 'module';
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

router.get('/health/zk', async (_req, res) => {
  try {
    const url = new URL('/health', env.PROOF_SERVER_URL).toString();
    const resp = await fetch(url);
    if (!resp.ok) return res.status(502).json({ error: 'proof server unhealthy' });
    const json = await resp.json().catch(() => ({ status: 'ok' }));
    return res.json({ proofServer: 'ok', details: json });
  } catch (e: any) {
    return res.status(500).json({ error: 'proof server not reachable', message: e?.message });
  }
});

const escrowCreateSchema = z.object({
  requestId: z.string(),
  verifier: z.string(),
  amount: z.string(),
});

router.post('/escrow/create', authMiddleware, async (req, res) => {
  const parse = escrowCreateSchema.safeParse(req.body);
  if (!parse.success) return res.status(400).json({ error: 'Missing fields' });
  const { requestId, verifier, amount } = parse.data;
  const cwd = path.resolve(process.cwd(), '..', 'midnightcontract');
  try {
    const requireFromMc = createRequire(path.join(cwd, 'package.json'));
    const mod = requireFromMc(path.join(cwd, 'interfaces', 'managed', 'escrow', 'contract', 'index.cjs'));
    const rt = requireFromMc('@midnight-ntwrk/compact-runtime');
    const hexToBytes = (h: string, len: number) => {
      const s = h.startsWith('0x') ? h.slice(2) : h;
      const b = Buffer.from(s, 'hex');
      if (b.length !== len) throw new Error(`Invalid length for hex: expected ${len}, got ${b.length}`);
      return new Uint8Array(b);
    };
    const contract = new mod.Contract({});
    const init = contract.initialState({ initialPrivateState: {}, initialZswapLocalState: {} });
    const ctx = {
      originalState: init.currentContractState,
      currentPrivateState: init.currentPrivateState,
      currentZswapLocalState: init.currentZswapLocalState,
      transactionContext: new rt.QueryContext(init.currentContractState.data, rt.dummyContractAddress()),
    };
    contract.circuits.createEscrow(ctx, hexToBytes(requestId,32), hexToBytes(verifier,32), hexToBytes(amount,8));
    const led = mod.ledger(ctx.transactionContext.state);
    return res.json({ status: 'ok', ledger: {
      created: Number(led.created),
      lastRequest: '0x' + Buffer.from(led.lastRequest).toString('hex'),
      lastVerifier: '0x' + Buffer.from(led.lastVerifier).toString('hex'),
      lastAmount: '0x' + Buffer.from(led.lastAmount).toString('hex'),
    }, via: 'compact' });
  } catch (e: any) {
    return res.status(500).json({ error: e?.message ?? 'escrow create failed' });
  }
});

router.post('/escrow/lock', authMiddleware, async (req, res) => {
  const requestId = String(req.body?.requestId || '');
  if (!requestId) return res.status(400).json({ error: 'Missing requestId' });
  const cwd = path.resolve(process.cwd(), '..', 'midnightcontract');
  try {
    const requireFromMc = createRequire(path.join(cwd, 'package.json'));
    const mod = requireFromMc(path.join(cwd, 'interfaces', 'managed', 'escrow', 'contract', 'index.cjs'));
    const rt = requireFromMc('@midnight-ntwrk/compact-runtime');
    const hexToBytes = (h: string, len: number) => new Uint8Array(Buffer.from(h.startsWith('0x')?h.slice(2):h,'hex'));
    const contract = new mod.Contract({});
    const init = contract.initialState({ initialPrivateState: {}, initialZswapLocalState: {} });
    const ctx = {
      originalState: init.currentContractState,
      currentPrivateState: init.currentPrivateState,
      currentZswapLocalState: init.currentZswapLocalState,
      transactionContext: new rt.QueryContext(init.currentContractState.data, rt.dummyContractAddress()),
    };
    contract.circuits.markLocked(ctx, hexToBytes(requestId,32));
    const led = mod.ledger(ctx.transactionContext.state);
    return res.json({ status: 'ok', ledger: { locked: Number(led.locked) }, via: 'compact' });
  } catch (e: any) {
    return res.status(500).json({ error: e?.message ?? 'escrow lock failed' });
  }
});

router.post('/escrow/release', authMiddleware, async (req, res) => {
  const requestId = String(req.body?.requestId || '');
  const verifier = String(req.body?.verifier || '');
  if (!requestId || !verifier) return res.status(400).json({ error: 'Missing fields' });
  const cwd = path.resolve(process.cwd(), '..', 'midnightcontract');
  try {
    const requireFromMc = createRequire(path.join(cwd, 'package.json'));
    const mod = requireFromMc(path.join(cwd, 'interfaces', 'managed', 'escrow', 'contract', 'index.cjs'));
    const rt = requireFromMc('@midnight-ntwrk/compact-runtime');
    const toBytes = (h: string, len: number) => new Uint8Array(Buffer.from(h.startsWith('0x')?h.slice(2):h,'hex'));
    const contract = new mod.Contract({});
    const init = contract.initialState({ initialPrivateState: {}, initialZswapLocalState: {} });
    const ctx = {
      originalState: init.currentContractState,
      currentPrivateState: init.currentPrivateState,
      currentZswapLocalState: init.currentZswapLocalState,
      transactionContext: new rt.QueryContext(init.currentContractState.data, rt.dummyContractAddress()),
    };
    contract.circuits.release(ctx, toBytes(requestId,32), toBytes(verifier,32));
    const led = mod.ledger(ctx.transactionContext.state);
    return res.json({ status: 'ok', ledger: { released: Number(led.released) }, via: 'compact' });
  } catch (e: any) {
    return res.status(500).json({ error: e?.message ?? 'escrow release failed' });
  }
});

router.post('/escrow/refund', authMiddleware, async (req, res) => {
  const requestId = String(req.body?.requestId || '');
  if (!requestId) return res.status(400).json({ error: 'Missing requestId' });
  const cwd = path.resolve(process.cwd(), '..', 'midnightcontract');
  try {
    const requireFromMc = createRequire(path.join(cwd, 'package.json'));
    const mod = requireFromMc(path.join(cwd, 'interfaces', 'managed', 'escrow', 'contract', 'index.cjs'));
    const rt = requireFromMc('@midnight-ntwrk/compact-runtime');
    const toBytes = (h: string, len: number) => new Uint8Array(Buffer.from(h.startsWith('0x')?h.slice(2):h,'hex'));
    const contract = new mod.Contract({});
    const init = contract.initialState({ initialPrivateState: {}, initialZswapLocalState: {} });
    const ctx = {
      originalState: init.currentContractState,
      currentPrivateState: init.currentPrivateState,
      currentZswapLocalState: init.currentZswapLocalState,
      transactionContext: new rt.QueryContext(init.currentContractState.data, rt.dummyContractAddress()),
    };
    contract.circuits.refund(ctx, toBytes(requestId,32));
    const led = mod.ledger(ctx.transactionContext.state);
    return res.json({ status: 'ok', ledger: { refunded: Number(led.refunded) }, via: 'compact' });
  } catch (e: any) {
    return res.status(500).json({ error: e?.message ?? 'escrow refund failed' });
  }
});

router.post('/compile', authMiddleware, async (_req, res) => {
  const cwd = path.resolve(process.cwd(), '..', 'midnightcontract');
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
  const cwd = path.resolve(process.cwd(), '..', 'midnightcontract');
  try {
    const requireFromMc = createRequire(path.join(cwd, 'package.json'));
    const mod = requireFromMc(path.join(cwd, 'interfaces', 'managed', 'registry', 'contract', 'index.cjs'));
    const rt = requireFromMc('@midnight-ntwrk/compact-runtime');
    const hexToBytes = (h: string, len: number) => {
      const s = h.startsWith('0x') ? h.slice(2) : h;
      const b = Buffer.from(s, 'hex');
      if (b.length !== len) throw new Error(`Invalid length for hex: expected ${len}, got ${b.length}`);
      return new Uint8Array(b);
    };
    const contract = new mod.Contract({});
    const init = contract.initialState({ initialPrivateState: {}, initialZswapLocalState: {} });
    const ctx = {
      originalState: init.currentContractState,
      currentPrivateState: init.currentPrivateState,
      currentZswapLocalState: init.currentZswapLocalState,
      transactionContext: new rt.QueryContext(init.currentContractState.data, rt.dummyContractAddress()),
    };
    const k = hexToBytes(key, 32);
    const v = hexToBytes(value, 64);
    contract.circuits.setRecord(ctx, k, v);
    return res.json({ status: 'ok', via: 'compact' });
  } catch (e) {
    const args = ['--loader','ts-node/esm','./scripts/registry-set.ts','--key', key, '--value', value];
    const proc = spawn(process.execPath, args, { cwd });
    let out = '';
    let err = '';
    proc.stdout.on('data', (d) => (out += d.toString()));
    proc.stderr.on('data', (d) => (err += d.toString()));
    proc.on('error', () => {
      try {
        const fs = require('fs');
        const file = path.join(cwd, 'interfaces', 'managed', 'local-registry.json');
        fs.mkdirSync(path.dirname(file), { recursive: true });
        let obj = {} as Record<string,string>;
        try { obj = JSON.parse(fs.readFileSync(file, 'utf-8')); } catch {}
        obj[key] = value;
        fs.writeFileSync(file, JSON.stringify(obj, null, 2));
        return res.json({ status: 'ok', via: 'local' });
      } catch {
        return res.status(500).json({ error: 'registry set failed' });
      }
    });
    proc.on('exit', (code) => {
      if (code === 0) {
        try { return res.json(JSON.parse(out)); } catch { return res.json({ status: 'ok' }); }
      }
      try {
        const fs = require('fs');
        const file = path.join(cwd, 'interfaces', 'managed', 'local-registry.json');
        let obj = {} as Record<string,string>;
        try { obj = JSON.parse(fs.readFileSync(file, 'utf-8')); } catch {}
        obj[key] = value;
        fs.writeFileSync(file, JSON.stringify(obj, null, 2));
        return res.json({ status: 'ok', via: 'local' });
      } catch (e) {
        return res.status(500).json({ error: 'registry set failed' });
      }
    });
  }
});

router.get('/registry/get', authMiddleware, async (req, res) => {
  const key = String(req.query.key || '');
  if (!key) return res.status(400).json({ error: 'Missing key' });
  const cwd = path.resolve(process.cwd(), 'midnightcontract');
  try {
    const requireFromMc = createRequire(path.join(cwd, 'package.json'));
    const mod = requireFromMc(path.join(cwd, 'interfaces', 'managed', 'registry', 'contract', 'index.cjs'));
    const rt = requireFromMc('@midnight-ntwrk/compact-runtime');
    const hexToBytes = (h: string, len: number) => {
      const s = h.startsWith('0x') ? h.slice(2) : h;
      const b = Buffer.from(s, 'hex');
      if (b.length !== len) throw new Error(`Invalid length for hex: expected ${len}, got ${b.length}`);
      return new Uint8Array(b);
    };
    const contract = new mod.Contract({});
    const init = contract.initialState({ initialPrivateState: {}, initialZswapLocalState: {} });
    const ctx = {
      originalState: init.currentContractState,
      currentPrivateState: init.currentPrivateState,
      currentZswapLocalState: init.currentZswapLocalState,
      transactionContext: new rt.QueryContext(init.currentContractState.data, rt.dummyContractAddress()),
    };
    const k = hexToBytes(key, 32);
    const out = contract.circuits.getRecord(ctx, k);
    const value = out.result?.[0] as Uint8Array;
    const hex = '0x' + Buffer.from(value).toString('hex');
    return res.json({ status: 'ok', value: hex, via: 'compact' });
  } catch (e) {
    const args = ['--loader','ts-node/esm','./scripts/registry-get.ts','--key', key];
    const proc = spawn(process.execPath, args, { cwd });
    let out = '';
    let err = '';
    proc.stdout.on('data', (d) => (out += d.toString()));
    proc.stderr.on('data', (d) => (err += d.toString()));
    proc.on('error', () => {
      try {
        const fs = require('fs');
        const file = path.join(cwd, 'interfaces', 'managed', 'local-registry.json');
        fs.mkdirSync(path.dirname(file), { recursive: true });
        let obj = {} as Record<string,string>;
        try { obj = JSON.parse(fs.readFileSync(file, 'utf-8')); } catch {}
        return res.json({ status: 'ok', value: obj[key] ?? null, via: 'local' });
      } catch {
        return res.status(500).json({ error: 'registry get failed' });
      }
    });
    proc.on('exit', (code) => {
      if (code === 0) {
        try { return res.json(JSON.parse(out)); } catch { return res.json({ status: 'ok' }); }
      }
      try {
        const fs = require('fs');
        const file = path.join(cwd, 'interfaces', 'managed', 'local-registry.json');
        let obj = {} as Record<string,string>;
        try { obj = JSON.parse(fs.readFileSync(file, 'utf-8')); } catch {}
        return res.json({ status: 'ok', value: obj[key] ?? null, via: 'local' });
      } catch (e) {
        return res.status(500).json({ error: 'registry get failed' });
      }
    });
  }
});

export default router;

