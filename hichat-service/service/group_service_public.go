package service

import (
	"sync"
)

var GroupLock sync.Mutex
