package main

import (
	"lets-jaegar/lib"
	"log"
	"net/http"

	"github.com/opentracing/opentracing-go"
)

func main() {
	tracer, closer := lib.InitializeJaeger("service-two")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		span := lib.StartSpanFromRequest(tracer, r)
		defer span.Finish()

		w.Write([]byte("service-two"))
	})
	log.Printf("Listening on localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
