package main

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"time"

	"net/http"
	"os/signal"

	"golang.org/x/sync/errgroup"
)


func main(){
	//group.GO只要一个返回error，这里的ctx就会被cancel
	//group的error只会记录一次（记录第一次的error值），是通过sync.Once方式赋值error，具体细节查看源码
	group,ctx := errgroup.WithContext(context.Background())

	stopChan := make(chan struct{})
	group.Go(func() error {
		return serverApp(ctx,stopChan)
	})
	group.Go(func() error {
		return serverDebug(ctx,stopChan)

	})
	go watchSignal(stopChan)
	if err := group.Wait(); err != nil {
		fmt.Println("server error and exit all")
	}else {
		fmt.Println("exit all server successfully")
	}

}

func serverApp(ctx context.Context,stopChan <-chan struct{})(err error){
	server:=http.Server{
		Addr:              "0.0.0.0:8080",
		Handler:           http.DefaultServeMux,
	}
	go func(){
		select {
			case <- ctx.Done():
				fmt.Println("ctx cancel")
			case <- stopChan :
				fmt.Println("get a stop signal")
			case <-time.After(10*time.Second)://模拟一个退出，其他也退出
				fmt.Println("timeout")
		}

		fmt.Println("exit serverApp")
		if err = server.Shutdown(context.Background());err !=nil {
			fmt.Printf("serverApp Shutdown error %s \n  ",err.Error())
			return
		}
	}()
	fmt.Printf("serverApp start at %s\n",server.Addr)
	if err = server.ListenAndServe();err !=nil {
		fmt.Printf("serverApp exit %s \n",err.Error())
		return
	}
	return nil
}

func serverDebug(ctx context.Context,stopChan <-chan struct{})(err error){
	server:=http.Server{
		Addr:              "0.0.0.0:8081",
		Handler:           http.DefaultServeMux,
	}
	go func(){
		select {
		case <- ctx.Done():
			fmt.Println("ctx cancel")
		case <- stopChan :
			fmt.Println("get a stop signal")
		}
		fmt.Println("exit serverDebug")
		if err = server.Shutdown(context.Background());err !=nil {
			fmt.Printf("serverDebug Shutdown error %s \n  ",err.Error())
			return
		}
	}()
	fmt.Printf("serverDebug start at %s\n",server.Addr)
	if err = server.ListenAndServe();err !=nil {
		fmt.Printf("serverDebug exit %s \n",err.Error())
		return
	}
	return nil
}

func watchSignal(stopChan chan struct{}){
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		fmt.Printf("get a signal %s \n", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			close(stopChan)
			fmt.Println("do server exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}