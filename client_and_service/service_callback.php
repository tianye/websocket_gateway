<?php

$requestJson = file_get_contents("php://input");

//gateway进程启动事件
//{"gateway_ip":"127.0.0.1","http_port":"8302","socket_port":"8301","connection_id":"0","event_time":1593577514,"event_type":"EventProcessRun","event_data":"{\"pid\":63005,\"prefix\":\"000000007f0000016e201df6\"}","data_length":50}

//gateway进程开始死亡事件
//{"gateway_ip":"127.0.0.1","http_port":"8302","socket_port":"8301","connection_id":"0","event_time":1593577599,"event_type":"EventProcessStartKill","event_data":"{\"sig\":\"interrupt\",\"pid\":63005,\"prefix\":\"000000007f0000016e201df6\"}","data_length":68}

//gateway进程完全死亡事件
//{"gateway_ip":"127.0.0.1","http_port":"8302","socket_port":"8301","connection_id":"0","event_time":1593577599,"event_type":"EventProcessEndKill","event_data":"{\"pid\":63005,\"prefix\":\"000000007f0000016e201df6\"}","data_length":50}

//管道上线事件
//{"gateway_ip":"127.0.0.1","http_port":"8302","socket_port":"8301","connection_id":"000000007f0000016e201df6bfd701000000","event_time":1593577533,"event_type":"EventOnline","data_length":1}

//管道离线事件
//{"gateway_ip":"127.0.0.1","http_port":"8302","socket_port":"8301","connection_id":"000000007f0000016e201df6bfd701000000","event_time":1593577552,"event_type":"EventOffline","data_length":1}

//接受管道消息事件
//{"gateway_ip":"127.0.0.1","http_port":"8302","socket_port":"8301","connection_id":"000000007f0000016e201df6436302000000","event_time":1593577582,"event_type":"EventMessage","event_data":"hello","data_length":6}

if (empty($requestJson)) {
    $requestJson = $HTTP_RAW_POST_DATA;
}

//请求日志
file_put_contents('./service_callback.log', $requestJson . PHP_EOL, FILE_APPEND);

$requestData = json_decode($requestJson, true);

$request = [
    'gateway_ip'    => (string) (!empty($requestData['gateway_ip']) ? $requestData['gateway_ip'] : ''),       //GATEWAY 内网IP
    'http_port'     => (string) (!empty($requestData['http_port']) ? $requestData['http_port'] : ''),         //http端口
    'socket_port'   => (string) (!empty($requestData['socket_port']) ? $requestData['socket_port'] : ''),     //socket端口
    'connection_id' => (string) (!empty($requestData['connection_id']) ? $requestData['connection_id'] : ''), //管道ID
    'event_time'    => (int) (!empty($requestData['event_time']) ? $requestData['event_time'] : 0),           //事件接受时间
    //事件接受类型:
    //EventOnline:上线事件 | EventOffline:下线事件 | EventMessage: 消息事件
    //EventProcessRun: gateway进程启动 | EventProcessStartKill:进程开始死亡 | EventProcessEndKill:进程已经死亡
    'event_type'    => (string) (!empty($requestData['event_type']) ? $requestData['event_type'] : ''),
    'event_data'    => (string) (!empty($requestData['event_data']) ? $requestData['event_data'] : ''),    //事件数据
    'data_length'   => (int) (!empty($requestData['data_length']) ? $requestData['data_length'] : 0),      //事件数据长度
];

$responseData = '{"message":"hello", "custom_1":"value", "custom_2":"value2"}';

$response = [
    'worker_ip'     => (string) '127.0.0.1',                                                         //本机IP
    'connection_id' => (string) !empty($request['connection_id']) ? $request['connection_id'] : '',  //管道ID
    'event_time'    => (int) !empty($request['event_time']) ? $request['event_time'] : '',           //事件接受时间
    'request_time'  => (int) time(),                        //响应时间
    'response_code' => (int) '200',                         //注: 200:把response_data会推送给管道 | 204:不推送数据给管道 | 401:身份验证失败踢出管道 | 500:处理失败
    'response_data' => (string) $responseData,              //code:200 会把整个response_data推送给管道
    'data_length'   => (int) strlen($responseData),         //response_data 长度
];

//响应回调
echo json_encode($response, JSON_UNESCAPED_UNICODE);