<template>
  <v-col md="12">
    <v-row>
      <div class="text-center pa-5" style="width: 100%">
        <v-icon x-large class="mb-4" style="font-size: 100px;" color="blue darken-1">assignment</v-icon>
        <h2 class="headline">{{ welcomeText }}</h2>
        <div class="subheading pa-2">
          To begin a survey session, select "<strong>Run Survey</strong>" in the menu on the left.
          <br>
          To review or modify completed surveys, select "<strong>Completed Surveys</strong>" in the menu on the left.
        </div>
      </div>
    </v-row>
  </v-col>
</template>

<script>
  import { mapState } from 'vuex'
  export default {
    computed: {
      ...mapState('config', ['friendlyName', 'welcomeText', 'contactName', 'contactEmail']),
      ...mapState('user', ['username']),
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>