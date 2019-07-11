package dingtalk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	commoncfg "github.com/prometheus/common/config"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/notify"
	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/alertmanager/types"
)


// ----------------------------------------------------------------------
// add Dingtalk notify plugin
// ----------------------------------------------------------------------
type Dingtalk struct {
	conf   *config.DingtalkConfig
	tmpl   *template.Template
	logger log.Logger
}

type TextContent struct {
	Content string `json:"content"`
}

type AtContent struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type DingtalkMessage struct {
	Msgtype string      `json:"msgtype"`
	Text    TextContent `json:"text"`
	At      *AtContent  `json:"at"`
}

func NewDingtalk(conf *config.DingtalkConfig, tmpl *template.Template, l log.Logger) *Dingtalk {
	return &Dingtalk{conf: conf, tmpl: tmpl, logger: l}
}

// Notify implements the Notifier interface.
func (d *Dingtalk) Notify(ctx context.Context, alerts ...*types.Alert) (bool, error) {
	var tmplErr error
	data := d.tmpl.Data(receiverName(ctx, d.logger), groupLabels(ctx, d.logger), alerts...)
	tmpl := tmplText(d.tmpl, data, &tmplErr)

	token := d.conf.Token
	if token == "" {
		token, _ = d.conf.GroupToken[d.conf.Group]
	}
	if token == "" {
		level.Error(d.logger).Log("msg", "invalid group")
		return false, fmt.Errorf("invalid group=%s", d.conf.Group)
	}

	msg := DingtalkMessage{
		Msgtype: "text",
		Text:    TextContent{Content: tmpl(d.conf.Content)},
		At:      &AtContent{IsAtAll: true},
	}
	if tmplErr != nil {
		return false, fmt.Errorf("failed to template: %v", tmplErr)
	}

	var msgBuf bytes.Buffer
	if err := json.NewEncoder(&msgBuf).Encode(msg); err != nil {
		return false, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?access_token=%s",
		config.DefaultGlobalConfig().DingtalkAPIURL.String(), token), &msgBuf)
	if err != nil {
		return true, err
	}
	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("User-Agent", userAgentHeader)

	c, err := commoncfg.NewClientFromConfig(*d.conf.HTTPConfig, "dingtalk")
	if err != nil {
		return false, err
	}

	resp, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return true, err
	}
	resp.Body.Close()

	return d.retry(resp.StatusCode)
}

func (d *Dingtalk) retry(statusCode int) (bool, error) {
	// Webhooks are assumed to respond with 2xx response codes on a successful
	// request and 5xx response codes are assumed to be recoverable.
	if statusCode/100 != 2 {
		return (statusCode/100 == 5), fmt.Errorf("unexpected status code %v from %s", statusCode, d.conf.APIURL)
	}

	return false, nil
}

func truncate(s string, n int) (string, bool) {
	r := []rune(s)
	if len(r) <= n {
		return s, false
	}
	if n <= 3 {
		return string(r[:n]), true
	}
	return string(r[:n-3]) + "...", true
}