package handlers

import (
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dchote/survey-engine/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

var (
	surveyManager *models.SurveyManager
)

// SurveySettingsObject struct
type SurveySettingsObject struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// SurveyFullObject struct
type SurveyFullObject struct {
	Name       string `json:"name"`
	Modified   string `json:"modified"`
	Author     string `json:"author"`
	SurveyJSON string `json:"surveyJSON"`
}

// SurveyListObject struct
type SurveyListObject struct {
	Name     string `json:"name"`
	Modified string `json:"modified"`
	Author   string `json:"author"`
}

// SurveyObjectList struct
type SurveyObjectList struct {
	Surveys []*SurveyListObject `json:"surveys"`
}

// SurveyUpdateRequest struct
type SurveyUpdateRequest struct {
	SurveyJSON string `json:"surveyJSON"`
}

// SurveySessionRequest struct
type SurveySessionRequest struct {
	SurveyType   string `json:"surveyType"`
	Username     string `json:"username"`
	IsTest       bool   `json:"isTest"`
	ReferralCode string `json:"referralCode"`
}

// SurveySessionUpdateRequest struct
type SurveySessionUpdateRequest struct {
	SurveyState map[string]interface{} `json:"surveyState"`
}

// SurveySessionObject struct
type SurveySessionObject struct {
	SessionID    string                 `json:"sessionID"`
	Created      string                 `json:"created"`
	SurveyType   string                 `json:"surveyType"`
	SurveyState  map[string]interface{} `json:"surveyState"`
	IsTest       bool                   `json:"isTest"`
	ReferralCode string                 `json:"referralCode"`
}

// SurveySessionFullObject struct
type SurveySessionFullObject struct {
	SessionID    string                 `json:"sessionID"`
	SurveyType   string                 `json:"surveyType"`
	SurveyState  map[string]interface{} `json:"surveyState"`
	IsTest       bool                   `json:"isTest"`
	ReferralCode string                 `json:"referralCode"`
	Username     string                 `json:"username"`
	Created      string                 `json:"created"`
	Modified     string                 `json:"modified"`
}

// SurveySessionFullObjectList struct
type SurveySessionFullObjectList struct {
	SurveySessions []SurveySessionFullObject `json:"surveySessions"`
}

// SurveySessionDataRequest struct
type SurveySessionDataRequest struct {
	ReferenceID string                 `json:"referenceID"`
	SessionID   string                 `json:"sessionID"`
	SurveyName  string                 `json:"surveyName"`
	DataKey     string                 `json:"dataKey"`
	Data        map[string]interface{} `json:"data"`
	Username    string                 `json:"username"`
}

// SurveySessionDataUpdateRequest struct
type SurveySessionDataUpdateRequest struct {
	Data map[string]interface{} `json:"data"`
}

// SurveySessionDataObject struct
type SurveySessionDataObject struct {
	SessionDataID string                 `json:"sessionDataID"`
	ReferenceID   string                 `json:"referenceID"`
	SessionID     string                 `json:"sessionID"`
	SurveyName    string                 `json:"surveyName"`
	DataKey       string                 `json:"dataKey"`
	Data          map[string]interface{} `json:"data"`
	Username      string                 `json:"username"`
	Created       string                 `json:"created"`
	Modified      string                 `json:"modified"`
}

// SurveySessionDataObjectList struct
type SurveySessionDataObjectList struct {
	SessionData []SurveySessionDataObject `json:"sessionData"`
}

// SurveyExportObject struct
type SurveyExportObject struct {
	Name   string        `json:"name"`
	Export []interface{} `json:"export"`
}

// InitSurveyHandler initialize the survey manager
func InitSurveyHandler(db *gorm.DB) {
	surveyManager = models.NewSurveyManager(db)
}

// SurveySettings return the survey settings
func SurveySettings() echo.HandlerFunc {
	return func(c echo.Context) error {
		surveySettings, err := surveyManager.SurveySettings()

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var config map[string]interface{}
		json.Unmarshal([]byte(surveySettings.Config), &config)

		surveySettingsResponse := SurveySettingsObject{
			Type:   surveySettings.Type,
			Config: config,
		}

		return c.JSON(http.StatusOK, surveySettingsResponse)
	}
}

// UpdateSurveySettings update the survey settings
func UpdateSurveySettings() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ValidateAdmin(c.Get("user").(*jwt.Token))

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Permission Denied")
		}

		surveySettingsUpdate := new(SurveySettingsObject)
		if err = c.Bind(surveySettingsUpdate); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		config, err := json.Marshal(surveySettingsUpdate.Config)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		surveyManager.UpdateSurveySettings(surveySettingsUpdate.Type, []byte(config))

		return c.JSON(http.StatusOK, surveySettingsUpdate)
	}
}

// ListSurveys list all surveys
func ListSurveys() echo.HandlerFunc {
	return func(c echo.Context) error {
		surveys := surveyManager.AllSurveys()

		surveyList := new(SurveyObjectList)

		for _, survey := range surveys {
			surveyListObject := &SurveyListObject{
				Name:     survey.Name,
				Modified: survey.UpdatedAt.Format("Mon, 2 Jan 2006 3:04 PM"),
				Author:   survey.Author,
			}
			surveyList.Surveys = append(surveyList.Surveys, surveyListObject)
		}

		return c.JSON(http.StatusOK, surveyList)
	}
}

// UpdateSurvey update survey
func UpdateSurvey() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := ValidateAdmin(c.Get("user").(*jwt.Token))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Permission Denied")
		}

		surveyName := c.Param("name")

		surveyUpdateRequest := new(SurveyUpdateRequest)
		if err = c.Bind(surveyUpdateRequest); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if surveyUpdateRequest.SurveyJSON == "" {
			return echo.NewHTTPError(http.StatusInternalServerError, "surveyJSON not attached to payload")
		}

		survey := surveyManager.FindSurvey(surveyName)

		if survey.Name == "" {
			survey, err = surveyManager.NewSurvey(surveyName, user.Username, []byte(surveyUpdateRequest.SurveyJSON))

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		} else {
			survey.SurveyJSON = []byte(surveyUpdateRequest.SurveyJSON)
			surveyManager.UpdateSurvey(survey)
		}

		return c.JSON(http.StatusOK, false)
	}
}

// FetchSurvey retrieve single survey
func FetchSurvey() echo.HandlerFunc {
	return func(c echo.Context) error {
		surveyName := c.Param("name")

		if surveyName == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Please specify survey name")
		}

		survey := surveyManager.FindSurvey(surveyName)

		if survey.Name == "" {
			log.Printf("Survey '%s' not found", surveyName)
			return echo.NewHTTPError(http.StatusNotFound, "Survey not found")
		}

		surveyFullObject := &SurveyFullObject{
			Name:       survey.Name,
			Modified:   survey.UpdatedAt.Format(time.RFC822),
			Author:     survey.Author,
			SurveyJSON: string(survey.SurveyJSON),
		}

		return c.JSON(http.StatusOK, surveyFullObject)
	}
}

// DeleteSurvey delete survey
func DeleteSurvey() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ValidateAdmin(c.Get("user").(*jwt.Token))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Permission Denied")
		}

		surveyName := c.Param("name")

		survey := surveyManager.FindSurvey(surveyName)

		if survey.Name == "" {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		surveyManager.DeleteSurvey(survey)

		return c.JSON(http.StatusOK, true)
	}
}

// NewSurveySession create new survey session
func NewSurveySession() echo.HandlerFunc {
	return func(c echo.Context) error {
		surveySessionRequest := new(SurveySessionRequest)
		if err := c.Bind(surveySessionRequest); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		surveySession, err := surveyManager.NewSurveySession(surveySessionRequest.SurveyType, surveySessionRequest.IsTest, surveySessionRequest.ReferralCode, surveySessionRequest.Username, c.RealIP(), c.Request().Header.Get("User-Agent"))

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var state map[string]interface{}
		json.Unmarshal(surveySession.SurveyState, &state)

		surveySessionResponse := SurveySessionObject{
			SessionID:    surveySession.SessionID,
			Created:      surveySession.CreatedAt.Format("Mon, 2 Jan 2006 3:04 PM"),
			SurveyType:   surveySession.SurveyType,
			SurveyState:  state,
			IsTest:       surveySession.IsTest,
			ReferralCode: surveySession.ReferralCode,
		}

		return c.JSON(http.StatusOK, surveySessionResponse)
	}
}

// UpdateSurveySession update survey session
func UpdateSurveySession() echo.HandlerFunc {
	return func(c echo.Context) error {
		surveySessionID := c.Param("id")

		if surveySessionID == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Please specify surveySessionID")
		}

		surveySessionUpdateRequest := new(SurveySessionUpdateRequest)
		if err := c.Bind(surveySessionUpdateRequest); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		log.Printf("Survey State: %s", surveySessionUpdateRequest.SurveyState)

		if len(surveySessionUpdateRequest.SurveyState) == 0 {
			return echo.NewHTTPError(http.StatusInternalServerError, "surveyState not attached to payload")
		}

		surveySession := surveyManager.FindSurveySession(surveySessionID)

		if surveySession.SessionID == "" {
			log.Printf("surveySessionID: '%s' not found", surveySessionID)
			return echo.NewHTTPError(http.StatusNotFound, "surveySessionID not found")
		}

		state, err := json.Marshal(surveySessionUpdateRequest.SurveyState)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		surveySession.SurveyState = state

		_, err = surveyManager.UpdateSurveySession(surveySession)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var surveyState map[string]interface{}
		json.Unmarshal(surveySession.SurveyState, &surveyState)

		surveySessionResponse := SurveySessionObject{
			SessionID:    surveySession.SessionID,
			Created:      surveySession.CreatedAt.Format("Mon, 2 Jan 2006 3:04 PM"),
			SurveyType:   surveySession.SurveyType,
			SurveyState:  surveyState,
			IsTest:       surveySession.IsTest,
			ReferralCode: surveySession.ReferralCode,
		}

		return c.JSON(http.StatusOK, surveySessionResponse)
	}
}

// FetchSurveySession fetch particular session
func FetchSurveySession() echo.HandlerFunc {
	return func(c echo.Context) error {
		surveySessionID := c.Param("id")

		if surveySessionID == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Please specify surveySessionID")
		}

		surveySession := surveyManager.FindSurveySession(surveySessionID)

		if surveySession.SessionID == "" {
			log.Printf("surveySessionID: '%s' not found", surveySessionID)
			return echo.NewHTTPError(http.StatusNotFound, "surveySessionID not found")
		}

		var state map[string]interface{}
		json.Unmarshal(surveySession.SurveyState, &state)

		surveySessionResponse := SurveySessionObject{
			SessionID:    surveySession.SessionID,
			Created:      surveySession.CreatedAt.Format("Mon, 2 Jan 2006 3:04 PM"),
			SurveyType:   surveySession.SurveyType,
			SurveyState:  state,
			IsTest:       surveySession.IsTest,
			ReferralCode: surveySession.ReferralCode,
		}

		return c.JSON(http.StatusOK, surveySessionResponse)
	}
}

// ListSurveySessions list all survey sessions
func ListSurveySessions() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := ValidateUser(c.Get("user").(*jwt.Token))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Permission Denied")
		}

		surveySessions := surveyManager.AllSurveySessions()

		surveySessionList := new(SurveySessionFullObjectList)

		for _, surveySession := range surveySessions {
			// filter by user, or give all if an admin
			if user.IsAdmin || user.Username == surveySession.User {
				var state map[string]interface{}
				json.Unmarshal(surveySession.SurveyState, &state)

				surveySessionResponse := SurveySessionFullObject{
					SessionID:    surveySession.SessionID,
					SurveyType:   surveySession.SurveyType,
					SurveyState:  state,
					IsTest:       surveySession.IsTest,
					ReferralCode: surveySession.ReferralCode,
					Username:     surveySession.User,
					Created:      surveySession.CreatedAt.Format("Mon, 2 Jan 2006 3:04 PM"),
					Modified:     surveySession.UpdatedAt.Format("Mon, 2 Jan 2006 3:04 PM"),
				}
				surveySessionList.SurveySessions = append(surveySessionList.SurveySessions, surveySessionResponse)
			}
		}

		return c.JSON(http.StatusOK, surveySessionList)
	}
}

/*
* SurveySessionData
 */

// NewSurveySessionData create new session data object
func NewSurveySessionData() echo.HandlerFunc {
	return func(c echo.Context) error {
		surveySessionDataRequest := new(SurveySessionDataRequest)

		if err := c.Bind(surveySessionDataRequest); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if surveySessionDataRequest.SessionID == "" {
			log.Printf("surveySessionID: '%s' not found", surveySessionDataRequest.SessionID)
			return echo.NewHTTPError(http.StatusNotFound, "surveySessionID not attached to payload")
		}

		data, err := json.Marshal(surveySessionDataRequest.Data)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		surveySessionData, err := surveyManager.NewSurveySessionData(surveySessionDataRequest.ReferenceID, surveySessionDataRequest.SessionID, surveySessionDataRequest.SurveyName, surveySessionDataRequest.DataKey, data, surveySessionDataRequest.Username, c.RealIP(), c.Request().Header.Get("User-Agent"))

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		surveySessionDataResponse := SurveySessionDataObject{
			SessionDataID: surveySessionData.SessionDataID,
			ReferenceID:   surveySessionData.ReferenceID,
			SessionID:     surveySessionData.SessionID,
			SurveyName:    surveySessionData.SurveyName,
			DataKey:       surveySessionData.DataKey,
			Data:          surveySessionDataRequest.Data, // meh, easier to just reuse than rebuild
			Username:      surveySessionData.User,
		}

		log.Printf("hmmm %s", surveySessionDataResponse)

		return c.JSON(http.StatusOK, surveySessionDataResponse)
	}
}

// UpdateSurveySessionData update session data object
func UpdateSurveySessionData() echo.HandlerFunc {
	return func(c echo.Context) error {
		surveySessionDataID := c.Param("id")

		if surveySessionDataID == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Please specify surveySessionDataID")
		}

		surveySessionDataUpdateRequest := new(SurveySessionDataUpdateRequest)
		if err := c.Bind(surveySessionDataUpdateRequest); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		log.Printf("Survey Session Data: %s", surveySessionDataUpdateRequest.Data)

		if len(surveySessionDataUpdateRequest.Data) == 0 {
			return echo.NewHTTPError(http.StatusInternalServerError, "data not attached to payload")
		}

		surveySessionData := surveyManager.FindSurveySessionData(surveySessionDataID)

		if surveySessionData.SessionDataID == "" {
			log.Printf("surveySessionDataID: '%s' not found", surveySessionDataID)
			return echo.NewHTTPError(http.StatusNotFound, "surveySessionDataID not found")
		}

		data, err := json.Marshal(surveySessionDataUpdateRequest.Data)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		surveySessionData.Data = data

		_, err = surveyManager.UpdateSurveySessionData(surveySessionData)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		surveySessionDataResponse := SurveySessionDataObject{
			SessionDataID: surveySessionData.SessionDataID,
			ReferenceID:   surveySessionData.ReferenceID,
			SessionID:     surveySessionData.SessionID,
			SurveyName:    surveySessionData.SurveyName,
			DataKey:       surveySessionData.DataKey,
			Data:          surveySessionDataUpdateRequest.Data, // meh, easier
			Username:      surveySessionData.User,
		}

		return c.JSON(http.StatusOK, surveySessionDataResponse)
	}
}

// FetchSurveySessionData retrieve particular session data object
func FetchSurveySessionData() echo.HandlerFunc {
	return func(c echo.Context) error {
		surveySessionDataID := c.Param("id")

		if surveySessionDataID == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Please specify surveySessionDataID")
		}

		surveySessionData := surveyManager.FindSurveySessionData(surveySessionDataID)

		if surveySessionData.SessionDataID == "" {
			log.Printf("surveySessionDataID: '%s' not found", surveySessionDataID)
			return echo.NewHTTPError(http.StatusNotFound, "surveySessionDataID not found")
		}

		var data map[string]interface{}
		json.Unmarshal(surveySessionData.Data, &data)

		surveySessionDataResponse := SurveySessionDataObject{
			SessionDataID: surveySessionData.SessionDataID,
			ReferenceID:   surveySessionData.ReferenceID,
			SessionID:     surveySessionData.SessionID,
			SurveyName:    surveySessionData.SurveyName,
			DataKey:       surveySessionData.DataKey,
			Data:          data,
			Username:      surveySessionData.User,
		}

		return c.JSON(http.StatusOK, surveySessionDataResponse)
	}
}

// DeleteSurveySessionData deletion of session data
func DeleteSurveySessionData() echo.HandlerFunc {
	return func(c echo.Context) error {
		surveySessionDataID := c.Param("id")

		if surveySessionDataID == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Please specify surveySessionDataID")
		}

		surveySessionData := surveyManager.FindSurveySessionData(surveySessionDataID)

		if surveySessionData.SessionDataID == "" {
			log.Printf("surveySessionDataID: '%s' not found", surveySessionDataID)
			return echo.NewHTTPError(http.StatusNotFound, "surveySessionDataID not found")
		}

		err := surveyManager.DeleteSurveySessionData(surveySessionData)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, true)
	}
}

// FetchAllSurveySessionData retrieve all session data for a session
func FetchAllSurveySessionData() echo.HandlerFunc {
	return func(c echo.Context) error {
		surveySessionID := c.Param("id")
		allData := c.QueryParam("all")

		if surveySessionID == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Please specify surveySessionID")
		}

		sessionDataObjects := surveyManager.FindAllSurveySessionData(surveySessionID)

		surveySessionDataObjectList := new(SurveySessionDataObjectList)

		specialKeys := map[string]bool{
			"entry":      true,
			"followup":   true,
			"completion": true,
		}

		for _, surveySessionData := range sessionDataObjects {
			if allData == "true" || !specialKeys[surveySessionData.DataKey] {
				var data map[string]interface{}
				json.Unmarshal(surveySessionData.Data, &data)

				surveySessionDataObject := SurveySessionDataObject{
					SessionDataID: surveySessionData.SessionDataID,
					ReferenceID:   surveySessionData.ReferenceID,
					SessionID:     surveySessionData.SessionID,
					SurveyName:    surveySessionData.SurveyName,
					DataKey:       surveySessionData.DataKey,
					Data:          data,
					Username:      surveySessionData.User,
					Created:       surveySessionData.CreatedAt.Format("Mon, 2 Jan 2006 3:04 PM"),
					Modified:      surveySessionData.UpdatedAt.Format("Mon, 2 Jan 2006 3:04 PM"),
				}
				surveySessionDataObjectList.SessionData = append(surveySessionDataObjectList.SessionData, surveySessionDataObject)
			}
		}

		return c.JSON(http.StatusOK, surveySessionDataObjectList)
	}
}

// ExportData export all survey sessions and data
func ExportData() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ValidateAdmin(c.Get("user").(*jwt.Token))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Permission Denied")
		}

		data := surveyManager.Export()

		surveyExportObject := SurveyExportObject{
			Name:   "export",
			Export: data,
		}

		return c.JSON(http.StatusOK, surveyExportObject)
	}
}

// ResetData remove all survey sessions and data
func ResetData() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ValidateAdmin(c.Get("user").(*jwt.Token))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Permission Denied")
		}

		err = surveyManager.ResetData()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, true)
	}
}
