import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
// The HttpClient service makes use of observables for all transactions. You must import the RxJS observable and operator symbols that appear in the example snippets.
import { Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';
import { User, Token, HouseholdEnergyData, BuyEnergyRequest, ProdForecastRequest, SellEnergyRequest, Validator } from './classes';

/** This file will allow the frontend to communicate with the backend
* using Angular's HTTP Client
*/


@Injectable({
  providedIn: 'root'
})
export class ConfigService {

  private _configUrl: string = "http://localhost:8080/";
  private _registerUrl: string = this._configUrl + "Register";
  private _loginUrl: string = this._configUrl + "Login"
  private _verifyToken: string = this._configUrl+"VerifyToken"
  private _addHouseholdPrice: string = this._configUrl+"AddHouseholdData"
  private _buyRequestUrl: string = this._configUrl + 'CreateBuyRequest'
  private _sellRequestUrl: string = this._configUrl + 'CreateSellRequest'
  private _getBuyRequestUrl: string = this._configUrl + 'GetBuyRequests'
  private _energyForecastURL: string = this._configUrl+ 'RunBuyEnergyForecast'
  private _latestBuyForecastURL: string = this._configUrl+ 'GetLatestBuyForecast'
  private _sellEnergyForecastURL: string = this._configUrl+ 'RunSellEnergyForecast'
  private _latestSellForecastURL: string = this._configUrl+ 'GetLatestSellForecast'
  private _closeBuyRequestURL:string = this._configUrl+'CloseBuyRequest'
  private _runDoubleAuction: string = this._configUrl + 'RunDoubleAuction'
  private _addValidator:string = this._configUrl + 'AddValidator'

  //blockchain urls
  private _createGenesis: string = this._configUrl+ 'CreateGenesisBlock'
  private _updateBlockchain: string = this._configUrl +'UpdateBlockchain'
  private _getBlockchain: string = this._configUrl + 'GetBlockchain'

  // url to get all user data
  private _getAllUsers: string = this._configUrl + 'GetAllUsers'

  // url to make a clerk
  private _makeClerk:string = this._configUrl + 'MakeClerk'
  private _clerkINTCheck :string = this._configUrl + 'ClerkIntegrityCheck'
 
  // to get profits for normal users
  private _getUserIncome :string = this._configUrl + 'GetUserIncome'

  // to get profits for tnb
  private _getTNBProfit :string = this._configUrl + 'GetTNBIncome'

  //get buy requests for the user
  private _getUserBuyRequests :string = this._configUrl + 'GetUserBuyRequests'

  //get sell requests for the user
  private _getUserSellRequests :string = this._configUrl + 'GetUserSellRequests'

  // get sell requests for the user

  TOKEN_KEY = 'token';


  //inject the HttpClient service as a dependency 
  constructor(private http: HttpClient) { }

  // add a user to the database
  addNewUser(data: User): Observable<any> {
    const body = JSON.stringify(data)
    return this.http.post<User>(this._registerUrl, body)
  }

  addNewValidator(data:Validator) :Observable<any>{
    const body = JSON.stringify(data)
    return this.http.post<Validator>(this._addValidator, body)  
  }

  // authenticate a user when they want to login  
  authUser(data: User): Observable<any> {
    const body = JSON.stringify(data)
    return this.http.post<User>(this._loginUrl, body)
  }

  // verify the jwt from backend
  verifyToken (){
    const body = JSON.stringify(localStorage.getItem('token'))
    return this.http.post(this._verifyToken, body)
  }

  // get energy price data for today
  getHouseholdData (data:HouseholdEnergyData){
    const body = JSON.stringify(data)
    return this.http.post(this._addHouseholdPrice, body)
  }

  // store buy request data in the database
  makeBuyRequest (data: BuyEnergyRequest){
    const body = JSON.stringify(data)
    //console.log("Buy request data to send to backend", body)
    return this.http.post(this._buyRequestUrl, body)
  }

  makeSellRequest (data: SellEnergyRequest){
    const body = JSON.stringify(data)
    console.log("Sell request data that is sent to backend", body)
    return this.http.post(this._sellRequestUrl, body)
  }


  // to get all the open buy requests made by all users
  getBuyRequests(){
    const body = JSON.stringify("Get energy data")
    //console.log("Buy request data to send to backend", body)
    return this.http.post(this._getBuyRequestUrl, body)
  }

  // to run the python script that will do energy forecasting via golang
  // this will produce prediction for the amount of energy that needs to be consumed
  runBuyEnergyForecast(userId: string){
    const body = JSON.stringify(userId)
    return this.http.post(this._energyForecastURL, body)
  }

  // to get the last energy forecast for users making a buy order
  getBuyEnergyForecast(date: string){
    const body = JSON.stringify(date)
    return this.http.post(this._latestBuyForecastURL, body)
  }

  //to run py script that will produce prediction for the amount of energy that can be produced
  runSellEnergyForecast(userId: string){
    const body = JSON.stringify(userId)
    return this.http.post(this._sellEnergyForecastURL, body)
  }

  // to get the latest energy production forecast for the current users making a buy order
  getSellEnergyForecast(data: ProdForecastRequest ){
    const body = JSON.stringify(data)
    return this.http.post(this._latestSellForecastURL, body)
  }

  closeBuyRequest(data: string){
    const body = JSON.stringify(data)
    return this.http.post(this._closeBuyRequestURL, body)
  }

  runDoubleAuction(){
    const body = JSON.stringify("RUN DOUBLE AUCTION")
    return this.http.post(this._runDoubleAuction, body)
  }

  createGenesisBlock(){
    const body = JSON.stringify("Genesis")
    return this.http.post(this._createGenesis, body)
  }

  updateBlockchain(){
    const body = JSON.stringify("Update Blockchain")
    return this.http.post(this._updateBlockchain, body)

  }

  getCurrentBlockchain(){
    const body = JSON.stringify("Get Blockchain")
    return this.http.post(this._getBlockchain, body)
  }

  getAllUsers(){
    const body = JSON.stringify("Get All Users")
    return this.http.post(this._getAllUsers, body)
  }

  // convert a normal user to a clerk
  convertToClerk(userId: string){
    const body = JSON.stringify(userId)
    return this.http.post(this._makeClerk, body)
  }

  // to invoke the clerk integrity check
  initClerkINTChk(){
    const body = JSON.stringify("Clerk integrity check")
    return this.http.post(this._clerkINTCheck, body)
  }

  // to get user profits
  getUserIncome(userId:string){
    const body = JSON.stringify(userId)
    console.log(body)
    return this.http.post(this._getUserIncome, body)

  }

  getTNBIncome(){
    const body = JSON.stringify("tnbIncome")
    return this.http.post(this._getTNBProfit, body)
  }


  // get the buy requests made by the user
  getUserBuyRequests(userId: string){
    const body = JSON.stringify(userId)
    return this.http.post(this._getUserBuyRequests, body)
  }

  
  // get the sell/bids made by the user
  getUserSellRequests(userId: string){
    const body = JSON.stringify(userId)
    return this.http.post(this._getUserSellRequests, body)
  }

}