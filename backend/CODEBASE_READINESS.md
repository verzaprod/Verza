# Verza Codebase Review & Readiness Plan

This document summarizes the current state of the backend and smart contracts, highlights what’s needed to make the frontend fully functional, and provides a practical deployment plan to Hedera Testnet.

## Current Status Overview

- Backend API (TypeScript/Express)
  - Implemented: `POST /escrow/initiate`, `GET /escrow/status/:escrowId`, `GET /verification/results/:escrowId`
  - Implemented: Verifiers catalog `GET /verifiers`, `GET /verifiers/:id` with on-chain metadata resolution
  - Implemented: Uploads `POST /uploads/presign`, `POST /uploads/verification/:escrowId/documents` (local storage)
  - Auth: Clerk JWT verification middleware with `AUTH_BYPASS` option for dev
  - Chain worker: Subscribes to `EscrowContract` events; mirrors escrow lifecycle in DB
- Smart Contracts (Upgradeable: VCRegistry, VerifierMarketplace, EscrowContract)
  - Hedera networks configured in Hardhat (`hederaTestnet`, chainId `296`)
  - Deployment script deploys proxies and writes addresses to `contracts/contract-config.json`
  - ABI artifacts available; backend loads them from `contracts/artifacts`
- Frontend
  - Uses real backend (`USE_MOCK=false`), calling `/escrow/initiate`, `/escrow/status/:escrowId`, `/verification/results/:escrowId`
  - Verifier details currently read from mock; intended to be wired to `/verifiers/:id`

## What’s Needed (Backend)

- Verification pipeline integration
  - Accept uploaded docs/selfie and route to a provider (e.g., Sumsub/Persona/Onfido) or internal verifier
  - Persist verification result and update `Verification` and `Escrow` status
  - On success, trigger credential issuance via `VCRegistry` and store `tokenId`, `tokenURI` in `Credential`
- VCRegistry event indexing (optional but recommended)
  - Extend chain worker to subscribe to VC issuance events to reconcile DB state for `Credential`
- Wallet/DID mapping
  - Persist user wallet (`wallet_address`) and optional DID; required for VC issuance and results display
- Ratings (optional)
  - Endpoints for `POST /verifiers/:id/ratings`, `GET /verifiers/:id/ratings` to support UX enhancements
- Storage provider abstraction
  - Replace local disk with S3/GCS and presigned URLs for production
- Security & Ops
  - Rate limiting, audit logging, PII redaction, consistent error handling

## What’s Needed (Smart Contracts)

- Confirm release/refund/dispute paths are callable by appropriate roles
  - `FundsReleased`, `RefundIssued`, `EscrowCancelled` are emitted and currently indexed by worker
  - Ensure verifier/admin gating is correct via `AccessControl`
- VCRegistry integration flow
  - Backend needs to call `issueCredential` on success; ensure issuer is registered and token URI schema matches frontend display
- VerifierMarketplace metadata
  - `getVerifierDetails` used when available; fallback to `getVerifier` works for older versions

## Frontend Wiring Gaps

- Verifier details
  - Replace mock data in `useVerifierDetails.tsx` with call to `GET /verifiers/:id`
- Auth token
  - Ensure frontend passes a valid Clerk Bearer token (dev can set `AUTH_BYPASS=true` on backend)
- Non-custodial flow wallet address
  - If `ESCROW_MODE=noncustodial`, include `wallet_address` in `POST /escrow/initiate` body

## Hedera Deployment Plan

### Prerequisites

- Contracts `.env`
  - `PRIVATE_KEY=0x...` (deployer EVM key)
  - `EVM_ADDRESS=0x...` (derived from the private key)
  - `ACCOUNT_ID=0.0.xxxxxxx` (Hedera account ID)
  - `HEDERA_RPC_URL=https://testnet.hashio.io/api`
- Backend `.env`
  - `RPC_URL=https://testnet.hashio.io/api`
  - `CHAIN_ID=296`
  - `NETWORK=hederaTestnet`
  - `CONTRACTS_CONFIG_PATH=../contracts/contract-config.json`
  - `AUTH_BYPASS=true` (for local dev) or `false` (production)
  - `ESCROW_MODE=custodial` with `SERVER_PRIVATE_KEY=0x...` OR `ESCROW_MODE=noncustodial`
  - `ENABLE_WORKER=true` to index events

### Deploy Contracts to Hedera Testnet

- Compile and deploy
  - `cd contracts`
  - `npm install`
  - `npx hardhat compile`
  - `npx hardhat run scripts/deploy.js --network hederaTestnet`
- Verify deployment artifacts
  - Contracts addresses written to `contracts/contract-config.json`
  - Detailed deployment info saved to `contracts/deployments/hederaTestnet.json`

### Start Backend Against Hedera

- Initialize DB
  - `cd backend`
  - `npm install`
  - `npm run prisma:generate`
  - `npm run prisma:push`
- Run API (dev)
  - `npm run dev`
- Run API (production)
  - `npm run build`
  - `npm start`

### Functional Test Flow

- Initiate escrow
  - `POST /escrow/initiate` (custodial or non-custodial)
- Upload documents/selfie
  - `POST /uploads/presign` → upload via returned URLs
- Poll status
  - `GET /escrow/status/:escrowId`
- Complete verification and issue VC
  - After verification, backend calls `VCRegistry.issueCredential`
  - `GET /verification/results/:escrowId` returns `credential` object on success

## Recommendations & Next Steps

- Extend chain worker to index VC issuance events and reconcile `Credential`
- Add a simple verifier adapter to emulate verification success/failure for end-to-end testing
- Wire frontend verifier details to `/verifiers/:id` and ensure Clerk auth for all calls
- Move uploads to S3/GCS with presigned URLs for production
- Add rate limiting and audit logging for user endpoints

## Checklist

- Contracts deployable to Hedera Testnet via Hardhat
- Backend reading `contract-config.json` and connecting to Hedera RPC
- Chain worker enabled and indexing `EscrowContract` events
- Verification pipeline integrated and VC issuance invoked
- Frontend calling backend endpoints with proper auth