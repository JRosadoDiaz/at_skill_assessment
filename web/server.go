package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/JRosadoDiaz/AT_Skill_Assessment/api"
	"github.com/JRosadoDiaz/AT_Skill_Assessment/internal/pinger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
)

type Server struct {
	pinger *pinger.PingManager
	router *chi.Mux
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://"+r.Host
	},
}

func NewServer(pm *pinger.PingManager) *Server {
	s := &Server{
		pinger: pm,
		router: chi.NewRouter(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.Use(middleware.Logger)

	fileServer := http.FileServer(http.Dir("web/static"))
	s.router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	templateServer := http.FileServer(http.Dir("web/templates"))
	s.router.Handle("/templates/*", http.StripPrefix("/templates/", templateServer))

	s.router.Get("/", s.handleHome)
	s.router.Get("/ws", s.socketHandler)
}

func (s *Server) Start(port string) {
	fmt.Printf("Server is now listening at http://localhost:%v \n", port)
	log.Fatal(http.ListenAndServe(":"+port, s.router))
}

func (s *Server) socketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Websocket connection recieved")
	conn, err := upgrader.Upgrade(w, r, nil) // Upgrades http request to a websocket for back and forth communication
	if err != nil {
		log.Printf("Websocket upgrade failed: %v", err)
	}

	go func(conn *websocket.Conn) {
		log.Printf("Starting data push with interval: %v", s.pinger.Interval)
		ticker := time.NewTicker(s.pinger.Interval)
		defer func() {
			ticker.Stop()
			conn.Close()
		}()

		for range ticker.C {
			metrics := s.pinger.GetMetrics()

			if len(metrics) == 0 {
				log.Println("Metrics map is empty...")
				continue
			}

			data, err := json.Marshal(metrics)
			if err != nil {
				log.Println("JSON encoding error:", err)
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("Websocket write error, closing connection:", err)
				return
			}
		}
	}(conn)
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./web/static/index.html")
	if err != nil {
		api.InternalErrorHandler(w)
		log.Printf("Error loading template: %v", err)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		api.InternalErrorHandler(w)
		log.Printf("Rendering the template has failed: %v", err)
	}
}
