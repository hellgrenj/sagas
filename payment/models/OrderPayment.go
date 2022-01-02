package models

type OrderPayment struct {
	Amount   float64 `bson:"amount"`
	Item     string  `bson:"item"`
	Quantity float64 `bson:"quantity"`
}
