from locust import HttpUser, task, between
import jwt
import gevent.lock
from datetime import datetime, timedelta

# Create a counter and a lock
user_counter = 0
user_counter_lock = gevent.lock.Semaphore()

class UserBehavior(HttpUser):
    wait_time = between(1, 2.5)
    host = "http://localhost:80"
    jwt_token = None
    user_id = None
    task_id = None
    user_count = 0
    task_counter = 0

    def on_start(self):
        global user_counter
        with user_counter_lock:
            self.user_count = user_counter
            user_counter += 1
        self.task_counter = 0
        self.register_and_login()

    def register_and_login(self):
        first_name = f"FN{self.user_count}"
        last_name = f"LN{self.user_count}"
        email = f"somemail{self.user_count}@mail.com"
        password = f"pass{self.user_count}"
        response = self.client.post("/api/v1/auth/register", json={"first_name": first_name,"last_name": last_name,"email": email, "password": password})
        if response.status_code != 201:
            print("Failed to create user")
        else:
            print("User created successfully with email: {}".format(email))
            response = self.client.post("/api/v1/auth/login", json={"email": email, "password": password})
            if response.status_code == 200:
                self.jwt_token = response.json()['token']
                print("Got token: {}".format(self.jwt_token))
                decoded_token = jwt.decode(self.jwt_token, options={"verify_signature": False})
                self.user_id = decoded_token['id']
            else:
                print("Failed to login")

    @task(4)
    def get_tasks(self):
        print("Getting tasks for user_id: {}".format(self.user_id))
        
        headers = {'Authorization': f'Bearer {self.jwt_token}', 'X-User-Id': str(self.user_id)}
        response = self.client.get("/api/v1/tasks", headers=headers)
        if response.status_code != 200:
            print(f"Failed to get tasks. Status code: {response.status_code}, Response: {response.text}")
        else:
            print(f"Tasks for user_id {self.user_id}: {response.json()}")
        
    @task(3)
    def add_task(self):
        print("Add task for user_id: {}".format(self.user_id))
        
        self.task_counter += 1
        task_name = f"task_{self.user_id}_{self.task_counter}"
        deadline = (datetime.now() + timedelta(hours=24)).isoformat() + "Z"
        
        headers = {'Authorization': f'Bearer {self.jwt_token}', 'X-User-Id': str(self.user_id)}
        response = self.client.post("/api/v1/tasks", json={"name": task_name,"deadline": deadline,"completed": "false"}, headers=headers)
        if response.status_code == 200:
            self.task_id = response.json()['id']
        else:
            print(f"Failed to get tasks. Status code: {response.status_code}, Response: {response.text}")

    @task(2)
    def modify_task(self):
        print("Complete task for user_id: {}".format(self.user_id))
        
        if self.task_id:
            headers = {'Authorization': f'Bearer {self.jwt_token}', 'X-User-Id': str(self.user_id)}
            self.client.put(f"/api/v1/tasks/{self.task_id}", json={"completed": True}, headers=headers)
            

    @task(1)
    def delete_task(self):
        print("Delete task for user_id: {}".format(self.user_id))
        
        if self.task_id:
            headers = {'Authorization': f'Bearer {self.jwt_token}', 'X-User-Id': str(self.user_id)}
            self.client.delete(f"/api/v1/tasks/{self.task_id}", headers=headers)