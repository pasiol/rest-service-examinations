package main

import (
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type (
	Examination struct {
		Id                    int64     `bson:"id" json:"id" validate:"required"`
		Url                   string    `bson:"url" json:"url" validate:"required"`
		Exam                  string    `bson:"exam" json:"tutkinto" validate:"required"`
		MarketingName         string    `bson:"marketingName" json:"markkinointiNimi" validate:"required"`
		Campus                string    `bson:"campus" json:"kampus" validate:"required"`
		StartDate             string    `bson:"startDate" json:"alkaa" validate:"required"`
		EndDate               string    `bson:"endDate" json:"päättyy" validate:"required"`
		RegistrationStartDate string    `bson:"registrationStartDate" json:"ilmoittautuminenAlkaa" validate:"required"`
		RegistrationEndDate   string    `bson:"registrationEndDate" json:"ilmoittautuminenPäättyy" validate:"required"`
		Form                  string    `bson:"form" json:"lomake" validate:"required"`
		TargetGroup           string    `bson:"targetGroup" json:"kohderyhmä"`
		UpdatedAt             time.Time `bson:"updateddAt" json:"päivitetty" validate:"required"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
