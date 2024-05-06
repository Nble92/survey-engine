<template>
  <div>
    <v-flex pt-5 v-if="surveyJSON == ''">
      <v-progress-linear :indeterminate="true" ></v-progress-linear>
    </v-flex>
    <survey :survey="survey" v-if="surveyJSON"></survey>
  </div>
</template>

<script>
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
      Survey
    },
    data () {
      var tempModel = new SurveyVue.Model({})
      
      return {
        survey: tempModel
      }
    },
    props: {
      readOnly: {
        type: Boolean,
        required: false,
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
      }
    },
    watch: {
      readOnly(isReadOnly) {
        if (isReadOnly) {
          this.survey.mode = "display"
        }
      },
      surveyJSON() {
        this.loadSurvey()
      },
      surveyData(data) {
        this.survey.data = data
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
  .sv_body img {
    max-width: 100%;
  }
  .sv_main .sv_container .sv_body .sv_p_root .sv_q .sv_q_other input, .sv_main .sv_container .sv_body .sv_p_root .sv_q .sv_q_dropdown_control, .sv_main .sv_container .sv_body .sv_p_root .sv_q input:not([type="button"]):not([type="reset"]):not([type="submit"]):not([type="image"]):not([type="checkbox"]):not([type="radio"]), .sv_main .sv_container .sv_body .sv_p_root .sv_q select {
    height: unset;
  }
</style>