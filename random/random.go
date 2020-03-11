package random

import (
	"github.com/google/uuid"
)

func UUID() (id string, err error) {
	newId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return newId.String(), nil
}