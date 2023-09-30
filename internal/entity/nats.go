package entity

type MailingWithClients struct {
	Mailing *Mailing `json:"mailing"`
	Clients Clients  `json:"clients"`
	Try     int      `json:"try"`
}
