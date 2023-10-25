# HTMX Demo

A simple demo using [HTMX](https://htmx.org/docs/), [Gorilla](https://github.com/gorilla/mux) and [Slog](https://golang.org/x/exp/slog).

## Getting Started

```bash
$ go run *.go 
time=2023-05-30T11:16:15.487+10:00 level=INFO source=/Users/nickglynn/Projects/htmx-demo/server.go:41 msg="Starting server..." SERVER=http://localhost:8080
```

## URLs

- [Static Index](http://localhost:8080/static/index.html)
- [Template Fragment Index](http://localhost:8080/)
