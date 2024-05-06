export function parseAnswers(surveyJSON, surveyData) {
  if (!surveyJSON || !surveyData) {
    // eslint-disable-next-line
    console.log('parseAnswers() requires both surveyJSON and surveyData to be passed')
    return false
  }
  
  if (surveyData.data) {
    // reset the answers
    var answers = []

    var basicAnswers = ['text', 'boolean', 'rating']
    var dynamicAnswers = ['matrixdynamic']
    var multiChoiceAnswers = ['dropdown', 'checkbox', 'radiogroup']
    var matrixAnswers = ['matrix', 'matrixdropdown']
    var ignoreAnswers = ['html']
    
    var elementChoiceID = false
    var thisAnswer = {}
    var foundElement = false
    var largestDynamicMatrix = 0

    var elementColumn, columnID, elementRow, rowID, answerValue, answerValueID, answerSubID, subAnswer, subAnswers, subAnswerName, subAnswerTitle, subAnswerValue, subAnswerText, subAnswerColumn


    for (var pageID in surveyJSON.pages) {
      var page = surveyJSON.pages[pageID]
  
      for (var elementID in page.elements) {
        var element = page.elements[elementID]
        
        // iterate over all the answers to find a match
        for (var answerKey in surveyData.data) {
          if (element.name == answerKey) {
            var answer = surveyData.data[answerKey]
            
            // basic
            if (basicAnswers.includes(element.type)) {
              var answerText = answer
              if (answer === true) {
                answerText = element.labelTrue
                answer = '1'
              } else if (answer === false) {
                answerText = element.labelFalse
                answer = '0'
              }
              
              answers.push({'type': 'single', 'name': element.name, 'title': titleOrName(element), 'text': answerText, 'value': answer})
            // matrix dynamic
            } else if (dynamicAnswers.includes(element.type)) {
              thisAnswer = {'type': 'multiple', 'name': element.name, 'title': titleOrName(element), 'values': []}
              
              for (answerSubID in answer) {
                // NOTE: fucking parseInt because 10 & 11 are not numbers?
                if (parseInt(answerSubID) > largestDynamicMatrix) {
                  largestDynamicMatrix = parseInt(answerSubID)
                }
                
                subAnswer = answer[answerSubID]
                subAnswers = []
            
                for (subAnswerColumn in subAnswer) {
                  for (columnID in element.columns) {
                
                    if (element.columns[columnID].name == subAnswerColumn) {
                      subAnswerValue = subAnswer[subAnswerColumn]
                      elementColumn = element.columns[columnID]
                      
                      // need to add row number here...
                      subAnswerName = answerSubID + '.' + elementColumn.name
                      subAnswerTitle = (elementColumn.title) ? elementColumn.title : elementColumn.name
                      
                      if (elementColumn.cellType) {
                        if (basicAnswers.includes(elementColumn.cellType)) {
                          subAnswerText = subAnswerValue
                          if (subAnswerValue === true) {
                            subAnswerValue = '1'
                            subAnswerText = 'Yes'
                          } else if (subAnswerValue === false) {
                            subAnswerValue = '0'
                            subAnswerText = 'No'
                          }
                          
                          subAnswers.push({'name': subAnswerName, 'title': subAnswerTitle, 'text': subAnswerText, 'value': subAnswerValue})
                        } else if (multiChoiceAnswers.includes(elementColumn.cellType)) {
                          for (var elementColumnChoiceID in elementColumn.choices) {
                            if (elementColumn.choices[elementColumnChoiceID].value == subAnswerValue || elementColumnChoiceID == subAnswerValue) {
                              subAnswers.push({'name': subAnswerName, 'title': subAnswerTitle, 'text': elementColumn.choices[elementColumnChoiceID].text, 'value': subAnswerValue})
                            }
                          }
                        }
                      } else {
                        // eslint-disable-next-line
                        console.log('NO CELL TYPE FOR SUB ANSWER COLUMN', element, answer)
                      }
                  
                    }
                  }
                }
                thisAnswer.values.push(subAnswers)
              }
              answers.push(thisAnswer)
            // multi choice
            } else if (multiChoiceAnswers.includes(element.type)) {
              thisAnswer = {'type': 'single', 'name': element.name, 'title': titleOrName(element), 'values': []}
              
              if (Array.isArray(answer)) {
                for (answerValueID in answer) {
                  answerValue = answer[answerValueID]
            
                  foundElement = false
                  for (elementChoiceID in element.choices) {
                    if (element.choices[elementChoiceID] == answerValue || element.choices[elementChoiceID].value == answerValue) {
                      foundElement = true
                      thisAnswer.values.push({'text': (element.choices[elementChoiceID].text) ? element.choices[elementChoiceID].text : element.choices[elementChoiceID], 'value': answerValue})
                    }
                  }
                  
                  if (answerValue == 'other' && surveyData.data.hasOwnProperty(answerKey + '-Comment')) {
                    answerValue = surveyData.data[answerKey + '-Comment']
                  }
              
                  if (!foundElement) {
                    thisAnswer.values.push({'text': answerValue, 'value': answerValue})
                  }
                }
              } else {
                foundElement = false
                for (elementChoiceID in element.choices) {
                  if (element.choices[elementChoiceID] == answer || element.choices[elementChoiceID].value == answer) {
                    foundElement = true
                    thisAnswer.values.push({'text': (element.choices[elementChoiceID].text) ? element.choices[elementChoiceID].text : element.choices[elementChoiceID], 'value': answer})
                  }
                }
                
                if (answer == 'other' && surveyData.data.hasOwnProperty(answerKey + '-Comment')) {
                  answer = surveyData.data[answerKey + '-Comment']
                }
                
                if (!foundElement) {
                  thisAnswer.values.push({'text': answer, 'value': answer})
                }
              }
          
              if (thisAnswer.values.length == 1) {
                thisAnswer.value = thisAnswer.values[0].value
                thisAnswer.text = thisAnswer.values[0].text
              } else {
                thisAnswer.value = ''
                thisAnswer.text = ''
                for (var value of thisAnswer.values) {
                  thisAnswer.value += value.value + ' '
                  thisAnswer.text += value.text + ' '
                }
              }
              
              answers.push(thisAnswer)
            // matrix
            } else if (matrixAnswers.includes(element.type)) {
              thisAnswer = {'type': 'multiple', 'name': element.name, 'title': titleOrName(element), 'values': []}
                              
              for (answerSubID in answer) {
                subAnswer = answer[answerSubID]
                subAnswers = []
            
                if (typeof(subAnswer) == 'object') {
                  for (rowID in element.rows) {
                    if (element.rows[rowID] == answerSubID || element.rows[rowID].value == answerSubID) {
                      elementRow = element.rows[rowID]
                  
                      for (subAnswerColumn in subAnswer) {
                        for (columnID in element.columns) {
                          if (element.columns[columnID].name == subAnswerColumn) {                            
                            subAnswerValue = subAnswer[subAnswerColumn]
                            elementColumn = element.columns[columnID]
                        
                            subAnswerName = (typeof(elementRow) == 'object') ? elementRow.value + '.' + elementColumn.name : elementRow + '.' + elementColumn.name
                            subAnswerTitle = (typeof(elementRow) == 'object') ? elementRow.text + ' ' + titleOrName(elementColumn) : elementRow + ' ' + titleOrName(elementColumn)
                            
                            if (elementColumn.choices) {
                              for (elementChoiceID in elementColumn.choices) {
                                if (elementColumn.choices[elementChoiceID] == subAnswerValue || elementColumn.choices[elementChoiceID].value == subAnswerValue) {
                                  subAnswers.push({'name': subAnswerName, 'title': subAnswerTitle, 'text': (elementColumn.choices[elementChoiceID].text) ? elementColumn.choices[elementChoiceID].text : elementColumn.choices[elementChoiceID], 'value': subAnswerValue})
                                  break
                                }
                              }
                            } else {
                              for (elementChoiceID in element.choices) {
                                if (element.choices[elementChoiceID] == subAnswerValue || element.choices[elementChoiceID].value == subAnswerValue) {
                                  subAnswers.push({'name': subAnswerName, 'title': subAnswerTitle, 'text': (element.choices[elementChoiceID].value) ? element.choices[elementChoiceID].value : element.choices[elementChoiceID], 'value': subAnswerValue})
                                  break
                                }
                              }
                            }
                          }
                        }
                      }
                    }
                  }
              
                  thisAnswer.values.push(subAnswers)
                } else {
                  for (rowID in element.rows) {
                    if (element.rows[rowID].value == answerSubID) {
                      elementRow = element.rows[rowID]
                  
                      for (columnID in element.columns) {
                        if (element.columns[columnID] == subAnswer || element.columns[columnID].value == subAnswer) {
                          elementColumn = element.columns[columnID]
                          subAnswers.push({ 'name': elementRow.value, 'title': elementRow.text, 'text': (elementColumn.text) ? elementColumn.text : elementColumn, 'value': subAnswer})
                          break
                        }
                      }
                    }
                  }
              
                  thisAnswer.values.push(subAnswers)
                }
              }
          
              answers.push(thisAnswer)
            } else if (ignoreAnswers.includes(element.type)) {
              // do nothing... we're ignoring these
            } else {
              // eslint-disable-next-line
              console.log('UNHANDLED ELEMENT TYPE', element.type, element, answer)
            } /// conditional
        
          }
        }
      }
    }
    
    return [answers, largestDynamicMatrix]
  }
}

export function titleOrName(element) {
  return (element.title) ? element.title : element.name
}