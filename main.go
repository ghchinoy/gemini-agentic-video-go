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

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"google.golang.org/genai"
)

// ExampleDefinition holds the metadata and parameters for a video understanding example.
type ExampleDefinition struct {
	ID            int
	Title         string
	Description   string
	VideoURI      string
	Prompt        string
	Mode          genai.MediaProcessing
	ThinkingLevel genai.ThinkingLevel
}

// Built-in examples directly mirroring the scenarios in intro_agentic_video.py.
var builtInExamples = []ExampleDefinition{
	{
		ID:            1,
		Title:         "Agentic Video Understanding (YouTube Long-form)",
		Description:   "Demonstrates dynamic timeline navigation and transcript triage on a long YouTube presentation.",
		VideoURI:      "https://www.youtube.com/watch?v=LzExSq9DU9w",
		Prompt:        "What were the key revenue figures mentioned by the presenter, and at what timestamp do they appear?",
		Mode:          genai.MediaProcessingAgentic,
		ThinkingLevel: genai.ThinkingLevelMedium,
	},
	{
		ID:            2,
		Title:         "Static Video Processing Comparison (YouTube Long-form)",
		Description:   "Runs the exact same query with static 1 FPS ingestion for latency and token consumption comparison.",
		VideoURI:      "https://www.youtube.com/watch?v=LzExSq9DU9w",
		Prompt:        "What were the key revenue figures mentioned by the presenter, and at what timestamp do they appear?",
		Mode:          genai.MediaProcessingStatic,
		ThinkingLevel: genai.ThinkingLevelMedium,
	},
	{
		ID:            3,
		Title:         "Detailed Timestamp-Based Descriptions (GCS Video)",
		Description:   "High-thinking visual analysis on a trail camera video stored in Google Cloud Storage.",
		VideoURI:      "https://storage.googleapis.com/generativeai-downloads/videos/Jukin_Trailcam_Videounderstanding.mp4",
		Prompt:        "Describe what happens in this video in detail, with timestamps.",
		Mode:          genai.MediaProcessingAgentic,
		ThinkingLevel: genai.ThinkingLevelHigh,
	},
	{
		ID:            4,
		Title:         "YouTube Shorts Video Understanding",
		Description:   "Rapid visual and audio explanation of short-form vertical video content.",
		VideoURI:      "https://www.youtube.com/shorts/y-mrGw1wW8E",
		Prompt:        "Explain this video",
		Mode:          genai.MediaProcessingAgentic,
		ThinkingLevel: genai.ThinkingLevelHigh,
	},
}

// Result captures the outcome of a video understanding request.
type Result struct {
	ExampleID int
	Title     string
	Mode      genai.MediaProcessing
	Duration  time.Duration
	Text      string
	Usage     *genai.GenerateContentResponseUsageMetadata
	Error     error
}

func main() {
	var (
		exampleFlag  = flag.String("example", "1", "Example to run: 1, 2, 3, 4, 'all', or 'compare' (runs 1 & 2)")
		videoFlag    = flag.String("video", "", "Custom video URI (e.g. YouTube URL or gs:// / https:// GCS URI)")
		promptFlag   = flag.String("prompt", "", "Custom prompt for video understanding")
		modeFlag     = flag.String("mode", "agentic", "Processing mode for custom video: 'agentic' or 'static'")
		thinkingFlag = flag.String("thinking", "medium", "Thinking level: 'low', 'medium', or 'high'")
		modelFlag    = flag.String("model", "gemini-3.7-flash", "Gemini model ID (e.g. gemini-3.7-flash, gemini-3.6-flash, gemini-3.5-flash-lite)")
		projectFlag  = flag.String("project", "", "Google Cloud project ID (defaults to GOOGLE_CLOUD_PROJECT or active gcloud project)")
		locationFlag = flag.String("location", "", "Google Cloud location/region (defaults to GOOGLE_CLOUD_LOCATION or 'global')")
		backendFlag  = flag.String("backend", "enterprise", "Client backend: 'enterprise', 'vertex', or 'gemini'")
	)
	flag.Parse()

	ctx := context.Background()

	// Resolve project ID
	projectID := resolveProjectID(*projectFlag)

	// Resolve location
	location := *locationFlag
	if location == "" {
		location = os.Getenv("GOOGLE_CLOUD_LOCATION")
		if location == "" {
			location = os.Getenv("GOOGLE_CLOUD_REGION")
		}
		if location == "" {
			location = "global"
		}
	}

	// Initialize GenAI client
	client, err := initClient(ctx, *backendFlag, projectID, location)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing Gemini client: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("=================================================================\n")
	fmt.Printf("  Gemini Agentic Video Understanding (Go)\n")
	fmt.Printf("  Model:    %s\n", *modelFlag)
	fmt.Printf("  Backend:  %s\n", *backendFlag)
	if *backendFlag != "gemini" {
		fmt.Printf("  Project:  %s\n", projectID)
		fmt.Printf("  Location: %s\n", location)
	}
	fmt.Printf("=================================================================\n\n")

	// Custom query branch
	if *videoFlag != "" {
		runCustom(ctx, client, *modelFlag, *videoFlag, *promptFlag, *modeFlag, *thinkingFlag)
		return
	}

	// Built-in examples branch
	switch strings.ToLower(*exampleFlag) {
	case "1":
		runExample(ctx, client, *modelFlag, builtInExamples[0])
	case "2":
		runExample(ctx, client, *modelFlag, builtInExamples[1])
	case "3":
		runExample(ctx, client, *modelFlag, builtInExamples[2])
	case "4":
		runExample(ctx, client, *modelFlag, builtInExamples[3])
	case "compare":
		fmt.Println("Running comparative benchmark: Agentic vs. Static on long-form YouTube video...")
		r1 := runExample(ctx, client, *modelFlag, builtInExamples[0])
		r2 := runExample(ctx, client, *modelFlag, builtInExamples[1])
		printComparison(r1, r2)
	case "all":
		fmt.Println("Running all 4 built-in examples...")
		var results []*Result
		for _, ex := range builtInExamples {
			res := runExample(ctx, client, *modelFlag, ex)
			results = append(results, res)
		}
		if len(results) >= 2 {
			printComparison(results[0], results[1])
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown example option: %q. Valid options: 1, 2, 3, 4, compare, all\n", *exampleFlag)
		os.Exit(1)
	}
}

// initClient configures the Google GenAI client based on requested backend and credentials.
func initClient(ctx context.Context, backendStr, projectID, location string) (*genai.Client, error) {
	cfg := &genai.ClientConfig{}

	switch strings.ToLower(backendStr) {
	case "enterprise":
		if projectID == "" {
			return nil, fmt.Errorf("project ID is required for Enterprise Agent Platform. Set GOOGLE_CLOUD_PROJECT or use -project")
		}
		cfg.Backend = genai.BackendEnterprise
		cfg.Project = projectID
		cfg.Location = location
	case "vertex":
		if projectID == "" {
			return nil, fmt.Errorf("project ID is required for Vertex AI. Set GOOGLE_CLOUD_PROJECT or use -project")
		}
		cfg.Backend = genai.BackendVertexAI
		cfg.Project = projectID
		cfg.Location = location
	case "gemini":
		cfg.Backend = genai.BackendGeminiAPI
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			apiKey = os.Getenv("GOOGLE_API_KEY")
		}
		if apiKey != "" {
			cfg.APIKey = apiKey
		}
	default:
		return nil, fmt.Errorf("unsupported backend %q; use 'enterprise', 'vertex', or 'gemini'", backendStr)
	}

	return genai.NewClient(ctx, cfg)
}

// resolveProjectID determines the GCP project ID from flag, env, or gcloud CLI.
func resolveProjectID(flagValue string) string {
	if flagValue != "" && flagValue != "[your-project-id]" {
		return flagValue
	}
	if envVal := os.Getenv("GOOGLE_CLOUD_PROJECT"); envVal != "" && envVal != "[your-project-id]" {
		return envVal
	}
	// Fallback to active gcloud config
	out, err := exec.Command("gcloud", "config", "get-value", "project").Output()
	if err == nil {
		p := strings.TrimSpace(string(out))
		if p != "" && p != "(unset)" {
			return p
		}
	}
	return ""
}

// parseThinkingLevel maps string representation to genai.ThinkingLevel.
func parseThinkingLevel(level string) genai.ThinkingLevel {
	switch strings.ToLower(level) {
	case "low":
		return genai.ThinkingLevelLow
	case "high":
		return genai.ThinkingLevelHigh
	case "minimal":
		return genai.ThinkingLevelMinimal
	default:
		return genai.ThinkingLevelMedium
	}
}

// parseProcessingMode maps string representation to genai.MediaProcessing.
func parseProcessingMode(mode string) genai.MediaProcessing {
	switch strings.ToLower(mode) {
	case "static":
		return genai.MediaProcessingStatic
	default:
		return genai.MediaProcessingAgentic
	}
}

// runExample executes a built-in example scenario and displays results.
func runExample(ctx context.Context, client *genai.Client, modelID string, ex ExampleDefinition) *Result {
	fmt.Printf("-----------------------------------------------------------------\n")
	fmt.Printf("Example %d: %s\n", ex.ID, ex.Title)
	fmt.Printf("Description: %s\n", ex.Description)
	fmt.Printf("Video URI:   %s\n", ex.VideoURI)
	fmt.Printf("Mode:        %s\n", ex.Mode)
	fmt.Printf("Thinking:    %s\n", ex.ThinkingLevel)
	fmt.Printf("Prompt:      %q\n", ex.Prompt)
	fmt.Printf("-----------------------------------------------------------------\n")

	res := executeRequest(ctx, client, modelID, ex.VideoURI, ex.Prompt, ex.Mode, ex.ThinkingLevel)
	res.ExampleID = ex.ID
	res.Title = ex.Title

	printResult(res)
	return res
}

// runCustom executes a custom user query against the model.
func runCustom(ctx context.Context, client *genai.Client, modelID, videoURI, prompt, modeStr, thinkingStr string) {
	mode := parseProcessingMode(modeStr)
	thinking := parseThinkingLevel(thinkingStr)

	if prompt == "" {
		prompt = "Describe this video in detail."
	}

	fmt.Printf("-----------------------------------------------------------------\n")
	fmt.Printf("Custom Video Query\n")
	fmt.Printf("Video URI:   %s\n", videoURI)
	fmt.Printf("Mode:        %s\n", mode)
	fmt.Printf("Thinking:    %s\n", thinking)
	fmt.Printf("Prompt:      %q\n", prompt)
	fmt.Printf("-----------------------------------------------------------------\n")

	res := executeRequest(ctx, client, modelID, videoURI, prompt, mode, thinking)
	res.Title = "Custom Video Understanding"
	printResult(res)
}

// executeRequest constructs the Part with MediaProcessing configuration and invokes GenerateContent.
func executeRequest(ctx context.Context, client *genai.Client, modelID, videoURI, prompt string, mode genai.MediaProcessing, thinking genai.ThinkingLevel) *Result {
	// Construct the video Part with MediaProcessing mode (Agentic or Static)
	videoPart := genai.NewPartFromURI(videoURI, "video/mp4")
	videoPart.MediaProcessing = mode

	// Construct message content containing video and text parts
	contents := []*genai.Content{
		genai.NewContentFromParts([]*genai.Part{
			videoPart,
			genai.NewPartFromText(prompt),
		}, genai.RoleUser),
	}

	// Configure generation with reasoning thinking level
	genConfig := &genai.GenerateContentConfig{
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingLevel: thinking,
		},
	}

	start := time.Now()
	resp, err := client.Models.GenerateContent(ctx, modelID, contents, genConfig)
	duration := time.Since(start)

	if err != nil {
		return &Result{
			Mode:     mode,
			Duration: duration,
			Error:    err,
		}
	}

	return &Result{
		Mode:     mode,
		Duration: duration,
		Text:     resp.Text(),
		Usage:    resp.UsageMetadata,
	}
}

// printResult outputs the generated response, token metrics, and elapsed time.
func printResult(res *Result) {
	if res.Error != nil {
		fmt.Fprintf(os.Stderr, "\n[ERROR] Request failed after %v: %v\n\n", res.Duration.Round(time.Millisecond), res.Error)
		return
	}

	fmt.Println("\n--- Model Response ---")
	fmt.Println(res.Text)
	fmt.Println("----------------------")

	if res.Usage != nil {
		fmt.Printf("Total Token Count:      %d\n", res.Usage.TotalTokenCount)
		fmt.Printf("  • Prompt Tokens:      %d\n", res.Usage.PromptTokenCount)
		fmt.Printf("  • Candidates Tokens:  %d\n", res.Usage.CandidatesTokenCount)
		if res.Usage.ThoughtsTokenCount > 0 {
			fmt.Printf("  • Thoughts Tokens:    %d\n", res.Usage.ThoughtsTokenCount)
		}
	}
	fmt.Printf("Processing Duration:    %v\n\n", res.Duration.Round(time.Millisecond))
}

// printComparison displays a summary table comparing Agentic vs. Static performance.
func printComparison(agenticRes, staticRes *Result) {
	fmt.Printf("=================================================================\n")
	fmt.Printf("  PERFORMANCE COMPARISON: AGENTIC vs. STATIC PROCESSING\n")
	fmt.Printf("=================================================================\n")

	if agenticRes.Error != nil || staticRes.Error != nil {
		fmt.Println("Cannot generate comparison: one or both runs encountered an error.")
		return
	}

	aTokens := int64(0)
	sTokens := int64(0)
	if agenticRes.Usage != nil {
		aTokens = int64(agenticRes.Usage.TotalTokenCount)
	}
	if staticRes.Usage != nil {
		sTokens = int64(staticRes.Usage.TotalTokenCount)
	}

	var tokenSavingsPct float64
	if sTokens > 0 {
		tokenSavingsPct = float64(sTokens-aTokens) / float64(sTokens) * 100.0
	}

	fmt.Printf("| Metric                | Agentic Video         | Static Ingestion      | Delta / Benefit        |\n")
	fmt.Printf("|-----------------------|-----------------------|-----------------------|------------------------|\n")
	fmt.Printf("| Total Tokens          | %-21d | %-21d | %+.1f%% reduction     |\n", aTokens, sTokens, tokenSavingsPct)
	if agenticRes.Usage != nil && staticRes.Usage != nil {
		fmt.Printf("| Prompt Tokens         | %-21d | %-21d | %-22s |\n",
			agenticRes.Usage.PromptTokenCount, staticRes.Usage.PromptTokenCount,
			fmt.Sprintf("%+d tokens", agenticRes.Usage.PromptTokenCount-staticRes.Usage.PromptTokenCount))
	}
	fmt.Printf("| Duration (Latency)    | %-21v | %-21v | %-22s |\n",
		agenticRes.Duration.Round(time.Millisecond),
		staticRes.Duration.Round(time.Millisecond),
		fmt.Sprintf("diff: %v", (agenticRes.Duration-staticRes.Duration).Round(time.Millisecond)))
	fmt.Printf("=================================================================\n\n")
}
