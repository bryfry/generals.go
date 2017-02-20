package main

type Bot struct {
	Updates <-chan Update
	Map
}

func (ai *Bot) RecieveUpdates() {
	for {
		u := <-ai.Updates
		ai.Map.Patch(u)
		ai.Map.Print()
	}

}
