package data

import (
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
)

type roles string
type status string

const (
	SuperAdmin roles  = "super_admin"
	Staff      roles  = "staff"
	Agen       roles  = "agent"
	Customer   roles  = "customer"
	Active     status = "active"
	NotActive  status = "not_active"
)

func (r *roles) Scan(value interface{}) error {
	*r = roles(value.([]byte))
	return nil
}

func (r roles) Value() (driver.Value, error) {
	return string(r), nil
}

func (s *status) Scan(value interface{}) error {
	*s = status(value.([]byte))
	return nil
}

func (s status) Value() (driver.Value, error) {
	return string(s), nil
}

type User struct {
	*gorm.Model
	Name            string    `gorm:"column:name;type:varchar(255)"`
	Email           string    `gorm:"column:email;unique;type:varchar(255)"`
	PIN             string    `gorm:"column:pin;type:varchar(255)"`
	Password        string    `gorm:"column:password;type:varchar(255)"`
	Role            string    `gorm:"column:role;type:roles"`
	DateOfBirth     time.Time `gorm:"column:date_of_birth"`
	PhoneNumber     string    `gorm:"column:phone_number;type:varchar(255)"`
	IsPhoneVerified bool      `gorm:"column:is_phone_verified;default:false" json:"is_phone_verified"`
	IsEmailVerified bool      `gorm:"column:is_email_verified;default:false" json:"is_email_verified"`
	AccountStatus   string    `gorm:"column:account_status;type:status"`
	LastLoginDate   time.Time `gorm:"column:last_login_date"`
	IP              string    `gorm:"column:ip"`
	UserAgent       string    `gorm:"column:user_agent"`
}

type DataRow struct {
	Zone              string
	AreaName          string
	TimePeriod        string
	Source            string
	Nb                string
	Sector            string
	Subsector         string
	IndicatorMetadata string
	Unit              string
	StartYear         int
	EndYear           int
	DataValue         float64
}

func (User) TableName() string {
	return "users"
}
