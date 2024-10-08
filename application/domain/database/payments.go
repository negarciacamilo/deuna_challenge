package database

/*
Status defines the payment status. Normally I wouldn't use an ORM and, I'm being honest here, I don't know if it's possible to create an enum with bun, so instead I will just assume
0 - Pending status
1 - Approved
2 - Cancelled
3 - Rejected
4 - Refunded
5 - Reversal

This would be like CREATE TYPE payment_status as ENUM ('approved', 'cancelled', 'rejected', 'pending', 'refunded', 'reversed')
And assigning type payment_status to status
*/

type Payment struct {
	Base
	Amount     float64   `json:"amount" bun:",notnull,type:numeric(12,2)"`
	Status     int       `json:"status" bun:",notnull,type:int,default:0"`
	CustomerID uint64    `bun:",notnull"`
	Customer   *Customer `bun:"rel:belongs-to,join:customer_id=id"`
	MerchantID uint64    `bun:",notnull"`
	Merchant   *Merchant `bun:"rel:belongs-to,join:merchant_id=id"`
	BankID     uint64    `bun:",notnull"`
	Bank       *Bank     `bun:"rel:belongs-to,join:bank_id=id"`
	// This is the response code based on ISO 8583:2023 (https://www.iso.org/standard/79451.html) and is a string cause the 0 padding matters
	Code string `json:"code"`
	// This is the ID that both the bank and the payment platform use to identify the payment
	OperationID *string `json:"operation_id" bun:",nullzero"`
}
