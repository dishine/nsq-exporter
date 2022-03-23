package collector

import (
	"encoding/json"
	"net/http"
)

type statsResponse struct {
	StatusCode int    `json:"status_code"`
	StatusText string `json:"status_text"`
	Data       stats  `json:"data"`
}

type stats struct {
	Version   string   `json:"version"`
	Health    string   `json:"health"`
	StartTime int64    `json:"start_time"`
	Topics    []*topic `json:"topics"`
}

// see https://github.com/nsqio/nsq/blob/master/nsqd/stats.go
type topic struct {
	Name         string     `json:"topic_name"`
	Paused       bool       `json:"paused"`
	Depth        int64      `json:"depth"`
	BackendDepth int64      `json:"backend_depth"`
	MessageCount uint64     `json:"message_count"`
	E2eLatency   e2elatency `json:"e2e_processing_latency"`
	Channels     []*channel `json:"channels"`
}

type channel struct {
	Name          string     `json:"channel_name"`
	Paused        bool       `json:"paused"`
	Depth         int64      `json:"depth"`
	BackendDepth  int64      `json:"backend_depth"`
	MessageCount  uint64     `json:"message_count"`
	InFlightCount int        `json:"in_flight_count"`
	DeferredCount int        `json:"deferred_count"`
	RequeueCount  uint64     `json:"requeue_count"`
	TimeoutCount  uint64     `json:"timeout_count"`
	E2eLatency    e2elatency `json:"e2e_processing_latency"`
	Clients       []*client  `json:"clients"`
}

type client struct {
	ID            string `json:"client_id"`
	Hostname      string `json:"hostname"`
	Version       string `json:"version"`
	RemoteAddress string `json:"remote_address"`
	State         int32  `json:"state"`
	FinishCount   uint64 `json:"finish_count"`
	MessageCount  uint64 `json:"message_count"`
	ReadyCount    int64  `json:"ready_count"`
	InFlightCount int64  `json:"in_flight_count"`
	RequeueCount  uint64 `json:"requeue_count"`
	ConnectTime   int64  `json:"connect_ts"`
	SampleRate    int32  `json:"sample_rate"`
	Deflate       bool   `json:"deflate"`
	Snappy        bool   `json:"snappy"`
	TLS           bool   `json:"tls"`
}
type e2elatency struct {
	Count       int                  `json:"count"`
	Percentiles []map[string]float64 `json:"percentiles"`
}

func (e *e2elatency) percentileValue(idx int) float64 {
	if idx >= len(e.Percentiles) {
		return 0
	}
	return e.Percentiles[idx]["value"]
}

func getPercentile(t *Topics, percentile int) float64 {
	if len(t.E2EProcessingLatency.Percentiles) > 0 {
		if percentile == 99 {
			return t.E2EProcessingLatency.Percentiles[0]["value"]
		} else if percentile == 95 {
			return t.E2EProcessingLatency.Percentiles[1]["value"]
		}
	}
	return 0
}

func getNsqdStats(client *http.Client, nsqdURL string) (*StatsResponse, error) {
	resp, err := client.Get(nsqdURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sr StatsResponse
	if err = json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, err
	}
	return &sr, nil
}

type StatsResponse struct {
	Version   string    `json:"version"`
	Health    string    `json:"health"`
	StartTime int       `json:"start_time"`
	Topics    []*Topics `json:"topics"`
	Memory    struct {
		HeapObjects       int `json:"heap_objects"`
		HeapIdleBytes     int `json:"heap_idle_bytes"`
		HeapInUseBytes    int `json:"heap_in_use_bytes"`
		HeapReleasedBytes int `json:"heap_released_bytes"`
		GcPauseUsec100    int `json:"gc_pause_usec_100"`
		GcPauseUsec99     int `json:"gc_pause_usec_99"`
		GcPauseUsec95     int `json:"gc_pause_usec_95"`
		NextGcBytes       int `json:"next_gc_bytes"`
		GcTotalRuns       int `json:"gc_total_runs"`
	} `json:"memory"`
	Producers interface{} `json:"producers"`
}
type Topics struct {
	TopicName            string      `json:"topic_name"`
	Channels             []*Channels `json:"channels"`
	Depth                int         `json:"depth"`
	BackendDepth         int         `json:"backend_depth"`
	MessageCount         int         `json:"message_count"`
	MessageBytes         int         `json:"message_bytes"`
	Paused               bool        `json:"paused"`
	E2EProcessingLatency E2elatency  `json:"e2e_processing_latency"`
}
type Channels struct {
	ChannelName          string     `json:"channel_name"`
	Depth                int        `json:"depth"`
	BackendDepth         int        `json:"backend_depth"`
	InFlightCount        int        `json:"in_flight_count"`
	DeferredCount        int        `json:"deferred_count"`
	MessageCount         int        `json:"message_count"`
	RequeueCount         int        `json:"requeue_count"`
	TimeoutCount         int        `json:"timeout_count"`
	ClientCount          int        `json:"client_count"`
	Clients              []*Clients `json:"clients"`
	Paused               bool       `json:"paused"`
	E2EProcessingLatency E2elatency `json:"e2e_processing_latency"`
}
type Clients struct {
	ClientId                      string `json:"client_id"`
	Hostname                      string `json:"hostname"`
	Version                       string `json:"version"`
	RemoteAddress                 string `json:"remote_address"`
	State                         int    `json:"state"`
	ReadyCount                    int    `json:"ready_count"`
	InFlightCount                 int    `json:"in_flight_count"`
	MessageCount                  int    `json:"message_count"`
	FinishCount                   int    `json:"finish_count"`
	RequeueCount                  int    `json:"requeue_count"`
	ConnectTs                     int    `json:"connect_ts"`
	SampleRate                    int    `json:"sample_rate"`
	Deflate                       bool   `json:"deflate"`
	Snappy                        bool   `json:"snappy"`
	UserAgent                     string `json:"user_agent"`
	Tls                           bool   `json:"tls"`
	TlsCipherSuite                string `json:"tls_cipher_suite"`
	TlsVersion                    string `json:"tls_version"`
	TlsNegotiatedProtocol         string `json:"tls_negotiated_protocol"`
	TlsNegotiatedProtocolIsMutual bool   `json:"tls_negotiated_protocol_is_mutual"`
}
type E2elatency struct {
	Count       int                  `json:"count"`
	Percentiles []map[string]float64 `json:"percentiles"`
}

func (e *E2elatency) percentileValue(idx int) float64 {
	if idx >= len(e.Percentiles) {
		return 0
	}
	return e.Percentiles[idx]["value"]
}
