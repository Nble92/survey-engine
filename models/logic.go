package models

import (
	//"fmt"

	"github.com/jinzhu/gorm"
	// the blank import is the way GORM wants you to do it
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// LogicFunction struct
type LogicFunction struct {
	gorm.Model   `json:"-"`
	FunctionName string `gorm:"not null" json:"functionName"`
	FunctionJS   string `json:"functionJS"`
}

// LogicManager manager
type LogicManager struct {
	db *gorm.DB
}

// NewLogicManager for handling the DB operations
func NewLogicManager(db *gorm.DB) *LogicManager {
	logicManager := LogicManager{}
	logicManager.db = db

	logicManager.db.AutoMigrate(&LogicFunction{})

	return &logicManager
}

// LogicFunction

// FindLogicFunction simple search
func (logicManager *LogicManager) FindLogicFunction(name string) *LogicFunction {
	logicFunction := LogicFunction{}

	logicManager.db.Where("function_name=?", name).Find(&logicFunction)

	return &logicFunction
}

// AllLogicFunctions return all stored functions
func (logicManager *LogicManager) AllLogicFunctions() (logicFunctions []LogicFunction) {
	err := logicManager.db.Find(&logicFunctions).Error

	if err != nil {
		return logicFunctions
	}

	return logicFunctions
}

// NewLogicFunction create new function
func (logicManager *LogicManager) NewLogicFunction(functionName string, functionJS string) (logicFunction *LogicFunction, err error) {
	logicFunction = &LogicFunction{
		FunctionName: functionName,
		FunctionJS:   functionJS,
	}

	err = logicManager.db.Create(&logicFunction).Error
	if err != nil {
		return nil, err
	}

	return logicFunction, nil
}

// UpdateLogicFunction update existing function
func (logicManager *LogicManager) UpdateLogicFunction(updateLogicFunction *LogicFunction) (logicFunction *LogicFunction, err error) {
	err = logicManager.db.Save(&updateLogicFunction).Error
	if err != nil {
		return nil, err
	}

	return updateLogicFunction, nil
}

// DeleteLogicFunction delete a function
func (logicManager *LogicManager) DeleteLogicFunction(updateLogicFunction *LogicFunction) (err error) {
	err = logicManager.db.Delete(&updateLogicFunction).Error
	if err != nil {
		return err
	}

	return nil
}
