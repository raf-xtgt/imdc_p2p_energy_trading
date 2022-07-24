import { Component, OnInit } from '@angular/core';
import { JWTService } from '../userAuth.service';
import {BuyEnergyRequest, GraphData, HouseholdEnergyData} from '../classes';
import { DateService } from '../date.service';
// import the custom http service module to communicate with backend
import { ConfigService } from '../config.service';
import { ModalService } from '../modals.service';
import Swal from 'sweetalert2'
import { Router } from '@angular/router';
import { GraphService } from '../graph.service';
import { ChartDataSets, ChartOptions } from 'chart.js';
import { Label } from 'ng2-charts';
// for the loading
import {ThemePalette} from '@angular/material/core';
import {ProgressSpinnerMode} from '@angular/material/progress-spinner';



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
  private _buyRequest = new BuyEnergyRequest("", 0,0, false, "", "")
  public currentAvgPrice :number = 0 // average price per kWh for the current day
  public energyInput :number = 0; // amount of energy required/entered by user
  public priceToPay :number = 0; // amount the user needs to pay

  public username :string =""
  private _buyerId :string = ""
  // all of this data needs to come from the backend
  public completedTran :number = 10;
  public currentFiat :number = 2000
  public currentEnergy :number = 12000;

  // buy forecast graph
  public actual_x :string[] = []
  public actual_y :number[] = []
  public pred_x :string[] = []
  public pred_y :number[] = []
  public graphData :GraphData[] = []
  // y and x axis
  public chartData: ChartDataSets[] = [{data:[], label:'Consumption Graph'}];
  public xAxis: Label[] = [];  

  // data for the card above the graph
  public currentDate:string = this.dateService.getCurrentDate()
  public currentTime:string = ""
  public currentConsumption :number = 0
  public predictionTime:string = "" // time at which the prediction is made
  public prediction:number = 0

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

  // loading before graph and all data are available
  public isLoading: boolean = true;
  color: ThemePalette = 'primary';
  mode: ProgressSpinnerMode = 'indeterminate';
  value = 100;



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
          this.getEnergyPriceData()
          this.initEnergyForecast()
        }
        
      })
    }
    else{
      this.router.navigateByUrl('/login');
    }
  }


    // ask the backend to add forecast data for this user on the database
  initEnergyForecast(){
    this._config.runBuyEnergyForecast(this._buyerId).subscribe(data => {
      console.log("Request sent to initiate forecasting", data)
      console.log("getting the data")
      this.getForecastForBuying()
    })
  }


  // get the average price in kWh from backend for the current day
  getEnergyPriceData(){
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

  // function to the get the data for the buy energy forecast
  getForecastForBuying(){
    let dateToday = this.dateService.getCurrentDate()
    //console.log("data as sent in request", dateToday)
    this._config.getBuyEnergyForecast(dateToday).subscribe(data => {
      //console.log("data to plot graph when making buy order", data)
      if (data != null){
        let response = JSON.parse(JSON.stringify(data))
        console.log("What graph data looks like", response)
        // prepare graph data for actual plot
        this.actual_x = response[0].Actual_X
        this.actual_y = response[0].Actual_Y
        let actualGraphData: GraphData = new GraphData(this.actual_y, this.actual_x, "Actual Power Consumption")
        this.graphData.push(actualGraphData)
        
        // prepare graph data for predicted plot
        this.pred_x = response[0].Pred_X
        this.pred_y = response[0].Pred_Y
        let predictedGraphData: GraphData = new GraphData(this.pred_y, this.pred_x, "Predicted Power Consumption")
        this.graphData.push(predictedGraphData)

        // draw the graph
        let makeGraph = new GraphService()
        makeGraph.data = this.graphData
        let plot = makeGraph.drawGraph()        
        this.chartData = plot.y
        this.xAxis = plot.x[1] // use the timestamps that includes the prediction
        console.log(plot)
        console.log(this.chartData)
        // get the card data
        this.currentConsumption = this.actual_y[this.actual_y.length-1] // last point is the current one
        this.prediction =  (response[0].Current_Pred).toFixed(2) 
        this.currentTime = (this.xAxis[this.xAxis.length-2]).toString()
        this.predictionTime = (this.xAxis[this.xAxis.length-1]).toString()

        // disable loading since all data has been received for now
        this.isLoading = false
      }
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
              //console.log("Buy request", this._buyRequest)
              
              if (this.orderValidation(this._buyRequest)){
                Swal.fire('Your request has been placed on the market!!', '', 'success')  
                  console.log(this._buyRequest)
                  this._config.makeBuyRequest(this._buyRequest).subscribe(data => {
                  //console.log("Response from backend for buy energy request", data)
                  this.router.navigateByUrl('/market');
                })
              }
              else{
                Swal.fire({
                  icon: 'error',
                  title: 'Oops...',
                  text: 'Please enter valid energy amount and ensure you have sufficient fiat balance!',
                })
              }
              
              

            } else if (result.isDismissed) {
              Swal.fire('Request Cancelled!', '', 'info')
            
            }
          })
  }


  // to check if the energy amount is within the predicted consumption rate
  orderValidation(req: BuyEnergyRequest){
    if (req.energyAmount <= 0 || req.energyAmount > this.prediction){
      return false
    }
    else if (this.currentFiat < req.fiatAmount){
      return false
    }
    else{
      return true
    }

  }

  /** When user clikcs on buy we need to update their account balance
   * and store the energy request in the database.
   */

  getUserDetails(data: JSON){

  }

}
