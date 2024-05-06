import Vue from 'vue'
import Router from 'vue-router'
import store from './store'

import BackendError from '@/views/BackendError.vue'
import NotImplemented from '@/views/NotImplemented.vue'

import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'

import Welcome from '@/views/Welcome.vue'
import RunSurvey from '@/views/RunSurvey.vue'
import SurveyCompleted from '@/views/SurveyCompleted.vue'

import CompletedSurveysList from '@/views/CompletedSurveysList.vue'
import ViewCompletedSurvey from '@/views/ViewCompletedSurvey.vue'


import SurveyConfig from '@/views/SurveyConfig.vue'

import LogicList from '@/views/LogicList.vue'
import EditLogic from '@/views/EditLogic.vue'

import SurveyList from '@/views/SurveyList.vue'
import EditSurvey from '@/views/EditSurvey.vue'

import InstructionList from '@/views/InstructionList.vue'
import EditInstruction from '@/views/EditInstruction.vue'

import ImagesList from '@/views/ImagesList.vue'

import UserList from '@/views/UserList.vue'
import EditUser from '@/views/EditUser.vue'

import DataExport from '@/views/DataExport.vue'
import DataReset from '@/views/DataReset.vue'

import Tests from '@/views/Tests.vue'

Vue.use(Router)

var router = new Router({
  routes: [
    {
      path: '/',
      redirect: '/welcome'
    },
    {
      path: '/welcome',
      name: 'Welcome',
      component: Welcome,
      meta: {
        requiresAuth: true
      }
    },
    {
      path: '/not_implemented',
      name: 'Not Implemented',
      component: NotImplemented,
      meta: {
        requiresAuth: true
      }
    },
    {
      path: '/login',
      name: 'Login',
      component: Login,
    },
    {
      path: '/register',
      name: 'Register',
      component: Register,
    },
    {
      path: '/survey/run/:name?',
      name: 'Run Survey',
      component: RunSurvey,
    },
    {
      path: '/survey/completed',
      name: 'Completed Survey',
      component: SurveyCompleted,
    },
    {
      path: '/survey/list/',
      name: 'Completed Surveys',
      component: CompletedSurveysList,
    },
    {
      path: '/survey/view/:id',
      name: 'View Completed Survey',
      component: ViewCompletedSurvey,
    },
    {
      path: '/backend_error',
      name: 'Backend Error',
      component: BackendError,
    },
    {
      path: '/admin/config',
      name: 'Survey Config',
      component: SurveyConfig,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/logic/list',
      name: 'Survey Logic',
      component: LogicList,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/logic/new',
      name: 'New Function',
      component: EditLogic,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/logic/edit/:name',
      name: 'Edit Function',
      component: EditLogic,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/survey/list',
      name: 'Manage Surveys',
      component: SurveyList,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/survey/new',
      name: 'New Survey',
      component: EditSurvey,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/survey/edit/:name',
      name: 'Edit Survey',
      component: EditSurvey,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/instructions/list',
      name: 'Survey Instructions',
      component: InstructionList,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/instructions/new',
      name: 'New Instruction',
      component: EditInstruction,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/instructions/edit/:name',
      name: 'Edit Instructions',
      component: EditInstruction,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/images/list',
      name: 'Images',
      component: ImagesList,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/user/list',
      name: 'User List',
      component: UserList,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/user/edit/:name',
      name: 'Edit User',
      component: EditUser,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/export',
      name: 'Data Export',
      component: DataExport,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/reset',
      name: 'Data Reset',
      component: DataReset,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/test',
      name: 'Tests',
      component: Tests,
      meta: {
        requiresAdmin: true
      }
    },
    {
      path: '/admin/backup',
      name: 'Backup & Restore',
      component: NotImplemented,
      meta: {
        requiresAdmin: true
      }
    },
  ]
})



// this validateAccess/router.beforeEach logic is a bit messy, but its the simplist way I can think of right now....
function validateAccess(to) {
  if (to.matched.some(record => record.meta.requiresAdmin)) {
    if (store.getters['user/isAdmin']) {
      return true
    } else {
      return false
    }
  } else if (to.matched.some(record => record.meta.requiresAuth)) {
    if (store.getters['user/isLoggedIn']) {
      return true
    } else {
      return false
    }
  } else {
    return true
  }
}

router.beforeEach((to, from, next) => {
  if (store.getters['config/isReady']) {
    if (validateAccess(to) == true) {
      next()
    } else {
      if (store.getters['surveySettings/guestSurveys']) {
        next({
          path: '/survey/run',
        })
      } else {
        next({
          path: '/login',
          params: { nextUrl: to.fullPath }
        })
      }
    }
  } else {
    store.dispatch('config/getConfig').then(() => {
      // do auto login with token if available
      if (localStorage.getItem('survey-token')) {
        store.dispatch('user/tokenLogin').then(() => {
          next()
        }).catch(() => {
          store.dispatch('user/logout').then(() => {
            if (validateAccess(to) == true) {
              next()
            } else {
              next({
                path: '/login',
                params: { nextUrl: to.fullPath }
              })
            }
          })
        })
      } else {
        if (validateAccess(to) == true) {
          next()
        } else {
          if (store.getters['surveySettings/guestSurveys']) {
            next({
              path: '/survey/run',
            })
          } else {
            next({
              path: '/login',
              params: { nextUrl: to.fullPath }
            })
          }
        }
      }
    }).catch(function(){
      router.push('/backend_error')
      next()
    })
  }
})

export default router