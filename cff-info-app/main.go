package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"code.cloudfoundry.org/lager"

	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/http_server"
	"github.com/tedsuo/ifrit/sigmon"
	"github.com/tedsuo/rata"
)

const (
	noTLSPort = 8080

	PhotosRoute = "Photos"
	PhotosDir   = "photos"

	RandomInfoRoute = "RandomInfo"
)

var Routes = rata.Routes{
	{Name: PhotosRoute, Method: "GET", Path: "/photos/"},
	{Name: RandomInfoRoute, Method: "GET", Path: "/random"},
}

func main() {
	logger := lager.NewLogger("cff-info")
	logger.RegisterSink(lager.NewPrettySink(os.Stdout, lager.INFO))
	photosPrefix, err := Routes.CreatePathForRoute(PhotosRoute, nil)

	photosServer := http.FileServer(http.Dir(PhotosDir))
	stripped := http.StripPrefix(photosPrefix, photosServer)

	var fetcher MemberFetcher

	url := os.Getenv("MEMBER_URL")
	if url == "" {
		fetcher = LocalMemberFetcher{}
	} else {
		fetcher = NewRemoteMemberFetcher(url)
	}

	infoHandler := NewRandomInfoHandler(logger, fetcher)

	handler, err := rata.NewRouter(Routes, rata.Handlers{
		PhotosRoute:     stripped,
		RandomInfoRoute: infoHandler,
	})

	if err != nil {
		logger.Fatal("failed-to-construct-router", err)
	}

	plainHTTPPort := os.Getenv("PORT")
	if plainHTTPPort == "" {
		plainHTTPPort = strconv.Itoa(noTLSPort)
	}

	plainHTTPServer := http_server.New("0.0.0.0:"+plainHTTPPort, handler)

	members := grouper.Members{
		{"plain", plainHTTPServer},
	}

	group := grouper.NewOrdered(os.Interrupt, members)

	monitor := ifrit.Invoke(sigmon.New(group))
	logger.Info("ready")

	err = <-monitor.Wait()
	if err != nil {
		logger.Error("exited-with-failure", err)
		os.Exit(1)
	}

	logger.Info("exited")
}

type randomInfoHandler struct {
	logger  lager.Logger
	fetcher MemberFetcher
}

func NewRandomInfoHandler(logger lager.Logger, fetcher MemberFetcher) http.Handler {
	return &randomInfoHandler{
		logger:  logger.Session("random"),
		fetcher: fetcher,
	}
}

func (h randomInfoHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	result, err := h.fetcher.Fetch(h.logger)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("failed-to-fetch-member", err)
		return
	}

	payload, err := json.Marshal(result)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("failed-to-marshal-data", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
