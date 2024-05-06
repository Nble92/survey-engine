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
            Saving Instructions...
            <v-progress-linear
              indeterminate
              color="white"
              class="mb-0"
            ></v-progress-linear>
          </v-card-text>
        </v-card>
      </v-dialog>
    
      <v-flex mb-5 align-start>
        <v-card style="width: 100%">
          <v-card-title class="headline font-weight-regular blue-grey white--text"><v-btn fab small text color="white" class="mr-3" to="/admin/instructions/list"><v-icon>keyboard_arrow_left</v-icon></v-btn>  Instructions Editor</v-card-title>
          <v-card-text>
            <v-form>
              <v-card flat style="width: 100%">
                <v-card-title>Survey View</v-card-title>
                <v-card-text>
                  <v-select
                    :items="availableSurveyList"
                    v-model="surveyName"
                    item-text="name"
                    item-value="name"
                    label="Instructions for survey"
                    required
                    :disabled="existingInstruction"
                  ></v-select>
              
                  <v-checkbox
                    label="Open window for full instructions"
                    v-model="popoverWindow"
                  ></v-checkbox>
                </v-card-text>
                
                <v-card-title>Instruction Summary</v-card-title>
                <v-card-text>
                  <div style="text-align: justify; text-justify: inner-word" class="mb-4">
                    Please provide a brief <i>(no more than 2 lines)</i> of instruction summary to display at the top of the survey. A button will provide a popup with the full instruction content from below.
                  </div>
                  <rich-editor 
                    v-model="summaryInstruction"  
                    label="Summary Instructions."
                    id="summary"
                  />
                </v-card-text>
                
                <v-card-title>Instruction Content</v-card-title>
                <v-card-text>
                  <div style="text-align: justify; text-justify: inner-word" class="mb-4">
                    To have full instructions please create a new page and set a page name and page content.  Click 'Update Page' to add the page and then click 'Save' to save the instructions.
                    <i>Adding and Deleting pages only persists after saving.</i>
                  </div>
                  <v-layout row wrap>
                    <v-flex sm6 md3>
                      <v-list style="min-height: 150px">
                        <v-toolbar color="grey" dark>
                          <v-toolbar-title>Instruction Pages</v-toolbar-title>
                          <v-spacer></v-spacer>
                          <v-btn icon @click="newPage">
                            <v-icon>add</v-icon>
                          </v-btn>
                        </v-toolbar>
                        
                        <v-list-item
                          v-for="page in instructionPages"
                          :key="page.pageName"
                          @click="selectPage(page.pageName)"
                          :value="selectedInstructionPageName == page.pageName"
                        >
                          <v-list-item-content>
                            <v-list-item-title>{{ page.pageName }}</v-list-item-title>
                          </v-list-item-content>
                        </v-list-item>
                        <v-list-item v-if="instructionPages.length == 0">
                          <v-list-item-content>
                            <v-list-item-title>No Pages</v-list-item-title>
                          </v-list-item-content>
                        </v-list-item>
                      </v-list>
                    </v-flex>
                    
                    <v-flex sm6 md9 pl-4 v-if="instructionPageSelected">
                      <v-form>
                        <v-text-field
                          v-model="selectedInstructionPageName"
                          label="Instruction Page Name"
                          box
                          required
                        ></v-text-field>
                    
                        <rich-editor 
                          v-model="selectedInstructionPageContent"  
                          label="Full Instructions."
                          id="full"
                        />
                      
                        <v-layout justify-end>
                          <v-btn color="error" text outlined class="ma-2" @click="deletePage">Delete Page</v-btn>
                          <v-btn color="primary" text outlined class="ma-2" @click="savePage">Update Page</v-btn>
                        </v-layout>
                      
                      </v-form>
                    </v-flex>
                    <v-flex sm6 md8 pl-4 v-if="!instructionPageSelected">
                      <v-layout align-center justify-space-around row fill-height pb-4>
                        <h3 class="grey--text">Select a page or create a new one</h3>
                      </v-layout>
                    </v-flex>
                  </v-layout>
                </v-card-text>
                
                <v-layout justify-end>
                  <v-btn color="primary" @click="instructionSave">Save Instructions</v-btn>
                </v-layout>
              </v-card>
            
            </v-form>
            
          </v-card-text>
        </v-card>
      </v-flex>
    </v-row>
  </v-col>
</template>

<script>
  import { mapState } from 'vuex'
  
  import RichEditor from '../components/RichEditor.vue'  
  
  export default {
    components: { RichEditor },
    data () {
      return {
        valid: false,
        error: false,
        errorText: 'an unknown error occurred',
        
        savingDialog: false,
        
        existingInstruction: false,
        
        surveyName: (this.$route.params.name) ? this.$route.params.name : '',
        
        popoverWindow: false,
        summaryInstruction: '',
                
        instructionPageSelected: false,
        instructionPages: [],
        
        selectedInstructionPageName: '',
        selectedInstructionPageContent: '',
        
        // raw object
        instructionObject: {}
        
      }
    },
    computed: {
      headerText() {
        return (this.surveyName) ? this.surveyName : 'New Instruction'
      },
      
      availableSurveyList() {
        var all = [{name: 'calendar'}].concat(this.surveys)
        
        for (var a in all) {
          for (var b in this.instructions) {
            // remove survey names where the instructions already exist
            if (all[a].name == this.instructions[b].surveyName) {
              all.splice(a, 1)
              break
            }
          }  
        }
        
        // add the current survey name so the selection shows
        if (this.surveyName) {
          all.push({name: this.surveyName})
        }
        
        return all
      },
      
      ...mapState('surveys', ['surveys']),
      ...mapState('instructions', ['instructions']),
    },
    mounted() {
      this.instruction()
            
      this.$store.dispatch('surveys/fetchSurveys').then(() => {
        this.$store.dispatch('instructions/fetchInstructions')
      })
    },
    methods: {      
      newPage() {
        this.instructionPageSelected = true
        this.selectedInstructionPageName = ''
        this.selectedInstructionPageContent = ''
      },
      
      selectPage(selectedPage) {
        this.instructionPageSelected = true
        this.selectedInstructionPageName = selectedPage
        
        for (var index in this.instructionPages) {
          if (this.instructionPages[index].pageName == selectedPage) {
            this.selectedInstructionPageContent = this.instructionPages[index].pageContent
            break
          }
        }
      }, 
      
      savePage() {
        if (this.selectedInstructionPageName) {
          for (var index in this.instructionPages) {
            if (this.instructionPages[index].pageName == this.selectedInstructionPageName) {
              this.instructionPages[index].pageContent = this.selectedInstructionPageContent
              return
            }
          }
          
          var page = {pageName: this.selectedInstructionPageName, pageContent: this.selectedInstructionPageContent}
          this.instructionPages.push(page)
        }
      }, 
      
      deletePage() {
        if (this.selectedInstructionPageName) {
          for (var index in this.instructionPages) {
            if (this.instructionPages[index].pageName == this.selectedInstructionPageName) {
              this.instructionPages.splice(index, 1)
              
              this.selectedInstructionPageName = ''
              this.selectedInstructionPageContent = ''
          
              this.instructionPageSelected = false
              return
            }
          }
        }
      },
      
      
      instruction() {
        var that = this
        
        if (this.surveyName) {
          this.axios.get('/instructions/' + this.surveyName)
          .then(response => {
            if (response.data.data) {
              that.existingInstruction = true
              that.instructionObject = response.data.data
              
              that.popoverWindow = that.instructionObject.popoverWindow
              that.summaryInstruction = that.instructionObject.summaryInstruction
              
              if (that.instructionObject.instructionPages) {
                that.instructionPages = that.instructionObject.instructionPages
              }
              
            } else {
              that.errorText = 'Instruction data not returned'
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
      instructionSave() {
        if (this.surveyName == '') {
          return
        }
        
        var that = this
        
        this.savingDialog = true
        
        this.axios.post('/instructions/' + this.surveyName, {Data: {popoverWindow: this.popoverWindow, summaryInstruction: this.summaryInstruction, instructionPages: this.instructionPages}})
        .then(() => {
          this.$store.dispatch('instructions/fetchInstructions').then(() => {
            this.functionNameDialog = false
            
            // use a setTimeout so you can see it did something
            setTimeout(function() {
              that.savingDialog = false
            }, 1000)
            
            // reload the page with the correct path
            if (this.$route.path == '/admin/instructions/new') {
              this.$router.push('/admin/instructions/edit/' + this.surveyName)
            }
            
            this.existingInstruction = true
            
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

</style>