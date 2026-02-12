# SPEC.md — Project Specification

> **Status**: `FINALIZED`

## Vision
Livelylivecaptions is a zero-network, offline, real-time speech transcription system designed to deliver visual feedback with <100ms latency. The system leverages Go for hardware control and UI rendering, Python for AI inference, and employs inter-process communication via stdin/stdout pipes to achieve an "instant" feel. The project prioritizes minimal latency, offline operation, and resource efficiency to run on consumer-grade hardware including RTX 3050 4GB GPUs and CPU-only systems.

## Goals
1.  **Zero-latency Experience:** Visual feedback <100ms, full transcription <250ms.
2.  **Offline Operation:** No internet connectivity required.
3.  **Resource Efficient:** Run on RTX 3050 4GB GPU or modern CPUs.
4.  **High Accuracy:** English audio transcription with production-quality accuracy.
5.  **Minimal Dependencies:** Self-contained deployment with minimal external dependencies.

## Non-Goals (Out of Scope)
- Multi-language support (Phase 1 focuses on English only).
- Cloud-based transcription services.
- Mobile platform support (desktop-first approach).
- Real-time translation or diarization.

## Users
- **Desktop Users:** Individuals who need real-time captions for audio/video content on their local machine without internet access.
- **Privacy-Conscious Users:** Users who require their audio data to remain local and secure.
- **Hardware Owners:** Users with consumer-grade GPUs (RTX 3050+) or modern CPUs who want to leverage their hardware for AI tasks.

## Constraints
- **Hardware:** Must run on RTX 3050 4GB GPU or modern CPUs (Intel i5/i7 8th gen+ / AMD Ryzen 5/7).
- **Latency:** Visual feedback <100ms, transcription <250ms.
- **Memory:** <2GB VRAM (GPU) or <1GB RAM (CPU).
- **OS:** Linux (Ubuntu 22.04+), macOS 12+, Windows 10/11.

## Success Criteria
- [ ] Visual latency <100ms.
- [ ] Transcription latency <250ms.
- [ ] GPU VRAM usage <2GB.
- [ ] Works offline.
- [ ] Stable for 1+ hour continuous use.
