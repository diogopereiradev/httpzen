package benchmark_module

import (
	"sync"
	"testing"
	"time"

	request_module "github.com/diogopereiradev/httpzen/internal/request"
)

func TestInitialResultModel(t *testing.T) {
	result := initialResultModel()
	if result.Metrics.TotalRequests != 0 {
		t.Errorf("Expected TotalRequests 0, got %d", result.Metrics.TotalRequests)
	}
	if result.Metrics.RequestsMinLatency != 0 {
		t.Errorf("Expected RequestsMinLatency 0, got %f", result.Metrics.RequestsMinLatency)
	}
}

func TestRunBenchmark(t *testing.T) {
	origRunRequest := request_module.RunRequest
	request_module.RunRequest = func(options request_module.RequestOptions) request_module.RequestResponse {
		return request_module.RequestResponse{
			StatusCode:    200,
			ExecutionTime: 10,
			Result:        "response body",
		}
	}
	defer func() { request_module.RunRequest = origRunRequest }()

	metrics := &Metrics{}
	options := BenchmarkOptions{
		Request: request_module.RequestOptions{
			Method: "GET",
			Url:    "https://google.com",
		},
		ThreadsAmount: 1,
		Duration:      1,
	}
	RunBenchmark(options, metrics)
	if metrics.TotalRequests == 0 {
		t.Errorf("Expected TotalRequests > 0, got %d", metrics.TotalRequests)
	}
	if metrics.TotalSuccess == 0 {
		t.Errorf("Expected TotalSuccess > 0, got %d", metrics.TotalSuccess)
	}
	if metrics.RequestsMinLatency == 0 {
		t.Errorf("Expected RequestsMinLatency > 0, got %f", metrics.RequestsMinLatency)
	}
	if metrics.RequestsMaxLatency == 0 {
		t.Errorf("Expected RequestsMaxLatency > 0, got %f", metrics.RequestsMaxLatency)
	}
}

func TestDoRequestSuccessAndError(t *testing.T) {
	model := initialResultModel()
	options := BenchmarkOptions{
		Request: request_module.RequestOptions{
			Method: "GET",
			Url:    "https://google.com",
		},
	}

	origRunRequest := request_module.RunRequest
	request_module.RunRequest = func(options request_module.RequestOptions) request_module.RequestResponse {
		return request_module.RequestResponse{StatusCode: 200, ExecutionTime: 5, Result: "ok"}
	}
	resp := options.doRequest(model)
	if resp.StatusCode != 200 {
		t.Errorf("Expected StatusCode 200, got %d", resp.StatusCode)
	}
	if model.Metrics.TotalSuccess != 1 {
		t.Errorf("Expected TotalSuccess 1, got %d", model.Metrics.TotalSuccess)
	}

	request_module.RunRequest = func(options request_module.RequestOptions) request_module.RequestResponse {
		return request_module.RequestResponse{StatusCode: 500, ExecutionTime: 5, Result: "fail"}
	}
	resp = options.doRequest(model)
	if resp.StatusCode != 500 {
		t.Errorf("Expected StatusCode 500, got %d", resp.StatusCode)
	}
	if model.Metrics.TotalErrors != 1 {
		t.Errorf("Expected TotalErrors 1, got %d", model.Metrics.TotalErrors)
	}
	request_module.RunRequest = origRunRequest
}

func TestRunThreads(t *testing.T) {
	origRunRequest := request_module.RunRequest
	request_module.RunRequest = func(options request_module.RequestOptions) request_module.RequestResponse {
		time.Sleep(10 * time.Millisecond)
		return request_module.RequestResponse{StatusCode: 200, ExecutionTime: 1, Result: "ok"}
	}
	defer func() { request_module.RunRequest = origRunRequest }()

	metrics := &Metrics{}
	options := BenchmarkOptions{
		Request: request_module.RequestOptions{
			Method: "GET",
			Url:    "https://google.com",
		},
		ThreadsAmount: 2,
		Duration:      1,
	}
	options.runThreads(initialResultModel(), metrics)
	if metrics.TotalRequests == 0 {
		t.Errorf("Expected TotalRequests > 0, got %d", metrics.TotalRequests)
	}
}

func TestBenchmarkResultMutex(t *testing.T) {
	model := initialResultModel()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		model.mutex.Lock()
		model.Metrics.TotalRequests++
		model.mutex.Unlock()
		wg.Done()
	}()
	go func() {
		model.mutex.Lock()
		model.Metrics.TotalRequests++
		model.mutex.Unlock()
		wg.Done()
	}()
	wg.Wait()
	if model.Metrics.TotalRequests != 2 {
		t.Errorf("Expected TotalRequests 2, got %d", model.Metrics.TotalRequests)
	}
}
