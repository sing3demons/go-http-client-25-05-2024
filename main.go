package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tlog "github.com/opentracing/opentracing-go/log"
	"github.com/sing3demons/go-http-client/tracing"
)

func main() {
	tracer, closer := tracing.Init("service-1")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	http.HandleFunc("/joke", func(w http.ResponseWriter, r *http.Request) {

		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("fetch-joke", opentracing.ChildOf(spanCtx))
		defer span.Finish()

		// Span 2
		span2 := tracer.StartSpan("step1", opentracing.ChildOf(span.Context()))
		span2.SetTag("step", "step 1")
		defer span2.Finish()

		// Span 3
		ctx := opentracing.ContextWithSpan(context.Background(), span2)

		type Joke struct {
			Type      string `json:"type"`
			SetUp     string `json:"setup"`
			Punchline string `json:"punchline"`
		}

		response, err := callService[Joke](ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(response)
	})
	log.Fatal(http.ListenAndServe(":8000", nil))

}

func callService[T any](ctx context.Context) (T, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "callService")
	defer span.Finish()

	url := "https://official-joke-api.appspot.com/random_joke"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	defer resp.Body.Close()
	var result T
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	json.Unmarshal(body, &result)

	span.LogFields(
		tlog.String("event", "call_target"),
		tlog.String("value", string(body)),
	)

	return result, nil
}
