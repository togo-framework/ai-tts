// Package tts adds Text-to-Speech to togo. Drivers (elevenlabs, openai, …) register
// via init(); select one with TTS_DRIVER. Mirrors the ai plugin's driver pattern, so
// adding a provider is one self-contained driver. Mount the REST handler under /api/ai/tts.
package tts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/togo-framework/togo"
)

// Request is a synthesis request.
type Request struct {
	Text   string `json:"text"`
	Voice  string `json:"voice,omitempty"`  // provider voice id/name
	Model  string `json:"model,omitempty"`  // provider model
	Format string `json:"format,omitempty"` // mp3 (default) | wav
}

// Result is synthesized audio + its MIME type.
type Result struct {
	Audio       []byte
	ContentType string
}

// Provider is the TTS driver interface every driver implements.
type Provider interface {
	Synthesize(ctx context.Context, req Request) (Result, error)
}

// DriverFactory builds a Provider from the kernel/env.
type DriverFactory func(k *togo.Kernel) (Provider, error)

var (
	regMu   sync.RWMutex
	drivers = map[string]DriverFactory{}
)

// RegisterDriver registers a TTS driver (called from init()).
func RegisterDriver(name string, f DriverFactory) {
	regMu.Lock()
	drivers[name] = f
	regMu.Unlock()
}

// Drivers returns the registered driver names.
func Drivers() []string {
	regMu.RLock()
	defer regMu.RUnlock()
	out := make([]string, 0, len(drivers))
	for n := range drivers {
		out = append(out, n)
	}
	return out
}

func init() {
	RegisterDriver("elevenlabs", func(*togo.Kernel) (Provider, error) { return newElevenLabs() })
	RegisterDriver("openai", func(*togo.Kernel) (Provider, error) { return newOpenAI() })

	togo.RegisterProviderFunc("ai-tts", togo.PriorityService, func(k *togo.Kernel) error {
		name := os.Getenv("TTS_DRIVER")
		if name == "" {
			name = "openai"
		}
		regMu.RLock()
		f, ok := drivers[name]
		regMu.RUnlock()
		if !ok {
			return fmt.Errorf("ai-tts: unknown driver %q (set TTS_DRIVER; e.g. elevenlabs, openai)", name)
		}
		p, err := f(k)
		if err != nil {
			return err
		}
		k.Set("ai-tts", &Service{provider: p, driver: name})
		return nil
	})
}

// Service is the kernel-bound TTS service.
type Service struct {
	provider Provider
	driver   string
}

// Synthesize turns text into audio using the configured driver.
func (s *Service) Synthesize(ctx context.Context, req Request) (Result, error) {
	return s.provider.Synthesize(ctx, req)
}

// Driver returns the active driver name.
func (s *Service) Driver() string { return s.driver }

// FromKernel returns the TTS service bound to the kernel.
func FromKernel(k *togo.Kernel) (*Service, bool) {
	v, ok := k.Get("ai-tts")
	if !ok {
		return nil, false
	}
	s, ok := v.(*Service)
	return s, ok
}

// ── ElevenLabs driver ────────────────────────────────────────────────────────

type elevenLabs struct {
	key    string
	client *http.Client
}

func newElevenLabs() (Provider, error) {
	key := os.Getenv("ELEVENLABS_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("ai-tts/elevenlabs: ELEVENLABS_API_KEY not set")
	}
	return &elevenLabs{key: key, client: &http.Client{Timeout: 60 * time.Second}}, nil
}

func (e *elevenLabs) Synthesize(ctx context.Context, req Request) (Result, error) {
	voice := req.Voice
	if voice == "" {
		voice = "21m00Tcm4TlvDq8ikWAM" // "Rachel" default
	}
	model := req.Model
	if model == "" {
		model = "eleven_multilingual_v2"
	}
	body, _ := json.Marshal(map[string]any{"text": req.Text, "model_id": model})
	url := "https://api.elevenlabs.io/v1/text-to-speech/" + voice
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Result{}, err
	}
	r.Header.Set("xi-api-key", e.key)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "audio/mpeg")
	resp, err := e.client.Do(r)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()
	audio, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("ai-tts/elevenlabs: %s: %s", resp.Status, string(audio))
	}
	return Result{Audio: audio, ContentType: "audio/mpeg"}, nil
}

// ── OpenAI TTS driver ─────────────────────────────────────────────────────────

type openAI struct {
	key    string
	client *http.Client
}

func newOpenAI() (Provider, error) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("ai-tts/openai: OPENAI_API_KEY not set")
	}
	return &openAI{key: key, client: &http.Client{Timeout: 60 * time.Second}}, nil
}

func (o *openAI) Synthesize(ctx context.Context, req Request) (Result, error) {
	model := req.Model
	if model == "" {
		model = "tts-1"
	}
	voice := req.Voice
	if voice == "" {
		voice = "alloy"
	}
	format := req.Format
	if format == "" {
		format = "mp3"
	}
	body, _ := json.Marshal(map[string]any{
		"model": model, "input": req.Text, "voice": voice, "response_format": format,
	})
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/audio/speech", bytes.NewReader(body))
	if err != nil {
		return Result{}, err
	}
	r.Header.Set("Authorization", "Bearer "+o.key)
	r.Header.Set("Content-Type", "application/json")
	resp, err := o.client.Do(r)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()
	audio, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("ai-tts/openai: %s: %s", resp.Status, string(audio))
	}
	ct := "audio/mpeg"
	if format == "wav" {
		ct = "audio/wav"
	}
	return Result{Audio: audio, ContentType: ct}, nil
}
