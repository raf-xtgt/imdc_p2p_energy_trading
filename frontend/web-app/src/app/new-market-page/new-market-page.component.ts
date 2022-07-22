import { Component, OnInit, ViewChild} from '@angular/core';
import { ConfigService } from '../config.service';
import { BuyEnergyRequest } from '../classes';
import {Router} from '@angular/router';
import { SendDataService } from '../send-data.service';
import { JWTService } from '../userAuth.service';
import Swal from 'sweetalert2'
import { TimerComponent } from '../timer/timer.component';


// imports required for the pagination
import { MatTableDataSource } from '@angular/material/table';
import { MatPaginator } from '@angular/material/paginator';
import { closedRequests, openRequests } from '../classes';

@Component({
  selector: 'app-new-market-page',
  templateUrl: './new-market-page.component.html',
  styleUrls: ['./new-market-page.component.css']
})
export class NewMarketPageComponent implements OnInit {

  constructor(private _config:ConfigService, private router: Router, private reqData: SendDataService, private _jwtServ:JWTService) { }

  
  public allOpenBuyRequests:Array<BuyEnergyRequest>=[];
  public allClosedBuyRequests:Array<BuyEnergyRequest>=[];
  public noOpenBuyRequests: boolean = true;
  public noClosedBuyRequests: boolean = true;
  private requestForBid :BuyEnergyRequest = new BuyEnergyRequest("", 0, 0, false, "","") // this will hold the buy energy request data on which the prosumer makes a bid

  public buyerId: string = ""
  public message: string = "";
  private _loggedInUserId : string = "" //id of the user that is logged in

  // for the pagination of closed buy requests
  closedRequestDisplayedCols: string[] = ['buyer', 'energyAmount', 'fiatAmount', 'reqId','remTime']
  closesRequestDataSource = new MatTableDataSource<closedRequests>(allClosedRequests)
    // add the paginator
    @ViewChild(MatPaginator) closedReqPaginator: MatPaginator | any




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
  }

  
  async getBuyRequests(){
    this._config.getBuyRequests().subscribe(data => {
      //console.log("Buy requests data for market page", data)
      let response = JSON.parse(JSON.stringify(data))
      console.log("Buy requests data for market page", response)
      //this.allBuyRequests = response.Requests
      let reqArr = response.Requests
      allClosedRequests = []
      for(let i = 0; i < reqArr.length; i++) {
        //console.log("All buy requests")
        this._jwtServ.gerUsername(reqArr[i].BuyerId).subscribe(data => {
          let response = JSON.parse(JSON.stringify(data))
          
          // for closed requests
          if (reqArr[i].RequestClosed){

            let closedRequest:closedRequests={
              buyer: "("+response.User.UserName+")\n"+reqArr[i].BuyerId,
              energyAmount: reqArr[i].EnergyAmount,
              fiatAmount: (reqArr[i].FiatAmount).toFixed(2),
              reqId: reqArr[i].ReqId,
              remTime: "Closed"
            }
            allClosedRequests.push(closedRequest)
            if (allClosedRequests.length>0){
              this.noClosedBuyRequests = false
            }
            
          }

          if (i==reqArr.length-1){

            // instantiate pagination list for closed requests
            this.closesRequestDataSource = new MatTableDataSource<closedRequests>(allClosedRequests)
            this.closesRequestDataSource.paginator = this.closedReqPaginator
          }

          
        })
        
    }

    })
  }



}

// for the closed requests
let allClosedRequests: closedRequests[] = []