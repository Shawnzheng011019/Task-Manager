# Task Manager

An intelligent task management system built with Go + Vue3, integrated with LLM API for automatic generation of technical solutions and workflows.

## 🚀 Features

### Core Functionality
- **Task Management**: Create, view, update, and delete tasks
- **Priority Management**: Support for Low, Medium, High, and Urgent priority levels
- **Status Tracking**: Manage task states including Pending, In Progress, Completed, and Cancelled
- **Deadline Management**: Task deadline tracking with urgency score calculation
- **Sub-tasks**: Support task breakdown into sub-tasks with dependency management

### AI-Powered Features
- **Technical Plan Generation**: Automatically generate detailed technical implementation plans based on task descriptions
- **Workflow Generation**: Create structured development workflows automatically
- **Task Decomposition**: Intelligently break down complex tasks into executable sub-tasks
- **Time Estimation**: Automatically estimate work hours for each sub-task

### Technology Stack
- **Backend**: Go 1.20 + Gin Framework
- **Frontend**: Vue 3 + Element Plus + Pinia
- **Database**: SQLite (PostgreSQL supported)
- **AI Integration**: OpenAI API (supports custom API endpoints)
- **Build Tool**: Vite

## 📁 Project Structure

```
Task-Manager/
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── configs/
│   └── config.yaml         # Configuration file
├── internal/
│   ├── config/             # Configuration management
│   ├── database/           # Database connection
│   ├── handlers/           # HTTP handlers
│   ├── llm/               # LLM client
│   ├── models/            # Data models
│   └── services/          # Business logic
└── frontend/
    ├── src/
    │   ├── components/    # Vue components
    │   ├── views/         # Page views
    │   ├── stores/        # Pinia state management
    │   └── router/        # Route configuration
    └── package.json       # Frontend dependencies
```

## 🛠️ Installation and Setup

### Prerequisites
- Go 1.20+
- Node.js 16+
- npm or yarn

### Backend Setup

1. Clone the repository
```bash
git clone <repository-url>
cd Task-Manager
```

2. Install Go dependencies
```bash
go mod tidy
```

3. Configure OpenAI API
Edit the `configs/config.yaml` file and set your OpenAI API key:
```yaml
openai:
  api_key: "your-openai-api-key"
  model: "gpt-4o-mini"
  base_url: "https://api.openai.com/v1"
  use_llm: true
```

4. Start the backend service
```bash
go run main.go
```

The service will start at `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory
```bash
cd frontend
```

2. Install dependencies
```bash
npm install
```

3. Start the development server
```bash
npm run dev
```

The frontend will start at `http://localhost:5173`

## 📋 API Endpoints

### Task Management
- `POST /api/v1/tasks` - Create a new task
- `GET /api/v1/tasks` - Get all tasks
- `GET /api/v1/tasks/:id` - Get a specific task
- `PUT /api/v1/tasks/:id/status` - Update task status
- `DELETE /api/v1/tasks/:id` - Delete a task

### Configuration Management
- `GET /api/v1/config` - Get configuration information
- `PUT /api/v1/config` - Update configuration

### Health Check
- `GET /health` - Service health status

## ⚠️ Current Issues

### 1. LLM Integration Issues
- **API Key Configuration**: Default configuration uses example API key, users need to manually configure a real OpenAI API key
- **Error Handling**: Error handling for LLM API call failures can be further optimized
- **Rate Limiting**: No implementation of API call rate limiting and retry mechanisms

### 2. Database Issues
- **Production Environment**: Currently using in-memory database (`:memory:`), data will be lost after restart
- **Data Migration**: Missing database migration scripts and version management

### 3. Frontend Issues
- **Error Handling**: Frontend error handling and user feedback can be more comprehensive
- **Responsive Design**: Mobile adaptation needs optimization
- **Data Validation**: Form validation logic needs strengthening

### 4. Security Issues
- **CORS Configuration**: Current CORS settings are too permissive (`*`)
- **Input Validation**: Need to strengthen server-side input validation
- **Authentication & Authorization**: Missing user authentication and permission management

### 5. Test Coverage
- **Unit Tests**: Missing backend unit tests
- **Integration Tests**: Missing API integration tests
- **Frontend Tests**: Missing frontend component tests

## 🔄 Development Roadmap

### Short-term Goals (1-2 weeks)
1. **Fix Configuration Issues**
   - Create configuration file templates
   - Add environment variable support
   - Replace in-memory database with file-based database

2. **Improve Error Handling**
   - Standardize error response format
   - Add detailed error logging
   - Improve frontend error notifications

3. **Optimize LLM Integration**
   - Add API call retry mechanisms
   - Implement request caching
   - Support more LLM providers

### Medium-term Goals (1-2 months)
1. **Add User System**
   - User registration and login
   - JWT authentication
   - Permission management

2. **Enhanced Features**
   - Task search and filtering
   - Task export functionality
   - Email notifications

3. **Testing and Documentation**
   - Write unit tests
   - Generate API documentation
   - Deployment documentation

### Long-term Goals (3-6 months)
1. **Performance Optimization**
   - Database query optimization
   - Frontend performance optimization
   - Caching strategies

2. **Extended Features**
   - Team collaboration features
   - Project management
   - Data analytics and reporting

3. **Deployment and Operations**
   - Docker containerization
   - CI/CD pipelines
   - Monitoring and logging systems

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Contact

For questions or suggestions, please create an Issue or contact the project maintainers.
