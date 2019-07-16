import { increment, decrement, set, get } from "automate-redux"
import { cloneDeep } from "lodash"
import { notification } from "antd"
import service from "./index"
import history from "./history"
import store from "./store"
import { defaultDbConnectionStrings } from "./constants"

export const openProject = (projectId) => {
  history.push(`/mission-control/projects/${projectId}`)
  const projects = get(store.getState(), "projects", [])
  const config = projects.find(project => project.id === projectId)
  if (!config) {
    notify("error", "Error", "Project does not exist")
    return
  }

  const adjustedConfig = adjustConfig(config)
  store.dispatch(set("config", adjustedConfig))
  store.dispatch(set("savedConfig", cloneDeep(adjustedConfig)))
}

export const adjustConfig = (config) => {
  // Adjust database rules
  if (config.modules && config.modules.crud) {
    Object.keys(config.modules.crud).forEach(db => {
      if (config.modules.crud[db].collections) {
        Object.keys(config.modules.crud[db].collections).forEach(col => {
          config.modules.crud[db].collections[col] = JSON.stringify(config.modules.crud[db].collections[col], null, 2)
        })
      }
    })
  }

  // Adjust function rules
  if (config.modules && config.modules.functions && config.modules.functions.rules) {
    Object.keys(config.modules.functions.rules).forEach(service => {
      config.modules.functions.rules[service] = JSON.stringify(config.modules.functions.rules[service], null, 2)
    })
  }

  // Adjust file storage rules
  if (config.modules && config.modules.fileStore && config.modules.fileStore.rules) {
    Object.keys(config.modules.fileStore.rules).forEach(rule => {
      config.modules.fileStore.rules[rule] = JSON.stringify(config.modules.fileStore.rules[rule], null, 2)
    })
  }

  return config
}

export const unAdjustConfig = (c) => {
  let config = cloneDeep(c)
  let result = { ack: true, errors: { crud: {}, functions: [], fileStore: [] } }
  // Unadjust database rules
  if (config.modules && config.modules.crud) {
    Object.keys(config.modules.crud).forEach(db => {
      if (config.modules.crud[db].collections) {
        Object.keys(config.modules.crud[db].collections).forEach(col => {
          try {
            config.modules.crud[db].collections[col] = JSON.parse(config.modules.crud[db].collections[col])
          } catch (error) {
            result.ack = false
            result.errors.crud[db] = result.errors.crud[db] ? [...result.errors.crud[db], col] : [col]
          }
        })
      }
    })
  }

  // Unadjust function rules
  if (config.modules && config.modules.functions && config.modules.functions.rules) {
    Object.keys(config.modules.functions.rules).forEach(service => {
      try {
        config.modules.functions.rules[service] = JSON.parse(config.modules.functions.rules[service])
      } catch (error) {
        result.ack = false
        result.errors.functions.push(service)
      }
    })
  }

  // Unadjust file storage rules
  if (config.modules && config.modules.fileStore && config.modules.fileStore.rules) {
    Object.keys(config.modules.fileStore.rules).forEach(rule => {
      try {
        config.modules.fileStore.rules[rule] = JSON.parse(config.modules.fileStore.rules[rule])
      } catch (error) {
        result.ack = false
        result.errors.functions.push(rule)
      }
    })
  }

  result.config = config
  return result
}

export const notify = (type, title, msg, duration) => {
  notification[type]({ message: title, description: msg, duration: duration });
}

const generateProjectId = (projectName) => {
  return projectName.toLowerCase().split(" ").join("-")
}

const getConnString = (dbType) => {
  const connString = defaultDbConnectionStrings[dbType]
  return connString ? connString : "localhost"
}

export const generateProjectConfig = (name, dbType) => ({
  name: name,
  id: generateProjectId(name),
  secret: 'some-secret',
  modules: {
    crud: {
      [dbType]: {
        enabled: true,
        conn: getConnString(dbType),
        collections: {
          users: {
            isRealtimeEnabled: false,
            rules: {
              create: {
                rule: 'allow'
              },
              read: {
                rule: 'allow'
              },
              update: {
                rule: 'allow'
              },
              delete: {
                rule: 'allow'
              }
            }
          }
        }
      }
    },
    auth: {},
    functions: {
      enabled: true,
      broker: "nats",
      conn: "nats://localhost:4222",
      rules: {}
    },
    realtime: {
      enabled: false,
      broker: "nats",
      conn: "nats://localhost:4222"
    },
    fileStore: {
      enabled: false,
      storeType: "local",
      conn: "./",
      rules: []
    },
    static: {
      enabled: false,
      routes: [
        { prefix: "/", path: "./public" }
      ]
    }
  }
})

/* 
* Fetches does all those tasks that are to be done after a successful login with Mission Control UI
* or to be done if the person was logged in already (i.e a token was present in localStorage).
* It fetches all the projects, deployConfig, operationConfig and saves them in redux.
* It also redirects to the welcome page if the user has no projects.
* If the user has projects then it opens the last project (if provided) or the first project
* Else it redirects the user to welcome page if user has no projects.
*/
export const handleClusterLoginSuccess = (token, lastProjectId) => {
  store.dispatch(increment("pendingRequests"))
  service.setToken(token)
  Promise.all([service.fetchProjects(), service.fetchDeployCofig(), service.fetchOperationCofig()])
    .then(([projects, deployConfig, operationConfig]) => {

      // Save deploy config
      const adjustedDeployConfig = adjustConfig(deployConfig)
      store.dispatch(set("deployConfig", adjustedDeployConfig))
      store.dispatch(set("savedDeployConfig", cloneDeep(adjustedDeployConfig)))

      // Save operation config
      store.dispatch(set("operationConfig", operationConfig))

      // Save projects
      store.dispatch(set("projects", projects))
      if (projects.length === 0) {
        history.push(`/mission-control/welcome`)
        return
      }

      // Open last project
      if (!lastProjectId) {
        lastProjectId = projects[0].id
      }
      openProject(lastProjectId)
    })
    .catch(error => {
      console.log("Error", error)
      notify("error", "Error", 'Could not fetch config')
    })
    .finally(() => store.dispatch(decrement("pendingRequests")))
}

export const handleSpaceUpLoginSuccess = (token) => {
  store.dispatch(increment("pendingRequests"))
  service.setSpaceApiToken(token)
  service.fetchSpaceProfile().then(user => {
    const date = new Date()
    const month = date.getMonth()
    const year = date.getFullYear()
    Promise.all([service.fetchCredits(user._id), service.fetchBilling(user._id, month, year)]).then(([credits, billing]) => {
      store.dispatch(set("credits", credits))
      store.dispatch(set("billing", billing))
    })
    store.dispatch(set("user", user))
  }).catch(error => {
    console.log("Error", error)
    notify("error", "Error", 'Could not fetch profile')
  }).finally(() => store.dispatch(decrement("pendingRequests")))
}

export const generateId = () => {
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, function (c) {
    var r = (Math.random() * 16) | 0,
      v = c == "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

export const isUserSignedIn = () => {
  const userId = get(store.getState(), "user._id", "")
  return userId && userId.length > 0
}

export const triggerSignin = () => {
  store.dispatch(set("uiState.isSigninModalVisible", true))
}

export const openPlansPage = () => {
  const projectId = get(store.getState(), "config.id")
  history.push(`/mission-control/projects/${projectId}/plans`)
}