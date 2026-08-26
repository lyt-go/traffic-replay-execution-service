package handler

import (
	"net/http"
	"strconv"

	"trafficreplay/internal/model"
	"trafficreplay/pkg/httpx"
)

func (s *Server) registerTrafficSampleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/traffic-samples", s.createTrafficSample)
	mux.HandleFunc("GET /api/traffic-samples", s.listTrafficSamples)
	mux.HandleFunc("GET /api/traffic-samples/{id}", s.getTrafficSample)
	mux.HandleFunc("PUT /api/traffic-samples/{id}", s.updateTrafficSample)
	mux.HandleFunc("DELETE /api/traffic-samples/{id}", s.deleteTrafficSample)
	mux.HandleFunc("POST /api/traffic-samples/batch", s.batchCreateTrafficSamples)
	mux.HandleFunc("POST /api/traffic-samples/capture", s.captureTraffic)
}

type createTrafficSampleRequest struct {
	RecordTaskID string `json:"record_task_id"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	Headers      string `json:"headers"`
	Body         string `json:"body"`
	StatusCode   int    `json:"status_code"`
}

func (s *Server) createTrafficSample(w http.ResponseWriter, r *http.Request) {
	var req createTrafficSampleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sample, err := s.svc.CreateTrafficSample(model.TrafficSample{
		RecordTaskID: req.RecordTaskID,
		Method:       req.Method,
		Path:         req.Path,
		Headers:      req.Headers,
		Body:         req.Body,
		StatusCode:   req.StatusCode,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sample)
}

func (s *Server) listTrafficSamples(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	statusCode, _ := strconv.Atoi(r.URL.Query().Get("status_code"))
	filter := model.TrafficSampleFilter{
		RecordTaskID: r.URL.Query().Get("record_task_id"),
		Method:       r.URL.Query().Get("method"),
		Path:         r.URL.Query().Get("path"),
		StatusCode:   statusCode,
	}
	items, total, err := s.svc.ListTrafficSamples(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTrafficSample(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sample, err := s.svc.GetTrafficSample(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sample)
}

type updateTrafficSampleRequest struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	Headers    string `json:"headers"`
	Body       string `json:"body"`
	StatusCode int    `json:"status_code"`
}

func (s *Server) updateTrafficSample(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateTrafficSampleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sample, err := s.svc.UpdateTrafficSample(id, model.TrafficSample{
		Method:     req.Method,
		Path:       req.Path,
		Headers:    req.Headers,
		Body:       req.Body,
		StatusCode: req.StatusCode,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sample)
}

func (s *Server) deleteTrafficSample(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteTrafficSample(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type batchCreateTrafficSamplesRequest struct {
	Samples []createTrafficSampleRequest `json:"samples"`
}

func (s *Server) batchCreateTrafficSamples(w http.ResponseWriter, r *http.Request) {
	var req batchCreateTrafficSamplesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inputs := make([]model.TrafficSample, 0, len(req.Samples))
	for _, sampleReq := range req.Samples {
		inputs = append(inputs, model.TrafficSample{
			RecordTaskID: sampleReq.RecordTaskID,
			Method:       sampleReq.Method,
			Path:         sampleReq.Path,
			Headers:      sampleReq.Headers,
			Body:         sampleReq.Body,
			StatusCode:   sampleReq.StatusCode,
		})
	}
	samples, err := s.svc.BatchCreateTrafficSamples(inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, samples)
}

type captureTrafficRequest struct {
	RecordTaskID string `json:"record_task_id"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	Headers      string `json:"headers"`
	Body         string `json:"body"`
	StatusCode   int    `json:"status_code"`
}

func (s *Server) captureTraffic(w http.ResponseWriter, r *http.Request) {
	var req captureTrafficRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sample, err := s.svc.CaptureTraffic(req.RecordTaskID, model.TrafficSample{
		Method:     req.Method,
		Path:       req.Path,
		Headers:    req.Headers,
		Body:       req.Body,
		StatusCode: req.StatusCode,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if sample == nil {
		httpx.OK(w, map[string]string{"message": "sample skipped by filter or rate"})
		return
	}
	httpx.Created(w, sample)
}
