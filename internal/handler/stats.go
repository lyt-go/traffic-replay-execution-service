package handler

import (
	"net/http"
	"strconv"

	"trafficreplay/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.getOverviewStats)
	mux.HandleFunc("GET /api/stats/record-task-status", s.getRecordTaskStatusStats)
	mux.HandleFunc("GET /api/stats/replay-task-status", s.getReplayTaskStatusStats)
	mux.HandleFunc("GET /api/stats/sample-count-by-task", s.getSampleCountByTask)
	mux.HandleFunc("GET /api/stats/result-count-by-replay-task", s.getResultCountByReplayTask)
	mux.HandleFunc("GET /api/stats/top-latency", s.getTopLatencyResults)
}

func (s *Server) getOverviewStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GetOverviewStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) getRecordTaskStatusStats(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetRecordTaskStatusStats()
	httpx.OK(w, stats)
}

func (s *Server) getReplayTaskStatusStats(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetReplayTaskStatusStats()
	httpx.OK(w, stats)
}

func (s *Server) getSampleCountByTask(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetSampleCountByTask()
	httpx.OK(w, stats)
}

func (s *Server) getResultCountByReplayTask(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetResultCountByReplayTask()
	httpx.OK(w, stats)
}

func (s *Server) getTopLatencyResults(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}
	results := s.svc.GetTopLatencyResults(limit)
	httpx.OK(w, results)
}
