package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"

	authboss "gopkg.in/authboss.v1"
	_ "gopkg.in/authboss.v1/auth"
)

var (
	ab = authboss.New()
)

func main() {
	initAuth()
	mux := mux.NewRouter()
	mux.Handle("/auth", ab.NewRouter())
	// mux.Handle("/", index())

	err := mux.Walk(gorillaWalkFn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("serving up on 8080")
	log.Println(http.ListenAndServe(":8080", mux))
}

func initAuth() {
	ab.MountPath = "/auth"
	ab.LogWriter = os.Stdout
	ab.ViewsPath = "ab_views"

	ab.Storer = &MemStorer{}
	ab.SessionStoreMaker = NewSessionStorer
	ab.CookieStoreMaker = NewCookieStorer
	ab.XSRFName = "csrf_token"
	ab.XSRFMaker = func(_ http.ResponseWriter, r *http.Request) string {
		return csrf.Token(r)
	}
	ab.LayoutDataMaker = layoutData
	if err := ab.Init(); err != nil {
		// Handle error, don't let program continue to run
		log.Fatalln(err)
	}
}

func layoutData(w http.ResponseWriter, r *http.Request) authboss.HTMLData {
	userInter, _ := ab.CurrentUser(w, r) //FIXME: error handling

	return authboss.HTMLData{
		"loggedin":               userInter != nil,
		authboss.FlashSuccessKey: ab.FlashSuccess(w, r),
		authboss.FlashErrorKey:   ab.FlashError(w, r),
	}
}

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hey")
	})
}

func gorillaWalkFn(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	path, _ := route.GetPathTemplate()
	log.Printf(path)
	return nil
}
