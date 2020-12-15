package server

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var closeCh = make(chan struct{})

func s1() error {
	mux := http.ServeMux{}
	mux.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("双月同天"))
	})

	server := http.Server{
		Addr:":9080",
		Handler: &mux,
	}

	go func() {
		<- closeCh
		log.Println("server1 received close signal")
		timeout, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()
		server.Shutdown(timeout)
		log.Println("server1 has shutdown")
	}()

	log.Println("server1: listening on port :9080....")
	return server.ListenAndServe()
}

func s2() error {
	mux := http.ServeMux{}
	mux.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("九星连珠"))
	})

	server := http.Server{
		Addr:":9090",
		Handler: &mux,
	}

	go func() {
		<- closeCh
		log.Println("server2 received close signal")
		timeout, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()
		server.Shutdown(timeout)
		log.Println("server2 has shutdown")
	}()
	log.Println("server2: listening on port :9090....")
	return server.ListenAndServe()
}

func sigRegister() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Print(err)
			}
		}()

		si := <- sig
		log.Println("Got signal :", si)
		if closeCh != nil {
			closeCh <- struct{}{}
		}
	}()
}


func Run() {
	fmt.Println("pid = ", os.Getpid())
	group, ctx := errgroup.WithContext(context.Background())

	group.Go(s1)
	group.Go(s2)
	sigRegister()

	<-ctx.Done()
	if err := ctx.Err(); err != nil {
		log.Println("select received: ", err)
		if closeCh != nil {
			closeCh <- struct{}{}
		}
	}
	close(closeCh)
	time.Sleep(time.Second * 2)
	log.Println("All server has shutdown")
}

/*
kill -s SIGTERM 73139

pid =  73139
2020/12/08 15:48:41 server1: listening on port :9080....
2020/12/08 15:48:41 server2: listening on port :9090....

2020/12/08 15:48:51 Got signal : terminated

2020/12/08 15:48:51 server2 received close signal
2020/12/08 15:48:51 select received:  context canceled
2020/12/08 15:48:51 server1 received close signal
2020/12/08 15:48:51 server1 has shutdown
2020/12/08 15:48:51 server2 has shutdown
2020/12/08 15:48:53 All server has shutdown
*/




