package handler

import (
	"net/http"

	"trafficreplay/internal/model"
	"trafficreplay/pkg/httpx"
)

func (s *Server) registerReplayResultRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/replay-results", s.createReplayResult)
	mux.HandleFunc("GET /api/replay-results", s.listReplayResults)
	mux.HandleFunc("GET /api/replay-results/{id}", s.getReplayResult)
	mux.HandleFunc("PUT /api/replay-results/{id}", s.updateReplayResult)
	mux.HandleFunc("DELETE /api/replay-results/{id}", s.deleteReplayResult)
	mux.HandleFunc("POST /api/replay-results/execute", s.executeReplay)
}

type createReplayResultRequest struct {
	ReplayTaskID   string `json:"replay_task_id"`
	SampleID       string `json:"sample_id"`
	ResponseStatus int    `json:"response_status"`
	LatencyMs      int    `json:"latency_ms"`
	Matched        bool   `json:"matched"`
	Diff           string `json:"diff"`
}

func (s *Server) createReplayResult(w http.ResponseWriter, r *http.Request) {
	var req createReplayResultRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.CreateReplayResult(model.ReplayResult{
		ReplayTaskID:   req.ReplayTaskID,
		SampleID:       req.SampleID,
		ResponseStatus: req.ResponseStatus,
		LatencyMs:      req.LatencyMs,
		Matched:        req.Matched,
		Diff:           req.Diff,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, result)
}

func (s *Server) listReplayResults(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	matchedStr := r.URL.Query().Get("matched")
	filter := model.ReplayResultFilter{
		ReplayTaskID: r.URL.Query().Get("replay_task_id"),
	}
	if matchedStr == "true" {
		v := true
		filter.Matched = &v
	} else if matchedStr == "false" {
		v := false
		filter.Matched = &v
	}
	items, total, err := s.svc.ListReplayResults(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getReplayResult(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := s.svc.GetReplayResult(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

type updateReplayResultRequest struct {
	ResponseStatus int    `json:"response_status"`
	LatencyMs      int    `json:"latency_ms"`
	Matched        bool   `json:"matched"`
	Diff           string `json:"diff"`
}

func (s *Server) updateReplayResult(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateReplayResultRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.UpdateReplayResult(id, model.ReplayResult{
		ResponseStatus: req.ResponseStatus,
		LatencyMs:      req.LatencyMs,
		Matched:        req.Matched,
		Diff:           req.Diff,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) deleteReplayResult(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteReplayResult(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type executeReplayRequest struct {
	ReplayTaskID string `json:"replay_task_id"`
}

func (s *Server) executeReplay(w http.ResponseWriter, r *http.Request) {
	var req executeReplayRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	results, err := s.svc.ExecuteReplay(req.ReplayTaskID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, results)
}
