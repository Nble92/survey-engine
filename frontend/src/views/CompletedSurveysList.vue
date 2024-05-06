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
      <v-card style="width: 100%">
        <v-toolbar class="blue-grey white--text">
          <v-toolbar-title class="headline font-weight-regular">Completed Surveys</v-toolbar-title>
          <v-spacer></v-spacer>
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
          :items="surveySessionItems"
          class="elevation-1"
          :search="search"
          :custom-sort="tableSort"
          sort-by="created"
          :sort-desc="true"
          :footer-props="{
              itemsPerPageOptions: [20, 50, 100],
            }"
        >
          <template v-slot:item.sessionID="{ item }">
            <a class="font-weight-medium" @click="$router.push('/survey/view/' + item.sessionID)">{{ item.sessionID }}</a>
          </template>
          
          <template v-slot:item.isTest="{ item }">
            <v-chip color="green" dark small v-if="item.isTest">Test</v-chip>
          </template>
          
        </v-data-table>

      </v-card>
    </v-row>
  </v-col>
</template>


<script>
  //import { mapState } from 'vuex'

  export default {
    data () {
      return {
        valid: false,
        error: false,
        errorText: 'an unknown error occurred',
        
        search: '',
        
        
        surveyListState: {
          descending: true,
          sortBy: 'created'
        },
        surveyListHeaders: [
          { text: 'Survey Session ID', value: 'sessionID', sortable: false },
          { text: 'Test', value: 'isTest' },
          { text: 'Created', value: 'created' },
          { text: 'Last Modified', value: 'modified' },
          { text: 'Survey User', value: 'username' },
        ],
        
        surveySessionItems: []
      }
    },
    computed: {
      
    },
    mounted() {
      var that = this
      
      this.axios.get('/survey/sessionList')
      .then(response => { 
        that.surveySessionItems = response.data.surveySessions
        
        // eslint-disable-next-line
        console.log('surveySessionItems', that.surveySessionItems)
      }).catch(function (error) {
        that.error = true
        if (error.response && error.response.data) {
          that.errorText = error.response.data
        } else if (error.message) {
          that.errorText = error.message
        } else {
          that.errorText = 'Error fetching survey session list...'
        }
      })
    },
    methods: {
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