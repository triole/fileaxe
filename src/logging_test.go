package main

import "testing"

func TestLogTarget(t *testing.T) {
	cases := []struct {
		file   string
		dryRun bool
		want   string
	}{
		{"/var/log/fileaxe.log", true, "stdout"},
		{"/dev/stdout", true, "stdout"},
		{"", true, "stdout"},
		{"/var/log/fileaxe.log", false, "/var/log/fileaxe.log"},
		{"/dev/stdout", false, "/dev/stdout"},
		{"", false, ""},
	}
	for _, c := range cases {
		if got := logTarget(c.file, c.dryRun); got != c.want {
			t.Errorf(
				"logTarget(%q, dryRun=%v): got %q, want %q",
				c.file, c.dryRun, got, c.want,
			)
		}
	}
}
