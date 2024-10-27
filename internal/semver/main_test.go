package semver

import "testing"

func TestFromString(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected SemVer
	}{
		"simple, with prefix": {
			input:    "v1.2.3",
			expected: SemVer{Prefix: "v", Major: 1, Minor: 2, Patch: 3, Rc: 0, CommitDistance: 0},
		},
		"simple, no prefix": {
			input:    "1.2.3",
			expected: SemVer{Prefix: "", Major: 1, Minor: 2, Patch: 3, Rc: 0, CommitDistance: 0},
		},
		"release candidate": {
			input:    "1.2.3-rc4",
			expected: SemVer{Prefix: "", Major: 1, Minor: 2, Patch: 3, Rc: 4, CommitDistance: 0},
		},
		"commit distance": {
			input:    "1.2.3-5-gabcdef",
			expected: SemVer{Prefix: "", Major: 1, Minor: 2, Patch: 3, Rc: 0, CommitDistance: 5},
		},
		"release candidate and commit distance": {
			input:    "1.2.3-rc4-5-gabcdef",
			expected: SemVer{Prefix: "", Major: 1, Minor: 2, Patch: 3, Rc: 4, CommitDistance: 5},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual, err := FromString(test.input)
			if err != nil {
				t.Fatalf("FromString(%q) failed: %q", test.input, err)
			}
			if actual != test.expected {
				t.Fatalf("FromString(%q) returned %q; expected %q", test.input, &actual, &test.expected)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := map[string]struct {
		input    SemVer
		expected string
	}{
		"simple, with prefix": {
			input:    SemVer{Prefix: "v", Major: 1, Minor: 2, Patch: 3, Rc: 0, CommitDistance: 0},
			expected: "v1.2.3",
		},
		"simple, no prefix": {
			input:    SemVer{Prefix: "", Major: 1, Minor: 2, Patch: 3, Rc: 0, CommitDistance: 0},
			expected: "1.2.3",
		},
		"release candidate": {
			input:    SemVer{Prefix: "", Major: 1, Minor: 2, Patch: 3, Rc: 4, CommitDistance: 0},
			expected: "1.2.3-rc4",
		},
		"commit distance": {
			input:    SemVer{Prefix: "", Major: 1, Minor: 2, Patch: 3, Rc: 0, CommitDistance: 5},
			expected: "1.2.3",
		},
		"release candidate and commit distance": {
			input:    SemVer{Prefix: "", Major: 1, Minor: 2, Patch: 3, Rc: 4, CommitDistance: 5},
			expected: "1.2.3-rc4",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual := test.input.String()
			if actual != test.expected {
				t.Fatalf("String(%+v) returned %q; expected %q", test.input, actual, test.expected)
			}
		})
	}
}

func TestBump(t *testing.T) {
	tests := map[string]struct {
		initial  SemVer
		major    bool
		minor    bool
		patch    bool
		rc       bool
		noRc     bool
		expected SemVer
	}{
		// major/minor/patch flag only
		"major bump, with zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 0},
			major:    true,
			expected: SemVer{Major: 2, Minor: 0, Patch: 0, Rc: 0},
		},
		"major bump, with non-zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 4},
			major:    true,
			expected: SemVer{Major: 2, Minor: 0, Patch: 0, Rc: 0},
		},
		"minor bump, with zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 0},
			minor:    true,
			expected: SemVer{Major: 1, Minor: 3, Patch: 0, Rc: 0},
		},
		"minor bump, with non-zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 4},
			minor:    true,
			expected: SemVer{Major: 1, Minor: 3, Patch: 0, Rc: 0},
		},
		"patch bump, with zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 0},
			patch:    true,
			expected: SemVer{Major: 1, Minor: 2, Patch: 4, Rc: 0},
		},
		"patch bump, with non-zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 4},
			patch:    true,
			expected: SemVer{Major: 1, Minor: 2, Patch: 4, Rc: 0},
		},

		// rc flag only
		"rc bump, from zero": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 0},
			rc:       true,
			expected: SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 1},
		},
		"rc bump, from non-zero": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 1},
			rc:       true,
			expected: SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 2},
		},

		// major/minor/patch flag AND rc flag
		"major bump + rc, with zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 0},
			major:    true,
			rc:       true,
			expected: SemVer{Major: 2, Minor: 0, Patch: 0, Rc: 1},
		},
		"major bump + rc, with non-zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 4},
			major:    true,
			rc:       true,
			expected: SemVer{Major: 2, Minor: 0, Patch: 0, Rc: 1},
		},
		"minor bump + rc, with zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 0},
			minor:    true,
			rc:       true,
			expected: SemVer{Major: 1, Minor: 3, Patch: 0, Rc: 1},
		},
		"minor bump + rc, with non-zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 4},
			minor:    true,
			rc:       true,
			expected: SemVer{Major: 1, Minor: 3, Patch: 0, Rc: 1},
		},
		"patch bump + rc, with zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 0},
			patch:    true,
			rc:       true,
			expected: SemVer{Major: 1, Minor: 2, Patch: 4, Rc: 1},
		},
		"patch bump + rc, with non-zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 4},
			patch:    true,
			rc:       true,
			expected: SemVer{Major: 1, Minor: 2, Patch: 4, Rc: 1},
		},

		// no-rc flag only
		"no-rc, with zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 0},
			noRc:     true,
			expected: SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 0},
		},
		"no-rc, with non-zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 4},
			noRc:     true,
			expected: SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 0},
		},

		// major/minor/patch flag AND no-rc flag
		"major bump + no-rc, with zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 0},
			major:    true,
			noRc:     true,
			expected: SemVer{Major: 2, Minor: 0, Patch: 0, Rc: 0},
		},
		"major bump + no-rc, with non-zero rc": {
			initial:  SemVer{Major: 1, Minor: 2, Patch: 3, Rc: 4},
			major:    true,
			noRc:     true,
			expected: SemVer{Major: 2, Minor: 0, Patch: 0, Rc: 0},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual := test.initial.Bump(test.major, test.minor, test.patch, test.rc, test.noRc)
			if actual != test.expected {
				t.Fatalf("%+v.Bump(%v, %v, %v, %v, %v) returned %+v, expected %+v", test.initial, test.major, test.minor, test.patch, test.rc, test.noRc, actual, test.expected)
			}
		})
	}
}
