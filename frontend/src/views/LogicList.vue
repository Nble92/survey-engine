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
      
      <v-dialog v-model="logicDeleteDialog" persistent max-width="600px">
        <v-card>
          <v-card-title>
            <span class="headline">Delete Function</span>
          </v-card-title>
          <v-card-text class="text-xs-center">
            Are you sure you want to delete the function '<strong>{{ deleteLogicName }}</strong>'?
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="logicDeleteDialog = false">No</v-btn>
            <v-btn color="blue darken-1" text @click="deleteLogicConfirm()">Yes, Delete it</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      
      <v-card style="width: 100%">
        <v-toolbar class="blue-grey white--text">
          <v-toolbar-title class="headline font-weight-regular"><v-icon dark class="mr-2">code</v-icon> Functions</v-toolbar-title>
          <v-spacer></v-spacer>
          <v-btn @click="$router.push('/admin/logic/new')" outlined color="white">New Function</v-btn>
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
          :headers="logicListHeaders"
          :items="logicFunctions"
          class="elevation-1"
          :search="search"
          :footer-props="{
              itemsPerPageOptions: [20, 50, 100],
            }"
        >
          
          <template v-slot:item.functionName="{ item }">
            <span class="font-weight-medium">{{ item.functionName }}</span>
          </template>
        
          <template v-slot:item.actions="{ item }">
            <v-icon
              small
              class="mr-2"
              @click="editLogic(item.functionName)"
            >
              edit
            </v-icon>
            <v-icon
              small
              @click="deleteLogic(item.functionName)"
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
        
        logicListHeaders: [
          { text: 'Function Name', value: 'functionName' },
          { text: 'Actions', value: 'actions', sortable: false, align: 'center' }
        ],

        // deletion stuff
        logicDeleteDialog: false,
        deleteLogicName: '',
      }
    },
    computed: {
      ...mapState('logicFunctions', ['logicFunctions']),
    },
    mounted() {
      var that = this
      this.$store.dispatch('logicFunctions/fetchLogicFunctions').catch(function (error) {
        if (error.response && error.response.data) {
          that.errorText = error.response.data.message
        } else {
          that.errorText = error.message
        }
        that.error = true
      })
    },
    methods: {
      editLogic(name) {
        this.$router.push('/admin/logic/edit/' + name)
      },
      deleteLogic(name) {
        this.deleteLogicName = name
        this.logicDeleteDialog = true
      },
      deleteLogicConfirm() {
        var that = this
        
        this.axios.delete('/logic/' + this.deleteLogicName)
        .then(() => {
          this.$store.dispatch('logicFunctions/fetchLogicFunctions').catch(function (error) {
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
        
        this.deleteLogicName = ''
        this.logicDeleteDialog = false
      }
    },
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>