import Vue from 'vue'

// initial state
const state = {
  status: 'new',
  sessionID: localStorage.getItem('survey-sessionID') || '',
  surveyType: '',
  surveyCreated: '',
  surveyState: {dataKey: '', completed: false, dataRefs: [], notes: []},
  isTest: false,
  referralCode: '',
  editMode: localStorage.getItem('survey-editMode') || false,
}

// getters
const getters = {}

// actions
const actions = {
  fetchSessionState({commit, state}) {
    return new Promise((resolve, reject) => {
      if (state.sessionID) {
        commit('fetch_request')
      
        Vue.axios.get('/survey/session/' + state.sessionID)
        .then(response => {
          commit('load_state', response.data)
        
          resolve(response)
        }).catch(function(error) {
          commit('fetch_error')
        
          reject(error)
        }) 
      } else {
        resolve('no session')
      }
    })
  },
  newSurveySession({commit}, session) {
    return new Promise((resolve, reject) => {
      commit('reset_state')
      
      localStorage.removeItem('survey-sessionID')
      localStorage.removeItem('survey-editMode')
      
      commit('fetch_request')
      Vue.axios.post('/survey/newSession', {surveyType: session.surveyType, username: session.username, isTest: session.isTest, referralCode: session.referralCode})
      .then(response => {
        const sessionID = response.data.sessionID
        
        if (sessionID) {
          localStorage.setItem('survey-sessionID', sessionID)
        }
        
        commit('load_state', response.data)
        
        resolve(response)
      }).catch(function(error) {
        commit('fetch_error')
      
        reject(error)
      }) 
    })
  },
  
  editSurveySession({commit}, sessionID) {
    return new Promise((resolve, reject) => {
      localStorage.removeItem('survey-sessionID')
      
      if (sessionID) {
        commit('fetch_request')
      
        Vue.axios.get('/survey/session/' + sessionID)
        .then(response => {
          localStorage.setItem('survey-sessionID', sessionID)
          localStorage.setItem('survey-editMode', true)
          
          commit('load_state', response.data)
          commit('edit_mode')
          
          resolve(response)
        }).catch(function(error) {
          commit('fetch_error')
        
          reject(error)
        }) 
      } else {
        resolve('no session')
      }
    })
  },
  
  updateSurveySessionState({commit}, args) {
    return new Promise((resolve, reject) => {
      if (state.sessionID) {
        commit('set_state', args)
        
        // eslint-disable-next-line
        console.log('updateSurveySessionState', state.surveyState)
        
        Vue.axios.post('/survey/session/' + state.sessionID, {surveyState: state.surveyState})
        .then(response => {
          commit('load_state', response.data)
        
          resolve(response)
        }).catch(function(error) {
          reject(error)
        }) 
      } else {
        resolve('no session')
      }
    })
  },
  resetSessionState({commit}) {
    return new Promise((resolve) => {
      commit('reset_state')
      
      localStorage.removeItem('survey-sessionID')
      localStorage.removeItem('survey-editMode')
      
      resolve()
    })
  },
  
  newSurveySessionData({commit}, sessionData) {
    return new Promise((resolve, reject) => {
      Vue.axios.post('/survey/newSessionData', sessionData)
      .then(response => {
        // store simple ref for sake of commit
        commit('session_data_id', {dataKey: sessionData.dataKey, referenceID: sessionData.referenceID, surveyName: sessionData.surveyName, sessionDataID: response.data.sessionDataID})
        
        resolve(response.data)
      }).catch(function(error) {
        reject(error)
      }) 
    })
  },
  
  updateSurveySessionData({commit}, sessionData) {
    return new Promise((resolve, reject) => {
      Vue.axios.post('/survey/sessionData/' + sessionData.sessionDataID, sessionData)
      .then(response => {
        commit('session_data_id', {dataKey: sessionData.dataKey, referenceID: sessionData.referenceID, surveyName: sessionData.surveyName, sessionDataID: response.data.sessionDataID})
        
        resolve(response.data)
      }).catch(function(error) {
        reject(error)
      }) 
    })
  },
  
  deleteSurveySessionData({commit}, sessionDataID) {
    return new Promise((resolve, reject) => {
      Vue.axios.delete('/survey/sessionData/' + sessionDataID)
      .then(response => {
        commit('delete_session_data_id', sessionDataID)
        
        resolve(response.data)
      }).catch(function(error) {
        commit('delete_session_data_id', sessionDataID)
        reject(error)
      }) 
    })
  },

}

// mutations
const mutations = {
  fetch_request(state) {
    state.status = 'loading'
  },
  fetch_error(state) {
    state.status = 'error'
  },
  load_state(state, payload) {
    // eslint-disable-next-line
    console.log('surveySessionState', payload)
    
    state.sessionID = payload.sessionID
    state.surveyCreated = payload.created
    state.surveyType = payload.surveyType
    
    state.surveyState = payload.surveyState
    
    state.isTest = payload.isTest
    
    if (!payload.surveyState) {
      state.status = 'error'
    } else {
      if (payload.surveyState.completed && !state.editMode) {
        state.status = 'completed'
      } else {
        state.status = 'in_process'
      }
    }
  },
  set_state(state, args) {
    if (args) {
      // eslint-disable-next-line
      console.log('setState', args)    
      state.surveyState[args.name] = args.value 
    }
  },
  edit_mode(state) {
    state.editMode = true
  }, 
  reset_state(state) {
    state.status = 'new'
    state.sessionID = ''
    state.surveyCreated = ''
    state.surveyType = ''
    state.surveyState = {completed: false}
    state.isTest = false
    state.editMode = false
  },
  session_data_id(state, refs) {
    if (refs.sessionDataID) {
      var hasRef = false
      
      for (var id in state.surveyState.dataRefs) {
        if (state.surveyState.dataRefs[id].sessionDataID == refs.sessionDataID) {
          state.surveyState.dataRefs[id] = refs
          hasRef = true
          break
        }
      }
      
      if (!hasRef) {
        state.surveyState.dataRefs.push(refs)
      }
    }
    
    // eslint-disable-next-line
    console.log('sessionDataID', state.surveyState)
  },
  delete_session_data_id(state, sessionDataID) {
    if (sessionDataID) {
      for (var id in state.surveyState.dataRefs) {
        if (state.surveyState.dataRefs[id].sessionDataID == sessionDataID) {
          state.surveyState.dataRefs.splice(id, 1)
          break
        }
      }
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