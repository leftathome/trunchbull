# Alternative Authentication Approach - Credential-Based Access

## Executive Summary

Based on extensive research, we've discovered that **official API access requiring district approval is not necessary** for Trunchbull. Instead, we can use parent/student credentials directly, just like the official mobile apps and web portals do.

### Key Finding

Multiple successful projects have demonstrated that both Schoology and PowerSchool can be accessed programmatically using regular parent/student login credentials without requiring district-level API approval.

---

## Evidence of Viability

### Case Study: "The Source: SPS" App

A student-built app that:
- Accessed PowerSchool/The Source data for **17,000+ monthly active users**
- Used regular parent/student credentials (not official API)
- Operated successfully for **3 years** (2021-2024)
- Achieved 30,000+ total downloads
- Only became "unauthorized" when acquired and monetized

**Lesson**: Credential-based access is technically viable and was accepted by the district for years when used responsibly.

### Existing Open-Source Projects

#### PowerSchool
1. **psscraper** (https://github.com/Desperationis/psscraper)
   - Browser-based scraping with Selenium + BeautifulSoup
   - Uses parent/student username and password
   - Available on PyPI: `pip install psscraper`
   - Accesses: grades, courses, assignments, GPA

2. **ps.py** (https://github.com/ouiliame/ps.py)
   - Authentication & data fetching library
   - Provides methods for student info, GPA, courses, assignments
   - Pure authentication (not browser automation)

3. **ps_scraper** (https://github.com/jarulsamy/ps_scraper)
   - CLI tool that exports grades to Excel
   - Query any PowerSchool website

#### Schoology
1. **SchoologyMessageWebScraper** (https://github.com/Saptak625/SchoologyMessageWebScraper)
   - Uses `requests.Session()` to maintain login
   - Extracts messages from Schoology
   - Works with standard credentials

2. **sgy-sgy** (https://github.com/SheepTester/sgy-sgy)
   - Uses session cookies from browser
   - Requires: HOST, UID, CSRF_KEY, CSRF_TOKEN, SESS_ID
   - Can be extracted from Chrome DevTools

---

## Two Viable Approaches

### Approach 1: Session Token Extraction (Recommended for MVP)

**How it works:**
1. User logs into Schoology/PowerSchool via their regular browser
2. User extracts session cookies using browser DevTools or extension
3. User provides cookies to Trunchbull via configuration
4. Trunchbull uses these cookies to make API requests
5. When cookies expire, user repeats the process

**Advantages:**
- ✅ Simple to implement
- ✅ No reverse engineering required
- ✅ User controls authentication
- ✅ Minimal Terms of Service concerns
- ✅ No password storage

**Disadvantages:**
- ❌ Requires manual cookie refresh (every few days/weeks)
- ❌ Less convenient for users
- ❌ Need to document the cookie extraction process

**Implementation:**
```yaml
# User extracts these from browser and adds to config
schoology:
  host: "meanyms.schoology.com"
  sess_id: "abc123..."
  csrf_token: "xyz789..."
  uid: "12345"

powerschool:
  pstoken: "token_value..."
  session_cookie: "sess123..."
```

### Approach 2: Credential-Based Automation (Better UX)

**How it works:**
1. User provides username and password (encrypted at rest)
2. Trunchbull uses headless browser or HTTP client to log in
3. Obtains session cookies programmatically
4. Automatically refreshes session when needed
5. Fully automated synchronization

**Advantages:**
- ✅ Better user experience (set and forget)
- ✅ Automatic session refresh
- ✅ No manual cookie management
- ✅ Works like official mobile apps

**Disadvantages:**
- ❌ Must store credentials (even if encrypted)
- ❌ More complex implementation
- ❌ Requires understanding login flow
- ❌ May break if login process changes
- ❌ Higher ToS risk (automated access)

**Implementation Options:**

**Option A: Headless Browser (Selenium)**
```go
// Use Selenium WebDriver
// Navigate to login page
// Fill credentials
// Extract cookies
// Use for API requests
```

**Option B: Reverse-Engineered Authentication**
```go
// Analyze login POST request
// Replicate authentication flow
// Handle CSRF tokens
// Maintain session
```

---

## Technical Implementation Details

### For Schoology

#### Session Cookies Needed:
- `SESS<hash>` - Main session cookie
- `CSRF_KEY` - CSRF protection key
- `CSRF_TOKEN` - CSRF protection token
- `UID` - User ID
- `HOST` - Schoology subdomain

#### API Endpoints:
Schoology's web interface uses internal API endpoints that don't require OAuth:
```
GET https://{school}.schoology.com/iapi2/site-navigation/courses
GET https://{school}.schoology.com/iapi2/gradebook/{section_id}
GET https://{school}.schoology.com/iapi2/assignment/{assignment_id}
```

#### Authentication Flow:
1. POST to login endpoint with credentials
2. Receive session cookies
3. Include cookies + CSRF token in subsequent requests
4. Session typically lasts 7-14 days

### For PowerSchool (The Source)

#### Session Cookies Needed:
- `pstoken` - Main authentication token
- Session cookies from login

#### Internal API Endpoints:
PowerSchool's parent portal uses JSON APIs:
```
GET /guardian/home.html
GET /ws/schema/query/com.pearson.core.student.grades
GET /ws/v1/district/grade/{id}
GET /ws/v1/district/attendance/{student_id}
```

#### Authentication Flow:
1. POST credentials to login endpoint
2. Receive authentication tokens
3. Use tokens for subsequent API calls
4. PowerSchool Mobile app uses OAuth-style tokens

---

## Recommended Architecture Changes

### Modified Tech Stack

**Backend:**
- Go for API gateway (unchanged)
- **Add**: Chromedp or Rod (Go headless browser library) OR
- **Add**: Colly (Go web scraping framework)

**For Session Management:**
```go
type SessionManager struct {
    cookies    map[string]*http.Cookie
    csrfToken  string
    userID     string
    expiresAt  time.Time
}

func (sm *SessionManager) RefreshIfNeeded() error {
    if time.Now().After(sm.expiresAt) {
        return sm.Login()
    }
    return nil
}
```

### Configuration Updates

```yaml
# New configuration structure
authentication:
  method: "session"  # or "credentials"

# Method 1: Session tokens
schoology_session:
  host: ""
  sess_id: ""
  csrf_token: ""
  csrf_key: ""
  uid: ""
  expires_at: ""

# Method 2: Credentials (encrypted)
schoology_credentials:
  username: ""  # Encrypted
  password: ""  # Encrypted
  auto_refresh: true

powerschool_session:
  pstoken: ""
  session_cookie: ""

powerschool_credentials:
  username: ""
  password: ""
  auto_refresh: true
```

---

## Security Considerations

### Credential Storage
If using Approach 2 (credentials):
- **MUST** encrypt credentials at rest
- Use strong encryption (AES-256)
- Store encryption key separately (environment variable)
- Document security implications clearly

```go
import "golang.org/x/crypto/nacl/secretbox"

func EncryptCredentials(username, password, key string) ([]byte, error) {
    // Implement proper encryption
}
```

### Session Token Storage
If using Approach 1 (sessions):
- Still encrypt session tokens
- Tokens are bearer tokens - same risk as credentials
- Shorter expiration reduces risk

### Network Security
- **MUST** use HTTPS for all requests
- Implement TLS certificate validation
- No unencrypted data transmission

---

## Legal & Terms of Service

### Is This Legal?

**Yes, with important caveats:**

1. **You own the right to access your child's educational records** (FERPA)
2. **You're using your own credentials** (not unauthorized access)
3. **Similar to how mobile apps work** (using parent portal access)
4. **For personal, non-commercial use**

### Terms of Service Considerations

**Gray Areas:**
- Most ToS prohibit "automated access" or "bots"
- BUT: Official mobile apps use the same approach
- Personal use for your own data is generally acceptable
- Don't abuse rate limits or scrape bulk data

**Seattle Public Schools Specifics:**
- Took action against "The Source: SPS" when it was **monetized**
- App was fine for 3 years as a free student project
- Issue was commercialization + use of district name/logo
- Personal use was never the concern

**Best Practices:**
- Keep it personal (family use only)
- Don't share credentials
- Respect rate limits
- Don't redistribute/commercialize
- Be prepared to stop if asked

---

## Implementation Recommendations

### Phase 1: MVP (Recommended)
**Use Session Token Approach**

Why:
- Faster to implement
- Lower legal/ToS risk
- User controls authentication
- Good proof of concept

Implementation:
1. Create session configuration structure
2. Build HTTP client with cookie support
3. Implement API wrappers
4. Document cookie extraction process
5. Build basic dashboard

### Phase 2: Enhanced UX
**Add Credential-Based Option**

Why:
- Better user experience
- Automated refresh
- More like production apps

Implementation:
1. Add headless browser support (Chromedp)
2. Implement login automation
3. Secure credential storage
4. Automatic session refresh
5. Fallback to manual session

### Phase 3: Polish
**Hybrid Approach**

Allow users to choose:
- Session tokens (manual, more control)
- Credentials (automatic, convenience)
- OAuth (if district provides access)

---

## Updated Project Structure

```
internal/
├── auth/
│   ├── session/         # Session token management
│   │   ├── schoology.go
│   │   └── powerschool.go
│   ├── browser/         # Headless browser automation
│   │   ├── login.go
│   │   └── extract.go
│   └── credentials/     # Credential encryption/storage
│       └── secure.go
├── client/
│   ├── schoology/       # Schoology HTTP client
│   │   ├── client.go
│   │   ├── assignments.go
│   │   ├── grades.go
│   │   └── messages.go
│   └── powerschool/     # PowerSchool HTTP client
│       ├── client.go
│       ├── grades.go
│       └── attendance.go
```

---

## Required Go Libraries

```go
// For HTTP client with cookie support
"net/http"
"net/http/cookiejar"

// For headless browser (optional, Phase 2)
"github.com/chromedp/chromedp"

// For web scraping (if needed)
"github.com/gocolly/colly/v2"

// For encryption
"golang.org/x/crypto/nacl/secretbox"
"golang.org/x/crypto/argon2"

// For HTML parsing
"github.com/PuerkitoBio/goquery"
```

---

## Documentation Needs

### User Documentation
1. **Cookie Extraction Guide**
   - Step-by-step with screenshots
   - For Chrome, Firefox, Safari
   - Using DevTools
   - Automated extraction tools

2. **Security Best Practices**
   - Why cookies must be protected
   - Not sharing credentials
   - Secure home lab setup

3. **Troubleshooting**
   - Session expired
   - Invalid cookies
   - Login changes

### Developer Documentation
1. **Authentication Flow**
2. **API Endpoints**
3. **Session Management**
4. **Error Handling**

---

## Testing Strategy

### Manual Testing
1. Extract session from real account
2. Test API calls with session
3. Verify data returned
4. Test session expiration

### Automated Testing
```go
func TestSchoologySession(t *testing.T) {
    session := &SchoologySession{
        SessID: os.Getenv("TEST_SESS_ID"),
        CSRF:   os.Getenv("TEST_CSRF"),
    }

    client := NewSchoologyClient(session)
    courses, err := client.GetCourses()
    // assertions
}
```

---

## Comparison: Official API vs Credential-Based

| Aspect | Official API | Credential-Based |
|--------|-------------|------------------|
| **District Approval** | Required | Not required |
| **Setup Complexity** | High | Low |
| **Rate Limits** | Defined | Undefined (be conservative) |
| **Documentation** | Official | Reverse engineered |
| **Stability** | High | Medium (can break) |
| **Legal/ToS** | Fully compliant | Gray area |
| **Access Scope** | District-level possible | Personal only |
| **Maintenance** | Lower | Higher |

---

## Next Steps

1. **Immediate:**
   - Build session-based authentication for Schoology
   - Test with real parent account
   - Document cookie extraction process

2. **Short-term:**
   - Build session-based authentication for PowerSchool
   - Create HTTP clients for both platforms
   - Implement basic data retrieval

3. **Medium-term:**
   - Add credential-based option
   - Implement automatic session refresh
   - Build monitoring for session expiration

4. **Long-term:**
   - Consider official API as fallback
   - Support multiple auth methods
   - Community contributions for other platforms

---

## Risks & Mitigation

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Login process changes | Medium | High | Version detection, fallback methods |
| Session invalidation | High | Medium | Automatic re-auth, user notification |
| Terms of Service violation | Low | High | Personal use only, clear documentation |
| Security breach | Low | Critical | Encryption, security best practices |
| District blocks access | Low | High | Respect rate limits, good citizenship |

---

## Conclusion

**The credential-based approach is viable and recommended** for Trunchbull's initial release. It:

- ✅ Eliminates need for district API approval
- ✅ Works with standard parent/student credentials
- ✅ Has proven viability (multiple existing projects)
- ✅ Provides same data as official portals
- ✅ Suitable for personal, family use

**Recommended Path Forward:**
1. Start with session token approach (MVP)
2. Add credential automation later (better UX)
3. Keep official OAuth as future option

This approach makes Trunchbull immediately usable without bureaucratic barriers while maintaining security and respecting terms of service for personal use.

---

**Document Version**: 2.0
**Last Updated**: 2025-10-23
**Status**: Ready for Implementation
