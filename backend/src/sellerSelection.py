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
    buyRequests = list(buyReqColl.find())
    sellRequests = list(sellReqColl.find())
    req_arr = []
    all_seller_ids = []
    for item in buyRequests:
        # for all open requests
        if item['requestclosed'] == False:
            all_bids = [] # list to hold all bids made on current buy request
            req_dict = {} # request dictionary
            buy_req_id = item['reqid']
            buyer_id = item['buyerid']
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
                    if bid['sellerid'] not in all_seller_ids:
                        all_seller_ids.append(bid['sellerid'])

            opt_payable = item['fiatamount'] * len(all_bids)
            req_dict[buy_req_id]= {
                'buyerId': buyer_id,
                'buyerPayable': item['fiatamount'], # amount of money buyer willing to pay,
                'buyerEnergyDemand': item['energyamount'], # amount of energy the buyer wants,
                "PAY": opt_payable, # amount to be used as Pay(E) for Optimal Bid Price
                'bids':all_bids
                
            }

            req_arr.append(req_dict)
    # print("data")
    # for trns in req_arr:
    #     print(trns)
    #     print("\n")
    # print(all_seller_ids) 
    all_receivable = optReceivable(req_arr, all_seller_ids)
    # print("Total Receivable for sellers")
    # for i in all_receivable:
    #     print(i)
    #     print("\n")

    trns = updateTrnsWtihSellerRew(all_receivable, req_arr)
    # for i in trns:
    #     print(i)
    #     print("\n")
    return trns

# to get the optimal reward that each seller in the current transaction pool can receive
def optReceivable(transactions, allSellerIds):
    # list to hold the seller total receivable amount for each seller
    all_receivable = []
    cost_factor = 0.3
    for i in range(len(allSellerIds)):
        current_seller = allSellerIds[i]
        total_receivable = 0 # receivable on all the bids that the seller made(this is the reward for OSP)
        receivable = 0 # receivable on the current bid
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
                        total_receivable += ((receivable*receivable))/(4*cost_factor)
                        break # break because seller found

                    
        sell_info = {
            # sellerId and corresponding receivable amount
            "sellerId": current_seller,
            "REW": total_receivable
        }
        all_receivable.append(sell_info)
    return all_receivable


# add the reward of each seller in the bids array of each transaction
def updateTrnsWtihSellerRew(sellerRew, transactions):
    for i in range(len(sellerRew)):
        sell_obj = sellerRew[i]
        current_seller = sell_obj['sellerId']

        for k in range(len(transactions)):
            for key in transactions[k]:
                trn = transactions[k][key]
                trn_bids = trn['bids']
                sorted_bids = sorted(trn_bids, key=lambda d: d['sellerReceivable'])
                selected_seller = sorted_bids[0]
                
                for j in range(len(trn_bids)):
                    trn_bid = trn_bids[j]
                    trn_seller = trn_bid['sellerId']
                    if current_seller == trn_seller:
                        reward = sell_obj['REW']
                        trn['bids'][j] = {
                            'sellerId': current_seller, 
                            'sellerReceivable': trn_bid['sellerReceivable'], 
                            'sellerEnergySupply': trn_bid['sellerEnergySupply'],
                            'REWARD': reward
                        }
                        trn['selectedSeller'] = selected_seller
                        transactions[k][key] = trn

    # for z in transactions:
    #     print(z)
    #     print('\n')
    return transactions
    #print(transactions)
