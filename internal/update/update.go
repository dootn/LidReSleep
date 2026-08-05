//go:build windows

// Package update checks for newer releases on GitHub. It queries the repository
// releases API (latest tag), compares it against the running version and reports
// whether an update is available.
package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// repo is the GitHub API endpoint for the latest release. The API requires a
// User-Agent header (otherwise it replies 403), which Check sets.
const repoURL = "https://api.github.com/repos/dootn/LidReSleep/releases/latest"

// release mirrors the fields of the GitHub release API response we need.
type release struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

// Result of a version check.
type Result struct {
	Current string // running version, e.g. "1.0.0"
	Latest  string // latest release tag, e.g. "v1.0.0"
	URL     string // release page (fallback when html_url is empty)
	Update  bool   // true when Latest is newer than Current
}

// Check queries GitHub for the latest release within timeout and reports whether
// it is newer than current.
func Check(current string, timeout time.Duration) (Result, error) {
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest(http.MethodGet, repoURL, nil)
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("User-Agent", "LidReSleep")

	resp, err := client.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var rel release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return Result{}, err
	}

	r := Result{Current: current, Latest: rel.TagName, Update: Newer(rel.TagName, current)}
	if rel.HTMLURL != "" {
		r.URL = rel.HTMLURL
	} else if r.Latest != "" {
		r.URL = "https://github.com/dootn/LidReSleep/releases/tag/" + r.Latest
	} else {
		r.URL = "https://github.com/dootn/LidReSleep/releases/latest"
	}
	return r, nil
}

// Newer reports whether tag (e.g. "v1.2.0") is a newer version than current
// (e.g. "1.0.0"). Tags that cannot be parsed are never considered newer.
func Newer(tag, current string) bool {
	va, okA := parseVersion(tag)
	vb, okB := parseVersion(current)
	if !okA {
		return false
	}
	if !okB {
		return true
	}
	for i := 0; i < len(va) && i < len(vb); i++ {
		if va[i] != vb[i] {
			return va[i] > vb[i]
		}
	}
	return len(va) > len(vb)
}

// parseVersion splits a dotted numeric version ("v1.2.3", "1.0.0-pre") into its
// numeric components, ignoring a leading "v" and any pre-release/build suffix.
func parseVersion(s string) ([]int, bool) {
	s = strings.TrimPrefix(strings.TrimPrefix(s, "v"), "V")
	parts := strings.Split(s, ".")
	if len(parts) == 0 {
		return nil, false
	}
	nums := make([]int, 0, len(parts))
	for _, p := range parts {
		if i := strings.IndexAny(p, "-+"); i >= 0 {
			p = p[:i]
		}
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return nil, false
		}
		nums = append(nums, n)
	}
	return nums, true
}
