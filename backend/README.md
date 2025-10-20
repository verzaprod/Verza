# Verza Backend (MVP)

TypeScript/Express backend providing APIs for escrow initiation, verifiers catalog, uploads, and verification results, integrated with Hedera smart contracts.

## Quick Start

1. Copy env and adjust

```
cp .env.example .env
```

2. Install and generate DB

```
npm install
npm run prisma:generate
npm run prisma:push
```

3. Run in dev

```
npm run dev
```

Server listens on `http://localhost:3001` by default.

## Endpoints

- `POST /escrow/initiate` — Initiate escrow
- `GET /escrow/status/:escrowId` — Escrow status
- `GET /verifiers` and `GET /verifiers/:id` — Verifiers catalog
- `POST /uploads/presign` — Presign upload (local dev)
- `POST /uploads/verification/:escrowId/documents` — Upload documents/selfie (multipart)
- `GET /verification/results/:escrowId` — Verification results

All user endpoints use Bearer auth. In dev you can set `AUTH_BYPASS=true`.

## Contracts

Addresses are read from `../contracts/contract-config.json` (Hedera testnet). ABIs are loaded from the Hardhat `artifacts` directory. Configure RPC via `RPC_URL`.

- Escrow: calls `createEscrow(bytes32,address)` in non-custodial flow by returning unsigned tx payload for the mobile app to sign.
- VCRegistry: used later when issuing credentials after successful verification.
- VerifierMarketplace: used to resolve active verifier info and fee.

## Next Steps

- Add worker to subscribe to `EscrowCreated`, call `lockFunds`, and mirror on-chain status to DB.
- Add verification provider integration and credential issuance via `VCRegistry.issueCredential`.
- Implement ratings endpoints.