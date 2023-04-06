package model

import (
	"TrackMaster/pkg"
	"database/sql/driver"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type TestResult string

const (
	SUCCESS   TestResult = "SUCCESS"
	FAIL      TestResult = "FAIL"
	UNCERTAIN TestResult = "UNCERTAIN"
	UNTESTED  TestResult = "UNTESTED"
)

func (r *TestResult) Scan(value interface{}) *pkg.Error {
	*r = TestResult(value.([]byte))
	return nil
}

func (r *TestResult) Value() (driver.Value, *pkg.Error) {
	return string(*r), nil
}

type Result struct {
	IOS     TestResult `sql:"type:ENUM('SUCCESS', 'FAIL', 'UNCERTAIN', 'UNTESTED')" gorm:"default:UNTESTED" json:"ios"`
	Android TestResult `sql:"type:ENUM('SUCCESS', 'FAIL', 'UNCERTAIN', 'UNTESTED')" gorm:"default:UNTESTED" json:"android"`
	Other   TestResult `sql:"type:ENUM('SUCCESS', 'FAIL', 'UNCERTAIN', 'UNTESTED')" gorm:"default:UNTESTED" json:"other"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type EventResult struct {
	Result
	RecordID string `gorm:"index;not null" json:"recordID"`
	EventID  string `gorm:"index;not null" json:"eventID"`
	ID       string `gorm:"type:varchar(191);primaryKey" json:"id" binding:"required,max=32"`
}

type FieldResult struct {
	Result
	RecordID string `gorm:"index;not null" json:"recordID"`
	FieldID  string `gorm:"index;not null" json:"fieldID"`
	ID       string `gorm:"type:varchar(191);primaryKey" json:"id" binding:"required,max=32"`
}

func (fr *FieldResult) BeforeSave(db *gorm.DB) (err error) {
	if fr.Android != SUCCESS && fr.Android != FAIL && fr.Android != UNCERTAIN && fr.Android != UNTESTED {
		return errors.New(fmt.Sprintf("field {%s} has invalid result value with Android", fr.FieldID))
	}

	if fr.IOS != SUCCESS && fr.IOS != FAIL && fr.IOS != UNCERTAIN && fr.IOS != UNTESTED {
		return errors.New(fmt.Sprintf("field {%s} has invalid result value with IOS", fr.FieldID))
	}

	if fr.Other != SUCCESS && fr.Other != FAIL && fr.Other != UNCERTAIN && fr.Other != UNTESTED {
		return errors.New(fmt.Sprintf("field {%s} has invalid result value with Other", fr.FieldID))
	}
	return
}
