package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
	"wb_project/pkg/handlers"
	"wb_project/pkg/logging"
	"wb_project/pkg/user"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
)

func main() {
	lg := logging.GetLogger()
	url := "nats://0.0.0.0:4222"
	lg.Info("Connect to nats")
	nc, err := nats.Connect(url)
	if err != nil {
		lg.Info("Can not connect to nats")
		lg.Fatal(err)
	}
	defer nc.Close()
	cache := user.NewCache(2*time.Minute, 3*time.Minute)
	cache.LoadFile("some path to DB")
	lg.Info("Registration subscribe")
	sub, err := nc.Subscribe("events.*", func(msg *nats.Msg) {
		fmt.Println(string(msg.Data))
		// запись данных в БД и кэш, проверка на корректность.
		goToBD(msg.Data, cache)
	})
	if err != nil {
		lg.Error("Can not registrate subscribe")
	}

	defer sub.Unsubscribe()
	templates := template.Must(template.ParseGlob("./static/*"))
	userHandler := &handlers.UserHandler{
		Tmpl:     templates,
		Lg:       lg,
		UserRepo: cache,
	}
	r := mux.NewRouter()
	addr := ":8080"
	r.HandleFunc("/", userHandler.Handler).Methods("GET")
	lg.Infof("start server with addres %v", addr)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		lg.Fatal("can not start server")
	}
}

func goToBD(data []byte, cache *user.UserCache) error {
	var u user.User
	err := json.Unmarshal(data, &u)
	if err != nil {
		return err
	}
	cache.Set(u.OrderUid, u, 3*time.Minute)

	return nil
}
