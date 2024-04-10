package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicateEmail = errors.New("用户邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)
var (
	GenderMenuMale   = int8(1)
	GenderMenuFemale = int8(2)
	GenderMenuMap    = map[int8]string{
		GenderMenuMale:   "男",
		GenderMenuFemale: "女",
	}
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreateDt = now
	u.UpdateDt = now

	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062 //唯一索引冲突
		if me.Number == duplicateErr {
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) Update(ctx context.Context, u User) error {
	u.UpdateDt = time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Save(&u).Error
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error

	return u, err
}

func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, id).Error

	return u, err
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique;size:255;comment:电子邮箱"`
	Password string `gorm:"size:255;not null;default:'';comment:密码"`
	Nickname string `gorm:"size:63;not null;default:'';comment:昵称"`
	Birthday int64  `gorm:"not null;default:0;comment:生日"`
	Gender   int8   `gorm:"type:enum('1', '2');default:1;comment:性别：1，男性；2，女性"`
	Phone    string `gorm:"size:31;not null;default:'';comment:手机号"`
	Profile  string `gorm:"size:255;not null;default:'';comment:个人简介"`

	CreateDt int64 `gorm:"not null;default:0;comment:创建时间，毫秒时间戳"`
	UpdateDt int64 `gorm:"not null;default:0;comment:更新时间，毫秒时间戳"`
	DeleteDt int64 `gorm:"not null;default:0;comment:删除时间，毫秒时间戳"`
}
