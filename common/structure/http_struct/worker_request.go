package http_struct

type RequestWorkerAddr struct {
	Addr string `json:"addr"`
}

type RequestWorkerAddrList struct {
	AddrList []string `json:"addr_list"`
}
