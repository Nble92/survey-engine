<template>
  <v-container fluid row>
    <div class="pb-4" v-if="error" style="width: 100%">
      <v-alert
        v-model="error"
        type="error"
        transition="scale-transition"
      >{{ errorText }}</v-alert>
    </div>
    <div class="pb-4" v-if="dialog" style="width: 100%">
      <v-alert
        v-model="dialog"
        type="info"
        transition="scale-transition"
      >{{ dialogText }}</v-alert>
    </div>
    
    
    <v-flex>
      <v-card xs12>
        <v-toolbar dark color="secondary">
          <v-toolbar-title class="headline font-weight-regular">Test</v-toolbar-title>
        </v-toolbar>
        <v-container>
          <v-flex xs-12 ma-4 class="text-md-center">
            <v-btn color="secondary" @click="runTest">Run Test</v-btn>
          </v-flex>
        </v-container>
      </v-card>
    </v-flex>
  </v-container>
</template>


<script>  
  export default {
    data () {
      return {
        dialog: false,
        dialogText: '',
        
        error: false,
        errorText: 'an unknown error occurred',
      }
    },

    methods: {
      runTest() {
        var that = this
        this.axios.post('/notify', {subject: 'Testing Notifications', message: 'Test message', additionalCCs: [{fullName: 'Daniel', email: 'test@mail.chote.com'}, {fullName: 'Daniel 2', email: 'test2@mail.chote.com'}]})
        .then((result) => { 
          that.dialog = true
          // eslint-disable-next-line
          console.log(result)
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
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>