package service

import (
	"TrackMaster/model"
	"TrackMaster/pkg"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type AccountService interface {
	CreateAccount(account *model.Account) *pkg.Error
	ListAccount(project *model.Project, account *model.Account, pager pkg.Pager) ([]model.Account, int64, *pkg.Error)
	DeleteAccount(account *model.Account) *pkg.Error
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
func (s accountService) accountExists(account *model.Account) (bool, *pkg.Error) {
	err := account.Get(s.db)
	if err != nil && !strings.Contains(err.Msg, "record not found") {
		return false, err
	}
	return err == nil, nil
}

// CreateAccount
// 如果创建成功或者已存在账号，返回nil
// 创建失败返回error
func (s accountService) CreateAccount(account *model.Account) *pkg.Error {
	exists, err := s.accountExists(account)
	if err != nil {
		return err
	}
	if exists {
		return pkg.NewError(pkg.BadRequest, fmt.Sprintf("account with id %s already exist", account.ID))
	} else if !exists {
		if err := account.Create(s.db); err != nil {
			return err
		}
	}
	return nil
}

func (s accountService) ListAccount(project *model.Project, account *model.Account, pager pkg.Pager) ([]model.Account, int64, *pkg.Error) {
	pageOffset := pkg.GetPageOffset(pager.Page, pager.PageSize)
	return account.List(s.db, project, pageOffset, pager.PageSize)
}

func (s accountService) DeleteAccount(account *model.Account) *pkg.Error {
	return account.Delete(s.db)
}
