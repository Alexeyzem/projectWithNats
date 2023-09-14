package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strings"
	"wb_project/pkg/cashe"
	"wb_project/pkg/fullrepo"
	"wb_project/pkg/logging"
)

type UserHandler struct {
	Tmpl     *template.Template
	Lg       *logging.Logger
	UserRepo *cashe.UserCache
	DB       *fullrepo.FullRepo
}

func (h *UserHandler) Handler(w http.ResponseWriter, r *http.Request) {
	order := r.URL.Query().Get("order_uid")
	h.Lg.Infof("uuid: %s", order)
	err := h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		Uid string
	}{
		order,
	})
	if err != nil {
		h.Lg.Errorf("error with template: %v", err)
	}
	u, ok := h.UserRepo.Get(order)
	if ok {
		out, err := json.Marshal(u)
		if err != nil {
			h.Lg.Errorf("can not marshal data:%v", err)
			w.Write([]byte(`internal error with this uid: "` + order + `", sorry. We are working to fix`))
			return
		}
		h.Lg.Info("Data from cache.")
		w.Write(out)
		return
	} else {
		err := Validation(order)
		if err != nil {
			w.Write([]byte("Wrong data. UUID must contain 19 characters and should not contain SQL-query. Try again!"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		u, err := h.DB.RepoU.FindOne(context.Background(), order)
		if err != nil {
			w.Write([]byte("We dont have info about this uuid:" + order))
			w.WriteHeader(http.StatusBadRequest)
			h.Lg.Error(err)
			return
		}
		if u.OrderUid != order {
			w.Write([]byte("We dont have info about this uuid:" + order))
			return
		}
		p, errP := h.DB.RepoP.FindOne(context.Background(), u.ID)
		if errP != nil {
			return
		}
		u.Payment = p
		i, err := h.DB.RepoI.FindAllOfOneUser(context.Background(), u.TrackNumber)
		if err != nil {
			return
		}
		u.Items = i
		d, err := h.DB.RepoD.FindOne(context.Background(), u.ID)
		u.Deliv = d
		if err != nil {
			w.Write([]byte("We dont have full info about this user."))
			return
		}
		out, err := json.Marshal(u)
		if err != nil {
			h.Lg.Errorf("can not marshal data:%v", err)
			w.Write([]byte(`internal error with this UID: "` + order + `", sorry. We are working to fix`))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.Lg.Info("data from db")
		w.Write(out)
	}
}

func Validation(data string) error {
	var err error = nil
	if len(data) != 19 {
		err = errors.New("wrong data")
	}
	d := strings.Trim(data, " ")
	if len(d) != 19 {
		err = errors.New("wrong data")
	}
	d = strings.ToLower(data)
	if strings.Contains(d, "drop table") {
		err = errors.New("wrong data")
	}
	return err
}
