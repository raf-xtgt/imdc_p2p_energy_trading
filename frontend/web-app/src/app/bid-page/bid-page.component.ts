import { Component, OnInit } from '@angular/core';
import { JWTService } from '../userAuth.service';
import { ConfigService } from '../config.service';
import { ModalService } from '../modals.service';
import Swal from 'sweetalert2'
import { Router } from '@angular/router';
import { DateService } from '../date.service';
import { GraphService } from '../graph.service';
import { ChartDataSets, ChartOptions } from 'chart.js';
import { Color, Label } from 'ng2-charts';
import {BuyEnergyRequest, GraphData, HouseholdEnergyData} from '../classes';


@Component({
  selector: 'app-bid-page',
  templateUrl: './bid-page.component.html',
  styleUrls: ['./bid-page.component.css']
})
export class BidPageComponent implements OnInit {

  
  private dateService = new DateService()
  public modalService = new ModalService()
  public pricingData = new HouseholdEnergyData("", 0, [0,0],  "", 0)
  
  // form elements
  public bidEnergyInput: number = 0;
  public currentAvgPrice :number = 0 // average price per kWh for the current day

  //user info and account info
  public completedTran :number = 10;
  public currentFiat :number = 2000
  public currentEnergy :number = 12000;
  public username :string =""
  private _sellerId :string = ""



   // sell forecast graph
   public actual_x :string[] = []
   public actual_y :number[] = []
   public pred_x :string[] = []
   public pred_y :number[] = []
   public graphData :GraphData[] = []
   // y and x axis
   public chartData: ChartDataSets[] = [];
   public xAxis: Label[] = [];  


   // data for the card above the graph
  public currentDate:string = this.dateService.getCurrentDate()
  public currentTime:string = ""
  public currentConsumption :number = 0
  public predictionTime:string = "" // time at which the prediction is made
  public prediction:number = 0

  // graph visuals
  public lineChartOptions = {
    responsive: true,
  };
  public chartColors: any[] = [
    {
        borderColor:"#2793FF",
        backgroundColor: "#B9DCFF",
        fill:false
    },

    // predicted
    {
      borderColor: "#FF6347",
      backgroundColor: "#B22222",
      fill:false
  },
  
  ];

  constructor(private _jwtServ:JWTService, private _config:ConfigService, private router: Router) { }

  ngOnInit(): void {
       // check if the jwt is stored in local storage or not
       if (this._jwtServ.checkToken()){
        this._jwtServ.verifyToken().subscribe(data => {
          console.log("Verified Token", data)
          let response = JSON.parse(JSON.stringify(data))
          //console.log(response.Username)
          if (data !=null){
            this.username = response.User.UserName
            this._sellerId = response.User.UId
            this.getEnergyPriceData()
            this.initSellEnergyForecast()
          }
          
        })
      }
      else{
        this.router.navigateByUrl('/login');
      }
  }

  
     // ask the backend to add forecast data for this user on the database
     initSellEnergyForecast(){
      this._config.runSellEnergyForecast(this._sellerId).subscribe(data => {
        console.log("Request sent to initiate forecasting", data)
        console.log("getting the data")
        //this.getForecastForBuying()
      })
    }

    // get the average price in kWh from backend for the current day
    getEnergyPriceData(){
      let dateToday = this.dateService.getCurrentDate()
      this.pricingData.dateStr = dateToday
      this._config.getHouseholdData(this.pricingData).subscribe(data => {
        let response = JSON.parse(JSON.stringify(data))
        if (this.pricingData.data != undefined){
          this.pricingData.data = response.Data.Data
        this.pricingData.day = response.Data.Day
        this.pricingData.dateStr = response.Data.DateStr
        this.pricingData.average = response.Data.Average  
        this.currentAvgPrice = this.pricingData.average 
        }
        else{
          console.log("lol")
        }
        console.log("House hold data for order page", response)
        //console.log("House hold data", response.Data.Data)
      })
    }


    createBid(){}

}
