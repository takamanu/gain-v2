package data

import (
	"context"
	"errors"
	"fmt"
	"gain-v2/features/users"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserData struct {
	db  *gorm.DB
	rds *redis.Client
	ctx context.Context
}

func NewData(db *gorm.DB, redis *redis.Client, ctx context.Context) users.UserDataInterface {
	return &UserData{
		db:  db,
		rds: redis,
		ctx: ctx,
	}
}

func (ud *UserData) AddUser(newData users.User) (*users.User, error) {

	dbData := &User{
		Name:            newData.Name,
		Email:           newData.Email,
		PIN:             newData.PIN,
		Password:        newData.Password,
		Role:            newData.Role,
		DateOfBirth:     newData.DateOfBirth,
		PhoneNumber:     newData.PhoneNumber,
		IsPhoneVerified: true,
		IsEmailVerified: true,
		AccountStatus:   "active",
		LastLoginDate:   time.Now(),
		IP:              newData.IP,
		UserAgent:       newData.UserAgent,
	}

	if err := ud.db.Create(dbData).Error; err != nil {
		return nil, err
	}

	resData := users.User{
		ID:              dbData.ID,
		Name:            dbData.Name,
		Email:           dbData.Email,
		PIN:             dbData.PIN,
		Password:        dbData.Password,
		Role:            dbData.Role,
		DateOfBirth:     dbData.DateOfBirth,
		PhoneNumber:     dbData.PhoneNumber,
		IsPhoneVerified: dbData.IsPhoneVerified,
		IsEmailVerified: dbData.IsEmailVerified,
		AccountStatus:   dbData.AccountStatus,
		LastLoginDate:   dbData.LastLoginDate,
		IP:              dbData.IP,
		UserAgent:       dbData.UserAgent,
	}

	return &resData, nil
}

func (ud *UserData) LoginAdmin(email, password, ip, userAgent string) (*users.User, error) {
	var dbData = new(User)
	var dataCount int64

	dbData.Email = email

	var qry = ud.db.Model(&User{}).
		Where("email = ? AND account_status = ?", dbData.Email, "active").
		Where("deleted_at IS NULL").
		First(dbData)

	err := qry.Count(&dataCount).Error

	if err != nil {
		logrus.Info("DB Error : ", err.Error())
		return nil, err
	}

	if dataCount == 0 {
		return nil, errors.New("Credentials not found")
	}

	passwordBytes := []byte(password)

	err = bcrypt.CompareHashAndPassword([]byte(dbData.Password), passwordBytes)
	if err != nil {
		return nil, errors.New("Incorrect Password")
	}

	_, err = ud.UpdateProfile(int(dbData.ID), users.UpdateProfile{
		IP:            ip,
		UserAgent:     userAgent,
		LastLoginDate: time.Now(),
	})

	if err != nil {
		return nil, errors.New("error updating data")
	}

	result := users.User{
		Name:            dbData.Name,
		Email:           dbData.Email,
		IsPhoneVerified: dbData.IsPhoneVerified,
		IsEmailVerified: dbData.IsEmailVerified,
		AccountStatus:   dbData.AccountStatus,
		LastLoginDate:   time.Now(),
		IP:              ip,
		UserAgent:       userAgent,
	}

	return &result, nil
}

func (ud *UserData) Login(email, password string) (*users.User, error) {
	var dbData = new(User)
	dbData.Email = email

	var qry = ud.db.Where("email = ? AND status = ?", dbData.Email, "active").First(dbData)

	var dataCount int64
	qry.Count(&dataCount)
	if dataCount == 0 {
		return nil, errors.New("Credentials not found")
	}

	if err := qry.Error; err != nil {
		logrus.Info("DB Error : ", err.Error())
		return nil, err
	}

	passwordBytes := []byte(password)

	err := bcrypt.CompareHashAndPassword([]byte(dbData.Password), passwordBytes)
	if err != nil {
		logrus.Info("Incorrect Password")
		return nil, errors.New("Incorrect Password")
	}

	var result = new(users.User)
	result.ID = dbData.ID
	result.Email = dbData.Email
	result.Name = dbData.Name
	result.Role = dbData.Role

	return result, nil
}

func (ud *UserData) LoginCustomer(email, password string) (*users.User, error) {
	var dbData = new(User)
	dbData.Email = email

	var qry = ud.db.Where("email = ? AND status = ?", dbData.Email, "active").First(dbData)

	var dataCount int64
	qry.Count(&dataCount)
	if dataCount == 0 {
		return nil, errors.New("Credentials not found")
	}

	if err := qry.Error; err != nil {
		logrus.Info("DB Error : ", err.Error())
		return nil, err
	}

	passwordBytes := []byte(password)

	err := bcrypt.CompareHashAndPassword([]byte(dbData.Password), passwordBytes)
	if err != nil {
		logrus.Info("Incorrect Password")
		return nil, errors.New("Incorrect Password")
	}

	var result = new(users.User)
	result.ID = dbData.ID
	result.Email = dbData.Email
	result.Name = dbData.Name
	result.Role = dbData.Role

	return result, nil
}

func (ud *UserData) GetByID(id int) (users.User, error) {
	var listUser users.User
	var qry = ud.db.Table("users").Select("users.*").
		Where("users.id = ?", id).
		Where("users.deleted_at is null").
		Scan(&listUser)

	if err := qry.Error; err != nil {
		return listUser, err
	}
	return listUser, nil
}

func (ud *UserData) GetByEmail(email string) (*users.User, error) {
	// var dbData = new(User)
	// dbData.Email = email
	var user users.User

	if err := ud.db.Table("users").Where("email = ?", email).First(&user).Error; err != nil {
		fmt.Println("This is the error:", err)
		return nil, err
	}

	var result = new(users.User)
	result.ID = user.ID
	result.Email = user.Email
	result.Name = user.Name
	result.Role = user.Role

	return result, nil
}

func (ud *UserData) ResetPassword(email string, password string) error {
	if err := ud.db.Table("users").Where("email = ?", email).Update("password", password).Error; err != nil {
		return err
	}

	return nil
}

func (ud *UserData) UpdateProfile(id int, newData users.UpdateProfile) (bool, error) {
	var qry = ud.db.Table("users").Where("id = ?", id).Updates(User{
		Name:          newData.Name,
		Email:         newData.Email,
		Password:      newData.Password,
		IP:            newData.IP,
		UserAgent:     newData.UserAgent,
		LastLoginDate: newData.LastLoginDate,
	})

	if err := qry.Error; err != nil {
		return false, err
	}

	if dataCount := qry.RowsAffected; dataCount < 1 {
		return false, nil
	}

	return true, nil
}

func (ud *UserData) InsertOTPCode(email, otpCode, otpType string) (bool, error) {
	// otp := helper.RandomInteger(6)
	key := fmt.Sprintf("%s:otp:%s", email, otpType)
	ttl, err := ud.rds.TTL(ud.ctx, key).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}

	if ttl > 0 {
		minutes := int(ttl.Minutes())
		seconds := int(ttl.Seconds()) % 60
		if minutes == 0 {
			return false, fmt.Errorf("you must wait %d seconds to request again", seconds)
		}
		return false, fmt.Errorf("you must wait %d minutes and %d seconds to request again", minutes, seconds)
	}

	err = ud.rds.Set(ud.ctx, key, otpCode, 5*time.Minute).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (ud *UserData) VerifyOTPCode(email, otpCode, otpType string) (bool, error) {
	key := fmt.Sprintf("%s:otp:%s", email, otpType)

	storedOTP, err := ud.rds.Get(ud.ctx, key).Result()

	if err == redis.Nil {
		log.Println("OTP not found or expired")
		return false, err
	}

	if err != nil {
		log.Fatalf("Error retrieving OTP: %v", err)
		return false, err
	}

	if storedOTP == otpCode {

		err := ud.rds.Del(ud.ctx, key).Err()

		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, errors.New("invalid otp")
}
