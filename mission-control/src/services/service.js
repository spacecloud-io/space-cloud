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
        resolve({ project: data.project, token: data.token })
      }).catch(ex => {
        reject({ error: ex })
      })
    })
  }

  loadConfig(project) {
    return new Promise((resolve, reject) => {
      this.client.getJSON(`/v1/api/${project}/config`).then(({ status, data }) => {
        if (status !== 200) {
          reject({ error: data.error })
          return
        }
        resolve({ config: data.config })
      }).catch(ex => {
        reject({ error: ex })
      })
    })
  }

  saveConfig(project, config) {
    return new Promise((resolve, reject) => {
      this.client.postJSON(`/v1/api/${project}/config`, config).then(({ status, data }) => {
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