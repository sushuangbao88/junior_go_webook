package repository

import (
	"context"
	"database/sql"
	"time"

	"example.com/basic-gin/webook/internal/domain"
	"example.com/basic-gin/webook/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDAO
}

var (
	ErrDuplicateUser  = dao.ErrDuplicateUser
	ErrRecordNotFound = dao.ErrRecordNotFound
)

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.toDao(u))
}

func (repo *UserRepository) Update(ctx context.Context, u domain.User) error {
	du, err := repo.dao.FindById(ctx, u.Id)
	if err != nil {
		return err
	}

	du.Nickname = u.Nickname
	du.Birthday = u.Birthday.Unix()
	du.Gender = u.Gender
	du.Profile = u.Profile

	return repo.dao.Update(ctx, du)
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return repo.toDomain(u), nil
}

func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}

	return repo.toDomain(u), nil
}

func (repo *UserRepository) FindById(ctx context.Context, Id int64) (domain.User, error) {
	u, err := repo.dao.FindById(ctx, Id)
	if err != nil {
		return domain.User{}, err
	}

	return repo.toDomain(u), nil
}

// 结构体在domain和dao之前切换：dao->domain
func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: time.Unix(u.Birthday, 0),
		Gender:   u.Gender,
		Phone:    u.Phone.String,
		Profile:  u.Profile,
	}
}

// 结构体在domain和dao之前切换：domain->dao
func (repo *UserRepository) toDao(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: u.Birthday.Unix(),
		Gender:   u.Gender,
		Profile:  u.Profile,
	}
}
