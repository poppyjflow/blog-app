package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strconv"
    "fmt"

    "github.com/gorilla/mux"
    "github.com/rs/cors"
    _ "github.com/lib/pq"
)

type App struct {
    Router *mux.Router
    DB     *sql.DB
}

func main() {
    app := App{}
    app.Initialize(
        os.Getenv("APP_DB_USERNAME"),
        os.Getenv("APP_DB_PASSWORD"),
        os.Getenv("APP_DB_NAME"))

    app.Run(":8080")
}

func (app *App) Initialize(user, password, dbname string) {
    connectionString :=
        fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

    var err error
    app.DB, err = sql.Open("postgres", connectionString)
    if err != nil {
        log.Fatal(err)
    }

    app.Router = mux.NewRouter()
    app.initializeRoutes()
}

func (app *App) Run(addr string) {
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type", "Authorization"},
    })

    handler := c.Handler(app.Router)
    log.Fatal(http.ListenAndServe(addr, handler))
}

func (app *App) initializeRoutes() {
    app.Router.HandleFunc("/posts", app.getPosts).Methods("GET")
    app.Router.HandleFunc("/posts", app.createPost).Methods("POST")
    app.Router.HandleFunc("/posts/{id:[0-9]+}", app.getPost).Methods("GET")
    app.Router.HandleFunc("/posts/{id:[0-9]+}", app.updatePost).Methods("PUT")
    app.Router.HandleFunc("/posts/{id:[0-9]+}", app.deletePost).Methods("DELETE")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

type post struct {
    ID      int    `json:"id"`
    UserID  int    `json:"user_id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}

func (app *App) getPosts(w http.ResponseWriter, r *http.Request) {
    count, _ := strconv.Atoi(r.FormValue("count"))
    start, _ := strconv.Atoi(r.FormValue("start"))

    if count > 10 || count < 1 {
        count = 10
    }
    if start < 0 {
        start = 0
    }

    posts, err := getPosts(app.DB, start, count)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, posts)
}

func (app *App) createPost(w http.ResponseWriter, r *http.Request) {
    var p post
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&p); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := p.createPost(app.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, p)
}

func (app *App) getPost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid post ID")
        return
    }

    p := post{ID: id}
    if err := p.getPost(app.DB); err != nil {
        switch err {
        case sql.ErrNoRows:
            respondWithError(w, http.StatusNotFound, "Post not found")
        default:
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }

    respondWithJSON(w, http.StatusOK, p)
}

func (app *App) updatePost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid post ID")
        return
    }

    var p post
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&p); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()
    p.ID = id

    if err := p.updatePost(app.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, p)
}

func (app *App) deletePost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid post ID")
        return
    }

    p := post{ID: id}
    if err := p.deletePost(app.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func getPosts(db *sql.DB, start, count int) ([]post, error) {
    rows, err := db.Query(
        "SELECT id, user_id, title, content FROM posts LIMIT $1 OFFSET $2",
        count, start)

    if err != nil {
        return nil, err
    }

    defer rows.Close()

    posts := []post{}

    for rows.Next() {
        var p post
        if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content); err != nil {
            return nil, err
        }
        posts = append(posts, p)
    }

    return posts, nil
}

func (p *post) getPost(db *sql.DB) error {
    return db.QueryRow("SELECT user_id, title, content FROM posts WHERE id=$1",
        p.ID).Scan(&p.UserID, &p.Title, &p.Content)
}

func (p *post) updatePost(db *sql.DB) error {
    _, err :=
        db.Exec("UPDATE posts SET title=$1, content=$2 WHERE id=$3",
            p.Title, p.Content, p.ID)

    return err
}

func (p *post) deletePost(db *sql.DB) error {
    _, err := db.Exec("DELETE FROM posts WHERE id=$1", p.ID)

    return err
}

func (p *post) createPost(db *sql.DB) error {
    err := db.QueryRow(
        "INSERT INTO posts(user_id, title, content) VALUES($1, $2, $3) RETURNING id",
        p.UserID, p.Title, p.Content).Scan(&p.ID)

    if err != nil {
        return err
    }

    return nil
}