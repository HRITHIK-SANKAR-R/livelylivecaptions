# ROADMAP.md

> **Current Phase**: Not started
> **Milestone**: v1.0

## Must-Haves (from SPEC)
- [ ] Zero-latency Experience (Visual <100ms, Transcription <250ms)
- [ ] Offline Operation
- [ ] Resource Efficient (RTX 3050 / Modern CPU)
- [ ] High Accuracy (English)
- [ ] Minimal Dependencies

## Phases

### Phase 1: Foundation
**Status**: ⬜ Not Started
**Objective**: Validate Audio Capture and On-Device Transcription (Go)
**Deliverable**: CLI/TUI tool that captures audio and prints transcripts using Sherpa-ONNX (Go bindings).
**Tasks**:
- [ ] Go Transcriber: Verify/Refine `internal/transcriber` with Sherpa-ONNX Go bindings.
- [ ] Go Audio Capture: Verify `internal/audio` (PortAudio/Malgo).
- [ ] Integration: Ensure `main.go` wires Audio -> Transcriber -> UI correctly.

### Phase 2: UI & Optimization
**Status**: ⬜ Not Started
**Objective**: Build TUI and optimize performance
**Deliverable**: Functional TUI with live transcription and audio visualization.
**Tasks**:
- [ ] Bubble Tea UI: Implement model, update, view.
- [ ] Real-time Updates: Handle JSON parsing and audio level metrics.
- [ ] Performance Tuning: Profile latency, optimize buffers.

### Phase 3: Model Integration
**Status**: ⬜ Not Started
**Objective**: Integrate recommended models and benchmarking
**Deliverable**: Multi-model support with automatic hardware detection.
**Tasks**:
- [ ] GPU Model Setup (Sherpa-ONNX Zipformer).
- [ ] CPU Model Setup (Faster-Whisper).
- [ ] Auto-Detection mechanism.

### Phase 4: Polish & Testing
**Status**: ⬜ Not Started
**Objective**: Production readiness
**Deliverable**: Production-ready binary with documentation.
**Tasks**:
- [ ] Error handling and configuration.
- [ ] Logging and debugging flags.
- [ ] Unit and Integration tests.
- [ ] Documentation and Packaging.
