# Task Plan: 服务新增向导与批量工具栏完善

## Goal
完成前端新增向导与批量工具栏的交互与视觉实现，接入后端创建 API，并通过深度测试与构建验收。

## Phases
- [x] Phase 1: 计划与环境确认
- [x] Phase 2: 研究与上下文收集
- [x] Phase 3: 实现与测试（前端为主）
- [x] Phase 4: 深度测试、构建与交付

## Key Questions
1. 新增向导的步骤结构与字段映射如何与现有服务编辑器复用？
2. 批量工具栏与新增入口的布局如何在各断点保持精致与易用？
3. 创建 API 的请求模型与前端类型如何对齐？

## Decisions Made
- 采用悬浮对话框形式的新增向导，与批量工具栏同一行显示新增入口。
- 通过当前管理器上下文创建服务，不额外选择管理器。

## Errors Encountered
- rg 命令不可用，已改用 grep 搜索。
- Context7 resolve-library-id 失败（TypeError: fetch failed），将使用本地经验作为降级。
- 设计文档首次写入因反引号被 shell 解析失败，已改为纯文本写入。
- python 命令不可用，改用 cat 写入测试文件。
- ServiceCreateWizard 模板字符串被错误写入导致渲染报错，已修复引号缺失问题。

## Status
**Currently in Phase 4** - 深度测试与构建完成，准备提交与交付
