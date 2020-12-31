package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/mikkel1156/Nextcloud-TrelloSync/internal/config"
	"github.com/mikkel1156/Nextcloud-TrelloSync/pkg/api/deck"
	"github.com/mikkel1156/Nextcloud-TrelloSync/pkg/api/trello"
	"github.com/sirupsen/logrus"
)

var conf config.Config

func syncTrello(sg config.SyncGroup) {
	logger := log.WithFields(logrus.Fields{
		"sg": sg.Name,
	})
	trelloClient := trello.NewClient(conf.Trello.Key, conf.Trello.Token)
	deckClient := deck.NewClient(conf.Nextcloud.API, conf.Nextcloud.Username, conf.Nextcloud.Password)

	trelloCards, err := trelloClient.GetCards(sg.Trello.List)
	if err != nil {
		logger.Error(err)
		return
	}

	deckStack, err := deckClient.GetStack(sg.Nextcloud.Board, sg.Nextcloud.Stack)
	if err != nil {
		logger.Error(err)
		return
	}

	for _, deckCard := range deckStack.Cards {
		matched := false
		var match trello.Card
		for _, trelloCard := range trelloCards {
			if trelloCard.Name == deckCard.Title {
				match = trelloCard
				matched = true
				break
			}
		}

		if !matched {
			card, err := trelloClient.NewCard(sg.Trello.List, deckCard.Title, deckCard.Description, nil)
			if err != nil {
				logger.Error(err)
				return
			}
			logger.Debug("created new trello card: " + card.Name)
		} else {
			if match.Description != deckCard.Description {
				trelloClient.UpdateCard(match.ID, match.Name, deckCard.Description, match.Labels)
				logger.Debug("updated trello card: " + match.Name)
			}
		}
	}

	for _, trelloCard := range trelloCards {
		found := false
		for _, deckCard := range deckStack.Cards {
			if deckCard.Title == trelloCard.Name {
				found = true
			}
		}

		if !found {
			trelloClient.DeleteCard(trelloCard.ID)
			logger.Debug("deleted old trello card: " + trelloCard.Name)
		}
	}

}

func task() {
	log.Info("syncronizing now.")
	for _, sg := range conf.Sync {
		switch sg.Type {
		case "trello":
			go syncTrello(sg)
		default:
			log.Warn(fmt.Sprintf("sync group '%s' has invalid sync type '%s'", sg.Name, sg.Type))
		}
	}
}

var log = logrus.New()

func main() {
	done := make(chan struct{})

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		close(done)
	}()

	conf = config.Load()

	if conf.Log != "" {
		f, err := os.Create(conf.Log)
		if err != nil {
			log.Panic(err)
		}
		defer f.Close()

		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)
	}

	if conf.Debug {
		log.SetLevel(logrus.DebugLevel)
	}

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(conf.Interval).Minute().Do(task)
	scheduler.StartAsync()
	defer scheduler.Clear()
	defer scheduler.Stop()

	<-done
	log.Info("exiting")
}
