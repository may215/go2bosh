package go2bosh

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

/* Custom timout dialer, and the main reason is for avoiding closing the response before getting its body */
func timeoutDialler() *http.Client {
	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: func(netw, addr string) (net.Conn, error) {
			deadline := time.Now().Add(time.Duration(config.RequestTimeOut) * time.Millisecond)
			c, err := net.DialTimeout(netw, addr, time.Second)
			if err != nil {
				return nil, err
			}
			c.SetDeadline(deadline)
			return c, nil
		}}
	httpclient := &http.Client{Transport: transport}
	return httpclient
}

/* Make request for the bosh server */
func makeRequest(data string, child string) ([]byte, *handlerError) {
	bosh_url := config.BoshServer

	client := timeoutDialler()
	req, _ := http.NewRequest("POST", bosh_url, bytes.NewBufferString(data))
	req.Header.Add("Accept", "text/xml")
	req.Header.Add("Content-Type", "text/xml")

	resp, err := client.Do(req)
	if err != nil {
		return nil, &handlerError{err, "Unable to make request", 100107}
	}

	defer resp.Body.Close()
	res_data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &handlerError{err, "Unable to read response data", 100108}
	}
	return res_data, nil
}
