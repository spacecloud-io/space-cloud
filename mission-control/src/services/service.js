class Service {
  constructor(client) {
    this.client = client
  }

  login(user, pass) {
    return new Promise((resolve, reject) => {
      this.client.postJSON('/v1/api/config/login', { user, pass }).then(({ status, data }) => {
        if (status !== 200) {
          reject({ error: data.error })
          return
        }
        this.client.setToken(data.token)
        resolve({ projects: data.projects, token: data.token })
      }).catch(ex => {
        reject({ error: ex })
      })
    })
  }

  fetchProjects() {
    return new Promise((resolve, reject) => {
      this.client.getJSON(`/v1/api/config`).then(({ status, data }) => {
        if (status !== 200) {
          reject({ error: data.error })
          return
        }
        resolve({ projects: data.projects })
      }).catch(ex => {
        reject({ error: ex })
      })
    })
  }

  // saveProjectConfig upserts a project config
  saveProjectConfig(projectConfig) {
    return new Promise((resolve, reject) => {
      this.client.postJSON(`/v1/api/config`, projectConfig).then(({ status, data }) => {
        if (status !== 200) {
          reject({ error: data.error })
          return
        }
        resolve()
      }).catch(ex => {
        reject({ error: ex })
      })
    })
  }

  deleteProject(projectId) {
    return new Promise((resolve, reject) => {
      this.client.delete(`/v1/api/config/${projectId}`).then(({ status, data }) => {
        if (status !== 200) {
          reject({ error: data.error })
          return
        }
        resolve()
      }).catch(ex => {
        reject({ error: ex })
      })
    })
  }
}

export default Service