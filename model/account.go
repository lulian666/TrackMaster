package model

import (
	"TrackMaster/pkg"
	"gorm.io/gorm"
	"time"
)

type Account struct {
	ProjectID   string    `gorm:"index;not null" json:"projectID" binding:"required,max=32"`
	ID          string    `gorm:"primaryKey" json:"id" binding:"required,max=32"`
	Description string    `gorm:"not null;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci" json:"description" binding:"required,min=2,max=50"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (a *Account) Create(db *gorm.DB) *pkg.Error {
	result := db.Create(a)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error())
	}
	return nil
}

func (a *Account) Get(db *gorm.DB) *pkg.Error {
	result := db.First(&a)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error())
	}
	return nil
}

func (a *Account) List(db *gorm.DB, project *Project, pageOffset, pageSize int) ([]Account, int64, *pkg.Error) {
	var accounts []Account
	result := db.Where("project_id = ?", project.ID).Find(&accounts)
	if result.Error != nil {
		return nil, 0, pkg.NewError(pkg.ServerError, result.Error.Error())
	}

	if pageOffset >= 0 && pageSize > 0 {
		result = result.Offset(pageOffset).Limit(pageSize).Find(&accounts)
	}

	if a.Description != "" {
		result = result.Where("description LIKE ?", "%"+a.Description+"%").Find(&accounts)
	}

	totalRow := result.RowsAffected

	return accounts, totalRow, nil
}

func (a *Account) GetSome(db *gorm.DB, accountIDs []string) ([]Account, int64, *pkg.Error) {
	var accounts []Account
	result := db.Model(a).Where("id in (?)", accountIDs).Find(&accounts)

	if result.Error != nil {
		return nil, 0, pkg.NewError(pkg.ServerError, result.Error.Error())
	}

	totalRow := result.RowsAffected

	return accounts, totalRow, nil
}

func (a *Account) Delete(db *gorm.DB) *pkg.Error {
	result := db.Delete(&a)
	if result.Error != nil {
		return pkg.NewError(pkg.ServerError, result.Error.Error())
	}
	return nil
}

type Accounts struct {
	Data  []*Account
	Pager *pkg.Pager
}
