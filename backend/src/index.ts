import express from 'express';
import cors from 'cors';
import morgan from 'morgan';
import { env } from './config/env';
import { logger } from './logger';
import healthRouter from './routes/health';
import verifiersRouter from './routes/verifiers';
import escrowRouter from './routes/escrow';
import resultsRouter from './routes/results';
import uploadsRouter from './routes/uploads';
import midnightRouter from './routes/midnight';

const app = express();
app.use(cors());
app.use(express.json({ limit: '2mb' }));
app.use(express.urlencoded({ extended: true }));
app.use(morgan('dev'));

app.use('/health', healthRouter);
app.use('/verifiers', verifiersRouter);
app.use('/escrow', escrowRouter);
app.use('/verification', resultsRouter);
app.use('/uploads', uploadsRouter);
app.use('/midnight', midnightRouter);

app.use((err: any, _req: express.Request, res: express.Response, _next: express.NextFunction) => {
  logger.error({ err }, 'Unhandled error');
  res.status(500).json({ error: 'Internal Server Error' });
});

const port = env.PORT;
app.listen(port, () => {
  logger.info({ port }, 'Backend API listening');
});

// Optional: start worker
if (env.ENABLE_WORKER) {
  import('./workers/chainWorker')
    .then(({ startChainWorker }) => startChainWorker())
    .catch((e) => logger.error({ e }, 'Failed to start worker'));
}
