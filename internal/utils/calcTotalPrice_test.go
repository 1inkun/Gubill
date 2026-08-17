package utils

import "testing"

func TestCalculateTotalPrice(t *testing.T) {
	const price = int64(500)
	cases := []struct {
		name     string
		duration int64 // 秒
		want     int
	}{
		{"零时长", 0, 0},
		{"30分钟内", 10 * 60, 500},
		{"恰好30分钟", 30 * 60, 500},
		{"31分钟进位", 31 * 60, 1000},
		{"5小时", 5 * 60 * 60, 5000},
		{"恰好12小时", 12 * 60 * 60, 5000},    // 封顶 10x
		{"12小时加1分钟", 12*60*60 + 60, 5500}, // 超出部分继续按单位计
		{"恰好24小时", 24 * 60 * 60, 10000},   // 封顶 20x
		{"25小时递归", 25 * 60 * 60, 11000},   // 24h(10000) + 1h(1000)
		{"30小时递归", 30 * 60 * 60, 15000},   // 24h(10000) + 6h封顶(5000)
		{"37小时递归", 37 * 60 * 60, 16000},   // 24h(10000) + 13h(6000)
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateTotalPrice(tc.duration, price)
			if got != tc.want {
				t.Errorf("CalculateTotalPrice(%d, %d) = %d, want %d", tc.duration, price, got, tc.want)
			}
		})
	}
}
