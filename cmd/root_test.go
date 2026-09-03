// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"testing"
)

func TestSubcommandRegistration(t *testing.T) {
	subcommands := []string{"example", "run", "compare", "multivideo", "multiturn"}

	for _, sub := range subcommands {
		cmd, _, err := RootCmd.Find([]string{sub})
		if err != nil {
			t.Errorf("RootCmd.Find(%q) error = %v, want nil", sub, err)
			continue
		}
		if cmd == nil || cmd.Name() != sub {
			t.Errorf("RootCmd.Find(%q) = %v, want command with name %q", sub, cmd, sub)
		}
	}
}

func TestRunCommandRequiredFlags(t *testing.T) {
	// Execute run command with no arguments or flags to verify --video requirement
	runVideoURI = ""
	err := runCmd.RunE(runCmd, []string{})
	if err == nil {
		t.Errorf("runCmd.RunE() with empty video URI = nil, want error")
	}
}

func TestResolveModelsList(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", nil},
		{"flash", []string{"gemini-3.6-flash", "gemini-3.7-flash", "gemini-3.8-flash"}},
		{"all", []string{"gemini-3.6-flash", "gemini-3.7-flash", "gemini-3.8-flash"}},
		{"3.6,3.8", []string{"gemini-3.6-flash", "gemini-3.8-flash"}},
		{"gemini-3.7-flash,gemini-3.8-flash", []string{"gemini-3.7-flash", "gemini-3.8-flash"}},
	}

	for _, tt := range tests {
		got := resolveModelsList(tt.input)
		if len(got) != len(tt.expected) {
			t.Errorf("resolveModelsList(%q) len = %d, want %d", tt.input, len(got), len(tt.expected))
			continue
		}
		for i := range got {
			if got[i] != tt.expected[i] {
				t.Errorf("resolveModelsList(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.expected[i])
			}
		}
	}
}
