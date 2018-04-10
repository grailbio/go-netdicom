package netdicom

import (
	"fmt"
	"sync/atomic"
)

var idSeq int32 = 32 // for generating unique ID

func newUID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, atomic.AddInt32(&idSeq, 1))
}

func doassert(cond bool, values ...interface{}) {
	if !cond {
		var s string
		for _, value := range values {
			s += fmt.Sprintf("%v ", value)
		}
		panic(s)
	}
}
