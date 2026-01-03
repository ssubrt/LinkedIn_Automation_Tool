# LinkedIn Automation - Testing Scenarios

This document outlines the expected behavior and logs for the automation scenarios.

## Scenario 2: Search Test

### Goal
Verify that the automation can navigate to the search page, execute a search, and correctly parse the results from the HTML.

### What You'll See - Browser
1.  **Navigation:** Browser goes to `linkedin.com/search/results/people/...`
2.  **Scrolling:** The page scrolls down gradually to load all results.
3.  **Pagination:** If configured, it clicks "Next" to go to page 2.

### Expected Log Output
```text
[INFO] Starting LinkedIn people search...
[INFO] Navigating to search URL: https://www.linkedin.com/search/results/people/...
[INFO] Scraping page 1/3
[INFO] ✓ Found 10 results with selector: ...
[INFO] Saved new profile: John Doe - Software Engineer
...
[INFO] Search completed: 10 total found, 10 new profiles...
```

---

## Scenario 3: Connection Requests (3 Profiles)

### Goal
Verify that the automation can visit a profile, find the "Connect" button (even if hidden), and send a personalized note.

### What You'll See - Browser (for EACH profile)

**Step 1: Profile Page (3-5 seconds)**
*   Navigates to `https://www.linkedin.com/in/[username]`
*   Shows full profile with photo, headline, experience
*   Automated scroll down the page (simulating reading)

**Step 2: Connect Button Detection (1-2 seconds)**
*   Searches for blue "Connect" button
*   OR searches "More" dropdown if 3rd-degree connection

**Step 3: Modal Opens (2 seconds)**
*   "Connect" button clicked
*   Modal pops up with "Add a note" option
*   **YOU'LL SEE:** Note textarea appears

**Step 4: Typing Note (2-4 seconds)**
*   Character by character typing (80-150ms per char)
*   Message appears: "Hi John, I'd love to connect..."
*   **NOT** instant paste - watch each character appear

**Step 5: Send (1 second)**
*   "Send" button clicked
*   Modal closes
*   Success indicator appears

**Cooldown (30 seconds)**
*   Browser stays on current page
*   No activity for 30 seconds
*   Then moves to next profile

### Expected Log (After 3 profiles)
```text
[INFO] Sending connection request to: John Doe...
[INFO] Clicking Connect button...
[INFO] Adding personalized note...
[INFO] Typing note (45 characters)...
[INFO] Connection request sent successfully to John Doe
...
[INFO] ✓ Sent 3/3 connection requests successfully
[INFO] Failed: 0
[INFO] Total execution time: 4 minutes 45 seconds
```



------------------------------>>>>>>>>>>>>>>>>>>>>>>     With Logs.  <<<<<<<<<<<<<<<<<----------------------------------------- 


🧪 Complete Testing Flow Guide - What You'll See On Screen
Step 1: Verify Your .env File
Step 2: Build the Project
 - go build -o linkedin-automation

🎬 Run Scenario 1: Login + Stealth Test (2-3 minutes)
What to Run:

./linkedin-automation

What You'll See - Terminal Logs:
[2026-01-01 18:05:32] INFO: Starting LinkedIn Automation
[2026-01-01 18:05:32] INFO: Loading environment variables...
[2026-01-01 18:05:32] INFO: Checking activity schedule...
[2026-01-01 18:05:32] INFO: Current time: Wednesday, 2026-01-01 18:05:32
[2026-01-01 18:05:32] INFO: Activity hours: 9:00 - 17:00 (weekdays only)
[2026-01-01 18:05:32] INFO: ✓ Within active hours
[2026-01-01 18:05:32] INFO: Initializing database...
[2026-01-01 18:05:32] INFO: Database initialized: ./data/linkedin_automation.db
[2026-01-01 18:05:32] INFO: Initializing rate limiter...
[2026-01-01 18:05:32] INFO: Rate limits: Connections=14/day, Messages=50/day, Searches=100/day
[2026-01-01 18:05:32] INFO: Checking for existing session...
[2026-01-01 18:05:32] INFO: No valid session found, will perform login
[2026-01-01 18:05:33] INFO: Initializing browser...
[2026-01-01 18:05:35] INFO: Browser launched successfully
[2026-01-01 18:05:35] INFO: Applying fingerprint masking...
[2026-01-01 18:05:35] INFO: ✓ WebDriver property masked
[2026-01-01 18:05:35] INFO: ✓ Chrome object masked
[2026-01-01 18:05:35] INFO: ✓ Plugins randomized
[2026-01-01 18:05:35] INFO: ✓ Languages set to [en-US, en]
[2026-01-01 18:05:35] INFO: ✓ Permissions API overridden
[2026-01-01 18:05:36] INFO: Attempting to log into LinkedIn...
[2026-01-01 18:05:36] INFO: Navigating to LinkedIn login page...



Stage 1: Login Page (5-10 seconds)


Chrome window opens (NOT headless, you can see it)
Navigates to https://www.linkedin.com/login
WATCH THE TYPING:
Email field fills character by character (NOT instant paste)
Each character has 80-150ms delay (human-like)
Password field fills character by character
Each character has 80-150ms delay
Expected Log:
[2026-01-01 18:05:38] INFO: Typing email into username field...
[2026-01-01 18:05:40] INFO: Typing password into password field...
[2026-01-01 18:05:42] INFO: Clicking Sign in button...

Stage 2: Post-Login (2-3 seconds)

[2026-01-01 18:05:45] INFO: ✓ Login successful - redirected to LinkedIn feed
[2026-01-01 18:05:45] INFO: Saving session state...
[2026-01-01 18:05:45] INFO: Session saved to: ./browser_data/state.json



Stage 3: Stealth Actions (10-15 seconds)

MOUSE MOVEMENTS (watch cursor):

[2026-01-01 18:05:46] INFO: Performing stealth actions - simulating human-like behavior
[2026-01-01 18:05:46] INFO: Executing random mouse movements...


YOU'LL SEE: Cursor moves smoothly across the page
Path: Bézier curves (natural acceleration/deceleration)
Count: 3-5 different positions
Pause: 300-800ms between each movement
NOT instant jumps - smooth gliding motion

PAGE SCROLLING (watch page):

[2026-01-01 18:05:48] INFO: Executing random page scrolling...


YOU'LL SEE: Page scrolls down in chunks
Distance: 200-600 pixels per scroll
Pause: 800-1500ms between scrolls (simulating reading)
Count: 3-5 scrolls total
Behavior: Natural, not robotic


HOVER ACTIONS (watch cursor + elements):

[2026-01-01 18:05:50] INFO: Hovering over interactive elements...


YOU'LL SEE: Cursor moves to buttons/links and hovers
Count: 2-3 random elements
Pause: 500-1000ms hover time
Elements: Notification bell, messaging icon, search bar, etc.



Final Output:

[2026-01-01 18:05:52] INFO: Stealth actions completed
[2026-01-01 18:05:52] INFO: ✓ LinkedIn automation completed successfully
[2026-01-01 18:05:52] INFO: Total execution time: 20.5 seconds


Browser stays open for 10 seconds, then closes automatically.


🎬 Run Scenario 2: Search Test (5 minutes)


Update .env:
SEARCH_KEYWORDS=software engineer
SEARCH_LOCATION=San Francisco Bay Area
ENABLE_CONNECTIONS=false  # Still disabled
ENABLE_MESSAGING=false    # Still disabled


Run:
./linkedin-automation

Expected Logs - Search Phase:
[2026-01-01 18:10:15] INFO: Step 8: LinkedIn People Search
[2026-01-01 18:10:15] INFO: Search keywords: software engineer
[2026-01-01 18:10:15] INFO: Search location: San Francisco Bay Area
[2026-01-01 18:10:15] INFO: Building search URL...
[2026-01-01 18:10:15] INFO: Search URL: /search/results/people/?keywords=software%20engineer&geoUrn=["urn:li:fs_geo:90000084"]
[2026-01-01 18:10:15] INFO: Navigating to search page...
[2026-01-01 18:10:18] INFO: Parsing search results...
[2026-01-01 18:10:18] INFO: Found 10 profiles on page 1
[2026-01-01 18:10:18] INFO: Profile 1: John Doe - Senior Software Engineer at Google
[2026-01-01 18:10:18] INFO: Profile 2: Jane Smith - Software Engineer at Meta
[2026-01-01 18:10:18] INFO: Profile 3: Mike Johnson - Staff Engineer at Apple
...
[2026-01-01 18:10:20] INFO: Saved 10 new profiles to database
[2026-01-01 18:10:20] INFO: Checking for next page...
[2026-01-01 18:10:20] INFO: Next page button found, navigating...
[2026-01-01 18:10:23] INFO: Found 10 profiles on page 2
...

What You'll See - Browser:
LinkedIn Search Results Page

URL: https://www.linkedin.com/search/results/people/?keywords=...
List of profiles with photos, names, titles
Scrolling through results (automated)
Pagination

After page 1, clicks "Next" button
Loads page 2, 3, etc.
Continues until no more pages or limit reached

Expected Database:

sqlite3 ./data/linkedin_automation.db "SELECT COUNT(*) FROM profiles;"
# Output: 20 (or however many profiles found)

sqlite3 ./data/linkedin_automation.db "SELECT name, headline, company FROM profiles LIMIT 5;"
# Output:
# John Doe|Senior Software Engineer|Google
# Jane Smith|Software Engineer|Meta
# ...



🎬 Run Scenario 3: Connection Requests (10 minutes)

Update .env:
ENABLE_CONNECTIONS=true
MAX_CONNECTIONS_PER_RUN=3  # Start with just 3 for testing
YOUR_NAME=Your Full Name
YOUR_TITLE=Your Job Title
YOUR_COMPANY=Your Company
CONNECTION_TEMPLATE=conn_generic


Run:
./linkedin-automation

Expected Logs - Connection Phase:
[2026-01-01 18:15:30] INFO: Step 9: Send Connection Requests
[2026-01-01 18:15:30] INFO: Connection requests enabled: true
[2026-01-01 18:15:30] INFO: Max connections per run: 3
[2026-01-01 18:15:30] INFO: Checking daily rate limit...
[2026-01-01 18:15:30] INFO: Today's connections: 0/14
[2026-01-01 18:15:30] INFO: Retrieving profiles not yet contacted...
[2026-01-01 18:15:30] INFO: Found 20 profiles, will send 3 connection requests
[2026-01-01 18:15:30] INFO: 
[2026-01-01 18:15:30] INFO: Processing profile 1/3: John Doe
[2026-01-01 18:15:30] INFO: Navigating to profile: https://www.linkedin.com/in/johndoe
[2026-01-01 18:15:33] INFO: Scrolling to simulate reading profile...
[2026-01-01 18:15:35] INFO: Checking if already connected...
[2026-01-01 18:15:35] INFO: Not connected, proceeding with request
[2026-01-01 18:15:35] INFO: Looking for Connect button...
[2026-01-01 18:15:35] INFO: Connect button found
[2026-01-01 18:15:35] INFO: Clicking Connect button...
[2026-01-01 18:15:36] INFO: Connection note modal opened
[2026-01-01 18:15:36] INFO: Rendering connection template: conn_generic
[2026-01-01 18:15:36] INFO: Typing connection note (49 characters)...
[2026-01-01 18:15:38] INFO: Clicking Send button...
[2026-01-01 18:15:39] INFO: ✓ Connection request sent successfully
[2026-01-01 18:15:39] INFO: Saved to database: status=pending
[2026-01-01 18:15:39] INFO: Cooldown: waiting 30 seconds before next action...



What You'll See - Browser (for EACH profile):
Step 1: Profile Page (3-5 seconds)

Navigates to https://www.linkedin.com/in/[username]
Shows full profile with photo, headline, experience
Automated scroll down the page (simulating reading)
Step 2: Connect Button Detection (1-2 seconds)

Searches for blue "Connect" button
OR searches "More" dropdown if 3rd-degree connection
Step 3: Modal Opens (2 seconds)

"Connect" button clicked
Modal pops up with "Add a note" option
YOU'LL SEE: Note textarea appears
Step 4: Typing Note (2-4 seconds)

Character by character typing (80-150ms per char)
Message appears: "Hi John, I'd love to connect..."
NOT instant paste - watch each character appear
Step 5: Send (1 second)

"Send" button clicked
Modal closes
Success indicator appears
Cooldown (30 seconds)

Browser stays on current page
No activity for 30 seconds
Then moves to next profile
Expected Log (After 3 profiles):
[2026-01-01 18:20:15] INFO: ✓ Sent 3/3 connection requests successfully
[2026-01-01 18:20:15] INFO: Failed: 0
[2026-01-01 18:20:15] INFO: Total execution time: 4 minutes 45 seconds


🎬 Run Scenario 4: Connection Status Check (5 minutes)


Wait 1-2 days for LinkedIn to accept/reject connections

Update .env:
CHECK_CONNECTION_STATUS=true
ENABLE_CONNECTIONS=false  # Disable new connections
ENABLE_MESSAGING=false    # Disable messaging for now


Run:
./linkedin-automation

Expected Logs:
[2026-01-03 10:15:30] INFO: Step 9.5: Check Connection Status
[2026-01-03 10:15:30] INFO: Checking pending connections for status updates...
[2026-01-03 10:15:30] INFO: Found 3 pending connections to check
[2026-01-03 10:15:30] INFO: Navigating to My Network page...
[2026-01-03 10:15:33] INFO: Checking profile 1/3: John Doe
[2026-01-03 10:15:33] INFO: Navigating to: https://www.linkedin.com/in/johndoe
[2026-01-03 10:15:36] INFO: Looking for 'Connected' indicator...
[2026-01-03 10:15:36] INFO: ✓ Connection ACCEPTED - updating status
[2026-01-03 10:15:36] INFO: Updated database: status=accepted
[2026-01-03 10:15:36] INFO: 
[2026-01-03 10:15:36] INFO: Checking profile 2/3: Jane Smith
[2026-01-03 10:15:36] INFO: Navigating to: https://www.linkedin.com/in/janesmith
[2026-01-03 10:15:39] INFO: Looking for 'Connected' indicator...
[2026-01-03 10:15:39] INFO: Still pending, no update needed
...
[2026-01-03 10:20:15] INFO: ✓ Status check complete: 2 newly accepted, 1 still pending

What You'll See - Browser:
Visits each pending connection's profile
Looks for "Connected" badge/text
Updates database if found




🎬 Run Scenario 5: Messaging Accepted Connections (10 minutes)
Update .env:
ENABLE_MESSAGING=true
MAX_MESSAGES_PER_RUN=2  # Start with 2
MESSAGE_TEMPLATE=msg_introduction
MESSAGE_CUSTOM_REASON=I have insights about software engineering



Run:
./linkedin-automation


Expected Logs:
[2026-01-03 10:30:15] INFO: Step 10: Send Messages to Accepted Connections
[2026-01-03 10:30:15] INFO: Messaging enabled: true
[2026-01-03 10:30:15] INFO: Max messages per run: 2
[2026-01-03 10:30:15] INFO: Retrieving accepted connections not yet messaged...
[2026-01-03 10:30:15] INFO: Found 2 accepted connections
[2026-01-03 10:30:15] INFO: 
[2026-01-03 10:30:15] INFO: Processing connection 1/2: John Doe
[2026-01-03 10:30:15] INFO: Navigating to profile: https://www.linkedin.com/in/johndoe
[2026-01-03 10:30:18] INFO: Looking for Message button...
[2026-01-03 10:30:18] INFO: Message button found, clicking...
[2026-01-03 10:30:19] INFO: Message modal opened
[2026-01-03 10:30:19] INFO: Rendering message template: msg_introduction
[2026-01-03 10:30:19] INFO: Subject: Quick introduction
[2026-01-03 10:30:19] INFO: Body length: 156 characters
[2026-01-03 10:30:19] INFO: Typing message...
[2026-01-03 10:30:23] INFO: Clicking Send button...
[2026-01-03 10:30:24] INFO: ✓ Message sent successfully
[2026-01-03 10:30:24] INFO: Saved to database: messages table
[2026-01-03 10:30:24] INFO: Cooldown: waiting 30 seconds...



What You'll See - Browser:
Profile page → "Message" button clicked
Message modal opens (like LinkedIn chat)
Typing animation (character by character)
Send button clicked
Modal closes → success


✅ Success Checklist
After running all scenarios, you should see:

 Browser opens and closes automatically
 Login completes (with or without 2FA)
 Mouse movements are visible and smooth (Bézier curves)
 Page scrolling happens naturally (not instant jumps)
 Character-by-character typing (NOT paste)
 Database file created: ./data/linkedin_automation.db
 Session file created: ./browser_data/state.json
 Profiles saved to database (check with sqlite3 command)
 Connection requests sent (visible in LinkedIn "Sent" tab)
 Messages sent (visible in LinkedIn inbox)
 All logs show INFO (not ERROR)
 Program exits with code 0 (success)


