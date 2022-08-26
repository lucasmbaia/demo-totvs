package request

import (
	"net/http/cookiejar"
	"encoding/json"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"bytes"
	"net/url"
	"context"
	"errors"
	"fmt"
	"os"
)

const (
	GET	= "GET"
	POST	= "POST"
	PUT	= "PUT"
	PATCH	= "PATCH"
	DELETE	= "DELETE"

	ContextBasicAuth  int = 2
)

type Response struct {
	Header	http.Header
	Code	int
	Body	[]byte
}

type Options struct {
	Ctx	    context.Context
	Body	    interface{}
	Headers	    map[string]string
	Params	    url.Values
}

type BasicAuth struct {
	Username  string  `json:"userName,omitempty"`
	Password  string  `json:"password,omitempty"`
}

type Client struct {
	client	*http.Client
}

func NewClient() (c *Client, err error) {
	var jar *cookiejar.Jar

	if jar, err = cookiejar.New(nil); err != nil {
		return
	}

	return &Client{
		client:	&http.Client{
			Transport: &http.Transport{
				TLSClientConfig:      &tls.Config{InsecureSkipVerify: true},
				MaxIdleConns:	      100,
				MaxIdleConnsPerHost:  100,
			},
			Jar:	jar,
		},
	}, nil
}

func (c *Client) SetCookies(u *url.URL, cookies	[]*http.Cookie) {
	c.client.Jar.SetCookies(u, cookies)
}

func (c *Client) GetCookies(u *url.URL) []*http.Cookie {
	return c.client.Jar.Cookies(u)
}

type ReadData struct {
	Data	Data  `json:"data,omitempty"`
}

type Data struct {
	HasMore	bool	      `json:"HasMore"`
	From	int	      `json:"From"`
	Hits	[]interface{} `json:"Hits"`
}

type ErrorMessage struct {
	Code	int	`json:"code,omitempty"`
	Message	string	`json:"message,omitempty"`
}

func (c *Client) Post(index string, payload interface{}) (err error) {
	var (
		path	  string
		response  Response
	)

	path = fmt.Sprintf("%s/v1/indices/%s/documents", os.Getenv("DATA_COLLECTOR_URL"), index)

	if response, err = c.Request(POST, path, Options{
		Body:     payload,
		Headers:  map[string]string{
		        "Content-Type": "application/json",
		        "Accept":       "application/json",
		        "User-Agent":   "",
		},
	}); err != nil {
		return
	}

	if response.Code != 201 {
		var em ErrorMessage

		if err = json.Unmarshal(response.Body, &em); err != nil {
			err = errors.New(string(response.Body))
		} else {
			err = errors.New(em.Message)
		}
	}

	return
}

func (c *Client) Get(index string, i interface{}, optionals map[string]interface{}) (err error) {
	var (
		path	  string
		response  Response
		params	  = url.Values{}
		rd	  ReadData
		body	  []byte
	)

	path = fmt.Sprintf("%s/v1/indices/%s/documents?size=1000", os.Getenv("DATA_COLLECTOR_URL"), index)

	if response, err = c.Request(GET, path, Options{
	        Headers:  map[string]string{
	                "Content-Type": "application/json",
	                "Accept":       "application/json",
	                "User-Agent":   "",
	        },
	        Params:   params,
	}); err != nil {
	        return
	}

	if response.Code == 200 {
		if err = json.Unmarshal(response.Body, &rd); err != nil {
			return
		}

		if body, err = json.Marshal(rd.Data.Hits); err != nil {
			return
		}

		if err = json.Unmarshal(body, i); err != nil {
			return
		}
	} else {
		var em ErrorMessage

		if err = json.Unmarshal(response.Body, &em); err != nil {
			err = errors.New(string(response.Body))
		} else {
			err = errors.New(em.Message)
		}
	}

	return
}

func (c *Client) Request(method, path string, o Options) (r Response, err error) {
	var (
		req	*http.Request
		resp	*http.Response
		b	[]byte
		pb	io.Reader
		query	url.Values
		uq	*url.URL
	)

	if o.Body != nil {
		var body []byte

		if body, err = json.Marshal(o.Body); err != nil {
			return
		}

		pb = bytes.NewReader(body)
	}

	if uq, err = url.Parse(path); err != nil {
		return
	}

	query = uq.Query()
	for k, v := range o.Params {
		for _, iv := range v {
			query.Add(k, iv)
		}
	}

	uq.RawQuery = query.Encode()

	if o.Body != nil {
		if req, err = http.NewRequest(method, uq.String(), pb); err != nil {
			return
		}
	} else {
		if req, err = http.NewRequest(method, uq.String(), nil); err != nil {
			return
		}
	}

	//req.Close = true

	if o.Headers != nil {
		for k, v := range o.Headers {
			req.Header.Set(k, v)
		}
	}

	if o.Ctx != nil {
		if auth, ok := o.Ctx.Value(ContextBasicAuth).(BasicAuth); ok {
			req.SetBasicAuth(auth.Username, auth.Password)
		}
	}

	if resp, err = c.client.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	r = Response{Header: resp.Header, Code: resp.StatusCode, Body: b}
	return
}
