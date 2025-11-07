package rimpay

import (
	"time"

	"github.com/CatoSystems/rim-pay/pkg/money"
)

type TransactionStatus struct {
	TransactionID     string                 `json:"transaction_id"`
	Status            PaymentStatus          `json:"status"`
	Amount            money.Money            `json:"amount,omitempty"`
	Reference         string                 `json:"reference"`
	ProviderReference string                 `json:"provider_reference,omitempty"`
	Message           string                 `json:"message,omitempty"`
	LastUpdated       time.Time              `json:"last_updated"`
	Events            []StatusEvent          `json:"events,omitempty"`
	ProviderData      map[string]interface{} `json:"provider_data,omitempty"`
}

// StatusEvent represents status change event
type StatusEvent struct {
	Status    PaymentStatus          `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Message   string                 `json:"message,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// AddEvent adds a status event
func (ts *TransactionStatus) AddEvent(status PaymentStatus, message string) {
	event := StatusEvent{
		Status:    status,
		Timestamp: time.Now(),
		Message:   message,
		Metadata:  make(map[string]interface{}),
	}
	ts.Events = append(ts.Events, event)
	ts.Status = status
	ts.LastUpdated = event.Timestamp
}

// GetLatestEvent returns the most recent status event
func (ts *TransactionStatus) GetLatestEvent() *StatusEvent {
	if len(ts.Events) == 0 {
		return nil
	}
	return &ts.Events[len(ts.Events)-1]
}

// IsCompleted returns true if transaction is completed
func (ts *TransactionStatus) IsCompleted() bool {
	return ts.Status.IsCompleted()
}

// IsSuccessful returns true if transaction was successful
func (ts *TransactionStatus) IsSuccessful() bool {
	return ts.Status.IsSuccessful()
}
