package handler

import (
	"net/http"
	"strconv"

	"trafficreplay/internal/model"
	"trafficreplay/pkg/httpx"
)

func (s *Server) registerReplayConfigRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/replay-configs", s.createReplayConfig)
	mux.HandleFunc("GET /api/replay-configs", s.listReplayConfigs)
	mux.HandleFunc("GET /api/replay-configs/{id}", s.getReplayConfig)
	mux.HandleFunc("PUT /api/replay-configs/{id}", s.updateReplayConfig)
	mux.HandleFunc("DELETE /api/replay-configs/{id}", s.deleteReplayConfig)
}

type createReplayConfigRequest struct {
	Name       string `json:"name"`
	TargetHost string `json:"target_host"`
	TimeoutMs  int    `json:"timeout_ms"`
	Retries    int    `json:"retries"`
	Enabled    bool   `json:"enabled"`
}

func (s *Server) createReplayConfig(w http.ResponseWriter, r *http.Request) {
	var req createReplayConfigRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateReplayConfig(model.ReplayConfig{
		Name:       req.Name,
		TargetHost: req.TargetHost,
		TimeoutMs:  req.TimeoutMs,
		Retries:    req.Retries,
		Enabled:    req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listReplayConfigs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ReplayConfigFilter{
		Keyword: r.URL.Query().Get("keyword"),
	}
	if enabledStr := r.URL.Query().Get("enabled"); enabledStr != "" {
		v, _ := strconv.ParseBool(enabledStr)
		filter.Enabled = &v
	}
	items, total, err := s.svc.ListReplayConfigs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getReplayConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.GetReplayConfig(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

type updateReplayConfigRequest struct {
	Name       string `json:"name"`
	TargetHost string `json:"target_host"`
	TimeoutMs  int    `json:"timeout_ms"`
	Retries    int    `json:"retries"`
	Enabled    bool   `json:"enabled"`
}

func (s *Server) updateReplayConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateReplayConfigRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateReplayConfig(id, model.ReplayConfig{
		Name:       req.Name,
		TargetHost: req.TargetHost,
		TimeoutMs:  req.TimeoutMs,
		Retries:    req.Retries,
		Enabled:    req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteReplayConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteReplayConfig(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
