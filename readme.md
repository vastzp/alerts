How to test it the manual mode:

# Put several alerts to the server with id my_test_service_id

curl -X POST -v -H "Content-Type: application/json" -d '{
"alert_id": "b950482e9911ec7e41f7ca5e5d9a424f",
"service_id": "my_test_service_id",
"service_name": "my_test_service",
"model": "my_test_model",
"alert_type": "anomaly",
"alert_ts": "1695644160",
"severity": "warning",
"team_slack": "slack_ch"
}' http://localhost:8088/alerts

curl -X POST -v -H "Content-Type: application/json" -d '{
"alert_id": "c950482e9911ec7e41f7ca5e5d9a424f",
"service_id": "my_test_service_id",
"service_name": "my_test_service",
"model": "my_test_model",
"alert_type": "anomaly",
"alert_ts": "1695644161",
"severity": "warning",
"team_slack": "slack_ch"
}' http://localhost:8088/alerts


# Get events for the service with ID my_test_service_id
curl "http://localhost:8088/alerts?service=my_test_service_id&start_ts=1695644159&end_ts=1695644161"
