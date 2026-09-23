// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"strings"
	"testing"
)

// TestLegacyGPUHardwareOnHost is opt-in because ordinary CI runners do not have
// the legacy GeForce adapters this fork is intended to validate.
//
// The test exercises the real platform GPU detector (DXGI on Windows,
// nvidia-smi/ghw on Linux) and then checks that the detected adapters are
// accepted by the fork's legacy compatibility metadata path.
func TestLegacyGPUHardwareOnHost(t *testing.T) {
	if os.Getenv("NVPAIR_HARDWARE_GPU_TEST") != "1" {
		t.Skip("set NVPAIR_HARDWARE_GPU_TEST=1 on a physical GPU runner")
	}

	expected := strings.TrimSpace(os.Getenv("NVPAIR_EXPECTED_GPU"))
	if expected == "" {
		t.Fatal("NVPAIR_EXPECTED_GPU must name the GPU family to validate")
	}

	gpus := detectGPUs()
	if len(gpus) == 0 {
		t.Fatal("no GPU adapters were detected")
	}

	expectedLower := strings.ToLower(expected)
	var matched bool
	for _, gpu := range gpus {
		if strings.Contains(strings.ToLower(gpu.Name), expectedLower) {
			matched = true
			if gpu.VramBytes == 0 {
				t.Fatalf("matched %q but reported zero VRAM", gpu.Name)
			}
		}
	}

	if !matched {
		names := make([]string, 0, len(gpus))
		for _, gpu := range gpus {
			names = append(names, gpu.Name)
		}
		t.Fatalf("expected GPU containing %q, detected: %s", expected, strings.Join(names, ", "))
	}

	t.Setenv(legacyGPUSupportEnv, "1")
	readyIDs, enabled := applyLegacyGPUCompatibility("hardware-test-host", gpus)
	if !enabled {
		t.Fatal("legacy GPU compatibility should be enabled")
	}
	if len(readyIDs) == 0 {
		t.Fatal("legacy compatibility produced no inference-ready NVIDIA GPU IDs")
	}
}
