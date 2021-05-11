package main

import (
  "bytes"
  "encoding/json"
  "errors"
  "net/http"
  "time"
)

const DefaultSlackTimeout = 5 * time.Second

type SlackClient struct {
  WebHookUrl string
  Bottoken   string
  UserName   string
  Channel    string
  TimeOut    time.Duration
}

type SlackPayload struct {
  Username    string       `json:"username,omitempty"`
  IconEmoji   string       `json:"icon_emoji,omitempty"`
  Channel     string       `json:"channel,omitempty"`
  Text        string       `json:"text,omitempty"`
  Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
  Color         string    `json:"color,omitempty"`
  AuthorName    string    `json:"author_name,omitempty"`
  Title         string    `json:"title,omitempty"`
  TitleLink     string    `json:"title_link,omitempty"`
  Pretext       string    `json:"pretext,omitempty"`
  Text          string    `json:"text,omitempty"`
  MarkdownIn  []string    `json:"mrkdwn_in,omitempty"`
  Ts          json.Number `json:"ts,omitempty"`
}

func (sc SlackClient) sendNotification (sp SlackPayload) error {
  slackRequest := SlackPayload {
    Text:        sp.Text,
    Username:    sc.UserName,
    IconEmoji:   sp.IconEmoji,
    Channel:     sc.Channel,
    Attachments: sp.Attachments,
  }
  return sc.sendHttpRequest(slackRequest)
}

func (sc SlackClient) sendHttpRequest(slackRequest SlackPayload) error {
  var slackRequestUrl string


  // For the case using slack bot instead of Webhook
  slackRequestUrl = sc.WebHookUrl
  if sc.Bottoken != "" && sc.WebHookUrl == "" { // If defined webhookurl so we don't need to send as bot, only pick 1
    slackRequestUrl = "https://slack.com/api/chat.postMessage"
  }

  slackBody, _ := json.Marshal(slackRequest)

  req, err := http.NewRequest(http.MethodPost, slackRequestUrl, bytes.NewBuffer(slackBody))
  if err != nil {
    return err
  }
  req.Header.Add("Content-Type", "application/json")

  // For the case using slack bot instead of Webhook
  if sc.Bottoken != "" && sc.WebHookUrl == "" { // If defined webhookurl so we don't need to send as bot, only pick 1
    req.Header.Add("Authorization", "Bearer " + sc.Bottoken)
  }

  if sc.TimeOut == 0 {
    sc.TimeOut = DefaultSlackTimeout
  }
  client := &http.Client{Timeout: sc.TimeOut}
  resp, err := client.Do(req)
  if err != nil {
    return err
  }

  // defer resp.Body.Close()
  //  dump, err := httputil.DumpResponse(resp, true)
	// if err != nil {
	//   log.Fatal(err)
	// }
  // fmt.Println(string(dump))


  if sc.Bottoken != "" && sc.WebHookUrl == "" { // If defined webhookurl so we don't need to send as bot, only pick 1
    type respContent struct {
      Ok bool                  `json:"ok"`
      X map[string]interface{} `json:"-"`
    }

    var r respContent

    defer resp.Body.Close()
    err = json.NewDecoder(resp.Body).Decode(&r)
    if err != nil {
      return err
    }
    if r.Ok != true {
      return errors.New("Non-ok response returned from Slack")
    }
  } else {
    buf := new(bytes.Buffer)
    _, err = buf.ReadFrom(resp.Body)
    if err != nil {
      return err
    }
    if buf.String() != "ok" {
      return errors.New("Non-ok response returned from Slack")
    }
  }
  return nil
}
