// Package example is a CoreDNS plugin that prints "example" to stdout on every packet received.
//
// It serves as an example CoreDNS plugin with numerous code comments.
package lograbbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"

	"io/ioutil"
	"net/http"
)

// Define log to be a logger with the plugin name in it. This way we can just use log.Info and
// friends to log.
var log = clog.NewWithPlugin("lograbbitmq")

// Example is an example plugin to show how to write a plugin.
type LogRabbitMQ struct {
	Next plugin.Handler
}

type QueryResponse struct {
	Ip string
}

// ServeDNS implements the plugin.Handler interface. This method gets called when example is used
// in a Server.
func (e LogRabbitMQ) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	q := r.Question[0]

	pw := NewResponsePrinter(w)

	resp, err := http.Get(fmt.Sprintf("http://host.docker.internal:8123?domain=%s&type=%d&class=%d&ip=%s", q.Name, q.Qtype, q.Qclass, state.IP()))
	if err != nil {
		log.Debug(err)
		plugin.NextOrFailure(e.Name(), e.Next, ctx, pw, r)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug(err)
		plugin.NextOrFailure(e.Name(), e.Next, ctx, pw, r)
	}

	sb := string(body)

	if sb == "null" {
		plugin.NextOrFailure(e.Name(), e.Next, ctx, pw, r)
		return 0, nil
	}

	var queryResponse QueryResponse
	json.Unmarshal([]byte(sb), &queryResponse)

	answer := new(dns.Msg)
	answer.SetReply(r)
	answer.Authoritative = true

	rr, _ := dns.NewRR(fmt.Sprintf("%s 3600 IN A %s", q.Name, queryResponse.Ip))
	answer.Answer = []dns.RR{rr}

	pw.WriteMsg(answer)

	return 0, nil
}

// Name implements the Handler interface.
func (e LogRabbitMQ) Name() string { return "lograbbitmq" }

// ResponsePrinter wrap a dns.ResponseWriter and will write example to standard output when WriteMsg is called.
type ResponsePrinter struct {
	dns.ResponseWriter
}

// NewResponsePrinter returns ResponseWriter.
func NewResponsePrinter(w dns.ResponseWriter) *ResponsePrinter {
	return &ResponsePrinter{ResponseWriter: w}
}

// WriteMsg calls the underlying ResponseWriter's WriteMsg method and prints "example" to standard output.
func (r *ResponsePrinter) WriteMsg(res *dns.Msg) error {
	return r.ResponseWriter.WriteMsg(res)
}
