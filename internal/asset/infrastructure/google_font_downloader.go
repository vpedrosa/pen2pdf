package infrastructure

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var ttfURLRegex = regexp.MustCompile(`url\(([^)]+\.ttf)\)`)

// GoogleFontDownloader downloads static font files from Google Fonts.
type GoogleFontDownloader struct {
	client *http.Client
}

func NewGoogleFontDownloader() *GoogleFontDownloader {
	return &GoogleFontDownloader{client: &http.Client{}}
}

// Download fetches a static TTF file for the given font family, weight, and style
// from Google Fonts and saves it to destDir.
// Returns the path to the saved file, or an error.
func (d *GoogleFontDownloader) Download(family, weight, style, destDir string) (string, error) {
	if weight == "" || weight == "normal" {
		weight = "400"
	}

	// Build Google Fonts CSS2 API URL
	cssURL := buildCSSURL(family, weight, style)

	// Fetch CSS to extract TTF URL (fall back to non-italic if italic not available)
	ttfURL, err := d.fetchTTFURL(cssURL)
	if err != nil && style == "italic" {
		cssURL = buildCSSURL(family, weight, "")
		ttfURL, err = d.fetchTTFURL(cssURL)
	}
	if err != nil {
		return "", fmt.Errorf("fetch font CSS for %q weight %s: %w", family, weight, err)
	}

	// Download the TTF file
	data, err := d.fetchBytes(ttfURL)
	if err != nil {
		return "", fmt.Errorf("download font %q: %w", family, err)
	}

	// Save to destDir
	filename := buildFilename(family, weight, style)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("create font directory: %w", err)
	}

	destPath := filepath.Join(destDir, filename)
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return "", fmt.Errorf("write font file: %w", err)
	}

	return destPath, nil
}

func buildCSSURL(family, weight, style string) string {
	// Google Fonts CSS2 API: family=Inter:ital,wght@0,700 or family=Inter:wght@700
	familyParam := url.QueryEscape(family)

	if style == "italic" {
		return fmt.Sprintf("https://fonts.googleapis.com/css2?family=%s:ital,wght@1,%s", familyParam, weight)
	}
	return fmt.Sprintf("https://fonts.googleapis.com/css2?family=%s:wght@%s", familyParam, weight)
}

func buildFilename(family, weight, style string) string {
	familyNoSpaces := strings.ReplaceAll(family, " ", "")
	suffix := weightToSuffix(weight)
	if style == "italic" && suffix != "" {
		suffix += "Italic"
	} else if style == "italic" {
		suffix = "Italic"
	}
	if suffix == "" {
		suffix = "Regular"
	}
	return familyNoSpaces + "-" + suffix + ".ttf"
}

func (d *GoogleFontDownloader) fetchTTFURL(cssURL string) (string, error) {
	req, err := http.NewRequest("GET", cssURL, nil)
	if err != nil {
		return "", err
	}
	// User-Agent must request TTF format (not woff2)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; pen2pdf)")

	resp, err := d.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("google fonts returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	matches := ttfURLRegex.FindSubmatch(body)
	if matches == nil {
		return "", fmt.Errorf("no TTF URL found in Google Fonts response")
	}

	return string(matches[1]), nil
}

func (d *GoogleFontDownloader) fetchBytes(rawURL string) ([]byte, error) {
	resp, err := d.client.Get(rawURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
