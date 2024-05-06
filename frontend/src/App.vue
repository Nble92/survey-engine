<template>
  <v-app id="survey">
    <div v-if="isLoggedIn">
      <v-navigation-drawer
        :clipped="$vuetify.breakpoint.mdAndUp"
        v-model="drawer"
        app
      >
        <v-list dense>
          <template v-for="item in items">
            <v-row
              v-if="item.heading"
              :key="item.heading"
              row
            >
              <v-col cols="6">
                <v-subheader v-if="item.heading">
                  {{ item.heading }}
                </v-subheader>
              </v-col>
            </v-row>
            <v-list-group
              v-else-if="item.children && !item.requiresAdmin || item.requiresAdmin && (item.requiresAdmin == isAdmin)"
              :key="item.text"
              v-model="item.model"
              :prepend-icon="item.model ? item.icon : item['icon-alt']"
            >
              <template v-slot:activator>
                <v-list-item-content>
                  <v-list-item-title>
                    {{ item.text }}
                  </v-list-item-title>
                </v-list-item-content>
              </template>
              <v-list-item
                v-for="(child, i) in item.children"
                :key="i"
                :to="child.route" 
                router
                class="ml-3"
              >
                <v-list-item-action v-if="child.icon">
                  <v-icon>{{ child.icon }}</v-icon>
                </v-list-item-action>
                <v-list-item-content>
                  <v-list-item-title>
                    {{ child.text }}
                  </v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </v-list-group>
            <v-list-item
              v-else-if="!item.requiresAdmin || item.requiresAdmin && (item.requiresAdmin == isAdmin)"
              :key="item.text"
              :to="item.route" 
              router
            >
              <v-list-item-action>
                <v-icon>{{ item.icon }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>
                  {{ item.text }}
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </template>
          
          <v-list-item @click="menuLogout()">
            <v-list-item-action>
              <v-icon>exit_to_app</v-icon>
            </v-list-item-action>
            <v-list-item-content>
              <v-list-item-title>
                Logout
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-navigation-drawer>
    </div>
    
    <v-app-bar
      :clipped-left="$vuetify.breakpoint.mdAndUp"
      color="blue darken-2"
      dark
      app
      fixed
    >
      <v-toolbar-title>
        <v-app-bar-nav-icon @click.stop="drawer = !drawer" v-if="isLoggedIn"></v-app-bar-nav-icon>
        <span class="ml-2"><v-icon color="blue lighten-4">assignment</v-icon> {{ friendlyName }}</span>
      </v-toolbar-title>
      
      <v-spacer></v-spacer>
      
      <v-btn text v-if="username" disabled>
        <v-icon>person</v-icon>
        <span>{{ username }}</span>
      </v-btn>
      <survey-actions></survey-actions>
    </v-app-bar>
    
    
    <v-content>
      <v-container align-start fluid>
        <router-view></router-view>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
  import { mapState, mapGetters, mapActions } from 'vuex'
  import SurveyActions from './components/SurveyActions.vue'
  
  export default {
    components: {
      SurveyActions
    },
    data: () => ({
      drawer: null,
      adminExpanded: null,
      items: [
        { icon: 'assignment', text: 'Run Survey', route: '/survey/run' },
        { icon: 'list', text: 'Completed Surveys', route: '/survey/list' },
        { 
          icon: 'perm_data_setting', 
          'icon-alt': 'perm_data_setting',
          text: 'Admin', 
          requiresAdmin: true,
          model: false,
          children: [
            {icon: 'ballot', text: 'Configuration', route: '/admin/config'},
            {icon: 'library_books', text: 'Edit Surveys', route: '/admin/survey/list'},
            {icon: 'code', text: 'Edit Functions', route: '/admin/logic/list'},
            {icon: 'live_help', text: 'Edit Instructions', route: '/admin/instructions/list'},
            {icon: 'collections', text: 'Manage Images', route: '/admin/images/list'},
            {icon: 'person', text: 'Manage Users', route: '/admin/user/list'},
            {icon: 'import_export', text: 'Export Data', route: '/admin/export'},
            {icon: 'backup', text: 'Backup & Restore', route: '/admin/backup'}
          ]
        }
      ]
    }),
    computed: {
      ...mapState('config', ['friendlyName']),
      ...mapState('user', ['username']),
      ...mapGetters('user', ['isLoggedIn', 'isAdmin']),
      ...mapActions('user', ['logout']),
    },
    props: {
      source: String
    },
    methods: {
      menuLogout () {
        var that = this
        this.$store.dispatch('surveySession/resetSessionState').then(() => {
          this.$store.dispatch('user/logout').then(() => {
            that.$router.push('/login')
          })
        })
      },
    }
  }
</script>
  
<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
  h1 p, h2 p, h3 p, h4 p, h5 p {
    margin-bottom: 0 !important;
  }
</style>