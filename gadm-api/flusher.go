package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Flusher struct {
	w       http.ResponseWriter
	flusher http.Flusher
	ctx     context.Context
}

func NewFlusher(w http.ResponseWriter, ctx context.Context) (*Flusher, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("ResponseWriter doesn't support flushing")
	}

	return &Flusher{
		w:       w,
		flusher: flusher,
		ctx:     ctx,
	}, nil
}

func (f *Flusher) Flush(data json.RawMessage) error {
	if _, err := f.w.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}
	if _, err := f.w.Write([]byte("\n")); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	f.flusher.Flush()

	select {
	case <-f.ctx.Done():
		return f.ctx.Err()
	default:
		return nil
	}
}
