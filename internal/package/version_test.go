package hapmpkg

import "testing"

func TestParseVersion(t *testing.T) {
	cases := []struct {
		in     string
		value  []int
		suffix []string
	}{
		{"1.0.0", []int{1, 0, 0}, nil},
		{"v1.0.0.0.0", []int{1, 0, 0, 0, 0}, nil},
		{"1.0.0-alpha", []int{1, 0, 0}, []string{"alpha"}},
		{"1.0.0.alpha.1", []int{1, 0, 0}, []string{"alpha", "1"}},
		{"v1-alpha.1", []int{1}, []string{"alpha", "1"}},
		{"1", []int{1}, nil},
	}

	for _, tc := range cases {
		parts, err := ParseVersion(tc.in)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.in, err)
		}
		if len(parts.Value) != len(tc.value) {
			t.Fatalf("unexpected value length for %q", tc.in)
		}
		for i := range tc.value {
			if parts.Value[i] != tc.value[i] {
				t.Fatalf("unexpected value for %q", tc.in)
			}
		}
		if len(parts.Suffix) != len(tc.suffix) {
			t.Fatalf("unexpected suffix length for %q", tc.in)
		}
		for i := range tc.suffix {
			if parts.Suffix[i] != tc.suffix[i] {
				t.Fatalf("unexpected suffix for %q", tc.in)
			}
		}
	}
}

func TestParseVersionInvalid(t *testing.T) {
	invalid := []string{"r1.0.0", "", "hello", "v1..0", "v1.0.0.", "1.0.0-"}
	for _, value := range invalid {
		if _, err := ParseVersion(value); err == nil {
			t.Fatalf("expected error for %q", value)
		}
	}
}

func TestCompareVersions(t *testing.T) {
	ltCases := [][2]string{
		{"1.0.0", "1.0.1"},
		{"1", "1.1.0"},
		{"1-beta1", "1-beta2"},
		{"1.0.0", "2.0.0-alpha.1"},
		{"0.1.9.1", "0.1.12.4"},
	}
	for _, pair := range ltCases {
		left := MustNewVersion(pair[0])
		right := MustNewVersion(pair[1])
		if left.Compare(right) >= 0 {
			t.Fatalf("expected %s < %s", pair[0], pair[1])
		}
	}

	gtCases := [][2]string{
		{"1.0.0", "1.0.0-alpha.1"},
		{"1.0.0", "1.0.0-rc.1"},
	}
	for _, pair := range gtCases {
		left := MustNewVersion(pair[0])
		right := MustNewVersion(pair[1])
		if left.Compare(right) <= 0 {
			t.Fatalf("expected %s > %s", pair[0], pair[1])
		}
	}
}

func TestFindLatest(t *testing.T) {
	tags := []string{"v1.0.0", "bad", "v1.2.0-rc.1", "v1.1.1"}
	if got := FindLatestVersion(tags, true); got != "v1.1.1" {
		t.Fatalf("unexpected stable latest: %s", got)
	}
	if got := FindLatestVersion(tags, false); got != "v1.2.0-rc.1" {
		t.Fatalf("unexpected latest: %s", got)
	}
}
