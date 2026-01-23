# Notes: 服务新增向导与批量工具栏完善

## Sources

### 现有代码与计划
- URL: 本地仓库
- Key points:
  - 现有后端已具备 /api/v1/{manager}/create 创建接口与模板能力。
  - 前端已有新增向导测试文件，但未实现 UI 与 API 接入。

## Synthesized Findings

### 前端设计约束
- 极简但精致，避免冗余文字。
- 悬浮对话框替代右侧抽屉；批量按钮同一行。
- 响应式与双语支持。

### 技术约束
- Vue3 + Vite + PNPM + shadcn + Tailwind。
- TDD：先写失败测试再实现。
- 每模块开始前至少 10 次上下文获取。

## Context Reads (>=10)
1. web/src/components/manager/BatchToolbar.vue
2. web/src/components/manager/ManagerPanel.vue
3. web/src/components/manager/ServiceEditor.vue
4. web/src/__tests__/service-create-wizard.test.ts
5. web/src/api/services.ts
6. web/src/types/service.ts
7. web/src/i18n/messages.ts
8. web/src/components/manager/ServiceList.vue
9. web/src/components/manager/LogViewer.vue
10. web/src/components/layout/MobileDrawer.vue
11. web/src/components/layout/TopBar.vue
12. web/package.json
13. web/src/__tests__/api-client.test.ts
14. web/src/__tests__/manager-panel.test.ts

## MCP
- Context7 resolve-library-id 失败：TypeError: fetch failed（待按离线经验继续）
