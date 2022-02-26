import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { BuyEnergyRequest } from './classes';
/** To send data between two components */

@Injectable({
  providedIn: 'root'
})
export class SendDataService {

  private buyEnergyRequest: BuyEnergyRequest = new BuyEnergyRequest("", 0, 0, false, "","")
  private msgSrc = new BehaviorSubject(this.buyEnergyRequest)
  // components that subscribe to this service will listen to changes in this msg
  currentMessage = this.msgSrc.asObservable();


  constructor() { }

  changeMessage(buyRequest: BuyEnergyRequest){
    // when function invoked, 
    this.msgSrc.next(buyRequest)
  }
}
