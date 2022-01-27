import { number } from "echarts";

export class User{
    constructor (public username:string, public email:string, public password:string, public address:string, public smartMeterNo:number){}
}


export class Token{
    constructor (public token:string){}
}


export class HouseholdEnergyData{
    constructor(public day:string, public average:number, public data:number[], public dateStr:string, public dateTime:number){}
}