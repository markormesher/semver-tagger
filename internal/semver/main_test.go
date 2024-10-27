package semver

import "testing"

func TestFromString(t *testing.T) {
	tests := map[string]struct {
		input  string
		expect SemVer
	}{
		"simple, with prefix": {
			input:  "v1.2.3",
			expect: SemVer{Prefix: "v", Major: 1, Minor: 2, Patch: 3, Rc: 0, CommitDistance: 0},
		},
		"simple, no prefix": {
			input:  "1.2.3",
			expect: SemVer{Prefix: "", Major: 1, Minor: 2, Patch: 3, Rc: 0, CommitDistance: 0},
		},
		"release candidate": {
			input:  "1.2.3-rc4",
			expect: SemVer{Prefix: "", Major: 1, Minor: 2, Patch: 3, Rc: 4, CommitDistance: 0},
		},
		"commit description": {
			input:  "1.2.3-5-gabcdef",
			expect: SemVer{Prefix: "", Major: 1, Minor: 2, Patch: 3, Rc: 0, CommitDistance: 5},
		},
		"release candidate and commit description": {
			input:  "1.2.3-rc4-5-gabcdef",
			expect: SemVer{Prefix: "", Major: 1, Minor: 2, Patch: 3, Rc: 4, CommitDistance: 5},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual, err := FromString(test.input)
			if err != nil {
				t.Fatalf("FromString(%q) failed: %q", test.input, err)
			}
			if actual != test.expect {
				t.Fatalf("FromString(%q) returned %q; expected %q", test.input, &actual, &test.expect)
			}
		})
	}
}
