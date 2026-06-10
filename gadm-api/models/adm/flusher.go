package adm

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type flusher struct {
	w       io.Writer
	flusher http.Flusher
	ctx     context.Context
}

func newFlusher(ctx context.Context, w io.Writer) (*flusher, error) {
	f, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("response_writer_does_not_support_flushing")
	}

	return &flusher{
		w:       w,
		flusher: f,
		ctx:     ctx,
	}, nil
}

func (f *flusher) flush(data []byte) error {
	dataWithNewline := append(data, []byte("\n")...)
	if _, err := f.w.Write(dataWithNewline); err != nil {
		return fmt.Errorf("failed_to_write_data: %w", err)
	}

	f.flusher.Flush()

	select {
	case <-f.ctx.Done():
		return f.ctx.Err()
	default:
		return nil
	}
}
