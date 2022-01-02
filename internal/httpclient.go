package internal

import (
	"time"

	"github.com/valyala/fasthttp"
)

func Get(url string, headers map[string]string, timeout ...time.Duration) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.Header.SetMethod("GET")
	if err := doRequest(req, resp, url, headers, timeout...); err != nil {
		return nil, err
	}

	out := fasthttp.AcquireResponse()
	resp.CopyTo(out)
	return out, nil
}

func Post(url string, headers map[string]string, payload []byte, timeout ...time.Duration) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.Header.SetMethod("POST")
	req.SetBody(payload)
	if err := doRequest(req, resp, url, headers, timeout...); err != nil {
		return nil, err
	}

	out := fasthttp.AcquireResponse()
	resp.CopyTo(out)
	return out, nil
}

func Put(url string, headers map[string]string, payload []byte, timeout ...time.Duration) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.Header.SetMethod("PUT")
	req.SetBody(payload)
	if err := doRequest(req, resp, url, headers, timeout...); err != nil {
		return nil, err
	}

	out := fasthttp.AcquireResponse()
	resp.CopyTo(out)
	return out, nil
}

func Delete(url string, headers map[string]string, timeout ...time.Duration) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.Header.SetMethod("DELETE")
	if err := doRequest(req, resp, url, headers, timeout...); err != nil {
		return nil, err
	}

	out := fasthttp.AcquireResponse()
	resp.CopyTo(out)
	return out, nil
}

func doRequest(req *fasthttp.Request, resp *fasthttp.Response, url string, headers map[string]string, timeout ...time.Duration) error {
	req.SetRequestURI(url)
	req.Header.Add("User-Agent", "bitkub-go/1.0")
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	timeOut := 10 * time.Second
	if len(timeout) > 0 {
		timeOut = timeout[0]
	}
	if err := fasthttp.DoTimeout(req, resp, timeOut); err != nil {
		return err
	}
	return nil
}
