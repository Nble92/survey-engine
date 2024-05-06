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
      
      <v-dialog v-model="savingDialog" hide-overlay persistent width="300">
        <v-card color="primary" dark class="pt-3">
          <v-card-text>
            Saving User...
            <v-progress-linear
              indeterminate
              color="white"
              class="mb-0"
            ></v-progress-linear>
          </v-card-text>
        </v-card>
      </v-dialog>
      
      <v-card style="width: 100%">
        <v-card-title class="headline font-weight-regular blue-grey white--text"><v-btn fab small text color="white" class="mr-3" to="/admin/user/list"><v-icon>keyboard_arrow_left</v-icon></v-btn>  Edit User</v-card-title>
        <v-card-text>
          <v-form>
            <v-layout justify-center mb-4 mt-4>
              <v-card-text class="pa-4">
                <v-flex xs12 ma-3>
                  <v-text-field
                    v-model="username"
                    label="Username"
                    required
                    disabled
                  ></v-text-field>
                  <v-checkbox
                    label="Administrator"
                    v-model="isAdmin"
                  ></v-checkbox>
                  <v-text-field
                    v-model="fullName"
                    label="Full Name"
                  ></v-text-field>
                  <v-text-field
                    v-model="email"
                    label="Email Address"
                  ></v-text-field>
                  <v-checkbox
                    label="Recieve Notifications"
                    v-model="notifications"
                  ></v-checkbox>
                </v-flex>
              </v-card-text>
            </v-layout>
            <v-layout justify-end>
              <v-btn @click="saveUser" color="primary">Save User</v-btn>
            </v-layout>
          </v-form>
        </v-card-text>
      </v-card>
    </v-row>
  </v-col>
</template>

<script>
  export default {
    data () {
      return {
        valid: false,
        error: false,
        errorText: 'an unknown error occurred',
        
        savingDialog: false,
        
        username: (this.$route.params.name) ? this.$route.params.name : '',
        isAdmin: false,
        fullName: '',
        email: '',
        notifications: false
      }
    },
    computed: {
      
    },
    mounted() {
      this.loadUser()
    },
    methods: {
      loadUser() {
        var that = this
      
        this.axios.get('/user/' + this.username)
        .then((response) => {
          this.username = response.data.username
          this.isAdmin = response.data.isAdmin
          this.fullName = response.data.fullName
          this.email = response.data.email
          this.notifications = response.data.notifications
        }).catch(function (error) {
            if (error.response && error.response.data) {
              that.errorText = error.response.data.message
            } else {
              that.errorText = error.message
            }
            that.error = true
        })
      },
      saveUser() {
        if (this.username == '') {
          return
        }
        
        var that = this
        
        this.savingDialog = true
        
        this.axios.post('/user/' + this.username, { isAdmin: this.isAdmin, fullName: this.fullName, email: this.email, notifications: this.notifications })
        .then(() => {
          this.loadUser()
          
          // use a setTimeout so you can see it did something
          setTimeout(function() {
            that.savingDialog = false
          }, 1000)
        }).catch(function (error) {
            if (error.response && error.response.data) {
              that.errorText = error.response.data.message
            } else {
              that.errorText = error.message
            }
            that.error = true
            
            that.savingDialog = false
        })
        
        
        
      }
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>