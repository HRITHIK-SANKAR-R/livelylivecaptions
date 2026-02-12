# RESEARCH.md - Phase 2: UI & Optimization

## Objectives
- Create a Terminal User Interface (TUI) using Bubble Tea.
- Implement a Split View layout:
  - **Left**: Audio Level Meter (Vertical).
  - **Right**: Transcription Log (Scrollable, "Subtitle" style).
- Optimize performance for <100ms visual latency.
- Implement runtime controls (Mute, Device Switch).

## Technical Approach

### 1. UI Framework: Bubble Tea + Lip Gloss
- **Model**: Stores application state (`TranscriptionHistory`, `AudioLevel`, `IsMuted`, `SelectedDevice`).
- **Update**: Handles messages:
  - `TranscriptionMsg`: New text from transcriber.
  - `AudioLevelMsg`: RMS value from audio capture.
  - `KeyMsg`: User input (q=quit, m=mute, d=device).
- **View**: Renders the split layout using `lipgloss.JoinHorizontal`.

### 2. Audio Level Meter
- **Source**: `internal/audio` needs to calculate RMS of the audio buffer.
- **Visualization**: A vertical bar rendering based on the RMS value (0.0 - 1.0).
- **Latency**: High frequency updates (e.g., every 50ms) sent to UI.

### 3. Architecture Refactor (Pure Go)
- **Current**: `transcriber` outputs string.
- **New**: `transcriber` outputs `TranscriptionEvent` struct.
  ```go
  type TranscriptionEvent struct {
      Text string
      IsFinal bool
      Confidence float64
  }
  ```
- **Concurrency**: Audio capture and Transcriber run in separate goroutines. UI runs in the main thread (Bubble Tea requirement).

### 4. Performance Considerations
- **Buffered Channels**: Use buffered channels for audio data and UI messages to prevent blocking.
- **Render Throttling**: TUI updates might need throttling if audio level messages are too frequent (limit to 60fps equivalent).

## Risks
- **Message Flood**: High frequency audio level updates could overwhelm the Bubble Tea event loop.
  - *Mitigation*: Batch updates or throttle rendering.
- **Complex Layout**: Lip Gloss layouts can be tricky with dynamic heights.
  - *Mitigation*: Keep the layout fixed height or use `viewport` for the log.
