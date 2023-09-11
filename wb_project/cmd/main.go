package main

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"time"
	"wb_project/pkg/client/postgresql"
	"wb_project/pkg/config"
	"wb_project/pkg/delivery"
	deliverydb "wb_project/pkg/delivery/delivery_db"
	"wb_project/pkg/handlers"
	"wb_project/pkg/items"
	itemsdb "wb_project/pkg/items/items_db"
	"wb_project/pkg/logging"
	"wb_project/pkg/pay"
	paydb "wb_project/pkg/pay/pay_db"
	"wb_project/pkg/user"
	"wb_project/pkg/user/db"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
)

type FullRepo struct {
	repoU user.Repository
	repoP pay.Repository
	repoI items.Repository
	repoD delivery.Repository
}

func main() {
	lg := logging.GetLogger()
	url := "nats://0.0.0.0:4222"
	lg.Info("Connect to nats")

	cfg := config.GetConfig()

	postrgeSQLClient, err := postgresql.NewClient(context.TODO(), cfg.Storage, 5)
	if err != nil {
		lg.Fatalf("%v", err)
	}
	repoU := db.NewRepository(postrgeSQLClient, lg)
	repoP := paydb.NewRepository(postrgeSQLClient, lg)
	repoI := itemsdb.NewRepository(postrgeSQLClient, lg)
	repoD := deliverydb.NewRepository(postrgeSQLClient, lg)

	repo := FullRepo{
		repoU: repoU,
		repoP: repoP,
		repoI: repoI,
		repoD: repoD,
	}

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
		// сделать валидацию данных
		err := Validation(msg.Data)
		if err != nil {
			lg.Error("Wrong Data")
		} else {
			err := goToBD(msg.Data, cache, repo)
			if err != nil {
				lg.Errorf("wrong data:%v with err:%v", msg.Data, err)
			}
		}
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

func Validation(msg []byte) error {
	return nil
}

func goToBD(data []byte, cache *user.UserCache, repo FullRepo) error {
	var u user.User

	err := json.Unmarshal(data, &u)
	if err != nil {
		return err
	}
	repo.repoU.Create(context.TODO(), &u)
	cache.Set(u.OrderUid, u, 3*time.Minute)
	return nil
}
