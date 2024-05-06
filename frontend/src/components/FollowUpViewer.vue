<template>
  <div style="width: 100%">
    <v-card
      color="#385F73"
      dark
    >
      <v-card-text class="white--text">
        <div class="headline">{{ titleText }}</div>
      </v-card-text>
    </v-card>
    
    <v-card class="pa-3 mb-2" v-if="logicValues">
      <v-card-text class="pa-0 ma-0">
        <v-chip outlined color="primary" v-for="(item, index) in logicValues" v-bind:key="index"><strong>{{ index }}: </strong> {{ item }} </v-chip>
      </v-card-text>
    </v-card>
        
    <survey-answers :surveyData="followUpData" mt-2></survey-answers>
    
    <div>
      <v-flex pt-5 v-if="surveyJSON == ''">
        <v-progress-linear :indeterminate="true" ></v-progress-linear>
      </v-flex>
      <survey :survey="survey" v-if="surveyJSON"></survey>
    </div>
  </div>
</template>

<script>
  import SurveyAnswers from '../components/SurveyAnswers.vue'
  import * as SurveyVue from 'survey-vue'
    
  var Survey = SurveyVue.Survey
    
  var cssOverride = {
      navigationButton: "v-btn v-btn--contained theme--light v-size--default primary", // this only seems to append, not replace
  };
  
  import './surveyjs/survey.css';
  
  function doMarkdown(survey, options) {
    var str = options.text;
      if (str.indexOf("<p>") == 0) {
        str = str.substring(3);
          str = str.substring(0, str.length - 4);
      }
      //set html
      options.html = str
  }
  
  export default {
    components: {
      SurveyAnswers,
      Survey,
    },
    data () {
      var tempModel = new SurveyVue.Model({})
      
      return {
        survey: tempModel,
        logicValues: false
      }
    },
    props: {
      surveyConfig: {
        type: Object,
        required: true
      },
      surveyName: {
        type: String,
        required: true
      },
      surveyJSON: {
        type: String,
        required: false
      },
      surveyData: {
        type: Object,
        required: false
      },
      followUpData: {
        type: Object,
        required: false
      }
    },
    watch: {
      surveyJSON() {
        this.loadSurvey()
      },
      surveyData(data) {
        this.survey.data = data
      },
      followUpData() {
        this.buildLogicValues()
      }
    },
    computed: {
      titleText() {
        var date = new Date(this.followUpData.dataKey)
        var formattedDate = date.toUTCString().split(' ').splice(0, 4).join(' ')
        
        return 'Follow up of ' + this.surveyConfig.eventLabel 
          + ' on ' + formattedDate
          + ' (' + this.followUpData.dataKey + ')'
      }
    },
    methods: {
      applyStyle() {
        SurveyVue.StylesManager.applyTheme("winterstone")
  
        var defaultThemeColors = SurveyVue.StylesManager.ThemeColors["default"];
  
        defaultThemeColors["$main-color"] = "#1976d2";
        defaultThemeColors["$main-hover-color"] = "#2196f3";
        defaultThemeColors["$text-color"] = "#000";
        defaultThemeColors["$header-color"] = "#7ff07f";

        defaultThemeColors["$header-background-color"] = "#4a4a4a";
        defaultThemeColors["$body-container-background-color"] = "#fafafa";
  
        SurveyVue.StylesManager.applyTheme(defaultThemeColors);
      },
      buildLogicValues() {
        this.logicValues = false

        if (this.followUpData) {
          var dataKeys = ['created', 'data', 'dataKey', 'modified', 'referenceID', 'sessionDataID', 'sessionID', 'surveyName', 'username']
          for (var dataProperty in this.followUpData) {
            if (!dataKeys.includes(dataProperty)) {
              if (!this.logicValues) {
                this.logicValues = {}
              }
              this.logicValues[dataProperty] = this.followUpData[dataProperty]
            }
          }
        }
      },
      loadSurvey() {
        if (this.surveyJSON) {
          var that = this
          
          var surveyObject = JSON.parse(this.surveyJSON)
          surveyObject.completedHtml = '&nbsp;'
          
          this.survey = new SurveyVue.Model(JSON.stringify(surveyObject))
          this.survey.css = cssOverride
        
          this.survey.onTextMarkdown.add(doMarkdown);
          
          this.survey.onComplete.add(function(sender) {
            that.surveyCompleted(sender.data)
          })
                    
          this.survey.sendResultOnPageNext = true;
          this.survey.onPartialSend.add(function(sender) {
            that.surveyPartialData(sender.data)
          })
          
          this.applyStyle()
          
          if (this.surveyData && Object.keys(this.surveyData).length > 0) {
            this.survey.data = this.surveyData
          }
        }
      },
      surveyPartialData(surveyResults) {
        this.$emit('partialData', {surveyName: this.surveyName, surveyResults: surveyResults})
      },
      surveyCompleted(surveyResults) {
        this.$emit('completed', {surveyName: this.surveyName, surveyResults: surveyResults})
      }
    },
    mounted () {
      if (this.survey) {
        this.loadSurvey()
      }
      this.buildLogicValues()
    }
  }
</script>

<style>
  .svd_commercial_container {
    display: block;
    float: right;
    margin-top: -14px;
  }
  .sv_nav {
    min-height: 40px !important;
  }
</style>