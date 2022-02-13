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

def initMatchmaking(client):
    final_list = getOrderData(client)
    # add the matchmaking data in the db
    cluster=client["IMDC-p2p-energy"]
    # collection to store selected sellers
    collection = cluster.selectedSellers
    for data in final_list:
        result = collection.insert_one(data)
        print("Added optimal seller for buy requests", result.inserted_id)
    return "Success"


def getOrderData(client):
    cluster=client["IMDC-p2p-energy"]
    buyReqColl = cluster.buyRequests
    sellReqColl = cluster.sellRequests
    buyRequests = buyReqColl.find()
    sellRequests = list(sellReqColl.find())
    req_arr = []
    for item in buyRequests:
        # for all open requests
        if item['requestclosed'] == False:
            req_dict = {} # request dictionary
            buy_req_id = item['reqid']
            all_bids = [] # list to hold all bids made on current buy request
            #print("Current buy request:", buy_req_id)
            for bid in sellRequests:
                # if bid made on current request
                #print("Current sell request for bid", bid['buyreqid'])
                if bid['buyreqid'] == buy_req_id:
                    #print("bid on request found")
                    bid_obj = {
                        "sellerId": bid['sellerid'],
                        "sellerReceivable": bid['fiatamount'], # amount of money seller wants to receive
                        "sellerEnergySupply": bid['energyamount'] # amount of energy seller wants to supply 
                    }
                    all_bids.append(bid_obj)

            req_dict[buy_req_id]= {
                'buyerId':item['buyerid'],
                'buyerPayable': item['fiatamount'], # amount of money buyer willing to pay,
                'buyerEnergyDemand': item['energyamount'],
                'bids':all_bids
            }

            req_arr.append(req_dict)
    
    # for trns in req_arr:
    #     print(trns)
    #     print("\n")
    # print() 
    #print(getOptimalPrices(req_arr))
    return getOptimalPrices(req_arr)

# to get optimal payable and receivable for each transaction
def getOptimalPrices(transactions):
    #print("running")
    all_seller_ids = [] # list of all sellers that are involved in the current transaction pool
    all_request_ids = [] # list of all buy request ids to be used to find the cheapest 
    # list to hold the buyer total payable amount for each seller
    all_payable = []
    for k in range(0, len(transactions)):
        # for each transaction
        for key in transactions[k]:
            obj = transactions[k][key]
            #print(obj)

            # optimal payable is the summation of the payable amount for all bids
            # since the payable to the sellers are same just multiply by number of bids on the buy request
            buyer_info = {
                "buyerId":obj['buyerId'],
                "buyerReq":key,
                "optimalPayable": obj['buyerPayable'] * len(obj['bids'])
            }
            all_payable.append(buyer_info)

            # if key not in all_request_ids:
            #     all_request_ids.append(key)

            # make a list of all sellerIds
            all_bids = obj['bids']
            for i in range(len(all_bids)):
                bid = all_bids[i]
                seller = bid['sellerId']
                if seller not in all_seller_ids:
                    all_seller_ids.append(seller)
    
    #the seller total receivable amount for each seller
    all_receivable = optReceivable(transactions, all_seller_ids)

    print("Total Payable for buyers")
    for i in all_payable:
        print(i)
        print("\n")
    
    print("Total Receivable for sellers")
    for i in all_receivable:
        print(i)
        print("\n")
    return selectSeller(all_receivable)


# to get optimal fiat money receivable by seller
def optReceivable(transactions, allSellerIds):
    # list to hold the seller total receivable amount for each seller
    all_receivable = []
    for i in range(len(allSellerIds)):
        current_seller = allSellerIds[i]
        total_receivable = 0 # receivable on all the bids that the seller made(this is the reward for OSP)
        receivable = 0 # receivable on the current bid
        buyReq = 0 # id of the buy request on which the bid is made
        # summation of the money they received from all the bids they made on the requests
        for k in range(0, len(transactions)):
            # for each transaction
            for key in transactions[k]:
                obj = transactions[k][key]
                all_bids = obj['bids']
                for i in range(len(all_bids)):
                    bid = all_bids[i]
                    seller = bid['sellerId']
                    # each seller can only make one bid per a specific buy request
                    if seller == current_seller:
                        receivable = bid['sellerReceivable']
                        total_receivable += receivable
                        buyReq = key
                        break # break because seller found

                    
        sell_info = {
            # sellerId and corresponding receivable amount
            "sellerId": current_seller,
            "buyRequest":buyReq,
            "receivableOnBid": receivable,
            "optimalReceivable": total_receivable
        }
        all_receivable.append(sell_info)
    return all_receivable

# get the seller with the minimum total receivable per request
def selectSeller(all_receivable):
    sorted_list = sorted(all_receivable, key=lambda d: d['receivableOnBid'])
    final_list = []
    # print("Sorted list")
    seenIds = []
    for i in sorted_list:
        reqId = i['buyRequest']
        if reqId not in seenIds:
            seenIds.append(reqId)
            # since it is sorted, we start from lowest to highest
            # Best seller per request
            data = {
                "buyRequestId": reqId,
                "sellerId": i['sellerId'],
                "fiatPayable": i["receivableOnBid"] # amomut buyer pays and receiver receives
            }
            final_list.append(data)
    
    print("Final list")
    print(final_list)
    return final_list





"""
# to get optimal fiat money payable by buyer
def optPayable(obj):
    # optimal payable by buyer
    buyer_payable = obj['buyerPayable']
    opt_payable = 0
    all_bids = obj['bids']
    for i in range(len(all_bids)):
        # selling_price = all_bids[i]['sellerReceivable']
        opt_payable += buyer_payable
    return opt_payable
"""