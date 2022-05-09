from locust import HttpUser, TaskSet, between, task
import json



order_data = { "buyerid": "6235eada63fc9762da8b7e54", "energyamount": 99,
        "fiatamount": 9.109,
        "requestclosed":True,
        "reqtime":"test_time",
        "reqid":"",
        "auctioned":False
    }

order = json.dumps(order_data)

class User(TaskSet):
   wait_time = between(5, 15)

   # the locust creates a dummy buy request first
   @task(1)
   def createBuyRequests(self):
       self.client.put("/CreateBuyRequest", order)


    
  
class WebsiteUser(HttpUser):
   tasks = [User]















