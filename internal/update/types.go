package update

import (
	"fmt"
	"sort"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func ParseVersion(v string) (Version, error) {
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("invalid version format: %s", v)
	}

	var version Version
	for i, part := range parts {
		var n int
		if _, err := fmt.Sscanf(part, "%d", &n); err != nil {
			return Version{}, fmt.Errorf("invalid version component: %s", part)
		}
		switch i {
		case 0:
			version.Major = n
		case 1:
			version.Minor = n
		case 2:
			version.Patch = n
		}
	}

	return version, nil
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v Version) Less(other Version) bool {
	if v.Major != other.Major {
		return v.Major < other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor < other.Minor
	}
	return v.Patch < other.Patch
}

func (v Version) CompatibleWith(other Version) bool {
	return v.Major == other.Major
}

type ComponentVersion struct {
	Core           Version
	ConfigSchema   int
	DatabaseSchema int
	PluginAPI      int
}

type VersionResolver struct {
	versions map[string]Version
}

func NewVersionResolver() *VersionResolver {
	return &VersionResolver{
		versions: make(map[string]Version),
	}
}

func (r *VersionResolver) Register(name string, version Version) {
	r.versions[name] = version
}

func (r *VersionResolver) Resolve(name string) (Version, bool) {
	version, ok := r.versions[name]
	return version, ok
}

func (r *VersionResolver) Compare(name string, target Version) (int, error) {
	current, ok := r.versions[name]
	if !ok {
		return 0, fmt.Errorf("component not found: %s", name)
	}

	if current.Less(target) {
		return -1, nil
	}
	if target.Less(current) {
		return 1, nil
	}
	return 0, nil
}

func (r *VersionResolver) AvailableUpdates() []string {
	var updates []string
	for name, current := range r.versions {
		if current.Major == 0 && current.Minor == 0 && current.Patch == 0 {
			continue
		}
		updates = append(updates, name)
	}
	sort.Strings(updates)
	return updates
}
