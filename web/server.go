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

var pingManager *pinger.PingManager

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartServer(port string, pm *pinger.PingManager) {
	pingManager = pm

	router := chi.NewRouter() // Returns a new Mux object that implements the Router interface
	router.Use(middleware.Logger)

	router.Get("/", handleHome)

	// router.Get("/ws", socketHandler)

	fmt.Printf("Server is now listening at http://localhost:%v \n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Websocket connection recieved")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		api.InternalErrorHandler(w)
		log.Printf("Websocket upgrade failed:", err)
	}

	go writeMetrics(conn)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
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
}

func writeMetrics(conn *websocket.Conn) {
	ticker := time.NewTicker(pingManager.Interval)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for i:=0; c range ticker{
		case <-ticker.C:
			metrics := pingManager.GetMetrics()

			data, err := json.Marshal(metrics)
			if err != nil {
				log.Println("JSON encoding error:", err)
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("WebSocket write error, closing connection:", err)
				return
			}
		}
	}
}
