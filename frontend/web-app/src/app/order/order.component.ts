import { Component, OnInit } from '@angular/core';
import { JWTService } from '../userAuth.service';
import {BuyEnergyRequest, HouseholdEnergyData} from '../classes';
import { DateService } from '../date.service';
// import the custom http service module to communicate with backend
import { ConfigService } from '../config.service';
import { ModalService } from '../modals.service';
import Swal from 'sweetalert2'
import { Router } from '@angular/router';

@Component({
  selector: 'app-order',
  templateUrl: './order.component.html',
  styleUrls: ['./order.component.css']
})
export class OrderComponent implements OnInit {

  constructor(private _jwtServ:JWTService, private _config:ConfigService, private router: Router) { }

  private dateService = new DateService()
  public modalService = new ModalService()
  public model = new HouseholdEnergyData("", 0, [0,0],  "", 0)
  private _buyRequest = new BuyEnergyRequest("", 0,0, false)
  public currentAvgPrice :number = 0 // average price per kWh for the current day
  public energyInput :number = 0; // amount of energy required/entered by user
  public priceToPay :number = 0; // amount the user needs to pay

  public username :string =""
  private _buyerId :string = ""
  // all of this data needs to come from the backend
  public completedTran :number = 10;
  public currentFiat :number = 2000
  public currentEnergy :number = 12000;

  ngOnInit(): void {
    
      // check if the jwt is stored in local storage or not
    if (this._jwtServ.checkToken()){
      this._jwtServ.verifyToken().subscribe(data => {
        console.log("Verified Token", data)
        let response = JSON.parse(JSON.stringify(data))
        //console.log(response.Username)
        if (data !=null){
          this.username = response.User.UserName
          this._buyerId = response.User.UId
          this.getEnergyData()
          this.initEnergyForecast()
        }
        
      })
    }
    else{
      this.router.navigateByUrl('/login');
    }
  }


  initEnergyForecast(){
    this._config.runEnergyForecast().subscribe(data => {
      console.log("Request sent to initiate forecasting", data)
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
      this.currentAvgPrice = this.model.average 
      }
      else{
        console.log("lol")
      }
      console.log("House hold data for order page", response)
      //console.log("House hold data", response.Data.Data)
    })
  }

  /** process the buy request,
   * If the user confirms the request, then store it in the database
   * If the user cancels then just refresh the page **/
  createBuyRequest(){  
        Swal.fire({
            title: 'Confirm Buy Request',
            showDenyButton: false,
            showCancelButton: true,
            confirmButtonText: 'Confirm',
            //denyButtonText: denyBtnTxt,
          }).then((result) => {
            /* Read more about isConfirmed, isDenied below */
            if (result.isConfirmed) {
              this._buyRequest.buyerId = this._buyerId
              this._buyRequest.energyAmount = this.energyInput
              this._buyRequest.fiatAmount = this.energyInput * this.currentAvgPrice
              console.log("Buy request", this._buyRequest)
              Swal.fire('Your request has been placed on the market!!', '', 'success')  
              this._config.makeBuyRequest(this._buyRequest).subscribe(data => {
                //console.log("Response from backend for buy energy request", data)
              })

            } else if (result.isDismissed) {
              Swal.fire('Request Cancelled!', '', 'info')
            
            }
          })
  }

  /** When user clikcs on buy we need to update their account balance
   * and store the energy request in the database.
   */

  getUserDetails(data: JSON){

  }

}
