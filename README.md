# Traffic Replay

流量录制回放系统，纯 Go 标准库实现，零第三方依赖。

## 运行说明

```bash
cd origin
go run ./cmd/server
```

默认监听 `:8080`，可通过环境变量 `PORT` 或 `ADDR` 修改。

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/record-tasks | 创建录制任务 |
| GET | /api/record-tasks | 录制任务列表（支持 status、keyword 筛选） |
| GET | /api/record-tasks/{id} | 获取录制任务 |
| PUT | /api/record-tasks/{id} | 更新录制任务（含状态机流转） |
| DELETE | /api/record-tasks/{id} | 删除录制任务 |
| POST | /api/traffic-samples | 创建流量样本 |
| GET | /api/traffic-samples | 流量样本列表（支持 record_task_id、method、path、status_code 筛选） |
| GET | /api/traffic-samples/{id} | 获取流量样本 |
| PUT | /api/traffic-samples/{id} | 更新流量样本 |
| DELETE | /api/traffic-samples/{id} | 删除流量样本 |
| POST | /api/traffic-samples/batch | 批量创建流量样本 |
| POST | /api/traffic-samples/capture | 样本捕获（按 filter_path 过滤 + sample_rate 抽样） |
| POST | /api/replay-tasks | 创建回放任务 |
| GET | /api/replay-tasks | 回放任务列表（支持 status、keyword 筛选） |
| GET | /api/replay-tasks/{id} | 获取回放任务 |
| PUT | /api/replay-tasks/{id} | 更新回放任务（含状态机流转） |
| DELETE | /api/replay-tasks/{id} | 删除回放任务 |
| POST | /api/replay-results | 创建回放结果 |
| GET | /api/replay-results | 回放结果列表（支持 replay_task_id、matched 筛选） |
| GET | /api/replay-results/{id} | 获取回放结果 |
| PUT | /api/replay-results/{id} | 更新回放结果 |
| DELETE | /api/replay-results/{id} | 删除回放结果 |
| POST | /api/replay-results/execute | 执行回放（生成 ReplayResult，计算 matched） |
| POST | /api/replay-configs | 创建回放配置 |
| GET | /api/replay-configs | 回放配置列表（支持 enabled、keyword 筛选） |
| GET | /api/replay-configs/{id} | 获取回放配置 |
| PUT | /api/replay-configs/{id} | 更新回放配置 |
| DELETE | /api/replay-configs/{id} | 删除回放配置 |
| POST | /api/schedules | 创建调度计划 |
| GET | /api/schedules | 调度计划列表（支持 status、keyword 筛选） |
| GET | /api/schedules/{id} | 获取调度计划 |
| PUT | /api/schedules/{id} | 更新调度计划（含状态机流转 active↔paused） |
| DELETE | /api/schedules/{id} | 删除调度计划 |
| POST | /api/schedules/{id}/run | 手动运行调度计划 |
| GET | /api/stats/overview | 统计概览 |
| GET | /api/stats/record-task-status | 录制任务状态分布 |
| GET | /api/stats/replay-task-status | 回放任务状态分布 |
| GET | /api/stats/sample-count-by-task | 各录制任务的样本数 |
| GET | /api/stats/result-count-by-replay-task | 各回放任务的结果数 |
| GET | /api/stats/top-latency | 延迟最高的回放结果 |
