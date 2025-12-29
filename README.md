# Make It Offline - 离线安装包生成器

本项目旨在简化为不同平台和格式创建离线安装包的过程。

## 功能

- **支持多种操作系统**: `ubuntu`, `centos`, `windows`
- **支持指定操作系统版本**
- **支持多种架构**: `x86_64`, `aarch64`
- **多格式支持**: `docker-compose`, `rpm`, `deb`,`exe`,`msi`
- **版本化**: 支持为指定软件版本创建安装包
- **一键安装**: 生成的离线包内置一键安装脚本
- **支持安装的应用**
  - nginx
  - penpot
  - postgresql
  - mattermost
  - ONLYOFFICE
  - wiki.js
  - nextcloud
  - jenkins
  - redis
  - kafka
  - rocketmq
  - elastic
  - grafana
  - prometheus
  - zabbix
  - harbor
  - gitlab
  - jenkins
  - sonarqube
  - nexus 
  - harbor
  - minio 
  - zookeeper
  - redis 
  - zookeeper
  - mysql
  - mongodb
  - rabbitmq

## 目录结构

```
.
├── README.md               # 本文档
├── prepare_package.sh      # 离线包生成主脚本
├── output/                   # 生成的离线安装包存放目录
└── scripts/                  # 各类辅助脚本
    ├── prepare_docker_compose.sh
    ├── prepare_rpm.sh
    └── prepare_deb.sh
```

## 如何使用

### 1. 首次设置

在第一次运行之前，请确保主脚本是可执行的：

```bash
chmod +x prepare_package.sh
```

### 2. 生成离线包

执行 `prepare_package.sh` 脚本并提供所需参数。

**参数:**

1. `ARCH`: 架构 (`x86_64` 或 `arm`)
2. `TYPE`: 安装方式 (`docker-compose`, `rpm`, 或 `deb`)
3. `APP_NAME`: 应用名称
4. `APP_VERSION`: 应用版本

**示例 (使用 Docker Compose):**

```bash
./prepare_package.sh x86_64 docker-compose my-app 1.0.0
```

### 3. 在目标服务器上安装

1. 将 `output/` 目录中生成的 `.tar.gz` 包传输到目标服务器。
2. 解压并执行内置的安装脚本。

```bash
tar -zxvf my-app-1.0.0-x86_64-offline.tar.gz
cd my-app-1.0.0
./install.sh
```
