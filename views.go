package main

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (a *App) getHealthz(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	err := a.Db.Client().Ping(ctx, readpref.Primary())
	defer cancel()
	if err != nil {
		return c.String(http.StatusInternalServerError, "non-operational")
	}
	return c.String(http.StatusOK, "operational")
}

func (a *App) getExaminations(c echo.Context) error {
	exams, err := examinations(a.Db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSONPretty(http.StatusOK, exams, "    ")
}
