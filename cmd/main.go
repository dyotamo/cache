package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/{key}", func(w http.ResponseWriter, r *http.Request) {
		body := r.Body
		defer body.Close()

		value, err := io.ReadAll(body)
		if err != nil || strings.EqualFold(string(value), "") {
			w.WriteHeader(400)
			w.Write([]byte("no value provided"))
			return
		}

		key := chi.URLParam(r, "key")
		err = rdb.Set(ctx, key, string(value), 0).Err()
		if err != nil {
			panic(err)
		}

		w.Write(value)
	})

	r.Get("/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		value, err := rdb.Get(ctx, key).Result()
		if err == redis.Nil {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf("%s is not set", key)))
			return
		}

		w.Write([]byte(value))
	})

	slog.Info("Server started!")
	http.ListenAndServe(":8080", r)
}
