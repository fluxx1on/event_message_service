package entity

import (
	"time"
)

// Identic with DB
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

type Mailings []*Mailing

// Client is the User based entity
//
// TimeZone simply is offset about UTC+0 with step equal 15 min = 1/4 hour.
// So if TimeZone value is 12 then timezone for client is UTC+3.
// TimeZone can be negative it's not mistake.
type Client struct {
	ID             int64  `json:"id"`
	PhoneNumber    int64  `json:"phone_number"`
	MobileOperator int    `json:"mobile_operator_code"`
	Tag            string `json:"tag"`
	TimeZone       int    `json:"time_zone"`
}

func (c Client) CheckTimeZone(IntervalStart, IntervalEnd time.Time) bool {
	currentTime := time.Now().Add(15 * time.Minute * time.Duration(c.TimeZone))

	return currentTime.After(IntervalStart) && currentTime.Before(IntervalEnd)
}

type Clients []*Client

type Message struct {
	ID               int64     `json:"id"`
	DateTimeCreation time.Time `json:"date_time_creation"`
	Try              int       `json:"try"`
	DeliveryStatus   bool      `json:"delivery_status"`
	MailingID        int64     `json:"mailing_id"`
	ClientID         int64     `json:"client_id"`
}

type Messages []*Message

// Grouped Data
type MailingStats struct {
	MailingID     int64     `json:"mailing_id"`
	DateTimeStart time.Time `json:"datetime_start"`
	DateTimeEnd   time.Time `json:"datetime_end"`
	Succesed      int       `json:"succesed"` // About Messages DeliveryStatus atribute
	Failed        int       `json:"failed"`   // About Messages DeliveryStatus atribute
}
