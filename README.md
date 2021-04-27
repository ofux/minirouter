# Minirouter

`go get github.com/ofux/minirouter`

Minirouter just adds support for middlewares on top of [httprouter](https://github.com/julienschmidt/httprouter). It was
inspired by [hitch](https://github.com/nbio/hitch), but minirouter also supports clean definition of sub-paths so that
routes can be grouped together.

It is designed to use only standard interfaces such as http.Handler and http.HandlerFunc to define handlers and
middlewares, and exclusively uses [context](https://golang.org/pkg/context/)
to retrieve query parameters. It does **not** pollute your code with weird custom stuff. Only standard interfaces.

## Usage

```go
package main

import (
	"fmt"
	"net/http"
	"log"

	"github.com/ofux/minirouter"
	"github.com/rs/cors"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, %s!\n", minirouter.Params(r).ByName("name"))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, user with ID %s!\n", minirouter.Params(r).ByName("id"))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	//TODO: do your thing
}

func checkIsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: check logged user is admin
		next.ServeHTTP(w, r)
	})
}

func main() {
	mr := minirouter.New()
	mr = mr.WithMiddleware(cors.Default().Handler)

	mr.GET("/", Index)
	mr.GET("/hello/:name", Hello)

	mrAdmin := mr.WithBasePath("/admin").WithMiddleware(checkIsAdmin)
	mrAdmin.GET("/users/:id", GetUser)
	mrAdmin.POST("/users", CreateUser)

	log.Fatal(http.ListenAndServe(":8080", mr))
}
```
