package main

import (
	"time"
	"fmt"
	"context"
	"net/http"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

func main() {
	r := chi.NewRouter();

	// Base middleware stack
	r.Use(middleware.RequestID) // all requests now get an id
	r.Use(middleware.RealIP) //
	r.Use(middleware.Logger) // log requests start and end time
	r.Use(middleware.Recoverer) //

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte ("Hello"))
	})

	r.Route("/User", func(r chi.Router){
		r.Get("/", func(w http.ResponseWriter, r *http.Request){
			w.Write([]byte("Get your users here"))
		})
		r.Post("/",  createUser)
		r.Route("/:userId", func(r chi.Router){
			r.Use(SetupUser)
			r.Get("/", getUser)
		})
	})

	r.Route("/List", func(r chi.Router){
		r.Get("/", func(w http.ResponseWriter, r *http.Request){
			w.Write([]byte("Get your lists here"))
		})
		r.Route("/:listId", func(r chi.Router){
			r.Use(List)
			r.Get("/", getList)
		})

	})

	r.Route("/Item", func(r chi.Router){
		r.Get("/", func(w http.ResponseWriter, r *http.Request){
			w.Write([]byte("Get your items here"))
		})
		r.Route("/:itemId", func(r chi.Router){
			r.Use(Item)
			r.Get("/", getItem)
		})
	})

	http.ListenAndServe(":3000", r)
}

type User struct{
	userId string
	userName string
}

func SetupUser(next http.Handler) http.Handler{
	return http.HandlerFunc (func (w http.ResponseWriter, r *http.Request){
		user := User{}
		user.userId = chi.URLParam(r, "userId")
		user.userName = chi.URLParam(r, "userName")
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUser(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	user := ctx.Value("user").(*User)
	w.Write([]byte(fmt.Sprintf("User:%s", user.userId)))
}

func createUser(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	user := ctx.Value("user").(*User)
	w.Write([]byte(fmt.Sprintf("User:%s, Username:%s", user.userId, user.userName)))
}

func List(next http.Handler) http.Handler{
	return http.HandlerFunc (func (w http.ResponseWriter, r *http.Request){
		listId := chi.URLParam(r, "listId")
		ctx := context.WithValue(r.Context(), "listId", listId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getList(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	w.Write([]byte(fmt.Sprintf("List:%s", ctx.Value("listId"))))
}

func Item(next http.Handler) http.Handler{
	return http.HandlerFunc (func (w http.ResponseWriter, r *http.Request){
		itemId := chi.URLParam(r, "itemId")
		ctx := context.WithValue(r.Context(), "itemId", itemId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getItem(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	w.Write([]byte(fmt.Sprintf("Item:%s", ctx.Value("itemId"))))
}
