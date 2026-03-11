package network

import (
	"context"
	"errors"
	"net"
	"net/netip"
)

var (
	ErrNoAvailableCIDRs = errors.New("no available CIDR blocks")
	ErrNoAvailableIPs   = errors.New("no available IP addresses in the subnet")
)

type networkService struct {
	reservedOffsets []int // 予約済みIP
}

func NewNetworkService() NetworkService {
	return &networkService{
		// 0: ネットワークアドレス, 1: ゲートウェイ, ブロードキャストアドレスは別途計算するためここでは指定しない
		reservedOffsets: []int{0, 1},
	}
}

func (s *networkService) CalculateNextAvailableVPCCIDR(
	ctx context.Context,
	usedCidrs []string,
) (string, error) {
	usedMap := make(map[string]struct{})
	for _, cidr := range usedCidrs {
		usedMap[cidr] = struct{}{}
	}

	// 10.x.0.0/16の範囲でCIDRを生成
	for x := 0; x <= 255; x++ {
		targetCidr := netip.PrefixFrom(netip.AddrFrom4([4]byte{10, byte(x), 0, 0}), 16).String()

		if _, occupied := usedMap[targetCidr]; !occupied {
			return targetCidr, nil
		}
	}
	return "", ErrNoAvailableCIDRs
}

func (s *networkService) CalculateNextAvailableSubnet(
	ctx context.Context,
	vpcCidr string,
	usedCidrs []string,
) (string, error) {
	parentPrefix, err := netip.ParsePrefix(vpcCidr)
	if err != nil {
		return "", err
	}

	usedMap := make(map[string]struct{})
	for _, cidr := range usedCidrs {
		usedMap[cidr] = struct{}{}
	}

	addr := parentPrefix.Addr()
	// vpc /16 から /24 のサブネットを順に走査
	for {
		targetSubnet := netip.PrefixFrom(addr, 24)

		if !parentPrefix.Overlaps(targetSubnet) {
			break
		}

		if _, occupied := usedMap[targetSubnet.String()]; !occupied {
			return targetSubnet.String(), nil
		}

		// 次の /24 へ
		// ちょっと脳筋すぎるので、いつかビットシフトに変更したい
		for i := 0; i < 256; i++ {
			addr = addr.Next()
		}
		if !addr.IsValid() {
			break
		}
	}
	return "", ErrNoAvailableCIDRs
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

	return "", ErrNoAvailableIPs
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
