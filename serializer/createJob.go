package serializer

import (
    "errors"
    "golang-project/models"
)

func ValidateJob(job *models.Job) error {
    if job.Title == "" {
        return errors.New("title is required")
    }
    if job.Description == "" {
        return errors.New("description is required")
    }
    if job.CompanyName == "" {
        return errors.New("company name is required")
    }
    return nil
}
