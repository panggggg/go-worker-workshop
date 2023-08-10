package repository

import (
	"encoding/json"

	"github.com/wisesight/go-api-template/pkg/adapter"
	"github.com/wisesight/go-api-template/pkg/entity"
)

func GetThreads(keyword string) (*entity.ThreadResponse, error) {
	var threads entity.ThreadResponse
	response := adapter.GetThreads(keyword)
	err := json.Unmarshal([]byte(response), &threads)
	if err != nil {
		return nil, err
	}

	return &threads, nil
}
