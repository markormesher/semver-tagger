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
