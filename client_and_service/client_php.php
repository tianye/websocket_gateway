<?php

#配置项目------------------------------------------------------------------------------------------#
define('SOCKET_IP', '127.0.0.1');        //socket连接的ip
define('SOCKET_PORT', '8301');           //socket端口
define('SOCKET_PATH', '/web_socket');    //socket路径
define('ERROR_LOG_PATH', './client_php_error.log'); //错误日志路径
#配置项目------------------------------------------------------------------------------------------#

#全局异常捕获---------------------------------------------------------------------------------------#
register_shutdown_function(function () {
    $error = error_get_last();
    if ($error) {
        $message = $error['message'] . ' ' . $error['file'] . ' LINE NO: ' . $error['line'] . ' .';

        $error_json = json_encode(['error' => $message]);
        Tools::rewrite_file(ERROR_LOG_PATH, $error_json);

        throw new Exception($error_json);
    }
});

set_exception_handler(function (Exception $exception) {
    /** @var \Exception $exception */

    $error            = [];
    $error['message'] = $exception->getMessage();
    $trace            = $exception->getTrace();
    if ('E' == $trace[0]['function']) {
        $error['file'] = $trace[0]['file'];
        $error['line'] = $trace[0]['line'];
    } else {
        $error['file'] = $exception->getFile();
        $error['line'] = $exception->getLine();
    }
    $error['trace'] = $exception->getTraceAsString();

    $error_json = json_encode(['error' => $error]);
    Tools::rewrite_file(ERROR_LOG_PATH, $error_json);
    throw new Exception($error_json);
});

set_error_handler(function ($err_no, $err_str, $err_file, $err_line) {
    switch ($err_no) {
        case E_ERROR:
        case E_PARSE:
        case E_CORE_ERROR:
        case E_COMPILE_ERROR:
        case E_USER_ERROR:
            ob_end_clean();
            $errorStr = "$err_str " . $err_file . " LINE NO: $err_line .";
            break;
        default:
            $errorStr = "[$err_no] $err_str " . $err_file . " LINE NO: $err_line .";
            break;
    }

    Tools::rewrite_file(ERROR_LOG_PATH, json_encode(['error' => $errorStr]));

    throw new Exception($errorStr);
});

#异常捕获------------------------------------------------------------------------------------------#

#初始化连接------------------------------------------------------------------------------------------#
try {
    $socket = Ws::getInstance()->init();
} catch (Exception $exception) {
    Tools::rewrite_file(ERROR_LOG_PATH, 'connect fail init:' . socket_strerror(socket_last_error()));
}
#初始化连接------------------------------------------------------------------------------------------#

#发送数据------------------------------------------------------------------------------------------#
try {
    Ws::getInstance()->send('你好, 我是php的webSocket客户端');
} catch (Exception $exception) {
    Tools::rewrite_file(ERROR_LOG_PATH, 'send message fail:' . socket_strerror(socket_last_error()));
}
#发送数据------------------------------------------------------------------------------------------#

do {
    try {
        Ws::getInstance()->read();
        while ($listen_info = array_shift(Ws::$response_message)) {
            //接受到的消息
            echo "Listen Message:", $listen_info, PHP_EOL;
        }
    } catch (Exception $exception) {
        Tools::rewrite_file(ERROR_LOG_PATH, 'listen message fail:' . socket_strerror(socket_last_error()));
    }

    usleep(1000);
} while (true);

#工具类------------------------------------------------------------------------------------------#
class Tools
{
    //REWRITE FILE
    public static function rewrite_file($file, $message) {
        if (empty($message) || !is_string($message)) {
            return false;
        }
        $fp = @fopen($file, "w+");
        @fwrite($fp, strval($message));
        @fclose($fp);

        return true;
    }
}

#工具类------------------------------------------------------------------------------------------#

#WebSocket
class Ws
{
    public static $ws;

    public static $response_message = [];

    public static $socket = null;

    public $options = ['fragment_size' => 4096];

    public $opcodes = [
        'continuation' => 0,
        'text'         => 1,
        'binary'       => 2,
        'close'        => 8,
        'ping'         => 9,
        'pong'         => 10,
    ];

    /**
     * @return \Ws
     */
    public static function getInstance() {
        if (!empty(self::$ws) && self::$ws instanceof Ws) {
            return self::$ws;
        }

        self::$ws = new self();

        return self::$ws;
    }

    /**
     * @return resource|null
     * @throws \Exception
     */
    public function init() {
        self::$socket = socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
        socket_set_option(self::$socket, SOL_SOCKET, SO_RCVTIMEO, ["sec" => 1, "usec" => 0]);
        socket_set_option(self::$socket, SOL_SOCKET, SO_SNDTIMEO, ["sec" => 6, "usec" => 0]);
        $result = socket_connect(self::$socket, SOCKET_IP, SOCKET_PORT);

        if ($result < 0) {
            throw new Exception(socket_strerror($result), $result);
        }

        $headers = [
            'host'                  => SOCKET_IP . ":" . SOCKET_PORT,
            'user-agent'            => 'websocket-client-php',
            'connection'            => 'Upgrade',
            'upgrade'               => 'websocket',
            'sec-websocket-key'     => $this->generateKey(),
            'sec-websocket-version' => '13',
        ];

        $header = "GET " . SOCKET_PATH . " HTTP/1.1\r\n"
            . implode(
                "\r\n", array_map(
                    function ($key, $value) {
                        return "$key: $value";
                    }, array_keys($headers), $headers
                )
            )
            . "\r\n\r\n";

        $this->write($header);

        return self::$socket;
    }

    public function generateKey() {
        $chars        = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!"$&/()=[]{}0123456789';
        $key          = '';
        $chars_length = strlen($chars);
        for ($i = 0; $i < 16; $i++) {
            $key .= $chars[mt_rand(0, $chars_length - 1)];
        }

        return base64_encode($key);
    }

    /**
     * @param        $payload
     * @param string $opcode
     * @param bool   $masked
     *
     * @throws \Exception
     */
    public function send($payload, $opcode = 'text', $masked = true) {
        // record the length of the payload
        $payload_length = strlen($payload);

        $fragment_cursor = 0;
        // while we have data to send
        while ($payload_length > $fragment_cursor) {
            // get a fragment of the payload
            $sub_payload = substr($payload, $fragment_cursor, $this->options['fragment_size']);

            // advance the cursor
            $fragment_cursor += $this->options['fragment_size'];

            // is this the final fragment to send?
            $final = $payload_length <= $fragment_cursor;

            // send the fragment
            $this->send_fragment($final, $sub_payload, $opcode, $masked);

            // all fragments after the first will be marked a continuation
            $opcode = 'continuation';
        }
    }

    /**
     * @param $final
     * @param $payload
     * @param $opcode
     * @param $masked
     *
     * @throws \Exception
     */
    public function send_fragment($final, $payload, $opcode, $masked) {
        // Binary string for header.
        $frame_head_binstr = '';

        // Write FIN, final fragment bit.
        $frame_head_binstr .= (bool) $final ? '1' : '0';

        // RSV 1, 2, & 3 false and unused.
        $frame_head_binstr .= '000';

        // Opcode rest of the byte.
        $frame_head_binstr .= sprintf('%04b', $this->opcodes[$opcode]);

        // Use masking?
        $frame_head_binstr .= $masked ? '1' : '0';

        // 7 bits of payload length...
        $payload_length = strlen($payload);
        if ($payload_length > 65535) {
            $frame_head_binstr .= decbin(127);
            $frame_head_binstr .= sprintf('%064b', $payload_length);
        } elseif ($payload_length > 125) {
            $frame_head_binstr .= decbin(126);
            $frame_head_binstr .= sprintf('%016b', $payload_length);
        } else {
            $frame_head_binstr .= sprintf('%07b', $payload_length);
        }

        $frame = '';

        // Write frame head to frame.
        foreach (str_split($frame_head_binstr, 8) as $binstr) {
            $frame .= chr(bindec($binstr));
        }

        $mask = '';
        // Handle masking
        if ($masked) {
            // generate a random mask:
            for ($i = 0; $i < 4; $i++) {
                $mask .= chr(rand(0, 255));
            }
            $frame .= $mask;
        }

        // Append payload to frame:
        for ($i = 0; $i < $payload_length; $i++) {
            $frame .= ($masked === true) ? $payload[$i] ^ $mask[$i % 4] : $payload[$i];
        }

        $this->write($frame);
    }

    /**
     * @param     $header
     * @param int $retry
     *
     * @throws \Exception
     */
    public function write($header, $retry = 0) {
        try {
            socket_write(self::$socket, $header, strlen($header));
        } catch (Exception $exception) {
            if (!$retry <= 0) {
                usleep(1000);
                $this->write($header, 1);
            }

            throw new Exception('fail to retry write:' . socket_strerror(socket_last_error()) . ',exception:' . $exception->getMessage());
        }
        usleep(1000);
    }

    public function read() {
        $listen_message = '';
        //PHP_BINARY_READ
        while ($callback = socket_read(self::$socket, 1024)) {
            $listen_message .= $callback;
        }

        if (!empty($listen_message)) {
            $delimiter   = chr(129);
            $listen_list = explode($delimiter, $listen_message);

            foreach ($listen_list as $listen_info) {
                if (ord($listen_info[0]) == 60) {
                    $listen_info = ltrim($listen_info, chr(60));
                }

                if (!empty($listen_info)) {
                    array_push(self::$response_message, $listen_info);
                }
            }
        }

        return self::$response_message;
    }
}