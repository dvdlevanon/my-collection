package opensubtitles

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func NewOpenSubtitles(apiKeys []string) *OpenSubtitiles {
	return &OpenSubtitiles{
		apiKeys: apiKeys,
	}
}

type OpenSubtitiles struct {
	apiKeys         []string
	defaultKeyIndex int
}

func (s *OpenSubtitiles) apiBaseUrl() string {
	return "https://api.opensubtitles.com/api/v1"
}

func (s *OpenSubtitiles) buildUrl(path string) (*url.URL, error) {
	baseUrl := s.apiBaseUrl()
	fullUrl, err := url.JoinPath(baseUrl, path)
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	return url, nil
}

func (s *OpenSubtitiles) buildRequest(method string, url *url.URL, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Test/1.0")
	return req, nil
}

func (s *OpenSubtitiles) fetch(origReq *http.Request) ([]byte, error) {
	if len(s.apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys available")
	}
	client := &http.Client{}
	var lastErr error

	var bodyBytes []byte
	if origReq.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(origReq.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		origReq.Body.Close()
	}

	for i := 0; i < len(s.apiKeys); i++ {
		keyIndex := (s.defaultKeyIndex + i) % len(s.apiKeys)
		apiKey := s.apiKeys[keyIndex]

		req := origReq.Clone(origReq.Context())
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		req.Header.Set("Api-Key", apiKey)

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read response: %w", err)
			}
			return body, nil
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			lastErr = fmt.Errorf("rate limit exceeded for API key %d", i)
			s.defaultKeyIndex++
			continue
		}

		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusNotAcceptable {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			bodyStr := string(body)
			if containsQuotaError(bodyStr) {
				lastErr = fmt.Errorf("quota exceeded for API key %d: %s", i, bodyStr)
				s.defaultKeyIndex++
				continue
			}

			return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, bodyStr)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	if lastErr != nil {
		return nil, fmt.Errorf("all API keys exhausted: %w", lastErr)
	}

	return nil, fmt.Errorf("all API keys failed")
}

func containsQuotaError(body string) bool {
	quotaKeywords := []string{
		"quota",
		"limit exceeded",
		"rate limit",
		"too many requests",
		"downloads remaining",
	}

	bodyLower := strings.ToLower(body)
	for _, keyword := range quotaKeywords {
		if strings.Contains(bodyLower, keyword) {
			return true
		}
	}
	return false
}
