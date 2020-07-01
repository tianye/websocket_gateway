package pack

import (
	"fmt"
	"testing"
)

func TestPack(t *testing.T) {
	p := new(Protocol)
	//这个例子 N8:4294967295  N4:2147483647  N2:65535
	p.Format = []string{"N8", "N2", "N4"}
	maxString := p.Pack16(4294967295, 65535, 2147483647)
	MaxUnInt64List := p.UnPack16(maxString)
	fmt.Println(maxString, "---", MaxUnInt64List)

	//这个例子 N8:-4294967295  N4:-2147483647  N2:0
	p.Format = []string{"N8", "N2", "N4"}
	MinString := p.Pack16(-4294967295, 0, -2147483647)
	MinUnInt64List := p.UnPack16(MinString)
	fmt.Println(MinString, "---", MinUnInt64List)
}
