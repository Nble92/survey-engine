<template>
  <div>
    <v-dialog v-model="showError" width="500px">
      <v-card color="error" dark class="pt-3">
        <v-card-text>
          <h4>Logic Error</h4>
          <p class="pt-3">{{ errorText }}</p>
        </v-card-text>
      </v-card>
    </v-dialog>
    
    <v-dialog v-model="clickConfirm" persistent max-width="500px">
      <v-card>
        <v-card-title>
          <span class="headline">Confirm Action</span>
        </v-card-title>
        <v-card-text class="text-xs-center">
          Are you sure you wish to <strong>{{ buttonLabel }}</strong>?
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text class="ma-2" @click="clickConfirm = false">No</v-btn>
          <v-btn color="blue darken-1" text class="ma-2" @click="buttonClick">Yes</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    
    <div v-if="error">
      <v-badge
        color="error"
        overlap
      >
        <template v-slot:badge>
          <span><strong>!</strong></span>
        </template>
        <v-btn color="error" outlined @click="showError = true">
          {{ buttonLabel }}
        </v-btn>
      </v-badge>
    </div>
    <div v-if="!error">
      <v-badge
        color="info"
        overlap
        v-if="itemCount > 0"
      >
        <template v-slot:badge>
          <span><strong>{{ itemCount }}</strong></span>
        </template>
        <v-btn color="primary" @click="clickConfirm = true">
          {{ buttonLabel }} 
        </v-btn>
      </v-badge>
      <v-btn v-else color="primary" @click="clickConfirm = true">
        {{ buttonLabel }} 
      </v-btn>
    </div>
  </div>
</template>

<script>
  
  
  export default {
    data () {
      return {
        valid: false,
        
        error: false,
        errorText: '',
        showError: false,
        
        clickConfirm: false,
        
        fullFunctionName: '',
        
        logicResult: {},
        itemCount: 1,
      }
    },
    props: {
      functionName: {
        type: String,
        required: true
      },
      functionType: {
        type: String,
        required: true
      },
      buttonLabel: {
        type: String,
        retuired: true
      },
      surveySessionData: {
        type: Object,
        retuired: true
      }
    },
    watch: {
      surveySessionData() {
        this.runLogic()
      }
    },
    methods: {
      runLogic() {
        this.error = false
        this.errorText = ''
        this.itemCount = 0
        
        if (this.valid) {
          try {
            this.logicResult = window[this.fullFunctionName](this.surveySessionData)
          }
          catch(err) {
            this.error = true
            this.errorText = 'The function "' + this.functionName + '" returned an error: ' + err.message
            
            // eslint-disable-next-line
            console.log('LogicButton Error', err)
            return
          }
          
          // eslint-disable-next-line
          console.log('runLogic result:', typeof this.logicResult, this.logicResult)
          
          this.$store.dispatch('runState/updateState', {lastRunLogic: Date.now()})
          
          
          if (this.functionType == 'followup') {
            if (!Array.isArray(this.logicResult)) {
              this.error = true
              this.errorText = 'The function "' + this.functionName + '" returned a ' + typeof this.logicResult + ', not an array as expected.'
            } else {
              this.itemCount = this.logicResult.length
            }
          } else if (this.functionType == 'completion') {
            if (typeof this.logicResult != 'boolean') {
              this.error = true
              this.errorText = 'The function "' + this.functionName + '" returned a ' + typeof this.logicResult + ', not a boolean as expected.'
            }
          }
        }
      },
      buttonClick() {
        this.clickConfirm = false
        
        if (this.valid && !this.error) {
          this.$emit('click', this.logicResult)
        }
      }
    },
    mounted () {
      // check if the function exists
      this.fullFunctionName = 'SurveyEngine_' + this.functionName
      
      if (window[this.fullFunctionName] == undefined) {
        this.error = true
        this.errorText = 'Function doesnt exist in scope, try reloading this website.'
        this.valid = false
      } else {
        this.valid = true
        this.runLogic()
      }
    }
  }
</script>

<style>

</style>