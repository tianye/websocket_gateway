package http_struct

type RequestConnectionId struct {
	ConnectionId string `json:"connection_id"`
}

type RequestConnectionMessage struct {
	ConnectionId string `json:"connection_id"`
	PushMessage  string `json:"push_message"`
}
