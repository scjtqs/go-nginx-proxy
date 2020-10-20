package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type realServer struct {
	Addr string
}

func (rs *realServer) HelloHandler(w http.ResponseWriter,r *http.Request){
	reqUrl := fmt.Sprintf("http://%s%s\n",rs.Addr,r.RequestURI)
	w.Write([]byte(reqUrl))
}

func (rs *realServer) Run(){
	fmt.Println("Http server tart to serve at :",rs.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/",rs.HelloHandler)
	server := &http.Server{
		Addr: rs.Addr,
		Handler: mux,
		WriteTimeout: time.Second * 3,
	}
	go func(){
		if err := server.ListenAndServe();err != nil{
			log.Fatal("Start http server failed,err:",err)
		}
	}()
}

func main() {
	rs1 := &realServer{Addr:"127.0.0.1:8081"}
	go rs1.Run()

	doneCh := make(chan os.Signal)
	signal.Notify(doneCh,syscall.SIGINT,syscall.SIGTERM)
	<- doneCh
}