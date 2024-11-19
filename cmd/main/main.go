package cmd

import (
	"github.com/satoufuyuki/kazulogy-bot/internal/client"
	"github.com/satoufuyuki/kazulogy-bot/internal/commands"
	"github.com/satoufuyuki/kazulogy-bot/internal/commands/utility"
	"github.com/satoufuyuki/kazulogy-bot/pkg/framework/config"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func Providers() fx.Option {
	return fx.Options(
		fx.Provide(
			zap.NewDevelopment,
			config.NewConfig,
			client.New,

			// Commands
			fx.Annotate(utility.NewStealCommand, fx.ResultTags(`group:"commands"`)),
			fx.Annotate(utility.NewPingCommand, fx.ResultTags(`group:"commands"`)),
			fx.Annotate(commands.CommandHandler, fx.ParamTags("", `group:"commands"`)),
		),
	)
}

func Entrypoint() fx.Option {
	return fx.Invoke(
		client.Connect,
	)
}

func Run() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		Providers(), Entrypoint()).Run()
}
