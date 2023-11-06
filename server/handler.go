package server

import (
	"encoding/json"
	"fmt"
	"github.com/vastzp/alerts/service"
	"log"
	"net/http"
	"strconv"
)

// responseGetAlerts is structure for the response data for GET /alerts.
type responseGetAlerts struct {
	Alerts []*service.Alert `json:"alerts"`
	Error  string           `json:"error"`
}

// responsePostAlerts is structure for the response data for POST /alerts.
type responsePostAlerts struct {
	AlertID string `json:"alert_id"`
	Error   string `json:"error"`
}

// handleGetAlerts handles GET /alerts.
func (s *server) handleGetAlerts(w http.ResponseWriter, r *http.Request) {

	// retrieve query parameters
	serviceID := r.URL.Query().Get("service")
	startTs := r.URL.Query().Get("start_ts")
	endTs := r.URL.Query().Get("end_ts")

	// validation for start_s
	startTS, err := strconv.Atoi(startTs)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		rGetAlerts(http.StatusBadRequest, &responseGetAlerts{
			Error: "failed to parse 'start_ts'",
		}, w)
		return
	}

	// validation for end_s
	endTS, err := strconv.Atoi(endTs)
	if err != nil {
		rGetAlerts(http.StatusBadRequest, &responseGetAlerts{
			Error: "failed to parse 'end_ts'",
		}, w)
		return
	}

	// get alerts from the service layer
	alerts, err := s.service.GetAlerts(serviceID, int64(startTS), int64(endTS))
	if err != nil {
		rGetAlerts(http.StatusInternalServerError, &responseGetAlerts{
			Error: "failed to interact with storage",
		}, w)
		return
	}

	rGetAlerts(http.StatusOK, &responseGetAlerts{
		Alerts: alerts,
	}, w)
}

// handlePostAlert handles POST /alerts.
func (s *server) handlePostAlert(w http.ResponseWriter, r *http.Request) {
	alert := &service.Alert{}
	err := json.NewDecoder(r.Body).Decode(alert)
	if err != nil {
		rPostAlerts(http.StatusBadRequest, &responsePostAlerts{
			AlertID: "",
			Error:   fmt.Errorf("failed to decode request body: %w", err).Error(),
		}, w)
		return
	}

	alertID, err := s.service.InsertAlert(alert)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &responsePostAlerts{
		AlertID: alertID,
	}

	responseJSON, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{
 "alert_id": "",
 "error": "failed to marshal JSON response"
}`))
		return
	}

	write, err := w.Write(responseJSON)
	if err != nil {
		log.Printf("failed to write response (%d): %w", write, err)
		return
	}
}

// rPostAlerts small wrapper that writes response for POST /alerts.
func rPostAlerts(httpStatus int, resp *responsePostAlerts, w http.ResponseWriter) {
	responseJSON, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Printf("failed to marshal response: %w\n", err)
		w.WriteHeader(httpStatus)
		return
	}
	w.WriteHeader(httpStatus)
	write, err := w.Write(responseJSON)
	if err != nil {
		log.Printf("failed to write response (%d): %w", write, err)
	}
}

// rGetAlerts small wrapper that writes response for GET /alerts.
func rGetAlerts(httpStatus int, resp *responseGetAlerts, w http.ResponseWriter) {
	responseJSON, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Printf("failed to marshal response: %w\n", err)
		w.WriteHeader(httpStatus)
		return
	}
	w.WriteHeader(httpStatus)
	write, err := w.Write(responseJSON)
	if err != nil {
		fmt.Printf("failed to write response (%d): %w", write, err)
	}
}
