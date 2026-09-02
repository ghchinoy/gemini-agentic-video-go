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

package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/ghchinoy/gemini-agentic-video-go/internal/catalog"
	"google.golang.org/genai"
)

func TestRenderBanner(t *testing.T) {
	out := RenderBanner("gemini-3.7-flash", "enterprise", "test-project", "global")
	if !strings.Contains(out, "gemini-3.7-flash") {
		t.Errorf("RenderBanner() missing model ID; output: %s", out)
	}
	if !strings.Contains(out, "test-project") {
		t.Errorf("RenderBanner() missing project; output: %s", out)
	}
}

func TestRenderCatalogTable(t *testing.T) {
	out := RenderCatalogTable(catalog.Scenarios)
	if out == "" {
		t.Errorf("RenderCatalogTable() returned empty string, want formatted table")
	}
	if !strings.Contains(out, "youtube-agentic") {
		t.Errorf("RenderCatalogTable() missing scenario name; output: %s", out)
	}
}

func TestRenderBenchmarkTable(t *testing.T) {
	agenticUsage := &genai.GenerateContentResponseUsageMetadata{
		TotalTokenCount:      8000,
		PromptTokenCount:     250,
		CandidatesTokenCount: 500,
		ThoughtsTokenCount:   7250,
	}
	staticUsage := &genai.GenerateContentResponseUsageMetadata{
		TotalTokenCount:      185000,
		PromptTokenCount:     184500,
		CandidatesTokenCount: 500,
		ThoughtsTokenCount:   0,
	}

	out := RenderBenchmarkTable(agenticUsage, staticUsage, 15*time.Second, 45*time.Second)
	if !strings.Contains(out, "PERFORMANCE BENCHMARK") {
		t.Errorf("RenderBenchmarkTable() missing title header; output: %s", out)
	}

	// Test graceful handling of nil usage metadata
	nilOut := RenderBenchmarkTable(nil, nil, 0, 0)
	if nilOut == "" {
		t.Errorf("RenderBenchmarkTable(nil, nil) returned empty string, want fallback table")
	}
}

func TestRenderTelemetryCard(t *testing.T) {
	usage := &genai.GenerateContentResponseUsageMetadata{
		TotalTokenCount:      5000,
		PromptTokenCount:     200,
		CandidatesTokenCount: 300,
		ThoughtsTokenCount:   4500,
	}

	out := RenderTelemetryCard(usage, genai.MediaProcessingAgentic, 10*time.Second)
	if !strings.Contains(out, "5000") {
		t.Errorf("RenderTelemetryCard() missing total token count; output: %s", out)
	}

	// Test graceful handling of nil usage
	nilOut := RenderTelemetryCard(nil, genai.MediaProcessingAgentic, 10*time.Second)
	if nilOut == "" {
		t.Errorf("RenderTelemetryCard(nil) returned empty string, want fallback card")
	}
}

func TestRenderBadges(t *testing.T) {
	agenticBadge := RenderModeBadge(genai.MediaProcessingAgentic)
	if !strings.Contains(agenticBadge, "AGENTIC") {
		t.Errorf("RenderModeBadge(Agentic) = %q, want AGENTIC", agenticBadge)
	}

	staticBadge := RenderModeBadge(genai.MediaProcessingStatic)
	if !strings.Contains(staticBadge, "STATIC") {
		t.Errorf("RenderModeBadge(Static) = %q, want STATIC", staticBadge)
	}
}
