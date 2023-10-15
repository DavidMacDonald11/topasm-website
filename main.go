package main

import (
	"bytes"
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

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := Data{Answer: `Try this code:
        ; Fibonacci Numbers
        ; using #7 and #8 to hold last and current, respectively
        move 0 into #7
        move 1 into #8

        loop:
            ; print current value and newline
            move #8 into #0
            call printi
            move 10 into #0
            call printc

            ; last = current and current = current + last
            move #8 into #6 ; temp value
            add #7 into #8
            move #6 into #7

            ; loop while current is less than or equal to 610
            comp #8 with 610
            jumpLTE loop
            `}
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
        cmd := exec.Command("./topasm", file.Name())

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
