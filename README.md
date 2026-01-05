# Make It Offline - 离线安装包生成器 (Go 版)

本项目旨在简化为不同平台和格式创建离线安装包的过程。现已转变为纯 CLI 工具。

## 使用方法

### CLI 使用

1. 编译 CLI:
   ```bash
   go build -o make-it-offline cmd/cli/main.go
   ```

2. 运行 CLI:
   ```bash
   ./make-it-offline nginx@1.21 ubuntu-20.04-x86_64 docker-compose
   ```

## 架构

- `cmd/cli/`: 命令行工具。
- `pkg/plugins/`: 插件框架目录（接口及基类）。
- `repos/`: 插件实现仓库，每个应用一个目录。
- `pkg/utils/`: 通用工具函数。

## 功能

- **支持多种操作系统**: `ubuntu`, `centos`, `windows`
- **支持指定操作系统版本**
- **支持多种架构**: `x86_64`, `aarch64`
- **多格式支持**: `docker-compose`, `rpm`, `deb`,`exe`,`msi`
- **版本化**: 支持为指定软件版本创建安装包
- **一键安装**: 生成的离线包内置一键安装脚本
- **支持安装的应用**
  - nginx
  - postgresql
  - redis
  - mattermost
  - mysql
  - prometheus
  - grafana
  - (更多应用正在移植中...)
