package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AXBoBka/wb_l0/internal/cache"
)

type serverHTTP struct {
	cache *cache.Cache
}

func New(cache *cache.Cache) *serverHTTP {
	return &serverHTTP{cache}
}

func (s *serverHTTP) Start() {
	http.HandleFunc("/", s.Serve)
	log.Println("Сервер запущен на localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func (s *serverHTTP) Serve(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		path := "./internal/static/index.html"

		http.ServeFile(w, r, path)
	case "POST":
		r.ParseMultipartForm(0)
		id := r.FormValue("message")

		log.Printf("Получение заказа с id: %s", id)

		order, err := s.cache.GetOrder(id)

		if err != nil {
			log.Printf("Ошибка получения заказа: %s", err)
			fmt.Fprintf(w, "Нет заказа с данным ID")
			return
		}

		fmt.Fprintf(w, "Информация о заказе: %s", order)
	}
}
