package models

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	// the blank import is the way GORM wants you to do it
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	uuid "github.com/satori/go.uuid"
)

// SurveySettings struct
type SurveySettings struct {
	gorm.Model `json:"-"`
	Type       string `gorm:"not null" json:"type"` // linear, calendar
	Config     []byte `sql:"type:blob" json:"config"`
}

// Survey struct
type Survey struct {
	gorm.Model `json:"-"`
	Name       string `gorm:"unique" json:"name"`
	Author     string `gorm:"not null;" json:"author"`
	SurveyJSON []byte `sql:"type:blob" json:"-"`
}

// SurveySession struct
type SurveySession struct {
	gorm.Model   `json:"-"`
	SessionID    string `gorm:"unique" json:"sessionID"`
	SurveyType   string `gorm:"not null" json:"surveyType"`
	SurveyState  []byte `sql:"type:blob" json:"-"`
	IsTest       bool   `gorm:"not null" sql:"default:false" json:"isTest"`
	ReferralCode string `gorm:"not null" sql:"default:''" json:"referralCode"`
	User         string `gorm:"not null" json:"user"`
	IPAddress    string `gorm:"not null" json:"ipAddress"`
	UserAgent    string `gorm:"not null" json:"userAgent"`
}

// SurveySessionData struct
type SurveySessionData struct {
	gorm.Model    `json:"-"`
	SessionDataID string `gorm:"unique" json:"sessionDataID"`
	ReferenceID   string `json:"referenceID"`
	SessionID     string `gorm:"not null" json:"sessionID"`
	SurveyName    string `gorm:"not null" json:"surveyName"`
	DataKey       string `gorm:"not null" json:"dataKey"`
	Data          []byte `sql:"type:blob" json:"data"`
	User          string `gorm:"not null" json:"user"`
	IPAddress     string `gorm:"not null" json:"ipAddress"`
	UserAgent     string `gorm:"not null" json:"userAgent"`
}

// SurveyManager manager
type SurveyManager struct {
	db *gorm.DB
}

// NewSurveyManager for handling the DB operations
func NewSurveyManager(db *gorm.DB) *SurveyManager {
	surveyManager := SurveyManager{}
	surveyManager.db = db

	surveyManager.db.AutoMigrate(&Survey{})
	surveyManager.db.AutoMigrate(&SurveySettings{})
	surveyManager.db.AutoMigrate(&SurveySession{})
	surveyManager.db.AutoMigrate(&SurveySessionData{})

	return &surveyManager
}

// SurveySettings object, singular
func (surveyManager *SurveyManager) SurveySettings() (surveySettings *SurveySettings, err error) {
	settings := &SurveySettings{}

	err = surveyManager.db.First(&settings).Error
	if err != nil {
		// lets generate a new blank settings
		settings.Type = "linear"

		err = surveyManager.db.Save(&settings).Error
		if err != nil {
			return nil, err
		}
	}

	return settings, nil
}

// UpdateSurveySettings update survey settings
func (surveyManager *SurveyManager) UpdateSurveySettings(surveyType string, surveyConfig []byte) (surveySettings *SurveySettings, err error) {
	settings := &SurveySettings{}
	surveyManager.db.First(&settings)

	settings.Type = surveyType
	settings.Config = surveyConfig

	err = surveyManager.db.Save(&settings).Error
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// FindSurvey simple search
func (surveyManager *SurveyManager) FindSurvey(name string) *Survey {
	survey := Survey{}

	surveyManager.db.Where("name=?", name).Find(&survey)

	return &survey
}

// AllSurveys return all surveys
func (surveyManager *SurveyManager) AllSurveys() (surveys []Survey) {
	err := surveyManager.db.Find(&surveys).Error

	if err != nil {
		return surveys
	}

	return surveys
}

// NewSurvey create new survey
func (surveyManager *SurveyManager) NewSurvey(name string, author string, surveyJSON []byte) (survey *Survey, err error) {
	survey = &Survey{
		Name:       name,
		Author:     author,
		SurveyJSON: surveyJSON,
	}

	err = surveyManager.db.Create(&survey).Error
	if err != nil {
		return nil, err
	}

	return survey, nil
}

// UpdateSurvey update existing survey
func (surveyManager *SurveyManager) UpdateSurvey(updateSurvey *Survey) (survey *Survey, err error) {
	err = surveyManager.db.Save(&updateSurvey).Error
	if err != nil {
		return nil, err
	}

	return updateSurvey, nil
}

// DeleteSurvey delete a survey
func (surveyManager *SurveyManager) DeleteSurvey(updateSurvey *Survey) (err error) {
	err = surveyManager.db.Delete(&updateSurvey).Error
	if err != nil {
		return err
	}

	return nil
}

// NewSurveySession create a new session
func (surveyManager *SurveyManager) NewSurveySession(surveyType string, isTest bool, referralCode string, username string, ipAddress string, userAgent string) (surveySession *SurveySession, err error) {
	sessionID := uuid.NewV4()

	surveySession = &SurveySession{
		SessionID:    sessionID.String(),
		SurveyType:   surveyType,
		SurveyState:  []byte("{\"completed\": false, \"dataRefs\": [], \"entrySurveyComplete\": false}"), // NOTE: not entirely sure this is a good idea yet
		IsTest:       isTest,
		ReferralCode: referralCode,
		User:         username,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	}

	err = surveyManager.db.Create(&surveySession).Error
	if err != nil {
		return nil, err
	}

	return surveySession, nil
}

// AllSurveySessions return all stored sessions
func (surveyManager *SurveyManager) AllSurveySessions() (surveySessions []SurveySession) {
	err := surveyManager.db.Find(&surveySessions).Error

	if err != nil {
		return surveySessions
	}

	return surveySessions
}

// FindSurveySession find a particular session
func (surveyManager *SurveyManager) FindSurveySession(id string) *SurveySession {
	surveySession := SurveySession{}

	surveyManager.db.Where("session_id=?", id).Find(&surveySession)

	return &surveySession
}

// UpdateSurveySession update existing session
func (surveyManager *SurveyManager) UpdateSurveySession(updateSurveySession *SurveySession) (surveySession *SurveySession, err error) {
	err = surveyManager.db.Save(&updateSurveySession).Error
	if err != nil {
		return nil, err
	}

	return updateSurveySession, nil
}

// NewSurveySessionData create a new session data object
func (surveyManager *SurveyManager) NewSurveySessionData(referenceID string, sessionID string, surveyName string, dataKey string, data []byte, username string, ipAddress string, userAgent string) (surveySession *SurveySessionData, err error) {
	sessionDataID := uuid.NewV4()

	surveySessionData := &SurveySessionData{
		SessionDataID: sessionDataID.String(),
		ReferenceID:   referenceID,
		SessionID:     sessionID,
		SurveyName:    surveyName,
		DataKey:       dataKey,
		Data:          data,
		User:          username,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
	}

	err = surveyManager.db.Create(&surveySessionData).Error
	if err != nil {
		return nil, err
	}

	return surveySessionData, nil
}

// FindSurveySessionData simple search for existing session data object
func (surveyManager *SurveyManager) FindSurveySessionData(sessionDataID string) *SurveySessionData {
	surveySessionData := SurveySessionData{}

	surveyManager.db.Where("session_data_id=?", sessionDataID).Find(&surveySessionData)

	return &surveySessionData
}

// FindAllSurveySessionData find all session data for the session
func (surveyManager *SurveyManager) FindAllSurveySessionData(sessionID string) (sessionData []SurveySessionData) {
	err := surveyManager.db.Where("session_id=?", sessionID).Find(&sessionData).Error

	if err != nil {
		return sessionData
	}

	return sessionData
}

// UpdateSurveySessionData update session data
func (surveyManager *SurveyManager) UpdateSurveySessionData(updateSurveySessionData *SurveySessionData) (surveySessionData *SurveySessionData, err error) {
	err = surveyManager.db.Save(&updateSurveySessionData).Error
	if err != nil {
		return nil, err
	}

	return updateSurveySessionData, nil
}

// DeleteSurveySessionData delete session data
func (surveyManager *SurveyManager) DeleteSurveySessionData(updateSurveySessionData *SurveySessionData) (err error) {
	err = surveyManager.db.Delete(&updateSurveySessionData).Error
	if err != nil {
		return err
	}

	return nil
}

// Export is used to export ALL survey sessions and data
func (surveyManager *SurveyManager) Export() []interface{} {
	var surveyExport []interface{}
	var surveySessions []SurveySession

	err := surveyManager.db.Order("created_at asc").Find(&surveySessions).Error

	if err != nil {
		fmt.Printf("error: %s", err)
		return nil
	}

	for _, surveySession := range surveySessions {
		var surveySessionData []SurveySessionData
		var surveySessionDataMap []interface{}

		err := surveyManager.db.Where("session_id=?", surveySession.SessionID).Order("created_at asc").Find(&surveySessionData).Error
		if err == nil {
			for _, surveySessionData := range surveySessionData {
				var surveySessionDataData map[string]interface{}
				json.Unmarshal(surveySessionData.Data, &surveySessionDataData)

				surveySessionDataMap = append(surveySessionDataMap, map[string]interface{}{
					"sessionID":     surveySessionData.SessionID,
					"sessionDataID": surveySessionData.SessionDataID,
					"referenceID":   surveySessionData.ReferenceID,
					"surveyName":    surveySessionData.SurveyName,
					"dataKey":       surveySessionData.DataKey,
					"data":          surveySessionDataData,
					"createdAt":     surveySessionData.CreatedAt,
					"updatedAt":     surveySessionData.UpdatedAt,
				})
			}
		}

		var surveyState map[string]interface{}
		json.Unmarshal(surveySession.SurveyState, &surveyState)

		sessionObject := map[string]interface{}{
			"sessionID":    surveySession.SessionID,
			"surveyType":   surveySession.SurveyType,
			"surveyState":  surveyState,
			"username":     surveySession.User,
			"isTest":       surveySession.IsTest,
			"referralCode": surveySession.ReferralCode,
			"sessionData":  surveySessionDataMap,
			"createdAt":    surveySession.CreatedAt,
			"updatedAt":    surveySession.UpdatedAt,
		}

		surveyExport = append(surveyExport, sessionObject)
	}

	return surveyExport
}

// ResetData removes all survey sessions and data
func (surveyManager *SurveyManager) ResetData() (err error) {
	surveySessionData := SurveySessionData{}
	err = surveyManager.db.Delete(&surveySessionData).Error
	if err != nil {
		return err
	}

	surveySession := SurveySession{}
	err = surveyManager.db.Delete(&surveySession).Error
	if err != nil {
		return err
	}

	return nil
}
