import { Component, OnInit } from '@angular/core';
import { ConfigService } from '../config.service';


@Component({
  selector: 'app-marketpage',
  templateUrl: './marketpage.component.html',
  styleUrls: ['./marketpage.component.css']
})
export class MarketpageComponent implements OnInit {

  constructor(private _config:ConfigService) { }

  ngOnInit(): void {
  }


  getBuyRequests(){
    this._config.getBuyRequests().subscribe(data => {
      console.log("Buy requests data for market page", data)
    })
  }

}
