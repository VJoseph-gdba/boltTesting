package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
	"networkmonitor/shared"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Monitor handles HTTP request monitoring
type Monitor struct {
	client       *http.Client
	targets      []shared.Target
	resultChan   chan shared.NetworkRequest
	stopChan     chan struct{}
	wg           sync.WaitGroup
	monitorMutex sync.Mutex
}

// NewMonitor creates a new HTTP request monitor
func NewMonitor() *Monitor {
	return &Monitor{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				DisableKeepAlives: false,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
				},
			},
		},
		resultChan: make(chan shared.NetworkRequest, 100),
		stopChan:   make(chan struct{}),
	}
}

// Start begins monitoring targets
func (m *Monitor) Start(targets []shared.Target) {
	m.monitorMutex.Lock()
	defer m.monitorMutex.Unlock()

	// Stop any existing monitoring
	m.Stop()

	// Set new targets
	m.targets = targets
	m.stopChan = make(chan struct{})

	// Start monitoring each target
	for _, target := range targets {
		if !target.Enabled {
			continue
		}

		m.wg.Add(1)
		go m.monitorTarget(target)
	}
}

// Stop stops all monitoring
func (m *Monitor) Stop() {
	close(m.stopChan)
	m.wg.Wait()
}

// GetResultChan returns the channel where results are sent
func (m *Monitor) GetResultChan() <-chan shared.NetworkRequest {
	return m.resultChan
}

// monitorTarget monitors a single target periodically
func (m *Monitor) monitorTarget(target shared.Target) {
	defer m.wg.Done()

	ticker := time.NewTicker(time.Duration(target.Interval) * time.Second)
	defer ticker.Stop()

	// Initial request
	m.makeRequest(target)

	for {
		select {
		case <-ticker.C:
			m.makeRequest(target)
		case <-m.stopChan:
			return
		}
	}
}

// makeRequest performs an HTTP request and records metrics
func (m *Monitor) makeRequest(target shared.Target) {
	req, err := http.NewRequest(http.MethodGet, target.URL, nil)
	if err != nil {
		m.recordError(target, err, "request_creation")
		return
	}

	var result shared.NetworkRequest
	result.ID = uuid.New().String()
	result.URL = target.URL
	result.Method = http.MethodGet
	result.StartTime = time.Now()
	result.TargetName = target.Name

	var dnsStart, connectStart, tlsStart, requestStart, responseStart time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			dnsStart = time.Now()
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			if !dnsStart.IsZero() {
				result.DNSTime = time.Since(dnsStart).Milliseconds()
			}
			if info.Err != nil {
				m.recordError(target, info.Err, "dns")
			}
		},
		ConnectStart: func(network, addr string) {
			connectStart = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			if !connectStart.IsZero() {
				result.TCPTime = time.Since(connectStart).Milliseconds()
			}
			if err != nil {
				m.recordError(target, err, "connect")
			}
		},
		TLSHandshakeStart: func() {
			tlsStart = time.Now()
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			if !tlsStart.IsZero() {
				result.TLSTime = time.Since(tlsStart).Milliseconds()
			}
			if err != nil {
				m.recordError(target, err, "tls")
			}
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			if info.Err != nil {
				m.recordError(target, info.Err, "request_write")
			}
			requestStart = time.Now()
		},
		GotFirstResponseByte: func() {
			responseStart = time.Now()
			if !requestStart.IsZero() {
				result.RequestTime = time.Since(requestStart).Milliseconds()
			}
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))

	resp, err := m.client.Do(req)
	result.EndTime = time.Now()
	result.TotalTime = result.EndTime.Sub(result.StartTime).Milliseconds()

	if !responseStart.IsZero() {
		result.ResponseTime = time.Since(responseStart).Milliseconds()
	}

	if err != nil {
		result.Error = err.Error()
		switch e := err.(type) {
		case *net.OpError:
			result.ErrorType = "network"
			if e.Timeout() {
				result.ErrorType = "timeout"
			}
		case net.Error:
			if e.Timeout() {
				result.ErrorType = "timeout"
			} else {
				result.ErrorType = "network"
			}
		case *url.Error:
			result.ErrorType = "url"
			if e.Timeout() {
				result.ErrorType = "timeout"
			}
		default:
			result.ErrorType = "unknown"
		}
	} else {
		result.StatusCode = resp.StatusCode
		resp.Body.Close()
	}

	m.resultChan <- result
}

// recordError records an error during the request process
func (m *Monitor) recordError(target shared.Target, err error, errorType string) {
	result := shared.NetworkRequest{
		ID:         uuid.New().String(),
		URL:        target.URL,
		Method:     http.MethodGet,
		StartTime:  time.Now(),
		EndTime:    time.Now(),
		Error:      err.Error(),
		ErrorType:  errorType,
		TargetName: target.Name,
	}
	m.resultChan <- result
}