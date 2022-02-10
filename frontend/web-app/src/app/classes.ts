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
    constructor(public buyerId:string, public energyAmount:number, public fiatAmount:number, public requestClosed: boolean){}
}

// the data that is sent to the graph service
export class GraphData{
    constructor(public yAxis: number[], public xAxis: string[], public label: string){}
}

// the class that plots the graph

export class Graph{
    constructor(public yAxis: ChartDataSets[], public xAxis: string[], public label: string){}
}