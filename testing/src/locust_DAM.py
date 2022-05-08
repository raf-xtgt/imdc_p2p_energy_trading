from locust import HttpUser, TaskSet, between, task
import json



class User(TaskSet):
   wait_time = between(5, 15)

   
   @task(1)
   def runDoubleAuction(self):
       # the locust creates a dummy buy request first
       data = json.dumps("run double auction")
       self.client.put("/RunDoubleAuction", data)



    
  
class WebsiteUser(HttpUser):
   tasks = [User]















