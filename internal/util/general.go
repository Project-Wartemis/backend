package util

import (
	"sync"
)

func Includes(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

type SafeCounter struct {
	sync.Mutex
	v int
}

func (this *SafeCounter) GetNext() int {
	this.Lock()
	defer this.Unlock()
	this.v++
	return this.v
}
