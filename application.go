// Package authentik generates and validates isolated Authentik application
// enrollment configuration without owning an Authentik installation.
package authentik

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// Application is the desired Authentik enrollment and access boundary for one
// consuming application.
type Application struct {
	Slug              string `json:"slug"`
	Name              string `json:"name"`
	RegistrationTitle string `json:"registration_title"`
	UserPath          string `json:"user_path"`
	AccessGroup       string `json:"access_group"`
}

// Validate rejects ambiguous identifiers, unsafe display text, and user paths
// outside Authentik's ordinary users hierarchy.
func (application Application) Validate() error {
	if !validSlug(application.Slug) {
		return fmt.Errorf("application slug is invalid")
	}
	if !validSlug(application.AccessGroup) {
		return fmt.Errorf("access group is invalid")
	}
	if !validText(application.Name, 1, 120) || !validText(application.RegistrationTitle, 1, 160) {
		return fmt.Errorf("application display text is invalid")
	}
	if !strings.HasPrefix(application.UserPath, "users/") || len(application.UserPath) > 160 {
		return fmt.Errorf("user path must be a bounded users/ path")
	}
	for _, segment := range strings.Split(application.UserPath, "/") {
		if !validSlug(segment) {
			return fmt.Errorf("user path contains an invalid segment")
		}
	}
	return nil
}

func validSlug(value string) bool {
	return len(value) <= 80 && slugPattern.MatchString(value)
}

func validText(value string, minimum, maximum int) bool {
	if !utf8.ValidString(value) {
		return false
	}
	length := utf8.RuneCountInString(value)
	return length >= minimum && length <= maximum && strings.IndexFunc(value, unicode.IsControl) < 0
}

// ValidateIssuerURL accepts one exact canonical Authentik OIDC issuer URL.
func ValidateIssuerURL(raw string) error {
	parsed, err := validateCanonicalURL(raw)
	if err != nil {
		return fmt.Errorf("invalid Authentik issuer: %w", err)
	}
	if parsed.Path == "/" || !strings.HasSuffix(parsed.Path, "/") {
		return fmt.Errorf("invalid Authentik issuer: path must identify a provider and end with a slash")
	}
	return nil
}

// ValidateEnrollmentURL requires the exact canonical flow URL on the expected
// Authentik origin.
func ValidateEnrollmentURL(raw, expectedOrigin, flowSlug string) error {
	if !validSlug(flowSlug) {
		return fmt.Errorf("invalid enrollment flow slug")
	}
	parsed, err := validateCanonicalURL(raw)
	if err != nil {
		return fmt.Errorf("invalid Authentik enrollment URL: %w", err)
	}
	origin, err := validateCanonicalURL(expectedOrigin)
	if err != nil || origin.Path != "/" {
		return fmt.Errorf("invalid expected Authentik origin")
	}
	wantPath := "/if/flow/" + flowSlug + "/"
	if parsed.Scheme != origin.Scheme || parsed.Host != origin.Host || parsed.Path != wantPath {
		return fmt.Errorf("Authentik enrollment URL does not match the expected origin and flow")
	}
	return nil
}

func validateCanonicalURL(raw string) (*url.URL, error) {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Opaque != "" || parsed.Host == "" || parsed.Hostname() == "" || parsed.Path == "" {
		return nil, fmt.Errorf("URL is not absolute and hierarchical")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("URL must use HTTP or HTTPS")
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.ForceQuery || parsed.Fragment != "" || parsed.String() != raw {
		return nil, fmt.Errorf("URL contains credentials, query, fragment, or noncanonical encoding")
	}
	return parsed, nil
}

// ValidateIsolation rejects sibling applications that accidentally share an
// application slug or access group.
func ValidateIsolation(applications []Application) error {
	seenSlugs := make(map[string]struct{}, len(applications))
	seenGroups := make(map[string]struct{}, len(applications))
	for index, application := range applications {
		if err := application.Validate(); err != nil {
			return fmt.Errorf("application %d: %w", index, err)
		}
		if _, exists := seenSlugs[application.Slug]; exists {
			return fmt.Errorf("application slug %q is duplicated", application.Slug)
		}
		if _, exists := seenGroups[application.AccessGroup]; exists {
			return fmt.Errorf("access group %q is shared by sibling applications", application.AccessGroup)
		}
		seenSlugs[application.Slug] = struct{}{}
		seenGroups[application.AccessGroup] = struct{}{}
	}
	return nil
}
