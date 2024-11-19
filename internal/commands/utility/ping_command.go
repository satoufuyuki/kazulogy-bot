package utility

import (
	"context"
	"fmt"
	"github.com/satoufuyuki/kazulogy-bot/internal/client"
	"github.com/satoufuyuki/kazulogy-bot/internal/commands"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
	"time"
)

type pingCommand struct {
	*client.Client
}

func (s pingCommand) Meta() commands.CommandMeta {
	return commands.CommandMeta{
		Name:        "ping",
		Aliases:     []string{"pong"},
		Description: "Check the bot's latency",
	}
}

func (s pingCommand) Execute(event *events.Message, args []string) error {
	start := time.Now()
	resp, err := s.SendMessage(context.Background(), event.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String("🏓 Pinging..."),
			ContextInfo: &waE2E.ContextInfo{
				QuotedMessage: event.Message,
			},
		},
	})

	if err != nil {
		return err
	}

	latency := time.Since(start)
	resp, err = s.SendMessage(context.Background(), event.Info.Chat, s.BuildEdit(event.Info.Chat, resp.ID, &waE2E.Message{
		Conversation: proto.String(fmt.Sprintf("🏓 Took me *%dms* to respond!", latency.Milliseconds())),
	}))

	if err != nil {
		return err
	}

	return nil
}

func NewPingCommand(client *client.Client) commands.Command {
	return &pingCommand{client}
}
