const fetchJSON = (url, options) => {
  return new Promise((resolve, reject) => {
    url = 'http://localhost:8080' + url
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
}

export default Client
