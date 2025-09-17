package main

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Mock Onfido server for testing
func setupMockOnfidoServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authorization header
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Token token=") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch {
		case r.Method == "POST" && r.URL.Path == "/applicants":
			// Create applicant
			var req KYCRequest
			json.NewDecoder(r.Body).Decode(&req)
			
			applicant := OnfidoApplicant{
				ID:        "test-applicant-123",
				CreatedAt: time.Now(),
				FirstName: req.FirstName,
				LastName:  req.LastName,
				Email:     req.Email,
				DOB:       req.DOB,
			}
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(applicant)

		case r.Method == "POST" && r.URL.Path == "/documents":
			// Upload document
			r.ParseMultipartForm(10 << 20) // 10MB
			applicantID := r.FormValue("applicant_id")
			docType := r.FormValue("type")
			side := r.FormValue("side")
			
			document := OnfidoDocument{
				ID:          "test-document-456",
				CreatedAt:   time.Now(),
				FileName:    "test-document.jpg",
				FileType:    "image/jpeg",
				FileSize:    1024,
				Type:        docType,
				Side:        side,
				ApplicantID: applicantID,
			}
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(document)

		case r.Method == "POST" && r.URL.Path == "/live_photos":
			// Upload live photo
			r.ParseMultipartForm(10 << 20) // 10MB
			applicantID := r.FormValue("applicant_id")
			
			livePhoto := OnfidoLivePhoto{
				ID:          "test-live-photo-789",
				CreatedAt:   time.Now(),
				FileName:    "test-selfie.jpg",
				FileType:    "image/jpeg",
				FileSize:    2048,
				ApplicantID: applicantID,
			}
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(livePhoto)

		case r.Method == "POST" && r.URL.Path == "/checks":
			// Create check
			var req map[string]interface{}
			json.NewDecoder(r.Body).Decode(&req)
			
			check := OnfidoCheck{
				ID:          "test-check-abc",
				CreatedAt:   time.Now(),
				Status:      "in_progress",
				Result:      "",
				ApplicantID: req["applicant_id"].(string),
				ReportIDs:   []string{"test-report-1", "test-report-2"},
			}
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(check)

		case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/checks/"):
			// Get check
			checkID := strings.TrimPrefix(r.URL.Path, "/checks/")
			
			check := OnfidoCheck{
				ID:          checkID,
				CreatedAt:   time.Now().Add(-5 * time.Minute),
				Status:      "complete",
				Result:      "clear",
				ApplicantID: "test-applicant-123",
				ReportIDs:   []string{"test-report-1", "test-report-2"},
			}
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(check)

		case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/reports/"):
			// Get report
			reportID := strings.TrimPrefix(r.URL.Path, "/reports/")
			
			report := OnfidoReport{
				ID:        reportID,
				CreatedAt: time.Now().Add(-3 * time.Minute),
				Name:      "document",
				Status:    "complete",
				Result:    "clear",
				SubResult: "clear",
				Properties: map[string]interface{}{
					"document_type": "passport",
					"issuing_country": "GBR",
					"document_numbers": []map[string]string{
						{"type": "document_number", "value": "502135326"},
					},
				},
			}
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(report)

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func setupTestMLService() *MLService {
	mockServer := setupMockOnfidoServer()
	
	config := OnfidoConfig{
		APIKey:  "test_api_key",
		BaseURL: mockServer.URL,
		Region:  "EU",
	}
	
	testLogger, _ := zap.NewDevelopment()
	return NewMLService(testLogger, config)
}

func setupTestRouter() *gin.Engine {
	// Override global mlService for testing
	mlService = setupTestMLService()
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	// Health check
	r.GET("/healthz", healthz)
	
	// KYC verification endpoints
	v1 := r.Group("/api/v1")
	{
		v1.POST("/applicants", createApplicant)
		v1.POST("/documents", uploadDocument)
		v1.POST("/live-photos", uploadLivePhoto)
		v1.POST("/checks", createCheck)
		v1.GET("/checks/:id", getCheck)
		v1.GET("/reports/:id", getReport)
		v1.POST("/kyc/verify", completeKYCVerification)
	}
	
	return r
}

func TestHealthz(t *testing.T) {
	router := setupTestRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", response["status"])
	}
	
	if response["service"] != "ml-svc" {
		t.Errorf("Expected service 'ml-svc', got %v", response["service"])
	}
}

func TestCreateApplicant(t *testing.T) {
	router := setupTestRouter()
	
	reqBody := KYCRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		DOB:       "1990-01-01",
	}
	
	jsonData, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/applicants", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
	
	var applicant OnfidoApplicant
	json.Unmarshal(w.Body.Bytes(), &applicant)
	
	if applicant.ID == "" {
		t.Error("Expected applicant ID to be set")
	}
	
	if applicant.FirstName != "John" {
		t.Errorf("Expected first name 'John', got %s", applicant.FirstName)
	}
	
	if applicant.LastName != "Doe" {
		t.Errorf("Expected last name 'Doe', got %s", applicant.LastName)
	}
}

func TestUploadDocument(t *testing.T) {
	router := setupTestRouter()
	
	// Create multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Add form fields
	writer.WriteField("applicant_id", "test-applicant-123")
	writer.WriteField("type", "passport")
	writer.WriteField("side", "front")
	
	// Add file
	part, _ := writer.CreateFormFile("file", "passport.jpg")
	part.Write([]byte("fake image data"))
	writer.Close()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/documents", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
	
	var document OnfidoDocument
	json.Unmarshal(w.Body.Bytes(), &document)
	
	if document.ID == "" {
		t.Error("Expected document ID to be set")
	}
	
	if document.Type != "passport" {
		t.Errorf("Expected document type 'passport', got %s", document.Type)
	}
	
	if document.ApplicantID != "test-applicant-123" {
		t.Errorf("Expected applicant ID 'test-applicant-123', got %s", document.ApplicantID)
	}
}

func TestUploadLivePhoto(t *testing.T) {
	router := setupTestRouter()
	
	// Create multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Add form fields
	writer.WriteField("applicant_id", "test-applicant-123")
	
	// Add file
	part, _ := writer.CreateFormFile("file", "selfie.jpg")
	part.Write([]byte("fake selfie data"))
	writer.Close()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/live-photos", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
	
	var livePhoto OnfidoLivePhoto
	json.Unmarshal(w.Body.Bytes(), &livePhoto)
	
	if livePhoto.ID == "" {
		t.Error("Expected live photo ID to be set")
	}
	
	if livePhoto.ApplicantID != "test-applicant-123" {
		t.Errorf("Expected applicant ID 'test-applicant-123', got %s", livePhoto.ApplicantID)
	}
}

func TestCreateCheck(t *testing.T) {
	router := setupTestRouter()
	
	reqBody := map[string]interface{}{
		"applicant_id": "test-applicant-123",
		"report_names": []string{"document", "facial_similarity_photo"},
	}
	
	jsonData, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/checks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
	
	var check OnfidoCheck
	json.Unmarshal(w.Body.Bytes(), &check)
	
	if check.ID == "" {
		t.Error("Expected check ID to be set")
	}
	
	if check.ApplicantID != "test-applicant-123" {
		t.Errorf("Expected applicant ID 'test-applicant-123', got %s", check.ApplicantID)
	}
	
	if check.Status != "in_progress" {
		t.Errorf("Expected status 'in_progress', got %s", check.Status)
	}
}

func TestGetCheck(t *testing.T) {
	router := setupTestRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/checks/test-check-abc", nil)
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	var check OnfidoCheck
	json.Unmarshal(w.Body.Bytes(), &check)
	
	if check.ID != "test-check-abc" {
		t.Errorf("Expected check ID 'test-check-abc', got %s", check.ID)
	}
	
	if check.Status != "complete" {
		t.Errorf("Expected status 'complete', got %s", check.Status)
	}
	
	if check.Result != "clear" {
		t.Errorf("Expected result 'clear', got %s", check.Result)
	}
}

func TestGetReport(t *testing.T) {
	router := setupTestRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/reports/test-report-1", nil)
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	var report OnfidoReport
	json.Unmarshal(w.Body.Bytes(), &report)
	
	if report.ID != "test-report-1" {
		t.Errorf("Expected report ID 'test-report-1', got %s", report.ID)
	}
	
	if report.Name != "document" {
		t.Errorf("Expected report name 'document', got %s", report.Name)
	}
	
	if report.Result != "clear" {
		t.Errorf("Expected result 'clear', got %s", report.Result)
	}
}

func TestCompleteKYCVerification(t *testing.T) {
	router := setupTestRouter()
	
	reqBody := map[string]string{
		"applicant_id": "test-applicant-123",
	}
	
	jsonData, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/kyc/verify", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["check_id"] == "" {
		t.Error("Expected check_id to be set")
	}
	
	if response["applicant_id"] != "test-applicant-123" {
		t.Errorf("Expected applicant_id 'test-applicant-123', got %v", response["applicant_id"])
	}
	
	if response["status"] != "in_progress" {
		t.Errorf("Expected status 'in_progress', got %v", response["status"])
	}
}

func TestMLServiceCreateApplicant(t *testing.T) {
	mlSvc := setupTestMLService()
	
	req := KYCRequest{
		FirstName: "Alice",
		LastName:  "Smith",
		Email:     "alice.smith@example.com",
		DOB:       "1985-05-15",
	}
	
	applicant, err := mlSvc.CreateApplicant(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to create applicant: %v", err)
	}
	
	if applicant.ID == "" {
		t.Error("Expected applicant ID to be set")
	}
	
	if applicant.FirstName != "Alice" {
		t.Errorf("Expected first name 'Alice', got %s", applicant.FirstName)
	}
}

func TestMLServiceUploadDocument(t *testing.T) {
	mlSvc := setupTestMLService()
	
	fakeImageData := []byte("fake passport image data")
	document, err := mlSvc.UploadDocument(context.Background(), "test-applicant-123", "passport", "front", fakeImageData, "passport.jpg")
	if err != nil {
		t.Fatalf("Failed to upload document: %v", err)
	}
	
	if document.ID == "" {
		t.Error("Expected document ID to be set")
	}
	
	if document.Type != "passport" {
		t.Errorf("Expected document type 'passport', got %s", document.Type)
	}
}

func TestMLServiceCreateCheck(t *testing.T) {
	mlSvc := setupTestMLService()
	
	reportNames := []string{"document", "facial_similarity_photo"}
	check, err := mlSvc.CreateCheck(context.Background(), "test-applicant-123", reportNames)
	if err != nil {
		t.Fatalf("Failed to create check: %v", err)
	}
	
	if check.ID == "" {
		t.Error("Expected check ID to be set")
	}
	
	if check.ApplicantID != "test-applicant-123" {
		t.Errorf("Expected applicant ID 'test-applicant-123', got %s", check.ApplicantID)
	}
}

func TestInvalidRequests(t *testing.T) {
	router := setupTestRouter()
	
	// Test invalid JSON for create applicant
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/applicants", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusBadRequest, w.Code)
	}
	
	// Test missing applicant_id for document upload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("type", "passport")
	part, _ := writer.CreateFormFile("file", "test.jpg")
	part.Write([]byte("test"))
	writer.Close()
	
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/documents", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for missing applicant_id, got %d", http.StatusBadRequest, w.Code)
	}
}