package models

import (
	//"fmt"

	"github.com/jinzhu/gorm"
	// the blank import is the way GORM wants you to do it
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Instruction struct
type Instruction struct {
	gorm.Model `json:"-"`
	SurveyName string `gorm:"not null" json:"surveyName"`
	Data       []byte `sql:"type:blob" json:"-"`
}

// InstructionManager manager
type InstructionManager struct {
	db *gorm.DB
}

// NewInstructionManager for handling the DB operations
func NewInstructionManager(db *gorm.DB) *InstructionManager {
	instructionManager := InstructionManager{}
	instructionManager.db = db

	instructionManager.db.AutoMigrate(&Instruction{})

	return &instructionManager
}

// Instruction

// FindInstruction simple search
func (instructionManager *InstructionManager) FindInstruction(name string) *Instruction {
	instruction := Instruction{}

	instructionManager.db.Where("survey_name=?", name).Find(&instruction)

	return &instruction
}

// AllInstructions return all stored functions
func (instructionManager *InstructionManager) AllInstructions() (instructions []Instruction) {
	err := instructionManager.db.Find(&instructions).Error

	if err != nil {
		return instructions
	}

	return instructions
}

// NewInstruction create new function
func (instructionManager *InstructionManager) NewInstruction(surveyName string, data []byte) (instruction *Instruction, err error) {
	instruction = &Instruction{
		SurveyName: surveyName,
		Data:       data,
	}

	err = instructionManager.db.Create(&instruction).Error
	if err != nil {
		return nil, err
	}

	return instruction, nil
}

// UpdateInstruction update existing function
func (instructionManager *InstructionManager) UpdateInstruction(updateInstruction *Instruction) (instruction *Instruction, err error) {
	err = instructionManager.db.Save(&updateInstruction).Error
	if err != nil {
		return nil, err
	}

	return updateInstruction, nil
}

// DeleteInstruction delete a function
func (instructionManager *InstructionManager) DeleteInstruction(updateInstruction *Instruction) (err error) {
	err = instructionManager.db.Delete(&updateInstruction).Error
	if err != nil {
		return err
	}

	return nil
}
