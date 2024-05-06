<template>
  <v-col md="6" offset-md="3">
    <v-row no-gutters>
      <v-card width="100%">
        <v-toolbar color="primary" dark>
          <v-toolbar-title>Account Creation</v-toolbar-title>
        </v-toolbar>
        <v-card-text>
          <v-form v-model="valid" ref="form" @submit.prevent="register({ email, password })">
            <div class="pb-3">
            <v-alert
              v-model="error"
              type="error"
              transition="scale-transition"
            >{{ errorText }}</v-alert>
            </div>
            <v-text-field
              label="Enter your username"
              v-model="username"
              :rules="usernameRules"
              required
            ></v-text-field>
            <v-text-field
              label="Enter your password"
              v-model="password"
              min="8"
              :append-icon="e1 ? 'visibility' : 'visibility_off'"
              @click:append="() => (e1 = !e1)"
              :type="e1 ? 'password' : 'text'"
              :rules="passwordRules"
              counter
              required
            ></v-text-field>
            <v-layout justify-end justify-space-between pt-4 align-center>
              <a @click="login">Use existing account</a>
              <v-btn @click="submit" :class=" { 'primary white--text' : valid, disabled: !valid }" text outlined>Create Account</v-btn>
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
        error: false,
        errorText: 'an unknown error occurred',
        valid: false,
        e1: true,
        name: '',
        nameRules: [
          (v) => !!v || 'Name is required'
        ],
        password: '',
        password_confirm: '',
        passwordRules: [
          (v) => !!v || 'Password is required'
        ],
        username: '',
        usernameRules: [
          (v) => !!v || 'Username is required',
        ],
      }
    },
    methods: {
      login () {
        this.$router.push('/login')
      }, 
      submit () {
        if (this.username && this.password) {
          this.$store.dispatch('user/register', {username: this.username, password: this.password}).then(() => {
             this.$router.push('/')
           }).catch((error) => {
            if (error.response && error.response.data) {
              this.errorText = error.response.data.message
            } else {
              this.errorText = error.message
            }
            
            this.error = true
          })
        }
      },
      clear () {
        this.$refs.form.reset()
      }
    },
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>