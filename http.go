package tts

import (
	"net/http"

	"github.com/togo-framework/togo"
)

// Handler exposes the TTS service over REST. Mount under /api/ai/tts in your app:
//
//	mux.Handle("/api/ai/tts/", http.StripPrefix("/api/ai/tts", tts.Handler(k)))
//
// POST /  with a JSON Request -> raw audio (audio/mpeg or audio/wav).
func Handler(k *togo.Kernel) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		svc, ok := FromKernel(k)
		if !ok {
			http.Error(w, "ai-tts not configured", http.StatusInternalServerError)
			return
		}
		var req Request
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res, err := svc.Synthesize(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", res.ContentType)
		_, _ = w.Write(res.Audio)
	})
	return mux
}
