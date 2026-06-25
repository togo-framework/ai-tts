# ai-tts — documentation

Text-to-Speech for togo — ElevenLabs + OpenAI TTS drivers (TTS_DRIVER)

## Install

```bash
togo install togo-framework/ai-tts
```

A capability plugin — it self-registers on boot; no driver selector needed.

## Configuration

Environment variables read by this plugin (extracted from the source — see the gateway/provider docs for each value):

| Env var |
|---|
| `ELEVENLABS_API_KEY` |
| `OPENAI_API_KEY` |
| `TTS_DRIVER` |

## Usage

```go
audio, err := tts.FromKernel(k).Synthesize(ctx, "Hello world", tts.Options{})
```

## Links

- Marketplace: https://to-go.dev/marketplace
- Source: https://github.com/togo-framework/ai-tts
- Full README: ../README.md
