# LinkedIn Automation

A sophisticated Go application for automating LinkedIn login, people search, and interactions with human-like behavior simulation using the Rod browser automation library.

## Overview

This project uses [Rod](https://go-rod.github.io/) to automate browser interactions with LinkedIn. It features advanced stealth techniques including random mouse movements, natural scrolling patterns, realistic typing speeds, and comprehensive anti-detection measures to avoid LinkedIn's bot detection systems.

## Features

- **Automated LinkedIn Login** - Logs into LinkedIn with email/password credentials
- **2FA/CAPTCHA Detection** - Pauses for manual intervention when challenges are detected
- **LinkedIn People Search**:
  - Filter by keywords, job title, company, and location
  - Pagination support (scrape multiple pages)
  - Duplicate detection (skip profiles visited in last 30 days)
  - Location-based targeting (50+ supported locations worldwide)
  - Profile data extraction (name, title, company, location, connection degree)
  - Database storage for persistent profile tracking
- **Human-like Behavior Simulation**:
  - Bézier curve mouse movements (natural acceleration/deceleration)
  - Random page scrolling with reading pauses
  - Realistic typing speeds with variance
  - Random delays between actions (simulating human reaction time)
  - Hover over interactive elements (2-3 random hovers per page)
- **Advanced Stealth Techniques**:
  - Comprehensive fingerprint masking (navigator.webdriver, canvas, WebGL, etc.)
  - Business hours scheduling (only run 9 AM - 5 PM weekdays)
  - Rate limiting (14 connections/day, 50 messages/day, 100 searches/day)
  - Cooldown periods (30 seconds between actions)
  - Persistent browser sessions (avoid repeated logins)
- **SQLite Database** - Tracks profiles, connections, messages, and rate limits
- **Session Persistence** - Saves browser state to skip login on subsequent runs
- **Comprehensive Logging** - Detailed logs of all actions performed
- **Environment-based Configuration** - Secure credential and settings management

## Prerequisites

- **Go 1.24.5** or higher
- **Chrome/Chromium browser** (Rod will auto-download if not present)
- **SQLite3** (included with Go)
- **.env file** with LinkedIn credentials and search configuration

## Installation

1. Clone the repository:
```bash
git clone https://github.com/YOUR_USERNAME/linkedin-automation.git
cd linkedin-automation
```

2. Install dependencies:
```bash
go mod tidy
```

3. Create a `.env` file in the project root:
```env
# LinkedIn Credentials
LINKEDIN_EMAIL=your_email@example.com
LINKEDIN_PASSWORD=your_password

# Database Configuration
DATABASE_PATH=./data/linkedin_automation.db

# Rate Limits
MAX_CONNECTIONS_PER_DAY=14
MAX_MESSAGES_PER_DAY=50
MAX_SEARCHES_PER_DAY=100
COOLDOWN_SECONDS=30

# Activity Scheduling
ACTIVE_HOURS_START=9
ACTIVE_HOURS_END=17
WEEKDAYS_ONLY=true

# Session Configuration
SESSION_VALIDITY_DAYS=7
BROWSER_DATA_DIR=./browser_data

# Search Configuration
SEARCH_KEYWORDS=software engineer
SEARCH_JOB_TITLE=
SEARCH_COMPANY=
SEARCH_LOCATION=San Francisco Bay Area
```

4. Build the project (optional):
```bash
go build
```

## Quick Start (5 Minutes)

### Prerequisites Check
```bash
# Verify Go is installed
go version  # Should be 1.24.5 or higher

# Verify you're in the project directory
cd linkedin-automation
ls -la  # Should show main.go, README.md, go.mod, etc.
```

### Setup Steps
```bash
# Step 1: Download dependencies
go mod tidy

# Step 2: Create credentials file (.env)
# On Windows (PowerShell):
@"
LINKEDIN_EMAIL=your_email@example.com
LINKEDIN_PASSWORD=your_password
"@ | Out-File -Encoding UTF8 .env

# On Linux/Mac:
cat > .env << EOF
LINKEDIN_EMAIL=your_email@example.com
LINKEDIN_PASSWORD=your_password
EOF

# Step 3: Verify .env file exists
ls -la .env  # Should show .env file

# Step 4: Run the application
go run main.go
```

### What to Expect (Timeline)
- **0-2 seconds**: "Starting LinkedIn Automation" log appears
- **2-3 seconds**: Browser window launches (Chrome opens)
- **3-5 seconds**: Navigates to LinkedIn login page
- **5-15 seconds**: Enters credentials and logs in (watch it type)
- **15-20 seconds**: "Login Successful" message appears
- **20-25 seconds**: Mouse movements visible on page (cursor moving)
- **25-30 seconds**: Page scrolls up and down automatically
- **30+ seconds**: Browser stays open showing home page

## Usage

Run the application:
```bash
go run main.go
```

### What the Application Does

1. **Loads credentials** from the `.env` file
2. **Launches a Chrome browser** in non-headless mode (visible)
3. **Navigates to LinkedIn login page**
4. **Performs automated login** with human-like delays
5. **Verifies successful login** by checking page URL
6. **Applies fingerprint masking** to avoid detection
7. **Simulates human behavior**:
   - Random mouse movements across the page
   - Random scrolling with natural pauses
   - Realistic timing between all actions
8. **Keeps browser open** for 30+ seconds showing the results

### Expected Output

```
[2025-12-24 19:08:28] INFO: Starting LinkedIn Automation
[2025-12-24 19:08:29] INFO: Opening Linkedin Login page
[2025-12-24 19:08:30] INFO: Locating email input field
[2025-12-24 19:08:31] INFO: Locating password input field
[2025-12-24 19:08:32] INFO: Locating and clicking on sign in button
[2025-12-24 19:08:33] INFO: Checking if login was successful
[2025-12-24 19:08:34] INFO: Login Successful - Redirected to home page
[2025-12-24 19:08:35] INFO: Applying fingerprint masking...
[2025-12-24 19:08:36] INFO: Starting human-like behavior simulation...
[2025-12-24 19:08:37] INFO: Performing stealth actions - simulating human-like behavior
[2025-12-24 19:08:38] INFO: Executing random mouse movements...
[2025-12-24 19:08:39] INFO: Executing random page scrolling...
[2025-12-24 19:08:40] INFO: Stealth actions completed
```

## First Run Checklist

Before running the application, verify:

- [ ] Go 1.24.5+ is installed: `go version`
- [ ] You're in the project directory: `pwd` or `cd linkedin-automation`
- [ ] Dependencies are installed: `go mod tidy`
- [ ] `.env` file exists with credentials: `cat .env`
- [ ] `.env` has valid LinkedIn email and password
- [ ] Chrome/Chromium is installed (or let Rod download it)
- [ ] No other applications are using port 9222 (Rod's default debugging port)
- [ ] Windows Defender isn't blocking Rod (add exclusion if needed)

## Build & Run Commands

### Development (with hot reload)
```bash
# Install go-task or use go run directly
go run main.go
```

### Production Build
```bash
# Build executable
go build -o linkedin-automation.exe

# Run the executable
./linkedin-automation.exe
```

### Build with Optimizations
```bash
# Smaller binary size
go build -ldflags="-s -w" -o linkedin-automation.exe

# Run it
./linkedin-automation.exe
```

### Cross-Platform Building
```bash
# Build for Windows (if on Windows)
go build -o linkedin-automation.exe

# Build for Linux (if on Linux/Mac)
go build -o linkedin-automation

# Build for macOS (if on Mac)
go build -o linkedin-automation

# Cross-compile Windows binary from Linux/Mac
GOOS=windows GOARCH=amd64 go build -o linkedin-automation.exe
```

## Verification Steps (Testing Actual Functionality)

### Step 1: Verify Compilation
```bash
go build
# Should complete with no errors
ls -la linkedin-automation.exe  # File should exist
```

### Step 2: Run and Verify Each Feature

**A. Verify Credential Loading**
```bash
go run main.go
# Check logs for: "Starting LinkedIn Automation"
# Check that browser launches
```

**B. Verify LinkedIn Login**
- Watch the browser window
- See email field being filled automatically
- See password field being filled automatically
- See the "Sign in" button being clicked
- Wait for page to load and check logs for "Login Successful"

**C. Verify Mouse Movements** ✅
- After login, watch the **mouse cursor** on the page
- You should see the cursor move to **3-5 different positions**
- Each movement is followed by a **300-800ms pause**
- Check logs for: "Executing random mouse movements..."

**D. Verify Page Scrolling** ✅
- After mouse movements, watch the **page scroll**
- Page should scroll down **3-5 times**
- Each scroll is 200-600 pixels
- Followed by 800-1500ms pause (reading time)
- Check logs for: "Executing random page scrolling..."

**E. Verify Logging**
- Check that **every action is logged**
- Look for log format: `[TIMESTAMP] INFO: Action name`
- Verify no ERROR messages in logs

### Step 3: Test with Different Credentials
```bash
# Edit .env with test account
# Run again
go run main.go

# Verify it works with different credentials
```

### Step 4: Test Build Executable
```bash
# Build the project
go build

# Run the executable instead of 'go run'
./linkedin-automation.exe

# Should work identically to 'go run main.go'
```

### Success Indicators ✅
- [x] Browser launches automatically
- [x] Credentials are entered without user interaction
- [x] Login completes successfully
- [x] Mouse movements are visible on the page
- [x] Page scrolling occurs automatically
- [x] All actions are logged with timestamps
- [x] No error messages in console
- [x] Program completes without crashing
- [x] Browser closes cleanly when done

## Project Structure

```
linkedin-automation/
├── main.go                    # Entry point - orchestrates the automation workflow
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
├── .env                       # Environment variables (credentials + config)
├── .env.example               # Example environment configuration
├── README.md                  # This file
│
├── data/
│   └── linkedin_automation.db # SQLite database (auto-created)
│
├── browser_data/              # Persistent browser session data (auto-created)
│
├── internal/
│   ├── automation/
│   │   ├── login.go           # LinkedIn login automation with 2FA detection
│   │   ├── search.go          # People search with pagination & duplicate detection
│   │   ├── scheduler.go       # Business hours scheduling (9 AM - 5 PM weekdays)
│   │   └── ratelimiter.go     # Rate limiting (connections, messages, searches)
│   │
│   ├── browser/
│   │   ├── browser.go         # Browser initialization with UserDataDir
│   │   └── fingerprint.go     # Anti-detection fingerprint masking (10+ techniques)
│   │
│   ├── logger/
│   │   └── logger.go          # Centralized logging utility
│   │
│   ├── stealth/
│   │   ├── delay.go           # Random delay generation for human-like timing
│   │   ├── mouse.go           # Bézier curve mouse movements + element hovering
│   │   ├── scroll.go          # Random page scrolling with natural patterns
│   │   └── typing.go          # Human-like typing speed simulation
│   │
│   └── storage/
│       ├── database.go        # SQLite database operations (profiles, connections, etc.)
│       └── state.go           # Session state persistence (JSON)
│
└── pkg/
    ├── models/
    │   └── models.go          # Public data models (Profile, Connection, Message)
    └── utils/
        ├── constants.go       # LinkedIn URLs, selectors, location URN codes
        ├── helpers.go         # Utility helper functions
        └── validators.go      # Input validation functions

tests/                         # Integration and unit tests
```

## Search Configuration

### Supported Locations

The search module supports 50+ locations worldwide (see `pkg/utils/constants.go`):

**United States:**
- Major cities: San Francisco Bay Area, New York City Area, Los Angeles, Chicago, Boston, Seattle, Austin, etc.
- States: California, New York, Texas, Florida, etc.
- Country-wide: United States

**International:**
- Countries: United Kingdom, Canada, Germany, France, India, Australia, Singapore, Japan, etc.
- Major cities: London, Toronto, Berlin, Paris, Sydney, Bangalore, Tokyo, Dubai, etc.

### Search Parameters

Configure search in your `.env` file:

```env
# Required - at least one must be set
SEARCH_KEYWORDS=software engineer    # General keywords
SEARCH_JOB_TITLE=Senior Engineer     # Filter by job title
SEARCH_COMPANY=Google                # Filter by company name

# Location filter (must match location name exactly)
SEARCH_LOCATION=San Francisco Bay Area

# The system will:
# - Search LinkedIn for profiles matching your criteria
# - Extract profile data (name, title, company, location)
# - Check for duplicates (profiles visited in last 30 days)
# - Save new profiles to SQLite database
# - Display statistics (total found, new saved, duplicates skipped)
```

### Search Examples

**Example 1: Find AI Researchers in San Francisco**
```env
SEARCH_KEYWORDS=artificial intelligence machine learning
SEARCH_JOB_TITLE=Research Scientist
SEARCH_LOCATION=San Francisco Bay Area
```

**Example 2: Find CTOs at Startups in New York**
```env
SEARCH_KEYWORDS=startup founder
SEARCH_JOB_TITLE=CTO
SEARCH_LOCATION=New York City Area
```

**Example 3: Find Product Managers at Google**
```env
SEARCH_KEYWORDS=product manager
SEARCH_COMPANY=Google
SEARCH_LOCATION=United States
```

### Search Output

The search module will:
1. Navigate to LinkedIn people search with your filters
2. Extract profiles from each page (up to 3 pages by default)
3. Check database for duplicates (profiles visited in last 30 days)
4. Save new profiles to SQLite database
5. Display statistics:

```
========== Search Statistics ==========
Total profiles found: 27
New profiles saved: 22
Duplicates skipped: 5
Pages scraped: 3
Errors encountered: 0
Duration: 1m 23s
=======================================
```

### Database Schema

Profiles are stored in SQLite with the following structure:

```sql
CREATE TABLE profiles (
    id TEXT PRIMARY KEY,          -- LinkedIn profile ID (e.g., "john-doe")
    name TEXT,                     -- Full name
    title TEXT,                    -- Current job title
    company TEXT,                  -- Current company
    location TEXT,                 -- Geographic location
    profile_url TEXT,              -- Full LinkedIn URL
    visited_at DATETIME,           -- Last visit timestamp
    created_at DATETIME            -- First scraped timestamp
);
```

## Configuration

### .env File Requirements

Create a `.env` file in the project root (see `.env.example` for template):

```env
# LinkedIn Credentials
LINKEDIN_EMAIL=your_email@example.com
LINKEDIN_PASSWORD=your_password

# Database Configuration
DATABASE_PATH=./data/linkedin_automation.db

# Rate Limits (safe defaults to avoid account restrictions)
MAX_CONNECTIONS_PER_DAY=14      # ~100 connections/week limit
MAX_MESSAGES_PER_DAY=50         # LinkedIn's typical message limit
MAX_SEARCHES_PER_DAY=100        # Conservative search limit
COOLDOWN_SECONDS=30             # Delay between actions

# Activity Scheduling (business hours only)
ACTIVE_HOURS_START=9            # Start at 9 AM
ACTIVE_HOURS_END=17             # End at 5 PM
WEEKDAYS_ONLY=true              # Only run Monday-Friday

# Session Configuration
SESSION_VALIDITY_DAYS=7         # Session expires after 7 days
BROWSER_DATA_DIR=./browser_data # Persistent browser data directory

# Search Configuration
SEARCH_KEYWORDS=software engineer
SEARCH_JOB_TITLE=
SEARCH_COMPANY=
SEARCH_LOCATION=San Francisco Bay Area
```

⚠️ **Security Warning**: Never commit `.env` files to version control. Add it to `.gitignore`.

### Browser Settings

Edit [internal/browser/browser.go](internal/browser/browser.go) to customize:

- **Headless mode**: Change `Headless(false)` to `Headless(true)` for invisible browsing
- **Browser binary**: Set custom Chrome path with `.Bin("/path/to/chrome")`
- **Additional flags**: Add launcher flags as needed

### Stealth Settings

Customize behavior simulation in:
- [internal/stealth/delay.go](internal/stealth/delay.go) - Adjust delay ranges
- [internal/stealth/mouse.go](internal/stealth/mouse.go) - Adjust movement ranges (lines 7-8)
- [internal/stealth/scroll.go](internal/stealth/scroll.go) - Adjust scroll ranges (lines 7-8)
- [internal/stealth/typing.go](internal/stealth/typing.go) - Typing speed parameters

## How Stealth Works

### Mouse Movement Simulation
- Performs 3-5 random mouse movements
- Coordinates range: X(100-800px), Y(100-600px)
- Pause between movements: 300-800ms
- Avoids corners to appear natural

### Scrolling Simulation
- Performs 3-5 random scroll actions
- Scroll distance per action: 200-600 pixels
- Pause between scrolls: 800-1500ms
- Mimics human reading behavior

### Typing Simulation
- Realistic keystroke delays between characters
- Variable speed based on character type
- Occasional typos and corrections (optional)

### Timing & Delays
- Login field detection: 800-1500ms delay
- Form submission: 1000-2000ms delay
- Page load verification: 3000-5000ms delay
- All delays are randomized to appear natural

## Troubleshooting

### Common Issues & Solutions

#### 1. ".env file not found" error

**Problem**: Application exits with error about missing `.env` file

**Solution**:
```bash
# Verify .env exists in project root
ls -la .env

# If missing, create it
echo "LINKEDIN_EMAIL=your_email@example.com" > .env
echo "LINKEDIN_PASSWORD=your_password" >> .env

# Verify contents
cat .env
```

#### 2. "go: command not found"

**Problem**: Go is not installed or not in PATH

**Solution**:
```bash
# Check if Go is installed
go version

# If not, download from https://golang.org/dl/
# Add Go to your PATH environment variable
# Restart PowerShell/Terminal after installation
```

#### 3. "Failed to start Browser" error

**Problem**: Chrome/Chromium not installed

**Solution**:
```bash
# Option 1: Let Rod download Chrome automatically
# The first run will download a Chromium instance

# Option 2: Install Chrome manually
# Download from https://www.google.com/chrome/

# Option 3: Set custom Chrome path in browser.go
.Bin("C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe")
```

#### 4. Windows Defender Blocking Rod

**Problem**: "Operation did not complete successfully because the file contains a virus or potentially unwanted software"

**Solution**:
```powershell
# Add exclusion to Windows Defender
Add-MpPreference -ExclusionPath "$env:TEMP\leakless-*" -Force

# Or manually:
# 1. Open Windows Security
# 2. Go to Virus & threat protection → Virus & threat protection settings
# 3. Click "Manage settings" under "Virus & threat protection settings"
# 4. Scroll down to "Exclusions"
# 5. Click "Add exclusions"
# 6. Select "Folder" and add: C:\Users\YourUsername\AppData\Local\Temp\leakless-*
```

#### 5. "Login Failed: email input not found"

**Problem**: Email input field selector is not found

**Solution**:
```bash
# LinkedIn HTML may have changed
# Edit internal/automation/login.go and update the selector
emailInput, err := page.Element(`input#username`)

# To find the correct selector:
# 1. Open LinkedIn login page in Chrome
# 2. Right-click on email field → Inspect
# 3. Find the correct ID or class name
# 4. Update the selector in login.go
```

#### 6. No Mouse Movements or Scrolling Visible

**Problem**: Application runs but mouse movements and scrolling not visible

**Solution**:
```bash
# Check 1: Browser must be visible (non-headless)
# In browser.go, ensure:
Headless(false)  # NOT true

# Check 2: Verify stealth actions are being called
# Check logs for "Executing random mouse movements..."
# If missing, check main.go has this line:
browser.PerformStealthActions(page)

# Check 3: Increase delay times in stealth files
# Edit internal/stealth/mouse.go and increase:
time.Sleep(time.Duration(1000+r.Intn(1000)) * time.Millisecond)

# Check 4: Verify page is fully loaded
# Wait longer before performing stealth actions
stealth.RandomDelay(3000, 5000)  // 3-5 second wait
```

#### 7. "Login verification timeout"

**Problem**: Application takes too long to verify login

**Solution**:
```bash
# Check internet connection - LinkedIn may be slow
# Increase timeout in internal/automation/login.go:
timeout := time.After(60 * time.Second)  // Increase to 60 seconds

# Check if LinkedIn is accessible
# Try opening LinkedIn in regular Chrome browser
```

#### 8. "Module not found" error

**Problem**: `linkedin-automation/internal/...` module not found

**Solution**:
```bash
# Ensure you're in the correct directory
pwd  # Should end with linkedin-automation

# Reinstall dependencies
go mod tidy
go mod download

# Verify go.mod exists
ls -la go.mod
```

#### 9. "Port 9222 already in use"

**Problem**: Another instance of the application is running

**Solution**:
```bash
# On Windows (PowerShell):
Get-Process | Where-Object {$_.ProcessName -like "*chrome*"} | Stop-Process

# On Linux/Mac:
pkill chrome

# Or just restart your terminal/PowerShell
```

#### 10. Application runs but immediately exits

**Problem**: No error message but program closes immediately

**Solution**:
```bash
# Check that .env file has valid credentials
cat .env

# Check that credentials are correct
# Try logging in manually to verify they work

# Check for silent errors:
go run main.go 2>&1 | tee output.log
# This saves all output to output.log for review
```

### Windows Defender Blocking Rod

If you see: `Operation did not complete successfully because the file contains a virus or potentially unwanted software`

Add an exclusion to Windows Defender:
```powershell
Add-MpPreference -ExclusionPath "$env:TEMP\leakless-*" -Force
```

Or manually:
1. Open Windows Security
2. Go to Virus & threat protection settings
3. Add exclusion for `%TEMP%\leakless-*`

### Browser Not Opening

- Ensure Chrome/Chromium is installed or let Rod download it
- Check that `Headless(false)` is set in browser.go
- Verify no firewall is blocking the browser process

### Login Fails

- Verify credentials in `.env` file are correct
- Check that LinkedIn is accessible from your network
- Ensure the LinkedIn login page HTML selectors haven't changed
- Check console logs for specific error messages

### No Mouse Movements or Scrolling Visible

- Ensure `Headless(false)` is set (browser must be visible)
- Check that `PerformStealthActions(page)` is being called in main.go
- Verify logger output shows stealth action messages
- Increase pause durations in stealth files if movements are too fast to see

## Dependencies

- `github.com/go-rod/rod` - Browser automation library
- `github.com/joho/godotenv` - .env file loading
- `github.com/ysmood/leakless` - Process management (included with rod)

See [go.mod](go.mod) for complete dependencies.

## Development

### Running with Debug Output

Add logging to see detailed execution:
```go
logger.Info("Your debug message")
```

### Adding New Features

1. Create new modules in the `internal/` directory
2. Follow the existing logging patterns
3. Use stealth functions to maintain natural behavior
4. Test thoroughly before deployment

### Extending Stealth Techniques

Add new stealth methods in `internal/stealth/`:
```go
// Example: New stealth function
func RandomHover(page *rod.Page) {
    // Implementation
}
```

Then call from `browser.PerformStealthActions()`.

## Important Notes

⚠️ **LinkedIn's Terms of Service**: Ensure your automation complies with LinkedIn's ToS. Unauthorized automation may result in account suspension.

✅ **Best Practices**:
- Use delays between actions
- Don't perform bulk operations rapidly
- Respect rate limits
- Log all activities for compliance
- Test with a dedicated account first

## Contributing

Feel free to submit issues and enhancement requests!

## License

MIT License - see LICENSE file for details (if added)

## Author

Created for LinkedIn automation tasks.

## Resources

- [Rod Documentation](https://go-rod.github.io/)
- [Go Language Docs](https://golang.org/doc/)
- [Chrome DevTools Protocol](https://chromedevtools.github.io/devtools-protocol/)
- [LinkedIn Developer](https://developer.linkedin.com/)

## Changelog

### v1.0.0 (2025-12-24)
- Initial release
- LinkedIn login automation
- Mouse movement simulation
- Page scrolling simulation
- Fingerprint masking
- Comprehensive logging

