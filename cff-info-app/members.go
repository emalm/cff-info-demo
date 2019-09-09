package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"code.cloudfoundry.org/lager"
)

const (
	maxRequests = 10
)

var ErrMemberServiceFailure error = errors.New("Member service failed to respond within retry limit")

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

type Metadata struct {
	Duration    time.Duration
	URL         string
	ResolvedIP  string
	ServerIP    string
	NumRequests int
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
		Duration:   1 * time.Millisecond,
		URL:        "<local>",
		ResolvedIP: "127.0.0.1",
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
	start := time.Now()
	resp, metadata, err := RequestMember(logger, f.url)
	end := time.Now()

	metadata.Duration = end.Sub(start).Round(time.Millisecond)

	if err != nil {
		logger.Error("request-failed", err, lager.Data{"metadata": metadata})
		return FetchResult{Metadata: metadata}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("readall-failed", err)
		return FetchResult{}, err
	}

	logger.Info("response", lager.Data{"body": string(body), "headers": resp.Header, "metadata": metadata})

	var response Response

	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Error("unmarshal-failed", err)
		return FetchResult{}, err
	}

	metadata.ServerIP = response.IP

	return FetchResult{Member: response.Member, Metadata: metadata}, nil
}

func RequestMember(logger lager.Logger, url string) (*http.Response, Metadata, error) {
	metadata := Metadata{URL: url}

	var response *http.Response
	var err error

	for numRequests := 1; numRequests <= maxRequests; numRequests++ {
		metadata.NumRequests = numRequests

		// cooling off time after failures
		time.Sleep(100 * time.Duration(numRequests-1) * time.Millisecond)

		// via https://stackoverflow.com/questions/49384786/how-to-capture-ip-address-of-server-providing-http-response
		client := &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					conn, err := net.Dial(network, addr)
					if err != nil {
						return conn, err
					}

					host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
					if err != nil {
						return conn, err
					}

					metadata.ResolvedIP = host

					return conn, err
				},
			},
		}

		response, err = client.Get(url)
		if err != nil {
			logger.Error("failed-request", err, lager.Data{"num": numRequests})
			continue
		}

		if response.StatusCode != http.StatusOK {
			logger.Error("failed-response-code", nil, lager.Data{"num": numRequests, "code": response.StatusCode})
			err = ErrMemberServiceFailure
			continue
		}

		logger.Info("response", lager.Data{"num": numRequests, "code": response.StatusCode})
		break
	}

	logger.Info("finished", lager.Data{"metadata": metadata, "error": err})
	return response, metadata, err
}
