package service

import (
	"TrackMaster/model"
	"TrackMaster/pkg"
	"errors"
	"gorm.io/gorm"
)

type AccountService interface {
	CreateAccount(account *model.Account) error
	ListAccount(account *model.Account, pager pkg.Pager) ([]model.Account, int64, error)
	DeleteAccount(account *model.Account) error
}

type accountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) AccountService {
	return &accountService{
		db: db,
	}
}

// 检查是否已存在账号
// 返回true表示已存在，false表示不存在
func (s accountService) accountExists(account *model.Account) (bool, error) {
	err := account.Get(s.db)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return err == nil, nil
}

// CreateAccount
// 如果创建成功或者已存在账号，返回nil
// 创建失败返回error
func (s accountService) CreateAccount(account *model.Account) error {
	if exists, err := s.accountExists(account); err != nil {
		return err
	} else if !exists {
		if err := account.Create(s.db); err != nil {
			return err
		}
	}
	return nil
}

func (s accountService) ListAccount(account *model.Account, pager pkg.Pager) ([]model.Account, int64, error) {
	pageOffset := pkg.GetPageOffset(pager.Page, pager.PageSize)
	return account.List(s.db, pageOffset, pager.PageSize)
}

func (s accountService) DeleteAccount(account *model.Account) error {
	return account.Delete(s.db)
}
