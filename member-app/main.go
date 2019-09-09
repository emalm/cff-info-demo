package main

import (
	"encoding/json"
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"

	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/http_server"
	"github.com/tedsuo/ifrit/sigmon"
)

const (
	defaultPort     = "8080"
	badListenerPort = "9999"
)

func main() {
	logger := lager.NewLogger("member")
	logger.RegisterSink(lager.NewPrettySink(os.Stdout, lager.INFO))

	m := Member{
		ID:    os.Getenv("MEMBER_ID"),
		Name:  os.Getenv("MEMBER_NAME"),
		Title: os.Getenv("MEMBER_TITLE"),
		Bio:   os.Getenv("MEMBER_BIO"),
	}

	ip := os.Getenv("CF_INSTANCE_INTERNAL_IP")
	if ip == "" {
		ip = "127.0.0.1"
	}

	handler := NewMemberHandler(logger, m, ip)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	badListener := os.Getenv("BAD_LISTENER") == "true"
	if badListener {
		port = badListenerPort
	}

	server := http_server.New("0.0.0.0:"+port, handler)

	members := grouper.Members{
		{"plain", server},
	}

	group := grouper.NewOrdered(os.Interrupt, members)

	monitor := ifrit.Invoke(sigmon.New(group))
	logger.Info("ready")

	err := <-monitor.Wait()
	if err != nil {
		logger.Error("exited-with-failure", err)
		os.Exit(1)
	}

	logger.Info("exited")
}

type Member struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
	Bio   string `json:"bio"`
}

type Response struct {
	Member Member `json:"member"`
	IP     string `json:"ip"`
}

type memberHandler struct {
	logger   lager.Logger
	response Response
}

func NewMemberHandler(logger lager.Logger, m Member, ip string) http.Handler {
	return &memberHandler{
		logger: logger.Session("handler"),
		response: Response{
			Member: m,
			IP:     ip,
		},
	}
}

func (h memberHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	buf, err := json.Marshal(h.response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("failed-to-encode", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
