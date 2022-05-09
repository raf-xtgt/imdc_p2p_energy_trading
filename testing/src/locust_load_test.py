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
    "buyreqid": "6278dd07aa3d8d61648e1fbe"
}
bid = json.dumps(bid_data)


household_data = {
    "day": "",
    "average":0,
    "data": [0,0],
    "dateStr": "09-05-2022",
    "datetime":0
}
homepage = json.dumps(household_data)


class User(TaskSet):
   wait_time = between(5, 15)

   # household data for homepage
   @task(1)
   def visitHomepage(self):
       self.client.put("/AddHouseholdData", homepage)

   # create energy order
   @task(2)
   def makeOrder(self):
       self.client.put("/CreateBuyRequest", order)

    # create energy order
   @task(3)
   def makeBid(self):
       self.client.put("/CreateSellRequest", bid)

   @task(4)
   def accessMarketPage(self):
       self.client.get("/GetBuyRequests")

   @task(5)
   def getBlockchain(self):
       self.client.get("/GetBlockchain")

    
  
class WebsiteUser(HttpUser):
   tasks = [User]


