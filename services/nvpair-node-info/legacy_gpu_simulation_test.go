// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

// TestLegacyGPUSimulation validates the compatibility layer with representative
// inventories for the exact GPUs this fork targets. It deliberately does not
// call detectGPUs, so it runs on ordinary GitHub-hosted runners without NVIDIA
// hardware.
func TestLegacyGPUSimulation(t *testing.T) {
	name := strings.TrimSpace(os.Getenv("NVPAIR_SIMULATED_GPU"))
	if name == "" {
		t.Fatal("NVPAIR_SIMULATED_GPU must name the simulated GPU")
	}
	vramMB, err := strconv.ParseUint(strings.TrimSpace(os.Getenv("NVPAIR_SIMULATED_VRAM_MB")), 10, 64)
	if err != nil || vramMB == 0 {
		t.Fatal("NVPAIR_SIMULATED_VRAM_MB must be a positive integer")
	}

	t.Setenv(legacyGPUSupportEnv, "1")
	gpus := []GPUInfo{{
		Name:      "NVIDIA GeForce " + name,
		VramBytes: vramMB * 1024 * 1024,
		statsKey:  "simulated-" + strings.ToLower(strings.ReplaceAll(name, " ", "-")),
	}}

	readyIDs, enabled := applyLegacyGPUCompatibility("simulated-host", gpus)
	if !enabled {
		t.Fatal("legacy compatibility should be enabled")
	}
	if len(readyIDs) != 1 {
		t.Fatalf("ready ids = %#v, want exactly one", readyIDs)
	}
	if gpus[0].ID == "" {
		t.Fatal("simulated GPU did not receive a stable ID")
	}
	if gpus[0].VramBytes != vramMB*1024*1024 {
		t.Fatalf("VRAM = %d, want %d", gpus[0].VramBytes, vramMB*1024*1024)
	}
}
