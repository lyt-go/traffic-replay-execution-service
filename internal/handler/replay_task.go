package handler

import (
	"net/http"

	"trafficreplay/internal/model"
	"trafficreplay/pkg/httpx"
)

func (s *Server) registerReplayTaskRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/replay-tasks", s.createReplayTask)
	mux.HandleFunc("GET /api/replay-tasks", s.listReplayTasks)
	mux.HandleFunc("GET /api/replay-tasks/{id}", s.getReplayTask)
	mux.HandleFunc("PUT /api/replay-tasks/{id}", s.updateReplayTask)
	mux.HandleFunc("DELETE /api/replay-tasks/{id}", s.deleteReplayTask)
}

type createReplayTaskRequest struct {
	Name         string `json:"name"`
	TargetURL    string `json:"target_url"`
	RecordTaskID string `json:"record_task_id"`
	Concurrency  int    `json:"concurrency"`
	TimeoutMs    int    `json:"timeout_ms"`
	SampleCount  int    `json:"sample_count"`
}

func (s *Server) createReplayTask(w http.ResponseWriter, r *http.Request) {
	var req createReplayTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateReplayTask(model.ReplayTask{
		Name:         req.Name,
		TargetURL:    req.TargetURL,
		RecordTaskID: req.RecordTaskID,
		Concurrency:  req.Concurrency,
		TimeoutMs:    req.TimeoutMs,
		SampleCount:  req.SampleCount,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listReplayTasks(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ReplayTaskFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListReplayTasks(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getReplayTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.GetReplayTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type updateReplayTaskRequest struct {
	Name         string `json:"name"`
	TargetURL    string `json:"target_url"`
	RecordTaskID string `json:"record_task_id"`
	Concurrency  int    `json:"concurrency"`
	TimeoutMs    int    `json:"timeout_ms"`
	SampleCount  int    `json:"sample_count"`
	Status       string `json:"status"`
}

func (s *Server) updateReplayTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateReplayTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateReplayTask(id, model.ReplayTask{
		Name:         req.Name,
		TargetURL:    req.TargetURL,
		RecordTaskID: req.RecordTaskID,
		Concurrency:  req.Concurrency,
		TimeoutMs:    req.TimeoutMs,
		SampleCount:  req.SampleCount,
		Status:       req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteReplayTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteReplayTask(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
