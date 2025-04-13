# Cloudflare DDNS

这是一个基于 Go 语言开发的 Cloudflare 动态 DNS 客户端，可以根据本地网络接口自动更新 IPv6 DNS 记录。

## 功能特点

- 自动检测指定网络接口的 IPv6 地址
- 使用 Cloudflare API 更新 DNS 记录
- 可配置的更新间隔时间
- 错误处理和重试机制
- 基于 YAML 的配置文件

## 配置说明

创建一个 `config.yaml` 配置文件，结构如下：

```yaml
cloudflare:
  api_token: "你的-API-令牌"  # 从 Cloudflare 控制面板获取
  zone_id: "你的-区域-ID"     # 从 Cloudflare 控制面板获取
  domain: "你的域名.com"      # 需要更新的域名

network:
  interface: "en0"  # 要监控的网络接口名称

interval:
  success: 600  # 成功更新后的等待时间（秒）
  error: 10     # 发生错误后的重试等待时间（秒）
```

## 使用方法

1. 从 Cloudflare 控制面板获取 API 令牌
2. 从 Cloudflare 控制面板获取区域 ID
3. 更新 `config.yaml` 文件中的配置信息
4. 运行程序：

```bash
go run main.go
```

或者编译后运行：

```bash
go build
./cloudflare-ddns
```

## 系统要求

- Go 1.24 或更高版本
- 具有 API 访问权限的 Cloudflare 账户
- 由 Cloudflare 管理的域名
- IPv6 网络连接