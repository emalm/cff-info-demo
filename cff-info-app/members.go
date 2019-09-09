package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"code.cloudfoundry.org/lager"
)

type Member struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
	Bio   string `json:"bio"`
}

type Metadata struct {
	Duration time.Duration
	URL      string
	Addr     string
}

type FetchResult struct {
	Member   Member   `json:"member"`
	Metadata Metadata `json:"metadata"`
}

type MemberFetcher interface {
	Fetch(logger lager.Logger) (FetchResult, error)
}

type LocalMemberFetcher struct{}

func (f LocalMemberFetcher) Fetch(logger lager.Logger) (FetchResult, error) {
	member := Member{
		ID:    "local",
		Name:  "L. Ocal",
		Title: "Factotum",
		Bio:   "Irrelevant",
	}

	metadata := Metadata{
		Duration: 1 * time.Millisecond,
		URL:      "<local>",
		Addr:     "127.0.0.1:8080",
	}

	return FetchResult{Member: member, Metadata: metadata}, nil
}

type RemoteMemberFetcher struct {
	url string
}

func NewRemoteMemberFetcher(url string) MemberFetcher {
	return &RemoteMemberFetcher{
		url: url,
	}
}

func (f RemoteMemberFetcher) Fetch(logger lager.Logger) (FetchResult, error) {
	// via https://stackoverflow.com/questions/49384786/how-to-capture-ip-address-of-server-providing-http-response
	request, err := http.NewRequest("GET", f.url, nil)
	if err != nil {
		return FetchResult{}, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				conn, err := net.Dial(network, addr)
				if err != nil {
					return conn, err
				}
				request.RemoteAddr = conn.RemoteAddr().String()
				return conn, err
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			request = req
			return nil
		},
	}

	start := time.Now()
	resp, err := client.Do(request)
	end := time.Now()

	if err != nil {
		return FetchResult{}, err
	}

	defer resp.Body.Close()

	metadata := Metadata{
		Duration: end.Sub(start),
		URL:      f.url,
		Addr:     resp.Request.RemoteAddr,
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("readall-failed", err)
		return FetchResult{}, err
	}

	logger.Info("response", lager.Data{"body": string(body), "duration": metadata.Duration, "headers": resp.Header, "request-addr": resp.Request.RemoteAddr})

	var member Member

	err = json.Unmarshal(body, &member)
	if err != nil {
		logger.Error("unmarshal-failed", err)
		return FetchResult{}, err
	}

	return FetchResult{Member: member, Metadata: metadata}, nil
}
