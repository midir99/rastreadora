package ws

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/midir99/rastreadora/doc"
	"golang.org/x/net/html"
)

func MakeClient(skipVerify bool) *http.Client {
	if skipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		return &http.Client{Transport: tr}
	} else {
		return http.DefaultClient
	}
}

func RetrieveDocument(url string, skipVerify bool) (*doc.Doc, error) {
	client := MakeClient(skipVerify)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d status code", resp.StatusCode)
	}
	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return &doc.Doc{Node: node}, nil
}
