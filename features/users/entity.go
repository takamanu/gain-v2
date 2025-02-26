package users

import (
	"time"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID              uint      `json:"id" gorm:"column:id"`
	Name            string    `json:"name" gorm:"column:name"`
	Email           string    `json:"email" gorm:"column:email"`
	PIN             string    `json:"pin" gorm:"column:pin"`
	Password        string    `json:"password" gorm:"column:password"`
	Role            string    `json:"role" gorm:"column:role"`
	DateOfBirth     time.Time `json:"date_of_birth" gorm:"column:date_of_birth"`
	PhoneNumber     string    `json:"phone_number" gorm:"column:phone_number"`
	IsPhoneVerified bool      `json:"is_phone_verified" gorm:"column:is_phone_verified"`
	IsEmailVerified bool      `json:"is_email_verified" gorm:"column:is_email_verified"`
	AccountStatus   string    `json:"account_status" gorm:"column:account_status"`
	LastLoginDate   time.Time `json:"last_login_date" gorm:"column:last_login_date"`
	IP              string    `json:"ip" gorm:"column:ip"`
	UserAgent       string    `json:"user_agent" gorm:"column:user_agent"`
}

type Admin struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Role           string `json:"role"`
	PhoneNumber    string `json:"phone_number"`
	TokenResetPass string `json:"token_reset_pass"`
	Status         string `json:"status"`
}

type UserInfo struct {
	Name   string         `json:"name"`
	Email  string         `json:"email"`
	Access map[string]any `json:"token"`
}

type UserCredential struct {
	Name   string         `json:"name"`
	Email  string         `json:"email"`
	Access map[string]any `json:"token"`
}

type UpdateAdmin struct {
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Role        string    `json:"role"`
	PhoneNumber string    `json:"phone_number"`
	Status      string    `json:"status"`
}

type UserResetPass struct {
	Email     string
	Code      string
	OTPType   string
	ExpiresAt time.Time
}

type UpdateProfile struct {
	Name          string
	Email         string
	Password      string
	LastLoginDate time.Time `json:"last_login_date" gorm:"column:last_login_date"`
	IP            string    `json:"ip" gorm:"column:ip"`
	UserAgent     string    `json:"user_agent" gorm:"column:user_agent"`
}

type UserHandlerInterface interface {
	AddUser() echo.HandlerFunc
	LoginAdmin() echo.HandlerFunc
	Register() echo.HandlerFunc
	RegisterCustomer() echo.HandlerFunc
	Login() echo.HandlerFunc
	LoginCustomer() echo.HandlerFunc
	ForgotPassword() echo.HandlerFunc
	ResetPassword() echo.HandlerFunc
	UpdateProfile() echo.HandlerFunc
	RefreshToken() echo.HandlerFunc
	GetProfile() echo.HandlerFunc
}

type UserServiceInterface interface {
	AddUser(newData User) (*User, error)
	LoginAdmin(email, password, ip, userAgent string) (*UserCredential, error)
	Register(newData User) (*User, error)
	RegisterCustomer(newData User) (*User, error)
	Login(email string, password string) (*UserCredential, error)
	LoginCustomer(email string, password string) (*UserCredential, error)
	GenerateJwt(email string) (*UserCredential, error)
	ForgotPassword(email string) error
	ResetPassword(code, email string, password string) error
	UpdateProfile(id int, newData UpdateProfile) (bool, error)
	GetProfile(id int) (*User, error)
}

type UserDataInterface interface {
	AddUser(newData User) (*User, error)
	LoginAdmin(email, password, ip, userAgent string) (*User, error)
	Login(email string, password string) (*User, error)
	LoginCustomer(email string, password string) (*User, error)
	GetByID(id int) (User, error)
	GetByEmail(email string) (*User, error)
	ResetPassword(email, password string) error
	UpdateProfile(id int, newData UpdateProfile) (bool, error)
	InsertOTPCode(email, otpCode, otpType string) (bool, error)
	VerifyOTPCode(email, otpCode, otpType string) (bool, error)
}
