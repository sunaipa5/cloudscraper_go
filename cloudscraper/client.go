package cloudscraper

import (
	"maps"

	"github.com/RomainMichau/CycleTLS/cycletls"
)

type CloudScrapper struct {
	client        cycletls.CycleTLS
	respChan      chan cycletls.Response
	defaultHeader map[string]string
	ja3           string
	userAgent     string
}

func (cs CloudScrapper) ExportSettings() map[string]string {
	settings := make(map[string]string)
	settings["ja3"] = cs.ja3
	settings["userAgent"] = cs.userAgent
	maps.Copy(settings, cs.defaultHeader)
	return settings
}

func (cs CloudScrapper) Do(url string, options cycletls.Options, method string) (cycletls.Response, error) {
	if options.Headers == nil {
		options.Headers = make(map[string]string)
	}
	maps.Copy(options.Headers, cs.defaultHeader)
	if options.UserAgent == "" {
		options.UserAgent = cs.userAgent
	}
	if options.Ja3 == "" {
		options.Ja3 = cs.ja3
	}
	return cs.client.Do(url, options, method)
}

// Queue Push a request in the Queue
// Read the result with RespChan()
// Queue pushes a request into the queue asynchronously and streams the response to respChan
func (cs CloudScrapper) Queue(url string, options cycletls.Options, method string) {
	if options.Headers == nil {
		options.Headers = make(map[string]string)
	}

	for k, v := range cs.defaultHeader {
		if _, exists := options.Headers[k]; !exists {
			options.Headers[k] = v
		}
	}
	if options.UserAgent == "" {
		options.UserAgent = cs.userAgent
	}
	if options.Ja3 == "" {
		options.Ja3 = cs.ja3
	}

	// Execute the request asynchronously using Do and forward to cs.respChan
	go func() {
		resp, err := cs.client.Do(url, options, method)
		if err != nil {
			return
		}
		cs.respChan <- resp
	}()
}

func (cs CloudScrapper) Get(url string, headers map[string]string, body string) (cycletls.Response, error) {
	options := cycletls.Options{
		Body:    body,
		Headers: headers,
	}
	return cs.Do(url, options, "GET")
}

func (cs CloudScrapper) Post(url string, headers map[string]string, body string) (cycletls.Response, error) {
	options := cycletls.Options{
		Body:    body,
		Headers: headers,
	}
	return cs.Do(url, options, "POST")
}

func (cs CloudScrapper) RespChan() chan cycletls.Response {
	return cs.respChan
}

func Init(mobile bool, workers bool) (*CloudScrapper, error) {
	browserConf, err := getUserAgents(mobile)
	if err != nil {
		return nil, err
	}

	cycleTstClient := cycletls.Init(workers)

	p := CloudScrapper{
		client:        cycleTstClient,
		defaultHeader: browserConf.Headers,
		ja3:           browserConf.Ja3,
		userAgent:     browserConf.UserAgent,
		respChan:      make(chan cycletls.Response, 100), // Buffered channel for async Queue responses
	}

	return &p, nil
}
