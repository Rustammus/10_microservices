package rest

import (
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"log/slog"
	"microservices/task_6/pkg/logger"
	"net/http"
	"runtime/debug"
	"time"
)

func (s *API) BaseMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqId, err := uuid.NewV7()
		if err != nil {
			s.l.Error("error generate uuidV7", slog.String("error", err.Error()))
		} else {
			r = r.WithContext(logger.AppendCtx(r.Context(), slog.String("req_id", reqId.String())))
		}

		defer func() {
			if err := recover(); err != nil {
				s.l.ErrorContext(r.Context(), "recover from panic",
					slog.Any("error", err), slog.String("stack", string(debug.Stack())))
			}
		}()
		next(w, r)

		// Increment requests counter
		s.m.RequestsIncrement()
	}
}

func (s *API) LoggerMiddleware(next *httprouter.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sTime := time.Now()
		next.ServeHTTP(w, r)
		reqTime := time.Since(sTime)
		s.l.DebugContext(r.Context(), "new request",
			slog.String("method", r.Method),
			slog.String("status", w.Header().Get("status")),
			slog.String("duration", reqTime.String()),
			slog.String("address", r.RemoteAddr),
			slog.String("url", r.URL.String()),
		)
		s.l.DebugContext(r.Context(), "logger middleware trace", slog.String("trace", string(debug.Stack())))
	}
}
