package models

import (
	"time"
	// "time_utils"
)

type User struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Address         string `json:"address"`
	UserType        string `json:"user_type"` // Applicant or Admin only string
	PasswordHash    string `json:"password_hash"`
	Password        string `json:"password,omitempty"` // error ,, remove this
	ProfileHeadline string `json:"profile_headline"`
}

type Profile struct {
	ResumeFileAddress string `json:"resume_file_address"`
	Skills            string `json:"skills"`
	Education         string `json:"education"`
	Experience        string `json:"experience"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
}

type Job struct {
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	PostedOn          time.Time `json:"posted_on"`
	TotalApplications int       `json:"total_applications"`
	CompanyName       string    `json:"company_name"`
	PostedBy          User      `json:"posted_by"`
}

type Resume struct {
	Id         int64  `json:"id,omitempty"`
	UserId     int64  `json:"user_id"`
	DocContent []byte `json:"doc_content"`
	DocType    string `json:"doc_type"`
}
