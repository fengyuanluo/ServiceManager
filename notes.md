# Notes: ServiceManager 代码审阅

## Sources

### Local repo
- 路径: /root/Coding/ServiceManager
- 说明: 逐文件静态审阅与结构梳理

## Synthesized Findings

### 模块划分
- 后端入口与CLI：main.go、cmd/*
- 后端配置与版本：internal/conf/*
- 后端数据层：internal/db/*、internal/model/*
- 后端业务操作：internal/op/*
- 后端服务层：internal/server/**（auth/handlers/middleware/resp/router/assets）
- 后端服务管理核心：internal/service/manager/*
- 后端工具与基础设施：internal/utils/*、tools/tools.go、docs/docs.go
- 构建与安装脚本：scripts/*
- 前端应用壳：web/index.html、web/src/main.ts、web/src/App.vue
- 前端 API 与数据：web/src/api/*、web/src/data/*、web/src/types/*
- 前端组件：web/src/components/**
- 前端国际化与样式：web/src/i18n/*、web/src/styles/*
- 前端测试：web/src/__tests__/*
- 前端构建配置：web/package.json、web/vite.config.ts、web/tailwind.config.ts、web/postcss.config.cjs、web/tsconfig.json、web/pnpm-lock.yaml
- 静态资源：static/index.html、internal/server/assets/static/**（疑似构建产物）

### 范围说明
- 不审阅第三方依赖与构建产物：web/node_modules/**（第三方依赖）
- internal/server/assets/static/** 作为构建产物仅做风险提示，不逐行审阅

### 问题清单
- 后端/配置
  - internal/conf/config.go：配置键使用 logging.level，但结构体字段为 Log，且缺少 mapstructure 标注，导致 logging 可能无法反序列化到 AppConfig.Log；默认值同样失效。
- 后端/启动与路由
  - cmd/start.go：conf.Load 返回值未处理，配置读取失败仍继续启动；router.RegisterAll 返回错误未处理。
  - cmd/start.go：服务先启动再初始化 DB/User，初始化失败时服务仍运行，可能进入半可用状态。
- 后端/用户与认证
  - internal/op/user.go：首次初始化默认账号 admin/admin 并记录日志，存在严重安全风险；任意 DB 错误都会触发创建账号，可能掩盖真实异常；全局 userCache 未加锁且不从 DB 重新加载，存在并发/一致性风险。
  - internal/server/auth/auth.go：JWT 密钥由用户名+密码派生且未校验签名算法，安全边界不清晰；密码变更后所有 token 失效（可能不是预期）。
  - internal/server/middleware/auth.go：缺少 Authorization 返回 400 而非 401；未严格校验 Bearer 格式。
- 后端/中间件
  - internal/server/middleware/cors.go：AllowAllOrigins=true 与 AllowCredentials=true 组合不合法且存在安全隐患。
  - internal/server/middleware/static.go：fileSystem.Open 未关闭；文件不存在时未调用 c.Next() 也未 fallback，可能导致空响应或 SPA 路由不可用。
- 后端/服务管理核心
  - internal/service/manager/service.go：pm2 使用 su -c 拼接命令字符串，服务名/参数来自用户输入时存在命令注入风险。
  - internal/service/manager/create.go：service 名称未校验，可能路径穿越写入任意文件；systemd/env 等文本未做换行/特殊字符约束，存在配置注入风险。
  - internal/service/manager/registry.go：DetectManagers 在 distro 非 unknown 时总返回 sysv，可能误报可用管理器。
- 后端/日志
  - internal/utils/log/log.go：日志级别未读取配置，固定 Info。
- 构建与发布
  - scripts/build.sh：LDFLAGS 使用 example/internal/conf，与实际模块路径不一致，版本注入失效；APP_NAME 与项目命名不一致可能导致产物/安装脚本不匹配。
  - go.mod：go 版本 1.24.4 可能与实际工具链不兼容（需确认）。
- 前端/功能完整性
  - web/src/components/manager/ServiceEditor.vue：保存/重置按钮无事件处理，功能空缺。
  - web/src/components/manager/ServiceCreateWizard.vue：steps/activeStep 逻辑未实现，步骤显示恒定。
  - web/src/data/services.ts：疑似 mock 数据，若未使用建议清理或明确用途。
- 前端/接口耦合
  - web/src/api/client.ts：严格依赖后端 code=200 结构，若接口返回非标准结构即抛错；对 204/非 JSON 场景兼容性不足。
