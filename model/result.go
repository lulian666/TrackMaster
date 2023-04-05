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
	value, _ := fr.Android.Value()
	v := value.(string)
	if v != "SUCCESS" && v != "FAIL" && v != "UNCERTAIN" && v != "UNTESTED" {
		return errors.New(fmt.Sprintf("field {%s} has invalid result value with Android", fr.FieldID))
	}

	value1, _ := fr.IOS.Value()
	v1 := value1.(string)
	if v1 != "SUCCESS" && v1 != "FAIL" && v1 != "UNCERTAIN" && v1 != "UNTESTED" {
		return errors.New(fmt.Sprintf("field {%s} has invalid result value with IOS", fr.FieldID))
	}

	value2, _ := fr.Other.Value()
	v2 := value2.(string)
	if v2 != "SUCCESS" && v2 != "FAIL" && v2 != "UNCERTAIN" && v2 != "UNTESTED" {
		return errors.New(fmt.Sprintf("field {%s} has invalid result value with Other", fr.FieldID))
	}
	return
}
