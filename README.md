# Minirouter

`go get github.com/ofux/minirouter`

Minirouter just adds support for middlewares on top of [httprouter](https://github.com/julienschmidt/httprouter).
It was inspired by [hitch](https://github.com/nbio/hitch), but minirouter also supports clean definition of sub-paths so that routes can
be grouped together.

It is designed to use only standard interfaces such as http.Handler and http.HandlerFunc to define handlers and middlewares, and exclusively uses [context](https://golang.org/pkg/context/)
to retrieve query parameters. It does **not** pollute your code with weird custom stuff. Only standard interfaces.
