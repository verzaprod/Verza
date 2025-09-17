package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Onfido API configuration
type OnfidoConfig struct {
	APIKey  string
	BaseURL string
	Region  string // "EU", "US", "CA"
}

// KYC verification request
type KYCRequest struct {
	ApplicantID string `json:"applicant_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email,omitempty"`
	DOB         string `json:"dob,omitempty"` // YYYY-MM-DD format
	Address     struct {
		Street      string `json:"street,omitempty"`
		Town        string `json:"town,omitempty"`
		Postcode    string `json:"postcode,omitempty"`
		Country     string `json:"country,omitempty"`
		State       string `json:"state,omitempty"`
	} `json:"address,omitempty"`
}

// Document upload request
type DocumentUploadRequest struct {
	ApplicantID  string `json:"applicant_id"`
	DocumentType string `json:"type"` // "passport", "driving_licence", "national_identity_card"
	Side         string `json:"side,omitempty"` // "front", "back"
}

// Biometric verification request
type BiometricRequest struct {
	ApplicantID string `json:"applicant_id"`
	Variant     string `json:"variant"` // "standard", "video"
}

// Onfido API responses
type OnfidoApplicant struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	DOB       string    `json:"dob"`
	Address   struct {
		Street      string `json:"street"`
		Town        string `json:"town"`
		Postcode    string `json:"postcode"`
		Country     string `json:"country"`
		State       string `json:"state"`
	} `json:"address"`
}

type OnfidoDocument struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	FileName     string    `json:"file_name"`
	FileType     string    `json:"file_type"`
	FileSize     int       `json:"file_size"`
	Type         string    `json:"type"`
	Side         string    `json:"side"`
	ApplicantID  string    `json:"applicant_id"`
}

type OnfidoLivePhoto struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	FileName    string    `json:"file_name"`
	FileType    string    `json:"file_type"`
	FileSize    int       `json:"file_size"`
	ApplicantID string    `json:"applicant_id"`
}

type OnfidoCheck struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Status      string    `json:"status"` // "in_progress", "complete"
	Result      string    `json:"result"` // "clear", "consider", "unidentified"
	ApplicantID string    `json:"applicant_id"`
	ReportIDs   []string  `json:"report_ids"`
}

type OnfidoReport struct {
	ID         string                 `json:"id"`
	CreatedAt  time.Time              `json:"created_at"`
	Name       string                 `json:"name"`
	Status     string                 `json:"status"`
	Result     string                 `json:"result"`
	SubResult  string                 `json:"sub_result"`
	Properties map[string]interface{} `json:"properties"`
}

// ML Service
type MLService struct {
	logger        *zap.Logger
	onfidoConfig  OnfidoConfig
	httpClient    *http.Client
}

func NewMLService(logger *zap.Logger, config OnfidoConfig) *MLService {
	return &MLService{
		logger:       logger,
		onfidoConfig: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Create Onfido applicant
func (ml *MLService) CreateApplicant(ctx context.Context, req KYCRequest) (*OnfidoApplicant, error) {
	url := fmt.Sprintf("%s/applicants", ml.onfidoConfig.BaseURL)
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Token token="+ml.onfidoConfig.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := ml.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("onfido API error: %d - %s", resp.StatusCode, string(body))
	}

	var applicant OnfidoApplicant
	if err := json.NewDecoder(resp.Body).Decode(&applicant); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &applicant, nil
}

// Upload document to Onfido
func (ml *MLService) UploadDocument(ctx context.Context, applicantID string, docType string, side string, fileData []byte, fileName string) (*OnfidoDocument, error) {
	url := fmt.Sprintf("%s/documents", ml.onfidoConfig.BaseURL)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	writer.WriteField("applicant_id", applicantID)
	writer.WriteField("type", docType)
	if side != "" {
		writer.WriteField("side", side)
	}

	// Add file
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = part.Write(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to write file data: %w", err)
	}

	writer.Close()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Token token="+ml.onfidoConfig.APIKey)
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := ml.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("onfido API error: %d - %s", resp.StatusCode, string(body))
	}

	var document OnfidoDocument
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &document, nil
}

// Upload live photo for biometric verification
func (ml *MLService) UploadLivePhoto(ctx context.Context, applicantID string, fileData []byte, fileName string) (*OnfidoLivePhoto, error) {
	url := fmt.Sprintf("%s/live_photos", ml.onfidoConfig.BaseURL)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	writer.WriteField("applicant_id", applicantID)

	// Add file
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = part.Write(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to write file data: %w", err)
	}

	writer.Close()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Token token="+ml.onfidoConfig.APIKey)
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := ml.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("onfido API error: %d - %s", resp.StatusCode, string(body))
	}

	var livePhoto OnfidoLivePhoto
	if err := json.NewDecoder(resp.Body).Decode(&livePhoto); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &livePhoto, nil
}

// Create verification check
func (ml *MLService) CreateCheck(ctx context.Context, applicantID string, reportNames []string) (*OnfidoCheck, error) {
	url := fmt.Sprintf("%s/checks", ml.onfidoConfig.BaseURL)

	reqBody := map[string]interface{}{
		"applicant_id": applicantID,
		"report_names": reportNames,
	}

	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Token token="+ml.onfidoConfig.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := ml.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("onfido API error: %d - %s", resp.StatusCode, string(body))
	}

	var check OnfidoCheck
	if err := json.NewDecoder(resp.Body).Decode(&check); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &check, nil
}

// Get check status and results
func (ml *MLService) GetCheck(ctx context.Context, checkID string) (*OnfidoCheck, error) {
	url := fmt.Sprintf("%s/checks/%s", ml.onfidoConfig.BaseURL, checkID)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Token token="+ml.onfidoConfig.APIKey)

	resp, err := ml.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("onfido API error: %d - %s", resp.StatusCode, string(body))
	}

	var check OnfidoCheck
	if err := json.NewDecoder(resp.Body).Decode(&check); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &check, nil
}

// Get report details
func (ml *MLService) GetReport(ctx context.Context, reportID string) (*OnfidoReport, error) {
	url := fmt.Sprintf("%s/reports/%s", ml.onfidoConfig.BaseURL, reportID)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Token token="+ml.onfidoConfig.APIKey)

	resp, err := ml.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("onfido API error: %d - %s", resp.StatusCode, string(body))
	}

	var report OnfidoReport
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &report, nil
}

// Global ML service instance
var mlService *MLService
var logger *zap.Logger

func init() {
	// Initialize logger
	logger, _ = zap.NewProduction()

	// Initialize Onfido configuration
	config := OnfidoConfig{
		APIKey:  getEnvOrDefault("ONFIDO_API_KEY", "test_api_key"),
		BaseURL: getEnvOrDefault("ONFIDO_BASE_URL", "https://api.eu.onfido.com/v3.6"),
		Region:  getEnvOrDefault("ONFIDO_REGION", "EU"),
	}

	mlService = NewMLService(logger, config)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// HTTP Handlers

// Health check
func healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "ml-svc"})
}

// Create KYC applicant
func createApplicant(c *gin.Context) {
	var req KYCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	logger.Info("Creating KYC applicant", zap.String("first_name", req.FirstName), zap.String("last_name", req.LastName))

	applicant, err := mlService.CreateApplicant(c.Request.Context(), req)
	if err != nil {
		logger.Error("Failed to create applicant", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create applicant"})
		return
	}

	logger.Info("Applicant created successfully", zap.String("applicant_id", applicant.ID))
	c.JSON(http.StatusCreated, applicant)
}

// Upload document
func uploadDocument(c *gin.Context) {
	applicantID := c.PostForm("applicant_id")
	docType := c.PostForm("type")
	side := c.PostForm("side")

	if applicantID == "" || docType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "applicant_id and type are required"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		logger.Error("Failed to get file from request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		logger.Error("Failed to read file data", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	logger.Info("Uploading document", 
		zap.String("applicant_id", applicantID),
		zap.String("type", docType),
		zap.String("filename", header.Filename))

	document, err := mlService.UploadDocument(c.Request.Context(), applicantID, docType, side, fileData, header.Filename)
	if err != nil {
		logger.Error("Failed to upload document", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload document"})
		return
	}

	logger.Info("Document uploaded successfully", zap.String("document_id", document.ID))
	c.JSON(http.StatusCreated, document)
}

// Upload live photo for biometric verification
func uploadLivePhoto(c *gin.Context) {
	applicantID := c.PostForm("applicant_id")

	if applicantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "applicant_id is required"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		logger.Error("Failed to get file from request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		logger.Error("Failed to read file data", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	logger.Info("Uploading live photo", 
		zap.String("applicant_id", applicantID),
		zap.String("filename", header.Filename))

	livePhoto, err := mlService.UploadLivePhoto(c.Request.Context(), applicantID, fileData, header.Filename)
	if err != nil {
		logger.Error("Failed to upload live photo", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload live photo"})
		return
	}

	logger.Info("Live photo uploaded successfully", zap.String("live_photo_id", livePhoto.ID))
	c.JSON(http.StatusCreated, livePhoto)
}

// Create verification check
func createCheck(c *gin.Context) {
	var req struct {
		ApplicantID   string   `json:"applicant_id"`
		ReportNames   []string `json:"report_names"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.ApplicantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "applicant_id is required"})
		return
	}

	// Default report names if not provided
	if len(req.ReportNames) == 0 {
		req.ReportNames = []string{"document", "facial_similarity_photo"}
	}

	logger.Info("Creating verification check", 
		zap.String("applicant_id", req.ApplicantID),
		zap.Strings("report_names", req.ReportNames))

	check, err := mlService.CreateCheck(c.Request.Context(), req.ApplicantID, req.ReportNames)
	if err != nil {
		logger.Error("Failed to create check", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create check"})
		return
	}

	logger.Info("Check created successfully", zap.String("check_id", check.ID))
	c.JSON(http.StatusCreated, check)
}

// Get check status
func getCheck(c *gin.Context) {
	checkID := c.Param("id")
	if checkID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "check ID is required"})
		return
	}

	logger.Info("Getting check status", zap.String("check_id", checkID))

	check, err := mlService.GetCheck(c.Request.Context(), checkID)
	if err != nil {
		logger.Error("Failed to get check", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get check"})
		return
	}

	c.JSON(http.StatusOK, check)
}

// Get report details
func getReport(c *gin.Context) {
	reportID := c.Param("id")
	if reportID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "report ID is required"})
		return
	}

	logger.Info("Getting report details", zap.String("report_id", reportID))

	report, err := mlService.GetReport(c.Request.Context(), reportID)
	if err != nil {
		logger.Error("Failed to get report", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get report"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// Complete KYC verification workflow
func completeKYCVerification(c *gin.Context) {
	var req struct {
		ApplicantID string `json:"applicant_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.ApplicantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "applicant_id is required"})
		return
	}

	logger.Info("Starting complete KYC verification", zap.String("applicant_id", req.ApplicantID))

	// Create comprehensive check with all verification types
	reportNames := []string{
		"document",                    // Document verification
		"facial_similarity_photo",     // Face matching
		"identity_enhanced",           // Enhanced identity verification
		"watchlist_enhanced",          // Watchlist screening
	}

	check, err := mlService.CreateCheck(c.Request.Context(), req.ApplicantID, reportNames)
	if err != nil {
		logger.Error("Failed to create comprehensive check", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create verification check"})
		return
	}

	logger.Info("Complete KYC verification initiated", 
		zap.String("check_id", check.ID),
		zap.String("applicant_id", req.ApplicantID))

	c.JSON(http.StatusCreated, gin.H{
		"check_id":     check.ID,
		"applicant_id": req.ApplicantID,
		"status":       check.Status,
		"message":      "KYC verification initiated. Use the check_id to monitor progress.",
	})
}

func main() {
	defer logger.Sync()

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Health check
	r.GET("/healthz", healthz)

	// KYC verification endpoints
	v1 := r.Group("/api/v1")
	{
		// Applicant management
		v1.POST("/applicants", createApplicant)
		
		// Document upload
		v1.POST("/documents", uploadDocument)
		
		// Biometric verification
		v1.POST("/live-photos", uploadLivePhoto)
		
		// Verification checks
		v1.POST("/checks", createCheck)
		v1.GET("/checks/:id", getCheck)
		
		// Reports
		v1.GET("/reports/:id", getReport)
		
		// Complete KYC workflow
		v1.POST("/kyc/verify", completeKYCVerification)
	}

	port := getEnvOrDefault("PORT", "8085")
	logger.Info("Starting ML service", zap.String("port", port))

	if err := r.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}