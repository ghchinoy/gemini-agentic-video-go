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
func PrintComparison(aUsage, sUsage *genai.GenerateContentResponseUsageMetadata, aDur, sDur time.Duration) {
	fmt.Printf("=========================================================================================\n")
	fmt.Printf("  PERFORMANCE BENCHMARK: AGENTIC VIDEO vs. STATIC (1 FPS) INGESTION\n")
	fmt.Printf("=========================================================================================\n")

	aTotal := int64(0)
	sTotal := int64(0)
	aPrompt := int32(0)
	sPrompt := int32(0)
	aCand := int32(0)
	sCand := int32(0)
	aThoughts := int32(0)
	sThoughts := int32(0)

	if aUsage != nil {
		aTotal = int64(aUsage.TotalTokenCount)
		aPrompt = aUsage.PromptTokenCount
		aCand = aUsage.CandidatesTokenCount
		aThoughts = aUsage.ThoughtsTokenCount
	}
	if sUsage != nil {
		sTotal = int64(sUsage.TotalTokenCount)
		sPrompt = sUsage.PromptTokenCount
		sCand = sUsage.CandidatesTokenCount
		sThoughts = sUsage.ThoughtsTokenCount
	}

	var tokenDeltaPct float64
	if sTotal > 0 {
		tokenDeltaPct = float64(sTotal-aTotal) / float64(sTotal) * 100.0
	}

	fmt.Printf("| Metric                     | Agentic Video        | Static Ingestion     | Observation / Delta       |\n")
	fmt.Printf("|:---------------------------|:---------------------|:---------------------|:--------------------------|\n")
	fmt.Printf("| Total Consumed Tokens      | %-20d | %-20d | %+.1f%% net token spend   |\n", aTotal, sTotal, tokenDeltaPct)
	fmt.Printf("| Prompt Tokens (Input)      | %-20d | %-20d | %-+25s |\n", aPrompt, sPrompt, fmt.Sprintf("%+d input tokens", aPrompt-sPrompt))
	fmt.Printf("| Candidate Output Tokens    | %-20d | %-20d | %-+25s |\n", aCand, sCand, fmt.Sprintf("%+d output tokens", aCand-sCand))
	fmt.Printf("| Reasoning Thoughts Tokens  | %-20d | %-20d | %-+25s |\n", aThoughts, sThoughts, "dynamic frame inspection")
	fmt.Printf("| Latency (Duration)         | %-20v | %-20v | %-+25s |\n",
		aDur.Round(time.Millisecond),
		sDur.Round(time.Millisecond),
		fmt.Sprintf("diff: %v", (aDur - sDur).Round(time.Millisecond)))
	fmt.Printf("=========================================================================================\n")
	fmt.Printf("💡 Telemetry Insight: In Agentic mode, initial prompt tokens are minimal because video frames\n")
	fmt.Printf("   are fetched dynamically during the model's Think ➔ Act timeline loop (accounted under thoughts).\n\n")
}
