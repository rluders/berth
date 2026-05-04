package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateStat(t *testing.T) {
	tests := []struct {
		name        string
		input       statsJSON
		wantCPU     float64
		wantMemUsed uint64
		wantMemLim  uint64
	}{
		{
			name: "normal CPU delta",
			input: statsJSON{
				CPUStats: struct {
					CPUUsage struct {
						TotalUsage  uint64   `json:"total_usage"`
						PercpuUsage []uint64 `json:"percpu_usage"`
					} `json:"cpu_usage"`
					SystemCPUUsage uint64 `json:"system_cpu_usage"`
					OnlineCPUs     uint32 `json:"online_cpus"`
				}{
					CPUUsage:       struct{ TotalUsage uint64 `json:"total_usage"`; PercpuUsage []uint64 `json:"percpu_usage"` }{TotalUsage: 2000},
					SystemCPUUsage: 10000,
					OnlineCPUs:     2,
				},
				PreCPUStats: struct {
					CPUUsage struct {
						TotalUsage uint64 `json:"total_usage"`
					} `json:"cpu_usage"`
					SystemCPUUsage uint64 `json:"system_cpu_usage"`
				}{
					CPUUsage:       struct{ TotalUsage uint64 `json:"total_usage"` }{TotalUsage: 1000},
					SystemCPUUsage: 5000,
				},
				MemoryStats: struct {
					Usage uint64            `json:"usage"`
					Limit uint64            `json:"limit"`
					Stats map[string]uint64 `json:"stats"`
				}{
					Usage: 512,
					Limit: 1024,
				},
			},
			// cpuDelta=1000, systemDelta=5000, numCPUs=2 → (1000/5000)*2*100 = 40%
			wantCPU:     40.0,
			wantMemUsed: 512,
			wantMemLim:  1024,
		},
		{
			name: "zero systemDelta → 0% CPU",
			input: statsJSON{
				CPUStats: struct {
					CPUUsage struct {
						TotalUsage  uint64   `json:"total_usage"`
						PercpuUsage []uint64 `json:"percpu_usage"`
					} `json:"cpu_usage"`
					SystemCPUUsage uint64 `json:"system_cpu_usage"`
					OnlineCPUs     uint32 `json:"online_cpus"`
				}{
					CPUUsage:       struct{ TotalUsage uint64 `json:"total_usage"`; PercpuUsage []uint64 `json:"percpu_usage"` }{TotalUsage: 2000},
					SystemCPUUsage: 5000,
					OnlineCPUs:     2,
				},
				PreCPUStats: struct {
					CPUUsage struct {
						TotalUsage uint64 `json:"total_usage"`
					} `json:"cpu_usage"`
					SystemCPUUsage uint64 `json:"system_cpu_usage"`
				}{
					CPUUsage:       struct{ TotalUsage uint64 `json:"total_usage"` }{TotalUsage: 1000},
					SystemCPUUsage: 5000, // same → delta = 0
				},
			},
			wantCPU:     0,
			wantMemUsed: 0,
			wantMemLim:  0,
		},
		{
			name: "zero cpuDelta → 0% CPU",
			input: statsJSON{
				CPUStats: struct {
					CPUUsage struct {
						TotalUsage  uint64   `json:"total_usage"`
						PercpuUsage []uint64 `json:"percpu_usage"`
					} `json:"cpu_usage"`
					SystemCPUUsage uint64 `json:"system_cpu_usage"`
					OnlineCPUs     uint32 `json:"online_cpus"`
				}{
					CPUUsage:       struct{ TotalUsage uint64 `json:"total_usage"`; PercpuUsage []uint64 `json:"percpu_usage"` }{TotalUsage: 1000},
					SystemCPUUsage: 10000,
					OnlineCPUs:     2,
				},
				PreCPUStats: struct {
					CPUUsage struct {
						TotalUsage uint64 `json:"total_usage"`
					} `json:"cpu_usage"`
					SystemCPUUsage uint64 `json:"system_cpu_usage"`
				}{
					CPUUsage:       struct{ TotalUsage uint64 `json:"total_usage"` }{TotalUsage: 1000}, // same → delta = 0
					SystemCPUUsage: 5000,
				},
			},
			wantCPU:     0,
			wantMemUsed: 0,
			wantMemLim:  0,
		},
		{
			name: "OnlineCPUs=0 falls back to len(percpu)",
			input: statsJSON{
				CPUStats: struct {
					CPUUsage struct {
						TotalUsage  uint64   `json:"total_usage"`
						PercpuUsage []uint64 `json:"percpu_usage"`
					} `json:"cpu_usage"`
					SystemCPUUsage uint64 `json:"system_cpu_usage"`
					OnlineCPUs     uint32 `json:"online_cpus"`
				}{
					CPUUsage: struct {
						TotalUsage  uint64   `json:"total_usage"`
						PercpuUsage []uint64 `json:"percpu_usage"`
					}{
						TotalUsage:  2000,
						PercpuUsage: []uint64{1000, 1000, 1000, 1000}, // 4 CPUs
					},
					SystemCPUUsage: 10000,
					OnlineCPUs:     0, // force fallback
				},
				PreCPUStats: struct {
					CPUUsage struct {
						TotalUsage uint64 `json:"total_usage"`
					} `json:"cpu_usage"`
					SystemCPUUsage uint64 `json:"system_cpu_usage"`
				}{
					CPUUsage:       struct{ TotalUsage uint64 `json:"total_usage"` }{TotalUsage: 1000},
					SystemCPUUsage: 5000,
				},
			},
			// cpuDelta=1000, systemDelta=5000, numCPUs=4 → (1000/5000)*4*100 = 80%
			wantCPU:     80.0,
			wantMemUsed: 0,
			wantMemLim:  0,
		},
		{
			name: "memory cache subtracted",
			input: statsJSON{
				CPUStats: struct {
					CPUUsage struct {
						TotalUsage  uint64   `json:"total_usage"`
						PercpuUsage []uint64 `json:"percpu_usage"`
					} `json:"cpu_usage"`
					SystemCPUUsage uint64 `json:"system_cpu_usage"`
					OnlineCPUs     uint32 `json:"online_cpus"`
				}{OnlineCPUs: 1},
				MemoryStats: struct {
					Usage uint64            `json:"usage"`
					Limit uint64            `json:"limit"`
					Stats map[string]uint64 `json:"stats"`
				}{
					Usage: 1024,
					Limit: 4096,
					Stats: map[string]uint64{"cache": 256},
				},
			},
			wantCPU:     0,
			wantMemUsed: 768, // 1024 - 256
			wantMemLim:  4096,
		},
		{
			name: "memory without cache key uses raw usage",
			input: statsJSON{
				CPUStats: struct {
					CPUUsage struct {
						TotalUsage  uint64   `json:"total_usage"`
						PercpuUsage []uint64 `json:"percpu_usage"`
					} `json:"cpu_usage"`
					SystemCPUUsage uint64 `json:"system_cpu_usage"`
					OnlineCPUs     uint32 `json:"online_cpus"`
				}{OnlineCPUs: 1},
				MemoryStats: struct {
					Usage uint64            `json:"usage"`
					Limit uint64            `json:"limit"`
					Stats map[string]uint64 `json:"stats"`
				}{
					Usage: 512,
					Limit: 2048,
					Stats: map[string]uint64{},
				},
			},
			wantCPU:     0,
			wantMemUsed: 512,
			wantMemLim:  2048,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateStat(tt.input)
			require.InDelta(t, tt.wantCPU, got.CPUPercent, 0.001)
			assert.Equal(t, tt.wantMemUsed, got.MemUsage)
			assert.Equal(t, tt.wantMemLim, got.MemLimit)
		})
	}
}
