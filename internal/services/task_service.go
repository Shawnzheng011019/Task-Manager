package services

import (
	"fmt"
	"math"
	"sort"
	"task-manager/internal/database"
	"task-manager/internal/llm"
	"task-manager/internal/models"
	"time"
)

type TaskService struct {
	llmClient *llm.LLMClient
}

func NewTaskService() *TaskService {
	return &TaskService{
		llmClient: llm.NewLLMClient(),
	}
}

func (s *TaskService) CreateTask(req models.TaskRequest) (*models.TaskResponse, error) {
	db := database.GetDB()

	// Calculate priority based on deadline
	priority := s.calculatePriority(req.Deadline)

	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		Priority:    priority,
		Status:      models.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Generate LLM content
	technicalPlan, err := s.llmClient.GenerateTechnicalPlan(req.Title, req.Description, req.Deadline)
	if err != nil {
		return nil, fmt.Errorf("unable to generate technical plan - %w", err)
	}
	task.TechnicalPlan = technicalPlan

	workflow, err := s.llmClient.GenerateWorkflow(req.Title, req.Description)
	if err != nil {
		return nil, fmt.Errorf("unable to generate workflow - %w", err)
	}
	task.Workflow = workflow

	// Save task to database
	if err := db.Create(&task).Error; err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Generate and save sub-tasks
	subTaskSuggestions, err := s.llmClient.GenerateSubTasks(req.Title, req.Description)
	if err != nil {
		return nil, fmt.Errorf("unable to generate sub-tasks - %w", err)
	}

	for _, suggestion := range subTaskSuggestions {
		subTaskPriority := s.parsePriority(suggestion.Priority)
		subTask := models.SubTask{
			TaskID:         task.ID,
			Title:          suggestion.Title,
			Description:    suggestion.Description,
			EstimatedHours: suggestion.EstimatedHours,
			Priority:       subTaskPriority,
			Status:         models.StatusPending,
			Order:          suggestion.Order,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := db.Create(&subTask).Error; err != nil {
			return nil, fmt.Errorf("failed to create sub-task: %w", err)
		}
	}

	// Load task with sub-tasks
	if err := db.Preload("SubTasks").First(&task, task.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load task with sub-tasks: %w", err)
	}

	urgencyScore := s.calculateUrgencyScore(task.Deadline)
	timeRemaining := s.formatTimeRemaining(task.Deadline)

	return &models.TaskResponse{
		Task:          task,
		UrgencyScore:  urgencyScore,
		TimeRemaining: timeRemaining,
	}, nil
}

func (s *TaskService) GetAllTasks() ([]models.TaskResponse, error) {
	db := database.GetDB()
	var tasks []models.Task

	if err := db.Preload("SubTasks").Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	var responses []models.TaskResponse
	for _, task := range tasks {
		urgencyScore := s.calculateUrgencyScore(task.Deadline)
		timeRemaining := s.formatTimeRemaining(task.Deadline)

		responses = append(responses, models.TaskResponse{
			Task:          task,
			UrgencyScore:  urgencyScore,
			TimeRemaining: timeRemaining,
		})
	}

	// Sort by urgency score (highest first)
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].UrgencyScore > responses[j].UrgencyScore
	})

	return responses, nil
}

func (s *TaskService) GetTaskByID(id uint) (*models.TaskResponse, error) {
	db := database.GetDB()
	var task models.Task

	if err := db.Preload("SubTasks").First(&task, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	urgencyScore := s.calculateUrgencyScore(task.Deadline)
	timeRemaining := s.formatTimeRemaining(task.Deadline)

	return &models.TaskResponse{
		Task:          task,
		UrgencyScore:  urgencyScore,
		TimeRemaining: timeRemaining,
	}, nil
}

func (s *TaskService) UpdateTaskStatus(id uint, status models.TaskStatus) error {
	db := database.GetDB()
	
	if err := db.Model(&models.Task{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	return nil
}

func (s *TaskService) DeleteTask(id uint) error {
	db := database.GetDB()
	
	if err := db.Delete(&models.Task{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func (s *TaskService) calculatePriority(deadline time.Time) models.Priority {
	now := time.Now()
	timeUntilDeadline := deadline.Sub(now)

	switch {
	case timeUntilDeadline <= 24*time.Hour:
		return models.PriorityUrgent
	case timeUntilDeadline <= 3*24*time.Hour:
		return models.PriorityHigh
	case timeUntilDeadline <= 7*24*time.Hour:
		return models.PriorityMedium
	default:
		return models.PriorityLow
	}
}

func (s *TaskService) calculateUrgencyScore(deadline time.Time) float64 {
	now := time.Now()
	timeUntilDeadline := deadline.Sub(now)
	
	if timeUntilDeadline <= 0 {
		return 100.0 // Overdue
	}

	hoursUntilDeadline := timeUntilDeadline.Hours()
	
	// Exponential decay function for urgency
	// Score approaches 100 as deadline approaches
	score := 100 * (1 - math.Exp(-10/hoursUntilDeadline))
	
	return math.Min(score, 100.0)
}

func (s *TaskService) formatTimeRemaining(deadline time.Time) string {
	now := time.Now()
	timeUntilDeadline := deadline.Sub(now)

	if timeUntilDeadline <= 0 {
		return "Overdue"
	}

	days := int(timeUntilDeadline.Hours() / 24)
	hours := int(timeUntilDeadline.Hours()) % 24

	if days > 0 {
		return fmt.Sprintf("%d days, %d hours", days, hours)
	}
	return fmt.Sprintf("%d hours", hours)
}

func (s *TaskService) parsePriority(priorityStr string) models.Priority {
	switch priorityStr {
	case "urgent":
		return models.PriorityUrgent
	case "high":
		return models.PriorityHigh
	case "medium":
		return models.PriorityMedium
	case "low":
		return models.PriorityLow
	default:
		return models.PriorityMedium
	}
}
