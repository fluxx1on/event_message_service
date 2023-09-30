package entity

type SendResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SendRequest struct {
	ID    int64  `json:"id"`
	Phone int64  `json:"phone"`
	Text  string `json:"text"`
}
