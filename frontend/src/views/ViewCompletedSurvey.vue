<template>
  <v-col cols="12">
    <v-row no-gutters>
      <div class="pb-4" v-if="error">
        <v-alert
          v-model="error"
          type="error"
          transition="scale-transition"
        >{{ errorText }}</v-alert>
      </div>

      <v-card style="width: 100%">
        <v-toolbar class="blue-grey white--text">
          <v-btn fab small text color="white" class="mr-3" to="/survey/list"><v-icon>keyboard_arrow_left</v-icon></v-btn> 
          <v-toolbar-title class="headline font-weight-regular">
            Completed <span v-if="isTest" >Test</span> Survey {{ sessionID }}
          </v-toolbar-title>
          <v-spacer></v-spacer>
          <v-btn v-if="isAdmin" @click="editSurvey">Edit Survey</v-btn>
        </v-toolbar>
        
        <div class="text-center pa-5" v-if="hasNoData">
          <h2 class="grey--text" ma-0>Survey session has no data</h2>
        </div>
        
        <v-card class="ma-5" v-for="surveyData in surveySessionData" :key="surveyData.sessionDataID">
          <v-card-text>
            <div class="mb-5">
              <strong>Data ID:</strong> {{ surveyData.sessionDataID }}<br/>
              <div v-if="surveyData.dataKey"><strong>Data Key:</strong> {{ surveyData.dataKey }}</div>
              <div v-if="surveyData.referenceID"><strong>Data Reference ID:</strong> {{ surveyData.referenceID }}</div>
              <strong>Created:</strong> {{ surveyData.created }}<br/>
              <strong>Modified:</strong> {{ surveyData.modified }}<br/>
              <strong>Username:</strong> {{ surveyData.username }}
            </div>
            
            <survey-answers :surveyData="surveyData" v-if="surveyData.data"></survey-answers>
          </v-card-text>
        </v-card>
      </v-card>
    </v-row>
  </v-col>
</template>


<script>
  import { mapGetters } from 'vuex'
  import SurveyAnswers from '../components/SurveyAnswers.vue'

  export default {
    components: {
      SurveyAnswers,
    },
    data () {
      return {
        valid: false,
        error: false,
        errorText: 'an unknown error occurred',
        
        sessionID: (this.$route.params.id) ? this.$route.params.id : '',
        surveyType: '',
        isTest: false,
        
        hasNoData: false,
        
        surveySessionData: {}
      }
    },
    computed: {
      ...mapGetters('user', ['isAdmin']),
    },
    mounted() {
      var that = this
      
      if (!this.sessionID) {
        return
      }
      
      this.axios.get('/survey/session/' + this.sessionID)
      .then(response => { 
        that.surveyType = response.data.surveyType
        that.isTest = response.data.isTest
        
        
        this.axios.get('/survey/sessionDump/' + this.sessionID + "?all=true")
        .then(response => { 
          if (response.data.sessionData) {
            that.hasNoData = false
            that.surveySessionData = response.data.sessionData
        
            // eslint-disable-next-line
            console.log('surveySessionData', that.surveySessionData)
          } else {
            that.hasNoData = true
          }
        
        }).catch(function (error) {
          that.error = true
          if (error.response && error.response.data) {
            that.errorText = error.response.data
          } else if (error.message) {
            that.errorText = error.message
          } else {
            that.errorText = 'Error fetching survey session data...'
          }
        })
        
      }).catch(function (error) {
        that.error = true
        if (error.response && error.response.data) {
          that.errorText = error.response.data
        } else if (error.message) {
          that.errorText = error.message
        } else {
          that.errorText = 'Error fetching survey data...'
        }
      })
    },
    methods: {
      editSurvey() {
        this.$store.dispatch('surveySession/editSurveySession', this.sessionID).then(() => {
          this.$router.push('/survey/run')
        })
      }
    },
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>