# Task Plan: 逐项修复审阅问题

## Goal
按 report.md 顺序建立任务列表并逐项高质量修复，确保测试通过。

## Phases
- [x] Phase 1: 规划与准备（确认范围、创建计划/笔记文件）
- [x] Phase 2: 任务分解（按模块拆分修复清单与验证点）
- [x] Phase 3: 实施修复（逐项修复 + 测试）
- [x] Phase 4: 回顾交付（汇总变更与验证结果）

## Key Questions
1. 哪些修复需要 API/前端联动调整以保证行为一致？
2. 是否需要在修复过程中补充/调整测试？

## Decisions Made
- 按 report.md 顺序推进修复，逐项完成并验证。

## Errors Encountered
- 使用 here-doc 写入包含反引号内容时触发命令替换，已改用 apply_patch 重新生成设计文档
- 静态资源 fallback 使用 /index.html 触发 301 重定向，改为指向目录路径以返回 index 内容
- create.go 引入 sanitizeLine 后出现未使用变量，已清理后通过构建

## Status
**Currently in Phase 4** - 已完成全部修复与回归验证，整理交付说明
