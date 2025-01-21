package main

import (
	"database/sql"
	"encoding/gob"
	"final-project/data"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

func main() {
	// connect to the database
	db := initDB()

	// create sessions
	session := initSession()

	// create logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create channels

	// create waitgroup
	wg := sync.WaitGroup{}

	// set up the application config
	app := Config{
		Session:       session,
		DB:            db,
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		Wait:          &wg,
		Models:        data.New(db),
		ErrorChan:     make(chan error),
		ErrorChanDone: make(chan bool),
	}

	// set up mail
	app.Mailer = app.createMail()
	go app.listenForMail()

	// listen for signals
	go app.listenForShutdown()

	// listen for error
	go app.listenForError()

	// listen for web connections
	app.serve()
}

// listenForError runs concurrently and wait for error to error channel to print it out
func (app *Config) listenForError() {
	for {
		select {
		case err := <-app.ErrorChan:
			app.ErrorLog.Println(err)
		case <-app.ErrorChanDone:
			return
		}
	}
}

// serve is used to start the web server
func (app *Config) serve() {
	// start http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	app.InfoLog.Println("Starting web server...")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// initDB sets up connection with database
func initDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to database!")
	}
	return conn
}

// connectToDB is used to "try to" connect to database
// if connection is failed, it will try 10 times before exit
func connectToDB() *sql.DB {
	counts := 0

	dsn := os.Getenv("DSN")
	if dsn == "" {
		fmt.Println("Cannot get DSN")
		dsn = "host=localhost port=5432 dbname=concurrency user=postgres password=24072001do sslmode=disable timezone=UTC connect_timeout=5"
	}

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println(err.Error())
			log.Println("Postgres not yet ready...")
		} else {
			log.Println("Connected to database!")
			return connection
		}

		if counts > 10 {
			return nil
		}
		log.Println("Backing off for 1 second")
		time.Sleep(1 * time.Second)
		counts++
		continue
	}
}

// openDB is used to connect to database
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// initSession sets up a session, using Redis for session store
func initSession() *scs.SessionManager {
	gob.Register(data.User{})

	// set up session
	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session
}

// initRedis returns a pool of connections to Redis
func initRedis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			redisdata := os.Getenv("REDIS")
			fmt.Println("REDIS = ", redisdata)
			if redisdata == "" {
				fmt.Println("Cannot get REDIS")
				redisdata = "127.0.0.1:6379"
			}
			return redis.Dial("tcp", redisdata)
		},
	}
	return redisPool
}

// listenForShutdown is a goroutine which concurrency running with serve()
// it's waiting for signals to stop the application
func (app *Config) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0)
}

// shutdown will check, close and wait for other sections to shutdown
// before shutdown the application
func (app *Config) shutdown() {
	// perform any cleanup tasks
	app.InfoLog.Println("Would run cleanup tasks...")

	// block until WaitGroup is empty
	app.Wait.Wait()

	app.Mailer.DoneChan <- true
	app.ErrorChanDone <- true

	app.InfoLog.Println("Closing channels and shutting down application...")
	close(app.Mailer.MailerChan)
	close(app.Mailer.ErrorChan)
	close(app.Mailer.DoneChan)
	close(app.ErrorChan)
	close(app.ErrorChanDone)
}

// createMail setup and return a Mail
func (app *Config) createMail() Mail {
	// create channels
	errorChan := make(chan error)
	mailerChan := make(chan Message, 1000)
	mailerDoneChan := make(chan bool)

	m := Mail{
		Domain:      "localhost",
		Host:        "localhost",
		Port:        1025,
		Encryption:  "none",
		FromAddress: "info@mycompany.com",
		FromName:    "Info",
		Wait:        app.Wait,
		ErrorChan:   errorChan,
		MailerChan:  mailerChan,
		DoneChan:    mailerDoneChan,
	}

	return m
}
