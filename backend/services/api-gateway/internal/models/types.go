package models

import "time"

// Simplified VC/VP models for initial scaffolding and tests

type IssueVCRequest struct {
    SubjectDID string                 `json:"subjectDID" binding:"required"`
    Claims     map[string]interface{} `json:"claims" binding:"required"`
    Expiry     *time.Time             `json:"expiry,omitempty"`
}

type IssueVCResponse struct {
    VC       map[string]interface{} `json:"vc"`
    AnchorTx AnchorTx               `json:"anchorTx"`
}

type AnchorTx struct {
    ChainID string `json:"chainId"`
    TxHash  string `json:"txHash"`
}

type RevokeVCRequest struct {
    VCHash string `json:"vcHash" binding:"required"`
    Reason string `json:"reason,omitempty"`
}

type StatusResponse struct {
    Anchored  bool      `json:"anchored"`
    Revoked   bool      `json:"revoked"`
    IssuedAt  time.Time `json:"issuedAt"`
    RevokedAt *time.Time `json:"revokedAt,omitempty"`
    AnchorTx  *AnchorTx `json:"anchorTx,omitempty"`
    URI       string    `json:"uri"`
}

type VPVerifyRequest struct {
    VP      map[string]interface{} `json:"vp" binding:"required"`
    Options struct {
        Challenge string `json:"challenge" binding:"required"`
        Domain    string `json:"domain" binding:"required"`
    } `json:"options"`
}

type VPVerifyResponse struct {
    Valid   bool                   `json:"valid"`
    Details map[string]interface{} `json:"details"`
}

type DIDResolveRequest struct { Did string `json:"did" binding:"required"` }

type DIDDocument struct {
    ID                 string                   `json:"id"`
    VerificationMethod []map[string]interface{} `json:"verificationMethod"`
}

type KYCSubmitRequest struct {
    SelfieRef string `json:"selfieRef" binding:"required"`
    IDFrontRef string `json:"idFrontRef" binding:"required"`
    IDBackRef string `json:"idBackRef" binding:"required"`
    UserDID string `json:"userDID" binding:"required"`
}

type KYCSubmitResponse struct { JobID string `json:"jobId"` }

type KYCResultResponse struct {
    Score    float64               `json:"score"`
    Liveness bool                  `json:"liveness"`
    DocValid bool                  `json:"docValid"`
    OCR      map[string]string     `json:"ocr"`
}

type DIDChallengeRequest struct { DID string `json:"did" binding:"required"` }

type DIDChallengeResponse struct { Challenge string `json:"challenge"`; Domain string `json:"domain"` }

type DIDResponseRequest struct { DID string `json:"did"`; Signature string `json:"signature"` }

type DIDResponseToken struct { Token string `json:"token"` }