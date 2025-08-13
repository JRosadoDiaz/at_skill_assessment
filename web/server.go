package web

import (
	"fmt"
	"net/http"

	"github.com/JRosadoDiaz/AT_Skill_Assessment/internal/pinger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var pingManager *pinger.PingManager

func StartServer(port string, pm *pinger.PingManager) {
	pingManager = pm

	router := chi.NewRouter() // Returns a new Mux object that implements the Router interface
	router.Use(middleware.Logger)

	router.Get("/hello", basicHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	fmt.Printf("Server is now listening at http://localhost:%v \n", port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("failed to listen to server", err)
	}
}

// // Mainpage
// func serveHome() {
// 	if r.URL.Path != "/" {
// 		http.NotFound(w, r)
// 		return
// 	}

// }

// func Handler(r *chi.Mux, pinger *pinger.PingManager) {
// 	r.Route("/metrics", func(router chi.Router) {
// 		router.Get("/", pinger.GetMetrics)
// 	})
// }

func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
