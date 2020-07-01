package client_worker

import (
	"encoding/json"
	"github.com/tianye/websocket_gateway/common/response_code"
	"github.com/tianye/websocket_gateway/common/structure/push_worker"
	"github.com/tianye/websocket_gateway/common/tools"
	"github.com/tianye/websocket_gateway/common/tools/http_request"
	"github.com/tianye/websocket_gateway/conf"
	"strings"
)

func EventOnline(request *push_worker.EventInfo) (responseInfo *push_worker.ResponseInfo, err error) {
	jsonRequest, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	responseJson, err := http_request.SendJson(conf.GetConf(conf.CallbackUrl), jsonRequest)

	if err != nil {
		return nil, err
	}

	responseInfo = &push_worker.ResponseInfo{
		WorkerIp:     request.GatewayIp,
		ConnectionId: request.ConnectionId,
		EventTime:    request.EventTime,
		RequestTime:  tools.GetNowTimeUnix(),
		ResponseCode: response_code.SUCCESS_NULL_CODE,
		ResponseData: response_code.ResponseMessage[response_code.SUCCESS_NULL_CODE],
		DataLength:   int64(strings.Count(response_code.ResponseMessage[response_code.SUCCESS_NULL_CODE], "")),
	}

	json.Unmarshal(responseJson, responseInfo)

	return responseInfo, nil
}

func EventOffline(request *push_worker.EventInfo) (responseInfo *push_worker.ResponseInfo, err error) {
	jsonRequest, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	responseJson, err := http_request.SendJson(conf.GetConf(conf.CallbackUrl), jsonRequest)

	if err != nil {
		return nil, err
	}

	responseInfo = &push_worker.ResponseInfo{
		WorkerIp:     request.GatewayIp,
		ConnectionId: request.ConnectionId,
		EventTime:    request.EventTime,
		RequestTime:  tools.GetNowTimeUnix(),
		ResponseCode: response_code.SUCCESS_NULL_CODE,
		ResponseData: response_code.ResponseMessage[response_code.SUCCESS_NULL_CODE],
		DataLength:   int64(strings.Count(response_code.ResponseMessage[response_code.SUCCESS_NULL_CODE], "")),
	}

	json.Unmarshal(responseJson, responseInfo)

	return responseInfo, nil
}

func EventMessage(request *push_worker.EventInfo) (responseInfo *push_worker.ResponseInfo, err error) {
	jsonRequest, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	responseJson, err := http_request.SendJson(conf.GetConf(conf.CallbackUrl), jsonRequest)

	if err != nil {
		return nil, err
	}

	responseInfo = &push_worker.ResponseInfo{
		WorkerIp:     request.GatewayIp,
		ConnectionId: request.ConnectionId,
		EventTime:    request.EventTime,
		RequestTime:  tools.GetNowTimeUnix(),
		ResponseCode: response_code.SUCCESS_NULL_CODE,
		ResponseData: response_code.ResponseMessage[response_code.SUCCESS_NULL_CODE],
		DataLength:   int64(strings.Count(response_code.ResponseMessage[response_code.SUCCESS_NULL_CODE], "")),
	}

	json.Unmarshal(responseJson, responseInfo)

	return responseInfo, nil
}

func ProcessMessage(request *push_worker.EventInfo) (responseInfo *push_worker.ResponseInfo, err error) {
	jsonRequest, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	responseJson, err := http_request.SendJson(conf.GetConf(conf.CallbackUrl), jsonRequest)

	if err != nil {
		return nil, err
	}

	responseInfo = &push_worker.ResponseInfo{
		WorkerIp:     request.GatewayIp,
		ConnectionId: request.ConnectionId,
		EventTime:    request.EventTime,
		RequestTime:  tools.GetNowTimeUnix(),
		ResponseCode: response_code.SUCCESS_NULL_CODE,
		ResponseData: response_code.ResponseMessage[response_code.SUCCESS_NULL_CODE],
		DataLength:   int64(strings.Count(response_code.ResponseMessage[response_code.SUCCESS_NULL_CODE], "")),
	}

	json.Unmarshal(responseJson, responseInfo)

	return responseInfo, nil
}
