// Package github provides an HTTP client for the GitHub REST API v3.
// It handles pagination, exponential-backoff retries, and automatic
// rate-limit back-off using the X-RateLimit-* response headers.
package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	baseURL        = "https://api.github.com"
	maxPerPage     = 100         // GitHub API max items per page
	maxConcurrent  = 5           // concurrent goroutine limit for batch fetches
	maxRetries     = 3           // retry budget per request (exponential backoff)
	retryBaseDelay = time.Second // base delay doubled on each retry (1s, 2s, 4s)
)

// Client is a thin wrapper over net/http that adds GitHub auth headers,
// automatic pagination, retry logic, and rate-limit handling.
type Client struct {
	httpClient *http.Client
	token      string     // GitHub personal access token or OAuth token
	mu         sync.Mutex // guards concurrent state if needed
}

// NewClient creates a Client with the given OAuth/PAT token.
// Pass an empty string for unauthenticated requests (60 req/hr limit).
func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      token,
	}
}

// doRequest executes an HTTP request with exponential-backoff retry.
// It automatically retries on network errors, 5xx server errors, and
// 403/429 rate-limit responses (sleeping until the reset window passes).
func (c *Client) doRequest(method, url string) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, reqErr := http.NewRequest(method, url, nil)
		if reqErr != nil {
			return nil, fmt.Errorf("creating request: %w", reqErr)
		}

		req.Header.Set("Accept", "application/vnd.github.v3+json")
		if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			if attempt < maxRetries {
				time.Sleep(retryBaseDelay * time.Duration(1<<uint(attempt)))
				continue
			}
			return nil, fmt.Errorf("request failed: %w", err)
		}

		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == 429 {
			c.handleRateLimit(resp)
			resp.Body.Close()
			continue
		}

		if resp.StatusCode >= 500 && attempt < maxRetries {
			resp.Body.Close()
			time.Sleep(retryBaseDelay * time.Duration(1<<uint(attempt)))
			continue
		}

		break
	}

	return resp, err
}

// handleRateLimit sleeps until the GitHub rate-limit window resets.
// It reads X-RateLimit-Remaining and X-RateLimit-Reset headers.
// If the reset time is unreasonably far away (>15 min), falls back to a 60s sleep.
func (c *Client) handleRateLimit(resp *http.Response) {
	remaining := resp.Header.Get("X-RateLimit-Remaining")
	resetStr := resp.Header.Get("X-RateLimit-Reset")

	if remaining == "0" && resetStr != "" {
		resetUnix, err := strconv.ParseInt(resetStr, 10, 64)
		if err == nil {
			resetTime := time.Unix(resetUnix, 0)
			sleepDuration := time.Until(resetTime) + time.Second
			if sleepDuration > 0 && sleepDuration < 15*time.Minute {
				log.Printf("Rate limited. Sleeping until %s (%v)", resetTime.Format(time.RFC3339), sleepDuration)
				time.Sleep(sleepDuration)
				return
			}
		}
	}

	time.Sleep(60 * time.Second)
}

// getJSON fetches a single JSON object from the GitHub API and decodes it into target.
func (c *Client) getJSON(url string, target interface{}) error {
	resp, err := c.doRequest("GET", url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d) for %s: %s", resp.StatusCode, url, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// getPaginated follows GitHub's Link-header pagination to collect all items
// across multiple pages into a single slice of raw JSON messages.
func (c *Client) getPaginated(url string) ([]json.RawMessage, error) {
	var allItems []json.RawMessage
	currentURL := url

	for currentURL != "" {
		resp, err := c.doRequest("GET", currentURL)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}

		var items []json.RawMessage
		if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("decoding response: %w", err)
		}
		resp.Body.Close()

		allItems = append(allItems, items...)
		currentURL = parseNextLink(resp.Header.Get("Link"))
	}

	return allItems, nil
}

// linkNextRegex extracts the URL from a Link header entry with rel="next".
var linkNextRegex = regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)

// parseNextLink extracts the "next" page URL from a GitHub Link header.
// Returns "" if there is no next page.
func parseNextLink(linkHeader string) string {
	if linkHeader == "" {
		return ""
	}

	for _, part := range strings.Split(linkHeader, ",") {
		matches := linkNextRegex.FindStringSubmatch(strings.TrimSpace(part))
		if len(matches) == 2 {
			return matches[1]
		}
	}
	return ""
}
