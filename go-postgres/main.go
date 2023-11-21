package main

import (
    "fmt"
    "go-postgres/router"
    "log"
    "net/http"
)

func main() {
    r := router.Router()
    // fs := http.FileServer(http.Dir("build"))
    // http.Handle("/", fs)
    fmt.Println("Iniciando servidor na porta 8080...")

    log.Fatal(http.ListenAndServe(":8080", r))
}