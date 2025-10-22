# Render Deployment Configuration

This document outlines the configuration needed to deploy the Verza backend on Render.

## Files Created/Modified

1. **render.yaml** - Main Render configuration file
2. **backend/package.json** - Added production scripts
3. **backend/prisma/schema.prisma** - Updated to use PostgreSQL

## Required Environment Variables

Set these environment variables in your Render dashboard:

### Database
- `DATABASE_URL` - Use the Internal Database URL: `postgresql://verzadb_user:vmgxVkv8gfTWazRk7Pyf6tr2mBQ9Hwu9@dpg-d3shrapr0fns73c5kiq0-a/verzadb`

### Blockchain Configuration
- `RPC_URL` - Hedera testnet RPC endpoint (e.g., `https://testnet.hashio.io/api`)
- `CHAIN_ID` - Network chain ID (default: `296` for Hedera testnet)
- `SERVER_PRIVATE_KEY` - Private key for server wallet operations

### Authentication
- `CLERK_JWKS_URL` - Clerk JWKS URL for JWT verification
- `AUTH_BYPASS` - Set to `false` for production

### Application Configuration
- `NODE_ENV` - Set to `production`
- `PORT` - Set to `3001` (or Render's default)
- `CONTRACTS_CONFIG_PATH` - Path to contract configuration (`../contracts/contract-config.json`)
- `ESCROW_MODE` - Set to `custodial`
- `ENABLE_WORKER` - Set to `true` to enable blockchain event monitoring

## Deployment Steps

1. **Connect Repository**: Connect your GitHub repository to Render
2. **Create PostgreSQL Database**: 
   - Name: `verza-db`
   - Plan: Starter (or higher)
   - Region: Oregon (or preferred)
3. **Create Web Service**:
   - Use the `render.yaml` configuration
   - Render will automatically detect and use this file
4. **Set Environment Variables**: Configure all required environment variables in Render dashboard
5. **Deploy**: Render will automatically build and deploy your application

## Build Process

The build process defined in `render.yaml`:
1. `cd backend && npm ci` - Install dependencies
2. `npm run build` - Compile TypeScript
3. `npm run prisma:generate` - Generate Prisma client
4. `npm run prisma:deploy` - Run database migrations

## Health Check

The service includes a health check endpoint at `/health`. Make sure your backend implements this endpoint.

## Database Migrations

Database migrations will run automatically during deployment using `prisma migrate deploy`.

## Troubleshooting

- Ensure all environment variables are set correctly
- Check that the contracts configuration file is accessible
- Verify database connection string is correct
- Monitor build logs for any compilation errors