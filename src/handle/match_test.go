package handle

import (
	"fmt"
	"testing"
)

type B struct {
	Num int32
}

func TestMatch(t *testing.T) {
	var b *B
	a := &B{Num: 1}
	b = a
	fmt.Println(b)
	c := b
	fmt.Println(c)
	d := &B{Num: 2}
	b = d
	fmt.Println(c)
	fmt.Println(b)

}
