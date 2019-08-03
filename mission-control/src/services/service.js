import Client from "./client";
import SpaceAPI from 'space-api';
import { SPACE_API_PROJECT, SPACE_API_URL } from "../constants";

const API = SpaceAPI.API
const cond = SpaceAPI.cond
const and = SpaceAPI.and

class Service {
  constructor() {
    this.client = new Client()
    this.spaceApi = new API(SPACE_API_PROJECT, SPACE_API_URL);
    this.db = this.spaceApi.Mongo()
  }

  setToken(token) {
    this.client.setToken(token)
  }

  setSpaceApiToken(token) {
    this.spaceApi.setToken(token)
  }

  fetchEnv() {
    return new Promise((resolve, reject) => {
      this.client.getJSON("/v1/api/config/env").then(({ status, data }) => {
        if (status !== 200) {
          reject("Internal server error")
          return
        }
        resolve(data.isProd)
      }).catch(ex => reject(ex.toString()))
    })
  }

  spaceUpRegister(name, email, pass) {
    return new Promise((resolve, reject) => {
      this.spaceApi.call("console-auth", "signup", { name, email, pass }).then(({ status, data }) => {
        if (status !== 200 || !data.result.ack) {
          reject(data.result.error)
          return
        }

        resolve({ token: data.result.token, user: data.result.user })
      }).catch(ex => reject(ex.toString()))
    })
  }

  spaceUpLogin(email, pass) {
    return new Promise((resolve, reject) => {
      this.spaceApi.call("console-auth", "login", { email, pass }).then(({ status, data }) => {
        if (status !== 200 || !data.result.ack) {
          reject(data.result.error)
          return
        }

        resolve({ token: data.result.token, user: data.result.user })
      }).catch(ex => reject(ex.toString()))
    })
  }

  login(user, pass) {
    return new Promise((resolve, reject) => {
      this.client.postJSON('/v1/api/config/login', { user, pass }).then(({ status, data }) => {
        if (status !== 200) {
          reject(data.error)
          return
        }

        resolve(data.token)
      }).catch(ex => reject(ex))
    })
  }

  fetchSpaceProfile() {
    return new Promise((resolve, reject) => {
      this.spaceApi.call("console-auth", "profile", {}).then(({ status, data }) => {
        if (status !== 200) {
          reject(data.error)
          return
        }

        resolve(data.result.user)
      }).catch(ex => reject(ex))
    })
  }

  fetchProjects() {
    return new Promise((resolve, reject) => {
      this.client.getJSON(`/v1/api/config/projects`).then(({ status, data }) => {
        if (status !== 200) {
          reject(data.error)
          return
        }
        resolve(data.projects)
      }).catch(ex => reject(ex))
    })
  }

  // saveProjectConfig upserts a project config
  saveProjectConfig(projectConfig) {
    return new Promise((resolve, reject) => {
      this.client.postJSON(`/v1/api/config/projects`, projectConfig).then(({ status, data }) => {
        if (status !== 200) {
          reject(data.error)
          return
        }
        resolve()
      }).catch(ex => {
        reject(ex)
      })
    })
  }

  deleteProject(projectId) {
    return new Promise((resolve, reject) => {
      this.client.delete(`/v1/api/config/${projectId}`).then(({ status, data }) => {
        if (status !== 200) {
          reject(data.error)
          return
        }
        resolve()
      }).catch(ex => reject(ex))
    })
  }

  fetchDeployCofig() {
    return new Promise((resolve, reject) => {
      this.client.getJSON(`/v1/api/config/deploy`).then(({ status, data }) => {
        if (status !== 200) {
          reject(data.error)
          return
        }
        resolve(data.deploy)
      }).catch(ex => reject(ex))
    })
  }

  saveDeployConfig(deployConfig) {
    return new Promise((resolve, reject) => {
      this.client.postJSON(`/v1/api/config/deploy`, deployConfig).then(({ status, data }) => {
        if (status !== 200) {
          reject(data.error)
          return
        }
        resolve()
      }).catch(ex => reject(ex))
    })
  }

  fetchOperationCofig() {
    return new Promise((resolve, reject) => {
      this.client.getJSON(`/v1/api/config/operation`).then(({ status, data }) => {
        if (status !== 200) {
          reject(data.error)
          return
        }
        resolve(data.operation)
      }).catch(ex => reject(ex))
    })
  }

  saveOperationConfig(operationConfig) {
    return new Promise((resolve, reject) => {
      this.client.postJSON(`/v1/api/config/operation`, operationConfig).then(({ status, data }) => {
        if (status !== 200) {
          reject(data.error)
          return
        }
        resolve()
      }).catch(ex => reject(ex))
    })
  }

  fetchStaticConfig() {
    return new Promise((resolve, reject) => {
      this.client.getJSON(`/v1/api/config/static`).then(({ status, data }) => {
        if (status !== 200) {
          reject(data.error)
          return
        }
        resolve(data.static)
      }).catch(ex => reject(ex))
    })
  }

  saveStaticConfig(staticConfig) {
    return new Promise((resolve, reject) => {
      this.client.postJSON(`/v1/api/config/static`, staticConfig).then(({ status, data }) => {
        if (status !== 200) {
          reject(data.error)
          return
        }
        resolve()
      }).catch(ex => reject(ex))
    })
  }

  requestPayment(email, name) {
    return new Promise((resolve, reject) => {
      this.spaceApi.call('space-site', 'request-payment', { email: email, name: name }, 5000)
        .then(({ status, data }) => {
          if (status !== 200 || !data.result.ack) {
            reject()
            return
          }

          resolve()
        }).catch(ex => reject(ex))
    })
  }

  fetchCredits(userId) {
    return new Promise((resolve, reject) => {
      this.db.getOne('credits').where(cond("userId", "==", userId)).apply().then(({ status, data }) => {
        if (status !== 200) {
          reject(data.error)
          return
        }
        resolve(data.result.credits)
      }).catch(ex => reject(ex))
    })
  }

  fetchBilling(userId, month, year) {
    return new Promise((resolve, reject) => {
      this.db.get('billing').where(
        and(
          cond("userId", "==", userId),
          cond("month", "==", month),
          cond("year", "==", year),
          cond("kind", "==", 0),
          cond("applied", "==", false)
        )
      ).apply().then(({ status, data }) => {
        if (status !== 200) {
          reject(data.error)
          return
        }
        resolve(data.result)
      }).catch(ex => reject(ex))
    })
  }

  execSpaceAPI(projectId, code, token) {
    return new Promise((resolve, reject) => {
      const url = process.env.NODE_ENV !== "production" ? "http://localhost:4122" : undefined
      const api = new API(projectId, url)
      if (token) {
        api.setToken(token)
      }
      const cond = SpaceAPI.cond
      const and = SpaceAPI.and
      const or = SpaceAPI.or
      try {
        const promise = eval(code)
        if (!promise || !promise.then) {
          reject("Not a valid Space Cloud API call")
        }
        promise.then(res => resolve(res)).catch(ex => reject(ex.toString()))
      } catch (error) {
        reject(error.toString())
      }
    })

  }
}

export default Service