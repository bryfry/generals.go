package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
)

var l *zap.SugaredLogger

func init() {
	conf := zap.NewDevelopmentConfig()
	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.InfoLevel)
	conf.Level = level
	conf.DisableStacktrace = true
	logger, _ := conf.Build()
	l = logger.Sugar()
}

func exitHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		l.Warn("Caught SIGINT, Exiting")
		return
	}
}

func main() {

	u, games := make(chan Update), make(chan Game)
	g := GeneralsIO{
		Updates: u,
		Games:   games,
	}
	ai := Bot{
		UserID:     "k_00004",
		Username:   "[Bot]k_0004",
		Updates:    u,
		Games:      games,
		GeneralsIO: g,
	}
	ai.GeneralsIO.Connect()

	err := ai.GeneralsIO.Emit("set_username", ai.UserID, ai.Username)
	if err != nil {
		l.Fatal(err)
	}
	err = ai.GeneralsIO.Emit("join_private", "test", ai.UserID)
	if err != nil {
		l.Fatal(err)
	}
	err = ai.GeneralsIO.Emit("set_force_start", "test", true)
	if err != nil {
		l.Fatal(err)
	}
	go ai.GeneralsIO.Ping(5)
	go ai.GeneralsIO.RecievePackets() // Packets > Messages > Events
	go ai.RecieveUpdates()
	exitHandler()
}
