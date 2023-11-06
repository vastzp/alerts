package service

import (
	"fmt"
)

// Alert represents an alert (I put here all the fields from the task as strings).
type Alert struct {
	AlertID     string `json:"alert_id"`
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	Model       string `json:"model"`
	AlertType   string `json:"alert_type"`
	AlertTs     string `json:"alert_ts"`
	Severity    string `json:"severity"`
	TeamSlack   string `json:"team_slack"`
}

// Service is an interface for the service layer.
type Service interface {
	// GetAlerts returns all alerts from the storage.
	GetAlerts(serviceID string, startTS, endTS int64) ([]*Alert, error)

	// InsertAlert inserts an alert into the storage.
	InsertAlert(alert *Alert) (string, error)
}

// Storage is an interface for the storage layer.
type Storage interface {
	// GetAlerts returns all alerts from the storage.
	GetAlerts(serviceID string, startTS, endTS int64) ([]*Alert, error)

	// InsertAlert inserts an alert into the storage.
	// Returns the ID of the newly inserted alert and the error if it was occurred.
	InsertAlert(alert *Alert) (string, error)

	// Close closes the storage.
	Close() error
}

// service is a structure for the service instance.
type service struct {
	storage Storage
}

// GetAlerts returns all alerts from the storage.
func (s service) GetAlerts(serviceID string, startTS, endTS int64) ([]*Alert, error) {
	alerts, err := s.storage.GetAlerts(serviceID, startTS, endTS)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %w", err)
	}
	return alerts, nil
}

// InsertAlert inserts an alert into the storage.
func (s service) InsertAlert(alert *Alert) (string, error) {
	alertID, err := s.storage.InsertAlert(alert)
	if err != nil {
		return "", err
	}
	return alertID, nil
}

// NewService creates a new service instance (constructor).
func NewService(storage Storage) Service {
	return &service{
		storage: storage,
	}
}
