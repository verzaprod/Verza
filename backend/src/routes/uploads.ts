import { Router } from 'express';
import { authMiddleware } from '../middleware/auth';
import multer from 'multer';
import fs from 'fs';
import path from 'path';
import { prisma } from '../db/client';

const router = Router();

router.post('/presign', authMiddleware, async (req, res) => {
  const { escrowId } = req.body as { escrowId?: string };
  if (!escrowId) return res.status(400).json({ error: 'escrowId required' });
  // For local dev, return direct upload endpoint
  res.json({
    docUploadUrl: `/uploads/verification/${escrowId}/documents`,
    selfieUploadUrl: `/uploads/verification/${escrowId}/documents`
  });
});

function ensureDir(dir: string) {
  if (!fs.existsSync(dir)) fs.mkdirSync(dir, { recursive: true });
}

const storage = multer.diskStorage({
  destination: function (req, file, cb) {
    const escrowId = req.params.escrowId;
    const dir = path.resolve(process.cwd(), 'uploads', 'verification', escrowId);
    ensureDir(dir);
    cb(null, dir);
  },
  filename: function (_req, file, cb) {
    const ts = Date.now();
    const ext = path.extname(file.originalname);
    const base = path.basename(file.originalname, ext);
    cb(null, `${base}-${ts}${ext}`);
  }
});

const upload = multer({ storage });

router.post('/verification/:escrowId/documents', authMiddleware, upload.fields([{ name: 'document', maxCount: 5 }, { name: 'selfie', maxCount: 1 }]), async (req, res) => {
  const escrowId = req.params.escrowId;
  const escrow = await prisma.escrow.findUnique({ where: { id: escrowId } });
  if (!escrow) return res.status(404).json({ error: 'Escrow not found' });

  const files = req.files as { [field: string]: Express.Multer.File[] };
  const docs = (files['document'] || []).map(f => f.path);
  const selfie = (files['selfie']?.[0]?.path) || null;

  await prisma.verification.upsert({
    where: { escrowId },
    update: { docUrls: docs, selfieUrl: selfie || undefined, status: 'received' },
    create: { escrowId, docUrls: docs, selfieUrl: selfie || undefined, status: 'received' }
  });

  res.json({ ok: true, received: { documents: docs.length, selfie: !!selfie } });
});

export default router;