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
	"fmt"
	"time"

	"github.com/ghchinoy/gemini-agentic-video-go/internal/runner"
	"github.com/ghchinoy/gemini-agentic-video-go/internal/ui"
	"github.com/spf13/cobra"
	"google.golang.org/genai"
)

var (
	compareVideoURI   string
	comparePrompt     string
	compareThinking   string
	compareConcurrent bool
)

// compareCmd represents the compare command
var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Benchmark Agentic vs. Static processing on the same video",
	Long: `Executes the identical query using Agentic video navigation and Static 1 FPS
frame ingestion, displaying an automated side-by-side performance comparison
with net token spend reduction and latency metrics.

By default, runs both pipelines concurrently in parallel goroutines with a live
side-by-side progress, elapsed time, and token counter.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := InitClient(ctx)
		if err != nil {
			return err
		}

		thinking := runner.ParseThinkingLevel(compareThinking)

		if compareConcurrent {
			tracker := ui.NewBenchmarkTracker(compareThinking)
			tracker.Start()

			benchRes := runner.ExecuteConcurrentBenchmark(ctx, client, runner.Request{
				ModelID:       modelFlag,
				VideoURI:      compareVideoURI,
				Prompt:        comparePrompt,
				ThinkingLevel: thinking,
			}, func(agenticDone, staticDone bool, aRes, sRes *runner.Result) {
				tracker.Update(agenticDone, staticDone, aRes, sRes)
			})

			tracker.SetErrors(benchRes.AgenticError, benchRes.StaticError)
			tracker.Stop()

			if benchRes.AgenticError != nil && benchRes.StaticError != nil {
				return fmt.Errorf("both benchmarks failed: agentic: %v; static: %v", benchRes.AgenticError, benchRes.StaticError)
			}
			if benchRes.AgenticError != nil {
				return fmt.Errorf("agentic benchmark failed: %w", benchRes.AgenticError)
			}
			if benchRes.StaticError != nil {
				return fmt.Errorf("static benchmark failed: %w", benchRes.StaticError)
			}

			aUsage := benchRes.AgenticResult.Usage
			sUsage := benchRes.StaticResult.Usage
			fmt.Println(ui.RenderBenchmarkTable(aUsage, sUsage, benchRes.AgenticResult.Duration, benchRes.StaticResult.Duration))
			return nil
		}

		// Sequential fallback
		fmt.Println(">>> Step 1/2: Executing with AGENTIC video processing...")
		aRes, err := runner.Execute(ctx, client, runner.Request{
			ModelID:       modelFlag,
			VideoURI:      compareVideoURI,
			Prompt:        comparePrompt,
			Mode:          genai.MediaProcessingAgentic,
			ThinkingLevel: thinking,
		})
		if err != nil {
			return fmt.Errorf("agentic run failed: %w", err)
		}
		fmt.Printf("    Agentic Run Complete (%v, %d tokens)\n\n",
			aRes.Duration.Round(time.Millisecond),
			func() int32 {
				if aRes.Usage != nil {
					return aRes.Usage.TotalTokenCount
				}
				return 0
			}())

		fmt.Println(">>> Step 2/2: Executing with STATIC 1-FPS frame ingestion...")
		sRes, err := runner.Execute(ctx, client, runner.Request{
			ModelID:       modelFlag,
			VideoURI:      compareVideoURI,
			Prompt:        comparePrompt,
			Mode:          genai.MediaProcessingStatic,
			ThinkingLevel: thinking,
		})
		if err != nil {
			return fmt.Errorf("static run failed: %w", err)
		}
		fmt.Printf("    Static Run Complete (%v, %d tokens)\n\n",
			sRes.Duration,
			func() int32 {
				if sRes.Usage != nil {
					return sRes.Usage.TotalTokenCount
				}
				return 0
			}())

		fmt.Println(ui.RenderBenchmarkTable(aRes.Usage, sRes.Usage, aRes.Duration, sRes.Duration))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(compareCmd)

	compareCmd.Flags().StringVarP(&compareVideoURI, "video", "v", "https://www.youtube.com/watch?v=LzExSq9DU9w", "Video URI for comparison benchmark")
	compareCmd.Flags().StringVar(&comparePrompt, "prompt", "What were the key revenue figures mentioned by the presenter, and at what timestamp do they appear?", "Prompt for comparison benchmark")
	compareCmd.Flags().StringVar(&compareThinking, "thinking", "medium", "Thinking level: 'low', 'medium', or 'high'")
	compareCmd.Flags().BoolVar(&compareConcurrent, "concurrent", true, "Execute Agentic and Static benchmarks concurrently in parallel goroutines with live telemetry")
}
