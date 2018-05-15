package service

import (
	"testing"
	"fmt"
)

func TestSearchEngine_BeginSafeFlush(t *testing.T) {
    notTmpkeys:=[]string{"1","2"}
	diffkeys:=make([]string,len(notTmpkeys)+1)
	diffkeys[0]="5"
	copy(diffkeys[1:len(diffkeys)],notTmpkeys)
	fmt.Println(diffkeys)
	fmt.Println(diffkeys[1:len(diffkeys)])
}
