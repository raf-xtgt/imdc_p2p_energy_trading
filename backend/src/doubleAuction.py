import numpy as np # linear algebra
from datetime import datetime
import random

def initDoubleAuction(client, biddings):
    data = [biddings]
    auctionData = prepareAuctionData(data, client)
    #print(auctionData)
    # for i in auctionData:
    #     print(i)
    # print('\n')
    auctionOutput = doubleAuction(auctionData)
    return auctionOutput


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

            buyer_fiat_balance = obj['buyerFiatBalance']
            buyer_energy_balance = obj['buyerEnergyBalance']
            
            #print("buyer account:", buyer_acc)
            bids = obj['bids']
            
            bidFiats = [] # list to hold all the bid(fiat amount) for each seller
            for bid in bids:
                fiat = bid['sellerReceivable']
                bidFiats.append(fiat)
            
            enFromSellers = []  #list to hold all the energy amount that each seller can provide
            for bid in bids:
                en = bid['sellerEnergySupply']
                enFromSellers.append(en)
            
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
                'bids': obj['bids'],
                'householdEnergyPrice': energy_price
            }
            all_auction_data.append(auction_data)
    return all_auction_data



def getAvgHouseholdPrice(client):
    now = datetime.now()
    date_str= str(now.strftime('%d-%m-%Y')) 
    #print("Date String:", date_str)
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
        # optimal energy allocation for buyer and seller
        opt_en = optimalAllocation(data)
        print("Optimal allocation", opt_en)
        print("\n")
    
    return opt_en



def optimalAllocation(data):
    # print("inside optimal allocation\n <=========================> ")
    # print(data)
    # print("inside optimal allocation\n <=========================> ")
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
    diff = "infinity"    
    willingness = 100 # always willing to charge
    min_energy = 85
    n = 0.085 # average charging efficiency
    buyer_id = data['buyerId']
    buyerPayable = data['buyerPayable']
    buyerDemand = data['buyerEnergyDemand']
    numOfBids = data['bids_j'] # number of bids made by the selected seller for all open requests
    pricing = data['householdEnergyPrice']
    bids = data['bids']
    sorted_bids = sorted(bids, key=lambda d: d['sellerReceivable'])
    auction_bids = []
    buyerEnReceivable = 0
    buyerOriginalPayable = buyerPayable # amount buyer agreed to pay when they made the buy request
    buyerOriginalDemand = buyerDemand

    while (diff=="infinity" or diff>0) and len(sorted_bids) !=0 :
         
        selected_seller = sorted_bids[0]
        seller_id = selected_seller['sellerId']
        seller_fiat_balance = selected_seller['sellerFiatBalance']
        seller_energy_balance = selected_seller['sellerEnergyBalance']
        seller_energy_supply = selected_seller['sellerEnergySupply']
        
        
        # optimal energy allocation for buyer from selected seller
        summation = 0
        for j in range(numOfBids):
            if buyerDemand < min_energy:
                summation += buyerDemand
            else:
                val = buyerDemand - min_energy
                summation += val
        
        num = ((n*summation) +1)*buyerPayable
        opt_en = num/(n*willingness) # optimal energy allocated for buyer

        # optimal energy produceable by selected seller for buyer
        c1 = 1
        c2 = 1
        opt_seller_en = (2*c1*seller_energy_supply) + c2
        increase = random.uniform(0.85, 0.90)
        if opt_seller_en >= opt_en:
            opt_seller_en = increase*opt_en

        buyerEnReceivable += opt_seller_en
        # deficit for buyer to be addressed in next round
        diff = buyerDemand - opt_en

        selected_seller_info = {
            'sellerId': seller_id,
            'optEnFromSeller': opt_seller_en,
            'optSellerReceivable':opt_seller_en*pricing,
            'sellerFiatBalance': seller_fiat_balance,
            'sellerEnergyBalance': seller_energy_balance
        }
        auction_bids.append(selected_seller_info)
        buyerDemand = diff
        buyerPayable = buyerDemand * pricing # payable for the new deficit demand
        sorted_bids.pop(0)


    total_seller_rec = 0    # total seller receivable
    for obj in auction_bids:
        total_seller_rec += obj['optSellerReceivable']
    TNBReceivable = buyerOriginalPayable - total_seller_rec
    date = getDateString()
    output = {
        'buyerId': buyer_id,
        'buyerPayable': buyerOriginalPayable,
        'buyerEnReceivableFromAuction': buyerEnReceivable,
        'buyerEnReceivableFromTNB': buyerOriginalDemand - buyerEnReceivable,
        'auctionBids': auction_bids, 
        'TNBReceivable': TNBReceivable,
        'verified': False, # whether the buyer has the required fiat amount for the transaction
        'chained':False, # whether the transaction is part of a block or not
        'tId': "", # id of the transaction in the database
        'checks': 0, # number of validators who have checked the transaction
        #'TNBReceivableFromBuyerDirect': (buyerOriginalDemand - buyerEnReceivable)*0.20,
        'date': date,
    }
    return output



def getBidData(client):
    cluster=client["IMDC-p2p-energy"]
    collection = cluster.auctionData
    # get the data without the '_id:Object(..)' part
    selectedSellers = list(collection.find({},{"_id":0}))
    return selectedSellers



def getDateString():
    now = datetime.now()
    return str(now.strftime('%d-%m-%Y')) 

