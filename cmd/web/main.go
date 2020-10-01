package  main

import (
    "database/sql"
    "flag"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"

    "vy.dao/snippetbox/pkg/models/postgres"

    _ "github.com/lib/pq"
)

type application struct {
    errorLog   *log.Logger
    infoLog    *log.Logger
    snippets   *postgres.SnippetModel
    templateCache map[string]*template.Template
}

const (
  host     = "localhost"
  port     = 5432
  user     = "postgres"
  dbname   = "snippetbox"
)

func main() {
    addr := flag.String("addr", ":4000", "HTTP Network address")
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
      "dbname=%s sslmode=disable",
      host, port, user, dbname)
    dsn := flag.String("dsn", psqlInfo, "PostgreSQL connection string")

    flag.Parse()

    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

    db, err := openDB(*dsn)
    if err != nil {
        errorLog.Fatal(err)
    }

    defer db.Close()

    templateCache, err := newTemplateCache("./ui/html")
    if err != nil {
        errorLog.Fatal(err)
    }

    app := &application{
        errorLog:   errorLog,
        infoLog:    infoLog,
        snippets:   &postgres.SnippetModel{DB: db},
        templateCache: templateCache,
    }

    srv := &http.Server{
        Addr:       *addr,
        ErrorLog:   errorLog,
        Handler:    app.routes(),
    }

    infoLog.Printf("Starting server on %s", *addr)
    err = srv.ListenAndServe()
    errorLog.Fatal(err)
}


// The openDB() function wraps sql.Open() and returns a sql.DB connection pool // for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return db, nil
}