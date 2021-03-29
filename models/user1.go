package models

import (
	"errors"
	"log"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User the user model
type User1 struct {
	UserId                string    `gorm:"primary_key" json:"user_id" description:"사용자"`
	UserPw                string    `json:"user_pw" description:"비밀번호"`
	UserNm                string    `json:"user_nm" description:"사용자"`
	TelNo                 string    `json:"tel_no" description:"전화 번호"`
	Email                 string    `json:"email" description:"이메일"`
	ShopId                int       `json:"shop_id" description:"가맹점"`
	ShopNm                string    `json:"shop_nm" description:"가맹점 명"`
	Enabled               bool      `json:"enabled" description:"계정 사용 가능"`
	AccountNonLocked      bool      `json:"account_non_locked" description:"계정 잠금 안됨"`
	AccountNonExpired     bool      `json:"account_non_expired" description:"계정 만료 안됨"`
	CredentialsNonExpired bool      `json:"credentials_non_expired" description:"비번 만료 안됨"`
	CreatedAt             time.Time `json:"created_at" description:"등록일시"`
	CreatedBy             string    `json:"created_by" description:"등록자"`
	UpdatedAt             time.Time `json:"updated_at" description:"수정일시"`
	UpdatedBy             string    `json:"updated_by" description:"수정자"`
}

// TableName for gorm
func (User1) TableName() string {
	return "user"
}

// GetFirstByID gets the user by his ID
func (u *User1) GetFirstByID(userId string) error {
	err := DB().Where("user_id=?", userId).First(u).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

// GetFirstByUserId gets the user by his email
func (u *User1) GetFirstByUserId(userId string) error {
	err := DB().Where("user_id=?", userId).First(u).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

// Create a new user
func (u *User1) Create() error {
	db := DB().Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

// Signup a new user
func (u *User1) Signup() error {
	var user User1
	err := user.GetFirstByUserId(u.UserId)

	if err == nil {
		return ErrUserExists
	} else if err != ErrDataNotFound {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.UserPw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// replace the plaintext password with ciphertext password
	u.UserPw = "{bcrypt}" + string(hash)

	return u.Create()
}

// Login a user
func (u *User1) Login(userPw string) error {
	aORb := regexp.MustCompile("}")
	matches := aORb.FindIndex([]byte(u.UserPw))
	// runes := []rune(u.UserPw)
	bcryptPw := string(u.UserPw[matches[1]:])
	log.Printf("user1.Login :: bcryptPw is %v\n", bcryptPw)

	err := bcrypt.CompareHashAndPassword([]byte(bcryptPw), []byte(userPw))
	if err != nil {
		log.Printf("user1.Login :: error message is %v\n", err.Error())
		return err
	}
	return nil
}

// LoginByEmailAndPassword login a user by his email and password
func LoginByEmailAndPassword(userId, userPw string) (*User1, error) {
	var user User1
	err := user.GetFirstByUserId(userId)
	if err != nil {
		return &user, err
	}

	return &user, user.Login(userPw)
}
