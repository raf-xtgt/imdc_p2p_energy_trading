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

#    def on_start(self):
#        self.client.get("/GetBuyRequests")
 
   
   # the locust creates a dummy buy request first
   @task(1)
   def createBuyRequests(self):
       self.client.put("/CreateBuyRequest", order)

    # then makes a sell request on the buy data
   @task(2)
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















