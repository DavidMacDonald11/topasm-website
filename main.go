package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type Code struct {
    Asm string
}

func main() {
    app := fiber.New(fiber.Config{Views: html.New("./views", ".html")})
    app.Static("/", "./public")

    app.Get("/", getRoot)
    app.Get("/example/:id", getExample)
    app.Post("/interpret", postInterpret)

    port := os.Getenv("PORT")
    if port == "" { port = "3000" }

    log.Println("Starting server...")
    log.Fatal(app.Listen(":" + port))
}

func getRoot(c *fiber.Ctx) error {
    return c.Render("index", fiber.Map{})
}

func getExample(c *fiber.Ctx) error {
    id := c.Params("id")
    name := fmt.Sprintf("./examples/%s.topasm", id)
    content, err := ioutil.ReadFile(name)

    if err != nil {
        log.Println(err)
        return c.Render("input", fiber.Map{})
    }

    return c.Render("input", fiber.Map{"Text": string(content)})
}

func postInterpret(c *fiber.Ctx) error {
    code := new(Code)
    err := c.BodyParser(code)

    if err != nil {
        return renderInternalError(c, err)
    }

    file, err := ioutil.TempFile("", "asm-*.topasm")

    if err != nil {
        return renderInternalError(c, err)
    }

    defer os.Remove(file.Name())

    _, err = file.WriteString(code.Asm)

    if err != nil {
        return renderInternalError(c, err)
    }

    file.Close()
    cmd := exec.Command("./topasm", file.Name())

    var buf bytes.Buffer
    cmd.Stdout = &buf
    cmd.Stderr = &buf

    cmd.Run()
    return c.Render("result", fiber.Map{"Output": buf.String()})
}

func renderInternalError(c *fiber.Ctx, err error) error {
    log.Println(err)
    return c.Render("result", fiber.Map{"Error": "Internal server error"})
}
