package outbound

type PaymentEvent struct {
	CorrelationId string  `json:"correlationId"`
	Name          string  `json:"name"`
	MessageId     string  `json:"messageId"`
	OrderId       float64 `json:"orderId"`
}
