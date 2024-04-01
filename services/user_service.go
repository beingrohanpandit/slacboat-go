package services

import (
	"fmt"
	"log"
	"strings"

	"example.com/database"
	"example.com/models"
	"github.com/slack-go/slack"
)

// Insert the work details into the data base
func AddWorkDetails(i *int, text string, attachment *slack.Attachment, details *models.WorkDetails) error {
	db = database.DB
	switch *i {
	case 1:
		// Get the projects from the database
		err := db.Debug().Raw(GetProjectsQuery).Scan(&projects).Error
		if err != nil {
			log.Println("error occured while fetching the projects")
			return err
		}

		attachment.Text = "Choose the relevant project from a list"
		// Print the list
		for index, val := range projects {
			attachment.Text += fmt.Sprintf("\n %d. %s", index+1, val)
		}
		attachment.Pretext = "Select Project"
		attachment.Color = "#4af030"
		*i++

		// Replace the user mention with an empty string
		cleanedText := strings.Replace(text, "<@u06njt6cs4e>", "", -1)
		details.WorkDetails = cleanedText

	case 2:
		attachment.Pretext = "Log Hours"
		attachment.Text = "Enter the number of hours worked on the project"
		attachment.Color = "#4af030"

		// Replace the user mention with an empty string
		cleanedText := strings.Replace(text, "<@u06njt6cs4e>", "", -1)
		details.ProjectName = cleanedText

		*i++
	case 3:
		attachment.Pretext = "Submission Confirmation"
		attachment.Text = "If are details are correct then Yes otherwise No!"
		attachment.Color = "#4af030"

		// Replace the user mention with an empty string
		cleanedText := strings.Replace(text, "<@u06njt6cs4e>", "", -1)
		val, err := ExtractNumeric(cleanedText)
		if err != nil {
			return err
		}
		details.LogHours = val

		*i++
	case 4:
		if strings.Contains(text, "yes") {
			attachment.Pretext = "Thank You"
			attachment.Color = "#4af030"
			err := InsertWorkData(details)
			if err != nil {
				attachment.Pretext = "Invalid"
				attachment.Text = "Please enter valid Text"
				attachment.Color = "#FF0000"
				return err
			}

			*i++
		} else {
			attachment.Pretext = "Welcome"
			attachment.Text = "If You want to re-enter the details then Write Hii + @Go-Bot"
			attachment.Color = "#4af030"
			*i++
		}
	}

	return nil
}

func InsertWorkData(details *models.WorkDetails) error {
	var user_id, project_id int

	// get the user id from the username
	err := db.Debug().Raw(GetUserByNameQuery, details.UserName).Scan(&user_id).Error
	if err != nil {
		return err
	}

	project_name := strings.ReplaceAll(details.ProjectName, " ", "")

	// get the project id from the project name
	err = db.Debug().Raw(GetProjectByNameQuery, project_name).Scan(&project_id).Error
	if err != nil {
		return err
	}

	// Insert values into the work_details table
	err = db.Debug().Exec(AddWorkDetailsQuery, user_id, project_id, details.LogHours, details.WorkDetails).Error
	if err != nil {
		return err
	}
	return nil
}
