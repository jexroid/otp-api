package pkg

import "time"

func Uuid() uint32 {
	return uint32(time.Now().UnixNano() / (1 << 22))
}
