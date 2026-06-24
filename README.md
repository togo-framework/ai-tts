<!-- togo-header -->
<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/ai-tts</h1>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/ai-tts"><img src="https://pkg.go.dev/badge/github.com/togo-framework/ai-tts.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/ai-tts
```

<!-- /togo-header -->

<p align="center"><img src="https://to-go.dev/togo-mark.svg" alt="togo" height="64"></p>
<h1 align="center">ai-tts</h1>
<p align="center">Text-to-Speech for <a href="https://to-go.dev">togo</a> — multi-provider, one interface.</p>

---

`ai-tts` adds **text-to-speech** to a togo app. It mirrors the togo driver pattern: a
single `Synthesize` interface with swappable provider drivers, selected by `TTS_DRIVER`.

## Install

```bash
togo install togo-framework/ai-tts
```

## Drivers

| Driver | Env | Notes |
|---|---|---|
| `openai` (default) | `OPENAI_API_KEY` | OpenAI TTS (`tts-1`, voices: alloy/echo/fable/onyx/nova/shimmer) |
| `elevenlabs` | `ELEVENLABS_API_KEY` | ElevenLabs (multilingual v2; `voice` = voice id) |

```bash
TTS_DRIVER=elevenlabs
ELEVENLABS_API_KEY=...
```

Add another provider by registering a driver in an `init()` — see `tts.RegisterDriver`.

## Use (Go)

```go
svc, _ := tts.FromKernel(k)
res, err := svc.Synthesize(ctx, tts.Request{Text: "Hello from togo", Voice: "nova"})
// res.Audio (bytes), res.ContentType ("audio/mpeg")
```

## Use (REST)

Mount the handler under `/api/ai/tts`:

```go
mux.Handle("/api/ai/tts/", http.StripPrefix("/api/ai/tts", tts.Handler(k)))
```

```bash
curl -X POST http://localhost:8080/api/ai/tts/ \
  -H 'content-type: application/json' \
  -d '{"text":"Hello from togo","voice":"nova"}' --output hello.mp3
```

Pairs with [`ai-stt`](https://github.com/togo-framework/ai-stt) (speech-to-text) and the
[`ai`](https://github.com/togo-framework/ai) plugin. MIT.

<!-- togo-sponsors -->
---

<div align="center">
  <h3>Premium sponsors</h3>
  <p>
    <a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp;
    <a href="https://one-studio.co"><strong>One Studio</strong></a>
  </p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
