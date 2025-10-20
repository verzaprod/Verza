import { NextFunction, Request, Response } from 'express';
import { env } from '../config/env';
import { jwtVerify, createRemoteJWKSet, JWTPayload } from 'jose';
import { URL } from 'url';

export interface AuthUser {
  id: string;
  email?: string;
  walletAddress?: string;
}

declare global {
  namespace Express {
    interface Request {
      user?: AuthUser;
    }
  }
}

const jwks = env.CLERK_JWKS_URL ? createRemoteJWKSet(new URL(env.CLERK_JWKS_URL)) : null;

export async function authMiddleware(req: Request, res: Response, next: NextFunction) {
  if (env.AUTH_BYPASS) {
    req.user = { id: 'dev-user' };
    return next();
  }

  try {
    const auth = req.headers.authorization || '';
    const token = auth.startsWith('Bearer ') ? auth.substring(7) : '';
    if (!token) return res.status(401).json({ error: 'Missing Bearer token' });

    if (!jwks) return res.status(500).json({ error: 'Auth not configured' });

    const { payload } = await jwtVerify(token, jwks);
    const user = mapClerkPayload(payload);
    req.user = user;
    next();
  } catch (e) {
    return res.status(401).json({ error: 'Invalid token' });
  }
}

function mapClerkPayload(payload: JWTPayload): AuthUser {
  const id = (payload.sub as string) || 'unknown';
  const email = (payload.email as string) || undefined;
  return { id, email };
}