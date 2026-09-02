// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Package telemetry provides metrics formatting and benchmark comparison for video processing.
package telemetry

import (
	"fmt"
	"time"

	"google.golang.org/genai"
)

// PrintUsage displays a structured breakdown of token consumption.
func PrintUsage(usage *genai.GenerateContentResponseUsageMetadata, mode genai.MediaProcessing, duration time.Duration) {
	if usage == nil {
		fmt.Printf("Processing Duration: %v (token usage telemetry not returned)\n\n", duration.Round(time.Millisecond))
		return
	}

	fmt.Printf("Token Telemetry Breakdown:\n")
	fmt.Printf("  • Total Token Count:     %d\n", usage.TotalTokenCount)
	fmt.Printf("  • Prompt Input Tokens:   %d", usage.PromptTokenCount)
	if mode == genai.MediaProcessingAgentic {
		fmt.Printf(" (text-only prompt; video frames dynamically fetched on-demand)")
	} else if mode == genai.MediaProcessingStatic {
		fmt.Printf(" (includes 100%% of video frames statically pre-ingested at 1 FPS)")
	}
	fmt.Println()

	fmt.Printf("  • Candidates Tokens:     %d\n", usage.CandidatesTokenCount)
	if usage.ThoughtsTokenCount > 0 {
		fmt.Printf("  • Thoughts Tokens:       %d", usage.ThoughtsTokenCount)
		if mode == genai.MediaProcessingAgentic {
			fmt.Printf(" (includes native timeline navigation & dynamic inspection)")
		}
		fmt.Println()
	}
	fmt.Printf("Processing Duration:       %v\n\n", duration.Round(time.Millisecond))
}

// PrintComparison displays a side-by-side performance comparison table.
func PrintComparison(agenticUsage, staticUsage *genai.GenerateContentResponseUsageMetadata, agenticDuration, staticDuration time.Duration) {
	fmt.Printf("=========================================================================================\n")
	fmt.Printf("  PERFORMANCE BENCHMARK: AGENTIC VIDEO vs. STATIC (1 FPS) INGESTION\n")
	fmt.Printf("=========================================================================================\n")

	agenticTotal := int64(0)
	staticTotal := int64(0)
	agenticPrompt := int32(0)
	staticPrompt := int32(0)
	agenticCandidates := int32(0)
	staticCandidates := int32(0)
	agenticThoughts := int32(0)
	staticThoughts := int32(0)

	if agenticUsage != nil {
		agenticTotal = int64(agenticUsage.TotalTokenCount)
		agenticPrompt = agenticUsage.PromptTokenCount
		agenticCandidates = agenticUsage.CandidatesTokenCount
		agenticThoughts = agenticUsage.ThoughtsTokenCount
	}
	if staticUsage != nil {
		staticTotal = int64(staticUsage.TotalTokenCount)
		staticPrompt = staticUsage.PromptTokenCount
		staticCandidates = staticUsage.CandidatesTokenCount
		staticThoughts = staticUsage.ThoughtsTokenCount
	}

	var tokenDeltaPct float64
	if staticTotal > 0 {
		tokenDeltaPct = float64(staticTotal-agenticTotal) / float64(staticTotal) * 100.0
	}

	fmt.Printf("| Metric                     | Agentic Video        | Static Ingestion     | Observation / Delta       |\n")
	fmt.Printf("|:---------------------------|:---------------------|:---------------------|:--------------------------|\n")
	fmt.Printf("| Total Consumed Tokens      | %-20d | %-20d | %+.1f%% net token spend   |\n", agenticTotal, staticTotal, tokenDeltaPct)
	fmt.Printf("| Prompt Tokens (Input)      | %-20d | %-20d | %-+25s |\n", agenticPrompt, staticPrompt, fmt.Sprintf("%+d input tokens", agenticPrompt-staticPrompt))
	fmt.Printf("| Candidate Output Tokens    | %-20d | %-20d | %-+25s |\n", agenticCandidates, staticCandidates, fmt.Sprintf("%+d output tokens", agenticCandidates-staticCandidates))
	fmt.Printf("| Reasoning Thoughts Tokens  | %-20d | %-20d | %-+25s |\n", agenticThoughts, staticThoughts, "dynamic frame inspection")
	fmt.Printf("| Latency (Duration)         | %-20v | %-20v | %-+25s |\n",
		agenticDuration.Round(time.Millisecond),
		staticDuration.Round(time.Millisecond),
		fmt.Sprintf("diff: %v", (agenticDuration-staticDuration).Round(time.Millisecond)))
	fmt.Printf("=========================================================================================\n")
	fmt.Printf("💡 Telemetry Insight: In Agentic mode, initial prompt tokens are minimal because video frames\n")
	fmt.Printf("   are fetched dynamically during the model's Think ➔ Act timeline loop (accounted under thoughts).\n\n")
}
