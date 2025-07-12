package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"real-estate-backend/internal/config"
)

// MPesaService handles M-Pesa API integrations
type MPesaService struct {
	config *config.MPesaConfig
	client *http.Client
}

// NewMPesaService creates a new M-Pesa service
func NewMPesaService(cfg *config.MPesaConfig) *MPesaService {
	return &MPesaService{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AuthResponse represents the M-Pesa authentication response
type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

// STKPushRequest represents the STK Push request
type STKPushRequest struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	TransactionType   string `json:"TransactionType"`
	Amount            string `json:"Amount"`
	PartyA            string `json:"PartyA"`
	PartyB            string `json:"PartyB"`
	PhoneNumber       string `json:"PhoneNumber"`
	CallBackURL       string `json:"CallBackURL"`
	AccountReference  string `json:"AccountReference"`
	TransactionDesc   string `json:"TransactionDesc"`
}

// STKPushResponse represents the STK Push response
type STKPushResponse struct {
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
}

// STKQueryRequest represents the STK Push query request
type STKQueryRequest struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	CheckoutRequestID string `json:"CheckoutRequestID"`
}

// STKQueryResponse represents the STK Push query response
type STKQueryResponse struct {
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResultCode          string `json:"ResultCode"`
	ResultDesc          string `json:"ResultDesc"`
}

// C2BCallbackData represents the C2B callback data
type C2BCallbackData struct {
	TransactionType   string `json:"TransactionType"`
	TransID           string `json:"TransID"`
	TransTime         string `json:"TransTime"`
	TransAmount       string `json:"TransAmount"`
	BusinessShortCode string `json:"BusinessShortCode"`
	BillRefNumber     string `json:"BillRefNumber"`
	InvoiceNumber     string `json:"InvoiceNumber"`
	OrgAccountBalance string `json:"OrgAccountBalance"`
	ThirdPartyTransID string `json:"ThirdPartyTransID"`
	MSISDN            string `json:"MSISDN"`
	FirstName         string `json:"FirstName"`
	MiddleName        string `json:"MiddleName"`
	LastName          string `json:"LastName"`
}

// GetAccessToken gets an access token from M-Pesa API
func (m *MPesaService) GetAccessToken() (*AuthResponse, error) {
	url := fmt.Sprintf("https://%s.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials", m.config.Environment)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Create basic auth header
	auth := base64.StdEncoding.EncodeToString([]byte(m.config.ConsumerKey + ":" + m.config.ConsumerSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication failed: %s", string(body))
	}

	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &authResp, nil
}

// GeneratePassword generates the password for STK Push
func (m *MPesaService) GeneratePassword(timestamp string) string {
	password := m.config.ShortCode + m.config.PassKey + timestamp
	return base64.StdEncoding.EncodeToString([]byte(password))
}

// GetTimestamp returns the current timestamp in the required format
func (m *MPesaService) GetTimestamp() string {
	return time.Now().Format("20060102150405")
}

// InitiateSTKPush initiates an STK Push payment
func (m *MPesaService) InitiateSTKPush(phoneNumber, amount, accountReference, description, callbackURL string) (*STKPushResponse, error) {
	// Get access token
	auth, err := m.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	timestamp := m.GetTimestamp()
	password := m.GeneratePassword(timestamp)

	// Prepare request
	stkRequest := STKPushRequest{
		BusinessShortCode: m.config.ShortCode,
		Password:          password,
		Timestamp:         timestamp,
		TransactionType:   "CustomerPayBillOnline",
		Amount:            amount,
		PartyA:            phoneNumber,
		PartyB:            m.config.ShortCode,
		PhoneNumber:       phoneNumber,
		CallBackURL:       callbackURL,
		AccountReference:  accountReference,
		TransactionDesc:   description,
	}

	jsonData, err := json.Marshal(stkRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://%s.safaricom.co.ke/mpesa/stkpush/v1/processrequest", m.config.Environment)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var stkResp STKPushResponse
	if err := json.Unmarshal(body, &stkResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("STK push failed: %s - %s", stkResp.ResponseCode, stkResp.ResponseDescription)
	}

	return &stkResp, nil
}

// QuerySTKPush queries the status of an STK Push transaction
func (m *MPesaService) QuerySTKPush(checkoutRequestID string) (*STKQueryResponse, error) {
	// Get access token
	auth, err := m.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	timestamp := m.GetTimestamp()
	password := m.GeneratePassword(timestamp)

	// Prepare request
	queryRequest := STKQueryRequest{
		BusinessShortCode: m.config.ShortCode,
		Password:          password,
		Timestamp:         timestamp,
		CheckoutRequestID: checkoutRequestID,
	}

	jsonData, err := json.Marshal(queryRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://%s.safaricom.co.ke/mpesa/stkpushquery/v1/query", m.config.Environment)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var queryResp STKQueryResponse
	if err := json.Unmarshal(body, &queryResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &queryResp, nil
}

// ValidatePhoneNumber validates and formats a Kenyan phone number for M-Pesa
func (m *MPesaService) ValidatePhoneNumber(phoneNumber string) (string, error) {
	// Remove any spaces, dashes, or plus signs
	cleaned := ""
	for _, char := range phoneNumber {
		if char >= '0' && char <= '9' {
			cleaned += string(char)
		}
	}

	// Handle different formats
	switch {
	case len(cleaned) == 10 && cleaned[:1] == "0":
		// 0712345678 -> 254712345678
		return "254" + cleaned[1:], nil
	case len(cleaned) == 9:
		// 712345678 -> 254712345678
		return "254" + cleaned, nil
	case len(cleaned) == 12 && cleaned[:3] == "254":
		// 254712345678 -> 254712345678
		return cleaned, nil
	default:
		return "", fmt.Errorf("invalid phone number format: %s", phoneNumber)
	}
}

// FormatAmount formats amount to string without decimals (M-Pesa expects whole numbers)
func (m *MPesaService) FormatAmount(amount float64) string {
	return fmt.Sprintf("%.0f", amount)
}

