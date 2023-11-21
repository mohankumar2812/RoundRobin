package servers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type ServerList struct {
	ports []int

}

func (s *ServerList) Populate(amount int) {

	if amount >= 10 {
		 log.Fatal("port amount exceed");
	}

	for x:= 0; x < amount; x++ {
		s.ports = append(s.ports, x)
	}

}

func (s *ServerList) Pop() int {
	port := s.ports[0]
	s.ports = s.ports[1:]
	return port
}

func RunServers(amount int) {
	var myServerList ServerList

	myServerList.Populate(amount)

	//waitGroup
	var wg sync.WaitGroup
	wg.Add(amount)
	defer wg.Wait()

	for x:= 0; x < amount; x++ {
		go makeServers(&myServerList, wg)
	}
}

func makeServers(s *ServerList, wg sync.WaitGroup) {
	defer wg.Done()
	r := http.NewServeMux()

	//port
	port := s.Pop()
	server := http.Server{
		Addr: fmt.Sprintf(":808%d",port),
		Handler: r,
	}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w,"server %d",port)
	})

	r.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("400 - Server Shut Down!"))
		server.Shutdown(context.Background())
	})

	server.ListenAndServe()

}

