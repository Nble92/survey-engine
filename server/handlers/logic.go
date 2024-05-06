package handlers

import (
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	"log"
	"net/http"

	"github.com/dchote/survey-engine/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

var (
	logicManager *models.LogicManager
)

// LogicFunctionObject struct
type LogicFunctionObject struct {
	FunctionName string `json:"functionName"`
	FunctionJS   string `json:"functionJS"`
}

// LogicFunctionList struct
type LogicFunctionList struct {
	LogicFunctions []*LogicFunctionObject `json:"logicFunctions"`
}

// LogicFunctionUpdateRequest struct
type LogicFunctionUpdateRequest struct {
	FunctionJS string `json:"functionJS"`
}

// InitLogicHandler initialize the logic manager
func InitLogicHandler(db *gorm.DB) {
	logicManager = models.NewLogicManager(db)
}

// ListLogicFunctions return all functions
func ListLogicFunctions() echo.HandlerFunc {
	return func(c echo.Context) error {
		logicFunctions := logicManager.AllLogicFunctions()

		logicFunctionList := new(LogicFunctionList)

		for _, logicFunction := range logicFunctions {
			logicFunctionObject := &LogicFunctionObject{
				FunctionName: logicFunction.FunctionName,
				FunctionJS:   logicFunction.FunctionJS,
			}
			logicFunctionList.LogicFunctions = append(logicFunctionList.LogicFunctions, logicFunctionObject)
		}

		return c.JSON(http.StatusOK, logicFunctionList)
	}
}

// FetchLogicFunction retrieve single function
func FetchLogicFunction() echo.HandlerFunc {
	return func(c echo.Context) error {
		functionName := c.Param("functionName")

		if functionName == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Please specify function name")
		}

		logicFunction := logicManager.FindLogicFunction(functionName)

		if logicFunction.FunctionName == "" {
			log.Printf("Function '%s' not found", functionName)
			return echo.NewHTTPError(http.StatusNotFound, "Function not found")
		}

		logicFunctionObject := &LogicFunctionObject{
			FunctionName: logicFunction.FunctionName,
			FunctionJS:   logicFunction.FunctionJS,
		}

		return c.JSON(http.StatusOK, logicFunctionObject)
	}
}

// UpdateLogicFunction update function
func UpdateLogicFunction() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ValidateAdmin(c.Get("user").(*jwt.Token))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Permission Denied")
		}

		logicFunctionName := c.Param("functionName")

		logicFunctionUpdateRequest := new(LogicFunctionUpdateRequest)
		if err = c.Bind(logicFunctionUpdateRequest); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if logicFunctionUpdateRequest.FunctionJS == "" {
			return echo.NewHTTPError(http.StatusInternalServerError, "FunctionJS not attached to payload")
		}

		logicFunction := logicManager.FindLogicFunction(logicFunctionName)

		if logicFunction.FunctionName == "" {
			_, err = logicManager.NewLogicFunction(logicFunctionName, logicFunctionUpdateRequest.FunctionJS)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		} else {
			logicFunction.FunctionJS = logicFunctionUpdateRequest.FunctionJS
			logicManager.UpdateLogicFunction(logicFunction)
		}

		return c.JSON(http.StatusOK, false)
	}
}

// DeleteLogicFunction delete function
func DeleteLogicFunction() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ValidateAdmin(c.Get("user").(*jwt.Token))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Permission Denied")
		}

		logicFunctionName := c.Param("functionName")

		logicFunction := logicManager.FindLogicFunction(logicFunctionName)

		if logicFunction.FunctionName == "" {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		logicManager.DeleteLogicFunction(logicFunction)

		return c.JSON(http.StatusOK, true)
	}
}

// BuildFunctionsJS builds a javascript file for including in the HTML HEAD
func BuildFunctionsJS() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJavaScript)
		c.Response().WriteHeader(http.StatusOK)

		logicFunctions := logicManager.AllLogicFunctions()

		for _, logicFunction := range logicFunctions {
			functionJS := "function SurveyEngine_" + logicFunction.FunctionName + "(surveySessionData) { \n" + logicFunction.FunctionJS + "\n}\n\n"
			c.Response().Write([]byte(functionJS))
		}

		c.Response().Flush()

		return nil
	}
}
