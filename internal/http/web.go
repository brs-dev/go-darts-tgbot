package local_server

import (
	"fmt"
	cfg "go-darts-tgbot/internal/config"
	"log/slog"
	"net/http"
)

func HttpLocal() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("http server responsed")
		fmt.Fprint(w, "OK")
	})

	go func() {
		http.ListenAndServe("0.0.0.0:"+cfg.GlobalConfig.Port, nil)
	}()
}
