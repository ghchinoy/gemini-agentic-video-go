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

package runner

import (
	"testing"

	"google.golang.org/genai"
)

func TestParseThinkingLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected genai.ThinkingLevel
	}{
		{"low", genai.ThinkingLevelLow},
		{"LOW", genai.ThinkingLevelLow},
		{"medium", genai.ThinkingLevelMedium},
		{"high", genai.ThinkingLevelHigh},
		{"minimal", genai.ThinkingLevelMinimal},
		{"unknown", genai.ThinkingLevelMedium},
	}

	for _, tt := range tests {
		got := ParseThinkingLevel(tt.input)
		if got != tt.expected {
			t.Errorf("ParseThinkingLevel(%q) = %v; want %v", tt.input, got, tt.expected)
		}
	}
}

func TestParseProcessingMode(t *testing.T) {
	tests := []struct {
		input    string
		expected genai.MediaProcessing
	}{
		{"agentic", genai.MediaProcessingAgentic},
		{"AGENTIC", genai.MediaProcessingAgentic},
		{"static", genai.MediaProcessingStatic},
		{"STATIC", genai.MediaProcessingStatic},
		{"anything_else", genai.MediaProcessingAgentic},
	}

	for _, tt := range tests {
		got := ParseProcessingMode(tt.input)
		if got != tt.expected {
			t.Errorf("ParseProcessingMode(%q) = %v; want %v", tt.input, got, tt.expected)
		}
	}
}

func TestBenchmarkResultStruct(t *testing.T) {
	res := BenchmarkResult{
		AgenticResult: &Result{
			Title: "Agentic",
			Mode:  genai.MediaProcessingAgentic,
		},
		StaticResult: &Result{
			Title: "Static",
			Mode:  genai.MediaProcessingStatic,
		},
	}

	if res.AgenticResult.Mode != genai.MediaProcessingAgentic {
		t.Errorf("res.AgenticResult.Mode = %v, want Agentic", res.AgenticResult.Mode)
	}
	if res.StaticResult.Mode != genai.MediaProcessingStatic {
		t.Errorf("res.StaticResult.Mode = %v, want Static", res.StaticResult.Mode)
	}
}
