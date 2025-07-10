package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ApplicationStatus represents the status of a rental application
type ApplicationStatus string

const (
	ApplicationStatusPending   ApplicationStatus = "pending"
	ApplicationStatusApproved  ApplicationStatus = "approved"
	ApplicationStatusRejected  ApplicationStatus = "rejected"
	ApplicationStatusWithdrawn ApplicationStatus = "withdrawn"
)

// LeaseStatus represents the status of a lease
type LeaseStatus string

const (
	LeaseStatusActive     LeaseStatus = "active"
	LeaseStatusExpired    LeaseStatus = "expired"
	LeaseStatusTerminated LeaseStatus = "terminated"
)

// PaymentType represents the type of payment
type PaymentType string

const (
	PaymentTypeRent        PaymentType = "rent"
	PaymentTypeDeposit     PaymentType = "deposit"
	PaymentTypeUtility     PaymentType = "utility"
	PaymentTypeMaintenance PaymentType = "maintenance"
)

// PaymentMethod represents the payment method
type PaymentMethod string

const (
	PaymentMethodMPesa        PaymentMethod = "mpesa"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodCash         PaymentMethod = "cash"
	PaymentMethodCheque       PaymentMethod = "cheque"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// References represents references as JSON
type References map[string]interface{}

// Value implements the driver.Valuer interface for database storage
func (r References) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Scan implements the sql.Scanner interface for database retrieval
func (r *References) Scan(value interface{}) error {
	if value == nil {
		*r = make(References)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, r)
	case string:
		return json.Unmarshal([]byte(v), r)
	default:
		return fmt.Errorf("cannot scan %T into References", value)
	}
}

// RentalApplication represents a rental application
type RentalApplication struct {
	ID               uuid.UUID         `json:"id" db:"id"`
	PropertyID       uuid.UUID         `json:"property_id" db:"property_id"`
	TenantID         uuid.UUID         `json:"tenant_id" db:"tenant_id"`
	ApplicationDate  time.Time         `json:"application_date" db:"application_date"`
	Status           ApplicationStatus `json:"status" db:"status"`
	MoveInDate       *time.Time        `json:"move_in_date,omitempty" db:"move_in_date"`
	Message          *string           `json:"message,omitempty" db:"message"`
	MonthlyIncome    *float64          `json:"monthly_income,omitempty" db:"monthly_income"`
	EmploymentStatus *string           `json:"employment_status,omitempty" db:"employment_status"`
	References       References        `json:"references" db:"references"`
	CreatedAt        time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at" db:"updated_at"`

	// Joined fields
	Property *Property `json:"property,omitempty"`
	Tenant   *User     `json:"tenant,omitempty"`
}

// CreateRentalApplicationRequest represents the request to create a rental application
type CreateRentalApplicationRequest struct {
	PropertyID       uuid.UUID  `json:"property_id" binding:"required"`
	MoveInDate       *time.Time `json:"move_in_date,omitempty"`
	Message          *string    `json:"message,omitempty"`
	MonthlyIncome    *float64   `json:"monthly_income,omitempty" binding:"omitempty,min=0"`
	EmploymentStatus *string    `json:"employment_status,omitempty"`
	References       References `json:"references"`
}

// UpdateApplicationStatusRequest represents the request to update application status
type UpdateApplicationStatusRequest struct {
	Status ApplicationStatus `json:"status" binding:"required"`
}

// Lease represents a lease agreement
type Lease struct {
	ID           uuid.UUID   `json:"id" db:"id"`
	PropertyID   uuid.UUID   `json:"property_id" db:"property_id"`
	TenantID     uuid.UUID   `json:"tenant_id" db:"tenant_id"`
	LandlordID   uuid.UUID   `json:"landlord_id" db:"landlord_id"`
	StartDate    time.Time   `json:"start_date" db:"start_date"`
	EndDate      time.Time   `json:"end_date" db:"end_date"`
	MonthlyRent  float64     `json:"monthly_rent" db:"monthly_rent"`
	DepositPaid  float64     `json:"deposit_paid" db:"deposit_paid"`
	Status       LeaseStatus `json:"status" db:"status"`
	LeaseTerms   *string     `json:"lease_terms,omitempty" db:"lease_terms"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`

	// Joined fields
	Property *Property `json:"property,omitempty"`
	Tenant   *User     `json:"tenant,omitempty"`
	Landlord *User     `json:"landlord,omitempty"`
}

// CreateLeaseRequest represents the request to create a lease
type CreateLeaseRequest struct {
	PropertyID  uuid.UUID `json:"property_id" binding:"required"`
	TenantID    uuid.UUID `json:"tenant_id" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	MonthlyRent float64   `json:"monthly_rent" binding:"required,min=0"`
	DepositPaid float64   `json:"deposit_paid" binding:"required,min=0"`
	LeaseTerms  *string   `json:"lease_terms,omitempty"`
}

// Payment represents a payment record
type Payment struct {
	ID                  uuid.UUID     `json:"id" db:"id"`
	LeaseID             uuid.UUID     `json:"lease_id" db:"lease_id"`
	Amount              float64       `json:"amount" db:"amount"`
	PaymentType         PaymentType   `json:"payment_type" db:"payment_type"`
	PaymentMethod       PaymentMethod `json:"payment_method" db:"payment_method"`
	MPesaTransactionID  *string       `json:"mpesa_transaction_id,omitempty" db:"mpesa_transaction_id"`
	PaymentDate         time.Time     `json:"payment_date" db:"payment_date"`
	DueDate             *time.Time    `json:"due_date,omitempty" db:"due_date"`
	Status              PaymentStatus `json:"status" db:"status"`
	Notes               *string       `json:"notes,omitempty" db:"notes"`
	CreatedAt           time.Time     `json:"created_at" db:"created_at"`

	// Joined fields
	Lease *Lease `json:"lease,omitempty"`
}

// CreatePaymentRequest represents the request to create a payment
type CreatePaymentRequest struct {
	LeaseID             uuid.UUID     `json:"lease_id" binding:"required"`
	Amount              float64       `json:"amount" binding:"required,min=0"`
	PaymentType         PaymentType   `json:"payment_type" binding:"required"`
	PaymentMethod       PaymentMethod `json:"payment_method" binding:"required"`
	MPesaTransactionID  *string       `json:"mpesa_transaction_id,omitempty"`
	DueDate             *time.Time    `json:"due_date,omitempty"`
	Notes               *string       `json:"notes,omitempty"`
}

// RentalApplicationRepository handles database operations for rental applications
type RentalApplicationRepository struct {
	db *sql.DB
}

// NewRentalApplicationRepository creates a new rental application repository
func NewRentalApplicationRepository(db *sql.DB) *RentalApplicationRepository {
	return &RentalApplicationRepository{db: db}
}

// Create creates a new rental application
func (r *RentalApplicationRepository) Create(application *RentalApplication) error {
	query := `
		INSERT INTO rental_applications (
			id, property_id, tenant_id, application_date, status, move_in_date,
			message, monthly_income, employment_status, references, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at`

	application.ID = uuid.New()
	application.ApplicationDate = time.Now()
	application.Status = ApplicationStatusPending
	now := time.Now()
	application.CreatedAt = now
	application.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		application.ID,
		application.PropertyID,
		application.TenantID,
		application.ApplicationDate,
		application.Status,
		application.MoveInDate,
		application.Message,
		application.MonthlyIncome,
		application.EmploymentStatus,
		application.References,
		application.CreatedAt,
		application.UpdatedAt,
	).Scan(&application.ID, &application.CreatedAt, &application.UpdatedAt)

	return err
}

// GetByID retrieves a rental application by ID
func (r *RentalApplicationRepository) GetByID(id uuid.UUID) (*RentalApplication, error) {
	application := &RentalApplication{}
	query := `
		SELECT id, property_id, tenant_id, application_date, status, move_in_date,
			   message, monthly_income, employment_status, references, created_at, updated_at
		FROM rental_applications
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&application.ID,
		&application.PropertyID,
		&application.TenantID,
		&application.ApplicationDate,
		&application.Status,
		&application.MoveInDate,
		&application.Message,
		&application.MonthlyIncome,
		&application.EmploymentStatus,
		&application.References,
		&application.CreatedAt,
		&application.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return application, nil
}

// GetByPropertyID retrieves applications for a property
func (r *RentalApplicationRepository) GetByPropertyID(propertyID uuid.UUID) ([]*RentalApplication, error) {
	query := `
		SELECT id, property_id, tenant_id, application_date, status, move_in_date,
			   message, monthly_income, employment_status, references, created_at, updated_at
		FROM rental_applications
		WHERE property_id = $1
		ORDER BY application_date DESC`

	rows, err := r.db.Query(query, propertyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []*RentalApplication
	for rows.Next() {
		application := &RentalApplication{}
		err := rows.Scan(
			&application.ID,
			&application.PropertyID,
			&application.TenantID,
			&application.ApplicationDate,
			&application.Status,
			&application.MoveInDate,
			&application.Message,
			&application.MonthlyIncome,
			&application.EmploymentStatus,
			&application.References,
			&application.CreatedAt,
			&application.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		applications = append(applications, application)
	}

	return applications, nil
}

// UpdateStatus updates the status of a rental application
func (r *RentalApplicationRepository) UpdateStatus(id uuid.UUID, status ApplicationStatus) error {
	query := `UPDATE rental_applications SET status = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(query, id, status, time.Now())
	return err
}

// LeaseRepository handles database operations for leases
type LeaseRepository struct {
	db *sql.DB
}

// NewLeaseRepository creates a new lease repository
func NewLeaseRepository(db *sql.DB) *LeaseRepository {
	return &LeaseRepository{db: db}
}

// Create creates a new lease
func (r *LeaseRepository) Create(lease *Lease) error {
	query := `
		INSERT INTO leases (
			id, property_id, tenant_id, landlord_id, start_date, end_date,
			monthly_rent, deposit_paid, status, lease_terms, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at`

	lease.ID = uuid.New()
	lease.Status = LeaseStatusActive
	now := time.Now()
	lease.CreatedAt = now
	lease.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		lease.ID,
		lease.PropertyID,
		lease.TenantID,
		lease.LandlordID,
		lease.StartDate,
		lease.EndDate,
		lease.MonthlyRent,
		lease.DepositPaid,
		lease.Status,
		lease.LeaseTerms,
		lease.CreatedAt,
		lease.UpdatedAt,
	).Scan(&lease.ID, &lease.CreatedAt, &lease.UpdatedAt)

	return err
}

// GetByID retrieves a lease by ID
func (r *LeaseRepository) GetByID(id uuid.UUID) (*Lease, error) {
	lease := &Lease{}
	query := `
		SELECT id, property_id, tenant_id, landlord_id, start_date, end_date,
			   monthly_rent, deposit_paid, status, lease_terms, created_at, updated_at
		FROM leases
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&lease.ID,
		&lease.PropertyID,
		&lease.TenantID,
		&lease.LandlordID,
		&lease.StartDate,
		&lease.EndDate,
		&lease.MonthlyRent,
		&lease.DepositPaid,
		&lease.Status,
		&lease.LeaseTerms,
		&lease.CreatedAt,
		&lease.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return lease, nil
}

// GetByTenantID retrieves leases for a tenant
func (r *LeaseRepository) GetByTenantID(tenantID uuid.UUID) ([]*Lease, error) {
	query := `
		SELECT id, property_id, tenant_id, landlord_id, start_date, end_date,
			   monthly_rent, deposit_paid, status, lease_terms, created_at, updated_at
		FROM leases
		WHERE tenant_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leases []*Lease
	for rows.Next() {
		lease := &Lease{}
		err := rows.Scan(
			&lease.ID,
			&lease.PropertyID,
			&lease.TenantID,
			&lease.LandlordID,
			&lease.StartDate,
			&lease.EndDate,
			&lease.MonthlyRent,
			&lease.DepositPaid,
			&lease.Status,
			&lease.LeaseTerms,
			&lease.CreatedAt,
			&lease.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		leases = append(leases, lease)
	}

	return leases, nil
}

// PaymentRepository handles database operations for payments
type PaymentRepository struct {
	db *sql.DB
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Create creates a new payment
func (r *PaymentRepository) Create(payment *Payment) error {
	query := `
		INSERT INTO payments (
			id, lease_id, amount, payment_type, payment_method, mpesa_transaction_id,
			payment_date, due_date, status, notes, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at`

	payment.ID = uuid.New()
	payment.PaymentDate = time.Now()
	payment.Status = PaymentStatusPending
	payment.CreatedAt = time.Now()

	err := r.db.QueryRow(
		query,
		payment.ID,
		payment.LeaseID,
		payment.Amount,
		payment.PaymentType,
		payment.PaymentMethod,
		payment.MPesaTransactionID,
		payment.PaymentDate,
		payment.DueDate,
		payment.Status,
		payment.Notes,
		payment.CreatedAt,
	).Scan(&payment.ID, &payment.CreatedAt)

	return err
}

// GetByLeaseID retrieves payments for a lease
func (r *PaymentRepository) GetByLeaseID(leaseID uuid.UUID) ([]*Payment, error) {
	query := `
		SELECT id, lease_id, amount, payment_type, payment_method, mpesa_transaction_id,
			   payment_date, due_date, status, notes, created_at
		FROM payments
		WHERE lease_id = $1
		ORDER BY payment_date DESC`

	rows, err := r.db.Query(query, leaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*Payment
	for rows.Next() {
		payment := &Payment{}
		err := rows.Scan(
			&payment.ID,
			&payment.LeaseID,
			&payment.Amount,
			&payment.PaymentType,
			&payment.PaymentMethod,
			&payment.MPesaTransactionID,
			&payment.PaymentDate,
			&payment.DueDate,
			&payment.Status,
			&payment.Notes,
			&payment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

