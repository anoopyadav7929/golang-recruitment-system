package routes

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	// "golang-project/utils"
	"bytes"
	"encoding/json"
	"golang-project/models"
	"net/http"
)

func UploadResume(c *gin.Context) {
    // Authenticate and check user type (implement authentication logic)

    file, _, err := c.Request.FormFile("resume")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
        return
    }

    filePath := "uploads/resumes/"// + file.Filename   checkkkkkkkkk--------------
    out, err := os.Create(filePath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
        return
    }
    defer out.Close()

    _, err = io.Copy(out, file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
        return
    }

    // Call third-party API to process the resume
    resumeData, err := processResume(filePath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not process resume"})
        return
    }

    // Save the processed resume data to the database (implement database logic)
    var profile models.Profile
    profile.ResumeFileAddress = filePath
    profile.Skills = resumeData.Skills
    profile.Education = resumeData.Education
    profile.Experience = resumeData.Experience
    profile.Name = resumeData.Name
    profile.Email = resumeData.Email
    profile.Phone = resumeData.Phone

    c.JSON(http.StatusOK, gin.H{"message": "Resume uploaded", "profile": profile})
}

func processResume(filePath string) (models.Profile, error) {
    var profile models.Profile

    file, err := os.Open(filePath)
    if err != nil {
        return profile, err
    }
    defer file.Close()

    req, err := http.NewRequest("POST", "https://api.apilayer.com/resume_parser/upload", file)
    if err != nil {
        return profile, err
    }
    req.Header.Add("Content-Type", "application/octet-stream")
    req.Header.Add("apikey", "YOUR_API_KEY")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return profile, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return profile, fmt.Errorf("failed to parse resume, status code: %d", resp.StatusCode)
    }

    var buffer bytes.Buffer
    _, err = io.Copy(&buffer, resp.Body)
    if err != nil {
        return profile, err
    }

    var resumeData struct {
        Skills     []string `json:"skills"`
        Education  []struct {
            Name string `json:"name"`
        } `json:"education"`
        Experience []struct {
            Name string `json:"name"`
        } `json:"experience"`
        Name  string `json:"name"`
        Email string `json:"email"`
        Phone string `json:"phone"`
    }

    err = json.Unmarshal(buffer.Bytes(), &resumeData)
    if err != nil {
        return profile, err
    }

    profile.Skills = strings.Join(resumeData.Skills, ", ")
    // profile.Education = strings.Join(resumeData.Education, ", ")
    // profile.Experience = strings.Join(resumeData.Experience, ", ")
    profile.Name = resumeData.Name
    profile.Email = resumeData.Email
    profile.Phone = resumeData.Phone

    return profile, nil
}
