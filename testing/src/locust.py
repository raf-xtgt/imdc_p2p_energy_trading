from locust import HttpUser, TaskSet, between, task

class User(TaskSet):
   wait_time = between(5, 15)

   def on_start(self):
       self.client.post("/GetBlockchain")
 
#    def on_stop(self):
#        self.client.post("url/name", {"data1":"data1", "data2":"2"})
 
   @task(2)
   def index(self):
       self.client.get("/GetAllUsers")
 
   @task(1)
   def test(self):
       self.client.get("/GetBuyRequests")
 
class WebsiteUser(HttpUser):
   tasks = [User]




"""
def on_start(self):
       self.client.post("/login", {"username":"admin", "password":"password"})
 
def on_stop(self):
       self.client.post("/logout", {"username":"admin", "password":"password"})
"""