import Vue from 'vue'

var configURL = location.protocol + '//' + location.hostname + ':' + location.port + '/config'


// initial state
const state = {
  status: 'pending',
  APIURL: '',
  Name: 'unknown',
  friendlyName: '',
  welcomeText: '',
  contactName: '',
  contactEmail: ''
}

// getters
const getters = {
  isReady: state => state.status == 'ready'
}

// actions
const actions = {
  getConfig({commit, dispatch}) {
    return new Promise((resolve, reject) => {
      Vue.axios.get(configURL)
      .then(response => {
        Vue.axios.defaults.baseURL =  '//' + location.hostname + ':' + location.port + response.data.APIURL

        commit('setConfig', response)
        
        // fetch settings/surveys
        dispatch('surveySettings/fetchState', {}, {root:true}).then(() => {
          dispatch('surveys/fetchSurveys', {}, {root:true}).then(() => {
            resolve(response)
          }).catch(function (error) {
            reject(error)
          })
        }).catch(function (error) {
          reject(error)
        })
      }).catch(function (error) {
        reject(error)
      })
    })
  }
}

// mutations
const mutations = {
  setConfig(state, config) {
    state.status = 'ready'
    state.APIURL = config.data.APIURL
    state.Name = config.data.Name
    state.friendlyName = config.data.friendlyName
    state.welcomeText = config.data.welcomeText
    state.contactName = config.data.contactName
    state.contactEmail = config.data.contactEmail
    //console.log('app config:', state)
  }
  
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}