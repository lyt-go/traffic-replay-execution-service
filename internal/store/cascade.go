package store

// 本文件集中实现删除时的级联清理。
//
// 约定：所有 *Locked 方法均假设调用方已持有 s.mu 写锁，方法自身不再加锁，
// 否则在递归调用 / 串行调用时会触发死锁。Go 规范允许在 range 遍历 map 的
// 同时 delete 当前或其它 key，迭代顺序未定义但不会漏删或 panic。

// deleteTrafficSamplesByTaskLocked 删除属于指定录制任务的全部流量样本，
// 并级联删除依赖这些样本的回放结果。调用方必须已持有 s.mu 写锁。
func (s *MemoryStore) deleteTrafficSamplesByTaskLocked(recordTaskID string) {
	for sid, sample := range s.trafficSamples {
		if sample.RecordTaskID == recordTaskID {
			delete(s.trafficSamples, sid)
			s.deleteReplayResultsBySampleLocked(sid)
		}
	}
}

// deleteReplayResultsByTaskLocked 删除属于指定回放任务的全部回放结果。
// 调用方必须已持有 s.mu 写锁。
func (s *MemoryStore) deleteReplayResultsByTaskLocked(replayTaskID string) {
	for rid, r := range s.replayResults {
		if r.ReplayTaskID == replayTaskID {
			delete(s.replayResults, rid)
		}
	}
}

// deleteReplayResultsBySampleLocked 删除依赖指定流量样本的全部回放结果。
// 调用方必须已持有 s.mu 写锁。
func (s *MemoryStore) deleteReplayResultsBySampleLocked(sampleID string) {
	for rid, r := range s.replayResults {
		if r.SampleID == sampleID {
			delete(s.replayResults, rid)
		}
	}
}

// deleteSchedulesByConfigLocked 删除依赖指定回放配置的全部调度计划。
// 调用方必须已持有 s.mu 写锁。
func (s *MemoryStore) deleteSchedulesByConfigLocked(configID string) {
	for sid, sch := range s.schedules {
		if sch.ConfigID == configID {
			delete(s.schedules, sid)
		}
	}
}
