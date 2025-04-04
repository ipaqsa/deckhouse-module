package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/x/module/internal/config"
	"github.com/x/module/internal/service/keyvalue"
	"github.com/x/module/internal/service/roller"
	"github.com/x/module/internal/service/sayer"
)

type Server struct {
	kvService     *keyvalue.Service
	rollerService *roller.Service
	sayerService  *sayer.Service
}

func New(conf config.Config) *Server {
	server := new(Server)

	if conf.KeyValue.Enabled {
		server.kvService = keyvalue.New()
	}

	if conf.Sayer.Enabled {
		server.sayerService = sayer.New(conf.Sayer)
	}

	if conf.Roller.Enabled {
		server.rollerService = roller.New(conf.Roller)
	}

	return server
}

func (s *Server) Serve(_ context.Context, address string) error {
	var routes []string

	if s.kvService != nil {
		log.Println("register kv handlers")
		http.HandleFunc("GET /kv/{key}", s.handleGetKey)
		http.HandleFunc("PUT /kv/{key}", s.handlePutKey)
		http.HandleFunc("DELETE /kv/{key}", s.handleDeleteKey)

		routes = append(routes, "GET /kv/{key}")
		routes = append(routes, "PUT /kv/{key}")
		routes = append(routes, "DELETE /kv/{key}")
	}

	if s.sayerService != nil {
		log.Println("register sayer handlers")
		http.HandleFunc("GET /say", s.handleSay)

		routes = append(routes, "GET /say")
	}

	if s.rollerService != nil {
		log.Println("register roller handlers")
		http.HandleFunc("GET /roll", s.handleRoll)

		routes = append(routes, "GET /roll")
	}

	http.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
		sort.Strings(routes)

		w.Header().Set("Content-Type", "text/plain")

		fmt.Fprintln(w, "Available endpoints:")
		for _, route := range routes {
			fmt.Fprintf(w, "  %s\n", route)
		}

		w.WriteHeader(http.StatusOK)
	})

	// for probes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("server listening on %s", address)
	return http.ListenAndServe(address, nil)
}

func (s *Server) handleGetKey(w http.ResponseWriter, r *http.Request) {
	log.Println("handle get key request")

	key := r.PathValue("key")
	if key == "" {
		log.Println("no key")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value := s.kvService.Get(key)
	if value == "" {
		log.Println("key not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(value)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handlePutKey(w http.ResponseWriter, r *http.Request) {
	log.Println("handle put key request")

	key := r.PathValue("key")
	if key == "" {
		log.Println("no key")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value := r.URL.Query().Get("value")
	if value == "" {
		log.Println("no value")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.kvService.Set(key, value)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleDeleteKey(w http.ResponseWriter, r *http.Request) {
	log.Println("handle delete key request")

	key := r.PathValue("key")
	if key == "" {
		log.Println("no key")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.kvService.Delete(key)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleSay(w http.ResponseWriter, _ *http.Request) {
	log.Println("handle say request")

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(s.sayerService.Say())); err != nil {
		log.Printf("failed to write say text: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleRoll(w http.ResponseWriter, _ *http.Request) {
	log.Println("handle roll request")

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(s.rollerService.RollDice())); err != nil {
		log.Printf("failed to write roll result: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
