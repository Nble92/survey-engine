<template>
  <div style="width: 100%">
    
    <v-dialog v-model="showDelete" max-width="360px">
      <v-card>
        <v-card-title>
          <span class="headline">Confirm Delete</span>
        </v-card-title>
        <v-card-text>
          <div v-if="selectedItem.note">
            Are you sure you want to delete the <strong>Note</strong> on {{ selectedItem.date }}?
          </div>
          <div v-else>
            Are you sure you want to delete the <strong>{{ this.surveyConfig.eventLabel }}</strong> on {{ selectedItem.date }}?
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="showDelete = false">Cancel</v-btn>
          <v-btn color="blue darken-1" text @click="confirmDeleteItem">Delete</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    
    <v-dialog v-model="showMoveCopy" max-width="200px">
      <v-list>
        <v-list-item @click="moveItem()">
          Move to {{ targetDate }}
        </v-list-item>
        <v-list-item @click="copyItem()">
          Copy to {{ targetDate }}
        </v-list-item>
      </v-list>
    </v-dialog>
    
    <v-sheet height="650">
      <v-calendar
        ref="calendar"
        v-model="currentMonth"
        color="primary"
        type="month"
        @click:date="dayClick"
        @contextmenu:date="dayRightClick"
      >
        <template v-slot:day="{ date }">
          <div
            style="overflow-y: auto; height: 70%;"
            :class="(highlights[date]) ? 'highlight' : ''"
            @drop="dropItem($event, date)" 
            @dragover="allowDrop"
          >
            <template v-for="note in notesMap[date]">
              <div class="aNote"  @click="noteClick(note.date)" @contextmenu="deleteItem($event, note)" draggable="true" @dragstart="startDraggingItem($event, note)">{{ note.title }}</div>
            </template>
            <template v-for="event in eventsMap[date]">
              <div class="anEvent" @click="eventClick(event.id)" @contextmenu="deleteItem($event, event)" draggable="true" @dragstart="startDraggingItem($event, event)">{{ event.title }}</div>
            </template>
          </div>
        </template>
      </v-calendar>
    </v-sheet>
    
    <v-layout>
      <v-flex sm4 xs12 class="text-sm-left text-center">
        <v-btn @click="$refs.calendar.prev()" text outlined>
          <v-icon dark left>
            keyboard_arrow_left
          </v-icon>
          Previous Month
        </v-btn>
      </v-flex>
      <v-flex offset-sm4 sm4 xs12 class="text-sm-right text-center">
        <v-btn @click="$refs.calendar.next()" text outlined>
          Next Month
          <v-icon right dark>
            keyboard_arrow_right
          </v-icon>
        </v-btn>
      </v-flex>
    </v-layout>
    
  </div>
</template>

<script>
  import { mapState } from 'vuex'
  
  export default {
    data () {
      return {
        showDelete: false,
        showMoveCopy: false,
        
        selectedItem: {},
        targetDate: false,
        
        currentMonth: '',
        
        highlights: {},
        events: [],
        notes: [],
      }
    },
    props: {
      surveyName: {
        type: String,
        required: true
      },
      surveyConfig: {
        type: Object,
        required: false
      },
      surveyState: {
        type: Object,
        required: false
      },
      surveyCreated: {
        type: String,
        required: false
      },
      surveySessionData: {
        type: Object,
        required: false
      }
    },
    watch: {
      surveyState: {
        handler() {
          // eslint-disable-next-line
          console.log('calendarStateUpdate surveyState')
          this.highlightDays()
          this.populateEvents()
          this.populateNotes()
        },
        deep: true
      },
      lastRunLogic() {
        // eslint-disable-next-line
        console.log('calendarStateUpdate surveyData')
        this.populateEvents()
        this.populateNotes()
      }
    },
    computed: {
      eventsMap () {
        const map = {}
        this.events.forEach(e => (map[e.date] = map[e.date] || []).push(e))
        return map
      },
      notesMap () {
        const map = {}
        this.notes.forEach(e => (map[e.date] = map[e.date] || []).push(e))
        return map
      },
      
      ...mapState('runState', ['lastRunLogic'])
    },
    methods: {
      pad(n) {
          return n<10 ? '0'+n : n;
      },
      
      highlightDays() {
        this.highlights = {}
        
        if (this.surveyCreated && this.surveyConfig.calendarHighlightPriorDays) {
          var priorDays = this.surveyConfig.calendarHighlightPriorDaysValue
          
          for (var i = 0; i <= priorDays; i++) {
            var d = new Date(this.surveyCreated);
            d.setDate(d.getDate() - i);
            
            var dateKey = d.getFullYear() + '-' + this.pad(d.getUTCMonth() + 1) + '-' + this.pad(d.getUTCDate())
            
            this.highlights[dateKey] = true
          }
        }
        // eslint-disable-next-line
        console.log('highlighting days', this.highlights)
      },
      populateEvents() {        
        this.events = []

        if (this.surveyState.dataRefs.length > 0) {
          this.surveyState.dataRefs.forEach(ref => {
            if (ref.dataKey.indexOf('-') != -1) {
              var eventLabel = (this.surveyConfig && this.surveyConfig.eventLabel) ? this.surveyConfig.eventLabel : 'Event'
              
              if (this.surveySessionData && this.surveyConfig && this.surveyConfig.eventLabelUseFunction && this.surveyConfig.eventLabelFunction) {
                eventLabel = this.runEventLabelFunction(ref.sessionDataID)
              } 
              
              var event = {id: ref.sessionDataID, title: eventLabel, date: ref.dataKey}
              this.events.push(event)
            }
          })
        }
      },
      runEventLabelFunction(sessionDataID) {
        if (this.surveyConfig && this.surveySessionData.sessionData && this.surveyConfig.eventLabelUseFunction && this.surveyConfig.eventLabelFunction) {
          var sessionData = this.surveySessionData.sessionData.find(ref => ref.sessionDataID == sessionDataID)
          
          if (sessionData) {
            try {
              var functionName = 'SurveyEngine_' + this.surveyConfig.eventLabelFunction
              var result = window[functionName](sessionData)
              return result
            }
            catch(err) {
              // eslint-disable-next-line
              console.log('runEventLabelFunction error: ', err)
              return 'error running function'
            }
          } else {
            return 'sessionDataID not found in surveySessionData'
          }
          
        } else {
          return 'eventLabelUseFunction not correctly defined'
        }
      },
      populateNotes() {
        this.notes = []
        
        if (this.surveyState.notes && this.surveyState.notes.length > 0) {
          this.surveyState.notes.forEach(ref => {
            var note = {title: ref.note, date: ref.date, note: ref.note}
            this.notes.push(note)
          })
        }
      },
      dayClick(day) {
        this.$emit('dayClicked', day.date)
      },
      dayRightClick(day) {
        this.$emit('dayRightClicked', day.date)
      },
      eventClick(event) {
        this.$emit('eventClicked', event)
      },
      noteClick(note) {
        this.$emit('noteClicked', note)
      },
      
      startDraggingItem(e, item) {
          e.stopPropagation(); // Prevent event from bubbling up
        e.dataTransfer.effectAllowed = 'move'
        e.dataTransfer.setData('item', JSON.stringify(item))
                console.log("ARE YOU RUSHING OR ARE YOU DRAGGING???")

      },
            
      allowDrop (e) {
        e.preventDefault()
      },
      dropItem(e, date) {
        this.selectedItem = JSON.parse(e.dataTransfer.getData('item'))
        if (this.selectedItem.date != date) {
          this.targetDate = date
          this.showMoveCopy = true
                  console.log("Dropping works")

        }
      },
      
      deleteItem(e, item) {
        e.preventDefault()
        this.selectedItem = item
        this.showDelete = true
      },
      confirmDeleteItem() {
        this.showDelete = false
        if (this.selectedItem.note) {
          this.$emit('noteDeleted', this.selectedItem)
        } else {
          this.$emit('eventDeleted', this.selectedItem)
        }
      },
      moveItem() {
        if (this.selectedItem.note) {
          this.$emit('noteMoved', this.selectedItem, this.targetDate)
          console.log("Note Move Fucking works!")

        } else {
          this.$emit('eventMoved', this.selectedItem, this.targetDate)
          console.log("Event Move Fucking works!")
        }
        this.showMoveCopy = false
      },
      copyItem() {
        if (this.selectedItem.note) {
          this.$emit('noteCopied', this.selectedItem, this.targetDate)
        } else {
          this.$emit('eventCopied', this.selectedItem, this.targetDate)
        }
        this.showMoveCopy = false
      }
    },
    mounted () {
      this.highlightDays()
      this.populateEvents()
      this.populateNotes()
      console.log("testing testing")
      
      if (this.surveyCreated) {
        var d = new Date(this.surveyCreated)
        this.currentMonth = d.getFullYear() + '-' + this.pad(d.getUTCMonth() + 1) + '-' + this.pad(d.getUTCDate())
      }
      const stupidElement = document.querySelector('#survey > div.v-application--wrap > main > div > div > div > div > div.flex.pa-0 > div.v-card__text > div > div.layout.align-start.justify-space-around.row.fill-height > div > div > div.v-sheet.theme--light > div')
      console.log("This stupid element:: ", stupidElement)
      stupidElement.dragstart = null;
    }
  }
</script>

<style>
  .anEvent {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    border-radius: 2px;
    background-color: #1867c0;
    color: #ffffff;
    border: 1px solid #1867c0;
    font-size: 12px;
    padding: 3px;
    cursor: pointer;
    margin-top: 1px;
    margin-bottom: 1px;
    left: 4px;
    margin-right: 8px;
    position: relative;
  }
  .aNote {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    border-radius: 2px;
    background-color: #8BC34A;
    color: #ffffff;
    border: 1px solid #8BC34A;
    font-size: 12px;
    padding: 3px;
    cursor: pointer;
    margin-top: 1px;
    margin-bottom: 1px;
    left: 4px;
    margin-right: 8px;
    position: relative;
  }
  
  .v-calendar-weekly__day-label {
    text-align: right !important;
  }
  
  .v-calendar-weekly__day-label .v-btn--fab.v-size--small {
    height: auto !important;
    width: auto !important;
    min-width: 24px !important;
    min-height: 24px !important;
    margin: 4px !important;
    margin-right: 10px !important;
  }
  
  .v-calendar-weekly__day-label .v-btn--round {
    padding: 2px !important; 
    border-radius: 10% !important;
  }
  
  .highlight {
    background-image: linear-gradient(0deg, rgba(139,195,74, 0.2), rgba(255,255,255,0.1));
  }
</style>