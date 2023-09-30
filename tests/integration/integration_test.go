package integration_test

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/Eun/go-hit"
)

type TestJson map[string]interface{}

const (
	app_host = "app:8080"
)

const (
	// Attempts connection
	healthPath = "http://" + app_host + "/healthz"
	attempts   = 5

	// HTTP REST
	basePath = "http://" + app_host + "/v1"
)

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", app_host, err)
	}

	log.Printf("Integration tests: host %s is available", app_host)

	code := m.Run()
	os.Exit(code)
}

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return err
}

func TestHTTP(t *testing.T) {

	// Client
	// PUT
	body := `{
		"mobile_operator_code": 900,
		"phone_number": 78819001121,
		"tag": "silver",
		"time_zone": 12
	}`
	Test(t,
		Description("Client creation succeeded"),
		Put(basePath+"/client"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusCreated),
	)

	// PATCH
	body = `{
		"id": 1,
		"mobile_operator_code": 701
	}`
	Test(t,
		Description("Client updating succeeded"),
		Method(http.MethodPatch, basePath+"/client"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusNoContent),
	)

	// DELETE
	body = `{
		"id": 1
	}`
	Test(t,
		Description("Client deletion succeeded"),
		Delete(basePath+"/client"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusNoContent),
	)

	type Mailing struct {
		ID             int64     `json:"id"`
		MessageText    string    `json:"message_text"`
		MobileOperator string    `json:"mobile_operator_code"`
		Tag            string    `json:"tag"`
		FilterChoice   string    `json:"filter_choice"`
		DateTimeStart  time.Time `json:"datetime_start"`
		DateTimeEnd    time.Time `json:"datetime_end"`
		IntervalStart  time.Time `json:"interval_start"`
		IntervalEnd    time.Time `json:"interval_end"`
	}

	// Mailing
	// PUT
	body = `{
		"message_text": "string",
		"mobile_operator_code": 900,
		"tag": "silver",
		"filter_choice": "tag",
		"DateTimeStart": "2023-09-23 01:51:58.466813165 +0300 MSK m=+10800.001254021",
		"DateTimeEnd": "2023-09-23 01:51:58.466813165 +0300 MSK m=+10800.001254021",
		"IntervalStart": "2023-09-23 01:51:58.466813165 +0300 MSK m=+10800.001254021",
		"IntervalEnd": "2023-09-23 01:51:58.466813165 +0300 MSK m=+10800.001254021",
	}`
	Test(t,
		Description("Mailing creation succeeded"),
		Put(basePath+"/mailing"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusCreated),
	)

	// PATCH
	body = `{
		"id": 1,
		"mobile_operator_code": 701
	}`
	Test(t,
		Description("Mailing updating succeeded"),
		Method(http.MethodPatch, basePath+"/mailing"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusNoContent),
	)

	// DELETE
	body = `{
		"id": 1
	}`
	Test(t,
		Description("Mailing deletion succeeded"),
		Delete(basePath+"/mailing"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusNoContent),
	)
}
