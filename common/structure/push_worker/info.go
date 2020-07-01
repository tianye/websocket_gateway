package push_worker

type EventInfo struct {
	GatewayIp    string `json:"gateway_ip,omitempty"`
	HttpPort     string `json:"http_port,omitempty"`
	SocketPort   string `json:"socket_port,omitempty"`
	ConnectionId string `json:"connection_id,omitempty"`
	EventTime    int64  `json:"event_time,omitempty"`
	EventType    string `json:"event_type,omitempty"`
	EventData    string `json:"event_data,omitempty"`
	DataLength   int64  `json:"data_length,omitempty"`
}

type ResponseInfo struct {
	WorkerIp     string `json:"worker_ip,omitempty"`
	ConnectionId string `json:"connection_id,omitempty"`
	EventTime    int64  `json:"event_time,omitempty"`
	RequestTime  int64  `json:"request_time,omitempty"`
	ResponseCode int64  `json:"response_code,omitempty"`
	ResponseData string `json:"response_data,omitempty"`
	DataLength   int64  `json:"data_length,omitempty"`
}
