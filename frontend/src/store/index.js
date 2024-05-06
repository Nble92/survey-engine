import Vue from 'vue'
import Vuex from 'vuex'

import config from './modules/config'
import user from './modules/user'
import surveySettings from './modules/surveySettings'
import surveys from './modules/surveys'
import surveySession from './modules/surveySession'
import runState from './modules/runState'

import logicFunctions from './modules/logicFunctions'
import instructions from './modules/instructions'

import images from './modules/images'

Vue.use(Vuex)

const debug = process.env.NODE_ENV !== 'production'

export default new Vuex.Store({
  plugins: [ ],
  modules: {
    config,
    user,
    surveySettings,
    surveys,
    surveySession,
    runState,
    logicFunctions,
    instructions,
    images
  },
  strict: debug
})