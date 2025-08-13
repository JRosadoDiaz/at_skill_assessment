package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/JRosadoDiaz/AT_Skill_Assessment/api"
	"github.com/JRosadoDiaz/AT_Skill_Assessment/internal/pinger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var pingManager *pinger.PingManager

func StartServer(port string, pm *pinger.PingManager) {
	pingManager = pm

	router := chi.NewRouter() // Returns a new Mux object that implements the Router interface
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./web/templates/index.html")
		if err != nil {
			api.InternalErrorHandler(w)
			log.Printf("Error loading template: %v", err)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			api.InternalErrorHandler(w)
			log.Printf("Error rendering template: %v", err)
			return
		}
	})

	fmt.Printf("Server is now listening at http://localhost:%v \n", port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Println("failed to listen to server", err)
		log.Printf("Failed to listen to server: %v", err)
	}
}

// // Mainpage
// func serveHome() {
// 	if r.URL.Path != "/" {
// 		http.NotFound(w, r)
// 		return
// 	}

// }

// func handler(r *chi.Mux, pinger *pinger.PingManager) { // Mux -> Multiplexer
// 	r.Route("/metrics", func(router chi.Router) {
// 		router.Get("/", pinger.GetMetrics)
// 	})
// }

func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
