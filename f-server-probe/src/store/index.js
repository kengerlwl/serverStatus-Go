// src/store/index.js

import { createStore } from 'vuex'

const store = createStore({
  state: {
    globalVariable: 'Hello, World!',
    backendServerUrl: 'http://'+ window.location.host
    // backendServerUrl: 'http://localhost:8080'
  },
  mutations: {
    setGlobalVariable(state, newValue) {
      state.globalVariable = newValue
    }
  },
  actions: {
    updateGlobalVariable({ commit }, newValue) {
      commit('setGlobalVariable', newValue)
    }
  },
  getters: {
    getGlobalVariable: (state) => state.globalVariable
  }
})

export default store
