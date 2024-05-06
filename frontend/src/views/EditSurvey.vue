<template>
  <v-col cols="12">
    <v-row no-gutters>
      <v-dialog v-model="surveyNameDialog" persistent max-width="600px">
        <v-card>
          <v-card-title>
            <span class="headline">Save Survey</span>
          </v-card-title>
          <v-card-text>
            <v-container grid-list-md>
              <v-layout wrap>
                <v-flex xs12>
                  <v-text-field label="Survey Name" v-model="surveySaveName" :rules="surveyNameRules" required></v-text-field>
                </v-flex>
              </v-layout>
            </v-container>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="surveyNameDialog = false">Close</v-btn>
            <v-btn color="blue darken-1" text @click="surveyNameSave()">Save Survey</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    
      <v-dialog v-model="savingDialog" hide-overlay persistent width="300">
        <v-card color="primary" dark class="pt-3">
          <v-card-text>
            Saving Survey...
            <v-progress-linear
              indeterminate
              color="white"
              class="mb-0"
            ></v-progress-linear>
          </v-card-text>
        </v-card>
      </v-dialog>
    
      <v-card style="width: 100%; height: 100%">
        <v-card-title class="headline font-weight-regular blue-grey white--text"><v-btn fab small text color="white" class="mr-3" to="/admin/survey/list"><v-icon>keyboard_arrow_left</v-icon></v-btn>  {{ headerText }} <div v-if="surveyName" class="ml-2 font-weight-bold"> {{ surveyName }} </div></v-card-title>
        <v-card-text class="pa-5">
          <v-progress-linear :indeterminate="true" class="pb-3" v-if="surveyName && surveyLoadJSON == ''"></v-progress-linear>
          <div class="pb-3" v-if="error">
            <v-alert
              v-model="error"
              type="error"
              transition="scale-transition"
            >{{ errorText }}</v-alert>
            </div>
          <survey-editor :survey-name="surveyName" @save="surveySave" :survey="surveyLoadJSON"></survey-editor>
        </v-card-text>
      </v-card>
    </v-row>
  </v-col>
</template>

<script>  
  import SurveyEditor from '../components/SurveyEditor.vue'

  export default {
    components: {
      SurveyEditor,
    },
    data () {
      return {
        error: false,
        errorText: 'an unknown error occurred',
        
        surveyNameDialog: false,
        savingDialog: false,
        
        surveyName: (this.$route.params.name) ? this.$route.params.name : '',
        surveyNameRules: [
          (v) => !!v || 'Name is required',
        ],
        surveySaveName: '',
        surveyLoadJSON: '',
        surveySaveJSON: '',
      }
    },
    computed: {
      headerText() {
        if (this.surveyName) {
          return 'Edit Survey'
        } else {
          return 'New Survey'
        }
      }
    },
    mounted() {
      this.surveyData()
    },
    methods: {
      surveyData() {
        var that = this
        
        if (this.surveyName) {
          this.axios.get('/survey/' + this.surveyName)
          .then(response => {
            if (response.data.surveyJSON) {
              that.surveyLoadJSON = response.data.surveyJSON
            } else {
              that.errorText = 'Survey data not attached to server response'
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
      surveySave(surveyJSON) {
        this.surveySaveJSON = surveyJSON
        
        if (this.surveyName == '') {
          this.surveyNameDialog = true
        } else {
          this.surveySaveName = this.surveyName
          this.saveSurveyNamed(this.surveySaveName, this.surveySaveJSON)
        }
      },
      surveyNameSave() {
        if (this.surveySaveName != '') {
          
          // TODO check for existing surveys with this name, maybe do it in surveyNameRules..
          
          this.saveSurveyNamed(this.surveySaveName, this.surveySaveJSON)
        }
      },
      saveSurveyNamed(surveyName, surveyJSON) {
        var that = this
        
        this.surveyName = surveyName
        this.surveyLoadJSON = surveyJSON
        this.surveySaveName = surveyName
        this.savingDialog = true
        
        this.axios.post('/survey/' + surveyName, {surveyJSON: surveyJSON})
        .then(() => {
          this.$store.dispatch('surveys/fetchSurveys').then(() => {
            this.surveyNameDialog = false
            
            // use a setTimeout so you can see it did something
            setTimeout(function() {
              that.savingDialog = false
            }, 1000)
            
            // reload the page with the correct path
            if (this.$route.path == '/admin/survey/new') {
              this.$router.push('/admin/survey/edit/' + surveyName)
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
</style>