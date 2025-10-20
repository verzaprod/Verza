import { Router } from 'express';

const router = Router();
router.get('/', (_req, res) => {
  res.json({ ok: true, service: 'verza-backend', ts: new Date().toISOString() });
});

export default router;