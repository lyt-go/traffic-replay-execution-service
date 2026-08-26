package handler

import (
	"net/http"

	"trafficreplay/internal/model"
	"trafficreplay/pkg/httpx"
)

func (s *Server) registerRecordTaskRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/record-tasks", s.createRecordTask)
	mux.HandleFunc("GET /api/record-tasks", s.listRecordTasks)
	mux.HandleFunc("GET /api/record-tasks/{id}", s.getRecordTask)
	mux.HandleFunc("PUT /api/record-tasks/{id}", s.updateRecordTask)
	mux.HandleFunc("DELETE /api/record-tasks/{id}", s.deleteRecordTask)
}

type createRecordTaskRequest struct {
	Name       string `json:"name"`
	SourceURL  string `json:"source_url"`
	FilterPath string `json:"filter_path"`
	SampleRate int    `json:"sample_rate"`
}

func (s *Server) createRecordTask(w http.ResponseWriter, r *http.Request) {
	var req createRecordTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateRecordTask(model.RecordTask{
		Name:       req.Name,
		SourceURL:  req.SourceURL,
		FilterPath: req.FilterPath,
		SampleRate: req.SampleRate,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listRecordTasks(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RecordTaskFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListRecordTasks(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRecordTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.GetRecordTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type updateRecordTaskRequest struct {
	Name       string `json:"name"`
	SourceURL  string `json:"source_url"`
	FilterPath string `json:"filter_path"`
	SampleRate int    `json:"sample_rate"`
	Status     string `json:"status"`
}

func (s *Server) updateRecordTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateRecordTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateRecordTask(id, model.RecordTask{
		Name:       req.Name,
		SourceURL:  req.SourceURL,
		FilterPath: req.FilterPath,
		SampleRate: req.SampleRate,
		Status:     req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteRecordTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRecordTask(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
