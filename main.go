package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shaineminkyaw/road-system-background/config"
	"github.com/shaineminkyaw/road-system-background/controller"
	"github.com/shaineminkyaw/road-system-background/ds"
)

func main() {
	//
	r := gin.Default()
	conf := config.Init()
	db := ds.ConnectToDB(conf.SQL.Host, conf.SQL.Port, conf.SQL.DB, conf.SQL.User, conf.SQL.Password)
	e, err := ds.NewRBAC(db)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	h := controller.Inject(r, e, db)

	server := http.Server{
		Addr:           fmt.Sprintf("%v:%v", conf.APP.Host, conf.APP.Port),
		Handler:        h.R,
		ReadTimeout:    300 * time.Second,
		WriteTimeout:   300 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		log.Printf("Starting server on %v:%v\n", conf.APP.Host, conf.APP.Port)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Println("failed to starting server....")
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	//block until ctrl+C is come
	<-c

	err = server.Shutdown(context.Background())
	if err != nil {
		log.Println("error on server closing....")
		log.Println(err.Error())
		os.Exit(1)
	}
	log.Println("server closing successfully....")

}
