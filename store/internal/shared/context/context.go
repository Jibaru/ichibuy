package context

import "context"

type ContextKey string

const (
	APITokenKey ContextKey = "ichibuy-api-token"
)

func AddToken(ctx context.Context, key any) context.Context {
	if token := ctx.Value(APITokenKey); token != nil {
		return context.WithValue(ctx, key, token)
	}
	return ctx
}
