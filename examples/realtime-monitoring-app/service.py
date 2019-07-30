import jwt
from space_api import API, COND

api = API('demo', 'localhost:4124')
SECRET = 'my_secret'
api.set_token(jwt.encode({"password": "super_secret_admin_password"}, SECRET, algorithm='HS256').decode('utf-8'))
db = api.my_sql()

service = api.service('login_service')


def login(params, auth, cb):
    response = db.get_one('demo_users').where(COND("username", "==", params["username"])).apply()
    if response.status == 200:
        res = response.result
        if res["username"] == params["username"] and res["password"] == params["password"]:
            cb('response',
               {'ack': True,
                'token': jwt.encode({"username": res["username"], "password": res["password"]}, SECRET,
                                    algorithm='HS256').decode('utf-8')})
        else:
            cb('response', {'ack': False})
    else:
        print(response.error)
        cb('response', {'ack': False})


def register(params, auth, cb):
    response = db.insert('demo_users').doc(
        {"username": params["username"], "password": params["password"]}).apply()
    if response.status != 200:
        print(response.error)
    cb('response', {'ack': response.status == 200})


service.register_func('login_func', login)
service.register_func('register_func', register)

service.start()
api.close()
