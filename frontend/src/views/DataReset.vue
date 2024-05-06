<template>
  <v-container fluid row>
    <div class="pb-4" v-if="error" style="width: 100%">
      <v-alert
        v-model="error"
        type="error"
        transition="scale-transition"
      >{{ errorText }}</v-alert>
    </div>
    <div class="pb-4" v-if="reset" style="width: 100%">
      <v-alert
        v-model="reset"
        type="info"
        transition="scale-transition"
      >Data has been reset</v-alert>
    </div>
    
    
    <v-flex>
      <v-card xs12>
        <v-toolbar dark color="error">
          <v-toolbar-title class="headline font-weight-regular">Data Reset</v-toolbar-title>
        </v-toolbar>
        <v-container>
          <v-flex xs-12 ma-4 class="text-md-center">
            <v-btn color="error" @click="resetData">Reset all data</v-btn>
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
        reset: false,
        
        error: false,
        errorText: 'an unknown error occurred',
      }
    },

    methods: {
      resetData() {
        var that = this
        this.axios.get('/survey/reset')
        .then(() => { 
          that.reset = true
        }).catch(function (error) {
          that.error = true
          that.errorText = error
        })
      }
    },
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>