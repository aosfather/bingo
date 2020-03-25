package service

import (
	"fmt"
	"testing"
)

func TestUnsaftSet_Gets(t *testing.T) {
	a := make(UnsaftSet)
	a.Add("1", "2", "3")
	fmt.Println(a.Gets(1, 2))
}
