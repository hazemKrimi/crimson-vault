package models

import (
	"strings"

	"github.com/google/uuid"
	"github.com/hazemKrimi/crimson-vault/internal/lib"
	"github.com/hazemKrimi/crimson-vault/internal/types"
)

func (db *DB) MigrateUsers() {
	db.instance.AutoMigrate(&types.User{})
}

func (db *DB) CreateUser(body types.CreateUserRequestBody) (types.User, error) {
	user := types.User{
		ID:         uuid.New().String(),
		SessionID:  uuid.New().String(),
		Name:       body.Name,
		FiscalCode: body.FiscalCode,
		Address:    body.Address,
		Zip:        body.Zip,
		Country:    body.Country,
		Phone:      body.Phone,
		Email:      body.Email,
	}

	result := db.instance.Create(&user)

	if result.Error != nil {
		return types.User{}, result.Error
	}

	return user, nil
}

func (db *DB) GetUsers() ([]types.User, error) {
	var users []types.User

	result := db.instance.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (db *DB) GetUserById(id uuid.UUID, user *types.User) error {
	result := db.instance.Where("id = ?", id).First(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) GetUserBySessionId(sessionId uuid.UUID, user *types.User) error {
	result := db.instance.Where("session_id = ?", sessionId).First(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) GetUserByUsername(username string, user *types.User) error {
	result := db.instance.Where("username = ?", username).First(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) UpdateUser(id uuid.UUID, body types.UpdateUserRequestBody, user *types.User) error {
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

func (db *DB) UpdateUserSecurityDetails(id uuid.UUID, body types.UpdateUserSecurityDetailsBody, user *types.User) error {
	result := db.instance.Where("id = ?", id).First(user, id)

	if result.Error != nil {
		return result.Error
	}

	hashedPassword, err := lib.HashPassword(body.Password)

	if err != nil {
		return err
	}

	result = db.instance.Model(user).Updates(types.User{
		Username: strings.ToLower(body.Username),
		Password: hashedPassword,
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

func (db *DB) UpdateUserSessionID(user *types.User) error {
	result := db.instance.Model(user).Updates(types.User{
		SessionID: uuid.New().String(),
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) DeleteUser(id uuid.UUID) error {
	result := db.instance.Unscoped().Delete(&types.User{}, id)

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

func (db *DB) DeleteUserSessionID(sessionId string) error {
	result := db.instance.Model(&types.User{}).Where("session_id = ?", sessionId).Update("session_id", "")

	if result.Error != nil {
		return result.Error
	}

	return nil
}
