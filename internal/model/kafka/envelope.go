package kafka

type Envelope struct {
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
	UserID  *int   `json:"userId"`
}

type KafkaResponse struct {
	CorrelationID string
	Payload       []byte
	Err           error
}
