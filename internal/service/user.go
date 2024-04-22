package service

import (
	"context"
	"errors"
	"time"

	"example.com/basic-gin/webook/internal/domain"
	"example.com/basic-gin/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateUser         = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("用户名或者密码不正确")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Signup(ctx context.Context, u domain.User) error {
	//将明文密码加密成密文密码
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	//设置默认「生日」
	if u.Birthday.IsZero() {
		defaultBirday, err := time.Parse("2006-01-02", "1949-10-01")
		if err == nil {
			u.Birthday = defaultBirday
		}
	}

	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrRecordNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	//检查密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword //密码不正确
	}

	return u, nil
}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

func (svc *UserService) Edit(ctx context.Context, u domain.User) error {
	return svc.repo.Update(ctx, u)
}

func (svc *UserService) FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)

	if err != repository.ErrRecordNotFound {
		//系统错误 或者 非「没有记录」错误
		return u, err
	}

	//没有找到用户，需要「新建」
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})

	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}

	//创建成功，或者创建时存在唯一索引冲突
	return svc.repo.FindByPhone(ctx, phone)

}
