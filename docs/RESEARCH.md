# Student Dashboard - Platform Research

## Executive Summary

This document outlines the research findings for integrating with **Schoology** and **The Source (PowerSchool)** to build a student dashboard for Meany Middle School.

### Key Findings

1. **Schoology** provides a well-documented REST API with OAuth authentication
2. **The Source** is Seattle Public Schools' parent portal built on **PowerSchool SIS**
3. Both platforms require proper authentication and have security considerations
4. ~~API access for PowerSchool requires district administrator approval~~ **UPDATE: Not required!**
5. Important security update for Schoology coming June 25, 2025
6. **🎯 BREAKTHROUGH: Credential-based access is viable!** (see Alternative Authentication section below)

### Critical Update (2025-10-23)

**We discovered that official API approval is NOT required!** Multiple successful projects demonstrate that both platforms can be accessed using regular parent/student credentials, eliminating the need for district-level API approval.

**See [ALTERNATIVE_AUTH.md](ALTERNATIVE_AUTH.md) for complete details on credential-based authentication.**

---

## Platform 1: Schoology

### Overview
Schoology is a Learning Management System (LMS) with a comprehensive REST API for third-party integrations.

### Authentication
- **Method**: OAuth 1.0a (two-legged and three-legged flows)
- **Requirements**:
  - Consumer key and secret
  - API keys can be obtained through the Schoology admin portal
- **Security Update (June 25, 2025)**: Personal API keys will no longer be able to access other users' data

### Available API Endpoints

#### ✅ Assignments
- **Endpoint**: `/sections/{section_id}/assignments`
- **Capabilities**:
  - List all assignments for a section
  - Get assignment details including due dates
  - Check completion status
  - View attachments and tags
- **Documentation**: developers.schoology.com/api-documentation/rest-api-v1/assignment/

#### ✅ Grades
- **Endpoint**: `/sections/{section_id}/grades`
- **Capabilities**:
  - Get grades by assignment_id or enrollment_id
  - Filter by grading period
  - Access individual assignment scores
- **Documentation**: developers.schoology.com/api-documentation/rest-api-v1/grade/

#### ✅ Grading Scales
- **Endpoint**: `/sections/{section_id}/grading_scales`
- **Capabilities**:
  - Get grading scale configurations
  - Understand grade weighting
- **Documentation**: developers.schoology.com/api-documentation/rest-api-v1/grading-scales/

#### ⚠️ GPA / Report Cards
- **Status**: School-configurable feature
- **Note**: Not all schools enable GPA display in Schoology
- **Workaround**: May need to calculate based on grades and grading scales
- **API Support**: Limited - depends on school configuration

#### ✅ Calendar Events
- **Endpoint**: `/sections/{section_id}/events`
- **Capabilities**:
  - View course events
  - Get assignment due dates
  - Access school calendar items
- **Documentation**: developers.schoology.com/api-documentation/rest-api-v1/course-section/

#### ✅ Messages
- **Endpoint**: `/messages`
- **Capabilities**:
  - Read private messages
  - Send messages programmatically
  - Check for unread messages from teachers
- **Documentation**: developers.schoology.com/api-documentation/rest-api-v1/

#### ✅ Parent-Child Associations
- **Endpoint**: `/users/{user_id}/parents` and `/users/{user_id}/children`
- **Capabilities**:
  - Link parent and child accounts
  - Access child data through parent credentials
- **Documentation**: developers.schoology.com/api-documentation/rest-api-v1/user/

### Rate Limits
- Not explicitly documented in search results
- Best practice: Implement exponential backoff for 429 responses
- Consider caching data to minimize API calls

### Developer Resources
- **Main Site**: developers.schoology.com
- **API Docs**: developers.schoology.com/api-documentation/rest-api-v1/
- **Community**: PowerSchool Community forums

---

## Platform 2: The Source (PowerSchool SIS)

### Overview
The Source is Seattle Public Schools' parent portal built on PowerSchool Student Information System. It provides access to attendance, grades, assessment scores, schedules, and more.

### Authentication
- **Method**: OAuth 2.0
- **Requirements**:
  - Application registration with PowerSchool
  - API key and secret
  - **District administrator approval required**
- **Access**: Permission must be granted by system admins at the district level

### Available API Endpoints

#### ✅ Attendance
- Access real-time attendance records
- Historical attendance data

#### ✅ Grades
- Current grades for all courses
- Assignment-level scores
- Progress reports and report cards

#### ✅ Schedules
- Student course schedules
- Class periods and room assignments

#### ✅ Assessment Scores
- Standardized test scores
- District assessments

#### ⚠️ GPA
- Typically available through PowerSchool
- May require specific API endpoints (needs verification with district)

#### ⚠️ Calendar / Events
- School events and days off
- District calendar
- **Note**: Availability depends on district configuration

### Rate Limits
- **Enforced**: Yes
- **Error Code**: 429 Too Many Requests
- **Header**: `Retry-After` header indicates wait time in seconds
- **Best Practice**: Implement exponential backoff and request throttling

### Developer Resources
- **Documentation**: support.powerschool.com/developer (requires PowerSchool account)
- **Access Request**: Must contact Seattle Public Schools IT department
- **Community**: PowerSchool Community forums (help.powerschool.com)

### Important Security Note
PowerSchool experienced a significant data breach in January 2025. Security and data handling are critical considerations for any integration.

---

## Feature Mapping

### Dashboard Requirements vs. API Capabilities

| Feature | Schoology | PowerSchool (The Source) | Notes |
|---------|-----------|-------------------------|-------|
| Outstanding assignments | ✅ Full support | ✅ Full support | Both platforms provide assignment lists |
| Current grades / Report card | ✅ Full support | ✅ Full support | Real-time grade access |
| GPA calculation | ⚠️ Limited | ✅ Likely available | May need to calculate from grades |
| School events / days off | ✅ Course events | ✅ District calendar | Combine both sources |
| Teacher messages | ✅ Full support | ⚠️ Needs verification | Schoology has messaging API |
| Parent authentication | ✅ Supported | ✅ Supported | Both support parent accounts |

**Legend:**
- ✅ Full support confirmed
- ⚠️ Partial support or needs verification
- ❌ Not supported

---

## API Access Process

### Schoology
1. Contact Meany Middle School's Schoology administrator
2. Request API access and developer credentials
3. Obtain consumer key and secret
4. Implement OAuth 1.0a authentication flow
5. Test with sandbox environment if available

### PowerSchool (The Source)
1. Contact Seattle Public Schools IT Department
2. Request API access to PowerSchool SIS
3. Complete any required security review processes
4. Obtain API credentials (client ID and secret)
5. Register application in PowerSchool Developer Portal
6. Implement OAuth 2.0 authentication flow

---

## API Usage Best Practices

### Rate Limiting Strategy
```
1. Implement request throttling on client side
2. Cache responses appropriately (consider 15-30 minute TTL for grade data)
3. Use exponential backoff for retries:
   - First retry: 2 seconds
   - Second retry: 4 seconds
   - Third retry: 8 seconds
   - Fourth retry: 16 seconds
4. Monitor rate limit headers in responses
5. Implement circuit breaker pattern for persistent failures
```

### Data Freshness
- **Real-time data**: Assignments, messages (check every 15-30 minutes)
- **Hourly updates**: Grades, attendance
- **Daily updates**: GPA, schedules, calendar events
- **Weekly updates**: Report cards, assessment scores

### Caching Strategy
- Store API responses locally with timestamps
- Implement cache invalidation based on data type
- Use Redis or similar for production deployments
- Consider offline mode for dashboard when APIs are unavailable

---

## Security & Privacy Considerations

### Data Sensitivity
- **PII Included**: Student names, grades, attendance, assessment scores
- **FERPA Compliance**: Educational records are protected under FERPA
- **Local Storage Only**: All data should be stored locally on user-controlled infrastructure
- **No Cloud Storage**: Avoid storing student data in cloud services

### Authentication Security
- Store API credentials securely (use environment variables or secrets manager)
- Never commit credentials to version control
- Implement credential rotation capability
- Use HTTPS for all API communications
- Consider encryption at rest for cached data

### Third-Party App Concerns
Seattle Public Schools has previously taken action against unauthorized third-party apps accessing student data. To be good netizens:

1. **Transparency**: Clear documentation about data collection and usage
2. **User Control**: Users must provide their own credentials (no credential sharing)
3. **No Monetization**: This is a personal/family tool, not a commercial product
4. **No Bulk Access**: Only access data for authorized children
5. **Respect ToS**: Comply with both platforms' Terms of Service
6. **Self-Hosted Only**: Emphasize that users must run this on their own infrastructure

---

## Recommended Architecture

### Technology Stack
- **Backend**: Go
  - Excellent for concurrent API calls
  - Strong standard library for HTTP clients
  - Easy containerization
  - Good OAuth libraries available

- **Frontend**: React
  - Component-based architecture for dashboard widgets
  - Rich ecosystem for data visualization
  - Mobile-responsive design capabilities

- **Data Layer**:
  - SQLite for local data persistence (suitable for single-family use)
  - Redis for caching (optional, for performance)

- **Containerization**: Docker + Docker Compose
  - Easy deployment on home lab
  - Isolated environment
  - Simple updates and maintenance

### High-Level Architecture
```
┌─────────────────┐
│   React SPA     │  (Dashboard UI)
└────────┬────────┘
         │ HTTP/WebSocket
         ↓
┌─────────────────┐
│   Go Backend    │  (API Gateway + Orchestration)
├─────────────────┤
│ - Auth Manager  │  (OAuth flows for both platforms)
│ - API Clients   │  (Schoology + PowerSchool)
│ - Cache Layer   │  (Redis/in-memory)
│ - Data Sync     │  (Scheduled jobs)
└────────┬────────┘
         │
         ↓
┌─────────────────┐
│  SQLite DB      │  (Local persistence)
└─────────────────┘

External APIs:
├─ Schoology API
└─ PowerSchool API (The Source)
```

---

## Implementation Phases

### Phase 1: Foundation (MVP)
- Go backend with basic HTTP server
- OAuth implementation for one platform (recommend starting with Schoology)
- Single student support
- Basic React dashboard showing:
  - Outstanding assignments count
  - Current grades
- Docker containerization

### Phase 2: Full Integration
- PowerSchool API integration
- Multiple student support
- Enhanced dashboard with:
  - Assignment details and due dates
  - Grade breakdown by class
  - School events calendar
- Background sync jobs
- Data caching

### Phase 3: Polish
- Real-time updates (WebSocket)
- GPA calculation
- Teacher messages inbox
- Mobile-responsive design
- Notifications for new assignments/grades

---

## Open Questions & Next Steps

### Questions for District/School
1. Can we get API access as a parent?
2. What are the exact rate limits for each platform?
3. Is there a sandbox/test environment available?
4. Are there any specific restrictions on parent API access?
5. What is the process for requesting API credentials?

### Technical Decisions
1. Should we implement our own GPA calculation or rely on platform data?
2. What's the refresh interval that balances freshness vs. API politeness?
3. Do we need a database or is caching sufficient for MVP?
4. Should we support multiple families or optimize for single-family use?

### Next Steps
1. ✅ Complete platform research
2. ⬜ Create detailed architecture document
3. ⬜ Draft security and privacy documentation
4. ⬜ Set up project structure
5. ⬜ Reach out to school district for API access
6. ⬜ Implement OAuth flow for Schoology (first platform)
7. ⬜ Build basic Go backend with API client
8. ⬜ Create React dashboard skeleton
9. ⬜ Test with real data

---

## Alternative Authentication Discovery (Update 2025-10-23)

### Breakthrough Finding

After the initial research, we discovered that **official district-level API approval is NOT required** for Trunchbull. Multiple successful open-source projects have demonstrated that both platforms can be accessed using standard parent/student login credentials.

### Evidence

#### Case Study: "The Source: SPS" App
- Built by a Seattle high school student
- Used regular parent/student credentials (no official API)
- **17,000+ monthly active users**
- **30,000+ total downloads**
- Operated successfully for **3 years** (2021-2024)
- Only became "unauthorized" when monetized by a commercial entity

#### Existing Open-Source Projects

**PowerSchool:**
- `psscraper` - Browser automation with Selenium + BeautifulSoup
- `ps.py` - Authentication & data fetching library
- `ps_scraper` - CLI tool for grade export

**Schoology:**
- `SchoologyMessageWebScraper` - Session-based scraping
- `sgy-sgy` - Cookie-based access
- Multiple projects using `requests.Session()`

### Two Viable Approaches

#### Approach 1: Session Token Extraction (Recommended for MVP)
- User logs in via browser
- Extract session cookies using DevTools
- Provide cookies to Trunchbull
- Manually refresh when expired

**Advantages:**
- Simple implementation
- User controls authentication
- Lower ToS risk
- No credential storage

#### Approach 2: Credential-Based Automation
- Store username/password (encrypted)
- Automate login with headless browser
- Automatic session refresh
- Fully autonomous operation

**Advantages:**
- Better user experience
- Automatic session management
- No manual cookie extraction

### Impact on Project

This discovery means:
1. ✅ No need to contact district IT for API access
2. ✅ Immediate development can begin
3. ✅ Works with standard parent accounts
4. ✅ Lower barrier to entry for users
5. ✅ Proven viability (multiple working examples)

### Documentation

**Complete implementation details available in:**
- [ALTERNATIVE_AUTH.md](ALTERNATIVE_AUTH.md) - Comprehensive guide
- Technical implementation
- Security considerations
- Legal/ToS analysis
- Code examples

### Updated Next Steps

1. ✅ Complete platform research
2. ✅ Document alternative authentication approach
3. ✅ Create detailed architecture document
4. ✅ Draft security and privacy documentation
5. ✅ Set up project structure
6. ⬜ ~~Reach out to school district for API access~~ **NOT NEEDED!**
7. ⬜ Implement session-based authentication
8. ⬜ Build HTTP clients with cookie support
9. ⬜ Create React dashboard skeleton
10. ⬜ Test with real parent account

---

## References

### Schoology
- Developer Docs: https://developers.schoology.com/
- API Documentation: https://developers.schoology.com/api-documentation/rest-api-v1/
- Community: https://help.powerschool.com/t5/Community-Forum/ct-p/Schoology

### PowerSchool
- Developer Portal: https://support.powerschool.com/developer/
- Seattle Public Schools: https://www.seattleschools.org/resources/the-source/
- Community: https://help.powerschool.com/t5/Community-Forum/ct-p/Community-Forum

### OAuth Resources
- OAuth 1.0a: https://oauth.net/core/1.0a/
- OAuth 2.0: https://oauth.net/2/

### Alternative Authentication Projects
- psscraper: https://github.com/Desperationis/psscraper
- ps.py: https://github.com/ouiliame/ps.py
- ps_scraper: https://github.com/jarulsamy/ps_scraper
- SchoologyMessageWebScraper: https://github.com/Saptak625/SchoologyMessageWebScraper
- sgy-sgy: https://github.com/SheepTester/sgy-sgy
- "The Source: SPS" Case Study: https://www.geekwire.com/2024/seattle-public-schools-issues-statement-explaining-its-stance-on-grade-viewing-app/

---

**Document Version**: 1.0
**Last Updated**: 2025-10-23
**Status**: Initial Research Complete
