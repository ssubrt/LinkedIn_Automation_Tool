package automation

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"

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
func SearchPeople(page *rod.Page, db *storage.Database, config SearchConfig) (*SearchStats, error) {
	logger.Info("Starting LinkedIn people search")
	logger.Info(fmt.Sprintf("Search parameters: keywords='%s', title='%s', company='%s', location='%s'",
		config.Keywords, config.JobTitle, config.Company, config.Location))

	stats := &SearchStats{
		StartTime: time.Now(),
	}

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
		return stats, fmt.Errorf("failed to build search URL: %w", err)
	}

	logger.Info("Navigating to search URL: " + searchURL)

	// Navigate to search page
	err = page.Navigate(searchURL)
	if err != nil {
		return stats, fmt.Errorf("failed to navigate to search page: %w", err)
	}

	// Wait for results to load
	page.MustWaitLoad()
	time.Sleep(2 * time.Second) // Additional wait for dynamic content

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
				}
			}
		}

		// Try to go to next page
		if pageNum < config.MaxPages {
			hasNext, err := HasNextPage(page)
			if err != nil {
				logger.Warning("Failed to check for next page: " + err.Error())
				break
			}

			if !hasNext {
				logger.Info("No more pages available, search complete")
				break
			}

			// Apply stealth delay before clicking
			stealth.RandomDelay(2000, 4000)

			// Click next page
			err = ClickNextPage(page)
			if err != nil {
				logger.Warning("Failed to navigate to next page: " + err.Error())
				stats.ErrorCount++
				break
			}

			// Wait for new page to load
			page.MustWaitLoad()
			time.Sleep(2 * time.Second)

			// Random scroll to simulate reading
			stealth.RandomScroll(page)
		}
	}

	stats.EndTime = time.Now()
	duration := stats.EndTime.Sub(stats.StartTime)

	logger.Info(fmt.Sprintf("Search completed: %d total found, %d new profiles, %d duplicates, %d pages scraped in %s",
		stats.TotalFound, stats.NewProfiles, stats.Duplicates, stats.PagesScraped, duration))

	return stats, nil
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

	// Wait for search results container
	resultContainers, err := page.Elements(utils.SearchResultItemSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to find search result items: %w", err)
	}

	if len(resultContainers) == 0 {
		// Check if page loaded correctly by looking for alternative selectors
		// This helps detect when LinkedIn changes their HTML structure
		alternativeSelector, _ := page.Element(".search-results-container")
		if alternativeSelector == nil {
			logger.Warning("No results found and page structure unrecognized - LinkedIn may have changed their HTML. Selectors may need updating.")
		}
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
	linkElements, err := container.Elements("a.app-aware-link")
	if err != nil || len(linkElements) == 0 {
		return nil, fmt.Errorf("no profile link found")
	}

	var profileURL string
	for _, link := range linkElements {
		href, err := link.Attribute("href")
		if err != nil || href == nil {
			continue
		}

		// Check if this is a profile link
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

	// Extract job title (primary subtitle)
	subtitleElement, err := container.Element(".entity-result__primary-subtitle")
	if err == nil {
		title, _ := subtitleElement.Text()
		result.Title = strings.TrimSpace(title)
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
func HasNextPage(page *rod.Page) (bool, error) {
	nextButton, err := page.Element(utils.PaginationNextButtonSelector)
	if err != nil {
		// Button not found means no next page
		return false, nil
	}

	// Check if button is disabled
	classes, err := nextButton.Attribute("class")
	if err != nil {
		return false, err
	}

	if classes != nil && strings.Contains(*classes, utils.PaginationDisabledClass) {
		return false, nil
	}

	return true, nil
}

// ClickNextPage clicks the next page button in pagination
func ClickNextPage(page *rod.Page) error {
	nextButton, err := page.Element(utils.PaginationNextButtonSelector)
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

	logger.Info("Navigated to next page")
	return nil
}
