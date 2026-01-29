package model

import "time"

type Record struct {
	ID        int64     `json:"id" db:"id"`
	Date      string    `json:"date" db:"date"`
	Content   string    `json:"content" db:"content"`
	Duration  int       `json:"duration" db:"duration"`
	Notes     string    `json:"notes" db:"notes"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateRecordRequest struct {
	Date     string `json:"date" validate:"required"`
	Content  string `json:"content" validate:"required"`
	Duration int    `json:"duration" validate:"required,min=1"`
	Notes    string `json:"notes"`
}

type UpdateRecordRequest struct {
	Date     string `json:"date" validate:"required"`
	Content  string `json:"content" validate:"required"`
	Duration int    `json:"duration" validate:"required,min=1"`
	Notes    string `json:"notes"`
}

type Stats struct {
	TotalRecords  int `json:"totalRecords"`
	TotalDuration int `json:"totalDuration"`
	ThisWeek      int `json:"thisWeek"`
	ThisMonth     int `json:"thisMonth"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
