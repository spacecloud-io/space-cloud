import service from "../index";
import store from "../store";
import history from "../history";
import { set, get, reset } from "automate-redux";
import { cloneDeep } from "lodash"
import { adjustConfig, unAdjustConfig, notify } from "./helpers"
import { defaultDbConnectionStrings } from "../constants"

const generateProjectId = (projectName) => {
  return projectName.toLowerCase().split(" ").join("-") + "-" + Math.random().toString(36).substring(6)
}

const getConnString = (dbType) => {
  const connString = defaultDbConnectionStrings[dbType]
  return connString ? connString : "localhost"
}

export const fetchProjects = () => {
  return new Promise((resolve, reject) => {
    service.fetchProjects().then(({ projects }) => {
      store.dispatch(set("projects", projects))
      resolve()
    }).catch(({ error }) => {
      console.log(error)
      notify("error", "Error", "Could not fetch projects")
      reject()
    })
  })
}

export const deleteProject = (projectId) => {
  service.deleteProject(projectId).then(() => {
    const updatedProjects = get(store.getState(), "projects", []).filter(project => project.id !== projectId)
    store.dispatch(set("projects", updatedProjects))
    const selectedProject = get(store.getState(), "config.id")
    if (selectedProject === projectId) {
      store.dispatch(reset("config"))
      store.dispatch(reset("savedConfig"))
      if (updatedProjects.length) {
        openProject(updatedProjects[0].id)
        return
      }
      history.push("/mission-control/welcome")
    }
  })
}

export const openProject = (projectId) => {
  history.push(`/mission-control/projects/${projectId}`)
  loadProject(projectId)
}

export const loadProject = (projectId) => {
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

export const login = (user, pass) => {
  service.login(user, pass).then(({ projects, token }) => {
    localStorage.setItem("token", token)
    if (projects.length === 0) {
      history.push(`/mission-control/welcome`)
      return
    }

    const projectId = projects[0].id
    history.push(`/mission-control/projects/${projectId}`)
    loadProject(projectId)
  }).catch(({ error }) => {
    notify("error", "Error in login", error)
  })
}

export const createProject = (name, dbType) => {
  const projectId = generateProjectId(name)
  const projectConfig = {
    name: name,
    id: projectId,
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
  }
  service.saveProjectConfig(projectConfig).then(() => {
    const updatedProjects = [...get(store.getState(), "projects", []), projectConfig]
    store.dispatch(set("projects", updatedProjects))
    history.push(`/mission-control/projects/${projectId}`)
    const adjustedConfig = adjustConfig(projectConfig)
    store.dispatch(set("config", adjustedConfig))
    store.dispatch(set("savedConfig", cloneDeep(adjustedConfig)))
  }).catch(({ error }) => {
    notify("error", "Error creating project", error)
  })
}

export const saveConfig = () => {
  const config = get(store.getState(), "config")
  const result = unAdjustConfig(config)
  if (!result.ack) {
    if (Object.keys(result.errors.crud).length) {
      let errorDesc = ''
      Object.keys(result.errors.crud).forEach(db => {
        errorDesc = `${errorDesc} ${db} - (${result.errors.crud[db].join(", ")})`
      })
      notify("error", 'Error in CRUD config', errorDesc, 0)
    }

    if (result.errors.functions.length) {
      notify("error", 'Error in Functions config', `Services - ${result.errors.functions.join(", ")}`, 0)
    }

    if (result.errors.fileStore.length) {
      notify("error", 'Error in File Storage config', `Rules - ${result.errors.fileStore.join(", ")}`, 0)
    }
    return
  }
  service.saveProjectConfig(result.config).then(() => {
    notify("success", 'Success', 'Config saved successfully')
    store.dispatch(set("savedConfig", store.getState().config))
  }).catch(({ error }) => {
    notify("error", 'Error saving config', error)
  })
}