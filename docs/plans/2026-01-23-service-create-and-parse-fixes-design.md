# 新增服务向导 + systemctl/sysv 解析修复 设计文档

## 目标
- 新增服务支持所有管理器（systemd/sysv/openrc/supervisor/pm2），提供独立新增向导与模板选择。
- 修复 systemctl 路径/执行与 sysv 详情解析，确保列表与详情字段稳定展示。
- 保持现有 UI 风格与 i18n，前端/后端均满足测试覆盖率要求。

## 范围与约束
- 新增入口位于批量操作按钮区，同一行；无需再选择管理器，使用当前已选管理器。
- 向导为精致悬浮对话框，非右侧抽屉。
- 编辑仍沿用“表单/原始文本切换”，新增向导复用该交互。
- 后端以 root 运行；缺失管理器需提示并由用户确认安装。

## 前端交互设计（悬浮向导）
- 入口：批量操作按钮左侧新增“新增”按钮，未选择服务也可用。
- 形态：居中悬浮对话框（玻璃质感/柔和阴影），不遮挡左侧列表上下文。
- 步骤：
  1) 模板选择：仅展示当前管理器的模板卡片（简短标题+一行摘要）。
  2) 表单配置：通用字段 + 扩展区；提供“表单/原始文本”切换。
  3) 确认与动作：勾选“立即启动/开机自启”，展示写入路径与 root 提示。
- 交互：顶部展示当前管理器徽标与名称；底部“上一步/下一步/创建”。

## 后端接口与数据流
- 新增接口：`POST /api/v1/{manager}/create`
- 请求体：
  - 通用字段：`name/description/exec/workdir/user/env`
  - `template_id`：模板标识
  - `raw_text`：可选，原始文本编辑保存
  - `extra`：管理器扩展字段（KV）
  - `options`：`start/autostart`
- 流程：
  1) 根据模板渲染配置文本（或使用 raw_text）。
  2) 写入默认路径（systemd `/etc/systemd/system`；sysv/openrc `/etc/init.d`；supervisor `/etc/supervisor/conf.d`；pm2 生态文件/命令）。
  3) 触发 reload/update；按 options 执行 start/enable。
  4) 返回统一 `ServiceDetail` + `extra`。

## 解析修复与兼容策略
- systemd：
  - `exec.LookPath` 探测 `systemctl/journalctl`，失败或输出含 “System has not been booted…” 即返回 `ErrManagerUnavailable`。
  - `systemctl show` 增加读取 `ExecStart` 与 `UnitFileState`，填充 `extra.exec` / `extra.autostart`。
- sysv：
  - 列表继续解析 `[ + ] name` 形式。
  - 详情改为独立解析：若输出包含 `Active:`/`Main PID:` 等 systemd 风格字段，则按 systemd 规则解析；否则解析传统 `start/running`、`stop/waiting`。

## 数据模型
- 复用 `ServiceDetail` 与 `extra` 扩展区。
- 前端 `ServiceDetail` 映射 `extra.exec` 为默认执行路径展示/编辑。

## 测试与验收
- Go：新增/修复解析单测；新增 create 流程单测；命令构造与路径探测单测。
- 前端：向导 UI 流程、模板切换、原始/表单切换单测。
- 验收：列表/详情可正确显示；新增后可启动/自启；缺失管理器提示准确。

## 风险与对策
- 非 systemd 环境：检测 init 与 LookPath 失败时返回不可用。
- sysv 详情输出差异：采用双模式解析并保留 unknown 回退。
- 文件写入权限：统一提示需要 root，并在 API 层明确错误信息。
