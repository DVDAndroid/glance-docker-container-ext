package main

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/gorilla/schema"
	"log/slog"
	"net/http"
	"os"
)

type params struct {
	WidgetTitle   string `schema:"title,default:Docker Containers"`
	Group         string `schema:"group"`
	AllContainers bool   `schema:"all,default:true"`
	Order         string `schema:"order,default:name"`
	SameTab       bool   `schema:"same-tab,default:false"`
	IgnoreStatus  bool   `schema:"ignore-status,default:false"`
}

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		slog.Error("error creating docker client", "err", err)
		panic(err)
	}

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var p params
		err := decoder.Decode(&p, r.URL.Query())
		if err != nil {
			slog.Error("error decoding params", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Widget-Title", p.WidgetTitle)
		w.Header().Set("Widget-Content-Type", "html")
		w.Header().Set("Content-Type", "text/html")

		containers, err := LoadContainers(dockerClient, p)
		if err != nil {
			slog.Error("cannot connect to docker engine", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`<div class="widget-error-header">
                <div class="color-negative size-h3">ERROR</div>
                <div class="widget-error-icon"></div>
            </div>
            <p class="break-all">Cannot connect to Docker Engine</p>`))
			return
		}

		if len(containers) == 0 {
			w.Write([]byte(`<p>No containers found</p>`))
			return
		}

		BuildHtml(w, containers)
	})

	slog.Info("starting webserver", "host", host, "port", port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil)
	if err != nil {
		slog.Error("error starting webserver", "err", err)
		return
	}
}
