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
	user_id := "k_00004"
	username := "[Bot]k_0004"

	var ai Bot
	var g GeneralsIO
	u, games := make(chan Update), make(chan Game)
	ai.Updates, ai.Games = u, games
	g.Updates, g.Games = u, games
	g.Connect()

	err := g.Emit("set_username", user_id, username)
	if err != nil {
		l.Fatal(err)
	}
	err = g.Emit("join_private", "test", user_id)
	if err != nil {
		l.Fatal(err)
	}
	err = g.Emit("set_force_start", "test", true)
	if err != nil {
		l.Fatal(err)
	}
	go g.Ping(5)
	go g.RecievePackets() // Packets > Messages > Events
	go ai.RecieveUpdates()
	exitHandler()
}
