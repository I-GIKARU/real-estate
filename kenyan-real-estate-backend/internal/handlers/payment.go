package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"kenyan-real-estate-backend/internal/models"
	"kenyan-real-estate-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	paymentRepo *models.PaymentRepository
	leaseRepo   *models.LeaseRepository
	mpesaService *services.MPesaService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentRepo *models.PaymentRepository, leaseRepo *models.LeaseRepository, mpesaService *services.MPesaService) *PaymentHandler {
	return &PaymentHandler{
		paymentRepo:  paymentRepo,
		leaseRepo:    leaseRepo,
		mpesaService: mpesaService,
	}
}

// InitiateRentPaymentRequest represents the request to initiate rent payment
type InitiateRentPaymentRequest struct {
	LeaseID     uuid.UUID `json:"lease_id" binding:"required"`
	Amount      float64   `json:"amount" binding:"required,min=1"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	PaymentType string    `json:"payment_type" binding:"required"`
}

// InitiateRentPayment initiates a rent payment via M-Pesa STK Push
func (h *PaymentHandler) InitiateRentPayment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	tenantID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	var req InitiateRentPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Verify lease exists and belongs to the tenant
	lease, err := h.leaseRepo.GetByID(req.LeaseID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Lease not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get lease",
		})
		return
	}

	if lease.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only make payments for your own lease",
		})
		return
	}

	// Validate payment type
	paymentType := models.PaymentType(req.PaymentType)
	if paymentType != models.PaymentTypeRent && paymentType != models.PaymentTypeDeposit && 
	   paymentType != models.PaymentTypeUtility && paymentType != models.PaymentTypeMaintenance {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid payment type",
		})
		return
	}

	// Validate and format phone number
	formattedPhone, err := h.mpesaService.ValidatePhoneNumber(req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid phone number format",
			"details": err.Error(),
		})
		return
	}

	// Create payment record
	payment := &models.Payment{
		LeaseID:       req.LeaseID,
		Amount:        req.Amount,
		PaymentType:   paymentType,
		PaymentMethod: models.PaymentMethodMPesa,
		Status:        models.PaymentStatusPending,
	}

	if err := h.paymentRepo.Create(payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create payment record",
		})
		return
	}

	// Initiate M-Pesa STK Push
	amount := h.mpesaService.FormatAmount(req.Amount)
	accountReference := fmt.Sprintf("LEASE-%s", req.LeaseID.String()[:8])
	description := fmt.Sprintf("Rent payment for lease %s", accountReference)
	callbackURL := fmt.Sprintf("%s/api/v1/payments/mpesa/callback", getBaseURL(c))

	stkResponse, err := h.mpesaService.InitiateSTKPush(
		formattedPhone,
		amount,
		accountReference,
		description,
		callbackURL,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to initiate M-Pesa payment",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":             "Payment initiated successfully",
		"payment_id":          payment.ID,
		"checkout_request_id": stkResponse.CheckoutRequestID,
		"customer_message":    stkResponse.CustomerMessage,
	})
}

// MPesaCallbackRequest represents the M-Pesa callback request
type MPesaCallbackRequest struct {
	Body struct {
		StkCallback struct {
			MerchantRequestID string `json:"MerchantRequestID"`
			CheckoutRequestID string `json:"CheckoutRequestID"`
			ResultCode        int    `json:"ResultCode"`
			ResultDesc        string `json:"ResultDesc"`
			CallbackMetadata  struct {
				Item []struct {
					Name  string      `json:"Name"`
					Value interface{} `json:"Value"`
				} `json:"Item"`
			} `json:"CallbackMetadata"`
		} `json:"stkCallback"`
	} `json:"Body"`
}

// HandleMPesaCallback handles M-Pesa payment callbacks
func (h *PaymentHandler) HandleMPesaCallback(c *gin.Context) {
	var callback MPesaCallbackRequest
	if err := c.ShouldBindJSON(&callback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid callback data",
		})
		return
	}

	stkCallback := callback.Body.StkCallback
	
	// Extract transaction details from callback metadata
	var transactionID, phoneNumber string
	var amount float64

	if stkCallback.ResultCode == 0 { // Success
		for _, item := range stkCallback.CallbackMetadata.Item {
			switch item.Name {
			case "MpesaReceiptNumber":
				if val, ok := item.Value.(string); ok {
					transactionID = val
				}
			case "PhoneNumber":
				if val, ok := item.Value.(float64); ok {
					phoneNumber = fmt.Sprintf("%.0f", val)
				}
			case "Amount":
				if val, ok := item.Value.(float64); ok {
					amount = val
				}
			}
		}

		// Update payment status to completed
		// Note: In a real implementation, you would need to find the payment by checkout_request_id
		// For now, we'll just log the successful payment
		fmt.Printf("Payment successful: TransactionID=%s, Phone=%s, Amount=%.2f\n", 
			transactionID, phoneNumber, amount)
	} else {
		// Payment failed
		fmt.Printf("Payment failed: %s\n", stkCallback.ResultDesc)
	}

	// Always respond with success to M-Pesa
	c.JSON(http.StatusOK, gin.H{
		"ResultCode": 0,
		"ResultDesc": "Success",
	})
}

// GetPaymentsByLease gets payments for a specific lease
func (h *PaymentHandler) GetPaymentsByLease(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	leaseIDStr := c.Param("lease_id")
	leaseID, err := uuid.Parse(leaseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid lease ID",
		})
		return
	}

	// Verify lease exists and user has access to it
	lease, err := h.leaseRepo.GetByID(leaseID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Lease not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get lease",
		})
		return
	}

	// Check if user is tenant or landlord of this lease
	if lease.TenantID != userUUID && lease.LandlordID != userUUID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to this lease",
		})
		return
	}

	payments, err := h.paymentRepo.GetByLeaseID(leaseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get payments",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
		"lease_id": leaseID,
	})
}

// QueryPaymentStatus queries the status of an M-Pesa payment
func (h *PaymentHandler) QueryPaymentStatus(c *gin.Context) {
	checkoutRequestID := c.Param("checkout_request_id")
	if checkoutRequestID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Checkout request ID is required",
		})
		return
	}

	queryResponse, err := h.mpesaService.QuerySTKPush(checkoutRequestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to query payment status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": queryResponse,
	})
}

// Helper function to get base URL
func getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}

