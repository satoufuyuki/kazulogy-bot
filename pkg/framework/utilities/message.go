package utilities

import (
	"bytes"
	"github.com/nfnt/resize"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.uber.org/zap/buffer"
	"golang.org/x/image/webp"
	"image"
	"image/jpeg"
	"strings"
)

func FindMessageContentFromMessage(v *waE2E.Message) string {
	switch {
	case v.GetConversation() != "":
		return v.GetConversation()
	case v.GetExtendedTextMessage() != nil:
		return v.GetExtendedTextMessage().GetText()
	case v.GetImageMessage() != nil:
		return v.GetImageMessage().GetCaption()
	case v.GetVideoMessage() != nil:
		return v.GetVideoMessage().GetCaption()
	case v.GetDocumentMessage() != nil:
		return v.GetDocumentMessage().GetCaption()
	case v.GetGroupInviteMessage() != nil:
		return v.GetGroupInviteMessage().GetCaption()
	case v.GetLiveLocationMessage() != nil:
		return v.GetLiveLocationMessage().GetCaption()
	case v.GetViewOnceMessage() != nil:
		return FindMessageContentFromMessage(v.GetViewOnceMessage().GetMessage())
	case v.GetViewOnceMessageV2() != nil:
		return FindMessageContentFromMessage(v.GetViewOnceMessageV2().GetMessage())
	case v.GetEphemeralMessage() != nil:
		return FindMessageContentFromMessage(v.GetEphemeralMessage().GetMessage())
	}

	return ""
}

func FindThumbnailFromMessage(v *waE2E.Message) []byte {
	switch {
	case v.GetImageMessage() != nil:
		return v.GetImageMessage().GetJPEGThumbnail()
	case v.GetStickerMessage() != nil:
		return v.GetStickerMessage().GetPngThumbnail()
	case v.GetVideoMessage() != nil:
		return v.GetVideoMessage().GetJPEGThumbnail()
	case v.GetViewOnceMessage() != nil:
		return FindThumbnailFromMessage(v.GetViewOnceMessage().GetMessage())
	case v.GetViewOnceMessageV2() != nil:
		return FindThumbnailFromMessage(v.GetViewOnceMessageV2().GetMessage())
	case v.GetEphemeralMessage() != nil:
		return FindThumbnailFromMessage(v.GetEphemeralMessage().GetMessage())
	}

	return nil
}

func GenerateThumbnail(buff []byte, mimetype string) (res []byte, err error) {
	var thumbnail buffer.Buffer

	var img image.Image
	switch {
	case strings.HasSuffix(mimetype, "webp"):
		img, err = webp.Decode(bytes.NewReader(buff))
		if err != nil {
			return nil, err
		}
		break
	default:
		img, _, err = image.Decode(bytes.NewReader(buff))
		if err != nil {
			return nil, err
		}
	}

	m := resize.Thumbnail(72, 72, img, resize.Lanczos3)
	if err = jpeg.Encode(&thumbnail, m, nil); err != nil {
		return nil, err
	}

	return thumbnail.Bytes(), err
}
