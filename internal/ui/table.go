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
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/ghchinoy/gemini-agentic-video-go/internal/catalog"
	"google.golang.org/genai"
)

// RenderCatalogTable formats the preset scenarios into a styled Lipgloss table.
func RenderCatalogTable(scenarios []catalog.Scenario) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorWhite).
		Background(ColorGeminiBlue).
		Align(lipgloss.Center).
		Padding(0, 1)

	baseCellStyle := lipgloss.NewStyle().Padding(0, 1)
	zebraStyle := baseCellStyle.Background(ColorZebraBg)

	tbl := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ColorGeminiBlue)).
		Headers("ID", "Mode", "Identifier", "Title & Capabilities").
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			style := baseCellStyle
			if row%2 == 0 {
				style = zebraStyle
			}
			switch col {
			case 0:
				return style.Bold(true).Foreground(ColorWhite).Align(lipgloss.Center)
			case 2:
				return style.Foreground(ColorGeminiCyan).Bold(true)
			default:
				return style.Foreground(ColorWhite)
			}
		})

	for _, sc := range scenarios {
		isMultiVideo := len(sc.MultiVideos) > 0
		isMultiTurn := len(sc.MultiTurns) > 0
		typeBadge := RenderScenarioTypeBadge(sc.Name, sc.Mode, isMultiVideo, isMultiTurn)

		tbl.Row(
			strconv.Itoa(sc.ID),
			typeBadge,
			sc.Name,
			sc.Title,
		)
	}

	header := SectionHeader.Render("📋 Built-In Tutorial Scenarios (Run with: ./bin/gemini-agentic-video-go example <id>)")
	return lipgloss.JoinVertical(lipgloss.Left, header, tbl.Render())
}

// RenderBenchmarkTable renders the side-by-side Agentic vs Static benchmark comparison.
func RenderBenchmarkTable(agenticUsage, staticUsage *genai.GenerateContentResponseUsageMetadata, agenticDuration, staticDuration time.Duration) string {
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

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorWhite).
		Background(ColorGeminiBlue).
		Align(lipgloss.Center).
		Padding(0, 1)

	baseCellStyle := lipgloss.NewStyle().Padding(0, 1)
	zebraStyle := baseCellStyle.Background(ColorZebraBg)

	tbl := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ColorGeminiBlue)).
		Headers("Performance Metric", "Agentic Video", "Static (1 FPS)", "Delta / Observation").
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			style := baseCellStyle
			if row%2 == 0 {
				style = zebraStyle
			}
			switch col {
			case 0:
				return style.Bold(true).Foreground(ColorWhite)
			case 1:
				return style.Foreground(ColorGeminiCyan).Align(lipgloss.Right)
			case 2:
				return style.Foreground(ColorWarning).Align(lipgloss.Right)
			case 3:
				if row == 0 {
					return style.Foreground(ColorSuccess).Bold(true)
				}
				return style.Foreground(ColorWhite)
			default:
				return style
			}
		})

	totalDeltaStr := fmt.Sprintf("%+.1f%% net spend", tokenDeltaPct)
	if tokenDeltaPct > 0 {
		totalDeltaStr = fmt.Sprintf("🟢 %+.1f%% reduction", tokenDeltaPct)
	}

	promptDeltaStr := fmt.Sprintf("%+d input tokens", agenticPrompt-staticPrompt)
	candDeltaStr := fmt.Sprintf("%+d output tokens", agenticCandidates-staticCandidates)
	latDeltaStr := fmt.Sprintf("%v diff", (agenticDuration - staticDuration).Round(time.Millisecond))

	tbl.Row(
		"Total Consumed Tokens",
		fmt.Sprintf("%d", agenticTotal),
		fmt.Sprintf("%d", staticTotal),
		totalDeltaStr,
	)
	tbl.Row(
		"Prompt Tokens (Input)",
		fmt.Sprintf("%d (text only)", agenticPrompt),
		fmt.Sprintf("%d (all frames)", staticPrompt),
		promptDeltaStr,
	)
	tbl.Row(
		"Candidate Tokens (Output)",
		fmt.Sprintf("%d", agenticCandidates),
		fmt.Sprintf("%d", staticCandidates),
		candDeltaStr,
	)
	tbl.Row(
		"Reasoning Thoughts Tokens",
		fmt.Sprintf("%d (dynamic frames)", agenticThoughts),
		fmt.Sprintf("%d", staticThoughts),
		"on-demand inspection",
	)
	tbl.Row(
		"Latency (Wall Clock)",
		fmt.Sprintf("%v", agenticDuration.Round(time.Millisecond)),
		fmt.Sprintf("%v", staticDuration.Round(time.Millisecond)),
		latDeltaStr,
	)

	title := SectionHeader.Render("⚡ PERFORMANCE BENCHMARK: AGENTIC VIDEO vs. STATIC (1 FPS) INGESTION")
	footer := MutedStyle.Render("💡 Insight: In Agentic mode, prompt tokens contain 0 video frames upfront. Frames are dynamically fetched\n   during Think ➔ Act timeline navigation and billed under reasoning thoughts.")

	return lipgloss.JoinVertical(lipgloss.Left, title, tbl.Render(), footer)
}
