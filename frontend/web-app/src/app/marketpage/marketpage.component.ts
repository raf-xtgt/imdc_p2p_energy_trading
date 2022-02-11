import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { ConfigService } from '../config.service';
import { BuyEnergyRequest } from '../classes';
import {Router} from '@angular/router';
import { BidPageComponent } from '../bid-page/bid-page.component';
import { SendDataService } from '../send-data.service';

@Component({
  selector: 'app-marketpage',
  templateUrl: './marketpage.component.html',
  styleUrls: ['./marketpage.component.css']
})
export class MarketpageComponent implements OnInit {

  constructor(private _config:ConfigService, private router: Router, private reqData: SendDataService) { }

  public allBuyRequests:Array<BuyEnergyRequest>=[];
  private requestForBid :BuyEnergyRequest = new BuyEnergyRequest("", 0, 0, false) // this will hold the buy energy request data on which the prosumer makes a bid

  public buyerId: string = ""
  public message: string = "";

  ngOnInit(): void {
    this.getBuyRequests()

    // subscribe to the message
    this.reqData.currentMessage.subscribe(message => this.requestForBid = message)
  }


  getBuyRequests(){
    this._config.getBuyRequests().subscribe(data => {
      //console.log("Buy requests data for market page", data)
      let response = JSON.parse(JSON.stringify(data))
      //console.log("Buy requests data for market page", response)
      //this.allBuyRequests = response.Requests
      let reqArr = response.Requests
      for(let i = 0; i < reqArr.length; i++) {
        let request = new BuyEnergyRequest(reqArr[i].BuyerId, reqArr[i].EnergyAmount, reqArr[i].FiatAmount, reqArr[i].RequestClosed)
        this.allBuyRequests.push(request)
    }
    })
  }

  // redirect user to the bidding page
  navigateToBidPage(){
    this.router.navigateByUrl('/bid');
  }

  bid(request: BuyEnergyRequest){
    console.log("the buy energy request", request.buyerId)
    //send the request to the bidpage that is listening on the msg
    this.reqData.changeMessage(request)
    this.router.navigateByUrl('/bid');
  }


}
