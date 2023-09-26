package telegram

type UpdatesResponse struct {
	Ok       bool     `json:"ok"`
	Response []Update `json:"response"`
}

type Update struct {
	ID      int    `json:"update_id"`
	Message string `json:"message"`
}
