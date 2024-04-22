package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicateUser  = errors.New("用户邮箱冲突")
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
		const duplicateErr uint16 = 1062 //唯一索引冲突,user表有email和phone两个唯一索引字段
		if me.Number == duplicateErr {
			return ErrDuplicateUser
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

func (dao *UserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone=?", phone).First(&u).Error

	return u, err
}

func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, id).Error

	return u, err
}

/*
因为email和phone都是唯一索引，而且是只要其中一个有值就满足条件，所以使用sql.NullString类型代替string。
sql.NullString类型实现了接口driver和Scanner连个接口

发散一下，我们可以自己定一个类型，并且实现driver和Scanner两个接口，就可以接管这个字段在数据读取和写入操作
我们就可以实现比如：json字段的序列化和反序列化、敏感字段（身份证ID）的加密和解密等
*/
type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique;size:255;comment:电子邮箱"`
	Phone    sql.NullString `gorm:"unique;size:31;comment:手机号"`
	Password string         `gorm:"size:255;not null;default:'';comment:密码"`
	Nickname string         `gorm:"size:63;not null;default:'';comment:昵称"`
	Birthday int64          `gorm:"not null;default:0;comment:生日"`
	Gender   int8           `gorm:"type:enum('1', '2');default:1;comment:性别：1，男性；2，女性"`
	Profile  string         `gorm:"size:255;not null;default:'';comment:个人简介"`
	CreateDt int64          `gorm:"not null;default:0;comment:创建时间，毫秒时间戳"`
	UpdateDt int64          `gorm:"not null;default:0;comment:更新时间，毫秒时间戳"`
	DeleteDt int64          `gorm:"not null;default:0;comment:删除时间，毫秒时间戳"`
}
