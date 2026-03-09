package network

import (
	"context"
	"errors"
	"net"
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
		blockedMap[ip.Unmap()] = struct{}{}
	}

	base := prefix.Addr()
	for _, offset := range s.reservedOffsets {
		target := base
		for j := 0; j < offset; j++ {
			target = target.Next()
		}
		blockedMap[target.Unmap()] = struct{}{}
	}

	lastAddr := s.broadcastAddress(prefix).Unmap()
	blockedMap[lastAddr] = struct{}{}

	addr := base

	for {
		current := addr.Unmap()

		if !prefix.Contains(current) {
			break
		}

		if _, blocked := blockedMap[current]; !blocked {
			return current.String(), nil
		}

		addr = addr.Next()
		if !addr.IsValid() {
			break
		}
	}

	return "", errors.New("no available IP addresses in subnet")
}

// helper: ブロードキャストアドレスの計算
func (s *networkService) broadcastAddress(prefix netip.Prefix) netip.Addr {
	addrBytes := prefix.Addr().As16()

	// アドレスがIPv4かIPv6かで、マスクの適用範囲を変える
	var mask net.IPMask
	if prefix.Addr().Is4() {
		// IPv4の場合、マスクは32ビット分
		mask = net.CIDRMask(prefix.Bits(), 32)
		// IPv4-mapped IPv6は末尾4バイト(12-15)に実体がある
		for i := range 4 {
			addrBytes[12+i] |= ^mask[i]
		}
	} else {
		// IPv6の場合
		mask = net.CIDRMask(prefix.Bits(), 128)
		for i := range 16 {
			addrBytes[i] |= ^mask[i]
		}
	}

	return netip.AddrFrom16(addrBytes)
}
