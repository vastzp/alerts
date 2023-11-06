package storage

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/vastzp/alerts/service"
	"testing"
	"time"
)

func TestSQLiteStorageInit(t *testing.T) {

	repo, err := NewSQLiteStorage("test.db")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	err = repo.Close()
	if err != nil {
		t.Fail()
	}
	timeNow := time.Now()
	timeNowTS := timeNow.Unix()
	serviceID := "my_service_id"
	alertId, err := repo.InsertAlert(&service.Alert{
		AlertID:     "alert_id" + uuid.New().String(),
		ServiceID:   serviceID,
		ServiceName: "service_name" + uuid.New().String(),
		Model:       "model" + uuid.New().String(),
		AlertType:   "alert_type" + uuid.New().String(),
		AlertTs:     fmt.Sprintf("%d", timeNowTS),
		Severity:    "severity" + uuid.New().String(),
		TeamSlack:   "team_slack" + uuid.New().String(),
	})
	if err != nil {
		t.Errorf("failed to insert alert: %s", err)
	}
	println(alertId)
	startTimeNowTS, endTimeNowTS := timeNowTS-10, timeNowTS+10
	fmt.Println("get alerts between", startTimeNowTS, "and", endTimeNowTS, "for service", serviceID, ":")
	alerts, err := repo.GetAlerts(serviceID, startTimeNowTS, endTimeNowTS)
	if err != nil {
		t.Errorf("failed to get alerts: %s", err)
	}
	for _, alert := range alerts {
		fmt.Printf("alert: %s\n%v\n\n", alert.AlertID, alert)
	}
}
