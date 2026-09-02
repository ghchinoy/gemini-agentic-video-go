# gemini-agentic-video-go

A Go application demonstrating Gemini agentic video understanding with active Think ➔ Act ➔ Observe dynamic timeline navigation, sub-second precision, and up to 96% token reduction compared to passive frame ingestion.

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Google GenAI SDK](https://img.shields.io/badge/Google%20GenAI%20SDK-v1.70.0-4285F4?logo=google)](https://pkg.go.dev/google.golang.org/genai)

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [Common Examples](#common-examples)
  - [Programmatic Go SDK Usage](#programmatic-go-sdk-usage)
- [Development](#development)
  - [Prerequisites](#prerequisites)
  - [Setup and Running](#setup-and-running)
  - [Makefile Targets](#makefile-targets)
  - [CLI Flags](#cli-flags)
- [Publishing](#publishing)
- [Contributing](#contributing)
- [License](#license)
- [Overview: Agentic Video Understanding](#overview-agentic-video-understanding)
- [Processing Modes & When to Use](#processing-modes--when-to-use)
- [Key Use Cases](#key-use-cases)
- [Migration Guide & Thinking Levels](#migration-guide--thinking-levels)
- [Built-In Scenarios](#built-in-scenarios)
- [Credits & Acknowledgements](#credits--acknowledgements)

## Installation

Install the binary directly using Go:

```bash
go install github.com/ghchinoy/gemini-agentic-video-go@latest
```

Or clone and build from source:

```bash
git clone https://github.com/ghchinoy/gemini-agentic-video-go.git
cd gemini-agentic-video-go
make build
```

The binary will be compiled to `./bin/gemini-agentic-video-go`.

## Usage

Run the agentic video query on a long-form YouTube presentation:

```bash
./bin/gemini-agentic-video-go -example=1
```

Or run via `make`:

```bash
make run
```

This analyzes an Alphabet earnings call video on YouTube using agentic timeline navigation and transcript triage, extracting revenue figures and timestamps in seconds without decoding all video frames upfront.

### Common Examples

Run a side-by-side performance comparison between Agentic and Static processing on the same video:

```bash
make compare
# or: ./bin/gemini-agentic-video-go -example=compare
```

Detailed visual timestamped analysis of a wildlife trail-cam video from Google Cloud Storage:

```bash
./bin/gemini-agentic-video-go -example=3
```

Rapid visual explanation of a YouTube Short:

```bash
./bin/gemini-agentic-video-go -example=4
```

Analyze any custom video (YouTube URL or Google Cloud Storage URI):

```bash
./bin/gemini-agentic-video-go \
  -video="https://www.youtube.com/watch?v=LzExSq9DU9w" \
  -prompt="What product announcements were made?" \
  -mode=agentic \
  -thinking=high
```

### Programmatic Go SDK Usage

To enable agentic video processing in your own Go applications using `google.golang.org/genai`:

```go
package main

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend:  genai.BackendEnterprise, // or genai.BackendVertexAI
		Project:  os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Location: "global",
	})
	if err != nil {
		panic(err)
	}

	// 1. Create a video Part and set MediaProcessing to Agentic
	videoPart := genai.NewPartFromURI("https://www.youtube.com/watch?v=LzExSq9DU9w", "video/mp4")
	videoPart.MediaProcessing = genai.MediaProcessingAgentic // Options: MediaProcessingAgentic, MediaProcessingStatic

	// 2. Wrap the video and prompt in content parts
	contents := []*genai.Content{
		genai.NewContentFromParts([]*genai.Part{
			videoPart,
			genai.NewPartFromText("What are the key revenue figures and timestamps?"),
		}, genai.RoleUser),
	}

	// 3. Configure reasoning thinking level (Low, Medium, High)
	config := &genai.GenerateContentConfig{
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingLevel: genai.ThinkingLevelMedium,
		},
	}

	// 4. Generate content
	resp, err := client.Models.GenerateContent(ctx, "gemini-3.7-flash", contents, config)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Text())
	fmt.Printf("Total Tokens: %d\n", resp.UsageMetadata.TotalTokenCount)
}
```

## Development

### Prerequisites

- [Go 1.24+](https://golang.org/dl/)
- A Google Cloud Project with the Vertex AI / Gemini API enabled
- Authenticated Google Cloud credentials via Application Default Credentials (`gcloud auth application-default login`)

### Setup and Running

```bash
git clone https://github.com/ghchinoy/gemini-agentic-video-go.git
cd gemini-agentic-video-go
go mod download
export GOOGLE_CLOUD_PROJECT="your-project-id"
export GOOGLE_CLOUD_LOCATION="global"
go run main.go -example=1
```

### Makefile Targets

| Target | Description |
| :--- | :--- |
| `make build` | Compiles binary to `./bin/gemini-agentic-video-go` |
| `make run` | Builds and runs Example 1 (use `ARGS="..."` for custom flags) |
| `make compare` | Runs the Agentic vs. Static performance comparison benchmark |
| `make clean` | Removes `./bin` directory and build artifacts |
| `make fmt` | Formats Go source code |
| `make vet` | Runs `go vet ./...` static analyzer |
| `make license` | Applies Apache 2.0 license headers using `addlicense` |
| `make license-check` | Verifies presence of license headers |
| `make help` | Displays list of available make targets |

### CLI Flags

| Flag | Default | Description |
| :--- | :--- | :--- |
| `-example` | `1` | Built-in scenario to run: `1`, `2`, `3`, `4`, `compare`, or `all` |
| `-video` | `""` | Custom video URI (YouTube URL, `gs://`, or `https://` Cloud Storage) |
| `-prompt` | `""` | Custom text prompt for video understanding |
| `-mode` | `agentic` | Processing mode: `agentic` or `static` |
| `-thinking` | `medium` | Reasoning thinking level: `low`, `medium`, or `high` |
| `-model` | `gemini-3.7-flash` | Gemini model ID (`gemini-3.7-flash`, `gemini-3.6-flash`, `gemini-3.5-flash-lite`) |
| `-project` | `""` | Google Cloud project ID (falls back to `GOOGLE_CLOUD_PROJECT` or `gcloud config`) |
| `-location` | `"global"` | Google Cloud location / region |
| `-backend` | `enterprise` | Backend: `enterprise`, `vertex`, or `gemini` (API key) |

## Publishing

To publish a new tagged release:

```bash
git tag v0.1.0
git push origin v0.1.0
```

To cross-compile binaries for multiple platforms:

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o bin/gemini-agentic-video-go-linux-amd64 .

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o bin/gemini-agentic-video-go-darwin-arm64 .

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o bin/gemini-agentic-video-go-windows-amd64.exe .
```

## Contributing

Pull requests are welcome! For major architectural changes or new feature proposals, please open an issue first to discuss what you would like to change.

Please make sure all Go code is properly formatted (`go fmt ./...`) and passes compilation before submitting a pull request.

## License

This project is licensed under the [Apache 2.0 License](LICENSE).

---

## Overview: Agentic Video Understanding

Traditional multimodal LLMs ingest video passively by extracting static frames at a fixed rate (1 FPS) alongside 1 Kbps audio. For a 60-minute video, static ingestion consumes over 400,000 tokens—leading to high costs, multi-minute latency, and context bloat even for simple verbal questions.

**Agentic Video** replaces passive frame extraction with an active **Think ➔ Act ➔ Observe** loop. Using a native timeline exploration tool, the model dynamically inspects only the required streams (captions/ASR, audio waveform, or visual frames), temporal windows, and sampling resolutions.

- **78%–96% Token Reduction**: Skips irrelevant visual frames (e.g., reducing a 1-hour video from ~407k to ~23k tokens).
- **Seconds Time-to-First-Byte (TTFB)**: Delivers rapid responses for transcript- and audio-answerable queries without decoding full video.
- **Sub-Second Precision (+1.1% to +4.6% QA, +17% Fast Action)**: Replaces uniform 1-FPS sampling with targeted high-FPS inspection (`fps=10+`).
- **Multi-Video Scalability**: Enables comparative analysis across multiple long-form videos within the 1M token context limit.

## Processing Modes & When to Use

| Mode | Description | Supported Models | Token Efficiency & Latency | When to Use |
| :--- | :--- | :--- | :--- | :--- |
| **Static**<br>*(default)* | Single-pass ingestion at fixed **1 FPS** (~300 tokens/sec) and **1 Kbps audio** with 1s timestamps. | All Gemini models | • **Baseline tokens:** ~300 tokens/sec<br>• **Latency:** Full ingestion required upfront | • Short clips (< 2 min)<br>• Requires simultaneous full audio & video<br>• Fixed-interval deterministic sampling<br>• Custom clipping (`start_offset`, `end_offset`) |
| **Agentic** | Iterative timeline navigation; dynamically loads frames and audio on-demand via reasoning budget. | • `gemini-3.7-flash`<br>• `gemini-3.6-flash`<br>• `gemini-3.5-flash-lite` | • **Tokens:** 70%–95% reduction on long content<br>• **Latency:** Fast TTFB (seconds) for transcript queries<br>• **Quality:** Higher overall answer accuracy | • Long videos > 2 min (lectures, meetings)<br>• Verbal / audio-first queries<br>• Localized visual search<br>• Fast-action clips needing high temporal precision<br>• Cross-video comparisons |

## Key Use Cases

| Use Case | Target Workloads | Why Static Fails | The Agentic Advantage |
| :--- | :--- | :--- | :--- |
| **Long-Form Q&A & Meetings** | • 30–90+ min meetings<br>• Lectures<br>• Earnings calls | • Decodes all 3,600+ frames (400k+ tokens) even for audio-only answers | • **78%–96% token reduction**<br>• **Seconds TTFB:** Uses transcript triage to bypass visual decoding |
| **Visual "Needle-in-a-Haystack"** | • Slide transitions<br>• Diagram edits<br>• Specific visual events | • Fixed 1 FPS dilutes attention<br>• Requires heavy external RAG pipelines | • **Coarse-to-fine zoom:** Probes at 0.5 FPS, then zooms into target 5s windows at 8–10 FPS<br>• **+3.5% to +4.6%** accuracy gain |
| **Fast-Action & Anomalies** | • Sports highlights<br>• Industrial QA<br>• UI glitches & telemetry | • 1 FPS misses sub-second actions (<1s)<br>• 10+ FPS full-video ingestion exceeds token limits | • **Adaptive High-FPS Replay:** Dynamically samples `fps=10+` only on critical 3–10s intervals<br>• **+17% accuracy** with 30% fewer tokens |
| **Cross-Video Synthesis** | • Multi-angle feeds<br>• Product reviews<br>• Lecture vs. lab footage | • Ingesting multiple 30+ min videos quickly exceeds the 1M token context limit | • **Multi-video feasibility:** Dynamically retrieves only relevant segments across files in a single prompt |

## Migration Guide & Thinking Levels

1. **Remove Preprocessing Workarounds**: Eliminate custom ffmpeg downsampling scripts (e.g. downsampling to 0.1 FPS). Pass the full video URI directly with `part.MediaProcessing = genai.MediaProcessingAgentic`.
2. **Configure `thinking_level`**:
   - `HIGH`: Dense visual QA, split-second action/sports analysis, or complex multi-step reasoning across 60+ minute videos.
   - `MEDIUM` (Default): General video Q&A, lecture summarization, and clip retrieval.
   - `LOW`: Fast transcript/caption searches and metadata extraction.

## Built-In Scenarios

The application includes 4 built-in scenarios directly adapted from the reference Python notebook (`sources/intro_agentic_video.py`):

1. **Example 1 (`-example=1`) — Agentic Video Understanding (YouTube Long-form)**:
   Extracts key revenue figures and exact timestamps from an Alphabet earnings call presentation (`https://www.youtube.com/watch?v=LzExSq9DU9w`) with `MEDIUM` thinking.
2. **Example 2 (`-example=2`) — Static Video Processing Comparison**:
   Runs the exact same query with static 1 FPS video frame ingestion for direct comparison of token counts and latency.
3. **Example 3 (`-example=3`) — Detailed Timestamp-Based Descriptions (GCS Video)**:
   Deep visual analysis with `HIGH` thinking on a wildlife trail-cam video (`https://storage.googleapis.com/generativeai-downloads/videos/Jukin_Trailcam_Videounderstanding.mp4`).
4. **Example 4 (`-example=4`) — YouTube Shorts Video Understanding**:
   Rapid visual and audio analysis of a YouTube Short (`https://www.youtube.com/shorts/y-mrGw1wW8E`) with `HIGH` thinking.

## Credits & Acknowledgements

- Based on the [Google Cloud Generative AI Agentic Video Tutorial](https://github.com/GoogleCloudPlatform/generative-ai/blob/main/gemini/agentic-video/intro_agentic_video.ipynb) by [Eric Dong](https://github.com/gericdong) and [Holt Skinner](https://github.com/holtskinner).
- README structure designed according to Mark Allen's ["How to Write a Great README for Your Public GitHub Project"](https://www.markcallen.com/how-to-write-a-great-readme-for-your-public-github-project/) (Everyday DevOps).
