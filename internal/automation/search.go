package automation

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"

	"linkedin-automation/internal/logger"
	"linkedin-automation/internal/stealth"
	"linkedin-automation/internal/storage"
	"linkedin-automation/pkg/utils"
)

// SearchConfig holds configuration for LinkedIn people search
type SearchConfig struct {
	Keywords string // General search keywords
	JobTitle string // Filter by job title
	Company  string // Filter by company name
	Location string // Location name (e.g., "San Francisco Bay Area")

	// Pagination settings
	MaxPages int // Maximum number of pages to scrape (0 = all available)

	// Duplicate handling
	SkipDuplicates bool // Skip profiles visited in last 30 days
	DuplicateDays  int  // Days to consider as duplicate (default: 30)
}

// SearchResult represents a parsed profile from search results
type SearchResult struct {
	ProfileID  string    // Extracted from URL
	Name       string    // Full name
	Title      string    // Current job title
	Company    string    // Current company
	Location   string    // Geographic location
	ProfileURL string    // Full LinkedIn profile URL
	Degree     string    // Connection degree (1st, 2nd, 3rd)
	ScrapedAt  time.Time // When this result was found
}

// SearchStats tracks statistics for a search session
type SearchStats struct {
	TotalFound   int
	NewProfiles  int
	Duplicates   int
	PagesScraped int
	ErrorCount   int
	StartTime    time.Time
	EndTime      time.Time
}

// SearchPeople performs a LinkedIn people search with the given configuration
func SearchPeople(page *rod.Page, db *storage.Database, config SearchConfig) ([]SearchResult, *SearchStats, error) {
	logger.Info("Starting LinkedIn people search")
	logger.Info(fmt.Sprintf("Search parameters: keywords='%s', title='%s', company='%s', location='%s'",
		config.Keywords, config.JobTitle, config.Company, config.Location))

	stats := &SearchStats{
		StartTime: time.Now(),
	}
	var allResults []SearchResult

	// Set default values
	if config.MaxPages == 0 {
		config.MaxPages = utils.MaxPaginationPages
	}
	if config.DuplicateDays == 0 {
		config.DuplicateDays = 30
	}
	if !config.SkipDuplicates {
		config.SkipDuplicates = true // Default to skip duplicates
	}

	// Build search URL
	searchURL, err := buildSearchURL(config)
	if err != nil {
		return nil, stats, fmt.Errorf("failed to build search URL: %w", err)
	}

	logger.Info("Navigating to search URL: " + searchURL)

	// Navigate to search page
	err = page.Navigate(searchURL)
	if err != nil {
		return nil, stats, fmt.Errorf("failed to navigate to search page: %w", err)
	}

	// Wait for results to load
	page.MustWaitLoad()
	time.Sleep(2 * time.Second) // Additional wait for dynamic content

	// Check for LinkedIn checkpoint/verification page
	currentURL := page.MustInfo().URL
	if utils.IsLinkedInCheckpoint(currentURL) {
		logger.Error("❌ LinkedIn checkpoint/verification detected at: " + currentURL)
		return nil, stats, fmt.Errorf("linkedin checkpoint detected, manual verification required")
	}

	// Apply stealth actions
	stealth.RandomDelay(500, 1000)

	// Scrape pages
	for pageNum := 1; pageNum <= config.MaxPages; pageNum++ {
		logger.Info(fmt.Sprintf("Scraping page %d/%d", pageNum, config.MaxPages))

		// Parse current page results
		results, err := ParseSearchResults(page)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to parse page %d: %s", pageNum, err.Error()))
			stats.ErrorCount++
			break
		}

		if len(results) == 0 {
			logger.Info("No results found on this page, stopping pagination")
			break
		}

		logger.Info(fmt.Sprintf("Found %d profiles on page %d", len(results), pageNum))
		stats.TotalFound += len(results)
		stats.PagesScraped++

		// Process each result
		for _, result := range results {
			// Check for duplicates if enabled
			if config.SkipDuplicates && db != nil {
				isDupe, err := db.IsDuplicateProfile(result.ProfileID, config.DuplicateDays)
				if err != nil {
					logger.Warning(fmt.Sprintf("Failed to check duplicate for %s: %s", result.ProfileID, err.Error()))
				} else if isDupe {
					logger.Info(fmt.Sprintf("Skipping duplicate profile: %s", result.Name))
					stats.Duplicates++
					continue
				}
			}

			// Save new profile to database
			if db != nil {
				profile := storage.Profile{
					ID:         result.ProfileID,
					Name:       result.Name,
					Title:      result.Title,
					Company:    result.Company,
					Location:   result.Location,
					ProfileURL: result.ProfileURL,
					VisitedAt:  result.ScrapedAt,
					CreatedAt:  result.ScrapedAt,
				}

				err := db.SaveProfile(profile)
				if err != nil {
					logger.Warning(fmt.Sprintf("Failed to save profile %s: %s", result.ProfileID, err.Error()))
					stats.ErrorCount++
				} else {
					logger.Info(fmt.Sprintf("Saved new profile: %s - %s", result.Name, result.Title))
					stats.NewProfiles++
					allResults = append(allResults, result)
				}
			}
		}

		// PAGINATION DISABLED FOR NOW - Just scrape first page to avoid getting stuck
		// LinkedIn has massive pagination that can cause the automation to hang
		logger.Info("Pagination disabled - only scraping first page")
		break
	}

	stats.EndTime = time.Now()
	duration := stats.EndTime.Sub(stats.StartTime)

	logger.Info(fmt.Sprintf("Search completed: %d total found, %d new profiles, %d duplicates, %d pages scraped in %s",
		stats.TotalFound, stats.NewProfiles, stats.Duplicates, stats.PagesScraped, duration))

	return allResults, stats, nil
}

// buildSearchURL constructs a LinkedIn people search URL with query parameters
func buildSearchURL(config SearchConfig) (string, error) {
	baseURL := utils.LinkedInSearchURL
	params := url.Values{}

	// Add keywords (main search query)
	if config.Keywords != "" {
		params.Add("keywords", config.Keywords)
	}

	// Add title filter
	if config.JobTitle != "" {
		params.Add("title", config.JobTitle)
	}

	// Add company filter
	if config.Company != "" {
		params.Add("company", config.Company)
	}

	// Add location filter (convert name to URN)
	if config.Location != "" {
		locationURN, found := utils.LinkedInLocations[config.Location]
		if found {
			params.Add("geoUrn", fmt.Sprintf("[\"urn:li:fs_geo:%s\"]", locationURN))
		} else {
			logger.Warning(fmt.Sprintf("Location '%s' not found in location map, skipping", config.Location))
		}
	}

	// Build final URL
	if len(params) == 0 {
		return "", fmt.Errorf("at least one search parameter is required")
	}

	fullURL := baseURL + "?" + params.Encode()
	return fullURL, nil
}

// ParseSearchResults extracts profile information from the current search results page
func ParseSearchResults(page *rod.Page) ([]SearchResult, error) {
	var results []SearchResult

	// Wait for search results container with timeout
	logger.Info("Waiting for search results to load...")

	// Debug: Check what's actually on the page
	logger.Info("Debugging: Checking page structure...")
	pageHTML, _ := page.HTML()
	logger.Info(fmt.Sprintf("Page HTML length: %d characters", len(pageHTML)))

	// Try to find what containers exist
	debugScript := `
		const containers = [
			'li.reusable-search__result-container',
			'.entity-result',
			'div.search-result',
			'.reusable-search-simple-insight__container',
			'[data-view-name="search-entity-result"]',
			'.scaffold-layout__list-container li',
			'.search-results-container li'
		];
		const found = [];
		containers.forEach(sel => {
			const count = document.querySelectorAll(sel).length;
			if (count > 0) found.push({selector: sel, count: count});
		});
		return found;
	`
	foundSelectors, _ := page.Eval(debugScript)
	if foundSelectors != nil {
		logger.Info(fmt.Sprintf("Found selectors on page: %v", foundSelectors.Value))
	}

	// Try multiple selectors since LinkedIn frequently changes their HTML structure
	var resultContainers rod.Elements
	var err error

	// Attempt 1: Modern LinkedIn structure (2024-2026)
	resultContainers, err = page.Timeout(5 * time.Second).Elements("li.reusable-search__result-container")
	if err == nil && len(resultContainers) > 0 {
		logger.Info(fmt.Sprintf("✓ Found %d results with selector: li.reusable-search__result-container", len(resultContainers)))
		goto parseResults
	}

	// Attempt 2: Try with data attribute
	logger.Info("Trying data attribute selector...")
	resultContainers, err = page.Timeout(5 * time.Second).Elements("[data-view-name=\"search-entity-result\"]")
	if err == nil && len(resultContainers) > 0 {
		logger.Info(fmt.Sprintf("✓ Found %d results with data-view-name selector", len(resultContainers)))
		goto parseResults
	}

	// Attempt 3: Scaffold layout list items
	logger.Info("Trying scaffold layout selector...")
	resultContainers, err = page.Timeout(5 * time.Second).Elements(".scaffold-layout__list-container li")
	if err == nil && len(resultContainers) > 0 {
		logger.Info(fmt.Sprintf("✓ Found %d results with scaffold-layout selector", len(resultContainers)))
		goto parseResults
	}

	// Attempt 4: Older structure
	logger.Info("Trying legacy selector...")
	resultContainers, err = page.Timeout(5 * time.Second).Elements(".entity-result")
	if err == nil && len(resultContainers) > 0 {
		logger.Info(fmt.Sprintf("✓ Found %d results with .entity-result selector", len(resultContainers)))
		goto parseResults
	}

	// Attempt 5: Generic List Item in Main (Fallback)
	logger.Info("Trying generic list selector...")
	resultContainers, err = page.Timeout(5 * time.Second).Elements("main ul li")
	if err == nil && len(resultContainers) > 0 {
		// Filter out small items (like dividers or loading spinners)
		var validContainers rod.Elements
		for _, c := range resultContainers {
			// Check if it has a link to a profile
			links, _ := c.Elements("a[href*='/in/']")
			if len(links) > 0 {
				validContainers = append(validContainers, c)
			}
		}
		if len(validContainers) > 0 {
			resultContainers = validContainers
			logger.Info(fmt.Sprintf("✓ Found %d results with generic main ul li selector", len(resultContainers)))
			goto parseResults
		}
	}

	// Attempt 6: JS-based Link Discovery (Nuclear Option)
	logger.Info("Trying JS-based link discovery...")
	{
		// This script finds all profile links and walks up to find their container (li or div)
		linkDiscoveryScript := `() => {
			const links = Array.from(document.querySelectorAll("main a[href*='/in/']"));
			const uniqueContainers = new Set();
			
			links.forEach(link => {
				// Ignore links that are too short or look like artifacts
				if (link.getAttribute('href').length < 25) return;
				
				// Walk up to find a likely container
				// We look for list items or divs that look like cards
				let container = link.closest('li');
				if (!container) {
					container = link.closest('div.entity-result');
				}
				if (!container) {
					container = link.closest('div[data-view-name="search-entity-result"]');
				}
				// Fallback: just use the parent div of the link's wrapper
				if (!container) {
					container = link.parentElement?.parentElement?.parentElement;
				}
				
				if (container) {
					uniqueContainers.add(container);
				}
			});
			return Array.from(uniqueContainers);
		}`

		resultContainers, err = page.ElementsByJS(rod.Eval(linkDiscoveryScript))
		if err == nil && len(resultContainers) > 0 {
			logger.Info(fmt.Sprintf("✓ Found %d results via JS link discovery", len(resultContainers)))
			goto parseResults
		}
	}

	logger.Warning("Could not find search results with any known selector")

	// Debug: Dump HTML structure to help identify the issue
	{
		logger.Info("DEBUG: Dumping page structure analysis...")
		structureAnalysis, _ := page.Eval(`() => {
			const main = document.querySelector('main');
			if (!main) return "No main tag found";
			
			const lists = main.querySelectorAll('ul');
			const analysis = [];
			lists.forEach((ul, i) => {
				analysis.push("List " + i + " classes: " + ul.className + ", items: " + ul.querySelectorAll('li').length);
			});
			return analysis.join('\n');
		}`)
		logger.Info(fmt.Sprintf("Page Structure:\n%v", structureAnalysis))
	}

parseResults:
	// Check if we have any results
	if len(resultContainers) == 0 {
		logger.Warning("No results found - page structure may have changed")
		return results, nil // Empty results, not an error
	}

	logger.Info(fmt.Sprintf("Parsing %d result containers", len(resultContainers)))

	for i, container := range resultContainers {
		result, err := parseProfileFromContainer(container)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to parse result %d: %s", i+1, err.Error()))
			continue
		}

		if result != nil {
			results = append(results, *result)
		}
	}

	return results, nil
}

// parseProfileFromContainer extracts profile data from a single result container
func parseProfileFromContainer(container *rod.Element) (*SearchResult, error) {
	result := &SearchResult{
		ScrapedAt: time.Now(),
	}

	// Extract profile URL and ID
	// Try generic selector for any profile link
	linkElements, err := container.Elements("a[href*='/in/']")
	if err != nil || len(linkElements) == 0 {
		// Fallback: try finding any link and checking href
		linkElements, err = container.Elements("a")
		if err != nil || len(linkElements) == 0 {
			return nil, fmt.Errorf("no profile link found")
		}
	}

	var profileURL string
	for _, link := range linkElements {
		href, err := link.Attribute("href")
		if err != nil || href == nil {
			continue
		}

		// Check if this is a profile link (and not a shared post or other noise)
		// Profile links usually look like /in/username
		if strings.Contains(*href, "/in/") {
			profileURL = *href
			break
		}
	}

	if profileURL == "" {
		return nil, fmt.Errorf("no valid profile URL found")
	}

	// Clean URL (remove query params)
	if idx := strings.Index(profileURL, "?"); idx != -1 {
		profileURL = profileURL[:idx]
	}

	result.ProfileURL = profileURL

	// Extract profile ID from URL (e.g., /in/john-doe/ -> john-doe)
	if strings.Contains(profileURL, "/in/") {
		parts := strings.Split(profileURL, "/in/")
		if len(parts) >= 2 {
			profileID := strings.TrimSuffix(parts[1], "/")
			result.ProfileID = profileID
		}
	}

	if result.ProfileID == "" {
		return nil, fmt.Errorf("could not extract profile ID from URL: %s", profileURL)
	}

	// Extract name (from title link)
	titleElement, err := container.Element(".entity-result__title-text a span[aria-hidden='true']")
	if err == nil {
		name, _ := titleElement.Text()
		result.Name = strings.TrimSpace(name)
	}

	// Fallback for name if first method fails
	if result.Name == "" {
		titleElement, err := container.Element(".entity-result__title-text")
		if err == nil {
			name, _ := titleElement.Text()
			result.Name = strings.TrimSpace(name)
		}
	}

	// Fallback 3: Try to get name from the profile link itself
	if result.Name == "" {
		for _, link := range linkElements {
			text, err := link.Text()
			if err == nil && len(text) > 3 && !strings.Contains(text, "LinkedIn Member") {
				result.Name = strings.TrimSpace(text)
				break
			}
		}
	}

	// Extract job title (primary subtitle)
	subtitleElement, err := container.Element(".entity-result__primary-subtitle")
	if err == nil {
		title, _ := subtitleElement.Text()
		result.Title = strings.TrimSpace(title)
	} else {
		// Fallback: Try to find any text that looks like a title (often in a div below the name)
		// This is a heuristic: look for text that is not the name and not the location
		// For now, we'll leave it empty if specific selector fails to avoid garbage data
	}

	// Extract company/location (secondary subtitle)
	secondaryElement, err := container.Element(".entity-result__secondary-subtitle")
	if err == nil {
		secondary, _ := secondaryElement.Text()
		secondary = strings.TrimSpace(secondary)

		// Often format is "Company | Location" or just "Location"
		if strings.Contains(secondary, " | ") {
			parts := strings.Split(secondary, " | ")
			if len(parts) >= 1 {
				result.Company = strings.TrimSpace(parts[0])
			}
			if len(parts) >= 2 {
				result.Location = strings.TrimSpace(parts[1])
			}
		} else {
			result.Location = secondary
		}
	} else {
		// Fallback: Try to find location in any secondary text
		// Often location is in a span with class containing 'location' or 'secondary'
	}

	// Extract connection degree (e.g., "1st", "2nd", "3rd")
	degreeElement, err := container.Element(".entity-result__badge-text .t-black--light")
	if err == nil {
		degree, _ := degreeElement.Text()
		result.Degree = strings.TrimSpace(degree)
	}

	return result, nil
}

// HasNextPage checks if there's a next page button available
/*
func HasNextPage(page *rod.Page) (bool, error) {
	logger.Info("Checking for next page button...")
	nextButton, err := page.Timeout(5 * time.Second).Element(utils.PaginationNextButtonSelector)
	if err != nil {
		// Button not found means no next page
		logger.Info("Next page button not found - no more pages")
		return false, nil
	}

	// Check if button is disabled
	classes, err := nextButton.Attribute("class")
	if err != nil {
		logger.Warning("Failed to get button classes: " + err.Error())
		return false, err
	}

	if classes != nil && strings.Contains(*classes, utils.PaginationDisabledClass) {
		logger.Info("Next page button is disabled - no more pages")
		return false, nil
	}

	logger.Info("Next page button found and enabled")
	return true, nil
}
*/

// ClickNextPage clicks the next page button in pagination
/*
func ClickNextPage(page *rod.Page) error {
	logger.Info("Clicking next page button...")
	nextButton, err := page.Timeout(5 * time.Second).Element(utils.PaginationNextButtonSelector)
	if err != nil {
		return fmt.Errorf("next page button not found: %w", err)
	}

	// Check if button is disabled before clicking
	classes, err := nextButton.Attribute("class")
	if err != nil {
		return err
	}

	if classes != nil && strings.Contains(*classes, utils.PaginationDisabledClass) {
		return fmt.Errorf("next page button is disabled")
	}

	// Scroll button into view
	err = nextButton.ScrollIntoView()
	if err != nil {
		return fmt.Errorf("failed to scroll button into view: %w", err)
	}

	// Add slight delay
	time.Sleep(500 * time.Millisecond)

	// Click the button
	err = nextButton.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("failed to click next page button: %w", err)
	}

	logger.Info("Successfully clicked next page button")
	return nil
}
*/
