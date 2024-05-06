<template>
  <v-col cols="8" offset="2">
    <v-dialog v-model="showPreview">
      <v-card class="pt-4">
        <v-card-title class="headline">{{ previewName }}</v-card-title>
        <v-card-text>
          <v-data-table
            :headers="previewHeaders"
            :items="previewRows"
            class="preview elevation-1 mb-5"
          >
          </v-data-table>
        </v-card-text>
      </v-card>
    </v-dialog>
    
    
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
          <v-toolbar-title class="headline font-weight-regular"><v-icon dark class="mr-2">import_export</v-icon> Data Download</v-toolbar-title>
        </v-toolbar>
        
        <v-progress-linear :indeterminate="true" v-if="loading"></v-progress-linear>
        
        <v-data-table
          :headers="exportListHeaders"
          :items="exportList"
          class="elevation-1"
          :search="search"
          :footer-props="{
              itemsPerPageOptions: [20, 50, 100],
            }"
        >
          <template v-slot:item.surveyName="{ item }">
            <span class="font-weight-medium">{{ item.surveyName }}</span>
          </template>
          
          <template v-slot:item.rows="{ item }">
            {{ item.rows }}
          </template>
          
          <template v-slot:item.actions="{ item }">
            <v-icon
              small
              @click="downloadBlob(item.surveyName)"
            >
              save_alt
            </v-icon> &nbsp; 
            <v-icon
              small
              @click="viewTable(item.surveyName)"
            >
              list_alt
            </v-icon>
          </template>
        </v-data-table>
      </v-card>
    </v-row>
  </v-col>
</template>


<script>
  import { mapState } from 'vuex'
  import * as answerHelper from '../helpers/answers.js'
  
  
  export default {
    data () {
      return {
        loading: false,
        
        error: false,
        errorText: 'an unknown error occurred',
        
        search: '',
        
        exportListHeaders: [
          { text: 'Survey Name', value: 'surveyName' },
          { text: 'Rows', value: 'rows' },
          { text: '', value: 'actions', sortable: false, align: 'center' }
        ],
        
        data: {},
        surveyData: {},
        
        dataBySurvey: {},
        
        surveyColumns: {},
        surveyLargestDynamicMatrix: {},
        surveyRows: {},
        
        exportList: [],
        
        showPreview: false,
        previewName: '',
        previewHeaders: [],
        previewRows: []
      }
    },
    computed: {
      ...mapState('surveySettings', ['surveyType', 'config', 'surveyTypes']),
      ...mapState('surveys', ['surveys']),
    },
    watch: {
      
    },
    mounted() {
      var that = this
      this.loading = true
      
      this.axios.get('/survey/export')
      .then(response => { 
        that.data = response.data
        setTimeout(function() {
          that.loadData()
        }, 1000);
      }).catch(function (error) {
        return 'error: ' + error
      })
    },
    methods: {
      downloadBlob(surveyName) {
        var csvData = this.objectsToCSV(this.surveyColumns[surveyName], this.surveyRows[surveyName])
        var csvBlob = new Blob([csvData]);
        
        var date = new Date()
        var timestamp = '' + date.getFullYear() + this.pad(date.getMonth()) + date.getDate();
        
        const link = document.createElement('a');
        link.href = URL.createObjectURL(csvBlob);
        link.download = timestamp + ' - ' + surveyName + '.csv'
        
        document.body.appendChild(link);
        
        link.click();
      },
      
      viewTable(surveyName) {
        this.previewName = surveyName
        
        this.previewHeaders = []
        
        for (var columnName of this.surveyColumns[surveyName]) {
          this.previewHeaders.push({'text': columnName, 'value': columnName, 'align': 'center', 'divider': true})
        }
        
        this.previewRows = []
        for (var row of this.surveyRows[surveyName]) {
          this.previewRows.push(row)
        }
        
        this.showPreview = true
      },
      
      updateExportList() {
        this.exportList = []
        
        for (var surveyName in this.surveyRows) {
          this.exportList.push({'surveyName': surveyName, 'rows': this.surveyRows[surveyName].length})
        }
      },
      
      answersToRows(surveyName, surveyJSON, surveyData, answers, rebuildColumns) {
        // eslint-disable-next-line
        //console.log(surveyName, surveyJSON, answers)
        
        var fixedColumns = ['SurveyName', 'SessionID', 'SessionDataID', 'DataKey', 'ReferenceID', 'IsTest', 'ReferralCode', 'Username', 'CreatedAt', 'UpdatedAt']
        
        // parse surveyJSON and determine ALL columns
        if (!this.surveyColumns[surveyName] || rebuildColumns) {          
          this.surveyColumns[surveyName] = []
          
          if (!this.surveyRows[surveyName]) {
            this.surveyRows[surveyName] = []
          }
        
          for (thisColumn of fixedColumns) {
            this.surveyColumns[surveyName].push(thisColumn)
          }  
          
          var buildDynamicMatrixColumns = false
          
          for (var pageID in surveyJSON.pages) {
            var page = surveyJSON.pages[pageID]
            
            for (var elementID in page.elements) {
              var element = page.elements[elementID]
              
              if (element.type == 'matrixdynamic') {
                buildDynamicMatrixColumns = true
              } else {
                buildDynamicMatrixColumns = false
              }
              
              if (element.rows) {
                for (var elementRowID in element.rows) {
                  var row = element.rows[elementRowID]
                  
                  var hasColumn = false
                  if (element.columns) {
                    for (var elementColumnID in element.columns) {
                      var column = element.columns[elementColumnID]
                      if (typeof(column.name) != 'undefined') {
                        hasColumn = true
                        this.surveyColumns[surveyName].push(element.name + '.' + this.rowOrValue(row) + '.' + column.name)
                      }
                    }
                  }
                  
                  if (!hasColumn) {
                    this.surveyColumns[surveyName].push(element.name + '.' + this.rowOrValue(row))
                  }
                }
              } else if (element.columns) {
                for (var elementAnswerColumnID in element.columns) {
                  var answerColumn = element.columns[elementAnswerColumnID]
                  
                  if (buildDynamicMatrixColumns && this.surveyLargestDynamicMatrix[surveyName] > 0) {
                    for (var x = 0; x <= this.surveyLargestDynamicMatrix[surveyName]; x++) {
                      this.surveyColumns[surveyName].push(element.name + '.' + x + '.' + answerColumn.name)
                    }
                  } else {
                    this.surveyColumns[surveyName].push(element.name + '.' + answerColumn.name)
                  }
                }
              } else {
                this.surveyColumns[surveyName].push(element.name)
              }
              //console.log(this.surveyColumns[surveyName][this.surveyColumns[surveyName].length - 1])
            }
          }
        }
        
        // run through the answers and create a single row
        var rowData = {}
        var columnName, thisColumn
        
        for (thisColumn of this.surveyColumns[surveyName]) {
          rowData[thisColumn] = ''
        }        
        
        // add fixed properties
        var survey = this.data.export.find(ref => ref.sessionID == surveyData.sessionID)        
        
        rowData['SurveyName'] = surveyData.surveyName
        rowData['SessionID'] = survey.sessionID
        rowData['SessionDataID'] = surveyData.sessionDataID
        rowData['DataKey'] = surveyData.dataKey
        rowData['ReferenceID'] = surveyData.referenceID
        rowData['IsTest'] = (survey.isTest) ? '1' : '0'
        rowData['ReferralCode'] = survey.referralCode
        rowData['Username'] = survey.username
        
        
        var createdDate = new Date(surveyData.createdAt)
        var updatedDate = new Date(surveyData.updatedAt)
        
        rowData['CreatedAt'] = createdDate.toString().split(' ').splice(0, 5).join(' ')
        rowData['UpdatedAt'] = updatedDate.toString().split(' ').splice(0, 5).join(' ')
        
        
        // populate answers in to the row
        for (var answerID in answers) {
          var answer = answers[answerID]
          
          if (answer.type == 'multiple') {
            for (var subAnswerID in answer.values) {
              var subAnswer = answer.values[subAnswerID]
              
              if (Array.isArray(subAnswer)) {
                for (var subSubAnswerID in subAnswer) {
                  var subSubAnswer = subAnswer[subSubAnswerID]
                  
                  columnName = answer.name + '.' + subSubAnswer.name
                  
                  if (typeof(rowData[columnName]) == 'undefined') {
                    // eslint-disable-next-line
                    console.log(columnName + ' not found in subSubAnswer row', answer)
                  } else {
                    rowData[columnName] = subSubAnswer.value
                  }
                }
              } else {
                columnName = answer.name + '.' + subAnswer.name
                
                if (typeof(rowData[columnName]) == 'undefined') {
                  // eslint-disable-next-line
                  console.log(columnName + ' not found in subAnswer row', answer)
                } else {
                  rowData[columnName] = subAnswer.value
                }
              }
            }
            
            
          } else {
            if (typeof(rowData[answer.name]) == 'undefined') {
              // eslint-disable-next-line
              console.log(answer.name + ' not found in answer row', answer)
            } else {
              rowData[answer.name] = answer.value
            }
          }
        }
        
        this.surveyRows[surveyName].push(rowData)
      },
      
      rowOrValue(row) {
        if (typeof(row) == 'string') {
          return row
        } else {
          return row.value
        }
      },
      
      loadData() {        
        var that = this
        
        // build a list of all the surveys with data
        for (var surveyIndex in this.data.export) {
          var survey = this.data.export[surveyIndex]
          //console.log(survey)
          
          if (survey.surveyState && survey.surveyState.notes) {
            if (!this.surveyRows['Survey Notes']) {
              this.surveyRows['Survey Notes'] = []
              this.surveyColumns['Survey Notes'] = ['SessionID', 'Username', 'NoteDate', 'Note']
            }
            
            for (var noteKey in survey.surveyState.notes) {
              var note = survey.surveyState.notes[noteKey]
              
              var rowData = {}
              rowData['SessionID'] = survey.sessionID
              rowData['Username'] = survey.username
              rowData['NoteDate'] = note.date
              rowData['Note'] = note.note.replace(/\n/g, " ");
              
              this.surveyRows['Survey Notes'].push(rowData)
            }
          }
          
          for (var sessionDataID in survey.sessionData) {
            var sessionData = survey.sessionData[sessionDataID]
            if (!this.dataBySurvey[sessionData.surveyName]) {
              this.dataBySurvey[sessionData.surveyName] = []
            }
            this.dataBySurvey[sessionData.surveyName].push(sessionData)
          }
        }
        
        for (var surveyName in this.dataBySurvey) {
          this.axios.get('/survey/' + surveyName)
          .then(response => {
            if (response.data.surveyJSON) {
              var surveyName = response.data.name
              var surveyJSON = JSON.parse(response.data.surveyJSON)
              var rebuildColumns = false
              
              for (var row in this.dataBySurvey[surveyName]) {
                var result = answerHelper.parseAnswers(surveyJSON, this.dataBySurvey[surveyName][row])
                
                var answers = result[0]
                var largestDynamicMatrix = result[1]
                
                if (largestDynamicMatrix > 0 && !that.surveyLargestDynamicMatrix[surveyName] || largestDynamicMatrix > that.surveyLargestDynamicMatrix[surveyName]) {
                  that.surveyLargestDynamicMatrix[surveyName] = largestDynamicMatrix
                  rebuildColumns = true
                } else {
                  rebuildColumns = false
                }
                
                that.answersToRows(surveyName, surveyJSON, this.dataBySurvey[surveyName][row], answers, rebuildColumns)
              }
              
              that.updateExportList()
              that.loading = false
              
            } else {
              that.errorText = 'Survey data not attached to server response'
              that.error = true
            }
          }).catch(function (error) {
            that.error = true
            if (error.response && error.response.data) {
              that.errorText = error.response.data
            } else if (error.message) {
              that.errorText = error.message
            } else {
              that.errorText = 'Error loading survey...'
            }
          })
        }
        
      },
      
      pad(n) {
          return n<10 ? '0'+n : n;
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
      }
    },
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  .preview {
    white-space: nowrap
  }
</style>