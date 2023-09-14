package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"wb_project/pkg/cashe"
	"wb_project/pkg/client/postgresql"
	"wb_project/pkg/config"
	"wb_project/pkg/fullrepo"

	deliverydb "wb_project/pkg/delivery/delivery_db"
	"wb_project/pkg/handlers"

	itemsdb "wb_project/pkg/items/items_db"
	"wb_project/pkg/logging"

	paydb "wb_project/pkg/pay/pay_db"
	"wb_project/pkg/user"
	"wb_project/pkg/user/db"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
)

func main() {
	lg := logging.GetLogger()
	url := "nats://0.0.0.0:4222"

	cfg := config.GetConfig()
	cfg.Storage.Lg = lg

	lg.Info("Connect to nats")
	nc, err := nats.Connect(url)
	if err != nil {
		lg.Info("Can not connect to nats")
		lg.Fatal(err)
	}
	lg.Info("Succesfull connect to nats")
	postrgeSQLClient, err := postgresql.NewClient(context.Background(), cfg.Storage, 5)
	if err != nil {
		lg.Fatalf("%v", err)
	}
	lg.Info("Succesfull connect to postgresql")
	repoU := db.NewRepository(postrgeSQLClient, lg)
	repoP := paydb.NewRepository(postrgeSQLClient, lg)
	repoI := itemsdb.NewRepository(postrgeSQLClient, lg)
	repoD := deliverydb.NewRepository(postrgeSQLClient, lg)

	repo := fullrepo.FullRepo{
		RepoU: repoU,
		RepoP: repoP,
		RepoI: repoI,
		RepoD: repoD,
	}

	defer nc.Close()

	cache := cashe.NewCache(2*time.Minute, 3*time.Minute)
	cache.LoadFile(&repo, 2*time.Minute)
	lg.Info("Registration subscribe")

	sub, err := nc.Subscribe("events.*", func(msg *nats.Msg) {
		err := goToBD(msg.Data, cache, &repo)
		if err != nil {
			lg.Errorf("wrong data:%s with err:%v", string(msg.Data), err)
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
		DB:       &repo,
	}
	r := mux.NewRouter()
	addr := ":9080"
	r.HandleFunc("/", userHandler.Handler).Methods("GET")
	lg.Infof("start server with addres %v", addr)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		lg.Fatal("can not start server", err)
	}
}

func goToBD(data []byte, cache *cashe.UserCache, repo *fullrepo.FullRepo) error {
	var u user.User
	if data[0] != '{' {
		a := strings.IndexByte(string(data), '{')
		data = data[a:]
	}
	fmt.Println(string(data))
	err := json.Unmarshal(data, &u)
	if err != nil {
		return err
	}
	fmt.Println(u)
	err = repo.RepoU.Create(context.Background(), &u)
	if err != nil {
		return err
	}
	for _, i := range u.Items {
		err = repo.RepoI.Create(context.Background(), &i)
		if err != nil {
			return err
		}
	}
	u.Deliv.UserId = u.ID
	err = repo.RepoD.Create(context.Background(), &u.Deliv)
	if err != nil {
		return err
	}
	u.Payment.UserId = u.ID
	err = repo.RepoP.Create(context.Background(), &u.Payment)
	if err != nil {
		return err
	}
	cache.Set(u.OrderUid, u, 3*time.Minute)
	return nil
}
