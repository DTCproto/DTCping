package base

import (
	"github.com/google/uuid"
	"os"
)

// V4版本UUID
func GetIdByV4() uuid.UUID {
	id, err := uuid.NewRandom()
	if err != nil {
		id, _ = uuid.NewRandom()
	}
	return id
}

// 进程ID
func GetPId() int {
	return os.Getpid()
}
