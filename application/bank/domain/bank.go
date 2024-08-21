package domain

type CardIsValidRequest struct {
	CardHash string `json:"card_hash"`
}

type ClientHasEnoughBalanceRequest struct {
	ClientID uint64  `json:"client_id"`
	Amount   float64 `json:"amount"`
}
