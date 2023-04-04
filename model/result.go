package model

import "database/sql/driver"

type TestResult string

const (
	SUCCESS   TestResult = "SUCCESS"
	FAIL      TestResult = "FAIL"
	UNCERTAIN TestResult = "UNCERTAIN"
	UNTESTED  TestResult = "UNTESTED"
)

func (r *TestResult) Scan(value interface{}) error {
	*r = TestResult(value.([]byte))
	return nil
}

func (r *TestResult) Value() (driver.Value, error) {
	return string(*r), nil
}

type Result struct {
	IOS     TestResult `sql:"type:ENUM('SUCCESS', 'FAIL', 'UNCERTAIN', 'UNTESTED')" gorm:"default:UNTESTED" json:"ios"`
	Android TestResult `sql:"type:ENUM('SUCCESS', 'FAIL', 'UNCERTAIN', 'UNTESTED')" gorm:"default:UNTESTED" json:"android"`
	Other   TestResult `sql:"type:ENUM('SUCCESS', 'FAIL', 'UNCERTAIN', 'UNTESTED')" gorm:"default:UNTESTED" json:"other"`
}

type EventResult struct {
	Result
	RecordID string `gorm:"not null" json:"recordID"`
	EventID  string `gorm:"not null" json:"eventID"`
	ID       string `gorm:"type:varchar(191);primaryKey" json:"id" binding:"required,max=32"`
}

type FieldResult struct {
	Result
	RecordID string `gorm:"not null" json:"recordID"`
	FieldID  string `gorm:"not null" json:"fieldID"`
	ID       string `gorm:"type:varchar(191);primaryKey" json:"id" binding:"required,max=32"`
}
