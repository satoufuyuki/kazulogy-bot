package utility

import (
	"context"
	"github.com/gabriel-vasile/mimetype"
	"github.com/satoufuyuki/kazulogy-bot/internal/client"
	"github.com/satoufuyuki/kazulogy-bot/internal/commands"
	"github.com/satoufuyuki/kazulogy-bot/pkg/framework/utilities"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
	"strings"
)

type stealCommand struct {
	*client.Client
}

func (s stealCommand) Meta() commands.CommandMeta {
	return commands.CommandMeta{
		Name:        "steal",
		Aliases:     []string{"s"},
		Description: "Steal any quoted media from the chat and status",
	}
}

func (s stealCommand) Execute(event *events.Message, args []string) error {
	authorMessageContext := &waE2E.ContextInfo{QuotedMessage: event.Message}
	quotedMessage := event.Message.GetExtendedTextMessage().GetContextInfo().QuotedMessage
	if quotedMessage == nil {
		if _, err := s.SendMessage(context.Background(), event.Info.Chat, &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text:        proto.String("You have to quote a message to steal."),
				ContextInfo: authorMessageContext,
			},
		}); err != nil {
			return err
		}

		return nil
	}

	ephemeralMessage := quotedMessage.EphemeralMessage
	messageCandidates := []*waE2E.Message{
		quotedMessage,
		quotedMessage.GetViewOnceMessageV2().GetMessage(),
		quotedMessage.GetViewOnceMessage().GetMessage(),
		ephemeralMessage.GetMessage(),
		ephemeralMessage.GetMessage().GetViewOnceMessageV2().GetMessage(),
		ephemeralMessage.GetMessage().GetViewOnceMessage().GetMessage(),
	}

	var result []byte
	var existingThumbnail []byte
	for i, candidate := range messageCandidates {
		if candidate == nil {
			continue
		}

		r, err := s.DownloadAny(candidate)
		if err != nil && i == len(messageCandidates)-1 {
			break
		} else if err == nil {
			result = r
			existingThumbnail = utilities.FindThumbnailFromMessage(candidate)
			break
		}
	}

	if len(result) == 0 {
		if _, err := s.SendMessage(context.Background(), event.Info.Chat, &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text:        proto.String("There is no media to steal."),
				ContextInfo: authorMessageContext,
			},
		}); err != nil {
			return err
		}

		return nil
	}

	caption := utilities.FindMessageContentFromMessage(quotedMessage)
	mimeType := mimetype.Detect(result).String()
	payload := &waE2E.Message{}
	if strings.HasPrefix(mimeType, "image") {
		if existingThumbnail == nil {
			thumbnail, err := utilities.GenerateThumbnail(result, mimeType)
			if err != nil {
				return err
			}

			existingThumbnail = thumbnail
		}

		media, err := s.Upload(context.Background(), result, whatsmeow.MediaImage)
		if err != nil {
			return err
		}

		payload.ImageMessage = &waE2E.ImageMessage{
			ContextInfo:   authorMessageContext,
			Mimetype:      proto.String(mimeType),
			Caption:       proto.String(caption),
			URL:           &media.URL,
			DirectPath:    &media.DirectPath,
			MediaKey:      media.MediaKey,
			FileEncSHA256: media.FileEncSHA256,
			FileSHA256:    media.FileSHA256,
			FileLength:    &media.FileLength,
			JPEGThumbnail: existingThumbnail,
		}
	} else if strings.HasPrefix(mimeType, "video") {
		media, err := s.Upload(context.Background(), result, whatsmeow.MediaVideo)
		if err != nil {
			return err
		}

		payload.VideoMessage = &waE2E.VideoMessage{
			ContextInfo:   authorMessageContext,
			Mimetype:      proto.String(mimeType),
			Caption:       proto.String(caption),
			URL:           &media.URL,
			DirectPath:    &media.DirectPath,
			MediaKey:      media.MediaKey,
			FileEncSHA256: media.FileEncSHA256,
			FileSHA256:    media.FileSHA256,
			FileLength:    &media.FileLength,
			JPEGThumbnail: existingThumbnail,
		}
	} else if strings.HasPrefix(mimeType, "audio") {
		media, err := s.Upload(context.Background(), result, whatsmeow.MediaAudio)
		if err != nil {
			return err
		}

		payload.AudioMessage = &waE2E.AudioMessage{
			ContextInfo:   authorMessageContext,
			Mimetype:      proto.String(mimeType),
			URL:           &media.URL,
			DirectPath:    &media.DirectPath,
			MediaKey:      media.MediaKey,
			FileEncSHA256: media.FileEncSHA256,
			FileSHA256:    media.FileSHA256,
			FileLength:    &media.FileLength,
		}
	} else {
		media, err := s.Upload(context.Background(), result, whatsmeow.MediaDocument)
		if err != nil {
			return err
		}

		payload.DocumentMessage = &waE2E.DocumentMessage{
			ContextInfo:   authorMessageContext,
			Mimetype:      proto.String(mimeType),
			Caption:       proto.String(caption),
			URL:           &media.URL,
			DirectPath:    &media.DirectPath,
			MediaKey:      media.MediaKey,
			FileEncSHA256: media.FileEncSHA256,
			FileSHA256:    media.FileSHA256,
			FileLength:    &media.FileLength,
		}
	}

	_, err := s.SendMessage(context.Background(), event.Info.Chat, payload)
	if err != nil {
		return err
	}

	return nil
}

func NewStealCommand(client *client.Client) commands.Command {
	return &stealCommand{client}
}
