package sign

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"sort"
	"strconv"
)

// GetSign get the sign info
func GetSign(data interface{}, appSecret string) string {
	md5ctx := md5.New()

	switch v := reflect.ValueOf(data); v.Kind() {
	case reflect.String:
		md5ctx.Write([]byte(v.String() + appSecret))
		return hex.EncodeToString(md5ctx.Sum(nil))
	case reflect.Struct:
		orderStr := StructToMapSing(v.Interface(), appSecret)
		md5ctx.Write([]byte(orderStr))
		return hex.EncodeToString(md5ctx.Sum(nil))
	case reflect.Ptr:
		originType := v.Elem().Type()
		if originType.Kind() != reflect.Struct {
			return ""
		}
		dataType := reflect.TypeOf(data).Elem()
		dataVal := v.Elem()
		orderStr := buildOrderStr(dataType, dataVal, appSecret)
		md5ctx.Write([]byte(orderStr))
		return hex.EncodeToString(md5ctx.Sum(nil))
	default:
		return ""
	}
}
func buildOrderStr(t reflect.Type, v reflect.Value, appSecret string) (returnStr string) {
	keys := make([]string, 0, t.NumField())

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("json") == "sign" {
			continue
		}
		data[t.Field(i).Tag.Get("json")] = v.Field(i).Interface()

		keys = append(keys, t.Field(i).Tag.Get("json"))
	}

	sort.Sort(sort.StringSlice(keys))

	var buf bytes.Buffer
	for _, k := range keys {
		if data[k] == "" {
			continue
		}
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}

		buf.WriteString(k)
		buf.WriteByte('=')
		switch vv := data[k].(type) {
		case string:
			buf.WriteString(vv)
		case int:
		case int8:
		case int16:
		case int32:
		case int64:
			buf.WriteString(strconv.FormatInt(int64(vv), 10))
		default:
			continue
		}
	}

	buf.WriteString("&sign=" + appSecret)
	returnStr = buf.String()

	return returnStr
}

func StructToMapSing(content interface{}, appSecret string) (returnStr string) {

	t := reflect.TypeOf(content)
	v := reflect.ValueOf(content)

	returnStr = buildOrderStr(t, v, appSecret)

	return returnStr
}
