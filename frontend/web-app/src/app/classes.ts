import { ChartDataSets } from "chart.js";
import { number } from "echarts";
import { Injectable } from '@angular/core';

export class User{
    constructor (public username:string, public email:string, public password:string, public address:string, public smartMeterNo:number, public uId:string, public type:string){}
}

export class Validator{
    constructor (public username:string, public email:string, public password:string, public address:string, public uId:string, public type:string, public fullName:string, public ICNum:number){}
}

// class for jwt
export class Token{
    constructor (public token:string){}
}


export class HouseholdEnergyData{
    constructor(public day:string, public average:number, public data:number[], public dateStr:string, public dateTime:number){}
}

export class BuyEnergyRequest{
    /**
     * 
     * @param buyerId 
     * @param energyAmount 
     * @param fiatAmount 
     * @param requestClosed 
     * @param reqId 
     * @param remTime String in Minutes and Seconds
     */
    constructor(public buyerId:string, public energyAmount:number, public fiatAmount:number, public requestClosed: boolean, public reqId: string, public remTime:string){}
}

// structure of bid(selling) energy request
export class SellEnergyRequest{
    /**
     * 
     * @param sellerId User id of the seller who is making the bid(sell request)
     * @param energyAmount Amount of energy that the seller can trade
     * @param fiatAmount Amount of energy the seller wants(or will receive)
     * @param sellReqId Id of the sell request
     * @param buyReqId Id of the buy request on which the bid is made
     */
    constructor(public sellerId:string, public energyAmount:number, public fiatAmount:number, public sellReqId: string, public buyReqId: string){}
}

export class ClosedBid{
    /**
     * 
     * @param buyerId Id of the buyer who made the order
     * @param buyReqId Id of the buy order request
     * @param bidIds Id of the SellEnergyRequests that were made as a bid on the buy order
     * @param bidId Id of this bid
     */
    constructor(public buyerId:string, public buyReqId:string, public bidIds:string[], public bidId: string){}
}



// the data that is sent to the graph service
export class GraphData{
    constructor(public yAxis: number[], public xAxis: string[], public label: string){}
}

// the class that plots the graph
export class Graph{
    constructor(public yAxis: ChartDataSets[], public xAxis: string[], public label: string){}
}


// request sent to backend for requesting data for plotting energy production prediction graph for sellers that make a bid
export class ProdForecastRequest {
    constructor(public userId: string, public date: string){}
}

export interface Bid {
    sellerId:string,
    optEnFromSeller:number,
    optSellerReceivable:number,
    sellerFiatBalance: number,
    sellerEnergyBalance:number
}

export interface Transaction {
    buyerId:string,
    buyerPayable:number,
    buyerEnReceivableFromAuction:number,
    buyerEnReceivableFromTNB:number,
    auctionBids:Bid[],
    TNBReceivable:number,

}

export interface Block {
    index:number,
	data:Transaction[],
	hash:string,       
	prevHash: string, 
	nonce:string,
    //info:string,
}

export interface PotentialClerks {
    username:string,
    userId:string,
    email:string,
    smartMeterNo:number
    button1: string,
}

export interface closedRequests{
    buyer:string,
    energyAmount:number, 
    fiatAmount:number,
    reqId: string,
    remTime:number|string
}

export interface openRequests{
    buyer:string,
    energyAmount:number, 
    fiatAmount:number,
    reqId: string,
    remTime:string,
    bidBtn:string
}


export interface TransactionInfo {
    buyerId:string,
    sellerId: string,
    energyTraded: number,
    fiatTraded:number,
    transactionId: string
}