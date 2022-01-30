import { Component, OnInit } from '@angular/core';
import { ConfigService } from '../config.service';
import { BuyEnergyRequest } from '../classes';

@Component({
  selector: 'app-marketpage',
  templateUrl: './marketpage.component.html',
  styleUrls: ['./marketpage.component.css']
})
export class MarketpageComponent implements OnInit {

  constructor(private _config:ConfigService) { }

  public allBuyRequests:Array<BuyEnergyRequest>=[];

  ngOnInit(): void {
    this.getBuyRequests()
  }


  getBuyRequests(){
    this._config.getBuyRequests().subscribe(data => {
      //console.log("Buy requests data for market page", data)
      let response = JSON.parse(JSON.stringify(data))
      console.log("Buy requests data for market page", response)
      //this.allBuyRequests = response.Requests
      let reqArr = response.Requests
      for(let i = 0; i < reqArr.length; i++) {
        let request = new BuyEnergyRequest(reqArr[i].BuyerId, reqArr[i].EnergyAmount, reqArr[i].FiatAmount, reqArr[i].RequestClosed)
        this.allBuyRequests.push(request)
        // Prints i-th element of the array
        //console.log(reqArr[i]);
    }
    })
  }

}
