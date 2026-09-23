// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"reflect"
	"testing"
)

func TestLegacyGPUSupportEnabled(t *testing.T) {
	const env = legacyGPUSupportEnv
	original, present := os.LookupEnv(env)
	t.Cleanup(func() {
		if present {
			_ = os.Setenv(env, original)
		} else {
			_ = os.Unsetenv(env)
		}
	})

	cases := []struct {
		value string
		want  bool
	}{
		{"", false},
		{"0", false},
		{"false", false},
		{"1", true},
		{"true", true},
		{"YES", true},
		{"on", true},
		{"unexpected", false},
	}
	for _, tc := range cases {
		t.Run(tc.value, func(t *testing.T) {
			_ = os.Setenv(env, tc.value)
			if got := legacyGPUSupportEnabled(); got != tc.want {
				t.Fatalf("legacyGPUSupportEnabled() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestApplyLegacyGPUCompatibility(t *testing.T) {
	const env = legacyGPUSupportEnv
	_ = os.Setenv(env, "1")
	t.Cleanup(func() { _ = os.Unsetenv(env) })

	gpus := []GPUInfo{
		{Name: "NVIDIA GeForce GTX 1070", statsKey: "GPU-1070"},
		{Name: "NVIDIA GeForce GTX 980", statsKey: "GPU-980"},
		{Name: "AMD Radeon RX 7900 XTX", statsKey: "GPU-amd"},
	}
	gotReady, enabled := applyLegacyGPUCompatibility("host-123", gpus)
	if !enabled {
		t.Fatal("legacy compatibility should be enabled")
	}

	wantIDs := []string{
		"host-123:gpu:GPU-1070",
		"host-123:gpu:GPU-980",
	}
	if !reflect.DeepEqual(gotReady, wantIDs) {
		t.Fatalf("ready ids = %#v, want %#v", gotReady, wantIDs)
	}
	if gpus[0].ID != wantIDs[0] || gpus[1].ID != wantIDs[1] {
		t.Fatalf("stable GPU IDs were not attached: %#v", gpus)
	}
	if gpus[2].ID != "host-123:gpu:GPU-amd" {
		t.Fatalf("non-NVIDIA GPU should still receive a stable inventory ID: %q", gpus[2].ID)
	}
}

func TestApplyLegacyGPUCompatibilityDisabled(t *testing.T) {
	const env = legacyGPUSupportEnv
	_ = os.Unsetenv(env)

	gpus := []GPUInfo{{Name: "NVIDIA GeForce GTX 1070", statsKey: "GPU-1070"}}
	gotReady, enabled := applyLegacyGPUCompatibility("host-123", gpus)
	if enabled {
		t.Fatal("legacy compatibility should be disabled")
	}
	if gotReady != nil {
		t.Fatalf("ready ids = %#v, want nil", gotReady)
	}
	if gpus[0].ID != "" {
		t.Fatalf("GPU ID changed while compatibility was disabled: %q", gpus[0].ID)
	}
}
