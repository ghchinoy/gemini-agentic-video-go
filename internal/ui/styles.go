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

// Package ui provides terminal styling, banners, badges, and table renderers for the CLI.
// NOTE: None of these presentation styles are required by the Gemini API.
// For core Agentic Video SDK usage, see internal/runner and internal/telemetry.
package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"google.golang.org/genai"
)

var (
	// Google & Gemini Branding Palette
	ColorGeminiBlue    = lipgloss.Color("#4285F4") // Google Blue
	ColorGeminiCyan    = lipgloss.Color("#00C4FF") // Gemini Sparkle Cyan
	ColorAgenticPurple = lipgloss.Color("#8E44AD") // Agentic Reasoning Purple
	ColorSuccess       = lipgloss.Color("#34A853") // Google Green
	ColorWarning       = lipgloss.Color("#FBBC05") // Google Amber
	ColorDanger        = lipgloss.Color("#EA4335") // Google Red
	ColorMuted         = lipgloss.Color("#70757A") // Google Medium Gray
	ColorDim           = lipgloss.Color("#3C4043") // Google Dark Slate
	ColorWhite         = lipgloss.Color("#FFFFFF")
	ColorZebraBg       = lipgloss.Color("#1A1D24") // Subtle dark row background

	// Typography & Container Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWhite).
			Background(ColorGeminiBlue).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorGeminiCyan).
			Bold(true)

	HeaderBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorGeminiBlue).
			Padding(0, 1).
			MarginBottom(1)

	CardBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorMuted).
		Padding(0, 1).
		MarginBottom(1)

	SectionHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorGeminiCyan).
			MarginTop(1).
			MarginBottom(0)

	MutedStyle = lipgloss.NewStyle().Foreground(ColorMuted)
	BoldWhite  = lipgloss.NewStyle().Bold(true).Foreground(ColorWhite)
	GreenStyle = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	CyanStyle  = lipgloss.NewStyle().Foreground(ColorGeminiCyan).Bold(true)

	// Mode Badges
	BadgeAgentic    = lipgloss.NewStyle().Bold(true).Foreground(ColorWhite).Background(ColorAgenticPurple).Padding(0, 1)
	BadgeStatic     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000000")).Background(ColorWarning).Padding(0, 1)
	BadgeMultiVideo = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000000")).Background(ColorGeminiCyan).Padding(0, 1)
	BadgeMultiTurn  = lipgloss.NewStyle().Bold(true).Foreground(ColorWhite).Background(ColorSuccess).Padding(0, 1)

	// Thinking Badges
	BadgeThinkingHigh    = lipgloss.NewStyle().Bold(true).Foreground(ColorWhite).Background(lipgloss.Color("#B02A37")).Padding(0, 1)
	BadgeThinkingMedium  = lipgloss.NewStyle().Bold(true).Foreground(ColorWhite).Background(lipgloss.Color("#0D6EFD")).Padding(0, 1)
	BadgeThinkingLow     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#6C757D")).Padding(0, 1)
	BadgeThinkingMinimal = lipgloss.NewStyle().Bold(true).Foreground(ColorWhite).Background(ColorDim).Padding(0, 1)
)

// RenderModeBadge returns a styled Lipgloss badge for a media processing mode.
func RenderModeBadge(mode genai.MediaProcessing) string {
	if mode == genai.MediaProcessingStatic {
		return BadgeStatic.Render("STATIC (1 FPS)")
	}
	return BadgeAgentic.Render("AGENTIC")
}

// RenderScenarioTypeBadge returns a styled badge describing the scenario category.
func RenderScenarioTypeBadge(name string, mode genai.MediaProcessing, isMultiVideo, isMultiTurn bool) string {
	if isMultiVideo {
		return BadgeMultiVideo.Render("MULTI-VIDEO")
	}
	if isMultiTurn {
		return BadgeMultiTurn.Render("MULTI-TURN")
	}
	if mode == genai.MediaProcessingStatic {
		return BadgeStatic.Render("STATIC")
	}
	return BadgeAgentic.Render("AGENTIC")
}

// RenderThinkingBadge returns a styled badge for reasoning thinking level.
func RenderThinkingBadge(level genai.ThinkingLevel) string {
	switch level {
	case genai.ThinkingLevelHigh:
		return BadgeThinkingHigh.Render("THINKING: HIGH")
	case genai.ThinkingLevelLow:
		return BadgeThinkingLow.Render("THINKING: LOW")
	case genai.ThinkingLevelMinimal:
		return BadgeThinkingMinimal.Render("THINKING: MINIMAL")
	default:
		return BadgeThinkingMedium.Render("THINKING: MEDIUM")
	}
}

// RenderTurnBadge formats a conversational turn marker.
func RenderTurnBadge(turnNum, totalTurns int) string {
	label := "Turn " + string(rune('0'+turnNum))
	if totalTurns > 0 {
		label = label + "/" + string(rune('0'+totalTurns))
	}
	return lipgloss.NewStyle().Bold(true).Foreground(ColorWhite).Background(ColorAgenticPurple).Padding(0, 1).Render(label)
}

// RenderVideoBadge formats a multi-video index badge.
func RenderVideoBadge(index int, label string) string {
	text := "Video " + string(rune('0'+index))
	if label != "" {
		text = text + " (" + strings.ToUpper(label) + ")"
	}
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000000")).Background(ColorGeminiCyan).Padding(0, 1).Render(text)
}
