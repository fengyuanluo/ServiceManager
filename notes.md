# Notes: 逐项修复审阅问题

## Sources

### Local report
- 路径: /root/Coding/ServiceManager/report.md
- 说明: 按报告顺序拆解修复任务

## Synthesized Findings

### 修复清单（按报告顺序）
- 入口与启动
  - cmd/start.go：处理 conf.Load 与 router.RegisterAll 错误；调整初始化顺序与失败退出策略。
- 配置与日志
  - internal/conf/config.go：修正 logging 配置键映射；补充 mapstructure 标签；默认值一致。
  - internal/utils/log/log.go：支持从配置读取日志级别。
- 用户与认证
  - internal/op/user.go：仅在记录不存在时初始化默认用户；移除明文密码日志；并发安全/一致性策略。
  - internal/server/auth/auth.go：校验签名算法；使用独立密钥或配置密钥；token 失效策略确认。
  - internal/server/middleware/auth.go：401 语义；严格解析 Bearer。
  - internal/server/handlers/user.go：区分 400/401/500 错误码与消息。
- 中间件
  - internal/server/middleware/cors.go：安全允许源策略；credentials 合法配置。
  - internal/server/middleware/static.go：关闭文件句柄；缺失文件走 Next()/SPA fallback。
- 服务管理核心
  - internal/service/manager/service.go：pm2 以参数安全传递，避免 su -c 字符串拼接。
  - internal/service/manager/create.go：服务名校验防路径穿越；过滤换行/注入字符。
  - internal/service/manager/registry.go：DetectManagers 逻辑更精准。
- 构建与发布
  - scripts/build.sh：修正 LDFLAGS 模块路径；统一 APP_NAME。
  - scripts/dockerfile/entrypoint.sh：启动命令与实际 CLI 对齐。
- 前端
  - web/src/components/manager/ServiceEditor.vue：保存/重置闭环或禁用按钮。
  - web/src/components/manager/ServiceCreateWizard.vue：补齐 steps/activeStep 逻辑。
  - web/src/api/client.ts：兼容 204/非 JSON 响应。
  - web/src/data/services.ts：标记或清理 mock 数据。

### 基线验证
- 后端：go test ./...
- 前端：pnpm -C web test（有预期的 ApiError 日志输出）

### 进度记录
- Task 5 完成：JWT 改为配置密钥 + 校验签名算法，缺失密钥时自动生成并尝试持久化。
- Task 6 完成：Auth 中间件严格解析 Bearer，缺失/非法授权返回 401；用户接口区分 400/401/500。
- Task 7 完成：CORS 支持环境变量白名单并启用凭证；静态资源关闭句柄并支持 SPA fallback。
- Task 8 完成：pm2 参数安全转义；服务名校验 + 字段去除换行；DetectManagers 更精准。
- Task 9 完成：构建脚本 LDFLAGS 与 APP_NAME 对齐；容器 entrypoint 对齐正确命令。
- Task 10 完成：前端新增确认步骤；编辑重置/保存状态完善；API 客户端兼容空响应并清理 mock 数据。
- Task 11 完成：全量回归（go test ./...、pnpm -C web test）通过。
