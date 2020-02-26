/**
* Copyright 2020 Comcast Cable Communications Management, LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
* SPDX-License-Identifier: Apache-2.0
 */

package proxy

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"time"

	translator "github.com/Comcast/DNS-over-HTTPs-translator"
	"github.com/miekg/dns"

	"go.uber.org/zap"
)

type proxy struct {
	// Instance of the zap logger.
	logger *zap.Logger

	// IP address and port of the resolver.
	resolver string

	// Address and port that the proxy will listen on.
	listen string
}

// New returns an instance of Translator.
func New(options ...Option) (result translator.Proxy, err error) {
	result = &proxy{
		logger: zap.NewNop(),
	}

	// Apply options.
	for _, opt := range options {
		opt(result.(*proxy))
	}

	return
}

func (p *proxy) Start(ctx context.Context) (err error) {
	http.HandleFunc("/", p.serveProxyRequest)
	p.logger.Info("Starting DOH translator HTTP proxy server")
	p.logger.Info("Resolver in use:", zap.String("resolver_string", p.resolver))
	if err := http.ListenAndServe(p.listen, nil); err != nil {
		panic(err)
	}
	return err
}

func (p *proxy) Stop() (err error) {
	return
}

func (p *proxy) serveProxyRequest(w http.ResponseWriter, r *http.Request) {

	var (
		buf   []byte
		err   error
		reply *dns.Msg
		rtt   time.Duration
	)

	switch r.Method {
	case http.MethodGet:
		dnsQuery := r.URL.Query().Get("dns")
		buf, err = base64.RawURLEncoding.DecodeString(dnsQuery)
		if len(buf) == 0 || err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

	case http.MethodPost:
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/dns-message" {
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}

		buf, err = ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	msg := new(dns.Msg)
	if err = msg.Unpack(buf); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	p.logger.Info("DNS request received", zap.String("dns_message", msg.String()))

	client := new(dns.Client)
	reply, rtt, err = client.Exchange(msg, p.resolver)
	if reply.Truncated {
		client.Net = "tcp"
		reply, rtt, err = client.Exchange(msg, p.resolver)
	}

	rttSeconds := rtt.Seconds()

	p.logger.Info("RTT for resolve in seconds", zap.Float64("rtt_sec", rttSeconds))

	if err != nil {
		p.logger.Error("Exchange failed", zap.Error(err))
		return
	}

	bytes, err := reply.Pack()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	w.Header().Set("Server", p.resolver)
	w.Header().Set("Content-Type", "application/dns-message")
	_, err = w.Write(bytes)
}
