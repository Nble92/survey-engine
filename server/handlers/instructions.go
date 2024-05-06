package handlers

import (
	//"fmt"
	//"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/dchote/survey-engine/models"

	"github.com/GeertJohan/go.rice"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

var (
	instructionManager *models.InstructionManager
	assetsBox          *rice.Box
)

// InstructionObject struct
type InstructionObject struct {
	SurveyName string                 `json:"surveyName"`
	Data       map[string]interface{} `json:"data"`
}

// InstructionList struct
type InstructionList struct {
	Instructions []*InstructionObject `json:"instructions"`
}

// InstructionUpdateRequest struct
type InstructionUpdateRequest struct {
	Data map[string]interface{} `json:"data"`
}

// InstructionPageObject object
type InstructionPageObject struct {
	PageName    string `json:"pageName"`
	PageContent string `json:"pageContent"`
}

// InstructionDataObject object
type InstructionDataObject struct {
	SummaryInstruction string                   `json:"summaryInstruction"`
	InstructionPages   []*InstructionPageObject `json:"instructionPages"`
}

// InitInstructionHandler initialize the logic manager
func InitInstructionHandler(db *gorm.DB, assets *rice.Box) {
	instructionManager = models.NewInstructionManager(db)
	assetsBox = assets
}

// ListInstructions return all functions
func ListInstructions() echo.HandlerFunc {
	return func(c echo.Context) error {
		instructions := instructionManager.AllInstructions()

		instructionList := new(InstructionList)

		for _, instruction := range instructions {
			var data map[string]interface{}
			json.Unmarshal(instruction.Data, &data)

			instructionObject := &InstructionObject{
				SurveyName: instruction.SurveyName,
				Data:       data,
			}
			instructionList.Instructions = append(instructionList.Instructions, instructionObject)
		}

		return c.JSON(http.StatusOK, instructionList)
	}
}

// FetchInstruction retrieve single function
func FetchInstruction() echo.HandlerFunc {
	return func(c echo.Context) error {
		instructionName := c.Param("name")

		if instructionName == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Please specify instruction name")
		}

		instruction := instructionManager.FindInstruction(instructionName)

		if instruction.SurveyName == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Instruction not found")
		}

		var data map[string]interface{}
		json.Unmarshal(instruction.Data, &data)

		instructionObject := &InstructionObject{
			SurveyName: instruction.SurveyName,
			Data:       data,
		}

		return c.JSON(http.StatusOK, instructionObject)
	}
}

// UpdateInstruction update function
func UpdateInstruction() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ValidateAdmin(c.Get("user").(*jwt.Token))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Permission Denied")
		}

		instructionName := c.Param("name")

		instructionUpdateRequest := new(InstructionUpdateRequest)
		if err = c.Bind(instructionUpdateRequest); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		data, err := json.Marshal(instructionUpdateRequest.Data)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		instruction := instructionManager.FindInstruction(instructionName)

		if instruction.SurveyName == "" {
			_, err = instructionManager.NewInstruction(instructionName, data)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		} else {
			instruction.Data = data
			instructionManager.UpdateInstruction(instruction)
		}

		return c.JSON(http.StatusOK, false)
	}
}

// DeleteInstruction delete function
func DeleteInstruction() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ValidateAdmin(c.Get("user").(*jwt.Token))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Permission Denied")
		}

		instructionName := c.Param("name")

		instruction := instructionManager.FindInstruction(instructionName)

		if instruction.SurveyName == "" {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		instructionManager.DeleteInstruction(instruction)

		return c.JSON(http.StatusOK, true)
	}
}

// BuildInstructionHTML builds a javascript file for including in the HTML HEAD
func BuildInstructionHTML() echo.HandlerFunc {
	return func(c echo.Context) error {
		instructionName := c.Param("name")
		pageName := c.Param("page")

		if instructionName == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Please specify instruction name")
		}

		instruction := instructionManager.FindInstruction(instructionName)

		if instruction.SurveyName == "" {
			log.Printf("Instruction '%s' not found", instructionName)
			return echo.NewHTTPError(http.StatusNotFound, "Instruction not found")
		}

		instructionHTML, err := assetsBox.String("instruction.html")
		if err != nil {
			log.Printf("Instruction template not found")
			return echo.NewHTTPError(http.StatusNotFound, "Instruction template not found")
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		c.Response().WriteHeader(http.StatusOK)

		var instructionDataObject InstructionDataObject
		json.Unmarshal(instruction.Data, &instructionDataObject)

		// i have to find a better way to do this...
		if len(instructionDataObject.InstructionPages) > 0 {
			linksHTML := ""
			contentHTML := ""

			for _, page := range instructionDataObject.InstructionPages {
				linksHTML += "<a href='/instructionHTML/" + instructionName + "/" + page.PageName + "' class='mui-btn mui-btn--small mui-btn--primary mui-btn--flat'>" + page.PageName + "</a> "

				if pageName == "" {
					pageName = page.PageName
				}

				if contentHTML == "" && pageName == page.PageName {
					contentHTML = "<hr>\n\n" + page.PageContent
				}
			}

			// {{.InstructionHTML}}

			instructionHTML = strings.Replace(instructionHTML, "{{.InstructionHTML}}", linksHTML+contentHTML, 1)
			c.Response().Write([]byte(instructionHTML))
		} else {
			c.Response().Write([]byte("instructionPages property has no content"))
		}

		c.Response().Flush()

		return nil
	}
}
