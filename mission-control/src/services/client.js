const fetchJSON = (url, options) => {
  if (process.env.NODE_ENV !== "production")  {
    url = "http://localhost:4122" + url
  }
  return new Promise((resolve, reject) => {
    fetch(url, options).then(res => {
      const status = res.status
      res.json().then(data => {
        resolve({ status, data })
      }).catch(ex => {
        reject(ex)
      })
    }).catch(ex => {
      reject(ex)
    })
  })
}


class Client {
  constructor() {
    this.options = {
      credentials: "include",
      headers: {
        "Content-Type": "application/json"
      }
    }
  }

  setToken(token) {
    this.options.headers.Authorization = "Bearer " + token;
  }

  getJSON(url) {
    return fetchJSON(url, Object.assign({}, this.options, { method: 'GET' }))
  }

  postJSON(url, obj) {
    return fetchJSON(url, Object.assign({}, this.options, { method: 'POST', body: JSON.stringify(obj) }))
  }

  delete(url) {
    return fetchJSON(url, Object.assign({}, this.options, { method: 'DELETE' }))
  }
}

export default Client
