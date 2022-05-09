from locust import HttpUser, TaskSet, between, task
import json
from locust import HttpUser, TaskSet, between, task
import json


bid_data = {
    "sellerid": "624892b4dd1dd97e648a86eb",
    "energyamount": 45,
    "fiatamount": 4.59,
    "sellreqid": "",
    "reqtime": "test-time",
    "buyreqid": ""
}


class User(TaskSet):
   wait_time = between(5, 15)

    # then makes a sell request on the buy data
   @task(1)
   def createSellRequests(self):
        lines  = []
        with open('/home/rafaquat/buyer_ids.txt') as f:
            lines = f.readlines()

        
        id_strs = []
        for id in lines:
            id = id.replace("\n", "")
            id_strs.append(id)

        for id in id_strs:
            bid_data["buyreqid"] = id
            bid =  json.dumps(bid_data)
            self.client.put("/CreateSellRequest", bid)
 

    
  
class WebsiteUser(HttpUser):
   tasks = [User]

