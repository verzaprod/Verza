# Verza Frontend/Backend Architecture Audit and Midnight Integration Impact

## Frontend (Expo + expo-router)

- Framework: React Native via Expo (`package.json`), routing with `expo-router`.
- Entry: `package.json:4` points to `expo-router/entry`.
- App structure: `src/app` segmented routes:
  - Auth flow: `src/app/(auth)/...` (e.g., `register.tsx`, `sign-in.tsx`).
  - KYC flow: `src/app/(kyc)/...` (e.g., `verification-tracker.tsx`).
  - Tabs: `src/app/(tabs)/...` (e.g., `home.tsx`).
- Components: UI and domain components under `src/components/*`.
- State: `src/store/*` using `zustand`.
- API integration:
  - Client: `src/services/api/apiService.ts` interacts with backend endpoints (`/escrow`, `/verification`).
  - Config: `src/services/api/config.ts` defines base URL and mock toggles.

## Backend (Node + Express + Prisma)

- Server entry: `backend/src/index.ts:12-22` sets up Express with routes.
- Routes: `backend/src/routes/*` including `escrow.ts:21-138`, `health.ts`, `results.ts`, `uploads.ts`.
- Contracts integration: `backend/src/contracts/index.ts:46-74` loads EVM contract artifacts and creates `ethers` Contract instances.
- Worker: `backend/src/workers/chainWorker.ts:6-161` subscribes to EVM events and persists to Prisma.
- Config/env: `backend/src/config/env.ts:13-31` centralizes environment variables and path resolution (`contractsConfigPath`).

## Dependencies and Integration Points

- Frontend → Backend:
  - `apiService.initiateEscrow` calls `POST /escrow/initiate` with auth header (`src/services/api/apiService.ts:41-47`).
  - Status and results via `GET /escrow/status/:id` and `GET /verification/results/:id` (`src/services/api/apiService.ts:50-56`).
- Backend → Chain:
  - Uses Ethers provider and contract ABIs (`backend/src/contracts/index.ts:47-63`).
  - Emits DB updates on chain events (`backend/src/workers/chainWorker.ts`).

## Midnight (Compact) Integration Strategy

- Contracts (Compact) compiled in WSL at `/home/ekko/.local/bin/compact`.
- Generated TypeScript bindings consumed by backend to perform transactions and read state.
- New backend route namespace `/midnight/*` for health, compile, and registry actions.
- Optional midnight worker to watch Compact-managed state using generated APIs.

## Impact Assessment

- Frontend changes: none required initially. New optional screens can consume `/midnight/*` endpoints when enabled.
- Backend changes: add routes and worker; extend `env.ts` with Midnight-specific variables (Compact path, proof server compose file, worker toggle).
- Deployment: proof server must be available before generating ZK proofs; CI must include Compact compile step and artifact packaging.
- Security:
  - No private keys stored in code; backend continues using env-based secrets.
  - Middleware validates schema and normalizes request/response payloads.
- Versioning and migration:
  - Contract migrations tracked in `midnightcontract/migrations` with versioned scripts.
  - OpenAPI and Postman docs provided under `midnightcontract/docs`.

## Requirements Summary

- WSL Ubuntu with Compact compiler at `/home/ekko/.local/bin/compact`.
- Docker available to run Midnight proof server (testnet) per example projects.
- Node 18+ and TypeScript for generated bindings and backend integration.

