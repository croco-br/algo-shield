import './assets/main.css'
// Import Font Awesome CSS
import '@fortawesome/fontawesome-free/css/all.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { library } from '@fortawesome/fontawesome-svg-core'
import { 
  faChartLine, 
  faExchangeAlt, 
  faTasks, 
  faUsersCog, 
  faPalette,
  faMagnifyingGlass,
  faBell,
  faChevronDown,
  faUser,
  faSignOutAlt,
  faChevronLeft,
  faChevronRight,
  faWindowMaximize,
  faImage,
  faStar,
  faHeading,
  faHashtag,
  faSave,
  faRotateLeft,
  faPlus,
  faPlusCircle,
  faUserSlash,
  faUserCheck,
  faShield,
  faXmark,
  faPencil,
  faTrash,
  faAlignLeft,
  faClose,
  faCheckCircle,
  faXmarkCircle,
  faEnvelope,
  faKey,
  faFilter,
  faDownload,
  faPaperPlane,
  faTriangleExclamation,
  faArrowUp,
  faCheck,
  faCircleInfo,
  faFlag,
  faCircleExclamation
} from '@fortawesome/free-solid-svg-icons'

import App from './App.vue'
import router from './router'
import { vuetify } from './plugins/vuetify'

// Add icons to library
library.add(
  faChartLine,
  faExchangeAlt,
  faTasks,
  faUsersCog,
  faPalette,
  faMagnifyingGlass,
  faBell,
  faChevronDown,
  faUser,
  faSignOutAlt,
  faChevronLeft,
  faChevronRight,
  faWindowMaximize,
  faImage,
  faStar,
  faHeading,
  faHashtag,
  faSave,
  faRotateLeft,
  faPlus,
  faPlusCircle,
  faUserSlash,
  faUserCheck,
  faShield,
  faXmark,
  faPencil,
  faTrash,
  faAlignLeft,
  faClose,
  faCheckCircle,
  faXmarkCircle,
  faEnvelope,
  faKey,
  faFilter,
  faDownload,
  faPaperPlane,
  faTriangleExclamation,
  faArrowUp,
  faCheck,
  faCircleInfo,
  faFlag,
  faCircleExclamation
)

const app = createApp(App)

app.component('font-awesome-icon', FontAwesomeIcon)

app.use(createPinia())
app.use(router)
app.use(vuetify)

app.mount('#app')
