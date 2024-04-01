package services

import (
	"fmt"
	"log"
	"strings"

	"example.com/database"
	"example.com/models"
	"github.com/slack-go/slack"
	"gorm.io/gorm"
)

var (
	is_member            bool
	user_name            string
	project_name         string
	db                   *gorm.DB
	projects             []string
	user_project_details = []models.UsersProjectDetails{}
	project_user_details = []models.ProjectUserDetails{}
)

func HandleAddProject(add *int, text string, attachment *slack.Attachment, project_details *models.ProjectDetails, retryAdd *bool) error {
	var count int
	db = database.DB
	switch *add {
	case 1:
		// Replace the user mention with an empty string
		cleanedText := strings.Replace(text, "<@u06njt6cs4e>", "", -1)
		project_name := strings.ReplaceAll(cleanedText, " ", "")

		project_details.ProjectName = project_name

		err := db.Debug().Raw(CheckProjectExistQuery, project_name).Scan(&count).Error
		if err != nil {
			return err
		}
		if count >= 1 {
			attachment.Pretext = "Add Project"
			attachment.Text = "This project is already exist! \n You want to add this project?\n Yes or No"
			attachment.Color = "#4af030"
		} else {
			attachment.Pretext = "Add Project"
			attachment.Text = "Enter the description of the project"
			attachment.Color = "#4af030"
		}

		*add++
	case 2:

		if strings.Contains(text, "yes") {
			err := db.Debug().Exec(ActivateProjectQuery, project_details.ProjectName).Error
			if err != nil {
				attachment.Pretext = "Issue!"
				attachment.Text = "Try Again!"
				attachment.Color = "#FF0000"
				*retryAdd = true
			}
			*add += 5
			attachment.Pretext = "Thank You"
			attachment.Color = "#4af030"
		} else if strings.Contains(text, "no") {
			attachment.Pretext = "Add Project"
			attachment.Text = "Change the project name"
			attachment.Color = "#4af030"
			*retryAdd = true
			*add = 1
		} else {
			// Replace the user mention with an empty string
			project_desc := strings.Replace(text, "<@u06njt6cs4e>", "", -1)

			project_details.ProjectDescription = project_desc

			attachment.Pretext = "You Want to add Project?"
			attachment.Text = "Yes or No?"
			attachment.Color = "#4af030"
			*add++
		}

	case 3:
		if strings.Contains(text, "yes") {
			attachment.Pretext = "Thank You"
			attachment.Color = "#4af030"
			err := AddNewProject(project_details)
			if err != nil {
				attachment.Pretext = "Invalid"
				attachment.Text = "Please enter valid Text"
				attachment.Color = "#FF0000"
				return err
			}
			*add++
		} else {
			attachment.Pretext = "Welcome"
			attachment.Text = "If You want to re-enter the details then Write Hii + @Go-Bot"
			attachment.Color = "#4af030"
			*add++
		}
	}
	return nil
}

func RemoveProject(text string, attachment *slack.Attachment) error {
	db = database.DB
	cleanedText := strings.Replace(text, "<@u06njt6cs4e>", "", -1)
	project_name := strings.ReplaceAll(cleanedText, " ", "")
	err := db.Debug().Exec(RemoveProjectQuery, project_name).Error
	if err != nil {
		attachment.Pretext = "Deleted!"
		attachment.Text = "Successfully removed the " + text + " project"
		attachment.Color = "#4af030"
		return err
	}
	attachment.Pretext = "Deleted!"
	attachment.Text = "Successfully removed the " + text + " project"
	attachment.Color = "#4af030"
	return nil
}

func AddNewProject(project_details *models.ProjectDetails) error {
	db = database.DB
	project_name := removeLastSpace(project_details.ProjectName)
	project_desc := removeLastSpace(project_details.ProjectDescription)
	err := db.Exec(AddProjectQuery, project_name, project_desc).Error
	if err != nil {
		return err
	}
	return nil
}
func removeLastSpace(input string) string {
	// Trim spaces from the beginning and end of the string
	trimmed := strings.TrimSpace(input)

	// Check if the string is not empty
	if len(trimmed) > 0 {
		// Remove the last character if it is a space
		if trimmed[len(trimmed)-1] == ' ' {
			trimmed = trimmed[:len(trimmed)-1]
		}
	}

	return trimmed
}

func ManageReports(report *int, text string, attachment *slack.Attachment) error {

	if strings.Contains(text, "member") || (is_member && *report == 2) {
		MemberProjectReport(report, text, attachment, &is_member)
	} else if strings.Contains(text, "project") || (!is_member && *report == 2) {
		ProjectReport(report, text, attachment, &is_member)
	}

	return nil
}
func MemberProjectReport(report *int, text string, attachment *slack.Attachment, is_member *bool) {
	*is_member = true
	db = database.DB
	var users []string
	switch *report {
	case 1:
		err := db.Debug().Raw(GetUsersQuery).Scan(&users).Error
		if err != nil {
			log.Println("error occured while fetching the projects")
		}
		for index, val := range users {
			attachment.Text += fmt.Sprintf("%d. %s\n", index+1, val)
		}
		attachment.Pretext = "List of Users"
		attachment.Color = "#4af030"
		*report++
	case 2:
		cleanedText := strings.Replace(text, "<@u06njt6cs4e>", "", -1)
		user_name := strings.ReplaceAll(cleanedText, " ", "")
		err := db.Debug().Raw(GetMemberProjectReport, user_name).Scan(&user_project_details).Error
		if err != nil {
			log.Println("error occured while fetching the projects")
		}
		attachment.Text = "Project Id - Project Name - Log Hours"
		for _, val := range user_project_details {
			attachment.Text += fmt.Sprintf("\n%d - %s - %d", val.ProjectId, val.ProjectName, val.LogHours)
		}
		attachment.Pretext = "List of Users"
		attachment.Color = "#4af030"
		*report++
	}
}
func ProjectReport(report *int, text string, attachment *slack.Attachment, is_member *bool) {
	*is_member = false
	db = database.DB
	var projects []string
	switch *report {
	case 1:
		err := db.Debug().Raw(GetProjectsQuery).Scan(&projects).Error
		if err != nil {
			log.Println("error occured while fetching the projects")
		}
		for index, val := range projects {
			attachment.Text += fmt.Sprintf("%d. %s\n", index+1, val)
		}
		attachment.Pretext = "List of Projects"
		attachment.Color = "#4af030"
		*report++
	case 2:
		cleanedText := strings.Replace(text, "<@u06njt6cs4e>", "", -1)
		project_name := strings.ReplaceAll(cleanedText, " ", "")
		err := db.Debug().Raw(GetProjectReportQuery, project_name).Scan(&project_user_details).Error
		if err != nil {
			log.Println("error occured while fetching the projects")
		}
		attachment.Text = "User Id - User Name - Log Hours"
		for _, val := range project_user_details {
			attachment.Text += fmt.Sprintf("\n%d - %s - %d", val.UserId, val.UserName, val.LogHours)
		}
		attachment.Pretext = "List of Users"
		attachment.Color = "#4af030"
		*report++
	}
}
func Customization(custom *int, text string, attachment *slack.Attachment) error {

	if strings.Contains(text, "individual") || (is_member && *custom >= 2) {
		MemberProjectReportTimeFrame(custom, text, attachment, &is_member, &user_name)
	} else if strings.Contains(text, "project") || (!is_member && *custom >= 2) {
		ProjectReportFocusTimeFrame(custom, text, attachment, &is_member, &project_name)
	}

	return nil
}

func MemberProjectReportTimeFrame(custom *int, text string, attachment *slack.Attachment, is_member *bool, user_name *string) {
	*is_member = true
	db = database.DB
	var users []string
	switch *custom {
	case 1:
		err := db.Debug().Raw(GetUsersQuery).Scan(&users).Error
		if err != nil {
			log.Println("error occured while fetching the projects")
		}
		for index, val := range users {
			attachment.Text += fmt.Sprintf("%d. %s\n", index+1, val)
		}
		attachment.Pretext = "List of Users"
		attachment.Color = "#4af030"
		*custom++
	case 2:
		attachment.Pretext = "Time Frame"
		attachment.Text = "Select the time frame\n 1. Weekly\n 2. Monthly\n 3. Quaterly"
		attachment.Color = "#4af030"

		cleanedText := strings.Replace(text, "<@u06njt6cs4e>", "", -1)
		temp := strings.ReplaceAll(cleanedText, " ", "")
		*user_name = temp

		*custom++
	case 3:
		if strings.Contains(text, "weekly") || strings.Contains(text, "week") {
			err := db.Debug().Raw(GetMemberReportWeekly, *user_name).Scan(&user_project_details).Error
			if err != nil {
				log.Println("error occured while fetching the projects")
			}
		} else if strings.Contains(text, "monthly") || strings.Contains(text, "month") {
			err := db.Debug().Raw(GetMemberReportMonthly, *user_name).Scan(&user_project_details).Error
			if err != nil {
				log.Println("error occured while fetching the projects")
			}
		} else if strings.Contains(text, "quaterly") || strings.Contains(text, "quater") {
			err := db.Debug().Raw(GetMemberReportQuaterly, *user_name).Scan(&user_project_details).Error
			if err != nil {
				log.Println("error occured while fetching the projects")
			}
		}

		attachment.Text = "Project Id - Project Name - Log Hours"
		for _, val := range user_project_details {
			attachment.Text += fmt.Sprintf("\n%d - %s - %d", val.ProjectId, val.ProjectName, val.LogHours)
		}
		attachment.Pretext = "List of Users"
		attachment.Color = "#4af030"
		*custom++
	}
}
func ProjectReportFocusTimeFrame(custom *int, text string, attachment *slack.Attachment, is_member *bool, project_name *string) {
	*is_member = false
	db = database.DB
	var projects []string
	switch *custom {
	case 1:
		err := db.Debug().Raw(GetProjectsQuery).Scan(&projects).Error
		if err != nil {
			log.Println("error occured while fetching the projects")
		}
		for index, val := range projects {
			attachment.Text += fmt.Sprintf("%d. %s\n", index+1, val)
		}
		attachment.Pretext = "List of Projects"
		attachment.Color = "#4af030"
		*custom++
	case 2:
		attachment.Pretext = "Time Frame"
		attachment.Text = "Select the time frame\n 1. Weekly\n 2. Monthly\n 3. Quaterly"
		attachment.Color = "#4af030"

		cleanedText := strings.Replace(text, "<@u06njt6cs4e>", "", -1)
		temp := strings.ReplaceAll(cleanedText, " ", "")
		*project_name = temp

		*custom++
	case 3:
		if strings.Contains(text, "weekly") || strings.Contains(text, "week") {
			err := db.Debug().Raw(GetProjectReportWeekly, project_name).Scan(&project_user_details).Error
			if err != nil {
				log.Println("error occured while fetching the projects")
			}
		} else if strings.Contains(text, "monthly") || strings.Contains(text, "month") {
			err := db.Debug().Raw(GetProjectReportMonthly, project_name).Scan(&project_user_details).Error
			if err != nil {
				log.Println("error occured while fetching the projects")
			}
		} else if strings.Contains(text, "quaterly") || strings.Contains(text, "quater") {
			err := db.Debug().Raw(GetProjectReportQuaterly, project_name).Scan(&project_user_details).Error
			if err != nil {
				log.Println("error occured while fetching the projects")
			}
		}

		attachment.Text = "User Id - User Name - Log Hours"
		for _, val := range project_user_details {
			attachment.Text += fmt.Sprintf("\n%d - %s - %d", val.UserId, val.UserName, val.LogHours)
		}
		attachment.Pretext = "List of Users"
		attachment.Color = "#4af030"
		*custom++
	}
}
