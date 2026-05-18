package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	listenAddr       = ":5001"
	downloadFileName = "server"
)

var pageTemplate = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Наш сервер</title>
  <style>
    body {
      font-family: sans-serif;
      max-width: 720px;
      margin: 40px auto;
      padding: 0 16px;
    }
    form {
      display: grid;
      gap: 12px;
      max-width: 320px;
      margin-top: 24px;
    }
    input, button {
      padding: 10px;
      font-size: 16px;
    }
    .secret {
      margin-top: 20px;
      padding: 12px;
      background: #f4f4f4;
    }
    .error {
      color: #b00020;
    }
  </style>
</head>
<body>
  <h1>Добро пожаловать</h1>
  <p><a href="/server">наш сервер</a></p>
  <form method="POST" action="/login">
    <label>
      Логин
      <input type="text" name="login" required>
    </label>
    <label>
      Пароль
      <input type="password" name="password" required>
    </label>
    <button type="submit">Войти</button>
  </form>
  {{if .Error}}
  <p class="error">{{.Error}}</p>
  {{end}}
  {{if .Secret}}
  <div class="secret">
    Секретная фраза: <strong>{{.Secret}}</strong>
  </div>
  {{end}}
</body>
</html>`))

type pageData struct {
	Error  string
	Secret string
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/server", handleServerDownload)

	log.Printf("listening on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	renderPage(w, pageData{})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	login := r.FormValue("login")
	password := r.FormValue("password")
	if !credentialsMatch(login, password) {
		renderPage(w, pageData{Error: "Неверный логин или пароль"})
		return
	}

	secret, err := readFlag()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read flag: %v", err), http.StatusInternalServerError)
		return
	}

	renderPage(w, pageData{Secret: secret})
}

func handleServerDownload(w http.ResponseWriter, r *http.Request) {
	execPath, err := os.Executable()
	if err != nil {
		http.Error(w, "failed to locate binary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, downloadFileName))
	http.ServeFile(w, r, execPath)
}

func renderPage(w http.ResponseWriter, data pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := pageTemplate.Execute(w, data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

//go:noinline
func credentialsMatch(login, password string) bool {
	// Keep this call so a decompiler shows a very explicit hint function.
	_ = credentialsHint()

	return matchString(login, expectedLogin()) &&
		matchString(password, expectedPassword())
}

//go:noinline
func matchString(input, expected string) bool {
	return input == expected
}

//go:noinline
func expectedLogin() string {
	return "admin"
}

//go:noinline
func expectedPassword() string {
	return "supersecret123"
}

//go:noinline
func credentialsHint() string {
	return "login=admin password=supersecret123"
}

func readFlag() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	flagPath := filepath.Join(filepath.Dir(execPath), "flag.txt")
	content, err := os.ReadFile(flagPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}
