<template>
  <div>
    <div v-if="!hidden">
      <v-card style="width: 100%" v-if="answers.length > 0" flat>
        <v-card-title>{{ surveyData.surveyName }}</v-card-title>
        <v-card-text class="pl-5">
          <div v-for="answer in answers" v-bind:key="answer.name">
            <h3 v-html="answer.title"></h3>
            <div v-if="answer.type == 'single'">
              <p>{{ answer.text }}</p>
            </div>
            <!-- eslint-disable -->
            <div class="pa-3" v-else-if="answer.type == 'multiple'">
              <div v-for="(row, index) in answer.values" v-bind:key="`row-${index}`">
                <div v-if="column" style="display: inline-block; white-space: nowrap;" class="ma-1" v-for="column in row" v-bind:key="column.title">
                  <b v-html="column.title"></b>: {{ column.text }} 
                </div> 
              </div>
            </div>
            <!-- eslint-disable -->
          </div>
        </v-card-text>
      </v-card>
      <v-card style="width: 100%" flat v-else-if="parsed == true">
        <v-card-text class="pl-5 text-center">
          <h3 class="red--text text--lighten-2" ma-0>No Answers</h3>
        </v-card-text>
      </v-card>
    </div>
  </div>
</template>

<script>
  import * as answerHelper from '../helpers/answers.js'
  
  export default {
    data () {
      return {
        error: false,
        errorText: '',
        
        loaded: false,
        parsed: false,
        
        surveyJSON: {},
        
        answers: [],
      }
    },
    props: {
      hidden: {
        type: Boolean,
        required: false
      },
      surveyData: {
        type: Object,
        required: false
      }
    },
    watch: {
      surveyData() {
        this.loadSurvey()
      }
    },
    methods: {
      parseAnswers() {
        var results = answerHelper.parseAnswers(this.surveyJSON, this.surveyData)
        this.answers = results[0]
        this.parsed = true
      },
      
      loadSurvey() {
        var that = this
        
        if (this.surveyData) {
          if (this.surveyData.surveyName && this.surveyData.surveyName == this.loaded) {
            this.parseAnswers()
          } else if (this.surveyData.surveyName) {
            this.axios.get('/survey/' + this.surveyData.surveyName)
            .then(response => {
              if (response.data.surveyJSON) {
                that.surveyJSON = JSON.parse(response.data.surveyJSON)
                that.loaded = response.data.name
                that.parseAnswers()
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
          
        }
      },
    },
    mounted () {
      if (this.surveyData) {
        this.loadSurvey()
      }
    }
  }
</script>

<style>
  b p {
    display: inline-block;
  }
</style>