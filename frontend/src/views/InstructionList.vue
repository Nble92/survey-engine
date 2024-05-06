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
      
      <v-dialog v-model="instructionDeleteDialog" persistent max-width="600px">
        <v-card>
          <v-card-title>
            <span class="headline">Delete Instruction</span>
          </v-card-title>
          <v-card-text class="text-xs-center">
            Are you sure you want to delete the instructions for '<strong>{{ deleteInstructionName }}</strong>'?
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="instructionDeleteDialog = false">No</v-btn>
            <v-btn color="blue darken-1" text @click="deleteInstructionConfirm()">Yes, Delete it</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      
      <v-card style="width: 100%">
        <v-toolbar class="blue-grey white--text">
          <v-toolbar-title class="headline font-weight-regular"><v-icon dark class="mr-2">live_help</v-icon> Instructions</v-toolbar-title>
          <v-spacer></v-spacer>
          <v-btn @click="$router.push('/admin/instructions/new')" outlined color="white">New Instruction Set</v-btn>
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
          :headers="instructionListHeaders"
          :items="instructions"
          class="elevation-1"
          :search="search"
          :footer-props="{
              itemsPerPageOptions: [20, 50, 100],
            }"
        >
          <template v-slot:item.surveyName="{ item }">
            <span class="font-weight-medium">{{ item.surveyName }}</span>
          </template>
        
          <template v-slot:item.actions="{ item }">
            <v-icon
              small
              class="mr-2"
              @click="editInstruction(item.surveyName)"
            >
              edit
            </v-icon>
            <v-icon
              small
              @click="deleteInstruction(item.surveyName)"
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
        
        instructionListHeaders: [
          { text: 'Survey Name', value: 'surveyName' },
          { text: 'Actions', value: 'actions', sortable: false, align: 'center' }
        ],

        // deletion stuff
        instructionDeleteDialog: false,
        deleteInstructionName: '',
      }
    },
    computed: {
      ...mapState('instructions', ['instructions']),
    },
    mounted() {
      var that = this
      this.$store.dispatch('instructions/fetchInstructions').catch(function (error) {
        if (error.response && error.response.data) {
          that.errorText = error.response.data.message
        } else {
          that.errorText = error.message
        }
        that.error = true
      })
    },
    methods: {
      editInstruction(name) {
        this.$router.push('/admin/instructions/edit/' + name)
      },
      deleteInstruction(name) {
        this.deleteInstructionName = name
        this.instructionDeleteDialog = true
      },
      deleteInstructionConfirm() {
        var that = this
        
        this.axios.delete('/instructions/' + this.deleteInstructionName)
        .then(() => {
          this.$store.dispatch('instructions/fetchInstructions').catch(function (error) {
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
        
        this.deleteInstructionName = ''
        this.instructionDeleteDialog = false
      }
    },
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>