# Existing Libraries Analysis & Recommendation

## Executive Summary

**Recommendation: Build new Go client libraries as standalone projects**

After evaluating existing libraries, none are suitable for reuse:
- All existing credential-based libraries are in Python
- None are actively maintained (last updates 2021 or earlier)
- No Go libraries exist for credential-based access
- Building Go libraries fills a gap in the ecosystem

**Proposed Solution:**
1. Create `powerschool-go` - standalone Go client library
2. Create `schoology-go` - standalone Go client library
3. Implement Approach 2 (credential automation) using chromedp
4. Use these libraries in Trunchbull

---

## Existing Library Evaluation

### PowerSchool Libraries

#### 1. psscraper (Python)
- **Repository**: https://github.com/Desperationis/psscraper
- **Last Update**: June 2021 (4+ years ago)
- **License**: GPL-3.0
- **Language**: Python
- **Approach**: Selenium + BeautifulSoup (browser automation)
- **Status**: ❌ **NOT MAINTAINED**

**Pros:**
- Works (as of 2021)
- Browser-based (handles JavaScript)
- Available on PyPI

**Cons:**
- Not maintained for 4 years
- GPL license (copyleft)
- Python, not Go
- Heavy (requires Firefox + geckodriver)
- Slow (browser automation)

#### 2. ps.py (Python)
- **Repository**: https://github.com/ouiliame/ps.py
- **Language**: Python
- **Approach**: HTTP requests + authentication
- **Status**: ❌ **NOT MAINTAINED**

**Author's Note:**
> "This was written more than a year ago, and last I checked it was working fine. However, PowerSchool may have undergone upgrades which could disrupt the functionality"

**Cons:**
- Explicitly noted as old by author
- May not work with current PowerSchool
- Only tested with one district (AUSD)
- Python, not Go

#### 3. ps_scraper (Python)
- **Repository**: https://github.com/jarulsamy/ps_scraper
- **Language**: Python
- **Approach**: CLI tool, web scraping
- **Status**: Unknown, likely not maintained

**Cons:**
- CLI tool, not a library
- Python, not Go
- Unknown maintenance status

### Schoology Libraries

#### 1. SchoologyMessageWebScraper (Python)
- **Repository**: https://github.com/Saptak625/SchoologyMessageWebScraper
- **Language**: Python
- **Approach**: requests.Session() + BeautifulSoup
- **Scope**: Messages only (not grades/assignments)
- **Status**: Unknown

**Cons:**
- Limited scope (messages only)
- Python, not Go
- Unclear maintenance

#### 2. sgy-sgy (JavaScript)
- **Repository**: https://github.com/SheepTester/sgy-sgy
- **Language**: JavaScript/Node.js
- **Approach**: Cookie-based access
- **Status**: Unknown

**Cons:**
- JavaScript, not Go
- Requires manual cookie extraction
- Unclear maintenance

#### 3. github.com/jbeuckm/schoology (Go)
- **Repository**: https://github.com/jbeuckm/schoology
- **Language**: Go ✅
- **Approach**: OAuth 1.0a (official API)
- **Status**: Unknown

**Cons:**
- Uses official OAuth (requires district approval)
- Not credential-based
- Unclear if maintained

### Summary Table

| Library | Language | Maintained | Approach | Suitable? |
|---------|----------|-----------|----------|-----------|
| psscraper | Python | ❌ No (2021) | Browser | ❌ |
| ps.py | Python | ❌ No | HTTP | ❌ |
| ps_scraper | Python | ❌ Unknown | CLI | ❌ |
| SchoologyMessageWebScraper | Python | ❌ Unknown | HTTP | ❌ |
| sgy-sgy | JavaScript | ❌ Unknown | Cookies | ❌ |
| jbeuckm/schoology | Go | ❌ Unknown | OAuth | ❌ |

**Conclusion: No suitable existing libraries for reuse**

---

## Proposed Architecture: Standalone Go Client Libraries

### Design Philosophy

1. **Separate, reusable libraries** - Not monolithic
2. **Go-native** - Leverage Go's strengths
3. **Well-tested** - Comprehensive test coverage
4. **Well-documented** - Easy for others to use
5. **MIT licensed** - Maximum reusability
6. **Approach 2 from the start** - Credential-based automation

### Library 1: `powerschool-go`

#### Repository Structure
```
powerschool-go/
├── README.md
├── LICENSE (MIT)
├── go.mod
├── client.go           # Main client
├── auth.go             # Authentication
├── student.go          # Student endpoints
├── grades.go           # Grades endpoints
├── assignments.go      # Assignment endpoints
├── attendance.go       # Attendance endpoints
├── calendar.go         # Calendar/events
├── types.go            # Data models
├── errors.go           # Error types
├── examples/
│   ├── basic/
│   ├── credentials/
│   └── session/
└── internal/
    ├── browser/        # Chromedp automation
    └── scraper/        # HTML parsing
```

#### API Design

```go
package powerschool

import "context"

// Client is the main PowerSchool client
type Client struct {
    baseURL     string
    credentials *Credentials
    session     *Session
    httpClient  *http.Client
}

// NewClient creates a new PowerSchool client
func NewClient(baseURL string, opts ...Option) (*Client, error)

// Option configures the client
type Option func(*Client) error

// WithCredentials sets username/password authentication
func WithCredentials(username, password string) Option

// WithSession sets session token authentication
func WithSession(token string, cookies []*http.Cookie) Option

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) Option

// Authenticate logs in and obtains session
func (c *Client) Authenticate(ctx context.Context) error

// GetStudents returns all students accessible to this account
func (c *Client) GetStudents(ctx context.Context) ([]*Student, error)

// GetGrades returns grades for a student
func (c *Client) GetGrades(ctx context.Context, studentID string) ([]*Grade, error)

// GetAssignments returns assignments for a student
func (c *Client) GetAssignments(ctx context.Context, studentID string) ([]*Assignment, error)

// GetGPA returns GPA for a student
func (c *Client) GetGPA(ctx context.Context, studentID string) (*GPA, error)

// GetAttendance returns attendance records
func (c *Client) GetAttendance(ctx context.Context, studentID string) ([]*Attendance, error)

// GetCalendar returns calendar events
func (c *Client) GetCalendar(ctx context.Context) ([]*Event, error)

// Types

type Student struct {
    ID          string
    Name        string
    GradeLevel  int
    SchoolID    string
    SchoolName  string
}

type Grade struct {
    CourseID       string
    CourseName     string
    Teacher        string
    Period         string
    CurrentGrade   string
    Percentage     float64
    LetterGrade    string
    GradingPeriod  string
    LastUpdated    time.Time
}

type Assignment struct {
    ID             string
    CourseID       string
    Title          string
    Description    string
    DueDate        time.Time
    Category       string
    Score          *float64
    MaxScore       float64
    Percentage     *float64
    Status         AssignmentStatus
}

type AssignmentStatus string

const (
    StatusPending   AssignmentStatus = "pending"
    StatusSubmitted AssignmentStatus = "submitted"
    StatusGraded    AssignmentStatus = "graded"
    StatusLate      AssignmentStatus = "late"
    StatusMissing   AssignmentStatus = "missing"
)

type GPA struct {
    Current    float64
    Cumulative float64
    Weighted   bool
}

type Attendance struct {
    Date        time.Time
    Status      AttendanceStatus
    Period      string
    CourseName  string
}

type AttendanceStatus string

const (
    AttendancePresent AttendanceStatus = "present"
    AttendanceAbsent  AttendanceStatus = "absent"
    AttendanceTardy   AttendanceStatus = "tardy"
    AttendanceExcused AttendanceStatus = "excused"
)

type Event struct {
    ID          string
    Title       string
    Description string
    Date        time.Time
    EventType   EventType
}

type EventType string

const (
    EventHoliday      EventType = "holiday"
    EventNoSchool     EventType = "no_school"
    EventEarlyRelease EventType = "early_release"
    EventOther        EventType = "other"
)
```

#### Example Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/leftathome/powerschool-go"
)

func main() {
    ctx := context.Background()

    // Option 1: Credential-based authentication
    client, err := powerschool.NewClient(
        "https://powerschool.seattleschools.org",
        powerschool.WithCredentials("parent_username", "password"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Authenticate
    if err := client.Authenticate(ctx); err != nil {
        log.Fatal(err)
    }

    // Get students
    students, err := client.GetStudents(ctx)
    if err != nil {
        log.Fatal(err)
    }

    for _, student := range students {
        fmt.Printf("Student: %s (Grade %d)\n", student.Name, student.GradeLevel)

        // Get grades
        grades, err := client.GetGrades(ctx, student.ID)
        if err != nil {
            log.Fatal(err)
        }

        for _, grade := range grades {
            fmt.Printf("  %s: %s (%.1f%%)\n",
                grade.CourseName, grade.LetterGrade, grade.Percentage)
        }

        // Get assignments
        assignments, err := client.GetAssignments(ctx, student.ID)
        if err != nil {
            log.Fatal(err)
        }

        pending := 0
        for _, a := range assignments {
            if a.Status == powerschool.StatusPending {
                pending++
            }
        }
        fmt.Printf("  Outstanding assignments: %d\n", pending)
    }
}
```

### Library 2: `schoology-go`

#### Repository Structure
```
schoology-go/
├── README.md
├── LICENSE (MIT)
├── go.mod
├── client.go           # Main client
├── auth.go             # Authentication
├── courses.go          # Course endpoints
├── assignments.go      # Assignment endpoints
├── grades.go           # Grades endpoints
├── messages.go         # Messages endpoints
├── calendar.go         # Calendar/events
├── types.go            # Data models
├── errors.go           # Error types
├── examples/
└── internal/
    └── browser/        # Chromedp automation
```

#### API Design

```go
package schoology

// Client is the main Schoology client
type Client struct {
    host        string
    credentials *Credentials
    session     *Session
    httpClient  *http.Client
}

// NewClient creates a new Schoology client
func NewClient(host string, opts ...Option) (*Client, error)

// WithCredentials sets username/password authentication
func WithCredentials(username, password string) Option

// WithSession sets session cookie authentication
func WithSession(sessID, csrfToken, csrfKey, uid string) Option

// Authenticate logs in and obtains session
func (c *Client) Authenticate(ctx context.Context) error

// GetCourses returns all courses
func (c *Client) GetCourses(ctx context.Context) ([]*Course, error)

// GetAssignments returns assignments for a course
func (c *Client) GetAssignments(ctx context.Context, courseID string) ([]*Assignment, error)

// GetGrades returns grades
func (c *Client) GetGrades(ctx context.Context) ([]*Grade, error)

// GetMessages returns messages
func (c *Client) GetMessages(ctx context.Context, unreadOnly bool) ([]*Message, error)

// GetCalendar returns calendar events
func (c *Client) GetCalendar(ctx context.Context) ([]*Event, error)

// Types similar to PowerSchool
```

---

## Implementation Plan

### Phase 1: `powerschool-go` (Week 1)

**Goals:**
- ✅ Repository setup
- ✅ Authentication with credentials (chromedp)
- ✅ Student listing
- ✅ Grade retrieval
- ✅ Basic assignment retrieval

**Tasks:**
1. Create repo with MIT license
2. Implement chromedp-based authentication
3. Parse HTML for student data
4. Parse gradebook HTML
5. Write tests
6. Document API
7. Publish v0.1.0

### Phase 2: `schoology-go` (Week 2)

**Goals:**
- ✅ Repository setup
- ✅ Authentication with credentials
- ✅ Course listing
- ✅ Assignment retrieval
- ✅ Grade retrieval

**Tasks:**
1. Create repo with MIT license
2. Implement chromedp-based authentication
3. Use internal API endpoints (iapi2)
4. Parse JSON responses
5. Write tests
6. Document API
7. Publish v0.1.0

### Phase 3: Enhance Libraries (Week 3)

**powerschool-go:**
- ✅ Attendance tracking
- ✅ Calendar/events
- ✅ Session persistence
- ✅ Error handling improvements

**schoology-go:**
- ✅ Messages
- ✅ Calendar integration
- ✅ Session persistence
- ✅ Error handling improvements

### Phase 4: Trunchbull Integration (Week 4)

**Goals:**
- ✅ Use both libraries in Trunchbull
- ✅ Build aggregation layer
- ✅ React dashboard
- ✅ Background sync

---

## Technical Decisions

### Authentication: chromedp vs. HTTP

**Recommendation: Start with chromedp, add HTTP later**

**chromedp Approach:**
```go
func (c *Client) authenticateWithBrowser(ctx context.Context, username, password string) error {
    allocCtx, cancel := chromedp.NewContext(ctx)
    defer cancel()

    var sessionCookies []*network.Cookie

    err := chromedp.Run(allocCtx,
        chromedp.Navigate(c.loginURL),
        chromedp.WaitVisible(`input[name="username"]`),
        chromedp.SendKeys(`input[name="username"]`, username),
        chromedp.SendKeys(`input[name="password"]`, password),
        chromedp.Click(`button[type="submit"]`),
        chromedp.WaitVisible(`#dashboard`),
        chromedp.ActionFunc(func(ctx context.Context) error {
            cookies, err := network.GetAllCookies().Do(ctx)
            sessionCookies = cookies
            return err
        }),
    )

    if err != nil {
        return fmt.Errorf("authentication failed: %w", err)
    }

    c.session = extractSession(sessionCookies)
    return nil
}
```

**Advantages:**
- Handles JavaScript
- Works with any login flow (SSO, 2FA, etc.)
- Easier to debug (can run non-headless)
- More reliable

**Disadvantages:**
- Heavier (requires Chrome)
- Slower initial login

**Future:** Add pure HTTP authentication as optimization for known login flows

### HTML Parsing: goquery

**Recommendation: goquery (jQuery-like selectors)**

```go
import "github.com/PuerkitoBio/goquery"

func parseGrades(html string) ([]*Grade, error) {
    doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
    if err != nil {
        return nil, err
    }

    var grades []*Grade
    doc.Find("table.grades tr.grade-row").Each(func(i int, s *goquery.Selection) {
        grade := &Grade{
            CourseName: s.Find(".course-name").Text(),
            // ...
        }
        grades = append(grades, grade)
    })

    return grades, nil
}
```

### Session Management

**Recommendation: Automatic refresh with fallback**

```go
type Session struct {
    Cookies   []*http.Cookie
    CSRFToken string
    ExpiresAt time.Time
    mu        sync.RWMutex
}

func (s *Session) IsValid() bool {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return time.Now().Before(s.ExpiresAt)
}

func (c *Client) ensureAuthenticated(ctx context.Context) error {
    if c.session != nil && c.session.IsValid() {
        return nil
    }

    if c.credentials == nil {
        return ErrSessionExpired
    }

    return c.Authenticate(ctx)
}
```

### Error Handling

```go
// Errors
var (
    ErrAuthFailed        = errors.New("authentication failed")
    ErrSessionExpired    = errors.New("session expired")
    ErrNotFound          = errors.New("not found")
    ErrRateLimited       = errors.New("rate limited")
    ErrInvalidCredentials = errors.New("invalid credentials")
)

// Error types
type Error struct {
    Code    ErrorCode
    Message string
    Err     error
}

func (e *Error) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

func (e *Error) Unwrap() error {
    return e.Err
}
```

---

## Benefits of This Approach

### For Trunchbull
1. ✅ Clean separation of concerns
2. ✅ Easy to test each library independently
3. ✅ Swap implementations easily
4. ✅ Simpler Trunchbull codebase

### For the Go Community
1. ✅ Fills gap - no Go libraries exist
2. ✅ Reusable by others
3. ✅ MIT license - maximum freedom
4. ✅ Well-documented, tested
5. ✅ Could become de facto standard

### For Maintenance
1. ✅ Separate versioning
2. ✅ Focused issues/PRs
3. ✅ Easier to onboard contributors
4. ✅ Better discoverability

---

## Required Dependencies

### powerschool-go
```go
require (
    github.com/chromedp/chromedp v0.9.5
    github.com/PuerkitoBio/goquery v1.9.0
)
```

### schoology-go
```go
require (
    github.com/chromedp/chromedp v0.9.5
)
```

### Trunchbull
```go
require (
    github.com/leftathome/powerschool-go v0.1.0
    github.com/leftathome/schoology-go v0.1.0
    github.com/gin-gonic/gin v1.10.0
    // ... rest of dependencies
)
```

---

## Timeline

| Week | Focus | Deliverables |
|------|-------|--------------|
| 1 | powerschool-go | v0.1.0 with core features |
| 2 | schoology-go | v0.1.0 with core features |
| 3 | Library enhancements | v0.2.0 for both |
| 4 | Trunchbull integration | Working dashboard |

Total: **4 weeks** to fully working Trunchbull with reusable libraries

---

## Next Steps

1. ✅ Get approval for this approach
2. ⬜ Create `powerschool-go` repository
3. ⬜ Implement authentication + basic features
4. ⬜ Write tests and documentation
5. ⬜ Publish v0.1.0
6. ⬜ Repeat for `schoology-go`
7. ⬜ Build Trunchbull using both libraries

---

## Questions for Discussion

1. **Repository Organization**
   - Separate repos under personal account?
   - Create `leftathome` GitHub org?
   - Keep everything in monorepo initially?

2. **Versioning**
   - Start at v0.1.0?
   - Semantic versioning?
   - When to hit v1.0.0?

3. **Testing**
   - Unit tests with mocked HTML?
   - Integration tests with real accounts?
   - CI/CD setup?

4. **Documentation**
   - GoDoc comments?
   - Separate documentation site?
   - Example repository?

---

**Document Version**: 1.0
**Last Updated**: 2025-10-23
**Status**: Proposal - Awaiting Approval
