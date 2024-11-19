package commands

import (
	"context"
	"fmt"
	"github.com/satoufuyuki/kazulogy-bot/internal/client"
	"github.com/satoufuyuki/kazulogy-bot/pkg/framework/utilities"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
	"strings"
)

type CommandMeta struct {
	Name        string
	Aliases     []string
	Description string
}
type Command interface {
	Meta() CommandMeta
	Execute(message *events.Message, args []string) error
}

func execCommand(client *client.Client, command Command, ev *events.Message, args []string) {
	go func() {
		authorName := ev.Info.Sender.User
		cached, err := client.Store.Contacts.GetContact(ev.Info.Sender)
		if err == nil {
			if cached.Found {
				if cached.FullName != "" {
					authorName = cached.FullName
				} else {
					authorName = cached.PushName
				}
			}
		}

		client.Log.Infof("[%s] executing command: %s", authorName, command.Meta().Name)
		err = command.Execute(ev, args)
		if err != nil {
			client.Log.Errorf("[%s] command execution throw an error: %v", command.Meta().Name, err.Error())
			if _, err := client.SendMessage(context.Background(), ev.Info.Chat, &waE2E.Message{
				ExtendedTextMessage: &waE2E.ExtendedTextMessage{
					Text: proto.String(fmt.Sprintf("An *error occurred* while executing the command: %s", err.Error())),
					ContextInfo: &waE2E.ContextInfo{
						QuotedMessage: ev.Message,
					},
				},
			}); err != nil {
				panic(err)
			}
		}
	}()
}

func CommandHandler(client *client.Client, commands ...Command) whatsmeow.EventHandler {
	prefix := "/"
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			content := utilities.FindMessageContentFromMessage(v.Message)

			if len(content) == 0 || !strings.HasPrefix(content, prefix) {
				return
			}

			arguments := strings.Split(content[len(prefix):], " ")
			commandName := arguments[0]
			for _, command := range commands {
				if command.Meta().Name == commandName {
					execCommand(client, command, v, arguments[1:])
					return
				} else {
					for _, alias := range command.Meta().Aliases {
						if alias == commandName {
							execCommand(client, command, v, arguments[1:])
							return
						}
					}
				}
			}
		}
	}
}
