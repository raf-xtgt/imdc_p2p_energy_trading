import json


# All buyer requests
B = [{"id":1,  # id of the buying bid
      "buyerId": 2345, # id of the user who made the request
    "bidPrice": 100,
    "energyAmount": '1000kW'},

   {"id":2,
    "buyerId": 23457,
    "bidPrice": 300,
    "energyAmount": '2000kW'},
     
   {"id":3,
    "buyerId": 2345123,
    "bidPrice": 50,
    "energyAmount": '100kW'}]

# All seller bids on the buyer's request
# sort by time later on, so first come first served basis
S = [{"id":6, # id of the sell bid
    "buyId":3, # id of the buyer request, that the seller is making a bid on
    "sellerId": 3335, # id of the user who wants to sell energy
    "bidPrice": 100,
    "energyAmount": '1000kW'},

   {"id":7,
    "buyId":2,
    "sellerId": 3332,
    "bidPrice": 300,
    "energyAmount": '2000kW'},
     
   {"id":12,
    "buyId":1,
    "sellerId": 3331,
    "bidPrice": 5000,
    "energyAmount": '100kW'},
     
     {"id":21,
    "buyId":1,
    "sellerId": 33334,
    "bidPrice": 50,
    "energyAmount": '100kW'}
     ]

def selectSeller(buyers, sellers):
    #print(len(buyers))
    #print(len(sellers))
    bid_pairs = []
    for i in range(len(buyers)):
        buyer = buyers[i]
        buy_price =  buyer["bidPrice"]
        buy_id = buyer["id"]
        # loop through seller bids made on the buyer bid
        sell_ids = []
        sell_bids = []
        data = {}
        min_sell_bid = float('inf')
        selected_seller = "None"
        for j in range(len(sellers)):
            seller = sellers[j]
            if seller["buyId"] == buy_id:
                sell_price = seller["bidPrice"]
                # first come first served basis
                # another seller with same price will not be accepted       
                if sell_price < min_sell_bid:
                    min_sell_bid = sell_price
                    selected_seller = seller

        data[str(buy_id)] = selected_seller
        json_data = json.dumps(data)
        json_obj = json.loads(json_data)
        bid_pairs.append(json_obj)


    for k in bid_pairs:
        print(k)

    print("\n")
    return bid_pairs


print(selectSeller(B,S))
                


