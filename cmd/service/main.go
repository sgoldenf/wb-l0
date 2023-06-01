package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/form"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sgoldenf/wb_l0/internal/application"
	"github.com/sgoldenf/wb_l0/internal/model"
	"github.com/sgoldenf/wb_l0/internal/templates"
)

var (
	addr  *string
	dbURL *string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("WARNING: No .env file found")
	}
	addr = flag.String("addr", ":8080", "HTTP network address")
	dbName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbURL = flag.String(
		"dbURL",
		"postgres://"+user+":"+password+"@localhost:5432/"+dbName,
		"PostgresSQL database URL",
	)
	flag.Parse()
}

func main() {
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Println("Template Cache created")
	db, err := dbConn()
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Println("Connected to database:", *dbURL)
	app := &application.Application{
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		TemplateCache: templateCache,
		FormDecoder:   form.NewDecoder(),
		Orders:        &model.OrderModel{Pool: db},
	}

	if err := app.InitOrdersCache(); err != nil {
		app.ErrorLog.Fatal(err)
	}

	if err := app.InitStanConnection(); err != nil {
		app.ErrorLog.Fatal(err)
	}

	if err := app.InitStanSubscription(); err != nil {
		app.ErrorLog.Fatal(err)
	}

	server := &http.Server{
		Addr:    *addr,
		Handler: app.Routes(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.ErrorLog.Fatal(err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	app.InfoLog.Println("Gracefully shutting down")
	server.Shutdown(context.Background())
	app.Shutdown()
}

func dbConn() (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), *dbURL)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return conn, err
}
