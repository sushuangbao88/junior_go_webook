package repository

import (
	"context"

	"example.com/basic-gin/webook/internal/domain"
	"example.com/basic-gin/webook/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDAO
}

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrRecordNotFound = dao.ErrRecordNotFound
)

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) Update(ctx context.Context, u domain.User) error {
	du, err := repo.dao.FindById(ctx, u.Id)
	if err != nil {
		return err
	}

	du.Nickname = u.Nickname
	du.Birthday = u.Birthday
	du.Gender = u.Gender
	du.Phone = u.Phone

	return repo.dao.Update(ctx, du)
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
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

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		Gender:   u.Gender,
		Phone:    u.Phone,
	}
}
