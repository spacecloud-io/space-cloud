import Client from "./client";
import { API, cond, and } from 'space-api';
import { SPACE_API_PROJECT, SPACE_API_URL } from "../constants";

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

  oauthLogin(uid) {
    return new Promise((resolve, reject) => {
      this.spaceApi.call("console-auth", "oauthComplete", { uid: uid }).then(({ status, data }) => {
        if (status !== 200 || !data.result.ack) {
          reject(data.error)
          return
        }

        resolve({ token: data.result.token, user: data.result.user })
      }).catch(ex => reject(ex))
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
}

export default Service