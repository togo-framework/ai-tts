# ai-tts — documentation

  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />

## Install

```bash
togo install togo-framework/ai-tts
```

A capability plugin — it self-registers on boot; no driver selector needed.

## Configuration

Environment variables read by this plugin (extracted from the source):

| Env var | Notes |
|---|---|
| `ELEVENLABS_API_KEY` | _see provider docs_ |
| `G` | _see provider docs_ |
| `OPENAI_API_KEY` | _see provider docs_ |
| `TTS_DRIVER` | _see provider docs_ |

## Usage

```go
audio, err := tts.FromKernel(k).Synthesize(ctx, "Hello world", tts.Options{})
```

## Links

- Marketplace: https://to-go.dev/marketplace
- Source: https://github.com/togo-framework/ai-tts
- README: ../README.md
