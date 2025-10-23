# Student Dashboard - Architecture Design

## Project Overview

**Project Name**: Trunchbull (Student Dashboard)
**Purpose**: Self-hosted dashboard for parents to monitor their children's academic progress across multiple learning platforms
**Target Platforms**: Schoology (LMS) and PowerSchool/The Source (SIS)
**Deployment**: Home lab, containerized

---

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Browser                               │
│  ┌────────────────────────────────────────────────────┐     │
│  │              React Frontend (SPA)                   │     │
│  │  ┌──────────┐ ┌──────────┐ ┌───────────────────┐  │     │
│  │  │Assignment│ │  Grades  │ │  Calendar/Events  │  │     │
│  │  │  Widget  │ │  Widget  │ │      Widget       │  │     │
│  │  └──────────┘ └──────────┘ └───────────────────┘  │     │
│  │  ┌──────────┐ ┌──────────┐ ┌───────────────────┐  │     │
│  │  │   GPA    │ │ Messages │ │  Student Selector │  │     │
│  │  │  Widget  │ │  Widget  │ │      Widget       │  │     │
│  │  └──────────┘ └──────────┘ └───────────────────┘  │     │
│  └────────────────────────────────────────────────────┘     │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTPS (REST API)
                       ↓
┌─────────────────────────────────────────────────────────────┐
│                   Docker Container                           │
│  ┌────────────────────────────────────────────────────┐     │
│  │               Go Backend Service                    │     │
│  │                                                      │     │
│  │  ┌─────────────────────────────────────────────┐   │     │
│  │  │           HTTP Server (Gin/Echo)            │   │     │
│  │  │  ┌────────────┐ ┌──────────────────────┐   │   │     │
│  │  │  │  REST API  │ │  WebSocket (future)  │   │   │     │
│  │  │  └────────────┘ └──────────────────────┘   │   │     │
│  │  └─────────────────────────────────────────────┘   │     │
│  │                                                      │     │
│  │  ┌─────────────────────────────────────────────┐   │     │
│  │  │         Authentication & OAuth Manager      │   │     │
│  │  │  ┌───────────────┐  ┌──────────────────┐   │   │     │
│  │  │  │  Schoology    │  │   PowerSchool    │   │   │     │
│  │  │  │ OAuth Client  │  │   OAuth Client   │   │   │     │
│  │  │  └───────────────┘  └──────────────────┘   │   │     │
│  │  └─────────────────────────────────────────────┘   │     │
│  │                                                      │     │
│  │  ┌─────────────────────────────────────────────┐   │     │
│  │  │            API Client Layer                 │   │     │
│  │  │  ┌───────────────┐  ┌──────────────────┐   │   │     │
│  │  │  │  Schoology    │  │   PowerSchool    │   │   │     │
│  │  │  │  API Client   │  │   API Client     │   │   │     │
│  │  │  └───────────────┘  └──────────────────┘   │   │     │
│  │  │  │                   │                      │   │     │
│  │  │  │ Rate Limiter      │ Rate Limiter         │   │     │
│  │  │  │ Circuit Breaker   │ Circuit Breaker      │   │     │
│  │  │  │ Retry Logic       │ Retry Logic          │   │     │
│  │  └─────────────────────────────────────────────┘   │     │
│  │                                                      │     │
│  │  ┌─────────────────────────────────────────────┐   │     │
│  │  │          Business Logic Layer               │   │     │
│  │  │  ┌──────────────────────────────────────┐   │   │     │
│  │  │  │  Assignment Aggregator               │   │   │     │
│  │  │  │  Grade Calculator (GPA)               │   │   │     │
│  │  │  │  Event/Calendar Aggregator            │   │   │     │
│  │  │  │  Message Aggregator                   │   │   │     │
│  │  │  └──────────────────────────────────────┘   │   │     │
│  │  └─────────────────────────────────────────────┘   │     │
│  │                                                      │     │
│  │  ┌─────────────────────────────────────────────┐   │     │
│  │  │            Data Sync Service                │   │     │
│  │  │  (Background jobs with cron/ticker)         │   │     │
│  │  └─────────────────────────────────────────────┘   │     │
│  │                                                      │     │
│  │  ┌─────────────────────────────────────────────┐   │     │
│  │  │             Cache Layer                     │   │     │
│  │  │  (In-memory map with TTL or Redis)          │   │     │
│  │  └─────────────────────────────────────────────┘   │     │
│  └────────────────────────────────────────────────────┘     │
│                                                               │
│  ┌────────────────────────────────────────────────────┐     │
│  │              SQLite Database                        │     │
│  │  ┌────────┐ ┌────────┐ ┌────────┐ ┌──────────┐    │     │
│  │  │Students│ │Courses │ │ Grades │ │Assignment│    │     │
│  │  └────────┘ └────────┘ └────────┘ └──────────┘    │     │
│  │  ┌────────┐ ┌────────┐ ┌────────┐                  │     │
│  │  │ Events │ │Messages│ │ Config │                  │     │
│  │  └────────┘ └────────┘ └────────┘                  │     │
│  └────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
                       │
                       │ HTTPS
                       ↓
        ┌──────────────────────────────┐
        │   External APIs               │
        │  ┌────────────┐ ┌───────────┐│
        │  │ Schoology  │ │PowerSchool││
        │  │    API     │ │    API    ││
        │  └────────────┘ └───────────┘│
        └──────────────────────────────┘
```

---

## Technology Stack

### Backend: Go

**Why Go?**
- Native concurrency (goroutines) perfect for concurrent API calls
- Strong standard library (`net/http`, `encoding/json`)
- Excellent performance and low memory footprint
- Simple deployment (single binary)
- Great OAuth libraries available
- Built-in HTTP client with good control over timeouts and retries

**Key Libraries:**
- **HTTP Framework**: `gin-gonic/gin` or `labstack/echo` (lightweight, fast)
- **OAuth 1.0a**: `dghubble/oauth1` (for Schoology)
- **OAuth 2.0**: `golang.org/x/oauth2` (for PowerSchool)
- **Database**: `mattn/go-sqlite3` with `database/sql`
- **Configuration**: `spf13/viper` (environment + config files)
- **Logging**: `sirupsen/logrus` or `uber-go/zap`
- **Caching**: `patrickmn/go-cache` (in-memory) or `go-redis/redis`
- **HTTP Client**: Standard `net/http` with custom transport for rate limiting
- **Job Scheduling**: `robfig/cron` or custom ticker-based scheduler

### Frontend: React

**Why React?**
- Component-based architecture ideal for dashboard widgets
- Large ecosystem for charts and data visualization
- Excellent developer experience
- Strong mobile-responsive capabilities

**Key Libraries:**
- **Framework**: `create-react-app` or `vite` (modern, fast builds)
- **UI Components**: `Material-UI` or `shadcn/ui` (accessible, polished)
- **Data Fetching**: `TanStack Query` (react-query) - caching, refetching
- **Routing**: `react-router` (if multi-page needed)
- **State Management**: React Context + hooks (simple) or `zustand` (if needed)
- **Charts**: `recharts` or `chart.js` (for GPA trends, grade visualization)
- **Date/Time**: `date-fns` (lightweight)
- **HTTP Client**: `axios` or native `fetch`

### Database: SQLite

**Why SQLite?**
- Serverless, zero-configuration
- Perfect for single-family use case
- Low overhead, fast for small datasets
- File-based (easy backups)
- No separate database server needed in Docker

**Schema** (see Database Schema section below)

### Containerization: Docker

**Components:**
- **Backend Container**: Go application
- **Frontend**: Static build served by backend or separate nginx container
- **Persistent Volume**: SQLite database file
- **Docker Compose**: Orchestrate services

---

## API Design

### Backend REST API Endpoints

#### Authentication
```
POST   /api/auth/schoology/init       - Initiate Schoology OAuth flow
GET    /api/auth/schoology/callback   - OAuth callback handler
POST   /api/auth/powerschool/init     - Initiate PowerSchool OAuth flow
GET    /api/auth/powerschool/callback - OAuth callback handler
GET    /api/auth/status                - Check authentication status
DELETE /api/auth/logout                - Clear stored credentials
```

#### Students
```
GET    /api/students                   - List all configured students
POST   /api/students                   - Add a new student
GET    /api/students/:id               - Get student details
DELETE /api/students/:id               - Remove student
```

#### Dashboard Data
```
GET    /api/dashboard/:studentId       - Get complete dashboard data
GET    /api/assignments/:studentId     - Get outstanding assignments
GET    /api/grades/:studentId          - Get current grades
GET    /api/gpa/:studentId             - Get calculated GPA
GET    /api/events                     - Get school calendar events
GET    /api/messages/:studentId        - Get unread messages
```

#### Sync
```
POST   /api/sync/:studentId            - Trigger manual sync
GET    /api/sync/status                - Get last sync time and status
```

#### Health
```
GET    /health                         - Health check endpoint
GET    /api/status                     - API status and version
```

### Request/Response Examples

#### Dashboard Data Response
```json
{
  "student": {
    "id": "student-1",
    "name": "Jane Doe",
    "grade": 7
  },
  "summary": {
    "outstandingAssignments": 5,
    "currentGPA": 3.67,
    "unreadMessages": 2,
    "upcomingEvents": 3
  },
  "assignments": [
    {
      "id": "assign-123",
      "title": "Math Homework Ch. 5",
      "course": "Algebra I",
      "dueDate": "2025-10-25T23:59:00Z",
      "status": "pending",
      "source": "schoology"
    }
  ],
  "grades": [
    {
      "course": "Algebra I",
      "currentGrade": "A-",
      "percentage": 91.5,
      "source": "schoology"
    }
  ],
  "events": [
    {
      "title": "No School - Teacher Planning",
      "date": "2025-10-27",
      "source": "powerschool"
    }
  ],
  "messages": [
    {
      "id": "msg-456",
      "from": "Ms. Smith",
      "subject": "Parent-Teacher Conferences",
      "date": "2025-10-22T10:30:00Z",
      "unread": true,
      "source": "schoology"
    }
  ],
  "lastSync": "2025-10-23T14:30:00Z"
}
```

---

## Database Schema

### SQLite Schema

```sql
-- Students table
CREATE TABLE students (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    grade_level INTEGER,
    schoology_user_id TEXT,
    powerschool_student_id TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Courses table
CREATE TABLE courses (
    id TEXT PRIMARY KEY,
    student_id TEXT NOT NULL,
    name TEXT NOT NULL,
    teacher TEXT,
    period TEXT,
    source TEXT NOT NULL, -- 'schoology' or 'powerschool'
    external_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
);

-- Assignments table
CREATE TABLE assignments (
    id TEXT PRIMARY KEY,
    student_id TEXT NOT NULL,
    course_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    due_date TIMESTAMP,
    status TEXT, -- 'pending', 'submitted', 'graded', 'late'
    source TEXT NOT NULL,
    external_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- Grades table
CREATE TABLE grades (
    id TEXT PRIMARY KEY,
    student_id TEXT NOT NULL,
    course_id TEXT NOT NULL,
    assignment_id TEXT, -- NULL for overall course grade
    score REAL,
    max_score REAL,
    percentage REAL,
    letter_grade TEXT,
    grading_period TEXT,
    source TEXT NOT NULL,
    external_id TEXT,
    recorded_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE CASCADE
);

-- Events table
CREATE TABLE events (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    event_date DATE NOT NULL,
    event_type TEXT, -- 'holiday', 'no-school', 'early-release', 'event'
    source TEXT NOT NULL,
    external_id TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Messages table
CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    student_id TEXT NOT NULL,
    from_name TEXT NOT NULL,
    from_email TEXT,
    subject TEXT,
    body TEXT,
    received_at TIMESTAMP,
    read BOOLEAN DEFAULT FALSE,
    source TEXT NOT NULL,
    external_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
);

-- Config table (for OAuth tokens and app settings)
CREATE TABLE config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    encrypted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sync status table
CREATE TABLE sync_status (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    student_id TEXT NOT NULL,
    source TEXT NOT NULL,
    data_type TEXT NOT NULL, -- 'assignments', 'grades', 'events', 'messages'
    last_sync TIMESTAMP,
    status TEXT, -- 'success', 'failed', 'in_progress'
    error_message TEXT,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    UNIQUE(student_id, source, data_type)
);

-- Indexes for performance
CREATE INDEX idx_assignments_student_due ON assignments(student_id, due_date);
CREATE INDEX idx_assignments_status ON assignments(status);
CREATE INDEX idx_grades_student_course ON grades(student_id, course_id);
CREATE INDEX idx_events_date ON events(event_date);
CREATE INDEX idx_messages_student_unread ON messages(student_id, read);
CREATE INDEX idx_sync_status_lookup ON sync_status(student_id, source, data_type);
```

---

## Component Design

### Backend Components

#### 1. HTTP Server
- **Responsibility**: Handle HTTP requests, routing, middleware
- **Framework**: Gin or Echo
- **Middleware**: CORS, logging, error handling, rate limiting

#### 2. OAuth Manager
- **Responsibility**: Manage OAuth flows for both platforms
- **Components**:
  - Schoology OAuth 1.0a client
  - PowerSchool OAuth 2.0 client
  - Token storage and refresh
  - Session management

#### 3. API Clients
- **Schoology Client**:
  - Wrapper around Schoology REST API
  - Rate limiting (client-side throttling)
  - Retry logic with exponential backoff
  - Circuit breaker pattern
  - Request/response logging

- **PowerSchool Client**:
  - Wrapper around PowerSchool REST API
  - Rate limiting with `Retry-After` header support
  - Retry logic with exponential backoff
  - Circuit breaker pattern
  - Request/response logging

#### 4. Business Logic Services
- **Assignment Service**: Aggregate assignments from both platforms
- **Grade Service**: Fetch and calculate grades, GPA
- **Event Service**: Aggregate calendar events
- **Message Service**: Fetch and aggregate messages
- **Sync Service**: Background data synchronization

#### 5. Data Layer
- **Repository Pattern**: Abstract database operations
- **Models**: Go structs matching database schema
- **Migrations**: Database schema versioning

#### 6. Cache Layer
- **In-memory cache** with TTL for API responses
- **Cache keys** by student ID + data type
- **TTL strategy**:
  - Assignments: 15 minutes
  - Grades: 30 minutes
  - Events: 24 hours
  - Messages: 10 minutes

### Frontend Components

#### Page Layout
```
App
├── Header (AppBar with logo, student selector)
├── Sidebar (navigation - future)
└── Dashboard
    ├── StudentSelector
    ├── SummaryCards (outstanding assignments, GPA, messages, events)
    ├── AssignmentsWidget
    ├── GradesWidget
    ├── CalendarWidget
    ├── MessagesWidget
    └── LastSyncStatus
```

#### Component Hierarchy
- **App**: Root component, routing setup
- **Dashboard**: Main dashboard page
- **StudentSelector**: Dropdown to switch between students
- **SummaryCards**: 4 cards with key metrics
- **AssignmentsWidget**: List of outstanding assignments with due dates
- **GradesWidget**: Current grades by course, GPA
- **CalendarWidget**: Upcoming events and days off
- **MessagesWidget**: Unread messages from teachers
- **LastSyncStatus**: Shows last sync time, manual refresh button

---

## Data Flow

### Initial Setup Flow
```
1. User starts Docker container
2. User visits web UI
3. User initiates Schoology OAuth flow
   ├─> Backend redirects to Schoology
   ├─> User authenticates
   ├─> Schoology redirects back with token
   └─> Backend stores OAuth tokens
4. User initiates PowerSchool OAuth flow (same as above)
5. User adds student(s) with their platform IDs
6. Backend triggers initial data sync
7. Dashboard displays data
```

### Data Sync Flow
```
Background Sync (every 15-30 minutes):
1. Sync Service wakes up (cron job)
2. For each student:
   ├─> Fetch assignments from Schoology
   ├─> Fetch grades from Schoology
   ├─> Fetch assignments from PowerSchool
   ├─> Fetch grades from PowerSchool
   ├─> Fetch events from both platforms
   ├─> Fetch messages from Schoology
   ├─> Update database
   ├─> Invalidate cache
   └─> Record sync status
3. Log sync completion
```

### Dashboard Load Flow
```
1. Frontend requests dashboard data for student
2. Backend checks cache
   ├─> If cache hit: return cached data
   └─> If cache miss:
       ├─> Fetch from database
       ├─> If data stale (> 1 hour): trigger background sync
       ├─> Cache result
       └─> Return data
3. Frontend renders dashboard
```

---

## Security Considerations

### Credential Storage
- OAuth tokens stored in database (config table)
- **Encryption at rest**: Use SQLite encryption extension or application-level encryption
- **Environment variables**: API client IDs and secrets (never commit to git)
- **.env file**: For local development (add to .gitignore)
- **Docker secrets**: For production deployment

### Network Security
- **HTTPS only**: Enforce HTTPS for production deployments
- **No CORS in production**: Frontend served from same origin
- **API rate limiting**: Prevent abuse of backend API
- **Authentication**: Future: Add basic auth or OAuth for dashboard access

### Data Privacy
- **Local storage only**: All data stays on user's infrastructure
- **No telemetry**: No data sent to external services
- **Audit logs**: Log all API calls for transparency
- **Data retention**: Configurable retention policy (default: 90 days)

### API Security
- **Token refresh**: Implement token refresh before expiry
- **Token rotation**: Support credential rotation
- **Error handling**: Never expose tokens/credentials in error messages
- **Secure logging**: Redact sensitive data in logs

---

## Configuration

### Environment Variables
```bash
# Server
PORT=8080
ENV=production  # development, staging, production

# Database
DATABASE_PATH=/data/trunchbull.db
DATABASE_ENCRYPTION_KEY=<secret>

# Schoology OAuth
SCHOOLOGY_CONSUMER_KEY=<key>
SCHOOLOGY_CONSUMER_SECRET=<secret>
SCHOOLOGY_BASE_URL=https://api.schoology.com/v1

# PowerSchool OAuth
POWERSCHOOL_CLIENT_ID=<id>
POWERSCHOOL_CLIENT_SECRET=<secret>
POWERSCHOOL_BASE_URL=https://meany.seattleschools.org  # Or district URL

# Sync Configuration
SYNC_INTERVAL=30m  # How often to sync data
CACHE_TTL_ASSIGNMENTS=15m
CACHE_TTL_GRADES=30m
CACHE_TTL_EVENTS=24h
CACHE_TTL_MESSAGES=10m

# Rate Limiting
RATE_LIMIT_SCHOOLOGY=60  # requests per minute
RATE_LIMIT_POWERSCHOOL=30  # requests per minute

# Logging
LOG_LEVEL=info  # debug, info, warn, error
LOG_FORMAT=json  # json, text
```

### Config File (config.yaml)
```yaml
server:
  port: 8080
  read_timeout: 10s
  write_timeout: 10s

database:
  path: /data/trunchbull.db
  max_open_conns: 10
  max_idle_conns: 5

sync:
  interval: 30m
  retry_attempts: 3
  retry_backoff: exponential

cache:
  enabled: true
  ttl:
    assignments: 15m
    grades: 30m
    events: 24h
    messages: 10m

rate_limits:
  schoology: 60  # per minute
  powerschool: 30  # per minute
```

---

## Deployment

### Docker Compose

```yaml
version: '3.8'

services:
  trunchbull:
    build: .
    container_name: trunchbull
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - ./data:/data
      - ./config:/config
    environment:
      - ENV=production
      - DATABASE_PATH=/data/trunchbull.db
      - PORT=8080
    env_file:
      - .env
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

### Dockerfile

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o trunchbull ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/trunchbull .
COPY --from=builder /app/frontend/build ./frontend/build

# Create data directory
RUN mkdir -p /data

EXPOSE 8080

CMD ["./trunchbull"]
```

---

## Development Phases

### Phase 1: MVP (Week 1-2)
**Goal**: Basic dashboard with Schoology integration

- [x] Research complete
- [ ] Project structure setup
- [ ] Go backend skeleton with Gin
- [ ] SQLite database setup with migrations
- [ ] Schoology OAuth implementation
- [ ] Schoology API client with rate limiting
- [ ] Assignment and grade endpoints
- [ ] React frontend skeleton
- [ ] Dashboard with assignments and grades widgets
- [ ] Docker containerization
- [ ] Basic documentation

**Deliverable**: Working dashboard showing Schoology assignments and grades for one student

### Phase 2: Full Integration (Week 3-4)
**Goal**: PowerSchool integration and multiple students

- [ ] PowerSchool OAuth implementation
- [ ] PowerSchool API client
- [ ] Data aggregation logic
- [ ] Multiple student support
- [ ] Calendar/events widget
- [ ] Background sync service
- [ ] Cache layer
- [ ] Error handling and retries
- [ ] Enhanced UI with Material-UI

**Deliverable**: Full-featured dashboard with both platforms and multiple students

### Phase 3: Polish & Features (Week 5-6)
**Goal**: Production-ready with additional features

- [ ] GPA calculation engine
- [ ] Message inbox widget
- [ ] Real-time updates (WebSocket or Server-Sent Events)
- [ ] Mobile-responsive design
- [ ] Dashboard customization
- [ ] Notification system (email/push)
- [ ] Data export (CSV/PDF)
- [ ] Comprehensive error logging
- [ ] Performance optimization
- [ ] Security hardening
- [ ] User documentation

**Deliverable**: Production-ready application with documentation

---

## Testing Strategy

### Backend Testing
- **Unit tests**: Business logic, data transformations
- **Integration tests**: API clients, database operations
- **End-to-end tests**: Full API flows
- **Mock external APIs**: Use mock servers for Schoology/PowerSchool

### Frontend Testing
- **Component tests**: React Testing Library
- **Integration tests**: Dashboard data flow
- **E2E tests**: Cypress or Playwright

### Load Testing
- Simulate multiple students and concurrent requests
- Test rate limiting and circuit breakers
- Verify cache performance

---

## Monitoring & Observability

### Logging
- Structured JSON logs
- Log levels: DEBUG, INFO, WARN, ERROR
- Log API requests/responses (redact sensitive data)
- Log sync job outcomes

### Metrics (Future)
- API response times
- Cache hit rates
- Sync job success/failure rates
- External API error rates

### Alerts (Future)
- Sync failures
- API authentication failures
- Rate limit violations
- High error rates

---

## Future Enhancements

### Phase 4+
- [ ] Multi-family support (separate databases or schema)
- [ ] Dashboard access control (authentication)
- [ ] Notification system (email, SMS, push)
- [ ] Historical data and trends (grade trends over time)
- [ ] Predictive analytics (at-risk assignments)
- [ ] Integration with additional platforms (Canvas, Google Classroom)
- [ ] Mobile app (React Native)
- [ ] Parent-teacher messaging
- [ ] Assignment submission tracking
- [ ] Calendar sync (iCal export)

---

## Risk Mitigation

### Technical Risks
| Risk | Mitigation |
|------|------------|
| API access denied by district | Document self-hosted nature, contact district early |
| Rate limiting too aggressive | Implement caching, reduce sync frequency |
| OAuth token expiry | Implement token refresh, notify user |
| API changes breaking integration | Version API clients, monitor for changes |
| Data breach | Encryption at rest, secure credential storage, audit logs |

### Operational Risks
| Risk | Mitigation |
|------|------------|
| User doesn't have API access | Provide clear instructions for requesting access |
| Home lab downtime | Document backup/restore procedures |
| Data loss | Automated backups, SQLite database file |

---

## Success Metrics

### MVP Success Criteria
- Successfully authenticate with Schoology
- Display current assignments and grades
- Run in Docker container
- Update data every 30 minutes
- Basic documentation complete

### Production Success Criteria
- Both platforms integrated
- Multiple students supported
- All dashboard features working
- Zero data breaches
- Comprehensive documentation
- Positive user feedback

---

## Open Questions

1. **Should we implement our own OAuth consent screens or use platform defaults?**
   - Recommendation: Use platform defaults for MVP, consider custom for better UX later

2. **How should we handle multiple grading periods?**
   - Recommendation: Default to current grading period, allow user to select

3. **Should we calculate GPA ourselves or rely on platform data?**
   - Recommendation: Try platform first, fall back to calculation if unavailable

4. **Do we need a separate database for each family?**
   - Recommendation: Single database for MVP, consider multi-tenancy for Phase 4

5. **Should the frontend be built into the Go binary or separate container?**
   - Recommendation: Embed in Go binary for simplicity (using embed.FS)

---

**Document Version**: 1.0
**Last Updated**: 2025-10-23
**Status**: Draft
