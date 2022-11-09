package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"awesome.forstes.go/internal/models"
	"awesome.forstes.go/internal/storage"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/minio/minio-go"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	pictures      *models.PictureModel
	objStorage    storage.ObjectStorage
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	connStr := flag.String("connStr", "postgres://postgres:12345@localhost:5432/go-awesome", "Postgres DB connection")

	objStorage := flag.String("objStorageConnStr", "localhost:9000/admin123/admin123", "Object storage address and credentials")
	objStorageArgs := strings.Split(*objStorage, "/")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	pool, err := pgxpool.Connect(context.Background(), *connStr)
	if err != nil {
		errorLog.Fatal(err)
		os.Exit(1)
	}
	defer pool.Close()

	minioClient, err := minio.New(objStorageArgs[0], objStorageArgs[1], objStorageArgs[2], false)
	if err != nil {
		errorLog.Fatal(err)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		pictures:      &models.PictureModel{DB: pool},
		objStorage:    &storage.MinioStore{Client: minioClient},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     "127.0.0.1" + *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %v", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
