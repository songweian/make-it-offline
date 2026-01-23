# Make It Offline - 离线安装包生成器 (Go 版)

本项目旨在简化为不同平台和格式创建离线安装包的过程。现已转变为纯 CLI 工具。

## 使用方法

### CLI 使用

1. 编译 CLI:
   ```bash
   go build -o make-it-offline cmd/cli/main.go
   ```

2. 运行 CLI (使用 @ 分隔应用、操作系统和版本):
   ```bash
   # 格式: ./make-it-offline <app@version> <os@version> <arch> <formats>
   
   # 示例 1: Ubuntu 上生成 nginx docker-compose 包
   ./make-it-offline nginx@1.21 ubuntu@20.04 x86_64 docker-compose
   
   # 示例 2: RedHat 7.9 上生成 nginx RPM 包
   ./make-it-offline nginx@1.21 redhat@7.9 x86_64 rpm
   
   # 示例 3: RedHat 7.9 上生成 mysql 混合包 (docker-compose + rpm)
   ./make-it-offline mysql@5.7 redhat@7.9 x86_64 docker-compose,rpm
   
   # 示例 4: aarch64 架构
   ./make-it-offline postgresql@12 redhat@7.9 aarch64 docker-compose
   ```

## 架构

- `cmd/cli/`: 命令行工具。
- `pkg/plugins/`: 插件框架目录（接口及基类）。
- `repos/`: 插件实现仓库，每个应用一个目录。
- `pkg/utils/`: 通用工具函数。

## 功能

- **支持多种操作系统**: `ubuntu`, `centos`, `redhat`, `windows`
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

## RedHat 7.9 离线安装包指南

### 在 macOS 上为 RedHat 7.9 生成离线包

本工具支持在 macOS 上为 RedHat 7.9 (x86_64 和 aarch64) 生成离线安装包。

#### 支持的格式

1. **docker-compose 格式** - 用于容器化部署
   - 生成 `docker-compose.yml` 配置文件
   - 需要目标系统已安装 Docker 和 docker-compose
   
2. **RPM 格式** - 用于传统 yum 安装 (仅限 nginx)
   - 下载对应版本的 RPM 包及其依赖
   - 需要在 RedHat 7.9 系统上执行 `yum localinstall`

#### 使用示例

```bash
# 生成 nginx docker-compose 离线包
./make-it-offline nginx@1.21 redhat@7.9 x86_64 docker-compose

# 生成 nginx RPM 离线包 (需要 Docker)
./make-it-offline nginx@1.21 redhat@7.9 x86_64 rpm

# 生成混合格式包
./make-it-offline mysql@5.7 redhat@7.9 x86_64 docker-compose,rpm
```

#### 系统要求

- **在 macOS 上**: 需要安装 Docker Desktop 或 Docker CLI
- **在 RedHat 7.9 上**: 
  - 对于 docker-compose 格式: 需要 Docker 和 docker-compose
  - 对于 RPM 格式: 仅需 yum 包管理工具

#### 部署步骤

1. 在 macOS 上生成离线包
2. 将 tar.gz 文件传输到 RedHat 7.9 系统
3. 解压包: `tar -xzf <package>.tar.gz`
4. 运行安装脚本: `./install.sh`
