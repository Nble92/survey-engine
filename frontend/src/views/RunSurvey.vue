<template>
  <v-col cols="12">
    <v-row no-gutters>  
      <!-- 
        -- 
        -- Test session --
        --
      -->
      <div v-if="specialHeaderText" style="width: 100%; text-align: center;">
        <h2 class="grey--text" ma-0>{{ specialHeaderText }}</h2>
      </div>
      
      <!-- 
        -- 
        -- Instructions --
        --
      -->
      <v-dialog v-model="fullInstructions" width="70%">
        <v-card class="pa-3" style="width: 100%">
          <v-tabs
            v-model="instructionTab"
            color="info"
            tabs
          >
            <v-tabs-slider color="black"></v-tabs-slider>
            <v-tab
              v-for="page in instructionObject.instructionPages"
              :key="page.pageName"
              :href="'#tab-' + page.pageName"
            >
              {{ page.pageName }}
            </v-tab>
          </v-tabs>
          
          <v-tabs-items v-model="instructionTab">
            <v-tab-item
              v-for="page in instructionObject.instructionPages"
              :key="page.pageName"
              :value="'tab-' + page.pageName"
            >
              <v-card flat>
                <v-card-text v-html="page.pageContent" class="pt-5 ql-viewer"></v-card-text>
              </v-card>
            </v-tab-item>
          </v-tabs-items>
          
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="fullInstructions = false">Close</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      
      <v-flex v-if="instructionObject && instructionObject.summaryInstruction" xs12 ma-3>
        <v-alert :value="true" type="info" outlined>
          <div style="float: right; text-align: right">
            <v-btn color="info" @click="fullInstructions = true"  v-if="!instructionObject.popoverWindow" outlined>Instructions</v-btn>
            <v-btn color="info" @click="popoutInstructions"  v-if="instructionObject.popoverWindow" outlined>Open Window</v-btn>
          </div>

          <div v-html="instructionObject.summaryInstruction" class="black--text"></div>
        </v-alert>
      </v-flex>
      

      <!-- 
        -- 
        -- NEW --
        --
      -->
      <v-container fluid v-if="status == 'new'">
        <v-row>
          <v-col cols="12">
            <v-alert
              v-model="error"
              type="error"
              transition="scale-transition"
              style="width: 100%"
            >{{ errorText }}</v-alert>
            <div class="text-center" style="width: 100%">
              <v-icon x-large class="mb-4" style="font-size: 100px;" color="blue darken-1">assignment</v-icon>
              
              <h2 class="headline" v-if="!config.startSurveyInstructions || config.startSurveyInstructions == ''">Welcome to the {{ friendlyName }} survey.</h2>
              <div v-else v-html="config.startSurveyInstructions"></div>
              
              <div class="subheading pa-2" v-if="isLoggedIn">
                To begin formal data collection, start new survey.
                <br>
                To practice or review the data collection flow, start test survey.
              </div>
              <div class="mt-4" v-if="isLoggedIn">
                <v-btn large color="primary" class="ma-2" @click="startSurvey(false)">Start New Survey</v-btn>
                <v-btn large color="primary" class="ma-2" @click="startSurvey(true)" outlined>Start Test Survey</v-btn>
              </div>
              <div class="mt-4" v-else>
                <v-btn large color="primary" class="ma-2" @click="startSurvey(false)">Start Survey</v-btn>
              </div>
            </div>
          </v-col>
        </v-row>
      </v-container>
      
      
      
      <!-- 
      -- 
      -- ERROR --
      --
      -->
      <v-container fluid v-else-if="status == 'error'">
        <v-row>
          <v-col cols="12">
            <div class="text-center text--red">
              <v-icon x-large class="mb-4" style="font-size: 100px;" color="red darken-1">assignment</v-icon>
              <h2 class="headline red--text" v-if="!errorText">There was a problem with this survey session.</h2>
              <h2 class="headline red--text" v-if="errorText">{{ errorText }}</h2>
              <div class="mt-4">
                <v-btn large color="error" outlined @click="resetSurvey">Reset Survey</v-btn>
              </div>
            </div>
          </v-col>
        </v-row>
      </v-container>
  
  
  
      <!-- 
        -- 
        -- SURVEY & COMPLETION SURVEY --
        --
      -->
      <v-flex v-else-if="status == 'in_process' && (currentView == 'survey' || currentView == 'completion')" pa-0>
        <v-card-text class="pa-0 pt-3 ql-viewer">
          <div class="pb-3" v-if="error">
            <v-alert
              v-model="error"
              type="error"
              transition="scale-transition"
            >{{ errorText }}</v-alert>
          </div>
    
          <v-layout wrap align-center v-if="surveyLoadJSON == ''">
            <v-flex pa-3>
              <v-progress-linear :indeterminate="true" v-if="!surveyLoadTimeout"></v-progress-linear>
              <div style="text-align: center;"><v-btn outlined color="primary" @click="resetSurvey" v-if="surveyLoadTimeout">Hung? Restart Survey Session</v-btn></div>
            </v-flex>
          </v-layout>
          
          <survey-viewer :survey-name="surveyName" @partialData="surveyPartialData" @completed="surveyCompleted" :surveyJSON="surveyLoadJSON" :surveyData="surveyData"  v-if="surveyLoadJSON != ''"></survey-viewer>
          
        </v-card-text>
      </v-flex>
    
    
      <!-- 
        -- 
        -- CALENDAR --
        --
      -->
      <v-flex v-else-if="status == 'in_process' && currentView == 'calendar'" pa-0>
        
        <v-dialog v-model="showNewNote" width="50%">
          <v-card class="pa-3 pt-5">
            <v-card-text>
              <v-textarea
                v-model="noteText"
                auto-grow
                label="Note"
              ></v-textarea>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="blue darken-1" text @click="cancelNote">Cancel</v-btn>
              <v-btn color="blue darken-1" text @click="saveNote">Save</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      
        <v-card-text>
          <v-container fluid grid-list-md pa-0>
            <v-layout align-start justify-space-around row fill-height>
              <v-flex xs12 md10>
                <div class="pb-3" v-if="error">
                  <v-alert
                    v-model="error"
                    type="error"
                    transition="scale-transition"
                  >{{ errorText }}</v-alert>
                </div>
            
                <calendar-viewer 
                  :survey-name="surveyName" 
                  :surveyConfig="config" 
                  :surveyState="surveyState"
                  :surveyCreated="surveyCreated"
                  :surveySessionData="surveySessionData"
                  @dayClicked="calendarDayClick" 
                  @dayRightClicked="calendarDayRightClick" 
                  @eventClicked="eventClick"
                  @eventDeleted="deleteEvent"
                  @eventMoved="moveEvent"
                  @eventCopied="copyEvent"
                  @noteClicked="calendarDayRightClick"
                  @noteDeleted="deleteNote"
                  @noteMoved="moveNote"
                  @noteCopied="copyNote"
                >
                </calendar-viewer>
              </v-flex>
            </v-layout>
            <v-layout align-start justify-space-around row mt-5>
              <v-flex xs12 md10>
                <v-toolbar>
                  <v-btn color="info" @click="popoutInstructions"  v-if="instructionObject && !instructionObject.summaryInstruction && instructionObject.popoverWindow" outlined>Open Instructions</v-btn>
                  <v-btn color="info" @click="fullInstructions = true"  v-if="instructionObject && !instructionObject.summaryInstruction && !instructionObject.popoverWindow" outlined>Instructions</v-btn>
                  
                  <v-spacer></v-spacer>
                  
                  <logic-button class="ma-2" :functionName="followUpFunction" functionType="followup" :buttonLabel="followUpButtonLabel" :surveySessionData="surveySessionData" @click="runFollowUpSurveys" v-if="followUpFunction"></logic-button>
                  <logic-button class="ma-2" :functionName="completionFunction" functionType="completion" :buttonLabel="completionButtonLabel" :surveySessionData="surveySessionData" @click="runCompletionSurvey" v-if="completionFunction"></logic-button>
                </v-toolbar>
              </v-flex>
            </v-layout>
          </v-container>
        </v-card-text>
      </v-flex>
    
    
      <!-- 
        -- 
        -- FOLLOWUP SURVEYS --
        --
      -->
      <v-flex v-else-if="status == 'in_process' && currentView == 'followup'" pa-0>
        <v-card-text class="pa-0 pt-3 ql-viewer">
          <v-layout column wrap align-center grey lighten-4 elevation-4 v-if="surveyLoadJSON == ''">
            <v-flex pa-3>
              <div class="pa-3 text-center">
                <v-progress-circular :indeterminate="true" color="primary" size="50"></v-progress-circular>
              </div>
              <v-btn outlined color="primary" @click="resetSurvey">Hung? Restart Survey Session</v-btn>
            </v-flex>
          </v-layout>

          <div class="pb-3" v-if="error">
            <v-alert
              v-model="error"
              type="error"
              transition="scale-transition"
            >{{ errorText }}</v-alert>
          </div>
          <follow-up-viewer :survey-config="config" :survey-name="surveyName" @partialData="surveyPartialData" @completed="surveyCompleted" :surveyJSON="surveyLoadJSON" :surveyData="surveyData" :followUpData="surveyFollowUpData"></follow-up-viewer>
        </v-card-text>
      </v-flex>
      
      <!-- 
        -- 
        -- COMPLETED ALL --
        --
      -->
      <v-container v-else-if="status == 'completed'">
        Completed, but you shouldnt be seeing this...
      </v-container>
      
      
    </v-row>
  </v-col>
</template>

<script>
  import { mapState, mapGetters } from 'vuex'
  
  import SurveyViewer from '../components/SurveyViewer.vue'
  import CalendarViewer from '../components/CalendarViewer.vue'
  import LogicButton from '../components/LogicButton.vue'
  import FollowUpViewer from '../components/FollowUpViewer.vue'
  
  
  export default {
    components: {
      SurveyViewer,
      CalendarViewer,
      LogicButton,
      FollowUpViewer
    },
    data () {
      return {
        currentSurveyName: '',
        currentSurveySessionDataID: '',
        
        surveyLoadJSON: '',
        surveyLoadTimeout: false,
        
        surveySessionData: {},
        surveyFollowUps: [],
        surveyFollowUpData: {},
        
        fullInstructions: false,
        instructionObject: false,
        instructionTab: null,
        
        showNewNote: false,
        noteDate: '',
        noteText: '',
      }
    },
    computed: {
      specialHeaderText() { 
        if (this.editMode && this.isTest) {
          return 'EDITING TEST SESSION FROM ' + this.surveyCreated
        } else if (this.editMode && !this.isTest) {
          return 'EDITING SESSION FROM ' + this.surveyCreated
        } else if (this.isTest) {
          return 'TEST SESSION'
        }
        return false
      },
      
      followUpFunction() {
        return (this.config.followUpFunction) ? this.config.followUpFunction : false
      },
      followUpButtonLabel() {
        return (this.config.followUpButtonLabel) ? this.config.followUpButtonLabel : 'NO LABEL'
      },
      completionFunction() {
        return (this.config.completionFunction) ? this.config.completionFunction : false
      },
      completionButtonLabel() {
        return (this.config.completionButtonLabel) ? this.config.completionButtonLabel : 'NO LABEL'
      },
      
      multipleRuns() {
        return this.config.userMultipleRuns
      },
      
      showRedirectButton() {
        return this.config.allowRedirectButton
      },
      redirectButtonInstruction() {
        return this.config.redirectButtonInstructions
      },
      redirectURL() {
        return this.config.redirectURL
      },
      
      ...mapState('config', ['friendlyName', 'contact']),
      ...mapState('surveySettings', ['surveyType', 'config', 'surveyTypes']),
      ...mapState('surveys', ['surveys']),
      ...mapGetters('user', ['isLoggedIn']),
      ...mapState('user', ['username']),
      ...mapState('surveySession', ['status', 'sessionID', 'surveyCreated', 'surveyState', 'isTest', 'editMode']),
      
      ...mapState('runState', ['error', 'errorText', 'currentView', 'surveyName', 'surveySessionDataID', 'surveyDataKey', 'surveyReferenceID', 'surveyData', 'allowReset', 'lastReset'])
    },
    mounted() {
      this.loadSurveyState()
    },
    watch: {
      status() {
        if (this.surveyState) {
          this.$store.dispatch('runState/statusChange')
        }
      },
      surveyName() {
        this.surveyLoadTimeout = false
        
        if (this.surveyName && (this.surveyName != this.currentSurveyName)) {
          if (this.surveySessionDataID) {
            this.loadExistingSurveyWithSessionDataID(this.surveyName, this.surveySessionDataID)
          } else if (this.surveyDataKey) {
            this.loadExistingSurveyWithDataKey(this.surveyName, this.surveyDataKey)
          }
        }
      },
      surveySessionDataID() {
        this.surveyLoadTimeout = false
        
        if (this.surveyName && this.surveySessionDataID && (this.surveySessionDataID != this.currentSurveySessionDataID)) {
          if (this.surveySessionDataID) {
            this.loadExistingSurveyWithSessionDataID(this.surveyName, this.surveySessionDataID)
          } else if (this.surveyDataKey) {
            this.loadExistingSurveyWithDataKey(this.surveyName, this.surveyDataKey)
          }
        }
      },
      lastReset() {
        // eslint-disable-next-line
        console.log('RESET due to lastReset change')
        this.surveyLoadTimeout = false
        this.surveyLoadJSON = ''
        this.currentSurveyName = ''
        this.currentSurveySessionDataID = ''
      }
    },
    methods: {
      loadSurveyState() {
        // we have a sessionID but not loaded state
        if (this.sessionID) {
          var that = this
          
          this.$store.dispatch('surveySession/fetchSessionState').then(() => {
            if (this.surveyName) {
              that.loadExistingSurveyWithDataKey(this.surveyName, this.surveyDataKey)
            }
            
            // load session data for logicButton processing
            that.updateSurveySessionData()
          }).catch(function(error) {
            var errorText = ''
            if (error.response && error.response.data) {
              if (error.response.data.message) {
                errorText = error.response.data.message
              } else {
                errorText = error.response.data
              }
            } else if (error.message) {
              errorText = error.message
            } else {
              errorText = 'Error loading session state...'
            }
            
            that.$store.dispatch('runState/updateState', {error: true, errorText: errorText})
          })
        }
      },
      
      // surveySessionData is used for logicButton data
      updateSurveySessionData() {
        var that = this
        if (!this.error) {
          this.axios.get('/survey/sessionDump/' + this.sessionID)
          .then(response => { 
            that.surveySessionData = response.data
            
            // eslint-disable-next-line
            console.log('updateSurveySessionData', that.surveySessionData)
          }).catch(function (error) {
            var errorText = ''
            if (error.response && error.response.data) {
              if (error.response.data.message) {
                errorText = error.response.data.message
              } else {
                errorText = error.response.data
              }
            } else if (error.message) {
              errorText = error.message
            } else {
              errorText = 'Error fetching session data...'
            }
            
            this.$store.dispatch('runState/updateState', {error: true, errorText: errorText})
          })
        }
      },
      
      updateInstructionObject() {
        var surveyName = (this.currentView == 'calendar') ? 'calendar' : this.surveyName
        if (surveyName) {
          var that = this
          this.axios.get('/instructions/' + surveyName)
          .then(response => {
            if (response.data.data) {
              that.instructionObject = response.data.data
            }
          }).catch(() => {
            // do nothing, we dont care
          })
        }
      },
      
      popoutInstructions() {
        var surveyName = (this.currentView == 'calendar') ? 'calendar' : this.surveyName
        if (surveyName && this.instructionObject) {
          window.open('//' + location.hostname + ':' + location.port + '/instructionHTML/' + surveyName, 'targetWindow', 'toolbar=no,location=no,status=no,menubar=no,scrollbars=yes,resizable=yes,width=400px,height=600px');
        }
      },
      
      resetSurvey() {
        this.$store.dispatch('surveySession/resetSessionState')
        this.$store.dispatch('runState/resetState')
      },
      
      startSurvey(isTest = false) {
        // currently all surveys need an entrySurvey defined, this will likely change as i evolve this component
        if (this.config.entrySurvey) {
          var that = this
          
          this.surveyLoadJSON = ''
          
          this.$store.dispatch('runState/resetState')
          
          var referralCode = (this.$route.query.referralCode) ? this.$route.query.referralCode : ''
          
          // create new survey session
          this.$store.dispatch('surveySession/newSurveySession', {surveyType: this.surveyType, username: this.username, isTest: isTest, referralCode: referralCode}).then(() => {
            // change our surveyName
            that.$store.dispatch('runState/updateState', {surveyName: this.config.entrySurvey})
            
            that.loadExistingSurveyWithDataKey(this.surveyName, this.surveyDataKey)
          }).catch(function(error) {
            var errorText = ''
            
            if (error.response && error.response.data) {
              if (error.response.data.message) {
                errorText = error.response.data.message
              } else {
                errorText = error.response.data
              }
            } else if (error.message) {
              errorText = error.message
            } else {
              errorText = 'Error creating new survey session...'
            }
            
            that.$store.dispatch('runState/updateState', {error: true, errorText: errorText})
          })
          
        } else {
          this.$store.dispatch('runState/updateState', {error: true, errorText: 'entry survey is not defined'})
        }
      },
      
      
      runFollowUpSurveys(surveySessionDataIds) {
        for (var index in surveySessionDataIds) {
          this.surveyFollowUps.push({'id': surveySessionDataIds[index], 'complete': false})
        }
        
        if (this.surveyFollowUps.length > 0) {
          // eslint-disable-next-line
          console.log('surveyFollowUps', this.surveyFollowUps)
          this.nextFollowUpSurvey()
        }
      },
      
      nextFollowUpSurvey() {
        var nextSessionDataID = this.surveyFollowUps.find(ref => ref.complete == false)
        
        if (nextSessionDataID && this.surveySessionData) {
          // eslint-disable-next-line
          console.log('nextFollowUpSurvey', nextSessionDataID)
          
          this.surveyLoadJSON = ''
          
          this.$store.dispatch('runState/resetState')
          this.$store.dispatch('runState/updateState', {surveyName: this.config.followUpSurvey, currentView: 'followup', surveyDataKey: 'followup', surveyReferenceID: nextSessionDataID.id})
          
          this.surveyFollowUpData = this.surveySessionData.sessionData.find(ref => ref.sessionDataID == this.surveyReferenceID)
                    
          var existingSessionData = this.surveyState.dataRefs.find(ref => ref.referenceID === this.surveyReferenceID)
          if (existingSessionData) {
            // eslint-disable-next-line
            console.log('existing session data for referenceID', existingSessionData.sessionDataID)
            this.loadExistingSurveyWithSessionDataID(this.surveyName, existingSessionData.sessionDataID)
          } else {
            // eslint-disable-next-line
            console.log('no existing session data for referenceID', this.surveyReferenceID)
            this.loadSurvey(this.surveyName)
          }
        } else if (nextSessionDataID == undefined) {
          // all done!
          this.runCompletionSurvey()
        } else {
          this.$store.dispatch('runState/updateState', {error: true, errorText: 'problem with nextFollowUp state...'})
        }
      },
      
      runCompletionSurvey() {
        // eslint-disable-next-line
        console.log('runCompletionSurvey')
        
        this.surveyLoadJSON = ''
        
        this.$store.dispatch('runState/resetState')
        this.$store.dispatch('runState/updateState', {surveyName: this.config.completionSurvey, surveyDataKey: 'completion', currentView: 'survey'})
        
        this.loadSurvey(this.surveyName)
      },
      
      existingDataRef(dataKey) {                
        if (this.surveyState.dataRefs && this.surveyState.dataRefs.length > 0) {                
          var existingSessionData = this.surveyState.dataRefs.find(ref => ref.surveyName === this.surveyName && ref.dataKey == dataKey)
          
          // eslint-disable-next-line
          console.log('existingDataRef', existingSessionData)
          
          return (existingSessionData) ? existingSessionData : false
        } 
        return false
      },
      
      existingSessionDataID(sessionDataID) {
        if (this.surveyState.dataRefs && this.surveyState.dataRefs.length > 0) {                
          var existingSessionData = this.surveyState.dataRefs.find(ref => ref.sessionDataID == sessionDataID)
          
          // eslint-disable-next-line
          console.log('existingSessionDataID', existingSessionData, this.surveyState.dataRefs)
          
          return (existingSessionData) ? existingSessionData : false
        } 
        return false
      },
      
      loadExistingSurveyWithDataKey(surveyName, dataKey) {
        this.$store.dispatch('runState/updateState', {surveyName: surveyName, surveyDataKey: dataKey})
        
        // load existing survey data if we have it
        var existingSessionData = this.existingDataRef(dataKey)
        
        if (existingSessionData) {          
          this.loadExistingSurveyWithSessionDataID(surveyName, existingSessionData.sessionDataID)
        } else {
          // new event
          this.loadSurvey(surveyName)
        }
        
      },
      
      loadExistingSurveyWithSessionDataID(surveyName, sessionDataID) {

        var that = this
        console.log("SURVEYYYY NAME!!! " + surveyName + " " + "ID!!!!! " +sessionDataID)
        this.currentSurveySessionDataID = sessionDataID
        this.$store.dispatch('runState/updateState', {surveyName: surveyName, surveySessionDataID: sessionDataID})
        
        // load existing survey data if we have it
        var existingSessionData = that.existingSessionDataID(sessionDataID)
        
        if (existingSessionData) {
          this.$store.dispatch('runState/updateState', {surveyDataKey: existingSessionData.dataKey})
          
          if (existingSessionData && existingSessionData.sessionDataID) {
            // eslint-disable-next-line
            console.log('loading existing survey data for:', existingSessionData.sessionDataID)
            
            this.axios.get('/survey/sessionData/' + existingSessionData.sessionDataID)
            .then(response => {
              var surveyData = response.data
              
              if (surveyData.sessionID == that.sessionID) {
                that.loadSurvey(surveyName, surveyData.data)
                
              } else {
                this.$store.dispatch('runState/updateState', {error: true, errorText: 'survey sessionID missmatch on sessionData'})
              }
            })
          }
        } else {
          this.$store.dispatch('runState/updateState', {currentView: 'calendar'})
          this.$store.dispatch('runState/updateState', {error: true, errorText: 'unable to load sessionData for this event'})
        }
        
      },
      
      loadSurvey(surveyName, surveyData = false) {
        // eslint-disable-next-line
        console.log('loadSurvey: ' + surveyName)
        
        var that = this
        
        if (surveyName) { 
          setTimeout(function(){
            if (that.surveyLoadJSON == '') {
              that.surveyLoadTimeout = true
            }
          }, 5000);
          
          this.axios.get('/survey/' + surveyName).then(response => {
            if (response.data.surveyJSON) {
              that.$store.dispatch('runState/updateState', {surveyName: surveyName})
              that.currentSurveyName = surveyName
              
              that.surveyLoadJSON = response.data.surveyJSON
              
              if (surveyData) {
                that.$store.dispatch('runState/updateState', {surveyData: surveyData})
              }
              
              that.updateInstructionObject()
            } else {
              that.$store.dispatch('runState/updateState', {error: true, errorText: 'Survey data not attached to server response'})
            }
          }).catch(function (error) {
            var errorText = ''
            if (error.response && error.response.data) {
              if (error.response.data.message) {
                errorText = error.response.data.message
              } else {
                errorText = error.response.data
              }
            } else if (error.message) {
              errorText = error.message
            } else {
              errorText = 'Error loading session state...'
            }
            
            that.$store.dispatch('runState/updateState', {error: true, errorText: errorText})
          })
        }
      },
      
      surveyPartialData(results) {
        this.partialSurveySessionData(results)
      },
      
      partialSurveySessionData(results) {
        var dispatchEndpoint = (this.surveySessionDataID) ? 'surveySession/updateSurveySessionData' : 'surveySession/newSurveySessionData'
        this.$store.dispatch(dispatchEndpoint, {sessionID: this.sessionID, sessionDataID: (this.surveySessionDataID) ? this.surveySessionDataID : false, surveyName: results.surveyName, dataKey: this.surveyDataKey, referenceID: this.surveyReferenceID, data: results.surveyResults, username: this.username}).then(response => {
          // eslint-disable-next-line
          console.log('partialSurveySessionData', response)
          
          this.$store.dispatch('runState/updateState', {surveySessionDataID: response.sessionDataID})
          this.$store.dispatch("surveySession/updateSurveySessionState")
        })
      },
      
      surveyCompleted(results) {
        var specialKeys = ['entry', 'completion']
        if (Object.keys(results.surveyResults).length == 0 && !specialKeys.includes(this.surveyDataKey)) {
          this.$store.dispatch('runState/updateState', {error: true, errorText: 'Completed without any answering any questions, ignoring submission...'})
          
          // eslint-disable-next-line
          console.log('Not storing results for ' + results.surveyName)
          this.$store.dispatch('runState/updateState', {currentView: 'calendar'})
          return
        }
        
        var that = this
        
        // easier to just clear this out each time
        this.instructionObject = false
        
        // TODO merge with partialSurveySessionData
        var dispatchEndpoint = (this.surveySessionDataID) ? 'surveySession/updateSurveySessionData' : 'surveySession/newSurveySessionData'
        this.$store.dispatch(dispatchEndpoint, {sessionID: this.sessionID, sessionDataID: (this.surveySessionDataID) ? this.surveySessionDataID : false, surveyName: results.surveyName, dataKey: this.surveyDataKey, referenceID: this.surveyReferenceID, data: results.surveyResults, username: this.username}).then(response => {
          // eslint-disable-next-line
          console.log('surveySessionData', response)
          
          // 
          // SURVEY TYPE NEXT STEP LOGIC HERE
          //
          if (that.surveyType == 'linear') {
            if (that.surveyDataKey == 'entry') {
              that.$store.dispatch("surveySession/updateSurveySessionState", {name: 'entrySurveyComplete', value: true}).then(() => {
                that.$store.dispatch("surveySession/updateSurveySessionState", {name: 'completed', value: true}).then(() => {
                  if (that.config.notifyOnComplete) {
                    that.axios.post('/notify', {subject: this.friendlyName + ' completion', message: 'There has been a new survey completed'})
                  }
                })
              })
            } else {
              that.$store.dispatch('runState/updateState', {error: true, errorText: 'surveyCompleted for a survey other than the entry survey'})
            }
            
          } else if (that.surveyType == 'calendar') {
            if (that.surveyDataKey == 'entry') {
              that.$store.dispatch("surveySession/updateSurveySessionState", {name: 'entrySurveyComplete', value: true}).then(() => {
                // display the calendar
                that.$store.dispatch('runState/updateState', {currentView: 'calendar'})
                
                that.updateInstructionObject()
                that.updateSurveySessionData()
                
                that.showReset = false
              })
            } else if (that.surveyDataKey == 'followup') {
              that.$store.dispatch("surveySession/updateSurveySessionState").then(() => {
                if (that.surveyReferenceID) {
                  // mark this surveyFollowUp as complete
                  for (var index in this.surveyFollowUps) {
                    if (that.surveyFollowUps[index].id == that.surveyReferenceID) {
                      that.surveyFollowUps[index].complete = true
                      that.surveyFollowUps[index].sessionDataID = that.surveySessionDataID 
                      
                      // eslint-disable-next-line
                      console.log('surveyFollowUps completed', that.surveyFollowUps[index])
                      break;
                    }
                  }
                }
                
                that.nextFollowUpSurvey()
              })
            } else if (that.surveyDataKey == 'completion') {
              that.$store.dispatch("surveySession/updateSurveySessionState", {name: 'completionSurveyComplete', value: true}).then(() => {
                that.$store.dispatch("surveySession/updateSurveySessionState", {name: 'completed', value: true}).then(() => {
                  if (that.config.notifyOnComplete) {
                    that.axios.post('/notify', {subject: this.friendlyName + ' completion', message: 'There has been a new survey completed'})
                  }
                })
              })
            } else {
              that.$store.dispatch("surveySession/updateSurveySessionState", {name: 'calendarHasEvents', value: true}).then(() => {
                // display the calendar
                that.$store.dispatch('runState/updateState', {currentView: 'calendar'})
                
                that.updateInstructionObject()
                that.updateSurveySessionData()
              })
            }
          }
          
        }).catch(function(error) {
          var errorText = ''
          if (error.response && error.response.data) {
            if (error.response.data.message) {
              errorText = error.response.data.message
            } else {
              errorText = error.response.data
            }
          } else if (error.message) {
            errorText = error.message
          } else {
            errorText = 'Error loading session state...'
          }
          
          that.$store.dispatch('runState/updateState', {error: true, errorText: errorText})
          
          if (that.surveyType == 'calendar') {
            that.$store.dispatch('runState/updateState', {currentView: 'calendar'})
          }
        })
      },
      
      deleteEvent(event) {
        // eslint-disable-next-line
        console.log('deleteEvent', event)
        
        var that = this
        
        this.$store.dispatch("surveySession/deleteSurveySessionData", event.id).then(() => {
          that.$store.dispatch("surveySession/updateSurveySessionState")
          that.updateSurveySessionData()
        }).catch(() => {
          that.$store.dispatch("surveySession/updateSurveySessionState")
        })
      },
      
      moveEvent(event, date) {
        // eslint-disable-next-line
        console.log('moveEvent', event, date)
        
        var that = this
        
        this.axios.get('/survey/sessionData/' + event.id)
        .then(response => {
          var surveyData = response.data
          
          if (surveyData.sessionID == that.sessionID) {
            that.$store.dispatch('surveySession/updateSurveySessionData', {sessionID: surveyData.sessionID, sessionDataID: surveyData.sessionDataID, surveyName: surveyData.surveyName, dataKey: date, referenceID: surveyData.referenceID, data: surveyData.data, username: surveyData.username}).then(() => {
              that.$store.dispatch("surveySession/updateSurveySessionState")
              that.updateSurveySessionData()
            })
          } else {
            this.$store.dispatch('runState/updateState', {error: true, errorText: 'survey sessionID missmatch on sessionData'})
          }
        })
      },
      
      copyEvent(event, date) {
        // eslint-disable-next-line
        console.log('copyEvent', event, date)
        
        var that = this
        
        this.axios.get('/survey/sessionData/' + event.id)
        .then(response => {
          var surveyData = response.data
          
          if (surveyData.sessionID == that.sessionID) {
            that.$store.dispatch('surveySession/newSurveySessionData', {sessionID: surveyData.sessionID, sessionDataID: false, surveyName: surveyData.surveyName, dataKey: date, referenceID: surveyData.referenceID, data: surveyData.data, username: surveyData.username}).then(() => {
              that.$store.dispatch("surveySession/updateSurveySessionState")
              that.updateSurveySessionData()
            })
          } else {
            this.$store.dispatch('runState/updateState', {error: true, errorText: 'survey sessionID missmatch on sessionData'})
          }
        })
      },
      
      calendarDayClick(date) {
        // eslint-disable-next-line
        console.log('calendarDayClick', date, this.config)
        
        this.instructionObject = false
        
        if (this.config && !this.config.allowMultipleEvents) {
          this.$store.dispatch('runState/updateState', {error: false})
          this.loadExistingSurveyWithDataKey(this.config.eventSurvey, date)
        } else {
          this.surveyLoadJSON = ''
          
          this.$store.dispatch('runState/resetState')
          this.$store.dispatch('runState/updateState', {surveyName: this.config.eventSurvey, surveyDataKey: date, currentView: 'survey', error: false})
          
          this.loadSurvey(this.surveyName)
        }
      },
      
      calendarDayRightClick(date) {
        this.noteDate = date
        
        if (this.surveyState.notes) {
          var existingNote = this.surveyState.notes.find(ref => ref.date == this.noteDate)
          if (existingNote) {
            this.noteText = existingNote.note
          }
          
        }
        
        this.showNewNote = true
      },
      
      eventClick(sessionDataID) {
        // eslint-disable-next-line
        console.log('eventClick', sessionDataID)
        
        this.instructionObject = false
        
        this.loadExistingSurveyWithSessionDataID(this.config.eventSurvey, sessionDataID)
        
        this.$store.dispatch('runState/updateState', {currentView: 'survey', error: false})
      },
      
      saveNote() {
        // such a HAX
        var notes = false
        
        if (this.surveyState.notes) {
          notes = JSON.parse(JSON.stringify(this.surveyState.notes))
        }
        
        if (!notes) {
          notes = []
        }
        
        // TODO fix this
        var hasNote = false
        for (var id in notes) {
          if (notes[id].date == this.noteDate) {
            notes[id].note = this.noteText
            hasNote = true
            break
          }
        }
        
        if (!hasNote) {
          notes.push({date: this.noteDate, note: this.noteText})
        }
        
        this.$store.dispatch("surveySession/updateSurveySessionState", {name: 'notes', value: notes})
        
        this.noteText = '' 
        this.showNewNote = false
      },
      
      cancelNote() {
        this.noteText = '' 
        this.showNewNote = false
      },
      
      deleteNote(note) {
        var notes = false
        
        if (this.surveyState.notes) {
          notes = JSON.parse(JSON.stringify(this.surveyState.notes))
        }
        
        if (!notes) {
          return
        }
        
        for (var id in notes) {
          if (notes[id].date == note.date) {
            notes.splice(id, 1)
            break
          }
        }
        
        this.$store.dispatch("surveySession/updateSurveySessionState", {name: 'notes', value: notes})
      },
      
      moveNote(note, date) {
        var notes = false
        
        if (this.surveyState.notes) {
          notes = JSON.parse(JSON.stringify(this.surveyState.notes))
        }
        
        if (!notes) {
          return
        }
        
        for (var id in notes) {
          if (notes[id].date == note.date) {
            notes[id].date = date
            break
          }
        }
        
        this.$store.dispatch("surveySession/updateSurveySessionState", {name: 'notes', value: notes})
      },
      
      copyNote(note, date) {
        var notes = false
        
        if (this.surveyState.notes) {
          notes = JSON.parse(JSON.stringify(this.surveyState.notes))
        }
        
        if (!notes) {
          return
        }
        
        notes.push({date: date, note: note.note})
        
        this.$store.dispatch("surveySession/updateSurveySessionState", {name: 'notes', value: notes})
      }
      
    },
  }
</script>

<!-- overwrite that damn limiting CSS from bootstrap/vuetify -->
<style>
  @media (min-width: 500px) {
    .container {
      max-width: 10000px;
    }
  }
  @media only screen and (min-width: 500px) {
    .container {
      max-width: 10000px;
    }
  }
  p {
    margin-bottom: 0;
  }
</style>