package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gproc/pkg/types"
)

type AggregationManager struct {
	config    *types.LogAggregationConfig
	providers map[string]LogProvider
}

type LogProvider interface {
	SendLogs(ctx context.Context, logs []LogEntry) error
	Query(ctx context.Context, query LogQuery) ([]LogEntry, error)
}

type LogEntry struct {
	Timestamp   time.Time         `json:"@timestamp"`
	Level       string            `json:"level"`
	Message     string            `json:"message"`
	ProcessID   string            `json:"process_id"`
	ProcessName string            `json:"process_name"`
	NodeID      string            `json:"node_id"`
}

type LogQuery struct {
	ProcessName string    `json:"process_name,omitempty"`
	Level       string    `json:"level,omitempty"`
	StartTime   time.Time `json:"start_time,omitempty"`
	EndTime     time.Time `json:"end_time,omitempty"`
	Limit       int       `json:"limit,omitempty"`
}

type ElasticSearchProvider struct {
	endpoint string
}

func NewAggregationManager(config *types.LogAggregationConfig) *AggregationManager {
	return &AggregationManager{
		config:    config,
		providers: make(map[string]LogProvider),
	}
}

func (a *AggregationManager) SendLog(ctx context.Context, entry LogEntry) error {
	for _, provider := range a.providers {
		provider.SendLogs(ctx, []LogEntry{entry})
	}
	return nil
}

func (e *ElasticSearchProvider) SendLogs(ctx context.Context, logs []LogEntry) error {
	fmt.Printf("Sending %d logs to ElasticSearch\n", len(logs))
	return nil
}

func (e *ElasticSearchProvider) Query(ctx context.Context, query LogQuery) ([]LogEntry, error) {
	return []LogEntry{}, nil
}