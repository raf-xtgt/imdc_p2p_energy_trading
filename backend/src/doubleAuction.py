import numpy as np # linear algebra
from datetime import datetime
import random

def initDoubleAuction(client):
    data = getBidData(client)
    auctionData = prepareAuctionData(data, client)
    #print(auctionData)
    for i in auctionData:
        print(i)
    print('\n')
    auction = doubleAuction(auctionData)

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
        print(data[i])
        print("\n")
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
            bidFiats = [] # list to hold all the bid(fiat amount) for each seller
            for bid in bids:
                fiat = bid['sellerReceivable']
                bidFiats.append(fiat)
            
            enFromSellers = []  #list to hold all the energy amount that each seller can provide
            for bid in bids:
                en = bid['sellerEnergySupply']
                enFromSellers.append(en)
            

            reward = sorted_bids[0]['REWARD'] 
            sellerBids = sorted_bids[0]['bidsInvolved']# number of bids the seller is involved in      
            #print(reward)
            energy_price = getAvgHouseholdPrice(client)
            auction_data = {
                'buyerId': buyer_id,
                'buyerFiatBalance': buyer_fiat_balance,
                'buyerEnergyBalance': buyer_energy_balance,
                'buyerEnergyDemand': obj['buyerEnergyDemand'],
                'buyerPayable': obj['buyerPayable'],
                'bids_j': len(bids), 
                'bidFiats': bidFiats, #list of all the bids in fiat amount made on the request
                'enFromSellers': enFromSellers, # list of all the energy amount that sellers can give
                'Pay': obj['PAY'],
                'sellerId': seller_id,
                'sellerFiatBalance': seller_fiat_balance,
                'sellerEnergyBalance': seller_energy_balance,
                'sellerReceivable':sorted_bids[0]['sellerReceivable'], 
                'sellerEnergySupply': sorted_bids[0]['sellerEnergySupply'], # amount of energy the seller wants to trade
                'bids_i': sellerBids,
                'Rew': reward,
                'householdEnergyPrice': energy_price   
            }
            all_auction_data.append(auction_data)
    return all_auction_data




def getAccData(uId, client):
    """
    Given the userId return the user's account data
    """
    cluster=client["IMDC-p2p-energy"]
    collection = cluster.accountBalance
    account = list(collection.find({'userid': uId}, {"_id":0}))
    return account


def getAvgHouseholdPrice(client):

    now = datetime.now()
    date_str= str(now.strftime('%d-%m-%Y')) 
    print("Date String:", date_str)
    cluster=client["IMDC-p2p-energy"]
    collection = cluster.householdEnergyPrice
    price_data = list(collection.find({'datestr': date_str}, {"_id":0}))
    #print("Household data:", price[0]['average'])
    price = price_data[0]['average']
    return price

def doubleAuction(auctionData):
    epsilon = 0.005
    charging_eff = 0.85
    for data in auctionData:
        en_balance = data['buyerEnergyBalance']
        flag = True
        t = 0
        bn = data['buyerPayable']
        en = data['buyerEnergyDemand']
        pn = data['bidFiats']
        sn = data['sellerEnergySupply']
        j = data['bids_j'] # number of bids made by the selected seller for all open requests
        i = data['bids_i'] # number of bids on the request
        #pay = data['Pay'] # total amount if buyer had to pay all the bidders
        #rew = data['Rew'] # total amount if seller can provide all the buyers they bidded on
        pricing = data['householdEnergyPrice']
    
        # optimal energy allocation for buyer and seller
        opt_en = optimalAllocation(bn,en, pn, sn, j, en_balance)
        buyer_payable = opt_en['buyerOptEn'] * pricing 
        seller_receivable = opt_en['sellerOptEn'] *pricing
        tnbReceivable = buyer_payable - seller_receivable

        final_output = {
            "optBuyerEnergy": opt_en['buyerOptEn'],
            "buyerPayable": buyer_payable,
            "optSellerEnergy": opt_en['sellerOptEn'],
            'sellerReceivable': seller_receivable,
            'tnbReceivable': tnbReceivable

        }

        print("Optimal allocation", final_output)
        print("\n")
    
    return



def optimalAllocation(bn, en, pn, sn, bids_j, en_balance):
    """
    bn amount of money the buyer wants to pay
    en amount of energy the buyer needs
    pn amount of money the seller wants (Array)
    sn amount of energy the seller can provide (array)
    note: for pn and sn, we will loop through all the sellers and to get optimal for the selected seller
    The buyer makes a single request which is sent to all sellers, so no need
    the outer loop
    
    Returns the optimal energy that buyer and seller can trade while 
    maximising social welfare
    """
    summation = 0
    # willingness = phi/en_balance
    willingness = 100 # always willing to charge
    min_energy = 85
    n = 0.085 # average charging efficiency
    
    # optimal bid
    for j in range(bids_j):
        val = en - min_energy
        summation += val
    
    num = ((n*summation) +1)*bn
    opt_en = num/(n*willingness) # optimal energy allocated for buyer
    if opt_en < en:
        diff = en - opt_en
        opt_en = opt_en + (0.5 * diff)
    print("Optimal energy allocated for buyer, en:", opt_en)

    # cost factors c1 and c2
    c1 = 1
    c2 = 1

    opt_seller_en = (2*c1*sn) + c2
    increase = random.uniform(0.85, 0.90)
    if opt_seller_en >= opt_en:
        opt_seller_en = increase*opt_en
    elif opt_seller_en < opt_en:
        diff = opt_en - opt_seller_en
        opt_seller_en = opt_seller_en + (0.3*diff)

    print("Optimal energy that can be provided:", opt_seller_en)
    output = {
        'buyerOptEn': opt_en,
        'sellerOptEn': opt_seller_en
    }
    return output






# the satisfaction function for a buyer (U)
def buyerSatisfaction(j, optEnergyValues):

    min_energy = 80 # minimum energy that can be consumed
    demand = optEnergyValues['buyerOptEn'] # optimised energy needed by the buyer from optimalAllocation
    n = 0.085 # average charging efficiecny 
    # charging willingness
    willingness = 1
    summation = 0
    for bids in range(j):
        val = demand - min_energy
        summation+=val
    inner_exp = (n*summation)+1
    satisfaction_val = willingness * np.log(inner_exp)
    return satisfaction_val



def sellerCost(sn, i, opt_en):
    opt_allc_en = opt_en['sellerOptEn']
    cost_sum1 = 0
    cost_sum2 = 0
    # cost factors c1 and c2
    c1 = 0
    c2 = 0.05
    for cost in range(i):
        cost_sum1 += opt_allc_en*opt_allc_en
        cost_sum2 += sn
    
    final_cost = (c1*cost_sum1) + (c2*cost_sum2)
    return final_cost


def getBidData(client):
    cluster=client["IMDC-p2p-energy"]
    collection = cluster.selectedSellers
    # get the data without the '_id:Object(..)' part
    selectedSellers = list(collection.find({},{"_id":0}))
    return selectedSellers


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


def optimalAllocation(bn, en, pn, sn, bids_j):
    
    bn amount of money the buyer wants to pay
    en amount of energy the buyer needs
    pn amount of money the seller wants (Array)
    sn amount of energy the seller can provide (array)
    note: for pn and sn, we will loop through all the sellers and to get optimal for the selected seller
    The buyer makes a single request which is sent to all sellers, so no need
    the outer loop
    
    Returns the optimal energy that buyer and seller can trade while 
    maximising social welfare

    final_arr = []
    summation = 0
    for j in range(bids_j):
        val = (bn*np.log(en)) - (pn[j]*sn[j])
        summation += val
        final_arr.append(val)

    print("Sum", summation)
    print("Energy values", final_arr)
    optAmount = max(final_arr)
    return optAmount

def optimalAllocation(bn, en, pn, sn, bids_j, en_balance):
    final_arr = []
    summation = 0
    phi = 1
    # willingness = phi/en_balance
    willingness = 100
    min_energy = 85
    n = 0.085 # average charging efficiency
    
    # optimal bid
    for j in range(bids_j):
        val = en - min_energy
        summation += val
    
    denom = (n*summation) +1
    print("denom", denom)
    opt_bid = (n*willingness*en)/denom
    print("Optimal bid, bn:", opt_bid)
    bn = opt_bid

    opt_seller_rec = []
    c1 = 0.9
    c2 = 0.9
    for j in range(bids_j):
        val = 2*c1*sn[j] - c2
        opt_seller_rec.append(val)
    
    print("Optimal energy, sn:", opt_seller_rec)




    for j in range(bids_j):
        val = (opt_bid*np.log(en)) - (pn[j]*opt_seller_rec[j])
        summation += val
        final_arr.append(val)


    print("Sum", summation)
    print("Energy values", final_arr)
    optAmount = max(final_arr)
    return optAmount
"""
