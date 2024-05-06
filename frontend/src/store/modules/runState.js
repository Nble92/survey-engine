import router from '../../router.js'

// initial state
const state = {
  error: false,
  errorText: '',
  
  // current view key
  currentView: 'survey',
  
  // survey & existing sessionDataID if any
  surveyName: '',
  surveySessionDataID: '',
  
  // dataKey and referenceID
  surveyDataKey: '',
  surveyReferenceID: '',
  
  // existing surveyData if any
  surveyData: {},
  
  // allow reset action
  allowReset: false,
  
  lastRunLogic: '',
  lastReset: ''
}

// getters
const getters = {}

// actions
const actions = {
  resetState({commit}) {
    commit('reset_state')
  },
  setSurveyType({commit}, type) {
    commit('set_survey_type', type)
  },
  statusChange({commit, dispatch, rootState}) {
    if (rootState.surveySettings.surveyType == 'linear') {
      if (!rootState.surveySession.surveyState.entrySurveyComplete) {
        commit('update_state', {surveyDataKey: 'entry', surveyName: rootState.surveySettings.config.entrySurvey, currentView: 'survey'})
      } else if (!rootState.surveySession.surveyState.completed) {
        dispatch("surveySession/updateSurveySessionState", {name: 'completed', value: true}, {root:true})
        commit('reset_state')
      } else {
        dispatch('surveySession/resetSessionState', false, {root:true})
        commit('reset_state')
        router.push('/survey/completed')
      }
    } else if (rootState.surveySettings.surveyType == 'calendar') {
      if (!rootState.surveySession.surveyState.entrySurveyComplete) {
        commit('update_state', {surveyDataKey: 'entry', surveyName: rootState.surveySettings.config.entrySurvey, currentView: 'survey', allowReset: true})
      } else if (rootState.surveySession.surveyState.completed && !localStorage.getItem('survey-editMode')) {
        dispatch('surveySession/resetSessionState', false, {root:true})
        commit('reset_state')
        router.push('/survey/completed')
      } else {
        commit('update_state', {currentView: 'calendar', allowReset: false})
      }
    }
  },
  updateState({commit}, state) {
    commit('update_state', state)
  }
}

// mutations
const mutations = {
  reset_state(state) {
    state.error = false
    state.errorText = ''
    
    state.currentView = 'survey'
    
    state.surveyName = ''
    state.surveySessionDataID = ''
    
    state.surveyDataKey = ''
    state.surveyReferenceID = ''
    
    state.surveyData = {}
        
    state.allowReset = false
    
    state.lastRunLogic = ''
    
    state.lastReset = Date.now()
  },
  set_survey_type(state, type) {
    state.surveyType = type
  },
  update_state(state, newState) {
    // eslint-disable-next-line
    console.log('runState update_state', newState)
    for (var key in newState) {
      state[key] = newState[key]
    }
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}