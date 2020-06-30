package util

import "sync/atomic"

func Int32Incr(addr int32) int32 {
	return atomic.AddInt32(&addr, 1)
}

func Int32Decr(addr int32) int32 {
	return atomic.AddInt32(&addr, -1)
}
