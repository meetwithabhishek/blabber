package blabber

type MessageResponse struct {
	Username string
	Message  string
}

type WSMessageType int

const (
	ClientDataMessage = WSMessageType(0)
	ServerDataMessage = WSMessageType(1)
)

type WebSocketMessage struct {
	MessageType WSMessageType
	Message     []byte
}
