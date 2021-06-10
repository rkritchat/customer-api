package repository

import "fmt"

type UserEntity struct {
	ID       int64
	Username string
	Password string
	Email    string
}

type CustomerRepo interface {
	Save(entity UserEntity) error
	FindByUsername(username string) (UserEntity, error)
}

type customerRepo struct {
}

func NewCustomer() CustomerRepo {
	return &customerRepo{
	}
}

func (repo customerRepo) Save(entity UserEntity) error {
	fmt.Printf("start save: %#v\n", entity)
	return nil
}

func (repo customerRepo) FindByUsername(username string) (UserEntity, error) {
	fmt.Printf("start get user by username: %v\n", username)
	return UserEntity{
		ID:       int64(1),
		Username: username,
		Password: "1235",
		Email:    "test@gmail.com",
	}, nil
}
