package main

import (
	"log"
	"net/http"
	"time"

	"alilacream/ecom/config"
	"alilacream/ecom/internal/db"
	"alilacream/ecom/store"

	"github.com/gin-gonic/gin"
)

type application struct {
	config Config
	store  *store.Storage
	router *gin.Engine
}

type Config struct {
	addr string
	db   *config.DBConfig
}

func Setup() *application {
	env_vars, err := config.AllEnvs()
	if err != nil {
		panic("Environement setup is not configured at all!")
	}

	dbConf := config.DBConfig{
		DSN:                env_vars.MainDB,
		MaxOpenConnections: config.ParseInt(env_vars.MaxOpenConnections),
		MaxIdleConnections: config.ParseInt(env_vars.MaxIdleConnections),
		MaxIdleTime:        env_vars.MaxIdleTime,
	}

	rdb, err := db.RedisNew(&dbConf)
	if err != nil {
		panic("Redis Database Connection Hasn't been configured")
	}

	pdb, err := db.PSQLNew(&dbConf)
	if err != nil {
		panic("Postgresql Database Connection Hasn't been configured")
	}
	defer pdb.Close()
	log.Println("Connected to the Database")

	store := store.NewPQStorage(pdb)
	return &application{
		config: Config{
			db:   &dbConf,
			addr: ":8080",
		},
		store: &store,
	}
}

func (a *application) routes(r *gin.Default) {
	// pub
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Api is working")
	})
}

// running the application, the core method for serving
func (a *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:    a.config.addr,
		Handler: mux,
		// security enhancing args in server interface
		WriteTimeout: time.Second * 30, // max timeout to write response to the client
		ReadTimeout:  time.Second * 10, // max timeout to read the request from the client
		IdleTimeout:  time.Minute,      // keeps the connection alive only in one minute
	}
	log.Println("Server listening in port: ", a.config.addr)
	return srv.ListenAndServe()
}
