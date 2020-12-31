package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mikkel1156/Nextcloud-TrelloSync/internal/config"
	"github.com/mikkel1156/Nextcloud-TrelloSync/pkg/api/deck"
	"github.com/mikkel1156/Nextcloud-TrelloSync/pkg/api/trello"
)

var conf config.Config

func displayNextcloud() {
	client := deck.NewClient(conf.Nextcloud.API, conf.Nextcloud.Username, conf.Nextcloud.Password)
	boards, err := client.GetBoards()
	if err != nil {
		panic(err)
	}
	for _, board := range boards {
		fmt.Printf("%s [%d]\n", board.Title, board.ID)
		stacks, err := client.GetStacks(board.ID)
		if err != nil {
			panic(err)
		}
		for _, stack := range stacks {
			fmt.Printf("- %s [%d]\n", stack.Title, stack.ID)
		}
	}
}

func displayTrello() {
	client := trello.NewClient(conf.Trello.Key, conf.Trello.Token)
	boards, err := client.GetBoards()
	if err != nil {
		panic(err)
	}
	for _, board := range boards {
		fmt.Printf("%s [%s]\n", board.Name, board.ID)
		stacks, err := client.GetLists(board.ID)
		if err != nil {
			panic(err)
		}
		for _, stack := range stacks {
			fmt.Printf("- %s [%s]\n", stack.Name, stack.ID)
		}
	}
}

func main() {
	conf = config.Load()

	args := os.Args
	if len(args) != 2 {
		fmt.Println("[Trello]")
		displayTrello()

		fmt.Println("")
		fmt.Println("[Nextcloud]")
		displayNextcloud()
		return
	}

	switch strings.ToLower(args[1]) {
	case "deck":
		displayNextcloud()
	case "nextcloud":
		displayNextcloud()
	case "trello":
		displayTrello()
	}
}
