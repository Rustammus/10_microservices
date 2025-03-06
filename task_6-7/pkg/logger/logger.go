package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

func NewLogger(lvl slog.Level, output string) *slog.Logger {

	var out io.Writer = os.Stdout

	if output != "" {
		var err error
		out, err = os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	}

	l := slog.New(&ContextHandler{slog.NewJSONHandler(out, &slog.HandlerOptions{
		AddSource:   false,
		Level:       lvl,
		ReplaceAttr: nil,
	})})
	slog.SetDefault(l)

	return l
}

type ctxKey string

const slogFields ctxKey = "slog_fields"

type ContextHandler struct {
	slog.Handler
}

// Handle adds contextual attributes to the Record before calling the underlying
// handler
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	return h.Handler.Handle(ctx, r)
}

// AppendCtx adds an slog attribute to the provided context so that it will be
// included in any Record created with such context
func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	v, _ := parent.Value(slogFields).([]slog.Attr)
	v = append(v, attr)
	return context.WithValue(parent, slogFields, v)

	//if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
	//	v = append(v, attr)
	//	return context.WithValue(parent, slogFields, v)
	//}

	//v := []slog.Attr{}
	//v = append(v, attr)
	//return context.WithValue(parent, slogFields, v)
}
