import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import Addresses from '../views/Addresses.vue'
import Rules from '../views/Rules.vue'
import Notifications from '../views/Notifications.vue'
import Alerts from '../views/Alerts.vue'

const routes = [
  { path: '/', name: 'Dashboard', component: Dashboard },
  { path: '/addresses', name: 'Addresses', component: Addresses },
  { path: '/rules', name: 'Rules', component: Rules },
  { path: '/notifications', name: 'Notifications', component: Notifications },
  { path: '/alerts', name: 'Alerts', component: Alerts },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router