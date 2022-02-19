
def initDoubleAuction(client):
    data = getBidData(client)
    auctionData = prepareAuctionData(data, client)
    print(auctionData)

def prepareAuctionData(data, client):
    """
    What the data json will have
    buyer Id*
    buyer account balance (fiat)
    buyer account balance (energy)
    buyer energy demand*
    PAY*

    seller Id
    seller account balance (fiat)
    seller account balance (energy)
    seller energy amount
    REW

    number of bids on the buy request (j)
    number of bids that the seller has made (i)
    """

    all_auction_data = [] # list to hold all the required data for auction
    for i in range(len(data)):
        # print(i)
        # print(data[i])
        # print("\n")
        for key in data[i]:
            
            obj = data[i][key]
            buyer_id = obj['buyerId']
            buyer_acc = getAccData(buyer_id, client)
            buyer_fiat_balance = buyer_acc[0]['fiatbalance']
            buyer_energy_balance = buyer_acc[0]['energybalance']
            
            #print("buyer account:", buyer_acc)
            seller_id = obj['selectedSeller']['sellerId']
            seller_acc = getAccData(seller_id, client)
            seller_fiat_balance = seller_acc[0]['fiatbalance']
            seller_energy_balance = seller_acc[0]['energybalance']
            bids = obj['bids']
            sorted_bids = sorted(bids, key=lambda d: d['sellerReceivable']) 
            reward = sorted_bids[0]['REWARD'] 
            sellerBids = sorted_bids[0]['bidsInvolved']# number of bids the seller is involved in      
            #print(reward)

            auction_data = {
                'buyerId': buyer_id,
                'buyerFiatBalance': buyer_fiat_balance,
                'buyerEnergyBalance': buyer_energy_balance,
                'buyerEnergyDemand': obj['buyerEnergyDemand'],
                'bids_i': len(bids), 
                'Pay': obj['PAY'],
                'sellerId': seller_id,
                'sellerFiatBalance': seller_fiat_balance,
                'sellerEnergyBalance': seller_energy_balance,
                'bids_j': sellerBids,
                'Rew': reward,
            }
            all_auction_data.append(auction_data)
    return all_auction_data




def getAccData(uId, client):
    cluster=client["IMDC-p2p-energy"]
    collection = cluster.accountBalance
    account = list(collection.find({'userid': uId}, {"_id":0}))
    return account








# the satisfaction function for a buyer (U)
def buyerSatisfaction(client):
    selectedSellers = 0

    min_energy = 20 # minimum energy that can be consumed
    demand = 0 # energy needed by the buyer
    current_energy_balance = 0 
    
    n = 0 # average charging efficiecny 
    phi = 3 # constant
    # charging willingness
    will = phi/current_energy_balance

    return

def getBidData(client):
    cluster=client["IMDC-p2p-energy"]
    collection = cluster.selectedSellers
    # get the data without the '_id:Object(..)' part
    selectedSellers = list(collection.find({},{"_id":0}))
    return selectedSellers

def optimalAllocation():
    en = 0 # optimal energy buyer needs
    sn = 0 # optimal energy seller wants to send

    bn = 0 # amount buyer wants to pay
    pn = 0 # amount seller wants to receive 
    
    i = 0 # number of buyer


def doubleAuction():
    epsilon = 0.005
    return

"""
When the user signs up, add an account balance for them with
1. current energy balance: 1200kWh
2. current fiat balance: RM 2000
"""


"""
def initDoubleAuction(client):
    data = getBidData(client)
    keys = []
    for i in range(len(data)):
        count = 0
        for key in data[i]:
            key_arr = key.split('\n')
            if count % 2 !=0:
                #print(key_arr) 
                keys.append(key_arr[0])
            count +=1
    print("keys", keys)
    print("\n")
    print("Data", data)
    # for key in keys:
    #     for i in range(len(data)):
    #         obj = data[i][key]
    #         print(obj)
    #         print("\n")
    #         break

"""