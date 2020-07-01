package tools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/tianye/websocket_gateway/common/tools/uuid"
	"net"
	"os"
	"strconv"
	"time"
)

//生成一个LOG ID
func GetLogId() (logId string) {
	uuId := uuid.NewV4()
	logId = uuId.String()

	if logId == "" {
		logId = strconv.Itoa(os.Getpid()) + "_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	return logId
}

//获取当前时间戳
func GetNowTimeUnix() int64 {
	now := time.Now().Unix() //获取时间戳

	return now
}

//ip2lang
func Ip2long(ip net.IP) uint32 {

	a := uint32(ip[12])
	b := uint32(ip[13])
	c := uint32(ip[14])
	d := uint32(ip[15])

	return uint32(a<<24 | b<<16 | c<<8 | d)
}

//lang2ip
func Long2ip(ip int64) net.IP {

	a := byte((ip >> 24) & 0xFF)
	b := byte((ip >> 16) & 0xFF)
	c := byte((ip >> 8) & 0xFF)
	d := byte(ip & 0xFF)

	return net.IPv4(a, b, c, d)
}

//生成MD5
func BuildMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//time33算法
func Time33(str string) int64 {
	hash := int64(5381) // 001 010 100 000 101 ,hash后的分布更好一些
	s := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	for i := 0; i < len(s); i++ {
		hash = hash*33 + int64(s[i])
	}
	return hash
}
