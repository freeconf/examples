package fcslack

import (
	"os"
	"testing"

	"github.com/freeconf/yang/fc"
)

func TestConnection(t *testing.T) {
	token := os.Getenv("SLACK_TEST_TOKEN")
	if token == "" {
		t.Skip("optional slack client test skipped")
	}
	client := newSlack()
	client.apply(slackOptions{
		ApiToken: token,
	})
	err := client.Send(msg{
		Channel: "testing-integration",
		Text:    "TestConnection unit test",
	})
	fc.AssertEqual(t, nil, err)
}
