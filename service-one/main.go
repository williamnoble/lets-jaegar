package main

import (
	"context"
	"fmt"
	"lets-jaegar/lib"
	"log"
	"net/http"
	"os"
	"github.com/opentracing/opentracing-go"
)

func main() {
	/*
	1. Initialize Jaeger and set as the global tracer. Open tracing can use several tracers.
	2. When calling the service-one "Ping" Handler start the Span.
	3. Create a context.Context to hold the Span.
	4. Call the Ping Function and start Span from the given context, injecting this Span to the new request.
	 */
	tracer, closer := lib.InitializeJaeger("service-one")
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	outboundHostPort, ok := os.LookupEnv("OUTBOUND_HOST_PORT")
	if !ok {
		outboundHostPort = "localhost:8082"
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		span := lib.StartSpanFromRequest(tracer, r)
		defer span.Finish()

		ctx := opentracing.ContextWithSpan(context.Background(), span)
		response, err := lib.Ping(ctx, outboundHostPort)
		if err != nil {
			log.Fatalf("Error occurred: %s", err)
		}
		w.Write([]byte(fmt.Sprintf("%s -> %s", "service-one", response)))
	})
	log.Printf("Listening on localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
