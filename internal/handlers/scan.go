package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/datmedevil17/go-vuln/internal/middleware"
	"github.com/datmedevil17/go-vuln/internal/services"
	"github.com/datmedevil17/go-vuln/internal/utils"
)

type ScannerHandler struct {
	scannerService *services.ScannerService
}

type ScanRequest struct {
	Target   string `json:"target" binding:"required"`
	Ports    string `json:"ports,omitempty"`
	Wordlist string `json:"wordlist,omitempty"`
}

func NewScannerHandler(scannerService *services.ScannerService) *ScannerHandler {
	return &ScannerHandler{
		scannerService: scannerService,
	}
}

// NmapScan initiates an Nmap scan
func (h *ScannerHandler) NmapScan(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req ScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error())
		return
	}

	ports := req.Ports
	if ports == "" {
		ports = "1-1000"
	}

	result, err := h.scannerService.NmapScan(c.Request.Context(), userID, req.Target, ports)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to start scan: "+err.Error())
		return
	}

	utils.SuccessMessageResponse(c, "Nmap scan started", result)
}

// NiktoScan initiates a Nikto scan
func (h *ScannerHandler) NiktoScan(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req ScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error())
		return
	}

	result, err := h.scannerService.NiktoScan(c.Request.Context(), userID, req.Target)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to start scan: "+err.Error())
		return
	}

	utils.SuccessMessageResponse(c, "Nikto scan started", result)
}

// GobusterScan initiates a Gobuster scan
func (h *ScannerHandler) GobusterScan(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req ScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error())
		return
	}

	result, err := h.scannerService.GobusterScan(c.Request.Context(), userID, req.Target, req.Wordlist)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to start scan: "+err.Error())
		return
	}

	utils.SuccessMessageResponse(c, "Gobuster scan started", result)
}

// GetScanResult retrieves a scan result
func (h *ScannerHandler) GetScanResult(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	scanID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid scan ID")
		return
	}

	result, err := h.scannerService.GetScanResult(scanID, userID)
	if err != nil {
		utils.NotFoundResponse(c, "Scan result not found")
		return
	}

	utils.SuccessResponse(c, result)
}

// ListScanResults lists all scan results
func (h *ScannerHandler) ListScanResults(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	results, err := h.scannerService.ListScanResults(userID)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to fetch scan results")
		return
	}

	utils.SuccessResponse(c, results)
}