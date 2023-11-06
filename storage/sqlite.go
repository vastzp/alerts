package storage

import (
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/vastzp/alerts/service"
	"strconv"
)

// SQLiteStorage is a structure for the SQLite storage. This is implementation of the service layer interface (Storage)
type SQLiteStorage struct {
	db *sqlx.DB
}

// NewSQLiteStorage creates a new SQLite storage (constructor).
func NewSQLiteStorage(storageFilename string) (*SQLiteStorage, error) {
	db, err := sqlx.Open("sqlite", storageFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	s := &SQLiteStorage{
		db: db,
	}
	err = s.migrate()
	if err != nil {
		return nil, fmt.Errorf("failed to migrate SQLite database: %w", err)
	}
	return s, nil
}

// GetAlerts returns all alerts from the storage.
func (s *SQLiteStorage) GetAlerts(serviceID string, startTS, endTS int64) ([]*service.Alert, error) {

	rows, err := s.db.Query(`SELECT 
	    * 
	FROM 
    	alerts 
	WHERE 
	    service_id = ? AND alert_ts BETWEEN ? AND ?`, serviceID, startTS, endTS)
	if err != nil {
		return nil, fmt.Errorf("failed to query alerts: %w", err)
	}

	alerts := []*service.Alert{}
	for rows.Next() {
		var alertID, serviceID, serviceName, model, alertType, alertTS, severity, teamSlack string
		err = rows.Scan(&alertID, &serviceID, &serviceName, &model, &alertType, &alertTS, &severity, &teamSlack)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alert: %w", err)
		}
		alerts = append(alerts, &service.Alert{
			AlertID:     alertID,
			ServiceID:   serviceID,
			ServiceName: serviceName,
			Model:       model,
			AlertType:   alertType,
			AlertTs:     alertTS,
			Severity:    severity,
			TeamSlack:   teamSlack})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan alert: %w", err)
	}

	return alerts, nil
}

// alertDTO is a structure for the alert data transfer object. I use it for the SQLite storage routine.
type alertDTO struct {
	AlertID     string
	ServiceID   string
	ServiceName string
	Model       string
	AlertType   string
	AlertTs     int64
	Severity    string
	TeamSlack   string
}

// InsertAlert inserts an alert into the storage.
// Returns the ID of the newly inserted alert and the error if it was occurred.
func (s *SQLiteStorage) InsertAlert(alert *service.Alert) (string, error) {
	alertTs, err := strconv.Atoi(alert.AlertTs)
	if err != nil {
		return "", fmt.Errorf("failed to convert alertTs: %w", err)
	}

	a := &alertDTO{
		AlertID:     alert.AlertID,
		ServiceID:   alert.ServiceID,
		ServiceName: alert.ServiceName,
		Model:       alert.Model,
		AlertType:   alert.AlertType,
		AlertTs:     int64(alertTs),
		Severity:    alert.Severity,
		TeamSlack:   alert.TeamSlack,
	}
	_, err = s.db.Exec(`INSERT INTO alerts (
                    alert_id, 
                    service_id, 
                    service_name, 
                    model, 
                    alert_type, 
                    alert_ts,
                    severity, 
                    team_slack) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		a.AlertID,
		a.ServiceID,
		a.ServiceName,
		a.Model,
		a.AlertType,
		a.AlertTs,
		a.Severity,
		a.TeamSlack)
	if err != nil {
		return "", fmt.Errorf("failed to insert alert: %w", err)
	}

	return alert.AlertID, nil
}

// Close closes the storage.
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

// migrate migrates the SQLite database. I implemented this function for quick "migration" to create an 'alerts' table.
func (s *SQLiteStorage) migrate() error {
	// Create a table
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS alerts (
			alert_id VARCHAR PRIMARY KEY,
			service_id VARCHAR NOT NULL,
			service_name VARCHAR NOT NULL,
			model VARCHAR NOT NULL,
			alert_type VARCHAR NOT NULL,
			alert_ts INTEGER NOT NULL,
			severity VARCHAR NOT NULL,
			team_slack VARCHAR NOT NULL
		);
	`
	_, err := s.db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	return nil
}
