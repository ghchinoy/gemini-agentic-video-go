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

package catalog

import (
	"testing"
)

func TestCatalogScenarios(t *testing.T) {
	if len(Scenarios) < 6 {
		t.Fatalf("expected at least 6 scenarios, got %d", len(Scenarios))
	}

	for i := 1; i <= len(Scenarios); i++ {
		sc, err := GetByID(i)
		if err != nil {
			t.Errorf("GetByID(%d) error: %v", i, err)
			continue
		}
		if sc.ID != i {
			t.Errorf("scenario ID mismatch: got %d, want %d", sc.ID, i)
		}
		if sc.Title == "" {
			t.Errorf("scenario %d missing Title", i)
		}
	}

	// Test invalid scenario ID
	_, err := GetByID(999)
	if err == nil {
		t.Errorf("expected error for non-existent scenario ID, got nil")
	}
}
