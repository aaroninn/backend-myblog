package uuid

import (
	"github.com/satori/go.uuid"
)

func NewV1() string {
	uuid, err := uuid.NewV1()
	if err != nil {
		return ""
	}

	return uuid.String()
}
