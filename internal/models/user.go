package models

import (
	"github.com/hazemKrimi/crimson-vault/internal/types"
)

func (db *DB) MigrateUsers() {
	db.instance.AutoMigrate(&types.User{})
}

func (db *DB) CreateUser(body types.CreateUserRequestBody) types.User {
	user := types.User{
		Name:       body.Name,
		FiscalCode: body.FiscalCode,
		Address:    body.Address,
		Zip:        body.Zip,
		Country:    body.Country,
		Phone:      body.Phone,
		Email:      body.Email,
	}

	db.instance.Create(&user)
	return user
}

func (db *DB) GetUsers() ([]types.User, error) {
	var users []types.User

	result := db.instance.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (db *DB) GetUser(id int, user *types.User) error {
	result := db.instance.Where("id = ?", id).First(user, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) UpdateUser(id int, body types.UpdateUserRequestBody, user *types.User) error {
	result := db.instance.Where("id = ?", id).First(user, id)

	if result.Error != nil {
		return result.Error
	}

	result = db.instance.Model(user).Updates(types.User{
		Name:       body.Name,
		FiscalCode: body.FiscalCode,
		Address:    body.Address,
		Zip:        body.Zip,
		Country:    body.Country,
		Phone:      body.Phone,
		Email:      body.Email,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) UpdateUserSecurityDetails(id int, body types.UpdateUserSecurityDetailsBody, user *types.User) error {
	result := db.instance.Where("id = ?", id).First(user, id)

	if result.Error != nil {
		return result.Error
	}

	result = db.instance.Model(user).Updates(types.User{
		Username: body.Username,
		Password: body.Password,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) UpdateUserLogo(path string, user *types.User) error {
	result := db.instance.Model(user).Updates(types.User{
		Logo: path,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) DeleteUser(id int) error {
	result := db.instance.Delete(&types.User{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) DeleteUserLogo(user *types.User) error {
	result := db.instance.Model(user).Updates(&types.User{
		Logo: "",
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
