package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
)

type Data struct {
    Answer string
}

func main() {
    mux := http.NewServeMux()
    tmpl := template.Must(template.ParseFiles("views/index.html"))

    fs := http.FileServer(http.Dir("./static"))
    mux.Handle("/static/", http.StripPrefix("/static/", fs))

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := Data{Answer: ""}
        tmpl.Execute(w, data)
    })

    mux.HandleFunc("/interpret", func(w http.ResponseWriter, r *http.Request) {
        asm := r.PostFormValue("asm")
        cmd := exec.Command("./topasm", fmt.Sprintf("--text=%s", asm))

        var buf bytes.Buffer
        cmd.Stdout = &buf
        cmd.Stderr = &buf

        cmd.Run()
        w.Write([]byte(buf.String()))
    })

    log.Println("Starting server...")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
