# SQLite 纯 Go 驱动替换设计

## 背景与目标
当前后端使用 `gorm.io/driver/sqlite`，该驱动内部依赖 `github.com/mattn/go-sqlite3`，需要 CGO。构建脚本默认 `CGO_ENABLED=0`，导致运行时出现 sqlite3 stub 报错，用户无法登录。目标是将 SQLite 驱动彻底切换到纯 Go 实现（`modernc.org/sqlite`），确保 `CGO_ENABLED=0` 的静态构建可稳定运行，并保持现有配置与数据文件兼容。

## 方案选择
采用方案 A：使用 `github.com/glebarez/sqlite` 作为 GORM Dialector（内部基于 `modernc.org/sqlite`）。该方案保留 GORM API 与现有数据模型，改动面小，风险低，性能与稳定性可接受。方案 B（继续用 gorm sqlite 并切换 DriverName）仍残留 go-sqlite3 依赖，不满足“彻底切换”的要求；方案 C 改为 `database/sql` 需要大量重构，不符合当前交付节奏。

## 架构与数据流
`internal/db` 仍作为数据库入口，`InitDB` 继续基于配置拼装 DSN 并初始化 GORM。改动点仅在 Dialector 引用：将 `gorm.io/driver/sqlite` 替换为 `github.com/glebarez/sqlite`。初始化流程保持不变：创建连接 → `AutoMigrate` → 提供全局 `db` 实例。数据文件 `data/data.db` 保持兼容，仍使用 WAL、同步、缓存等参数配置。

## 错误处理与回滚
保持现有错误处理逻辑：`InitDB` 返回错误时记录日志并终止启动。若迁移后发现兼容性问题，可回滚到旧驱动并启用 CGO 构建作为应急方案。为降低风险，将新增 `CGO_ENABLED=0` 的数据库初始化测试，提前阻止回归。

## 测试策略与验收
新增单元测试验证在 `CGO_ENABLED=0` 下数据库初始化与迁移可用；测试使用临时目录数据库文件，避免影响真实数据。验收包含：后端 `go test`、前端构建、后端二进制构建、服务重启，并验证登录可用。构建与 CI 将显式声明 `CGO_ENABLED=0`，确保纯 Go 驱动路径被持续验证。
