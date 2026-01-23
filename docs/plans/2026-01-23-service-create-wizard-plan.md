# 服务新增向导 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan.

**Goal:** 完成新增向导的悬浮对话框 UI、API 接入与测试，并通过深度测试与构建。

**Architecture:** 前端以 ManagerPanel 为入口打开 ServiceCreateWizard；向导提交后调用 createService，再刷新列表并选中新服务。UI 与现有卡片风格一致，支持响应式与双语。

**Tech Stack:** Vue 3, Vite, Vitest, Tailwind CSS

### Task 1: API 创建接口与类型

**Files:**
- Modify: web/src/api/services.ts
- Modify: web/src/types/service.ts
- Test: web/src/__tests__/api-client.test.ts 或新增 web/src/__tests__/services-create.test.ts

**Step 1: Write the failing test**
编写测试断言 createService 发起 POST 并携带 payload。

**Step 2: Run test to verify it fails**
Run: pnpm -C web test -- --run api-client
Expected: FAIL (createService 未定义或未调用)

**Step 3: Write minimal implementation**
在 types 与 api 中补齐 CreateRequest 与 createService。

**Step 4: Run test to verify it passes**
Run: pnpm -C web test -- --run api-client
Expected: PASS

**Step 5: Commit**
Commit: feat(web): add create service api client

### Task 2: 新增向导组件

**Files:**
- Create: web/src/components/manager/ServiceCreateWizard.vue
- Modify: web/src/components/manager/BatchToolbar.vue
- Modify: web/src/i18n/messages.ts
- Test: web/src/__tests__/service-create-wizard.test.ts

**Step 1: Write the failing test**
在测试中断言点击新增按钮后出现向导。

**Step 2: Run test to verify it fails**
Run: pnpm -C web test -- --run service-create-wizard
Expected: FAIL (create-service 按钮或向导不存在)

**Step 3: Write minimal implementation**
实现新增按钮与 ServiceCreateWizard 基础结构（对话框 + 关闭）。

**Step 4: Run test to verify it passes**
Run: pnpm -C web test -- --run service-create-wizard
Expected: PASS

**Step 5: Commit**
Commit: feat(web): add service create wizard shell

### Task 3: 向导表单与创建流程

**Files:**
- Modify: web/src/components/manager/ServiceCreateWizard.vue
- Modify: web/src/components/manager/ManagerPanel.vue
- Modify: web/src/api/services.ts
- Modify: web/src/i18n/messages.ts
- Test: web/src/__tests__/service-create-wizard.test.ts

**Step 1: Write the failing test**
补充测试：填写表单并提交后调用 createService，刷新列表。

**Step 2: Run test to verify it fails**
Run: pnpm -C web test -- --run service-create-wizard
Expected: FAIL (createService 未被调用)

**Step 3: Write minimal implementation**
向导内实现步骤与表单，提交触发 createService 并通知 ManagerPanel 刷新。

**Step 4: Run test to verify it passes**
Run: pnpm -C web test -- --run service-create-wizard
Expected: PASS

**Step 5: Commit**
Commit: feat(web): wire create service flow

### Task 4: 视觉与响应式微调

**Files:**
- Modify: web/src/components/manager/ServiceCreateWizard.vue
- Modify: web/src/components/manager/BatchToolbar.vue

**Step 1: Visual polish**
完善排版、间距、阴影与移动端布局。

**Step 2: Run targeted tests**
Run: pnpm -C web test -- --run service-create-wizard
Expected: PASS

**Step 3: Commit**
Commit: style(web): polish create wizard layout

### Task 5: 深度测试与构建

**Files:**
- None (verification)

**Step 1: Run frontend tests**
Run: pnpm -C web test
Expected: PASS

**Step 2: Run backend tests**
Run: go test ./...
Expected: PASS

**Step 3: Build frontend**
Run: pnpm -C web build
Expected: build success

**Step 4: Build backend**
Run: go build ./...
Expected: build success

**Step 5: Commit**
Commit: chore: verify create wizard build
