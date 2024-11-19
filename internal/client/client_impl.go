package client

import (
	"context"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"github.com/satoufuyuki/kazulogy-bot/pkg/framework/config"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
)

type Client struct {
	*whatsmeow.Client
}

func Connect(lc fx.Lifecycle, client *Client, cmdHandler whatsmeow.EventHandler) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) (err error) {
				go func() {
					client.AddEventHandler(cmdHandler)
					if client.Store.ID == nil {
						// No ID stored, new login
						qrChan, _ := client.GetQRChannel(context.Background())
						err = client.Connect()
						if err != nil {
							zap.L().Fatal("couldn't connect to WhatsApp", zap.Error(err))
						}

						for evt := range qrChan {
							if evt.Event == "code" {
								// Render the QR code here
								// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
								// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
								qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
							} else {
								zap.S().Info("login event:", evt.Event)
							}
						}
					} else {
						// Already logged in, just connect
						err = client.Connect()
						if err != nil {
							zap.L().Fatal("couldn't connect to WhatsApp", zap.Error(err))
						}
					}
				}()

				return
			},
			OnStop: func(ctx context.Context) error {
				client.Disconnect()
				return nil
			},
		})
}

func New(config config.Config) *Client {
	dbLog := waLog.Stdout("Database", config.DatabaseLogLevel, true)
	container, err := sqlstore.New("sqlite3", "file:data/bot.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", config.ClientLogLevel, true)
	return &Client{whatsmeow.NewClient(deviceStore, clientLog)}
}
