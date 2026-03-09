package network

import (
	"context"
	"errors"
	"net/netip"
)

type networkService struct {
	reservedOffsets []int // 予約済みIP
}

func NewNetworkService() NetworkService {
	return &networkService{
		// 0: ネットワークアドレス, 1: ブロードキャストアドレス
		reservedOffsets: []int{0, 1},
	}
}

func (s *networkService) CalculateNextAvailableIP(
	ctx context.Context,
	cidr string,
	usedIPs []string,
) (string, error) {
	// CIDRの解析
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return "", err
	}

	// Map化
	blockedMap := make(map[netip.Addr]struct{})
	for _, ipStr := range usedIPs {
		ip, err := netip.ParseAddr(ipStr)
		if err != nil {
			continue // 無効なIPはスキップ
		}
		blockedMap[ip] = struct{}{}
	}

	base := prefix.Addr()
	for _, offset := range s.reservedOffsets {
		target := base
		for j := 0; j < offset; j++ {
			target = target.Next()
		}
		blockedMap[target] = struct{}{}
	}

	lastAddr := s.broadcastAddress(prefix)
	blockedMap[lastAddr] = struct{}{}

	addr := base

	for {
		// サブネットの範囲を超えたら終了
		if !prefix.Contains(addr) {
			break
		}

		if s.isValidIP(addr, prefix, lastAddr, blockedMap) {
			return addr.String(), nil
		}

		addr = addr.Next()
	}

	return "", errors.New("no available IP addresses in subnet")
}

func (s *networkService) isValidIP(
	addr netip.Addr,
	prefix netip.Prefix,
	last netip.Addr,
	blockedMap map[netip.Addr]struct{},
) bool {
	if addr == prefix.Addr() {
		return false // ネットワークアドレス
	}
	if addr == last {
		return false // ブロードキャストアドレス
	}
	if _, exists := blockedMap[addr]; exists {
		return false // 使用済み
	}
	return true
}

// helper: ブロードキャストアドレスの計算
func (s *networkService) broadcastAddress(prefix netip.Prefix) netip.Addr {
	// ネットワークアドレスをバイトスライスで取得
	naddr := prefix.Addr().As16()
	// マスクを取得 (例: /24 なら 24)
	maskLen := prefix.Bits()

	// IPv4 (32bit) か IPv6 (128bit) かで処理を分ける
	isIPv4 := prefix.Addr().Is4()
	totalBits := 128
	if isIPv4 {
		totalBits = 32
	}

	// ホスト部のビットをすべて1にする操作
	for i := maskLen; i < totalBits; i++ {
		// IPv4の場合は末尾4バイト(12〜15)を操作する
		byteIdx := (128 - totalBits + i) / 8
		bitIdx := uint(7 - (i % 8))
		naddr[byteIdx] |= (1 << bitIdx)
	}

	return netip.AddrFrom16(naddr)
}
