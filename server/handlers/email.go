package handlers

import (
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/dchote/survey-engine/config"
	//"github.com/dchote/survey-engine/models"

	"github.com/labstack/echo"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	sendgridClient *sendgrid.Client
)

// EmailTo struct
type EmailTo struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
}

// EmailNotificationRequest struct
type EmailNotificationRequest struct {
	Subject       string    `json:"subject"`
	Content       string    `json:"message"`
	AdditionalCCs []EmailTo `json:"additionalCCs"`
}

// InitEmailHandler function
func InitEmailHandler() {
	if config.Config.SendGrid.APIKey != "" && config.Config.SendGrid.FromName != "" && config.Config.SendGrid.FromAddress != "" {
		sendgridClient = sendgrid.NewSendClient(config.Config.SendGrid.APIKey)
	}
}

// SendNotification function
func SendNotification() echo.HandlerFunc {
	return func(c echo.Context) error {
		notificationRequest := new(EmailNotificationRequest)
		if err := c.Bind(notificationRequest); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if sendgridClient != nil {

			if config.Config.Survey.ContactName != "" && config.Config.Survey.ContactEmail != "" {
				email := mail.NewV3Mail()

				to := mail.NewEmail(config.Config.Survey.ContactName, config.Config.Survey.ContactEmail)

				p := mail.NewPersonalization()
				p.AddTos(to)
				email.AddPersonalizations(p)

				// build list of ccs
				var ccs []*mail.Email

				// fetch all users and add users with notifications enabled
				users := userManager.ListUsers()
				for _, user := range users {
					if user.Notifications && user.Email != config.Config.Survey.ContactEmail {
						ccs = append(ccs, mail.NewEmail(user.FullName, user.Email))
					}
				}

				// add any CCs in the http request
				if len(notificationRequest.AdditionalCCs) > 0 {
					for _, cc := range notificationRequest.AdditionalCCs {
						if cc.Email != config.Config.Survey.ContactEmail {
							ccs = append(ccs, mail.NewEmail(cc.FullName, cc.Email))
						}
					}
				}

				if len(ccs) > 0 {
					p.AddCCs(ccs...)
				}

				from := mail.NewEmail(config.Config.SendGrid.FromName, config.Config.SendGrid.FromAddress)
				email.SetFrom(from)

				email.Subject = notificationRequest.Subject

				plainText := mail.NewContent("text/plain", notificationRequest.Content)
				email.AddContent(plainText)

				log.Println("Sending Email to " + config.Config.Survey.ContactEmail)
				response, err := sendgridClient.Send(email)

				log.Println("Sendgrid status code was " + strconv.Itoa(response.StatusCode))

				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}

				if response.StatusCode > 299 {
					return echo.NewHTTPError(http.StatusInternalServerError, response.Body)
				}

				return c.JSON(http.StatusOK, true)
			}

		}

		return c.JSON(http.StatusOK, false)
	}
}
