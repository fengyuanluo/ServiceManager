# 新增服务向导与解析修复 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan.

**Goal:** 修复 systemctl/sysv 解析问题，并为全部服务管理器补齐新增服务功能（后端+前端+测试）。

**Architecture:** 后端在 manager 包内新增 create 入口与模板渲染器，统一通过 ExecSource/FileSource 执行与落盘；前端增加悬浮新增向导并复用表单/原始切换。解析修复在 systemd/sysv parser 内完成，返回统一字段 + extra 扩展。

**Tech Stack:** Go（Gin）、Vue3 + Vite + pnpm、Tailwind、Vitest、Go test

---

### Task 1: systemctl 路径解析与执行探测

**Files:**
- Modify: `internal/service/manager/exec.go`
- Test: `internal/service/manager/exec_test.go`

**Step 1: Write the failing test**

```go
func TestResolveCommandFallback(t *testing.T) {
    got, err := resolveCommand("systemctl", func(name string) (string, error) {
        return "", errors.New("not found")
    }, []string{"/bin/systemctl"})
    if err != nil || got != "/bin/systemctl" {
        t.Fatalf("expected fallback path")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/service/manager -run TestResolveCommandFallback`
Expected: FAIL with "undefined: resolveCommand"

**Step 3: Write minimal implementation**

```go
func resolveCommand(name string, look func(string) (string, error), fallbacks []string) (string, error) {
    if path, err := look(name); err == nil {
        return path, nil
    }
    for _, candidate := range fallbacks {
        if _, err := os.Stat(candidate); err == nil {
            return candidate, nil
        }
    }
    return "", errCommandNotFound
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/service/manager -run TestResolveCommandFallback`
Expected: PASS

**Step 5: Integrate into RunCommand**

- Run `resolveCommand` only when `src.Run == nil`。
- 避免影响现有单测（mock Run 不做解析）。

Run: `go test ./internal/service/manager -run TestRunCommand`
Expected: PASS

---

### Task 2: systemd 详情补充 ExecStart / UnitFileState

**Files:**
- Modify: `internal/service/manager/service.go`
- Test: `internal/service/manager/systemd_test.go`

**Step 1: Write the failing test**

```go
func TestParseSystemdDetailExecStart(t *testing.T) {
    out := "Id=cron.service\nExecStart=/usr/sbin/cron -f\nUnitFileState=enabled\nActiveState=active"
    detail := parseSystemdDetail(out, "cron")
    if detail.Extra["exec"] != "/usr/sbin/cron -f" || detail.Extra["autostart"] != "enabled" {
        t.Fatalf("missing exec/autostart")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/service/manager -run TestParseSystemdDetailExecStart`
Expected: FAIL with missing fields

**Step 3: Write minimal implementation**

- `buildDetailCommand` 增加 `ExecStart`、`UnitFileState` 属性。
- `parseSystemdDetail` 写入 `detail.Extra["exec"]` 与 `detail.Extra["autostart"]`（保留原始字符串）。

**Step 4: Run test to verify it passes**

Run: `go test ./internal/service/manager -run TestParseSystemdDetailExecStart`
Expected: PASS

---

### Task 3: sysv 详情独立解析（兼容 systemd 风格输出）

**Files:**
- Modify: `internal/service/manager/sysv.go`
- Modify: `internal/service/manager/service.go`
- Test: `internal/service/manager/sysv_test.go`

**Step 1: Write the failing test**

```go
func TestParseSysVDetailSystemdStyle(t *testing.T) {
    out := "● cron.service - Regular background program\\n Active: active (running) since Thu\\n Main PID: 872 (cron)\\n Loaded: loaded (/lib/systemd/system/cron.service; enabled)"
    detail := ParseSysVDetail(out, "cron")
    if detail.Status != ServiceStatus(\"running\") || detail.PID != \"872\" {
        t.Fatalf(\"unexpected detail parse\")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/service/manager -run TestParseSysVDetailSystemdStyle`
Expected: FAIL with "undefined: ParseSysVDetail"

**Step 3: Write minimal implementation**

- 新增 `ParseSysVDetail(out, fallbackName)`：
  - 含 `Active:` 时解析状态
  - `Main PID:` 解析 PID
  - `Loaded:` 解析 unit path
  - `● xxx - desc` 解析描述
- `GetServiceDetail` 的 sysv 分支改为使用 `ParseSysVDetail`。

**Step 4: Run test to verify it passes**

Run: `go test ./internal/service/manager -run TestParseSysVDetailSystemdStyle`
Expected: PASS

---

### Task 4: 后端新增服务创建能力（模板 + 落盘 + 执行）

**Files:**
- Create: `internal/service/manager/create.go`
- Modify: `internal/service/manager/types.go`
- Test: `internal/service/manager/create_test.go`

**Step 1: Write the failing test**

```go
func TestCreateSystemdService(t *testing.T) {
    dir := t.TempDir()
    calls := []string{}
    src := ExecSource{Run: func(name string, args ...string) ([]byte, error) {
        calls = append(calls, name+\" \"+strings.Join(args, \" \"))
        return []byte(\"ok\"), nil
    }}
    fs := FileSource{WriteFile: os.WriteFile, MkdirAll: os.MkdirAll}
    req := CreateRequest{Name: \"demo\", Exec: \"/usr/bin/demo\", TemplateID: \"simple\", Options: CreateOptions{Start: true}}
    err := CreateService(\"systemd\", CreateSource{Exec: src, FS: fs, BaseDir: dir}, req)
    if err != nil { t.Fatalf(\"create failed\") }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/service/manager -run TestCreateSystemdService`
Expected: FAIL with undefined types/functions

**Step 3: Write minimal implementation**

- 新增类型：`CreateRequest`、`CreateOptions`、`CreateSource`、`FileSource`。
- `CreateService`：
  - 选择模板渲染字符串
  - 写入对应路径（BaseDir 允许覆盖，默认系统路径）
  - systemd `daemon-reload` + 可选 `enable`/`start`
- 为 sysv/openrc/supervisor/pm2 添加基础模板与落盘逻辑。

**Step 4: Run test to verify it passes**

Run: `go test ./internal/service/manager -run TestCreateSystemdService`
Expected: PASS

---

### Task 5: 新增创建接口与路由

**Files:**
- Modify: `internal/server/handlers/manager_common.go`
- Modify: `internal/server/handlers/manager_common_test.go`

**Step 1: Write the failing test**

```go
func TestCreateHandler(t *testing.T) {
    prevExec := managerExecSource
    managerExecSource = manager.ExecSource{Run: func(name string, args ...string) ([]byte, error) {
        return []byte(\"ok\"), nil
    }}
    defer func() { managerExecSource = prevExec }()
    body := bytes.NewBufferString(`{\"name\":\"demo\",\"template_id\":\"simple\",\"exec\":\"/bin/demo\"}`)
    rec := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(rec)
    c.Request = httptest.NewRequest(http.MethodPost, \"/\", body)
    c.Request.Header.Set(\"Content-Type\", \"application/json\")
    createHandler(\"systemd\")(c)
    if rec.Code != http.StatusOK { t.Fatalf(\"unexpected status\") }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/server/handlers -run TestCreateHandler`
Expected: FAIL with undefined handler

**Step 3: Write minimal implementation**

- 在 `registerManagerRoutes` 添加 `/create` 路由。
- 新增 `createHandler`：解析 JSON，调用 `manager.CreateService` 并返回结果。

**Step 4: Run test to verify it passes**

Run: `go test ./internal/server/handlers -run TestCreateHandler`
Expected: PASS

---

### Task 6: 前端新增向导 UI + 交互（使用 frontend-design）

**Files:**
- Create: `web/src/components/manager/ServiceCreateWizard.vue`
- Modify: `web/src/components/manager/BatchToolbar.vue`
- Modify: `web/src/components/manager/ManagerPanel.vue`
- Modify: `web/src/i18n/messages.ts`
- Test: `web/src/__tests__/service-create-wizard.test.ts`

**Step 1: Write the failing test**

```ts
it(\"opens wizard from batch toolbar\", async () => {
  render(ManagerPanel, { props: { managerKey: \"systemd\", lang: \"zh\" } });
  await fireEvent.click(screen.getByTestId(\"create-service\"));
  expect(screen.getByTestId(\"create-wizard\")).toBeInTheDocument();
});
```

**Step 2: Run test to verify it fails**

Run: `pnpm -C web test -- --run service-create-wizard`
Expected: FAIL with missing component/test ids

**Step 3: Write minimal implementation**

- 在 BatchToolbar 添加“新增”按钮（`data-testid="create-service"`）。
- 新增悬浮对话框组件（模板选择 → 表单配置 → 确认创建）。
- ManagerPanel 控制弹层开关与数据提交。
- i18n 增加新增向导文案。

**Step 4: Run test to verify it passes**

Run: `pnpm -C web test -- --run service-create-wizard`
Expected: PASS

---

### Task 7: 前端 API 接入与提交逻辑

**Files:**
- Modify: `web/src/api/services.ts`
- Modify: `web/src/types/service.ts`
- Test: `web/src/__tests__/api-client.test.ts`

**Step 1: Write the failing test**

```ts
it(\"creates service\", async () => {
  mockFetchOnce({ code: 200, message: \"success\", data: { name: \"demo\", status: \"running\" } });
  await createService(\"systemd\", { name: \"demo\", templateId: \"simple\" });
  expect(fetch).toHaveBeenCalled();
});
```

**Step 2: Run test to verify it fails**

Run: `pnpm -C web test -- --run api-client`
Expected: FAIL with missing createService

**Step 3: Write minimal implementation**

- 新增 `createService` API 方法。\n- Wizard 提交后刷新列表并选中新建服务。

**Step 4: Run test to verify it passes**

Run: `pnpm -C web test -- --run api-client`
Expected: PASS

---

### Task 8: 全量验证与提交

**Step 1: Run backend tests**

Run: `go test ./...`
Expected: PASS

**Step 2: Run frontend tests**

Run: `pnpm -C web test`
Expected: PASS

**Step 3: Commit**

```bash
git add internal web docs/plans
git commit -m "feat: add service creation wizard and fix manager parsing"
```
