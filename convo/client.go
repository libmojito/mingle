package convo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	DefaultServerAddress = "http://localhost:8080"
	ContentTypeJson      = "application/json"
)

type chatReq struct {
	ThreadID string            `json:"thread_id"`
	Content  string            `json:"content"`
	Security map[string]string `json:"security"`
}

type chatRes struct {
	Response string `json:"response"`
}

type threadRes struct {
	ID string `json:"id"`
}

type ChatClient struct {
	ThreadClient ThreadClient
	ThreadID     string
	Security     map[string]string
}

type ThreadClient struct {
	serverAddress string
}

func NewThreadClient(options ...func(*ThreadClient)) *ThreadClient {
	c := &ThreadClient{
		serverAddress: DefaultServerAddress,
	}

	for _, o := range options {
		o(c)
	}

	return c
}

func WithServerAddress(addr string) func(*ThreadClient) {
	return func(c *ThreadClient) {
		c.serverAddress = addr
	}
}

func (sc *ThreadClient) NewThread() string {
	res, err := http.Post(fmt.Sprintf("%s/thread", sc.serverAddress), ContentTypeJson, nil)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode > 299 {
		log.Fatalf(
			"Response failed with status code: %d and \nbody: %s\n",
			res.StatusCode,
			body,
		)
	}

	var resJson threadRes
	if err = json.Unmarshal(body, &resJson); err != nil {
		log.Fatal(err)
	}
	if resJson.ID == "" {
		log.Fatalf("Invalid return json thread %s", body)
	}
	return resJson.ID
}

func NewChatClient(server string, threadID string, security map[string]string) *ChatClient {
	return &ChatClient{
		ThreadClient: *NewThreadClient(
			WithServerAddress(server),
		),
		ThreadID: threadID,
		Security: security,
	}
}

func (s *ChatClient) SendMessageContent(c string) Message {
	req := chatReq{
		ThreadID: s.ThreadID,
		Content:  c,
		Security: s.Security,
	}
	reqJson, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post(
		fmt.Sprintf("%s/chat", s.ThreadClient.serverAddress),
		ContentTypeJson,
		bytes.NewBuffer(reqJson),
	)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode > 299 {
		log.Fatalf(
			"Response failed with status code: %d and \nbody: %s\n",
			res.StatusCode,
			body,
		)
	}

	var resJson chatRes
	if err = json.Unmarshal(body, &resJson); err != nil {
		log.Fatal(err)
	}

	return Message{
		Role:    RoleAssistant,
		Content: resJson.Response,
	}
}
