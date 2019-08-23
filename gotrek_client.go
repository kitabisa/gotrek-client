package gotrekclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/httpclient"
)

type TrekClient interface {
	Publish(auditID string, trail map[string]interface{}, timestamp int64, tag string) error
}

type trekClient struct {
	url          string
	clientSecret string
	httpClient   *httpclient.Client
}

type TrekHttpClient struct {
	Timeout               time.Duration
	BackoffInterval       time.Duration
	MaximumJitterInterval time.Duration
	RetryCount            int
}

func NewTrekClient(url, clientSecret string, t *TrekHttpClient) TrekClient {
	return &trekClient{
		url:          url,
		clientSecret: clientSecret,
		httpClient:   newTrekHttpClient(t),
	}
}

func newTrekHttpClient(t *TrekHttpClient) *httpclient.Client {
	var c *TrekHttpClient
	if t == nil {
		c = new(TrekHttpClient)
		c.BackoffInterval = 2 * time.Millisecond
		c.MaximumJitterInterval = 5 * time.Millisecond
		c.Timeout = 2000 * time.Millisecond
	} else {
		c = t
	}

	backoff := heimdall.NewConstantBackoff(c.BackoffInterval, c.MaximumJitterInterval)
	retrier := heimdall.NewRetrier(backoff)

	return httpclient.NewClient(
		httpclient.WithHTTPTimeout(c.Timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(c.RetryCount),
	)
}

func (c *trekClient) Publish(auditID string, trail map[string]interface{}, timestamp int64, tag string) error {

	payload := map[string]interface{}{
		"id":        auditID,
		"timestamp": timestamp,
		"trail":     trail,
		"tag":       tag,
	}

	payloadMarshal, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("%s%s", c.url, gotrekPublishUrl)
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(payloadMarshal))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.clientSecret))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyByte, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			return fmt.Errorf("Error code %d : %s", resp.StatusCode, string(bodyByte))
		}
	}

	return nil
}
