package repository

import (
	"encoding/json"

	"github.com/wisesight/go-api-template/pkg/adapter"
	"github.com/wisesight/go-api-template/pkg/entity"
)

type IAccountRepository interface {
	GetAccountInfo(thread entity.Thread) (*entity.Job, error)
}

type account struct {
	socialAPIAdapter adapter.ISocialAPIAdapter
}

func (a account) GetAccountInfo(thread entity.Thread) (*entity.Job, error) {
	var response entity.AccountInfo
	accountInfo, _ := a.socialAPIAdapter.GetAccountInfo(thread.UserID)
	err := json.Unmarshal([]byte(accountInfo), &response)
	if err != nil {
		return nil, err
	}

	result := entity.Job{
		Thread:  thread,
		Account: response,
	}

	return &result, nil

}
