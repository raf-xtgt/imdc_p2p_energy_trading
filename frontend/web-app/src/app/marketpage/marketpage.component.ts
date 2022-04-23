import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { ConfigService } from '../config.service';
import { BuyEnergyRequest } from '../classes';
import {Router} from '@angular/router';
import { BidPageComponent } from '../bid-page/bid-page.component';
import { SendDataService } from '../send-data.service';
import { JWTService } from '../userAuth.service';
import Swal from 'sweetalert2'
import { TimerComponent } from '../timer/timer.component';

@Component({
  selector: 'app-marketpage',
  templateUrl: './marketpage.component.html',
  styleUrls: ['./marketpage.component.css']
})
export class MarketpageComponent implements OnInit {

  constructor(private _config:ConfigService, private router: Router, private reqData: SendDataService, private _jwtServ:JWTService) { }

  public allOpenBuyRequests:Array<BuyEnergyRequest>=[];
  public allClosedBuyRequests:Array<BuyEnergyRequest>=[];
  public noOpenBuyRequests: boolean = true;
  public noClosedBuyRequests: boolean = true;
  private requestForBid :BuyEnergyRequest = new BuyEnergyRequest("", 0, 0, false, "","") // this will hold the buy energy request data on which the prosumer makes a bid

  public buyerId: string = ""
  public message: string = "";
  private _loggedInUserId : string = "" //id of the user that is logged in

  ngOnInit(): void {
    // check if the jwt is stored in local storage or not
    if (this._jwtServ.checkToken()){
      this._jwtServ.verifyToken().subscribe(data => {
        console.log("Verified Token", data)
        let response = JSON.parse(JSON.stringify(data))
        //console.log(response.Username)
        this.getBuyRequests()
        if (data !=null){
          this._loggedInUserId = response.User.UId
          
          // subscribe to the message
          this.reqData.currentMessage.subscribe(message => this.requestForBid = message)
        }        
      })
    }
    // else{
    //   this.router.navigateByUrl('/login');
    // }
    
  }

  
  getBuyRequests(){
    this._config.getBuyRequests().subscribe(data => {
      //console.log("Buy requests data for market page", data)
      let response = JSON.parse(JSON.stringify(data))
      console.log("Buy requests data for market page", response)
      //this.allBuyRequests = response.Requests
      let reqArr = response.Requests
      for(let i = 0; i < reqArr.length; i++) {
        //console.log("All buy requests")
        this._jwtServ.gerUsername(reqArr[i].BuyerId).subscribe(data => {
          let response = JSON.parse(JSON.stringify(data))
          //console.log("response", response)
          // concatenate username and buyer id
          // let request = new BuyEnergyRequest("("+response.User.UserName+")\n"+reqArr[i].BuyerId, reqArr[i].EnergyAmount, reqArr[i].FiatAmount, reqArr[i].RequestClosed, reqArr[i].ReqId, "")
          //console.log("Closed request",reqArr[i].RequestClosed)
          if (reqArr[i].RequestClosed){

            let request = new BuyEnergyRequest("("+response.User.UserName+")\n"+reqArr[i].BuyerId, reqArr[i].EnergyAmount, reqArr[i].FiatAmount, reqArr[i].RequestClosed, reqArr[i].ReqId, "Closed")
            this.allClosedBuyRequests.push(request)
            if (this.allClosedBuyRequests.length>0){
              this.noClosedBuyRequests = false
            }
          }
          // when request is not closed
          else{
            let timer = new TimerComponent()
            let remainingTime = timer.getTimeDiff(reqArr[i])
            let timeArr = remainingTime.split(' Min')
            let timeMin = parseInt(timeArr[0])
            console.log("Minutes elapsed", timeMin)
            if (timeMin>=2){
              //close the request
              this._config.closeBuyRequest(reqArr[i].ReqId).subscribe(data => {
                console.log("buy request closed", data)
                this._config.runDoubleAuction().subscribe(data=>{
                  console.log("Finished running double auction")
                })


              })

            }
            let request = new BuyEnergyRequest("("+response.User.UserName+")\n"+reqArr[i].BuyerId, reqArr[i].EnergyAmount, reqArr[i].FiatAmount, reqArr[i].RequestClosed, reqArr[i].ReqId, remainingTime)
            this.allOpenBuyRequests.push(request)
            console.log(this.allOpenBuyRequests)
            if (this.allOpenBuyRequests.length>0){
              this.noOpenBuyRequests = false
            }
          }
          
        })
        
    }
    })
  }

  // redirect user to the bidding page
  navigateToBidPage(){
    this.router.navigateByUrl('/bid');
  }

  bid(request: BuyEnergyRequest){
    // check if bidder is a different user or not
    let buyerIdArr = request.buyerId.split('\n')
    let buyerId = buyerIdArr[1]
    if (buyerId == this._loggedInUserId){
      Swal.fire({
        icon: 'error',
        title: 'Oops...',
        text: 'You cannot bid on your own buy request!!',
      })
    }
    else{
      //send the request to the bidpage that is listening on the msg
      this.reqData.changeMessage(request)
      this.router.navigateByUrl('/bid');
    }
  }


}
