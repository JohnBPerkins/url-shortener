package service

import (
	"strings"
	"testing"
)

func TestEncodeBase62(t *testing.T) {
	type tc struct {
		num  uint64
		want string
	}
	tests := []tc{
		{num: 0, want: strings.Repeat("0", codeLength)},
		{num: 1, want: strings.Repeat("0", codeLength-1) + "1"},
		{num: 61, want: strings.Repeat("0", codeLength-1) + string(base62Chars[61])},
		{num: func() uint64 {
			var v uint64 = 1
			for i := 0; i < codeLength-1; i++ {
				v *= 62
			}
			return v
		}(), want: "10000000"},
	}

	for _, tc := range tests {
		got := encodeBase62(tc.num)
		if got != tc.want {
			t.Errorf("encodeBase62(%d) = %q; want %q", tc.num, got, tc.want)
		}
	}

	oversize := func() uint64 {
		var v uint64 = 1
		for i := 0; i < codeLength+1; i++ {
			v *= 62
		}
		return v
	}()
	got := encodeBase62(oversize)
	if len(got) != codeLength {
		t.Errorf("encodeBase62(oversize) length = %d; want %d", len(got), codeLength)
	}
}

func TestIsValidURL_ASCII(t *testing.T) {
	// Contains an accented 'á'
	nonASCII := "http://exámple.com"
	if isValidURL(nonASCII) {
		t.Errorf("expected non-ASCII URL to be invalid, got valid: %q", nonASCII)
	}
	// Control character (newline) embedded
	controlChar := "http://example.com/foo\nbar"
	if isValidURL(controlChar) {
		t.Errorf("expected URL with control char to be invalid, got valid: %q", controlChar)
	}
}

func TestIsValidURL_Regex(t *testing.T) {
	longInput := strings.Repeat("a", maxURLLength+1)

	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"http://example.com", true},
		{"https://sub.domain.co.uk/path?query=1#frag", true},
		{"example.com", true},
		{"www.google.com", true},
		{"ftp://example.com", false},
		{"://missing.scheme.com", false},
		{"http:///nohost", false},
		{"http://toolong." + strings.Repeat("a", maxURLLength), false},
		{longInput, false},
	}

	for _, tt := range tests {
		got := isValidURL(tt.input)
		if got != tt.want {
			t.Errorf("isValidURL(%q) = %v; want %v", tt.input, got, tt.want)
		}
	}
}

func TestIsValidURL_LengthBoundary(t *testing.T) {
    const hostAndSlash = "example.com/"
    if len(hostAndSlash) >= maxURLLength {
        t.Fatalf("hostAndSlash length (%d) >= maxURLLength (%d)", len(hostAndSlash), maxURLLength)
    }

    pathLen := maxURLLength - len(hostAndSlash)
    valid := hostAndSlash + strings.Repeat("a", pathLen)
    if len(valid) != maxURLLength {
        t.Fatalf("setup error: valid URL length = %d; want %d", len(valid), maxURLLength)
    }

    if !isValidURL(valid) {
        t.Errorf("expected boundary-length ASCII URL to be valid, got invalid: %q", valid)
    }
}

func TestIsValidURL_TooLong(t *testing.T) {
    tooLong := strings.Repeat("a", maxURLLength+1)
    if len(tooLong) != maxURLLength+1 {
        t.Fatalf("setup error: string length = %d; want %d", len(tooLong), maxURLLength+1)
    }
    if isValidURL(tooLong) {
        t.Errorf("expected URL longer than %d to be invalid, got valid", maxURLLength)
    }
}