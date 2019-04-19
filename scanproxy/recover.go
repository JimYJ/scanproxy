package scanproxy

import "log"

// HandelRecover 异常捕获
func HandelRecover() {
	if r := recover(); r != nil {
		log.Println("[Fatal]", r)
	}
}
