package sign

import (
	"fmt"
	"testing"
)

func TestPack(t *testing.T) {
	test := struct {
		Test string `json:"test"`
	}{Test: "你好啊"}

	signString := GetSign(test, "testsecret")

	fmt.Println(signString)
}
