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
