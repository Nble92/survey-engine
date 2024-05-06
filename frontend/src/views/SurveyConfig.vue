<template>
  <v-col cols="12">
    <v-row no-gutters>
      <div class="pb-4" v-if="error" style="width: 100%">
        <v-alert
          v-model="error"
          type="error"
          transition="scale-transition"
        >{{ errorText }}</v-alert>
      </div>
      
      <v-dialog v-model="savingDialog" hide-overlay persistent width="300">
        <v-card color="primary" dark class="pt-3">
          <v-card-text>
            Saving Settings...
            <v-progress-linear
              indeterminate
              color="white"
              class="mb-0"
            ></v-progress-linear>
          </v-card-text>
        </v-card>
      </v-dialog>
      
      <v-card style="width: 100%">
        <v-card-title class="headline font-weight-regular blue-grey white--text"><v-icon dark class="mr-2">ballot</v-icon> Configuration</v-card-title>
        <v-card-text class="pa-5">
         <v-form v-model="valid" ref="form" @submit.prevent="">
            <v-select
              :items="surveyTypes"
              v-model="currentSurveyType"
              item-text="label"
              item-value="type"
              label="Master Survey Type"
              required
            ></v-select>
          
            <linear-survey-options v-if="currentSurveyType == 'linear'" :config="currentConfig" v-model="linearConfig"></linear-survey-options>
            <calendar-survey-options v-if="currentSurveyType == 'calendar'" :config="currentConfig" v-model="calendarConfig"></calendar-survey-options>
          
            <v-layout justify-end pt-3>
              <v-btn @click="updateSettings" color="primary">Save Config</v-btn>
            </v-layout>
          </v-form>
        </v-card-text>
      </v-card>
    </v-row>
  </v-col>
</template>

<script>
  import LinearSurveyOptions from '../components/LinearSurveyOptions.vue'
  import CalendarSurveyOptions from '../components/CalendarSurveyOptions.vue'
  
  import { mapState } from 'vuex'
  
  
  export default {
    components: {
      LinearSurveyOptions,
      CalendarSurveyOptions,
    },
    data () {
      return {
        valid: false,
        error: false,
        errorText: 'an unknown error occurred',
        
        savingDialog: false,
        
        surveyListHeaders: [
          { text: 'Survey Name', value: 'name' },
          { text: 'Last Modified', value: 'modified'},
          { text: 'Author', value: 'author', class: 'hidden-md-and-down'},
          { text: 'Actions', value: 'name', sortable: false, align: 'center' }
        ],
        
        currentSurveyType: '',
        currentConfig: {},
        
        // per type configs
        linearConfig: {},
        calendarConfig: {},
        
        // deletion stuff
        surveyDeleteDialog: false,
        deleteSurveyName: '',
      }
    },
    computed: {
      ...mapState('surveySettings', ['surveyType', 'config', 'surveyTypes']),
      ...mapState('surveys', ['surveys']),
    },
    watch: {
      surveyType() {
        this.currentSurveyType = this.surveyType
        this.currentConfig = this.config
      }
    },
    mounted() {
      this.$store.dispatch('logicFunctions/fetchLogicFunctions')
      
      this.currentSurveyType = this.surveyType
      this.currentConfig = this.config
    },
    methods: {
      updateSettings() {
        var config = {notDefined: true}
        var that = this
        
        this.savingDialog = true
        
        if (this.currentSurveyType == 'linear') {
          config = this.linearConfig
        } else if (this.currentSurveyType == 'calendar') {
          config = this.calendarConfig
        }
                
        this.$store.dispatch('surveySettings/update', {
          surveyType: this.currentSurveyType, 
          config: config
        }).then(() => {
          // use a setTimeout so you can see it did something
          setTimeout(function() {
            that.savingDialog = false
          }, 1000)
        })
      },
      editSurvey(name) {
        this.$router.push('/admin/survey/edit/' + name)
      },
      deleteSurvey(name) {
        this.deleteSurveyName = name
        this.surveyDeleteDialog = true
      },
      deleteSurveyConfirm() {
        var that = this
        
        this.axios.delete('/survey/' + this.deleteSurveyName)
        .then(() => {
          this.$store.dispatch('surveys/fetchSurveys').catch(function (error) {
            if (error.response && error.response.data) {
              that.errorText = error.response.data.message
            } else {
              that.errorText = error.message
            }
            that.error = true
          })
        }).catch(function (error) {
            if (error.response && error.response.data) {
              that.errorText = error.response.data.message
            } else {
              that.errorText = error.message
            }
            that.error = true
        })
        
        this.deleteSurveyName = ''
        this.surveyDeleteDialog = false
      }
    },
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>