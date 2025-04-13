package network

import (
	"net"
	"testing"
)

func TestGetIPv6Address(t *testing.T) {
	// 创建 NetworkManager 实例
	nm := NewNetworkManager()

	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		t.Fatalf("获取网络接口失败: %v", err)
	}

	// 测试每个网络接口
	for _, iface := range interfaces {
		t.Run(iface.Name, func(t *testing.T) {
			// 跳过回环接口和未启用的接口
			if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
				t.Skipf("跳过接口 %s: 回环或未启用", iface.Name)
			}

			// 获取IPv6地址
			ipv6, err := nm.GetIPv6Address(iface.Name)
			if err != nil {
				// 如果接口没有IPv6地址，这是预期的结果
				if err.Error() == "no IPv6 address found for interface "+iface.Name {
					t.Logf("接口 %s 没有IPv6地址，这是预期的结果", iface.Name)
					return
				}
				t.Errorf("获取IPv6地址失败: %v", err)
				return
			}

			// 验证获取到的地址是有效的IPv6地址
			ip := net.ParseIP(ipv6)
			if ip == nil {
				t.Errorf("获取到的地址 %s 不是有效的IP地址", ipv6)
				return
			}

			if ip.To4() != nil {
				t.Errorf("获取到的地址 %s 是IPv4地址，而不是IPv6地址", ipv6)
				return
			}

			if ip.IsLinkLocalUnicast() {
				t.Errorf("获取到的地址 %s 是链路本地地址", ipv6)
				return
			}

			t.Logf("成功获取接口 %s 的IPv6地址: %s", iface.Name, ipv6)
		})
	}
}

// 测试无效接口名称
func TestGetIPv6AddressInvalidInterface(t *testing.T) {
	nm := NewNetworkManager()
	_, err := nm.GetIPv6Address("non_existent_interface")
	if err == nil {
		t.Error("期望获取无效接口名称时返回错误，但没有返回错误")
	}
}
