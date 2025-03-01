package main

import (
	"fmt"
)

func main() {
	connectionString := "TODO!"
	context, err := NewContext(connectionString)
	if err != nil {
		fmt.Println(err.Error())
	}

	if err := context.AddPlayers([]string{"Joe", "Mark", "Jim", "Dave"}); err != nil {
		fmt.Println(err)
	}

	if err := context.SetDecisiveTeamMatch([]string{"Mark", "Jim"}, []string{"Joe", "Dave"}); err != nil {
		fmt.Println(err)
	}

	if err := context.SetDrawMatch("Mark", "Jim"); err != nil {
		fmt.Println(err)
	}

	result, err := context.GetPlayersOrderBy(Elo, Descending)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("All players:")
	fmt.Println(result)
	fmt.Println("Best player is:")
	fmt.Println(result[0].Key)
}
