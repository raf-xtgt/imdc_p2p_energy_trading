import { Component, OnInit } from '@angular/core';
import { JWTService } from '../userAuth.service';
// import the custom http service module
import { ConfigService } from '../config.service';
// import class
import {HouseholdEnergyData} from '../classes';


@Component({
  selector: 'app-homepage',
  templateUrl: './homepage.component.html',
  styleUrls: ['./homepage.component.css']
})
export class HomepageComponent implements OnInit {

  constructor(private _jwtServ:JWTService, private _config:ConfigService) { }
  
  // constructor(public day:string, public average:number, public data:[number], public dateStr:string, public dateTime:number)
  model = new HouseholdEnergyData("", 0, [0,0],  "", 0)
  ngOnInit(): void {
    let dateToday = this.getHouseholdMarketData()
    this.model.dateStr = dateToday
    try{
      this._config.getHouseholdData(this.model).subscribe(data => {
        console.log("House hold data", data)
      })
  
      this._jwtServ.verifyToken().subscribe(data => {
        console.log("Verified Token", data)
      })
    }catch (err){
      console.log(err)
      window.location.reload()
    }
    
  }

  getHouseholdMarketData(): string {
    let date :Date = new Date()
    let day = date.getDate()
    let month = date.getMonth()+1
    let year = date.getFullYear()
    let dayStr = "";
    let monthStr = "";
    if (day < 10){
      dayStr += "0"+day
    }else{
      dayStr += day
    }

    if(month <10){
      monthStr += "0" + month
    }else{
      monthStr += month
    }

    let finalDate = dayStr + "-"+monthStr +"-"+ year
    return finalDate
  }
  

}
