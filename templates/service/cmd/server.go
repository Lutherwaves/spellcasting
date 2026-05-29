package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tink3rlabs/magic/errors"
	"github.com/tink3rlabs/magic/health"
	"github.com/tink3rlabs/magic/middlewares"
	"github.com/tink3rlabs/magic/observability"
	"github.com/tink3rlabs/magic/storage"
)

func requestSizeMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			maxBytes := viper.GetInt64("request.max_body_size")
			if maxBytes <= 0 {
				maxBytes = 10 << 20
			}
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}

func initRoutes(storageAdapter storage.StorageAdapter, obs *observability.Observer) chi.Router {
	router := chi.NewRouter()

	router.Use(middlewares.ObservabilityWithOptions(obs, middlewares.ObservabilityOptions{
		SkipPaths:        []string{"/metrics"},
		SkipPathPrefixes: []string{"/health/"},
	}))
	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
		requestSizeMiddleware(),
		cors.Handler(cors.Options{
			AllowedOrigins:   viper.GetStringSlice("cors.allowed_origins"),
			AllowedMethods:   viper.GetStringSlice("cors.allowed_methods"),
			AllowedHeaders:   viper.GetStringSlice("cors.allowed_headers"),
			ExposedHeaders:   viper.GetStringSlice("cors.exposed_headers"),
			AllowCredentials: viper.GetBool("cors.allow_credentials"),
			MaxAge:           viper.GetInt("cors.max_age"),
		}),
	)

	router.Handle("/metrics", obs.MetricsHandler())

	h := middlewares.ErrorHandler{}

	authConfig := middlewares.EnsureValidTokenConfig{
		Enabled:          viper.GetBool("auth.enabled"),
		IssuerURL:        viper.GetString("auth.config.oauth_provider_url"),
		Audience:         viper.GetStringSlice("auth.config.oauth_audience"),
		AllowedClockSkew: viper.GetDuration("auth.config.allowed_clock_skew"),
	}
	authMiddleware := middlewares.EnsureValidToken(authConfig)

	middlewares.SetDefaultClaimsConfig(middlewares.ClaimsConfig{
		TenantIdKey: viper.GetString("auth.config.claims.tenant_id"),
		EmailKey:    viper.GetString("auth.config.claims.email"),
		RolesKey:    viper.GetString("auth.config.claims.roles"),
		GroupsKey:   viper.GetString("auth.config.claims.groups"),
	})

	healthChecker := health.NewHealthChecker(storageAdapter)
	router.Route("/health", func(r chi.Router) {
		r.Get("/liveness", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, map[string]string{"status": "ok"})
		})
		r.Get("/readiness", h.Wrap(func(w http.ResponseWriter, r *http.Request) error {
			if err := healthChecker.Check(viper.GetBool("health.storage"), viper.GetStringSlice("health.dependencies")); err != nil {
				slog.Error("readiness failed", slog.Any("error", err.Error()))
				return &errors.ServiceUnavailable{Message: err.Error()}
			}
			render.JSON(w, r, map[string]string{"status": "ok"})
			return nil
		}))
	})

	router.Route("/", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Use(middlewares.TenantRequestContext)
		r.Use(middlewares.UserRequestContext)
		// The casting-a-route skill appends mounts to this block.
	})

	return router
}

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Start the SERVICENAME server",
	RunE:  runServer,
}

func runServer(cmd *cobra.Command, args []string) error {
	storageAdapter, err := storage.StorageAdapterFactory{}.GetInstance(
		storage.StorageAdapterType(viper.GetString("storage.type")),
		viper.GetStringMapString("storage.config"),
	)
	if err != nil {
		return fmt.Errorf("storage adapter: %w", err)
	}

	storage.NewDatabaseMigration(storageAdapter).Migrate()

	ctx := context.Background()
	obs, err := setupObservability(ctx, "SERVICENAME")
	if err != nil {
		return fmt.Errorf("observability init: %w", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = obs.Shutdown(shutdownCtx)
	}()

	router := initRoutes(storageAdapter, obs)

	port := viper.GetString("service.port")
	if port == "" {
		port = "SERVICEPORT"
	}
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  viper.GetDuration("server.read_timeout"),
		WriteTimeout: viper.GetDuration("server.write_timeout"),
		IdleTimeout:  viper.GetDuration("server.idle_timeout"),
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		slog.Info("shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("forced shutdown", slog.Any("error", err))
		}
	}()

	slog.Info("starting HTTP server", slog.String("address", server.Addr))
	return server.ListenAndServe()
}
