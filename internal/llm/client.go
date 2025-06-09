package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"task-manager/internal/config"
	"time"
)

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type LLMClient struct {
	config *config.OpenAIConfig
	client *http.Client
}

func NewLLMClient() *LLMClient {
	cfg := config.GetConfig()
	return &LLMClient{
		config: &cfg.OpenAI,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *LLMClient) GenerateTechnicalPlan(taskTitle, taskDescription string, deadline time.Time) (string, error) {
	if !c.config.UseLLM {
		return c.generateMockTechnicalPlan(taskTitle), nil
	}

	prompt := fmt.Sprintf(`
Please generate a detailed technical plan for the following development task:

Task Title: %s
Task Description: %s
Deadline: %s

Please provide:
1. Technology stack recommendations
2. Architecture overview
3. Implementation phases
4. Key milestones
5. Risk assessment
6. Resource requirements

Format the response in markdown.
`, taskTitle, taskDescription, deadline.Format("2006-01-02 15:04:05"))

	return c.callOpenAI(prompt)
}

func (c *LLMClient) GenerateWorkflow(taskTitle, taskDescription string) (string, error) {
	if !c.config.UseLLM {
		return c.generateMockWorkflow(taskTitle), nil
	}

	prompt := fmt.Sprintf(`
Generate a detailed workflow for the following development task:

Task Title: %s
Task Description: %s

Please provide:
1. Step-by-step development process
2. Dependencies between tasks
3. Estimated time for each step
4. Quality checkpoints
5. Testing strategy

Format the response as a structured workflow in markdown.
`, taskTitle, taskDescription)

	return c.callOpenAI(prompt)
}

func (c *LLMClient) GenerateSubTasks(taskTitle, taskDescription string) ([]SubTaskSuggestion, error) {
	if !c.config.UseLLM {
		return c.generateMockSubTasks(taskTitle), nil
	}

	prompt := fmt.Sprintf(`
Break down the following development task into smaller sub-tasks:

Task Title: %s
Task Description: %s

Please provide a JSON array of sub-tasks with the following structure:
[
  {
    "title": "Sub-task title",
    "description": "Detailed description",
    "estimated_hours": 8,
    "priority": "high|medium|low",
    "order": 1,
    "dependencies": []
  }
]

Only return the JSON array, no additional text.
`, taskTitle, taskDescription)

	response, err := c.callOpenAI(prompt)
	if err != nil {
		return nil, err
	}

	var subTasks []SubTaskSuggestion
	if err := json.Unmarshal([]byte(response), &subTasks); err != nil {
		return c.generateMockSubTasks(taskTitle), nil
	}

	return subTasks, nil
}

func (c *LLMClient) callOpenAI(prompt string) (string, error) {
	request := OpenAIRequest{
		Model: c.config.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.config.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return "", fmt.Errorf("OpenAI API key is invalid or not properly configured. Please check your config.yaml file and add a valid API key")
		} else if resp.StatusCode == http.StatusTooManyRequests {
			return "", fmt.Errorf("OpenAI API rate limit exceeded. Please try again later or check your subscription limits")
		} else if resp.StatusCode == http.StatusBadRequest {
			return "", fmt.Errorf("OpenAI API request is invalid. Please check your input parameters and model availability: %s", string(body))
		}
		return "", fmt.Errorf("OpenAI API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response OpenAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return response.Choices[0].Message.Content, nil
}

type SubTaskSuggestion struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	EstimatedHours int    `json:"estimated_hours"`
	Priority       string `json:"priority"`
	Order          int    `json:"order"`
	Dependencies   []int  `json:"dependencies"`
}

func (c *LLMClient) generateMockTechnicalPlan(taskTitle string) string {
	return fmt.Sprintf(`# Technical Plan for %s

## Technology Stack
- Backend: Go with Gin framework
- Frontend: Vue.js 3
- Database: SQLite/PostgreSQL
- API: RESTful services

## Architecture Overview
- Microservices architecture
- Clean code principles
- MVC pattern implementation

## Implementation Phases
1. Database design and setup
2. Backend API development
3. Frontend component development
4. Integration and testing
5. Deployment and monitoring

## Key Milestones
- Week 1: Backend foundation
- Week 2: Core functionality
- Week 3: Frontend development
- Week 4: Testing and deployment

## Risk Assessment
- Medium complexity
- Standard technology stack
- Well-defined requirements

## Resource Requirements
- 1 Full-stack developer
- 4 weeks development time
- Standard development environment`, taskTitle)
}

func (c *LLMClient) generateMockWorkflow(taskTitle string) string {
	return fmt.Sprintf(`# Workflow for %s

## Development Process
1. **Requirements Analysis** (4 hours)
   - Analyze requirements
   - Create user stories
   
2. **System Design** (8 hours)
   - Database schema design
   - API endpoint design
   - UI/UX wireframes
   
3. **Backend Development** (40 hours)
   - Set up project structure
   - Implement database models
   - Create API endpoints
   - Add authentication
   
4. **Frontend Development** (32 hours)
   - Set up Vue.js project
   - Create components
   - Implement routing
   - Connect to backend API
   
5. **Testing** (16 hours)
   - Unit tests
   - Integration tests
   - End-to-end tests
   
6. **Deployment** (8 hours)
   - Set up production environment
   - Deploy application
   - Configure monitoring

## Quality Checkpoints
- Code review after each major component
- Testing at each phase
- Performance optimization
- Security audit

## Testing Strategy
- Unit tests for business logic
- API testing with Postman
- Frontend component testing
- End-to-end user testing`, taskTitle)
}

func (c *LLMClient) generateMockSubTasks(taskTitle string) []SubTaskSuggestion {
	return []SubTaskSuggestion{
		{
			Title:          "Project Setup and Configuration",
			Description:    "Initialize project structure, set up development environment, configure build tools",
			EstimatedHours: 8,
			Priority:       "high",
			Order:          1,
			Dependencies:   []int{},
		},
		{
			Title:          "Database Design and Implementation",
			Description:    "Design database schema, create models, set up migrations",
			EstimatedHours: 12,
			Priority:       "high",
			Order:          2,
			Dependencies:   []int{1},
		},
		{
			Title:          "Backend API Development",
			Description:    "Implement REST API endpoints, business logic, and data validation",
			EstimatedHours: 24,
			Priority:       "high",
			Order:          3,
			Dependencies:   []int{2},
		},
		{
			Title:          "Frontend Component Development",
			Description:    "Create Vue.js components, implement routing, and state management",
			EstimatedHours: 20,
			Priority:       "medium",
			Order:          4,
			Dependencies:   []int{3},
		},
		{
			Title:          "Integration and Testing",
			Description:    "Connect frontend to backend, implement tests, fix bugs",
			EstimatedHours: 16,
			Priority:       "medium",
			Order:          5,
			Dependencies:   []int{4},
		},
		{
			Title:          "Deployment and Documentation",
			Description:    "Deploy to production, write documentation, create user guides",
			EstimatedHours: 8,
			Priority:       "low",
			Order:          6,
			Dependencies:   []int{5},
		},
	}
}
