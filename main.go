package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
        file, err := ioutil.TempFile("", "asm-*.topasm")

        if err != nil {
            log.Println(err)
            w.Write([]byte("Post Error 1"))
            return
        }

        defer os.Remove(file.Name())

        _, err = file.WriteString(asm)

        if err != nil {
            log.Println(err)
            w.Write([]byte("Post Error 2"))
            return
        }

        file.Close()
        cmd := exec.Command("./topasm", fmt.Sprintf("--file=%s", file.Name()))

        var buf bytes.Buffer
        cmd.Stdout = &buf
        cmd.Stderr = &buf

        cmd.Run()
        w.Write([]byte(buf.String()))
    })

    port := os.Getenv("PORT")
    if port == "" { port = "8080" }

    log.Println("Starting server...")
    log.Fatal(http.ListenAndServe(":" + port, mux))
}
