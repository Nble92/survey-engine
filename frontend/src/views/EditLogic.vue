<template>
  <v-col cols="12">
    <v-row no-gutters>   
      <v-dialog v-model="functionNameDialog" persistent max-width="600px">
        <v-card>
          <v-card-title>
            <span class="headline">Save Function</span>
          </v-card-title>
          <v-card-text>
            <v-container grid-list-md>
              <v-layout wrap>
                <v-flex xs12>
                  <v-text-field label="Function Name" v-model="functionSaveName" :rules="functionNameRules" required></v-text-field>
                </v-flex>
              </v-layout>
            </v-container>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="functionNameDialog = false">Close</v-btn>
            <v-btn color="blue darken-1" text @click="functionNameSave">Save</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    
      <v-dialog v-model="showValidationOutput">
        <v-card color="black" dark class="pt-4">
          <v-card-text>
            <h4>Validation result:</h4>
            <pre class="pa-2">{{ validationOutput }}</pre>
          </v-card-text>
        </v-card>
      </v-dialog>
      
      <v-dialog v-model="error" width="500px">
        <v-card color="error" dark class="pt-3">
          <v-card-text>
            <h4>Logic Error</h4>
            <p class="pt-3">{{ errorText }}</p>
          </v-card-text>
        </v-card>
      </v-dialog>
      
      <v-dialog v-model="savingDialog" hide-overlay persistent width="300">
        <v-card color="primary" dark class="pt-3">
          <v-card-text>
            Saving Function...
            <v-progress-linear
              indeterminate
              color="white"
              class="mb-0"
            ></v-progress-linear>
          </v-card-text>
        </v-card>
      </v-dialog>
    
      <v-flex mb-5 align-start>
        <v-card xs12>
          <v-card-title class="headline font-weight-regular blue-grey white--text"><v-btn fab small text color="white" class="mr-3" to="/admin/logic/list"><v-icon>keyboard_arrow_left</v-icon></v-btn>  Edit Function</v-card-title>
          <v-card-text class="pa-3">
            <p class="mb-4">
              You can create javascript functions to use to validate data and collate results to then use in follow-up surveys.
              Functions are passed the entire survey session state and can either return a boolean or an array of filtered survey sessions.
              The function name is specified at save, please only enter the code for within the function, and not the function definition itself.
              <b class="red--text">You must validate before you can save.</b>
            </p>
            <v-form>
              <v-layout justify-center mb-4>
                <v-flex md11>
                  <v-toolbar class="blue darken-2">
                    <v-toolbar-title class="white--text"><h4>{{ headerText }}</h4></v-toolbar-title>
                  </v-toolbar>
                  <v-card>
                    <v-card-text class="pa-4">
                      <v-textarea
                        class="monospace"
                        label="Function Javascript"
                        v-model="functionJS"
                        :error="error"
                        @input="functionChange"
                        @keydown.tab.prevent="tabber($event)"
                        auto-grow
                      ></v-textarea>
                    </v-card-text>
                  </v-card>
                  <v-toolbar class="blue darken-2 mb-4">
                    <v-toolbar-title class="white--text"><h4>}</h4></v-toolbar-title>
                  </v-toolbar>
                
                  <v-layout justify-end>
                    <v-btn @click="functionValidate" color="primary" outlined class="ma-2" dark>Validate</v-btn>
                    <v-btn @click="functionSave" class="ma-2" color="primary" :class="{ 'blue-grey white--text' : valid, disabled: !valid }" :disabled="!valid">Save</v-btn>
                  </v-layout>
                
                  <h4>Example Function:</h4>
                  <pre class="mb-3 pa-3 monospace">for (var id in surveySessionData.sessionData) {<br>  var sessionData = surveySessionData.sessionData[id]<br>  console.log(sessionData)<br>}<br><br>return true</pre>
                  <div v-if="hasSurveySessionData">
                    <h4>Your current surveySessionData:</h4>
                    <pre class="mb-3 pa-3 monospace">{{ surveySessionDataString }}</pre>
                  </div>
                </v-flex>
              </v-layout>
            </v-form>
          </v-card-text>
        </v-card>
      </v-flex>
    
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
        
        hasSurveySessionData: false,
        surveySessionData: '',
        
        showValidationOutput: false,
        validationOutput: '',
        
        functionNameDialog: false,
        savingDialog: false,
        
        functionName: (this.$route.params.name) ? this.$route.params.name : '',
        functionNameRules: [
          (v) => !!v || 'Name is required',
        ],
        
        functionJS: '',
        
        functionSaveName: '',
        functionSaveJS: '',
      }
    },
    computed: {
      headerText() {
        return (this.functionName) ? this.functionName + '(surveySessionData) {' : 'newFunction(surveySessionData){'
      },
      surveySessionDataString() {
        return JSON.stringify(this.surveySessionData, null, 2)
      },
      ...mapState('surveySettings', ['surveyType', 'config', 'surveyTypes']),
      ...mapState('logicFunctions', ['logicFunctions']),
      ...mapState('surveySession', ['status', 'sessionID', 'surveyState']),
    },
    mounted() {
      var that = this
      this.logicFunction()
      
      this.$store.dispatch('surveySession/fetchSessionState').then(() => {
        this.axios.get('/survey/sessionDump/' + this.sessionID)
        .then(response => { 
          that.surveySessionData = response.data
          that.hasSurveySessionData = true
        }).catch(function () {
          // do nothing
        })
      })
    },
    methods: {
      tabber (event) {
        let text = this.functionJS,
            originalSelectionStart = event.target.selectionStart,
            textStart = text.slice(0, originalSelectionStart),
            textEnd =  text.slice(originalSelectionStart);

        this.functionJS = `${textStart}  ${textEnd}` // muahahaha.. oh well. 2 spaces instead of a tab makes monospace better
        event.target.value = this.functionJS // required to make the cursor stay in place.
        event.target.selectionEnd = event.target.selectionStart = originalSelectionStart + 2
      },
      logicFunction() {
        var that = this
        
        if (this.functionName) {
          this.axios.get('/logic/' + this.functionName)
          .then(response => {
            if (response.data.functionJS) {
              that.functionJS = response.data.functionJS
            } else {
              that.errorText = 'Function data not returned'
              that.error = true
            }
          }).catch(function (error) {
            if (error.response && error.response.data) {
              that.errorText = error.response.data.message
            } else {
              that.errorText = error.message
            }
            that.error = true
          })
        }
      },
      functionChange() {
        this.valid = false
      },
      functionValidate() {
        var scriptJS = 'function test(surveySessionData) {' + this.functionJS + '} test(' + this.surveySessionDataString + ')'
        
        try {
          this.validationOutput = eval(scriptJS)
          
          if (!this.validationOutput) {
            this.validationOutput = 'Undefined'
          }
          
          this.valid = true
          this.error = false
          
          this.showValidationOutput = true
          
          // eslint-disable-next-line
          console.log('validation output: ', this.validationOutput)
        }
        catch (err) {
          this.valid = false
          this.showValidationOutput = false
          
          this.error = true
          this.errorText = 'Validation error: ' + err
          
          // eslint-disable-next-line
          console.log('validate error: ', err)
        }
      },
      functionSave() {
        this.functionSaveJS = this.functionJS
        
        if (this.functionName == '') {
          this.functionNameDialog = true
        } else {
          this.functionSaveName = this.functionName
          this.saveFunctionNamed(this.functionSaveName, this.functionSaveJS)
        }
      },
      functionNameSave() {
        if (this.functionSaveName != '') {
          
          // TODO check for existing surveys with this name, maybe do it in surveyNameRules..
          
          this.saveFunctionNamed(this.functionSaveName, this.functionSaveJS)
        }
      },
      saveFunctionNamed(functionName, functionJS) {
        if (functionJS == '') {
          return
        }
        
        var that = this
        
        this.functionName = functionName
        this.functionJS = functionJS
        this.functionSaveName = functionName
        this.savingDialog = true
        
        this.axios.post('/logic/' + functionName, {functionJS: functionJS})
        .then(() => {
          this.$store.dispatch('logicFunctions/fetchLogicFunctions').then(() => {
            this.functionNameDialog = false
            
            // use a setTimeout so you can see it did something
            setTimeout(function() {
              that.savingDialog = false
            }, 1000)
            
            // reload the page with the correct path
            if (this.$route.path == '/admin/logic/new') {
              this.$router.push('/admin/logic/edit/' + functionName)
            }
            
            
          }).catch(function (error) {
            if (error.response && error.response.data) {
              that.errorText = error.response.data.message
            } else {
              that.errorText = error.message
            }
            that.error = true
            this.savingDialog = false
            
          })
        }).catch(function (error) {
            if (error.response && error.response.data) {
              that.errorText = error.response.data.message
            } else {
              that.errorText = error.message
            }
            that.error = true
            
            this.savingDialog = false
            
        })
      }
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  .monospace {
    font-family: 'DejaVu Sans Mono', monospace;
  }
</style>