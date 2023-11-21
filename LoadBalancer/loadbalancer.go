package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

var (
	baseURL = "http://localhost:808"
)

type Loadbalancer struct {
	RevrseProxy httputil.ReverseProxy
}

type Endpoints struct {
	List []*url.URL
}

func (e *Endpoints) shuffle() {
	temp := e.List[0]
	e.List = e.List[1:]
	e.List = append(e.List, temp)
}

func MakeLoadBalancer(amount int) {
	var lb Loadbalancer
	var ep Endpoints

	router := http.NewServeMux()
	server := http.Server{
		Addr:    ":8090",
		Handler: router,
	}

	for i := 0; i < amount; i++ {
		ep.List = append(ep.List, createEndpoint(baseURL, i))
	}
	
	router.HandleFunc("/loadbalancer", makeRequrest(&lb, &ep))
	log.Fatal(server.ListenAndServe())
}

func makeRequrest(lb *Loadbalancer, ep *Endpoints) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		for !testServer(ep.List[0].String()) {
			ep.shuffle()
		}

		lb.RevrseProxy = *httputil.NewSingleHostReverseProxy(ep.List[0])
		ep.shuffle()
		lb.RevrseProxy.ServeHTTP(w, r)
	}
}

func createEndpoint(endpoint string, port int) *url.URL {
	link := endpoint + strconv.Itoa(port)
	url, _ := url.Parse(link)
	return url
}

func testServer(endpoint string) bool {
	res, err := http.Get(endpoint)
	if err != nil {
		return false
	}

	if res.StatusCode != http.StatusOK {
		return false
	}
	return true
}
