package llms

import "context"

type AiModelsSvc interface {
	Embedder(ctx context.Context) error
	Transcriber(ctx context.Context) error
	Synthesizer(ctx context.Context) error
	GenerateSyntheticData(ctx context.Context) error
	ChatCompletor(ctx context.Context) error
}
