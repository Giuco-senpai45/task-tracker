from locust import HttpUser, task, between
import jwt

class UserBehavior(HttpUser):
    wait_time = between(1, 2.5)
    host = "http://localhost:80"
    jwt_token = None
    user_id = None

    def on_start(self):
        self.login()

    def login(self):
        response = self.client.post("/api/v1/auth/login", json={"email": "some@mail.com", "password": "pass"})
        if response.status_code == 200:
            self.jwt_token = response.json()['token']
            decoded_token = jwt.decode(self.jwt_token, options={"verify_signature": False})
            self.user_id = decoded_token['id']
        else:
            print("Failed to login")

    @task(2)
    def get_tasks(self):
        headers = {'Authorization': f'Bearer {self.jwt_token}', 'X-User-Id': str(self.user_id)}
        self.client.get("/api/v1/tasks", headers=headers)

    @task(1)
    def index(self):
        self.client.get("/")