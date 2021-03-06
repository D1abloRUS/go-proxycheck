package models

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/D1abloRUS/proxycheck-server/config"

	"github.com/julienschmidt/httprouter"
)

const (
	get = "GET"
)

//AllProxy get
func AllProxy(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		if r.Method != get {
			http.Error(w, http.StatusText(405), 405)
			return
		}

		bks, err := AllProxyReq(env.DB)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err = json.NewEncoder(w).Encode(bks); err != nil {
			w.WriteHeader(500)
		}
	}
}

//AllCountry get
func AllCountry(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		if r.Method != get {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		bks, err := AllCountryReq(env.DB)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err = json.NewEncoder(w).Encode(bks); err != nil {
			w.WriteHeader(500)
		}
	}
}

//FilterCountry get
func FilterCountry(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id, _ := strconv.Atoi(p.ByName("id"))

		if r.Method != get {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		bks, err := FilterCountryReq(env.DB, id)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if len(bks) == 0 {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err = json.NewEncoder(w).Encode(bks); err != nil {
			w.WriteHeader(500)
		}
	}
}

//FilterProxy get
func FilterProxy(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id, _ := strconv.Atoi(p.ByName("id"))

		if r.Method != get {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		bks, err := FilterProxyReq(env.DB, id)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if len(bks) == 0 {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err = json.NewEncoder(w).Encode(bks); err != nil {
			w.WriteHeader(500)
		}
	}
}

//UpdateProxyStatus post json
func UpdateProxyStatus(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id, _ := strconv.Atoi(p.ByName("id"))

		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}

		err := UpdateStatus(env.DB, id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(200)
	}
}

//AddProxy post json
func AddProxy(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var p ProxyRespone
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		err = AddToBase(env.DB, p.Country, p.IP, p.Port, p.Respone, p.Status)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(200)
	}
}
