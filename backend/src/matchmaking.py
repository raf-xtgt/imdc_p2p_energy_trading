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

