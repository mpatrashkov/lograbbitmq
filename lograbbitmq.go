// Package example is a CoreDNS plugin that prints "example" to stdout on every packet received.
//
// It serves as an example CoreDNS plugin with numerous code comments.
package lograbbitmq

import (
	"context"
	"encoding/json"
	"fmt"

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
}

type QueryResponse struct {
	Ip string
}

// ServeDNS implements the plugin.Handler interface. This method gets called when example is used
// in a Server.
func (e LogRabbitMQ) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	// This function could be simpler. I.e. just fmt.Println("example") here, but we want to show
	// a slightly more complex example as to make this more interesting.
	// Here we wrap the dns.ResponseWriter in a new ResponseWriter and call the next plugin, when the
	// answer comes back, it will print "example".

	// Debug log that we've have seen the query. This will only be shown when the debug plugin is loaded.
	q := r.Question[0]

	log.Debug(q)
	log.Debug(state.IP())
	log.Debug("test1")

	resp, err := http.Get(fmt.Sprintf("http://host.docker.internal:8123?domain=%s&type=%d&class=%d&ip=%s", q.Name, q.Qtype, q.Qclass, state.IP()))
	if err != nil {
		log.Debug(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug(err)
	}

	log.Debug("%d", len(body))
	log.Debug("%v", body)

	sb := string(body)

	log.Debug(body)

	var queryResponse QueryResponse
	json.Unmarshal([]byte(sb), &queryResponse)

	log.Debug(sb)

	answer := new(dns.Msg)
	answer.SetReply(r)
	answer.Authoritative = true

	rr, _ := dns.NewRR(fmt.Sprintf("%s 3600 IN A %s", q.Name, queryResponse.Ip))
	answer.Answer = []dns.RR{rr}

	// // Wrap.
	pw := NewResponsePrinter(w)

	pw.WriteMsg(answer)

	return 0, nil

	// // Export metric with the server label set to the current server handling the request.
	// // requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

	// // Call next plugin (if any).
	// return plugin.NextOrFailure(e.Name(), e.Next, ctx, pw, r)
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
	log.Info("example1234")
	return r.ResponseWriter.WriteMsg(res)
}
