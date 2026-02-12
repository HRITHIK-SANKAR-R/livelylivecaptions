# DECISIONS.md

> **Status**: Active

## Log
- [ ] [YYYY-MM-DD] Initial project setup
- [ ] [YYYY-MM-DD] **Architecture Change**: Switch from Python IPC to Pure Go implementation (User Request).
- [ ] [YYYY-MM-DD] **Phase 2 UI**: Split view with audio level meter + side logs.
- [ ] [YYYY-MM-DD] **Phase 2 Architecture**: Use `TranscriptionEvent` struct for internal messaging.
- [ ] [YYYY-MM-DD] **Phase 2 Controls**: Implement runtime controls (mute, device switch).
- [ ] [2026-02-12] **Phase 3 GPU Support**: Use `run.sh` to dynamically link CUDA 12 libraries from venv.
- [ ] [2026-02-12] **Phase 3 Model Choice**: Standardized on 2023-02-17 CPU model for all modes to ensure metadata compatibility.
- [ ] [2026-02-12] **Phase 3 Provider**: Implemented `internal/hardware` for automatic `nvidia-smi` detection.
- [ ] [2026-02-12] **Phase 4 Testing**: Adopt Table-Driven tests for logic and Interface-based mocking for hardware (AudioDevice/Provider).
- [ ] [2026-02-12] **Phase 4 Integration**: Implement a "Golden File" regression pipeline with Levenshtein Distance (WER) metrics.
- [ ] [2026-02-12] **Phase 4 Config**: Use Viper for layered configuration (CLI > Env > YAML > Defaults).
- [ ] [2026-02-12] **Phase 4 Reliability**: Implement a "Startup Probe" for hardware context before engine initialization.
- [ ] [2026-02-12] **Phase 4 Logging**: Implement a memory ring-buffer with a TUI-accessible debug tab.
