<template>
  <div id="surveyEditorContainer"></div>
</template>

<script>
  import Vue from 'vue'
  import RichEditor from '../components/RichEditor.vue'
  
  import * as SurveyVue from 'survey-vue'
  import * as SurveyCreator from 'survey-creator'
  
  import './surveyjs/survey-creator.css';
  import './surveyjs/fixes.css';
  
  function doMarkdown(survey, options) {
    var str = options.text;
      if (str.indexOf("<p>") == 0) {
        str = str.substring(3);
          str = str.substring(0, str.length - 4);
      }
      //set html
      options.html = str
  }

  var RichText_ModalEditor = {		
    afterRender: function (modalEditor, htmlElement) {
      var ComponentClass = Vue.extend(RichEditor)
      var editor = new ComponentClass({propsData: {id: htmlElement.id, value: modalEditor.editingValue}})
      
      editor.$on('input', function (newValue) {
        modalEditor.editingValue = newValue;
      });
      
      modalEditor.onValueUpdated = function(newValue) {
        editor.value = newValue;		
      }
      
      editor.$mount('#' + htmlElement.id)
    }	
  };
 

  export default {
    data() {
      return {}
    },
    props: {
      surveyName: {
        type: String,
        required: true
      },
      survey: {
        type: String,
        required: false
      }
    },
    watch: {
      survey(newVal) {
        this.editor.text = newVal
        this.applyStyle()
      }
    },
    methods: {
      applyStyle() {
        SurveyVue.StylesManager.applyTheme("winterstone")
        SurveyCreator.StylesManager.applyTheme("winterstone")
        
        var defaultThemeColors = SurveyVue.StylesManager.ThemeColors["default"];
        defaultThemeColors["$main-color"] = "#1976d2";
        defaultThemeColors["$main-hover-color"] = "#2196f3";
        defaultThemeColors["$text-color"] = "#000";
        defaultThemeColors["$header-color"] = "#7ff07f";

        defaultThemeColors["$header-background-color"] = "#4a4a4a";
        defaultThemeColors["$body-container-background-color"] = "#fafafa";
        
        var defaultThemeColorsEditor = SurveyCreator.StylesManager.ThemeColors["default"];
        defaultThemeColorsEditor["$primary-color"] =  "#1976d2";
        defaultThemeColorsEditor["$secondary-color"] =  "#1976d2";
        defaultThemeColorsEditor["$primary-hover-color"] = "#2196f3";
        defaultThemeColorsEditor["$primary-text-color"] = "#000";
        defaultThemeColorsEditor["$selection-border-color"] = "#1976d2";
        
        SurveyVue.StylesManager.applyTheme(defaultThemeColors);
        SurveyCreator.StylesManager.applyTheme(defaultThemeColorsEditor);
      },
    },
    mounted () {
      var that = this
      
      let editorOptions = { 
        showEmbededSurveyTab: false,
        showJSONEditorTab: true,
        showTestSurveyTab: true,
        questionTypes : ["html", "text", "checkbox", "radiogroup", "dropdown", "comment", "rating", "boolean", "matrix", "matrixdropdown", "matrixdynamic", "multipletext", "imagepicker"]
      };
            
      SurveyCreator.SurveyPropertyModalEditor.registerCustomWidget("html", RichText_ModalEditor);
      SurveyCreator.SurveyPropertyModalEditor.registerCustomWidget("text", RichText_ModalEditor);
      
      this.editor = new SurveyCreator.SurveyCreator('surveyEditorContainer', editorOptions);
      this.editor.text = this.survey
      
      this.editor.survey.onTextMarkdown.add(doMarkdown);
      
      this.editor.onDesignerSurveyCreated.add(function (editor, options) {
        options.survey.onTextMarkdown.add(doMarkdown);
      });
      
      this.editor.onSurveyInstanceCreated.add(function (editor, options) {
        options.survey.onTextMarkdown.add(doMarkdown);
      });
      
      this.editor.onTestSurveyCreated.add(function (editor, options) {
        options.survey.onTextMarkdown.add(doMarkdown);
      });
      
      this.editor.saveSurveyFunc = function() {
        that.$emit('save', this.text)
      };
      
      this.applyStyle()
    }
  }
</script>

<style>
  .svd_commercial_container {
    display: block;
    float: right;
    margin-top: -14px;
  }
  .svd_container #testDeviceSelector {
    margin-top: -6px;
    margin-bottom: 6px;
    margin-left: 3px;
  }
  .sv_body img {
    max-width: 100%;
  }
  .sv_main .sv_container .sv_body .sv_p_root .sv_q .sv_q_other input, .sv_main .sv_container .sv_body .sv_p_root .sv_q .sv_q_dropdown_control, .sv_main .sv_container .sv_body .sv_p_root .sv_q input:not([type="button"]):not([type="reset"]):not([type="submit"]):not([type="image"]):not([type="checkbox"]):not([type="radio"]), .sv_main .sv_container .sv_body .sv_p_root .sv_q select {
    height: unset;
  }
</style>