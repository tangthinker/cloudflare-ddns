# Cloudflare DDNS

一个基于 Go 语言的 Cloudflare 动态 DNS 客户端，可以自动根据本地网络接口更新 IPv6 DNS 记录。

## 功能特点

- 自动从指定网络接口检测 IPv6 地址
- 使用 Cloudflare API 更新多个 DNS 记录
- 可配置的更新间隔
- 错误处理和重试机制
- 基于 YAML 的配置文件
- 支持多域名更新
- 智能记录管理（根据记录是否存在自动创建或更新）
- 支持自定义配置文件

## 配置说明

创建 `config.yaml` 文件，结构如下：

```yaml
cloudflare:
  api_token: "your-api-token-here"  # 从 Cloudflare 控制面板获取
  zone_id: "your-zone-id-here"      # 从 Cloudflare 控制面板获取
  domains:                          # 需要更新的域名列表
    - "sub1.example.com"
    - "sub2.example.com"
    - "sub3.example.com"

network:
  interface: "en0"  # 要监控的网络接口名称

interval:
  success: 600  # 成功更新后的等待时间（秒）
  error: 10     # 发生错误后的等待时间（秒）
```

## 前置条件

1. 具有 API 访问权限的 Cloudflare 账户
2. 由 Cloudflare 管理的域名
3. IPv6 网络连接
4. Go 1.24 或更高版本

## 快速开始

1. 获取 Cloudflare API 令牌：
   - 登录 Cloudflare 控制面板
   - 进入 "我的个人资料" > "API 令牌"
   - 创建具有 DNS 编辑权限的新令牌
   - 复制令牌

2. 获取区域 ID：
   - 在 Cloudflare 控制面板中选择你的域名
   - 在右侧边栏找到 "区域 ID"
   - 复制 ID

3. 更新 `config.yaml` 文件：
   - 粘贴你的 API 令牌
   - 粘贴你的区域 ID
   - 添加需要更新的域名列表
   - 设置网络接口名称
   - 根据需要调整时间间隔

4. 运行程序：
```bash
# 使用默认的 config.yaml
go run main.go

# 使用自定义配置文件
go run main.go -config /path/to/your/config.yaml

# 构建并使用自定义配置文件运行
go build
./cloudflare-ddns -config /path/to/your/config.yaml
```

## 命令行选项

- `-config string`：配置文件路径（默认：`config.yaml`）
  ```bash
  # 示例：
  ./cloudflare-ddns -config custom_config.yaml
  ./cloudflare-ddns -config /etc/cloudflare-ddns/config.yaml
  ./cloudflare-ddns -config ../configs/production.yaml
  ```

## 工作原理

1. 程序从指定的配置文件（默认：`config.yaml`）读取配置
2. 每 10 分钟（可配置）检查指定网络接口的 IPv6 地址
3. 对于配置中的每个域名：
   - 检查是否存在 AAAA 记录
   - 如果不存在则创建新记录
   - 如果 IPv6 地址发生变化则更新记录
4. 如果发生错误，将在 10 秒后（可配置）重试
5. 记录所有成功更新和错误信息

## 错误处理

- 如果获取 IPv6 地址失败，将在错误间隔后重试
- 如果更新某个域名失败，将继续处理其他域名
- 如果部分域名更新成功，将报告部分成功状态
- 在日志中提供详细的错误信息

## 日志记录

程序提供详细的日志记录，包括：
- IPv6 地址检测
- DNS 记录查询
- 记录创建/更新
- 成功/失败状态
- 错误详情

## 贡献

欢迎提交问题和改进建议！