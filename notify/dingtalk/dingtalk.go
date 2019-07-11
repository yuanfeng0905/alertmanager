// Copyright 2015 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dingtalk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	commoncfg "github.com/prometheus/common/config"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/notify"
	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/alertmanager/types"
)

// Notifier add Dingtalk notify plugin
type Notifier struct {
	conf   *config.DingtalkConfig
	tmpl   *template.Template
	logger log.Logger
	client *http.Client
}

type textContent struct {
	Content string `json:"content"`
}

type atContent struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type dingtalkMessage struct {
	Msgtype string      `json:"msgtype"`
	Text    textContent `json:"text"`
	At      *atContent  `json:"at"`
}

// New returns a new Dingtalk notification handler.
func New(c *config.DingtalkConfig, t *template.Template, l log.Logger) (*Notifier, error) {
	client, err := commoncfg.NewClientFromConfig(*c.HTTPConfig, "dingtalk")
	if err != nil {
		return nil, err
	}

	return &Notifier{
		conf:   c,
		tmpl:   t,
		logger: l,
		client: client,
	}, nil
}

// Notify implements the Notifier interface.
func (n *Notifier) Notify(ctx context.Context, as ...*types.Alert) (bool, error) {
	var err error

	data := notify.GetTemplateData(ctx, n.tmpl, as, n.logger)
	tmpl := notify.TmplText(n.tmpl, data, &err)
	token := n.conf.Token

	if token == "" {
		token, _ = n.conf.GroupToken[n.conf.Group]
	}
	if token == "" {
		level.Error(n.logger).Log("msg", "invalid group")
		return false, fmt.Errorf("invalid group=%s", n.conf.Group)
	}

	msg := dingtalkMessage{
		Msgtype: "text",
		Text:    textContent{Content: tmpl(n.conf.Content)},
		At:      &atContent{IsAtAll: true},
	}
	if err != nil {
		return false, fmt.Errorf("failed to template: %v", err)
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(msg); err != nil {
		return false, err
	}

	url := fmt.Sprintf("%s?access_token=%s", config.DefaultGlobalConfig().DingtalkAPIURL.String(), token)
	resp, err := notify.PostJSON(ctx, n.client, url, &buf)
	if err != nil {
		return true, notify.RedactURL(err)
	}
	defer notify.Drain(resp)

	return n.retry(resp.StatusCode)
}

func (n *Notifier) retry(statusCode int) (bool, error) {
	// Webhooks are assumed to respond with 2xx response codes on a successful
	// request and 5xx response codes are assumed to be recoverable.
	if statusCode/100 != 2 {
		return (statusCode/100 == 5), fmt.Errorf("unexpected status code %v from %s", statusCode, n.conf.APIURL)
	}

	return false, nil
}
