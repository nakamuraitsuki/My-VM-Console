package network_test

import (
	"context"
	"testing"

	"example.com/m/internal/domain/network"
)

func TestCalculateNextAvailableIP(t *testing.T) {
	svc := network.NewNetworkService()
	ctx := context.Background()

	tests := []struct {
		name    string
		cidr    string
		usedIPs []string
		want    string
		wantErr bool
	}{
		{
			name:    "正常系: 未使用の最初のIP(.2)が返ること",
			cidr:    "10.0.0.0/24",
			usedIPs: []string{},
			want:    "10.0.0.2",
			wantErr: false,
		},
		{
			name:    "正常系: 中間の空きIPを埋めること",
			cidr:    "10.0.0.0/24",
			usedIPs: []string{"10.0.0.2", "10.0.0.4"},
			want:    "10.0.0.3",
			wantErr: false,
		},
		{
			name:    "正常系: 連続した使用済みIPの次を返すこと",
			cidr:    "10.0.0.0/24",
			usedIPs: []string{"10.0.0.2", "10.0.0.3", "10.0.0.4"},
			want:    "10.0.0.5",
			wantErr: false,
		},
		{
			name:    "異常系: サブネットが満杯(IP枯渇)の場合",
			cidr:    "192.168.1.0/30", // .0(net), .1(res), .2(host), .3(broad)
			usedIPs: []string{"192.168.1.2"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "異常系: 無効なCIDR形式",
			cidr:    "invalid-cidr",
			usedIPs: []string{},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.CalculateNextAvailableIP(ctx, tt.cidr, tt.usedIPs)

			// エラーの有無が期待通りかチェック
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateNextAvailableIP() error = %v, wantErr %v and got IP: %v", err, tt.wantErr, got)
				return
			}

			// 期待したIPが返ってきているかチェック
			if got != tt.want {
				t.Errorf("CalculateNextAvailableIP() got = %v, want %v", got, tt.want)
			}
		})
	}
}
