import time
import random
import jwt
from space_api import API

api = API('demo', 'localhost:4124')
SECRET = 'my_secret'
api.set_token(jwt.encode({"password": "super_secret_password"}, SECRET, algorithm='HS256').decode('utf-8'))
db = api.my_sql()

for i in range(10):
    response = db.insert('demo').doc({"id": i+100, "device": 2, "value": random.randint(11, 20)}).apply()
    if response.status == 200:
        print("Sent")
    else:
        print(response.error)
    time.sleep(2)
