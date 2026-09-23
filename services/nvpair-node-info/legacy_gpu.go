// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"strconv"
	"strings"
)

const legacyGPUSupportEnv = "NVPAIR_LEGACY_GPU_SUPPORT"

// legacyGPUSupportEnabled enables the fork-specific hardware-policy override.
// It is deliberately opt-in so an unmodified PAIR deployment keeps exactly
// the upstream behavior.
//
// Accepted true values are 1, true, yes, and on (case-insensitive).
func legacyGPUSupportEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(legacyGPUSupportEnv))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

// isNvidiaGPU reports whether the inventory name identifies an NVIDIA adapter.
// The detector already supplies vendor names on Windows and Linux; keeping the
// check name-based makes this helper platform-neutral and avoids inventing a
// second hardware database just for the legacy override.
func isNvidiaGPU(name string) bool {
	return strings.Contains(strings.ToLower(name), "nvidia")
}

// gpuHardwareID returns the stable UI identity used by the fork's readiness
// list. statsKey is UUID/LUID/IORegistry-derived on platforms with a dynamic
// GPU source, while the host UUID + index fallback keeps the contract usable
// for a static-only detector.
func gpuHardwareID(hostUUID string, gpu GPUInfo, index int) string {
	key := gpu.statsKey
	if key == "" {
		key = strconv.Itoa(index)
	}
	if hostUUID == "" {
		return "gpu:" + key
	}
	return hostUUID + ":gpu:" + key
}

// applyLegacyGPUCompatibility decorates the node inventory with stable GPU IDs
// and returns the IDs that the desktop may treat as inference-ready.
//
// This does NOT make an inference engine capable of using a GPU it cannot
// actually run on. It only removes PAIR's hardware-policy distinction for an
// explicitly opted-in NVIDIA node. The engine and model still have to load
// successfully.
func applyLegacyGPUCompatibility(hostUUID string, gpus []GPUInfo) ([]string, bool) {
	if !legacyGPUSupportEnabled() {
		return nil, false
	}

	ready := make([]string, 0, len(gpus))
	for i := range gpus {
		gpus[i].ID = gpuHardwareID(hostUUID, gpus[i], i)
		if isNvidiaGPU(gpus[i].Name) {
			ready = append(ready, gpus[i].ID)
		}
	}
	return ready, true
}
