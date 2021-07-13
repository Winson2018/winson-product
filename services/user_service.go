package services

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"winson-product/datamodels"
	"winson-product/repositories"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

func NewService(repository repositories.IUserRepository) IUserService{
	return &UserService{repository}
}

type UserService struct{
	UserRepository repositories.IUserRepository
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool){
	user, err := u.UserRepository.Select(userName)
	if err != nil {
		fmt.Println("IsPwdSuccess.Select error:")
		fmt.Println(err)
		return
	}
	isOk, err = ValidatePassword(pwd, user.HashPassword)

	if !isOk {
		fmt.Println("IsPwdSuccess.ValidatePassword error:")
		fmt.Println(err)
		return &datamodels.User{}, false
	}
	return
}

func ValidatePassword(userPassword string, hashed string) (isOk bool, err error){
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误！")
	}

	return true, nil
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error){
	pwdByte, err := GeneratePassword(user.HashPassword)
	if err != nil {
		return userId, err
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

func GeneratePassword(userPassword string) ([]byte, error){
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}
