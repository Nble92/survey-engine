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
          <v-card-title class="headline font-weight-regular primary white--text">Single Linear Survey Settings</v-card-title>
          <v-card-text class="pa-5">
            <h3>Step 1</h3>
            <div style="text-align: justify; text-justify: inner-word" class="mb-4">
              A single linear survey is the simplest type provided by this survey engine, it presents a single survey and records the submission in a unique survey session.
            </div>
            <v-form>
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
            </v-form>
          
            <h3>Step 2</h3>
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
                  label="Completion Logic Function"
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
            
            <h3>Step 3</h3>
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
        
        allowFollowUp: false,
        followUpFunction: '',
        followUpSurvey: '',
        
        completedSurveyInstructions: 'Thanks for completing the survey',
        allowCompletionSurvey: false,
        completionFunction: '',
        completionButtonLabel: 'Complete Survey',
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