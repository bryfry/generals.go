package main

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
)

const (
	GENIO_SETNAME      = "error_set_username"
	GENIO_NAMETAKEN    = "This username is already taken."
	GENIO_QUEUE        = "queue_update"
	GENIO_CHAT         = "chat_message"
	GENIO_PRESTART     = "pre_game_start"
	GENIO_START        = "game_start"
	GENIO_UPDATE       = "game_update"
	GENIO_WIN          = "game_won"
	GENIO_LOSE         = "game_lost"
	GENIO_OVER         = "game_over"
	GENIO_EMPTY        = -1
	GENIO_MOUNTAIN     = -2
	GENIO_FOG          = -3
	GENIO_FOG_OBSTACLE = -4
)

type Update struct {
	Scores      []PlayerScore `json:"scores"`
	Turn        int           `json:"turn"`
	AttackIndex int           `json:"attackIndex"`
	Generals    []int         `json:"generals"`
	MapDiff     []int         `json:"map_diff"`
	CitiesDiff  []int         `json:"cities_diff"`
}

type PlayerScore struct {
	PlayerIndex int  `json:"i"`
	Dead        bool `json:"dead"`
	Army        int  `json:"total"`
	Land        int  `json:"tiles"`
}

type ChatMessage struct {
	Text        string `json:"text"`
	Username    string `json:"username"`
	Prefix      string `json:"prefix"`
	PlayerIndex int    `json:"playerIndex"`
}

type Game struct {
	UserIndex  int      `json:"playerIndex"` // YOU!
	ReplayID   string   `json:"replay_id"`
	ChatID     string   `json:"chat_room"`
	TeamChatID string   `json:"team_chat_room"`
	Usernames  []string `json:"usernames"`
	Teams      []int    `json:"teams"`
}

func DecodeEvent(p []byte, updates chan<- Update) (err error) {
	var (
		m []json.RawMessage
		t string
	)
	dec := json.NewDecoder(bytes.NewBuffer(p))
	err = dec.Decode(&m)
	if err != nil {
		l.Error(zap.Error(err))
		return err
	}
	err = json.Unmarshal(m[0], &t)
	if err != nil {
		l.Error(zap.String("type", t), zap.Error(err))
		return err
	}

	switch t {

	case GENIO_SETNAME:
		var msg string

		err = json.Unmarshal(m[1], &msg)
		if err != nil {
			l.Error(zap.String("type", t), zap.Error(err))
			return err
		}
		if msg == GENIO_NAMETAKEN {
			// This message doesn't really make sense.
			// You get this "error" even if you are the owner of the name
			l.Info("Set Name Failed (expected) ", zap.String("msg", msg))
		} else if msg == "" {
			l.Info("Set Name Success")
		}

	case GENIO_CHAT:
		var room string
		var chat ChatMessage

		err = json.Unmarshal(m[1], &room)
		if err != nil {
			l.Error(zap.String("type", t), zap.Error(err))
			return err
		}
		err = json.Unmarshal(m[2], &chat)
		if err != nil {
			l.Error(zap.String("type", t), zap.Error(err))
			return err
		}
		l.Infof("Chat - %s: %s", chat.Username, chat.Text)

	case GENIO_QUEUE:
		// TODO if this is important someday
		l.Info("Queue Update ", string(p))

	case GENIO_PRESTART:
		l.Info("Game Starting")

	case GENIO_START:
		var game Game

		err = json.Unmarshal(m[1], &game)
		if err != nil {
			l.Error(zap.String("type", t), zap.Error(err))
			return err
		}
		l.Info("Game Start ", zap.Any("game", game))

	case GENIO_UPDATE:
		var u Update

		err = json.Unmarshal(m[1], &u)
		if err != nil {
			l.Error(zap.String("type", t), zap.Error(err))
			return err
		}
		updates <- u
		l.Debug("Update sent to AI ")

	case GENIO_WIN:
		l.Info("Game Won")

	case GENIO_LOSE:
		l.Info("Game Lost")

	case GENIO_OVER:
		l.Info("Game Over")

	default:
		l.Error("Unhandled Event", zap.String("type", t), zap.String("event", string(p)))
	}
	return err
}
