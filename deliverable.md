# Deliverable: 服务新增向导与批量工具栏完善

## 交付内容
- 新增“悬浮对话框”新增向导，支持表单/原始切换、模板提示与选项设置。
- 批量工具栏新增入口与交互，接入创建 API 并刷新列表。
- 中英双语文案补齐，前端静态构建输出到 static。

## 验证与构建
- pnpm -C web test
- go test ./...
- pnpm -C web build
- go build ./...
