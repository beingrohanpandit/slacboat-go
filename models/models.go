package models

type WorkDetails struct {
	UserName    string
	ProjectName string
	WorkDetails string
	LogHours    int
}

type ProjectDetails struct {
	ProjectName        string
	ProjectDescription string
}

type UsersProjectDetails struct {
	LogHours    int
	ProjectId   int
	ProjectName string
}
type ProjectUserDetails struct {
	UserId   int
	UserName string
	LogHours int
}
