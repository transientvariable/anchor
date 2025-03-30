package anchor

import "regexp"

var (
	EmailPattern       *regexp.Regexp
	IPv4Pattern        *regexp.Regexp
	StoragePathPattern *regexp.Regexp
	URISchemePattern   *regexp.Regexp
	UsernamePattern    *regexp.Regexp
)

func init() {
	EmailPattern = regexp.MustCompile(`^[a-zA-Z0-9_+&*-]+(?:\.[a-zA-Z0-9_+&*-]+)*@(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,7}$`)
	IPv4Pattern = regexp.MustCompile(`(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)
	StoragePathPattern = regexp.MustCompile(`([a-z0-9\-._~/]*)`)
	URISchemePattern = regexp.MustCompile(`^([a-z][a-z0-9+\-.]*):`)
	UsernamePattern = regexp.MustCompile(`^[a-zA-Z0-9]{5,15}$`)
}
