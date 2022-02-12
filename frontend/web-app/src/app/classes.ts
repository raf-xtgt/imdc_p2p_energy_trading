import { ChartDataSets } from "chart.js";
import { number } from "echarts";
import { Injectable } from '@angular/core';

export class User{
    constructor (public username:string, public email:string, public password:string, public address:string, public smartMeterNo:number, public uId:string){}
}

// class for jwt
export class Token{
    constructor (public token:string){}
}


export class HouseholdEnergyData{
    constructor(public day:string, public average:number, public data:number[], public dateStr:string, public dateTime:number){}
}

export class BuyEnergyRequest{
    constructor(public buyerId:string, public energyAmount:number, public fiatAmount:number, public requestClosed: boolean, public reqId: string){}
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