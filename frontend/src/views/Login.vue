<template>
  <v-col md="6" offset-md="3">
    <v-row no-gutters>
      <v-card width="100%">
        <v-toolbar color="primary" dark>
          <v-toolbar-title>Please Login</v-toolbar-title>
        </v-toolbar>
        <v-card-text>
          <v-form v-model="valid" @keyup.native.enter="submit">
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
              :append-icon="hidePass ? 'visibility' : 'visibility_off'"
              @click:append="() => (hidePass = !hidePass)"
              :type="hidePass ? 'password' : 'text'"
              :rules="passwordRules"
              counter
              required
            ></v-text-field>
            <v-layout justify-end justify-space-between pt-4 align-center>
              <a @click="register">Create new account</a>
              <v-btn @click="submit" :class=" { 'primary white--text' : valid, disabled: !valid }" text outlined>Login</v-btn>
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
        hidePass: true,
        password: '',
        passwordRules: [
          (v) => !!v || 'Password is required',
        ],
        username: '',
        usernameRules: [
          (v) => !!v || 'Username is required',
        ],
        }
      },
      methods: {
        register () {
          this.$router.push('/register')
        }, 
        submit () {
          if (this.username && this.password) {
            this.$store.dispatch('user/login', {username: this.username, password: this.password}).then(() => {
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