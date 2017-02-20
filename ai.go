package main

type Bot struct {
	Updates <-chan Update
	Games   <-chan Game
	Game
}

func (ai *Bot) RecieveUpdates() {
	ai.Game = <-ai.Games
	for {
		u := <-ai.Updates
		ai.Map.Patch(u)
		ai.Map.Print()
	}

}
