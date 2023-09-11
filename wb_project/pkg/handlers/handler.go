package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"wb_project/pkg/logging"
	"wb_project/pkg/user"
)

type UserHandler struct {
	Tmpl     *template.Template
	Lg       *logging.Logger
	UserRepo *user.UserCache
}

func (h *UserHandler) Handler(w http.ResponseWriter, r *http.Request) {
	order := r.URL.Query().Get("order_uid")
	err := h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		uid string
	}{
		uid: order,
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
		}
		w.Write(out)
	} else {
		find := false
		// обращение к БД
		var data user.User
		if find {
			out, err := json.Marshal(data)
			if err != nil {
				h.Lg.Errorf("can not marshal data:%v", err)
				w.Write([]byte(`internal error with this UID: "` + order + `", sorry. We are working to fix`))
			}
			w.Write(out)
		} else {
			w.Write([]byte(`We dont have info about uid: "` + order + `". Check that you entered correct UID`))
		}
	}
}
