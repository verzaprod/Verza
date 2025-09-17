-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    did TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at TIMESTAMPTZ,
    risk_score NUMERIC(5,2) DEFAULT 0.0,
    CONSTRAINT users_did_check CHECK (length(did) > 0),
    CONSTRAINT users_risk_score_check CHECK (risk_score >= 0.0 AND risk_score <= 100.0)
);

-- KYC jobs table
CREATE TYPE kyc_status AS ENUM ('pending', 'processing', 'passed', 'failed', 'expired');

CREATE TABLE kyc_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status kyc_status NOT NULL DEFAULT 'pending',
    score NUMERIC(5,2),
    liveness BOOLEAN,
    doc_valid BOOLEAN,
    result_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT kyc_jobs_score_check CHECK (score IS NULL OR (score >= 0.0 AND score <= 100.0))
);

-- Credentials table
CREATE TABLE credentials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    subject_did TEXT NOT NULL,
    issuer_did TEXT NOT NULL,
    vc_hash BYTEA UNIQUE NOT NULL,
    vc_jws BYTEA NOT NULL,
    type TEXT[] NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    revoked_at TIMESTAMPTZ,
    anchor_chain TEXT,
    anchor_tx TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT credentials_subject_did_check CHECK (length(subject_did) > 0),
    CONSTRAINT credentials_issuer_did_check CHECK (length(issuer_did) > 0),
    CONSTRAINT credentials_type_check CHECK (array_length(type, 1) > 0),
    CONSTRAINT credentials_revoked_check CHECK (
        (revoked = FALSE AND revoked_at IS NULL) OR 
        (revoked = TRUE AND revoked_at IS NOT NULL)
    ),
    CONSTRAINT credentials_expires_check CHECK (expires_at IS NULL OR expires_at > issued_at)
);

-- Revocations table
CREATE TABLE revocations (
    vc_hash BYTEA PRIMARY KEY,
    reason TEXT NOT NULL,
    revoked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT revocations_reason_check CHECK (length(reason) > 0)
);

-- Audit logs table
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    actor TEXT NOT NULL,
    action TEXT NOT NULL,
    obj TEXT NOT NULL,
    meta JSONB,
    ts TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT audit_logs_actor_check CHECK (length(actor) > 0),
    CONSTRAINT audit_logs_action_check CHECK (length(action) > 0),
    CONSTRAINT audit_logs_obj_check CHECK (length(obj) > 0)
);

-- Indexes for performance
CREATE INDEX idx_users_did ON users(did);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_last_seen_at ON users(last_seen_at);
CREATE INDEX idx_users_risk_score ON users(risk_score);

CREATE INDEX idx_kyc_jobs_user_id ON kyc_jobs(user_id);
CREATE INDEX idx_kyc_jobs_status ON kyc_jobs(status);
CREATE INDEX idx_kyc_jobs_created_at ON kyc_jobs(created_at);
CREATE INDEX idx_kyc_jobs_updated_at ON kyc_jobs(updated_at);

CREATE INDEX idx_credentials_subject_did ON credentials(subject_did);
CREATE INDEX idx_credentials_issuer_did ON credentials(issuer_did);
CREATE INDEX idx_credentials_vc_hash ON credentials(vc_hash);
CREATE INDEX idx_credentials_type ON credentials USING GIN(type);
CREATE INDEX idx_credentials_issued_at ON credentials(issued_at);
CREATE INDEX idx_credentials_expires_at ON credentials(expires_at);
CREATE INDEX idx_credentials_revoked ON credentials(revoked);
CREATE INDEX idx_credentials_anchor_chain ON credentials(anchor_chain);
CREATE INDEX idx_credentials_anchor_tx ON credentials(anchor_tx);

CREATE INDEX idx_revocations_revoked_at ON revocations(revoked_at);

CREATE INDEX idx_audit_logs_actor ON audit_logs(actor);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_obj ON audit_logs(obj);
CREATE INDEX idx_audit_logs_ts ON audit_logs(ts);
CREATE INDEX idx_audit_logs_meta ON audit_logs USING GIN(meta);

-- Triggers for updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_kyc_jobs_updated_at BEFORE UPDATE ON kyc_jobs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_credentials_updated_at BEFORE UPDATE ON credentials
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to automatically create revocation record when credential is revoked
CREATE OR REPLACE FUNCTION handle_credential_revocation()
RETURNS TRIGGER AS $$
BEGIN
    -- If credential is being revoked
    IF NEW.revoked = TRUE AND OLD.revoked = FALSE THEN
        INSERT INTO revocations (vc_hash, reason, revoked_at)
        VALUES (NEW.vc_hash, 'Credential revoked', NEW.revoked_at)
        ON CONFLICT (vc_hash) DO NOTHING;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER credential_revocation_trigger AFTER UPDATE ON credentials
    FOR EACH ROW EXECUTE FUNCTION handle_credential_revocation();