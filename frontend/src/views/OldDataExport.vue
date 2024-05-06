<template>
  <v-col cols="8" offset="2">
    <v-row no-gutters>
      <div class="pb-4" v-if="error" style="width: 100%">
        <v-alert
          v-model="error"
          type="error"
          transition="scale-transition"
        >{{ errorText }}</v-alert>
      </div>
      
      <v-card style="width: 100%">
        <v-toolbar class="blue-grey white--text">
          <v-toolbar-title class="headline font-weight-regular"><v-icon dark class="mr-2">import_export</v-icon> CSV Download</v-toolbar-title>
        </v-toolbar>
        <v-card-text>
          <form>
            <v-checkbox
              label="Prefix data column names with survey name"
              v-model="usePrefix"
            ></v-checkbox>
            <v-checkbox
              label="Sort column names alphabetically"
              v-model="sortAlpha"
            ></v-checkbox>
          </form>
        </v-card-text>
        <div class="v-data-table elevation-1 theme--light">
          <div class="v-data-table__wrapper">
            <table>
              <tbody>
                <tr>
                  <td class="font-weight-medium">Non-Event Survey Data</td>
                  <td >{{ nonEventRows }} rows</td>
                  <td class="layout justify-end pa-1 pr-5">
                    <v-btn icon inverted id="nonEventDownload" href="#">
                      <v-icon>save_alt</v-icon>
                    </v-btn>
                  </td>
                </tr>
                <tr>
                  <td class="font-weight-medium">Event Survey Data</td>
                  <td>{{ eventRows }} rows</td>
                  <td class="layout justify-end pa-1 pr-5">
                    <v-btn icon inverted id="eventDownload" href="#">
                      <v-icon>save_alt</v-icon>
                    </v-btn>
                  </td>
                </tr>
                <tr>
                  <td class="font-weight-medium">Event Follow-Up Survey Data</td>
                  <td>{{ followUpRows }} rows</td>
                  <td class="layout justify-end pa-1 pr-5">
                    <v-btn icon inverted id="followUpDownload" href="#">
                      <v-icon>save_alt</v-icon>
                    </v-btn>
                  </td>
                </tr>
                <tr>
                  <td class="font-weight-medium">Note Data</td>
                  <td>{{ notesRows }} rows</td>
                  <td class="layout justify-end pa-1 pr-5">
                    <v-btn icon inverted id="notesDownload" href="#">
                      <v-icon>save_alt</v-icon>
                    </v-btn>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </v-card>
    </v-row>
  </v-col>
</template>


<script>
  import { mapState } from 'vuex'
  
  export default {
    data () {
      return {
        valid: false,
        
        error: false,
        errorText: 'an unknown error occurred',
        
        usePrefix: false,
        sortAlpha: false,
        data: '',
        
        nonEventCSV: '',
        nonEventRows: 0,
        
        eventsCSV: '',
        eventRows: 0,
        
        followUpCSV: '',
        followUpRows: 0,
        
        notesCSV: '',
        notesRows: 0,
      }
    },
    computed: {
      ...mapState('surveySettings', ['surveyType', 'config', 'surveyTypes']),
      ...mapState('surveys', ['surveys']),
    },
    watch: {
      usePrefix() {
        if (this.data) {
          this.buildCSV()
        }
      },
      sortAlpha() {
        if (this.data) {
          this.buildCSV()
        }
      },
    },
    mounted() {
      var that = this
      this.axios.get('/survey/export')
      .then(response => { 
        that.data = response.data
        that.buildCSV()
      }).catch(function (error) {
        return 'error: ' + error
      })
    },
    methods: {
      targetArray(surveyName) {
        if (this.surveyType == 'calendar') {
          if (surveyName == this.config.eventSurvey) {
            return 'event'
          } else if (surveyName == this.config.followUpSurvey) {
            return 'followup'
          }
        }
        return 'nonevent'
      },
      
      pad(n) {
          return n<10 ? '0'+n : n;
      },
      
      buildColumnName(surveyName, columnKey) {
        if (this.usePrefix) {
          return surveyName.replace(/ /g, '_')  + '_' + columnKey.replace(/ /g, '_')
        } else {
          return columnKey.replace(/ /g, '_')
        }
      },
      
      objectsToCSV(columns, rows) {
        var columnDelimiter =  ',';
        var lineDelimiter = '\n';
        
        
        var csv = ''
        csv += columns.join(columnDelimiter)
        csv += lineDelimiter
                
        rows.forEach(function(item) {
          var index = 0
                    
          columns.forEach(function(key) {
            if (index > 0) {
              csv += columnDelimiter
            }
            
            var columnValue = (item[key]) ? item[key] : ''
            
            csv += '"' + columnValue + '"'

            index++
          })
          csv += lineDelimiter
        })
        
        return csv
      },
      
      buildCSV() {
        var eventColumns = {}
        var nonEventColumns = {}
        var followUpColumns = {}
        
        //
        // build a map of all the possible columns
        //
        for (var surveyID in this.data.export) {
          var survey = this.data.export[surveyID]
          
          for (var sessionDataID in survey.sessionData) {
            var sessionData = survey.sessionData[sessionDataID]
            
            for (var questionKey in sessionData.data) {
              var question = sessionData.data[questionKey]
              
              // the question is a dynamic table of rows
              if (typeof(question) == 'object') {
                if (Array.isArray(question)) {
                  for (var index = 0; index <= question.length; index++) {
                    for (var answerColumn in question[index]) {
                      if (sessionData.surveyName  && questionKey && answerColumn) {
                        var objectColummnName = this.buildColumnName(sessionData.surveyName, questionKey + '_row_' + index + '_' + answerColumn)
                      
                        if (this.targetArray(sessionData.surveyName) == 'event') {
                          eventColumns[objectColummnName] = ''
                        } else if (this.targetArray(sessionData.surveyName) == 'followup') {
                          followUpColumns[objectColummnName] = ''
                        } else {
                          nonEventColumns[objectColummnName] = ''
                        }
                      }
                    }
                  }
                } else {
                  //
                  // fucking matrix
                  //
                  //console.log(questionKey, question)
                  
                  for (var matrixItemName in question) {
                    var matrixItem = question[matrixItemName]
                    
                    if (typeof(matrixItem) == 'object') {
                      for (var matrixAnswerName in matrixItem) {
                        var matrixAnswerColumn = matrixItem[matrixAnswerName]
                        
                        if (sessionData.surveyName  && questionKey && matrixAnswerColumn) {
                          var matrixObjectColummnName = this.buildColumnName(sessionData.surveyName, questionKey + '_matrix_' + matrixItemName + '_' + matrixAnswerName)
                                                    
                          if (this.targetArray(sessionData.surveyName) == 'event') {
                            eventColumns[matrixObjectColummnName] = ''
                          } else if (this.targetArray(sessionData.surveyName) == 'followup') {
                            followUpColumns[matrixObjectColummnName] = ''
                          } else {
                            nonEventColumns[matrixObjectColummnName] = ''
                          }
                        }
                      }
                      
                      
                    } else {
                      if (sessionData.surveyName  && questionKey && matrixItemName) {
                        var matrixStringColummnName = this.buildColumnName(sessionData.surveyName, questionKey + '_matrix_' + matrixItemName)
                        
                        if (this.targetArray(sessionData.surveyName) == 'event') {
                          eventColumns[matrixStringColummnName] = ''
                        } else if (this.targetArray(sessionData.surveyName) == 'followup') {
                          followUpColumns[matrixStringColummnName] = ''
                        } else {
                          nonEventColumns[matrixStringColummnName] = ''
                        }
                      }
                    }
                  }
                }
              // the question returns a simple string
              } else {
                if (sessionData.surveyName  &&  questionKey) {
                  var stringColumnName = this.buildColumnName(sessionData.surveyName, questionKey)
                  
                  if (this.targetArray(sessionData.surveyName) == 'event') {
                    eventColumns[stringColumnName] = ''
                  } else if (this.targetArray(sessionData.surveyName) == 'followup') {
                    followUpColumns[stringColumnName] = ''
                  } else {
                    nonEventColumns[stringColumnName] = ''
                  }
                }
              }
            }
          }
        }
        
        var orderedNoteColumns = ['SessionID', 'NoteDate', 'Note', 'IsTest', 'ReferralCode', 'Username'].sort()
        
        var fixedColumns = ['SurveyName', 'SessionID', 'SessionDataID', 'IsTest', 'ReferralCode', 'Username', 'CreatedAt', 'UpdatedAt']
        var orderedEventFixedColumns = ['EventID']
        var orderedEventColumns = Object.keys(eventColumns).concat(fixedColumns).concat(orderedEventFixedColumns)
        
        var orderedNonEventColumns = Object.keys(nonEventColumns).concat(fixedColumns)
        
        var followUpFixedColumns = ['EventID', 'ReferenceID']
        var orderedFollowUpColumns = Object.keys(followUpColumns).concat(fixedColumns).concat(followUpFixedColumns)
        
        // only sort if opted in
        if (this.sortAlpha) {
          orderedEventColumns.sort()
          orderedNonEventColumns.sort()
          orderedFollowUpColumns.sort()
        }
        
        
        // eslint-disable-next-line
        //console.log('columns:', orderedEventColumns, orderedNonEventColumns)
        
        var noteRows = []
        var eventRows = []
        var nonEventRows = []
        var followUpRows = []
        //
        // build the CSV in to the precompiled columns
        ///
        for (surveyID in this.data.export) {
          survey = this.data.export[surveyID]
          // eslint-disable-next-line
          //console.log('survey:', survey)
          
          var rowData = {}
          var thisColumn = ''
          
          if (survey.surveyState && survey.surveyState.notes) {
            for (var noteKey in survey.surveyState.notes) {
              var note = survey.surveyState.notes[noteKey]
              
              for (thisColumn of orderedNoteColumns) {
                rowData[thisColumn] = ''
              }
              
              rowData['SessionID'] = survey.sessionID
              rowData['IsTest'] = (survey.isTest) ? '1' : '0'
              rowData['ReferralCode'] = survey.referralCode
              rowData['Username'] = survey.username
              rowData['NoteDate'] = note.date
              rowData['Note'] = note.note
              
              noteRows.push(rowData)
            }
          }
          
          for (sessionDataID in survey.sessionData) {
            sessionData = survey.sessionData[sessionDataID]
            
            if (this.targetArray(sessionData.surveyName) == 'event') {
              for (thisColumn of orderedEventColumns) {
                rowData[thisColumn] = ''
              }
            } else if (this.targetArray(sessionData.surveyName) == 'followup') {
              for (thisColumn of orderedFollowUpColumns) {
                rowData[thisColumn] = ''
              }
            } else {
              for (thisColumn of orderedNonEventColumns) {
                rowData[thisColumn] = ''
              }
            }
            
            // populate the static content
            rowData['SurveyName'] = sessionData.surveyName
            rowData['SessionID'] = survey.sessionID
            rowData['SessionDataID'] = sessionData.sessionDataID
            rowData['IsTest'] = (survey.isTest) ? '1' : '0'
            rowData['ReferralCode'] = survey.referralCode
            rowData['Username'] = survey.username
            
            // creating a linking column between events and followups
            if (this.targetArray(sessionData.surveyName) == 'event') {
              rowData['EventID'] = sessionData.sessionDataID
            } else if (this.targetArray(sessionData.surveyName) == 'followup') {
              rowData['EventID'] = sessionData.referenceID
              rowData['ReferenceID'] = sessionData.referenceID
            }
            
            
            var createdDate = new Date(sessionData.createdAt)
            var updatedDate = new Date(sessionData.updatedAt)
            rowData['CreatedAt'] = createdDate.toLocaleString()
            rowData['UpdatedAt'] = updatedDate.toLocaleString()
            
            // iterate over the data
            for (questionKey in sessionData.data) {
              question = sessionData.data[questionKey]
              // the question is a dynamic table of rows
              if (typeof(question) == 'object') {
                
                if (Array.isArray(question)) {
                  for (index = 0; index <= question.length; index++) {
                    for (answerColumn in question[index]) {
                      if (typeof(question[index]) == 'object') {  
                        if (sessionData.surveyName  && questionKey && answerColumn) {
                          objectColummnName = this.buildColumnName(sessionData.surveyName, questionKey + '_row' + index + '_' + answerColumn)
                        
                          rowData[objectColummnName] = question[index][answerColumn]
                        }
                      } else {
                        if (sessionData.surveyName && questionKey && answerColumn) {
                          stringColumnName = this.buildColumnName(sessionData.surveyName, questionKey + '_' + answerColumn)
                        
                          rowData[stringColumnName] = question[index]
                        }
                      }
                    }
                  }
                } else {
                  // fucking matrix
                  
                  for (matrixItemName in question) {
                    matrixItem = question[matrixItemName]
                    
                    if (typeof(matrixItem) == 'object') {
                      for (matrixAnswerName in matrixItem) {
                        matrixAnswerColumn = matrixItem[matrixAnswerName]
                        
                        if (sessionData.surveyName  && questionKey && matrixAnswerColumn) {
                          matrixObjectColummnName = this.buildColumnName(sessionData.surveyName, questionKey + '_matrix_' + matrixItemName + '_' + matrixAnswerName)
                          
                          rowData[matrixObjectColummnName] = matrixAnswerColumn
                        }
                      }
                      
                    } else {
                      if (sessionData.surveyName  && questionKey && matrixItemName) {
                        matrixStringColummnName = this.buildColumnName(sessionData.surveyName, questionKey + '_matrix_' + matrixItemName)
                        
                        rowData[matrixStringColummnName] = matrixItem
                      }
                    }
                  }
                }
                
              // the question returns a simple string
              } else {
                if (sessionData.surveyName && questionKey) {
                  stringColumnName = this.buildColumnName(sessionData.surveyName, questionKey)
                  
                  rowData[stringColumnName] = question
                }
              }
            }
            
            if (this.targetArray(sessionData.surveyName) == 'event') {
              eventRows.push(rowData)
            } else if (this.targetArray(sessionData.surveyName) == 'followup') {
              followUpRows.push(rowData)
            } else {
              nonEventRows.push(rowData)
            }
          }
        }
        
        //
        // Build CSV output
        //
        this.nonEventRows = nonEventRows.length
        this.nonEventCSV = this.objectsToCSV(orderedNonEventColumns, nonEventRows)
        
        this.eventRows = eventRows.length
        this.eventsCSV = this.objectsToCSV(orderedEventColumns, eventRows)
        
        this.followUpRows = followUpRows.length
        this.followUpCSV = this.objectsToCSV(orderedFollowUpColumns, followUpRows)
        
        this.notesRows = noteRows.length
        this.notesCSV = this.objectsToCSV(orderedNoteColumns, noteRows)
        
        var date = new Date()
        var timestamp = '' + date.getFullYear() + this.pad(date.getMonth()) + date.getDate();
        
        // download as file
        var nonEventData = new Blob([this.nonEventCSV]);
        var nonEventDownloadButton = document.getElementById('nonEventDownload')
        nonEventDownloadButton.href = URL.createObjectURL(nonEventData);
        nonEventDownloadButton.download = timestamp + '-surveys' + '.csv'
        
        
        var eventData = new Blob([this.eventsCSV]);
        var eventDownloadButton = document.getElementById('eventDownload')
        eventDownloadButton.href = URL.createObjectURL(eventData);
        eventDownloadButton.download = timestamp + '-events' + '.csv'
        
        var followUpData = new Blob([this.followUpCSV]);
        var followUpDownloadButton = document.getElementById('followUpDownload')
        followUpDownloadButton.href = URL.createObjectURL(followUpData);
        followUpDownloadButton.download = timestamp + '-followup' + '.csv'
        
        var notesData = new Blob([this.notesCSV]);
        var notesDownloadButton = document.getElementById('notesDownload')
        notesDownloadButton.href = URL.createObjectURL(notesData);
        notesDownloadButton.download = timestamp + '-notes' + '.csv'

      }
    },
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>