import Vue from 'vue'
import Vuex from 'vuex'
import createPersistedState from "vuex-persistedstate"
import message from './modules/message'
import group from './modules/group'
import Socket from '@/utils/socket'
import createSocketPulgin from './plugins/createSocketPlugin'

Vue.use(Vuex)

const socket = new Socket()
Vue.prototype.$socket = socket

export default new Vuex.Store({
  modules: {
    message,
    group
  },
  plugins: [
    createSocketPulgin(socket), 
    createPersistedState({
      paths: ['group', 'message.messages']
    })
  ]
})