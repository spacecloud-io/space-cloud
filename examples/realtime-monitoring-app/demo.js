var api = new Space.API("demo", "http://localhost:8080")
var db = api.MySQL()
var loginLayout = document.getElementById("login")
var registerLayout = document.getElementById("register")
var displayLayout = document.getElementById("display")
loginLayout.classList.add("hide")
displayLayout.classList.add("hide")
function login() {
    var username = document.getElementById("username").value
    var password = document.getElementById("password").value
    api.call('login_service', 'login_func', { username, password }, 10)
        .then(res => {
            console.log(res)
            if (res.status === 200) {
                if (res.data.ack) {
                    alert("Login Success")
                    api.setToken(res.data.token)
                    displayLayout.classList.remove("hide")
                    loginLayout.classList.add("hide")
                    const onSnapshot = (docs, type) => {
                        console.log(docs)
                        var table1 = document.getElementById("table1")
                        var table2 = document.getElementById("table2")
                        var avg1 = 0
                        var avg2 = 0
                        var len1 = 0
                        var len2 = 0
                        table1.innerHTML = ''
                        table2.innerHTML = ''
                        docs.forEach(element => {
                            if (element.device == 1) {
                                table1.innerHTML += '<tr><td>'+String(element.value)+'</td></tr>'
                                avg1 += element.value
                                len1 += 1
                            }
                            if (element.device == 2) {
                                table2.innerHTML += '<tr><td>'+String(element.value)+'</td></tr>'
                                avg2 += element.value
                                len2 += 1
                            }
                        })
                        document.getElementById("length1").innerText = len1
                        if(!Number.isNaN(avg1/len1)) {
                            document.getElementById("average1").innerText = String(avg1/len1)
                        }
                        document.getElementById("length2").innerText = len2
                        if(!Number.isNaN(avg2/len2)) {
                            document.getElementById("average2").innerText = String(avg2/len2)
                        }
                    }
                    const onError = (err) => {
                        console.log('Operation failed:', err)
                    }
                    db.liveQuery("demo").subscribe(onSnapshot, onError)
                } else {
                    alert("Login Failed")
                }
            } else {
                console.log(res.error)
                alert("Login Failed")
            }
        }).catch(ex => {
            console.log(ex)
            alert("Login Failed")
        })
}
function register() {
    var username = document.getElementById("usernameReg").value
    var password = document.getElementById("passwordReg").value
    if(username.trim() == "") {
        alert("Username cannot be empty")
        return
    }
    if(password.trim() == "") {
        alert("Password cannot be empty")
        return
    }
    api.call('login_service', 'register_func', { username, password }, 10)
        .then(res => {
            console.log(res)
            if (res.status === 200) {
                if (res.data.ack) {
                    alert("Successfully Registered")
                    showLogin()
                } else {
                    alert("Register Failed")
                }
            } else {
                console.log(res.error)
                alert("Register Failed")
            }
        }).catch(ex => {
            console.log(ex)
            alert("Register Failed")
        })
}
function showLogin() {
    loginLayout.classList.remove("hide")
    registerLayout.classList.add("hide")
    displayLayout.classList.add("hide")
}
document.getElementById("loginButton").onclick = login
document.getElementById("registerButton").onclick = register
document.getElementById("showRegister").onclick = () => {
    loginLayout.classList.add("hide")
    registerLayout.classList.remove("hide")
    displayLayout.classList.add("hide")
}
document.getElementById("showLogin").onclick = showLogin
