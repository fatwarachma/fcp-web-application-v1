package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByEmail(email string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
	GetUserTaskCategory() ([]model.UserTaskCategory, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByEmail(email string) (model.User, error) {
	var u model.User
	result := r.db.Where("email = ?", email).First(&u)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return model.User{}, nil
		}
		return model.User{}, result.Error
	}

	return u, nil
}

func (r *userRepository) CreateUser(itsAnUser model.User) (model.User, error) {
	err := r.db.Create(&itsAnUser).Error
	if err != nil {
		return itsAnUser, err
	}
	return itsAnUser, nil
}

func (r *userRepository) GetUserTaskCategory() ([]model.UserTaskCategory, error) {
	var userTask []model.UserTaskCategory
	result := r.db.Table("users").
		Select("users.id, users.fullname AS fullname, users.email AS email, tasks.title AS task, tasks.deadline AS deadline, tasks.priority AS priority, tasks.status AS status, categories.name AS category").
		Joins("JOIN tasks ON tasks.user_id = users.id").
		Joins("JOIN categories ON categories.id = tasks.category_id").
		Scan(&userTask)
	if result.Error != nil {
		return nil, result.Error
	}

	return userTask, nil
}
