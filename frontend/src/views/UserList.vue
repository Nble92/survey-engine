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
      
      <v-dialog v-model="userDeleteDialog" persistent max-width="600px">
        <v-card>
          <v-card-title>
            <span class="headline">Delete User</span>
          </v-card-title>
          <v-card-text class="text-xs-center">
            Are you sure you want to delete the user '<strong>{{ deleteUserName }}</strong>'?
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" text @click="userDeleteDialog = false">No</v-btn>
            <v-btn color="blue darken-1" text @click="deleteUserConfirm()">Yes, Delete it</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      
      
      <v-card style="width: 100%">
        <v-toolbar class="blue-grey white--text">
          <v-toolbar-title class="headline font-weight-regular"><v-icon dark class="mr-2">person</v-icon> Users</v-toolbar-title>
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
          :headers="userListHeaders"
          :items="users"
          class="elevation-1"
          :search="search"
          :footer-props="{
              itemsPerPageOptions: [20, 50, 100],
            }"
        >
          <template v-slot:item.username="{ item }">
            <span class="font-weight-medium"> {{ item.username }} </span> <v-chip small v-if="item.username == username">Me</v-chip> 
          </template>
        
          <template v-slot:item.isAdmin="{ item }">
            <v-chip color="primary" dark small v-if="item.isAdmin">Admin</v-chip>
          </template>
        
          <template v-slot:item.actions="{ item }">
            <v-icon
              small
              class="mr-2"
              @click="editUser(item.username)"
            >
              edit
            </v-icon>
            <v-icon
              small
              @click="deleteUser(item.username)"
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
        
        userListHeaders: [
          { text: 'Username', value: 'username' },
          { text: 'Admin', value: 'isAdmin'},
          { text: 'Actions', value: 'actions', sortable: false, align: 'center' }
        ],

        // deletion stuff
        userDeleteDialog: false,
        deleteUserName: '',
        
        users: [],
      }
    },
    computed: {
      ...mapState('user', ['username']),
    },
    mounted() {
      this.loadUsers()
    },
    methods: {
      loadUsers() {
        var that = this
        
        this.axios.get('/user/list')
        .then(response => { 
          that.users = response.data.users
        
          // eslint-disable-next-line
          console.log('users', that.users)
        }).catch(function (error) {
          that.error = true
          if (error.response && error.response.data) {
            that.errorText = error.response.data
          } else if (error.message) {
            that.errorText = error.message
          } else {
            that.errorText = 'Error fetching user list...'
          }
        })
      },
      editUser(name) {
        this.$router.push('/admin/user/edit/' + name)
      },
      deleteUser(name) {
        if (name == this.username) {
          return
        }
        
        this.deleteUserName = name
        this.userDeleteDialog = true
      },
      deleteUserConfirm() {
        var that = this
        
        this.axios.delete('/user/' + this.deleteUserName)
        .then(() => {
          that.loadUsers()
        }).catch(function (error) {
            if (error.response && error.response.data) {
              that.errorText = error.response.data.message
            } else {
              that.errorText = error.message
            }
            that.error = true
        })
        
        this.deleteUserName = ''
        this.userDeleteDialog = false
        
      }
    },
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>