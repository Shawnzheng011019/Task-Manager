package models

import (
	"time"
	"gorm.io/gorm"
)

type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
	PriorityUrgent
)

type TaskStatus int

const (
	StatusPending TaskStatus = iota
	StatusInProgress
	StatusCompleted
	StatusCancelled
)

type Task struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Deadline    time.Time      `json:"deadline"`
	Priority    Priority       `json:"priority"`
	Status      TaskStatus     `json:"status" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// LLM generated content
	TechnicalPlan string `json:"technical_plan" gorm:"type:text"`
	Workflow      string `json:"workflow" gorm:"type:text"`
	Documentation string `json:"documentation" gorm:"type:text"`
	
	// Relationships
	SubTasks []SubTask `json:"sub_tasks" gorm:"foreignKey:TaskID"`
}

type SubTask struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TaskID      uint           `json:"task_id"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	EstimatedHours int         `json:"estimated_hours"`
	Priority    Priority       `json:"priority"`
	Status      TaskStatus     `json:"status" gorm:"default:0"`
	Order       int            `json:"order"`
	Dependencies []uint        `json:"dependencies" gorm:"type:json"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type TaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Deadline    time.Time `json:"deadline" binding:"required"`
}

type TaskResponse struct {
	Task          Task      `json:"task"`
	UrgencyScore  float64   `json:"urgency_score"`
	TimeRemaining string    `json:"time_remaining"`
}

func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityMedium:
		return "medium"
	case PriorityHigh:
		return "high"
	case PriorityUrgent:
		return "urgent"
	default:
		return "unknown"
	}
}

func (s TaskStatus) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusInProgress:
		return "in_progress"
	case StatusCompleted:
		return "completed"
	case StatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}
