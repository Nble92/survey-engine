<template>  
  <v-container>
    <v-col cols="12">
      <v-row no-gutters>
        <v-card style="width: 100%">
          <v-card-title class="headline font-weight-regular primary white--text">General Survey Settings</v-card-title>
          <v-card-text class="pa-5">
            <v-form>
              <v-checkbox
                label="Allow users to run this survey multiple times."
                v-model="userMultipleRuns"
              ></v-checkbox>

              <v-checkbox
                label="Allow guests to run this survey."
                v-model="allowGuests"
              ></v-checkbox>              
            </v-form>
          </v-card-text>
        </v-card>
        
        <v-card style="width: 100%">
          <v-card-title class="headline font-weight-regular primary white--text">Calendar Based Survey Settings</v-card-title>
          <v-card-text class="pa-5">
            <v-form>
              <h3>Step 1</h3>
              <div style="text-align: justify; text-justify: inner-word" class="blue-grey--text text--darken-3">
                This survey is presented when you start a new survey flow.  <i>Generally this survey would include identification information.</i>
              </div>
            
              <v-flex xs12 ma-3>
                <v-select
                  :items="surveys"
                  v-model="entrySurvey"
                  item-text="name"
                  item-value="name"
                  label="Entry Survey (identification information)"
                  required
                ></v-select>
                <p>Please give some basic instructions to be displayed before the user starts the survey.</p>
                <rich-editor 
                  v-model="startSurveyInstructions"  
                  label="Start survey instructions."
                  id="startSurveyInstructions"
                />
              </v-flex>
            
              <h3>Step 2</h3>
              <div style="text-align: justify; text-justify: inner-word" class="blue-grey--text text--darken-3">
                The calendar is shown.
              </div>
            
              <div style="margin-top: 40px; margin-bottom: 40px">
                <v-layout ma-4 row wrap>
                  <v-flex md6>
                    <v-list-item>
                      <v-icon style="font-size: 80px" class="mr-3" color="blue-grey">calendar_today</v-icon>
                      <v-list-item-content>
                        <v-list-item-title>Calendar Component</v-list-item-title>
                        <v-list-item-subtitle>Calendar showing the days of the month.</v-list-item-subtitle>
                      </v-list-item-content>
                    </v-list-item>
                  </v-flex>
                  <v-flex md6 class="text-sm-center">
                    <v-dialog v-model="showCalendarSettings" persistent max-width="600px">
                      <template v-slot:activator="{ on }">
                        <v-btn color="primary" dark v-on="on">Calendar Settings</v-btn>
                      </template>
                      <v-card>
                        <v-card-title>
                          <span class="headline">Calendar Settings</span>
                        </v-card-title>
                        <v-card-text>
                          <v-flex xs12 ma-3>
                            <v-checkbox
                              label="Highlight N days in the past"
                              v-model="calendarHighlightPriorDays"
                            ></v-checkbox>
                            <v-text-field
                              v-if="calendarHighlightPriorDays"
                              label="Number of days"
                              v-model="calendarHighlightPriorDaysValue"
                            ></v-text-field>
                          </v-flex>
                          
                        </v-card-text>
                        <v-card-actions>
                          <v-spacer></v-spacer>
                          <v-btn color="blue darken-1" text @click="showCalendarSettings = false">Close</v-btn>
                          <v-btn color="blue darken-1" text @click="showCalendarSettings = false">Save</v-btn>
                        </v-card-actions>
                      </v-card>
                    </v-dialog>
                  </v-flex>
                </v-layout>
              </div>
            
              <h3>Step 3</h3>
              <div style="text-align: justify; text-justify: inner-word" class="blue-grey--text text--darken-3">
                The user clicks on a calendar day, or existing event.
              </div>
            
              <v-flex xs12 ma-3>
                <v-checkbox
                  label="Allow multiple events per day."
                  v-model="allowMultipleEvents"
                ></v-checkbox>
                <v-text-field
                  label="Event Label"
                  v-model="eventLabel"
                ></v-text-field>
                
                <v-checkbox
                  label="Use function to determine event label"
                  v-model="eventLabelUseFunction"
                ></v-checkbox>
                
                <v-select
                  :items="logicFunctions"
                  v-model="eventLabelFunction"
                  item-text="functionName"
                  item-value="functionName"
                  label="Event Label Function"
                  v-if="eventLabelUseFunction"
                ></v-select>
                
                <v-select
                  :items="surveys"
                  v-model="eventSurvey"
                  item-text="name"
                  item-value="name"
                  label="Event Survey (when you click a day)"
                  required
                ></v-select>
              </v-flex>
            
              <h3>Step 4</h3>
              <div style="text-align: justify; text-justify: inner-word" class="blue-grey--text text--darken-3">
                Follow-up surveys can run custom logic to determine whether particular events require extra data. Please enable this, and then select the survey you wish to use for follow up questioning.
              </div>
              <v-flex xs12 ma-3>
                <v-checkbox
                  label="Enable follow up questioning."
                  v-model="allowFollowUp"
                ></v-checkbox>
                <div v-if="allowFollowUp">
                  <p>The follow-up function determines surveys that require further questioning based off the Follow-up Survey. This function must return an array of the sessionDataIDs you wish to reference.</p>
                  <v-select
                    :items="logicFunctions"
                    v-model="followUpFunction"
                    item-text="functionName"
                    item-value="functionName"
                    label="Follow-up Function"
                    required
                  ></v-select>
              
                  <v-text-field
                    v-model="followUpButtonLabel"
                    label="Follow-up Button Label"
                    required
                  ></v-text-field>
              
                  <v-select
                    :items="surveys"
                    v-model="followUpSurvey"
                    item-text="name"
                    item-value="name"
                    label="Follow-up Survey"
                    required
                  ></v-select>
                </div>
              </v-flex>
            
              <h3>Step 5</h3>
              <div style="text-align: justify; text-justify: inner-word" class="blue-grey--text text--darken-3">
                Completion surveys are used when you click complete survey. Please enable this, and then select the survey you wish to use.
              </div>
              <v-flex xs12 ma-3>
                <v-checkbox
                  label="Enable completion questioning."
                  v-model="allowCompletionSurvey"
                ></v-checkbox>
                <div v-if="allowCompletionSurvey">
                  <p>The completion function to determine when the completion flow can be presented. This function must return an boolean.</p>
                  <v-select
                    :items="logicFunctions"
                    v-model="completionFunction"
                    item-text="functionName"
                    item-value="functionName"
                    label="Completion Function"
                    required
                  ></v-select>
              
                  <v-text-field
                    v-model="completionButtonLabel"
                    label="Completion Button Label"
                    required
                  ></v-text-field>
              
                  <v-select
                    :items="surveys"
                    v-model="completionSurvey"
                    item-text="name"
                    item-value="name"
                    label="Completion Survey"
                    required
                  ></v-select>
                </div>
                <v-checkbox
                  :label="notifyOnCompleteLabel"
                  v-model="notifyOnComplete"
                ></v-checkbox>
              </v-flex>
              
              
              <h3>Step 6</h3>
              <div style="text-align: justify; text-justify: inner-word" class="blue-grey--text text--darken-3">
                Once all surveys have been completed present the user with a button to click that takes them to an external site. The button says "<strong>Continue to external site</strong>".
              </div>
              <v-flex xs12 ma-3>
                <p>Please give some basic instructions to be displayed when the user finishes the survey.</p>
                <rich-editor 
                  v-model="completedSurveyInstructions"  
                  label="Completed survey instructions."
                  id="completedSurveyInstructions"
                />
                <v-checkbox
                  label="Enable redirect button."
                  v-model="allowRedirectButton"
                ></v-checkbox>
                <div v-if="allowRedirectButton">
                  <p>Please give some basic instructions to be displayed with the redirect button.</p>
                
                  <rich-editor 
                    v-model="redirectButtonInstructions"  
                    label="Redirect button instructions."
                    id="redirectButtonInstructions"
                  />
                
                  <v-text-field
                    v-model="redirectURL"
                    label="Redirect URL"
                    required
                  ></v-text-field>
                </div>
              </v-flex>
              
            </v-form>
          </v-card-text>
        </v-card>
      </v-row>
    </v-col>
  </v-container>
</template>

<script>  
  import { mapState } from 'vuex'
  import RichEditor from '../components/RichEditor.vue'  
  
  export default {
    components: { RichEditor },
    props: {
      config: {
        type: Object,
        required: false
      },
    },
    data () {
      return {
        valid: true,
        
        allowGuests: false,
        userMultipleRuns: false,
        
        startSurveyInstructions: 'Please click the Start Survey button below',
        entrySurvey: '',
        
        showCalendarSettings: false,
        calendarHighlightPriorDays: false,
        calendarHighlightPriorDaysValue: 28,
        
        eventSurvey: '',
        allowMultipleEvents: false,
        eventLabel: '',
        eventLabelUseFunction: false,
        eventLabelFunction: '',
        
        allowFollowUp: false,
        followUpFunction: '',
        followUpButtonLabel: 'Continue to follow-up survey(s)',
        followUpSurvey: '',
        
        completedSurveyInstructions: 'Thanks for completing the survey',
        allowCompletionSurvey: false,
        completionFunction: '',
        completionButtonLabel: 'Continue to completion survey',
        completionSurvey: '',
        notifyOnComplete: false,

        allowRedirectButton: false,
        redirectButtonInstructions: '',
        redirectURL: ''
      }
    },
    watch: {
      $data: {
        handler: function() {
          this.$emit('input', this.generateConfig())
        },
        deep: true
      }
    },
    computed: {
      notifyOnCompleteLabel() {
        return 'Email notification to ' + this.contactName + ' on completion'
      },
      ...mapState('config', ['contactName', 'contactEmail']),
      ...mapState('surveys', ['surveys']),
      ...mapState('logicFunctions', ['logicFunctions']),
    },
    methods: {
      generateConfig() {
        return this.$data
      }
    },
    mounted () {
      if (this.config) {
        for (var key in this.config) {
          if (this.hasOwnProperty(key)) {
            this[key] = this.config[key]
          }
        }
      }
    }
  }
</script>

<style scoped>

</style>