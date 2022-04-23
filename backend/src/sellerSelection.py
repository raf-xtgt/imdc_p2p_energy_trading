"""
First-price reverse auction mechanism-based smart contract to allow a buyer to select the bids having the lowest price value.
Double auction mechanism to handle the amount of traded energy and energy pricing.

Selection of the Service Provider.
Steps:
1. Buyer makes a buy request
2. Seller makes a bid on the buy request
3. After 30 mins, the request is closed
4. Selection of a seller for a buyer
    * dict to pair buyer and sellers
    * loop over all buy requests
    * for each buyRequest:
        * loop through all sellRequests(bids) on current buyRequest
        * for each buyRequest:
            * add buyerId, sellerId, buyRequest and sellRequest data in the dict.

    * array to hold all the prices that the seller wants
    * loop through each Buy and its corresponding bids
    * for each  {Buy and its corresponding bid}
        * compare the prices:
            * if payable_by_buyer >= price_by_seller:
                calculate opt_price_payable_by_buyer()
                calculate opt_price_recivable_by_seller()
                append(opt_price_recivable_by_seller())
        
            else:
                notify buyer nohting found
    return min[selling_request with min opt_price_payable_by_seller]

"""
from doubleAuction import *

def initMatchmaking(client):
    final_list = getOrderData(client)
    # add the matchmaking data in the db
    cluster=client["IMDC-p2p-energy"]
    for data in final_list:
        print("Running double auction")
        transaction = initDoubleAuction(client, data)
        
        # store the double auction in the database
        collection = cluster.transactions
        transactionId = collection.insert_one(transaction)
        myquery = { "tId": "" }
        newvalues = { "$set": { "tId": str(transactionId.inserted_id) } }
        result = collection.update_one(myquery, newvalues)
        print("Transaction stored on database successfully")
    return "Success"

"""
Retrieve all the transactions that have not undergone the double auction.
"""
def getOrderData(client):
    cluster=client["IMDC-p2p-energy"]
    buyReqColl = cluster.buyRequests
    sellReqColl = cluster.sellRequests
    buyRequests = list(buyReqColl.find())
    sellRequests = list(sellReqColl.find())
    req_arr = []
    all_seller_ids = []
    for item in buyRequests:
        # request_id = item['reqid']
        # for all requests that have not undergone auction
        if item['auctioned'] == False:
            all_bids = [] # list to hold all bids made on current buy request
            req_dict = {} # request dictionary
            buy_req_id = item['reqid']
            buyer_id = item['buyerid']
            print("Buy Request:", buy_req_id)
            for bid in sellRequests:
                # if bid made on current request
                #print("Current sell request for bid", bid['buyreqid'])
                if bid['buyreqid'] == buy_req_id:
                    seller_id = bid['sellerid']                
                    seller_acc = getAccData(seller_id, client)
                    seller_fiat_balance = seller_acc[0]['fiatbalance']
                    seller_energy_balance = seller_acc[0]['energybalance']
                    bid_obj = {
                        "sellerId": seller_id,
                        "sellerReceivable": bid['fiatamount'], # amount of money seller wants to receive
                        "sellerFiatBalance": seller_fiat_balance,
                        "sellerEnergySupply": bid['energyamount'], # amount of energy seller wants to supply 
                        "sellerEnergyBalance": seller_energy_balance,
                    }
                    all_bids.append(bid_obj)
                    if bid['sellerid'] not in all_seller_ids:
                        all_seller_ids.append(bid['sellerid'])

            opt_payable = item['fiatamount'] * len(all_bids)
            buyer_acc = getAccData(buyer_id, client)
            buyer_fiat_balance = buyer_acc[0]['fiatbalance']
            buyer_energy_balance = buyer_acc[0]['energybalance']
            req_dict[buy_req_id]= {
                'buyerId': buyer_id,
                'buyerPayable': item['fiatamount'], # amount of money buyer willing to pay,
                'buyerEnergyDemand': item['energyamount'], # amount of energy the buyer wants,
                "PAY": opt_payable, # amount to be used as Pay(E) for Optimal Bid Price
                'bids':all_bids,
                'buyerFiatBalance': buyer_fiat_balance,
                'buyerEnergyBalance': buyer_energy_balance
            }

            # update the buy request to be set as auctioned
            req_arr.append(req_dict)
            myquery = { "reqid": buy_req_id }
            newvalues = { "$set": { "auctioned": True } }
            buyReqColl.update_one(myquery, newvalues)
            
        

    return req_arr


def getAccData(uId, client):
    """
    Given the userId return the user's account data
    """
    cluster=client["IMDC-p2p-energy"]
    collection = cluster.accountBalance
    account = list(collection.find({'userid': uId}, {"_id":0}))
    return account
