package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"example.com/database"
	"example.com/models"
	"example.com/services"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"gorm.io/gorm"
)

var (
	i               int
	remove          int
	add             int
	report          int
	custom          int
	retryAdd        bool
	db              *gorm.DB
	projects        []string
	details         models.WorkDetails
	project_details models.ProjectDetails
)

func init() {
	database.ConnectDatabase()
	db = database.DB
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("error : ", err)
		return
	}
	token := os.Getenv("SLACK_BOT_TOKEN")
	// channelID := os.Getenv("Channel_ID")
	webSocketToken := os.Getenv("Web_Socket")

	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(webSocketToken))
	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-socketClient.Events:
				switch event.Type {
				case socketmode.EventTypeEventsAPI:
					eventApi, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						return
					}
					socketClient.Ack(*event.Request)
					err := HandleEventMessage(eventApi, client)
					if err != nil {
						log.Fatal(err)
						return
					}
				}
			}
		}
	}(ctx, client, socketClient)
	socketClient.Run()
}
func HandleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
	switch event.Type {
	case slackevents.CallbackEvent:
		innerEvent := event.InnerEvent
		switch evnt := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			err := HandleAppMentionEventToBot(evnt, client)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported error")
	}
	return nil
}
func HandleAppMentionEventToBot(event *slackevents.AppMentionEvent, client *slack.Client) error {
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}
	text := strings.ToLower(event.Text)
	attachment := slack.Attachment{}

	if strings.Contains(text, "help") {
		if user.IsAdmin {
			attachment.Pretext = "HELP!"
			attachment.Text = "1. Add projects - add + @go-bot\n 2. Remove projects - remove + @go-bot\n 3. Request member-wise report - report + go-bot\n 4. Generate report in timeframe - custom + @go-bot"
			attachment.Color = "#4af030"
		} else {
			attachment.Pretext = "HELP!"
			attachment.Text = "1. Add work details - work + @go-bot\n "
			attachment.Color = "#4af030"
		}
		// } else if strings.Contains(text, "work") && !user.IsAdmin {
	} else if strings.Contains(text, "work ") {
		attachment.Pretext = "Add Work Details"
		attachment.Text = "Provide a description of the work performed!"
		attachment.Color = "#4af030"

		details.UserName = user.Name
		i = 1
	} else if (strings.Contains(text, "add") || retryAdd) && user.IsAdmin {
		attachment.Pretext = "Add Project"
		attachment.Text = "Enter the name of the project"
		attachment.Color = "#4af030"
		add = 1
		retryAdd = false
	} else if strings.Contains(text, "remove") && user.IsAdmin {
		err := db.Debug().Raw(services.GetProjectsQuery).Scan(&projects).Error
		if err != nil {
			return err
		}

		attachment.Text = "Choose the relevant project from a list"
		// Print the list
		for index, val := range projects {
			attachment.Text += fmt.Sprintf("\n %d. %s", index+1, val)
		}
		attachment.Pretext = "Select Project"
		attachment.Color = "#4af030"
		remove = 1
	} else if (strings.Contains(text, "reports") || strings.Contains(text, "report")) && user.IsAdmin {
		attachment.Pretext = "Request Report"
		attachment.Text = "member-wise or project-wise report?"
		attachment.Color = "#4af030"
		report = 1
	} else if (strings.Contains(text, "custom") || strings.Contains(text, "customization")) && user.IsAdmin {
		attachment.Pretext = "Customization"
		attachment.Text = "Individual or Project?"
		attachment.Color = "#4af030"
		custom = 1
	} else {
		if i != 0 {
			err := services.AddWorkDetails(&i, text, &attachment, &details)
			if err != nil {
				return err
			}
		}
		if remove == 1 {
			err := services.RemoveProject(text, &attachment)
			if err != nil {
				return err
			}
			remove++
		}
		if add != 0 {
			err := services.HandleAddProject(&add, text, &attachment, &project_details, &retryAdd)
			if err != nil {
				return err
			}
		}
		if report != 0 {
			err := services.ManageReports(&report, text, &attachment)
			if err != nil {
				return err
			}
		}
		if custom != 0 {
			err := services.Customization(&custom, text, &attachment)
			if err != nil {
				return err
			}
		}
	}

	_, _, err = client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message : %w", err)
	}
	return nil
}
