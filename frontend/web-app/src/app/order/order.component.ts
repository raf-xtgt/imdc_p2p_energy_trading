import { Component, OnInit } from '@angular/core';
import { JWTService } from '../userAuth.service';
import {HouseholdEnergyData} from '../classes';
import { DateService } from '../date.service';
// import the custom http service module to communicate with backend
import { ConfigService } from '../config.service';
@Component({
  selector: 'app-order',
  templateUrl: './order.component.html',
  styleUrls: ['./order.component.css']
})
export class OrderComponent implements OnInit {

  constructor(private _jwtServ:JWTService, private _config:ConfigService) { }

  private dateService = new DateService()
  public model = new HouseholdEnergyData("", 0, [0,0],  "", 0)


  public username :string =""
  // all of this data needs to come from the backend
  public completedTran :number = 10;
  public currentFiat :number = 2000
  public currentEnergy :number = 12000;

  ngOnInit(): void {
    this._jwtServ.verifyToken().subscribe(data => {
      console.log("Verified Token", data)
      let response = JSON.parse(JSON.stringify(data))
      //console.log(response.Username)
      if (data !=null){
        this.username = response.Username
        this.getEnergyData()
      }
      
    })
  }

  // get the average price in kWh from backend for the current day
  getEnergyData(){
    let dateToday = this.dateService.getCurrentDate()
    this.model.dateStr = dateToday
    this._config.getHouseholdData(this.model).subscribe(data => {
      let response = JSON.parse(JSON.stringify(data))
      if (this.model.data != undefined){
        this.model.data = response.Data.Data
      this.model.day = response.Data.Day
      this.model.dateStr = response.Data.DateStr
      this.model.average = response.Data.Average  
      }
      else{
        console.log("lol")
      }
      
      console.log("House hold data for order page", response)
      //console.log("House hold data", response.Data.Data)
    })
  }


  getUserDetails(data: JSON){

  }

}
