package main

import (
	"os"

	"github.com/go-kit/kit/log"
)

func main() {
	w := log.NewSyncWriter(os.Stdout)
	l := log.NewLogfmtLogger(w)
	l = log.With(l, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	l.Log("msg", "Hello World!")
}
