package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey;column:id;default:(-)"`
	IsActive  bool      `gorm:"column:is_active;default:true"`
	CreatedAt time.Time `gorm:"column:created_at;default:(-)"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:(-)"`
	UserID    int64     `gorm:"column:user_id;not null"`
	FirstName string    `gorm:"column:first_name;not null"`
	LastName  *string   `gorm:"column:last_name"`
	Username  *string   `gorm:"column:username"`
	Score     int       `gorm:"column:score;default:0"`
}

type IncUserScoreParams struct {
	Ctx    context.Context
	UserID int64
	Points int
}

func (User) TableName() string {
	return "users"
}

func (d *Database) CreateUser(ctx context.Context, user *User) error {
	result := d.DB.WithContext(ctx).Omit("ID", "CreatedAt", "UpdatedAt").Create(user)
	return result.Error
}

func (d *Database) GetUserByTelegramID(ctx context.Context, telegramID int64) (*User, error) {
	var user User

	result := d.DB.WithContext(ctx).Where("user_id = ?", telegramID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

func (d *Database) IncrementUserScore(params IncUserScoreParams) (int, error) {
	var user User

	result := d.DB.WithContext(params.Ctx).
		Model(&user).
		Where("user_id = ?", params.UserID).
		Update("score", gorm.Expr("score + ?", params.Points)).
		First(&user)

	if result.Error != nil {
		slog.Error("increment user score ends with error", slog.Any("err", result.Error))
		return 0, result.Error
	} else {
		slog.Info(fmt.Sprintf("increment user %d score by %d", params.UserID, params.Points))
	}

	return user.Score, nil
}

func (d *Database) GetAllUsers(ctx context.Context) ([]User, error) {
	var users []User

	result := d.DB.WithContext(ctx).Order("score DESC").Find(&users)
	if result.Error != nil {
		slog.Error("failed to get all users", slog.Any("err", result.Error))
		return nil, result.Error
	}

	return users, nil
}
