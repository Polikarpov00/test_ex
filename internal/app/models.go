package app

import (
	"time"
)

// Структура для табоицы Question
type Question struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Text      string    `gorm:"type:text;not null" json:"text"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Answers   []Answer  `gorm:"constraint:OnDelete:CASCADE;" json:"answers,omitempty"`
}

// Структура для таблицы Answer
type Answer struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	QuestionID uint      `gorm:"index;not null" json:"question_id"`
	UserID     string    `gorm:"type:text;not null" json:"user_id"`
	Text       string    `gorm:"type:text;not null" json:"text"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
