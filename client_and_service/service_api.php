<?php

/**
 * 请求操作
 *
 * Class WebSocketGateway
 */
class WebSocketGateway
{
    /**
     * @param string $connectionId
     *
     * @return \ConnectionStruct
     * @throws \Exception
     */
    function DecodeConnection($connectionId) {
        /** @var ConnectionStruct $conection */
        $conection = (new Connection($connectionId))->Format(["N8", "N2", "N2", "N2", "N4"])->UnPack16()->DecodeConnection();

        return $conection;
    }

    /**
     * @param        $connectionId
     * @param string $message
     *
     * @return array
     * @throws \Exception
     */
    function pushConnectionMessage($connectionId, $message = '') {
        $request = [
            'connection_id' => $connectionId,
            'push_message'  => $message,
        ];

        $conection = $this->DecodeConnection($connectionId);

        $url = 'http://' . $conection->Ip . ':' . $conection->HttpPort . '/api/push_connection_message';

        $response = Http::Post($url, $request);

        return $response;
    }

    /**
     * @param $connectionId
     *
     * @return array
     * @throws \Exception
     */
    function getConnectionIsOnline($connectionId) {
        $request = [
            'connection_id' => $connectionId,
        ];

        $conection = $this->DecodeConnection($connectionId);

        $url = 'http://' . $conection->Ip . ':' . $conection->HttpPort . '/api/get_connection_is_online';

        $response = Http::Post($url, $request);

        return $response;
    }

    /**
     * @param $connectionId
     *
     * @return array
     * @throws \Exception
     */
    function kickedOutConnection($connectionId) {
        $request = [
            'connection_id' => $connectionId,
        ];

        $conection = $this->DecodeConnection($connectionId);

        $url = 'http://' . $conection->Ip . ':' . $conection->HttpPort . '/api/kicked_out_connection';

        $response = Http::Post($url, $request);

        return $response;
    }

    /**
     * @param $ip
     * @param $port
     *
     * @return array
     * @throws \Exception
     */
    function getOnlineNumByGateway($ip, $port) {
        $request = [];

        $url = 'http://' . $ip . ':' . $port . '/api/get_online_num';

        $response = Http::Post($url, $request);

        return $response;
    }
}

/**
 * 管道
 *
 * Class Connection
 */
class Connection
{
    private static $format       = [];
    private static $unPack16     = [];
    private static $connectionId = '';

    public function __construct($connectionId) {
        self::$connectionId = $connectionId;

        return $this;
    }

    public function Format($format) {
        self::$format = $format;

        return $this;
    }

    /**
     * @return $this
     * @throws \Exception
     */
    public function UnPack16() {
        if (strlen(self::$connectionId) != 36) {
            throw new Exception("非法的管道ID");
        }

        $response = $this->unPack(self::$connectionId);

        self::$unPack16 = $response;

        return $this;
    }

    public function DecodeConnection() {
        $ConnectionStruct               = new ConnectionStruct;
        $ConnectionStruct->Ip           = long2ip(self::$unPack16[0]);
        $ConnectionStruct->HttpPort     = self::$unPack16[1];
        $ConnectionStruct->Pid          = self::$unPack16[2];
        $ConnectionStruct->Range        = self::$unPack16[3];
        $ConnectionStruct->Cid          = self::$unPack16[4];
        $ConnectionStruct->ConnectionId = self::$connectionId;

        return $ConnectionStruct;
    }

    private function unPack($connectionId) {
        $formatLen = count(self::$format);
        $response  = [];
        $data      = $connectionId;

        if ($formatLen > 0) {
            for ($i = 0; $i < $formatLen; $i++) {
                if (self::$format[$i] == "N8") {
                    $response[$i] = $this->bytes8ToInt64(substr($data, 0, 16));
                    $data         = substr($data, 16);
                } elseif (self::$format[$i] == "N4") {
                    $response[$i] = $this->Bytes4ToInt64(substr($data, 0, 8));
                    $data         = substr($data, 8);
                } elseif (self::$format[$i] == "N2") {
                    $response[$i] = $this->Bytes2ToInt64(substr($data, 0, 4));
                    $data         = substr($data, 4);
                }
            }
        }

        return $response;
    }

    private function bytes8ToInt64($data) {
        $bytes = [
            substr($data, 0, 2),
            substr($data, 2, 2),
            substr($data, 4, 2),
            substr($data, 6, 2),
            substr($data, 8, 2),
            substr($data, 10, 2),
            substr($data, 12, 2),
            substr($data, 14, 2),
        ];

        $val = intval(hexdec($bytes[7])) | intval(hexdec($bytes[6])) << 8 | intval(hexdec($bytes[5])) << 16 | intval(hexdec($bytes[4])) << 24 |
            intval(hexdec($bytes[3])) << 32 | intval(hexdec($bytes[2])) << 40 | intval(hexdec($bytes[1])) << 48 | intval(hexdec($bytes[0])) << 56;

        return $val;
    }

    private function Bytes4ToInt64($data) {
        $bytes = [
            substr($data, 0, 2),
            substr($data, 2, 2),
            substr($data, 4, 2),
            substr($data, 6, 2),
        ];

        $val = intval(hexdec($bytes[0])) | intval(hexdec($bytes[1])) << 8 | intval(hexdec($bytes[2])) << 16 | intval(hexdec($bytes[3])) << 24;

        return $val;
    }

    private function Bytes2ToInt64($data) {
        $bytes = [
            substr($data, 0, 2),
            substr($data, 2, 2),
        ];

        $val = intval(hexdec($bytes[0])) | intval(hexdec($bytes[1])) << 8;

        return $val;
    }
}

/**
 * 管道解析结构体
 *
 * Class ConnectionStruct
 */
class ConnectionStruct
{
    public $Ip;
    public $HttpPort;
    public $Pid;
    public $Range;
    public $Cid;
    public $ConnectionId;
}

/**
 * HTTP请求
 *
 * Class Http
 */
class Http
{
    /**
     * @param       $url
     * @param array $json_array
     *
     * @return array
     * @throws \Exception
     */
    public static function Post($url, array $json_array) {
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_CUSTOMREQUEST, 'POST');
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
        curl_setopt($ch, CURLOPT_HTTPHEADER, [
            'Content-Type: application/json; charset=utf-8',
        ]);
        $body = json_encode($json_array);
        curl_setopt($ch, CURLOPT_POST, 1);
        curl_setopt($ch, CURLOPT_POSTFIELDS, $body);
        $response = curl_exec($ch);
        if (!$response) {
            throw new Exception(curl_error($ch), curl_errno($ch));
        }
        $response = ['status' => curl_getinfo($ch, CURLINFO_HTTP_CODE), 'body' => $response];
        curl_close($ch);

        return $response;
    }
}

//给管道推送消息
$response = (new WebSocketGateway())->pushConnectionMessage('000000007f0000016e20d6e59fa301000000', 'hello');
echo "pushConnectionMessage:", " status:", $response['status'], " body:", $response['body'], PHP_EOL;

//获取管道在线状态
$response = (new WebSocketGateway())->getConnectionIsOnline('000000007f0000016e20d6e59fa301000000');
echo "getConnectionIsOnline:", " status:", $response['status'], " body:", $response['body'], PHP_EOL;

//踢掉管道
$response = (new WebSocketGateway())->kickedOutConnection('000000007f0000016e20d6e59fa301000000');
echo "kickedOutConnection:", " status:", $response['status'], " body:", $response['body'], PHP_EOL;

//获取当前gateway的在线数量
$response = (new WebSocketGateway())->getOnlineNumByGateway('127.0.0.1', '8302');
echo "getOnlineNumByGateway:", " status:", $response['status'], " body:", $response['body'], PHP_EOL;
