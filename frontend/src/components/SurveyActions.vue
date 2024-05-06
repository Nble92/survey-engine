<template>
  <div>
    <v-dialog v-model="actionsPopup" persistent max-width="500px">
      <v-card>
        <v-card-title>
          <span class="headline">Survey Actions</span>
        </v-card-title>
        <v-card-text>
          <v-list rounded>
            <v-subheader>Current Actions</v-subheader>
            <v-list-item-group>
              <v-list-item @click="endSession" :disabled="!canEndSession">
                <v-list-item-icon>
                  <v-icon>highlight_off</v-icon>
                </v-list-item-icon>
                <v-list-item-content>
                  <v-list-item-title>End Session <span class="red--text" v-if="!canEndSession">(not available)</span></v-list-item-title>
                </v-list-item-content>
              </v-list-item>
              <v-list-item @click="viewCalendar" v-if="surveyType == 'calendar'">
                <v-list-item-icon>
                  <v-icon>calendar_today</v-icon>
                </v-list-item-icon>
                <v-list-item-content>
                  <v-list-item-title>View Calendar</v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </v-list-item-group>
          </v-list>
          
          <v-list rounded v-if="completedSurveys && completedSurveys.length > 0">
            <v-subheader>Edit Completed Survey</v-subheader>
            <v-list-item-group>
              <v-list-item
                v-for="item in completedSurveys"
                :key="item.sessionDataID"
                @click="editSurvey(item.surveyName, item.sessionDataID)"
              >
                <v-list-item-icon>
                  <v-icon>list</v-icon>
                </v-list-item-icon>
                <v-list-item-content>
                  <v-list-item-title v-text="item.title"></v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </v-list-item-group>
          </v-list>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text class="ma-2" @click="actionsPopup = false">Cancel</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    
    <v-btn text v-if="isLoggedIn && status == 'in_process'" @click="actionsPopup = true">
      <v-icon>list</v-icon>
      <span>Actions</span>
    </v-btn>
  </div>
</template>

<script>
  import { mapState, mapGetters } from 'vuex'
  
  export default {
    data () {
      return {
        actionsPopup: false,
        
        items: [
          { text: 'Real-Time', icon: 'mdi-clock' },
          { text: 'Audience', icon: 'mdi-account' },
          { text: 'Conversions', icon: 'mdi-flag' },
        ],
      }
    },
    watch: {
      
    },
    methods: {
      endSession() {
        // eslint-disable-next-line
        console.log('resetSession called')
        
        this.$store.dispatch('surveySession/resetSessionState')
        this.$store.dispatch('runState/resetState')
        
        this.actionsPopup = false
      },
      viewCalendar() {
        this.$store.dispatch('runState/updateState', {currentView: 'calendar'})
        this.actionsPopup = false
      },
      editSurvey(surveyName, surveySessionDataID) {
        // eslint-disable-next-line
        console.log('editSurvey ' + surveyName + ' ' + surveySessionDataID)
        
        this.$store.dispatch('runState/updateState', {surveySessionDataID: surveySessionDataID, surveyName: surveyName, currentView: 'survey'})
        this.actionsPopup = false
      }
    },
    computed: {
      canEndSession() {
        return (this.allowReset || this.editMode) ? true : false
      },
      completedSurveys() {
        var surveyArray = []
        
        if (this.surveyState.dataRefs) {
          for (var ref of this.surveyState.dataRefs) {
            surveyArray.push({surveyName: ref.surveyName, sessionDataID: ref.sessionDataID, dataKey: ref.dataKey, title: ref.surveyName + ' (' + ref.dataKey + ')'})
          }
        }
        
        return surveyArray
      },
      
      ...mapGetters('user', ['isLoggedIn', 'isAdmin']),
      ...mapState('surveySettings', ['surveyType', 'config', 'surveyTypes']),
      ...mapState('surveySession', ['status', 'sessionID', 'surveyState', 'isTest', 'editMode']),
      ...mapState('runState', ['error', 'errorText', 'currentView', 'surveyName', 'surveySessionDataID', 'surveyDataKey', 'surveyReferenceID', 'surveyData', 'allowReset'])
    },
    mounted () {
      
    }
  }
</script>

<style>

</style>