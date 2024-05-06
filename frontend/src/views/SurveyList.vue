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
      
      <v-dialog v-model="surveyDeleteDialog" persistent max-width="600px">
        <v-card>
          <v-card-title>
            <span class="headline">Delete Survey</span>
          </v-card-title>
          <v-card-text class="text-xs-center">
            Are you sure you want to delete the survey '<strong>{{ deleteSurveyName }}</strong>'?
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="surveyDeleteDialog = false">No</v-btn>
            <v-btn color="blue darken-1" text @click="deleteSurveyConfirm()">Yes, Delete it</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      
      <v-card style="width: 100%">
        <v-toolbar class="blue-grey white--text">
          <v-toolbar-title class="headline font-weight-regular"><v-icon dark class="mr-2">library_books</v-icon> Surveys</v-toolbar-title>
          <v-spacer></v-spacer>
          <v-btn @click="$router.push('/admin/survey/new')" outlined color="white">New Survey</v-btn>
        </v-toolbar>
        <v-card-title>
          <v-text-field
            v-model="search"
            append-icon="search"
            label="Search"
            single-line
            hide-details
          ></v-text-field>
        </v-card-title>
        <v-data-table
          :headers="surveyListHeaders"
          :items="surveys"
          class="elevation-1"
          :search="search"
          :custom-sort="tableSort"
          sort-by="modified"
          :sort-desc="true"
          :footer-props="{
              itemsPerPageOptions: [20, 50, 100],
            }"
        >
          <template v-slot:item.name="{ item }">
            <span class="font-weight-medium">{{ item.name }}</span>
          </template>
          
          <template v-slot:item.actions="{ item }">
            <v-icon
              small
              class="mr-2"
              @click="editSurvey(item.name)"
            >
              edit
            </v-icon>
            <v-icon
              small
              @click="deleteSurvey(item.name)"
            >
              delete
            </v-icon>
          </template>
        </v-data-table>
      </v-card>
    </v-row>
  </v-col>
</template>


<script>
  import { mapState } from 'vuex'

  export default {
    data () {
      return {
        valid: false,
        error: false,
        errorText: 'an unknown error occurred',
        
        search: '',
        
        surveyListState: {
          descending: true,
          sortBy: 'modified'
        },
        surveyListHeaders: [
          { text: 'Survey Name', value: 'name' },
          { text: 'Last Modified', value: 'modified'},
          { text: 'Author', value: 'author', class: 'hidden-md-and-down'},
          { text: 'Actions', value: 'actions', sortable: false, align: 'center' }
        ],

        // deletion stuff
        surveyDeleteDialog: false,
        deleteSurveyName: '',
      }
    },
    computed: {
      ...mapState('surveys', ['surveys']),
    },
    methods: {
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
      },
      tableSort(items, index, isDesc) {
        if (index && isDesc) {
          items.sort((a, b) => {
            if (index[0] === "created" || index[0] === "modified") {
              var d1 = new Date(a[index[0]])
              var d2 = new Date(b[index[0]])
              if (!isDesc[0]) {
                return (d1>d2)-(d1<d2);
              } else {
                return (d2>d1)-(d2<d1);
              }
            } else {
              if (!isDesc[0]) {
                return a[index[0]] < b[index[0]] ? -1 : 1;
              } else {
                return b[index[0]] < a[index[0]] ? -1 : 1;
              }
            }
          });
        }
        return items;
      }
    },
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>