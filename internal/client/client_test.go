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

package client

import (
	"os"
	"testing"
)

func TestResolveLocation(t *testing.T) {
	tests := []struct {
		name      string
		flagVal   string
		envLoc    string
		envReg    string
		expected  string
	}{
		{
			name:     "flag provided",
			flagVal:  "us-central1",
			expected: "us-central1",
		},
		{
			name:     "env location set",
			flagVal:  "",
			envLoc:   "europe-west1",
			expected: "europe-west1",
		},
		{
			name:     "fallback to global",
			flagVal:  "",
			expected: "global",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envLoc != "" {
				os.Setenv("GOOGLE_CLOUD_LOCATION", tt.envLoc)
				defer os.Unsetenv("GOOGLE_CLOUD_LOCATION")
			}
			got := ResolveLocation(tt.flagVal)
			if got != tt.expected {
				t.Errorf("ResolveLocation(%q) = %q; want %q", tt.flagVal, got, tt.expected)
			}
		})
	}
}

func TestResolveProject(t *testing.T) {
	got := ResolveProject("custom-project")
	if got != "custom-project" {
		t.Errorf("ResolveProject(custom-project) = %q; want custom-project", got)
	}
}
