<template>
  <v-container>
    <v-row>
      <v-col cols="10" offset="1">
        <div class="text-center">
          <v-icon x-large class="mb-4" style="font-size: 100px;" color="blue darken-1">assignment</v-icon>
          <h2 class="headline" v-if="!config.completedSurveyInstructions || config.completedSurveyInstructions == ''">Thanks for completing the {{ friendlyName }} survey.</h2>
          <div v-else v-html="config.completedSurveyInstructions"></div>
          
          <div v-if="showRedirectButton" class="ma-5">
            <v-row justify="center">
              <v-alert text color="warning" class="text-left">
                <h3>You are not finished!</h3>
                <div v-html="redirectButtonInstruction"></div>
                <v-divider class="my-4 warning" style="opacity: 0.22"></v-divider>
                <div class="mt-4 text-right">
                  <v-btn large color="warning" outlined class="ma-2" :href="redirectURL" target="_blank">Continue to external site</v-btn>
                </div>
              </v-alert>
            </v-row>
          </div>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
  import { mapState } from 'vuex'
  export default {
    computed: {
      showRedirectButton() {
        return this.config.allowRedirectButton
      },
      redirectButtonInstruction() {
        return this.config.redirectButtonInstructions
      },
      redirectURL() {
        return this.config.redirectURL
      },
      ...mapState('config', ['friendlyName']),
      ...mapState('surveySettings', ['surveyType', 'config', 'surveyTypes']),
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>