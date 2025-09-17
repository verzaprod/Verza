-- name: CreateUser :one
INSERT INTO users (did, risk_score)
VALUES ($1, $2)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByDID :one
SELECT * FROM users
WHERE did = $1;

-- name: UpdateUserLastSeen :exec
UPDATE users
SET last_seen_at = NOW()
WHERE id = $1;

-- name: UpdateUserRiskScore :exec
UPDATE users
SET risk_score = $2
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CreateKYCJob :one
INSERT INTO kyc_jobs (user_id, status, score, liveness, doc_valid, result_json)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetKYCJob :one
SELECT * FROM kyc_jobs
WHERE id = $1;

-- name: GetKYCJobsByUser :many
SELECT * FROM kyc_jobs
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateKYCJobStatus :exec
UPDATE kyc_jobs
SET status = $2, score = $3, liveness = $4, doc_valid = $5, result_json = $6
WHERE id = $1;

-- name: ListKYCJobs :many
SELECT * FROM kyc_jobs
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateCredential :one
INSERT INTO credentials (subject_did, issuer_did, vc_hash, vc_jws, type, issued_at, expires_at, anchor_chain, anchor_tx)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetCredential :one
SELECT * FROM credentials
WHERE id = $1;

-- name: GetCredentialByHash :one
SELECT * FROM credentials
WHERE vc_hash = $1;

-- name: GetCredentialsBySubject :many
SELECT * FROM credentials
WHERE subject_did = $1 AND revoked = FALSE
ORDER BY issued_at DESC;

-- name: GetCredentialsByIssuer :many
SELECT * FROM credentials
WHERE issuer_did = $1
ORDER BY issued_at DESC
LIMIT $2 OFFSET $3;

-- name: GetCredentialsByType :many
SELECT * FROM credentials
WHERE $1 = ANY(type) AND revoked = FALSE
ORDER BY issued_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateCredentialAnchor :exec
UPDATE credentials
SET anchor_chain = $2, anchor_tx = $3
WHERE id = $1;

-- name: RevokeCredential :exec
UPDATE credentials
SET revoked = TRUE, revoked_at = NOW()
WHERE vc_hash = $1;

-- name: GetExpiredCredentials :many
SELECT * FROM credentials
WHERE expires_at IS NOT NULL AND expires_at < NOW() AND revoked = FALSE;

-- name: GetCredentialsForRevocation :many
SELECT c.* FROM credentials c
INNER JOIN revocations r ON c.vc_hash = r.vc_hash
WHERE r.revoked_at >= $1;

-- name: CreateRevocation :one
INSERT INTO revocations (vc_hash, reason)
VALUES ($1, $2)
RETURNING *;

-- name: GetRevocation :one
SELECT * FROM revocations
WHERE vc_hash = $1;

-- name: ListRevocations :many
SELECT * FROM revocations
ORDER BY revoked_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateAuditLog :exec
INSERT INTO audit_logs (actor, action, obj, meta)
VALUES ($1, $2, $3, $4);

-- name: GetAuditLogs :many
SELECT * FROM audit_logs
WHERE actor = $1
ORDER BY ts DESC
LIMIT $2 OFFSET $3;

-- name: GetAuditLogsByAction :many
SELECT * FROM audit_logs
WHERE action = $1
ORDER BY ts DESC
LIMIT $2 OFFSET $3;

-- name: GetAuditLogsByObject :many
SELECT * FROM audit_logs
WHERE obj = $1
ORDER BY ts DESC
LIMIT $2 OFFSET $3;

-- name: GetAuditLogsByTimeRange :many
SELECT * FROM audit_logs
WHERE ts >= $1 AND ts <= $2
ORDER BY ts DESC
LIMIT $3 OFFSET $4;

-- name: GetUserStats :one
SELECT 
    COUNT(*) as total_users,
    COUNT(*) FILTER (WHERE last_seen_at > NOW() - INTERVAL '24 hours') as active_users_24h,
    COUNT(*) FILTER (WHERE last_seen_at > NOW() - INTERVAL '7 days') as active_users_7d,
    AVG(risk_score) as avg_risk_score
FROM users;

-- name: GetKYCStats :one
SELECT 
    COUNT(*) as total_jobs,
    COUNT(*) FILTER (WHERE status = 'pending') as pending_jobs,
    COUNT(*) FILTER (WHERE status = 'processing') as processing_jobs,
    COUNT(*) FILTER (WHERE status = 'passed') as passed_jobs,
    COUNT(*) FILTER (WHERE status = 'failed') as failed_jobs,
    AVG(score) FILTER (WHERE score IS NOT NULL) as avg_score
FROM kyc_jobs;

-- name: GetCredentialStats :one
SELECT 
    COUNT(*) as total_credentials,
    COUNT(*) FILTER (WHERE revoked = FALSE) as active_credentials,
    COUNT(*) FILTER (WHERE revoked = TRUE) as revoked_credentials,
    COUNT(*) FILTER (WHERE expires_at IS NOT NULL AND expires_at < NOW()) as expired_credentials,
    COUNT(*) FILTER (WHERE anchor_tx IS NOT NULL) as anchored_credentials
FROM credentials;

-- name: SearchCredentials :many
SELECT * FROM credentials
WHERE 
    (subject_did ILIKE '%' || $1 || '%' OR issuer_did ILIKE '%' || $1 || '%')
    AND ($2::boolean IS NULL OR revoked = $2)
    AND ($3::text IS NULL OR $3 = ANY(type))
ORDER BY issued_at DESC
LIMIT $4 OFFSET $5;

-- name: GetCredentialsByDateRange :many
SELECT * FROM credentials
WHERE issued_at >= $1 AND issued_at <= $2
ORDER BY issued_at DESC
LIMIT $3 OFFSET $4;

-- name: BulkUpdateCredentialStatus :exec
UPDATE credentials
SET revoked = $1, revoked_at = CASE WHEN $1 THEN NOW() ELSE NULL END
WHERE id = ANY($2::uuid[]);