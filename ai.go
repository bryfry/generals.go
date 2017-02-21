package main

type Bot struct {
	UserID   string // AKA Password
	Username string
	Updates  <-chan Update
	Games    <-chan Game
	Game
	GeneralsIO
}

func (ai *Bot) RecieveUpdates() {
	ai.Game = <-ai.Games
	for {
		u := <-ai.Updates
		ai.Map.Patch(u)
		ai.Map.Print()
		c := ai.Game.Map.Generals[ai.Game.UserIndex]
		l.Info(c)
		err := ai.GeneralsIO.Emit("attack", c, c+1)
		if err != nil {
			l.Fatal(err)
		}
	}

}
