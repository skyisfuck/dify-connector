package api

import (
	"context"
	"errors"
	"github.com/gin-contrib/graceful"
	"github.com/skyisfuck/dify-connector/internal/api/controller"
	"github.com/skyisfuck/dify-connector/internal/api/middleware"
	"github.com/skyisfuck/dify-connector/internal/misc"
	"sync"
)

func StartAPI(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	r, err := graceful.Default(
		graceful.WithAddr(":" + misc.GetEnv("API_PORT", "6789")),
	)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	middleware.RegisterMiddlewares(r)
	middleware.RegisterHealthCheck(r, ctx)

	registerCompatibleAPIs(r)

	if err := r.RunWithContext(ctx); err != nil && !errors.Is(err, context.Canceled) {
		panic(err)
	}
}

// registerCompatibleAPIs register OpenAI compatible APIs
func registerCompatibleAPIs(r *graceful.Graceful) {
	g := r.Group("/compatible/v1")

	g.GET("/models", controller.Model)
	g.POST("/chat/completions", controller.ChatCompletion)
}
