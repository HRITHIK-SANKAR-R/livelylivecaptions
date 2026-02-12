---
phase: 2
plan: 1
wave: 1
---

# Plan 2.1: UI Structure & Data Models

## Objective
Establish the foundational data structures and UI skeleton for the Split View TUI.

## Context
- .gsd/SPEC.md
- .gsd/phases/2/RESEARCH.md
- internal/ui/ui.go

## Tasks

<task type="auto">
  <name>Refactor Transcriber Output</name>
  <files>internal/transcriber/transcriber.go</files>
  <action>
    - Define `TranscriptionEvent` struct in a shared package (e.g., `internal/types`).
    - Update `Transcriber` to output `TranscriptionEvent` instead of string.
    - Update `Start` loop to populate `IsFinal` based on endpoint detection.
  </action>
  <verify>go build ./internal/transcriber</verify>
  <done>Transcriber outputs structured events.</done>
</task>

<task type="auto">
  <name>Implement Audio Level Calculation</name>
  <files>internal/audio/audio.go</files>
  <action>
    - Update `CaptureAudio` to calculate RMS (Root Mean Square) of the audio buffer.
    - Send `AudioLevelMsg` (float64) to a new UI channel alongside audio data.
  </action>
  <verify>Check `audio.go` compiles and calculates RMS.</verify>
  <done>Audio package emits level data.</done>
</task>

<task type="auto">
  <name>Create Basic TUI Model</name>
  <files>internal/ui/model.go, internal/ui/view.go</files>
  <action>
    - Create Bubble Tea model with `TranscriptionHistory` and `AudioLevel`.
    - Implement `Init`, `Update`, and `View` methods.
    - Use `lipgloss` to define the Split View layout (Meter | Log).
  </action>
  <verify>go run cmd/livelylivecaptions/main.go (check UI renders)</verify>
  <done>Split view UI renders with placeholder data.</done>
</task>

## Success Criteria
- [ ] `TranscriptionEvent` replaces string strings.
- [ ] Audio RMS determines level meter height.
- [ ] Split view layout renders correctly.
