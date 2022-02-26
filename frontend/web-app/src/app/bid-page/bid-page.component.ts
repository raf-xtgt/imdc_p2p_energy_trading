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
import {ProdForecastRequest, GraphData, HouseholdEnergyData, SellEnergyRequest} from '../classes';
import { SendDataService } from '../send-data.service';
import { BuyEnergyRequest } from '../classes';
import {ThemePalette} from '@angular/material/core';
import {ProgressSpinnerMode} from '@angular/material/progress-spinner';


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
   public chartData: ChartDataSets[] = [{data:[], label:'Production graph'}];
   public xAxis: Label[] = [];  


   // data for the card above the graph
  public currentDate:string = this.dateService.getCurrentDate()
  public currentTime:string = ""
  public currentProduction :number = 0
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

  //message from market page
  public message: string = ""
  // this will hold the buy energy request data for making the bid
  private requestForBid :BuyEnergyRequest = new BuyEnergyRequest("", 0, 0, false, "", "") 

  // loading before graph and all data are available
  public isLoading: boolean = true;
  color: ThemePalette = 'primary';
  mode: ProgressSpinnerMode = 'indeterminate';
  value = 100;

  // for showing the bid info
  public buyer: string = "";
  public energyAmnt: number = 0;
  public fiatOffer: number = 0;
  public totalBids: number = 0;

  // the data that will be used to make the bid(sell request)
  public energyTradeAmount: number = 0;
  private _fiatReceived: number =0; // amount of money the seller will receive
  private _buyReqId:string = "" // the id of the buy request on which the bid is being made 


  constructor(private _jwtServ:JWTService, private _config:ConfigService, private router: Router, private reqData: SendDataService) { }

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
            // subscribe to the message from the marketpage
            this.reqData.currentMessage.subscribe(message => this.requestForBid = message)
            // bid data is here 
            console.log(this.requestForBid)
            this.showRequestInfo()
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
        //console.log("Request sent to initiate forecasting", data)
        //console.log("getting the data")
        this.getForecastForSelling()
      })
    }



      // function to the get the data for the buy energy forecast
    getForecastForSelling(){
      let dateToday = this.dateService.getCurrentDate()
      //console.log("data as sent in request", dateToday)
      let dataRequest= new ProdForecastRequest(this._sellerId, dateToday)
      this._config.getSellEnergyForecast(dataRequest).subscribe(data => {
        //console.log("data to plot graph when making a bid", data)
        if (data != null){
          let response = JSON.parse(JSON.stringify(data))
          // prepare graph data for actual plot
          this.actual_x = response[0].Actual_X
          this.actual_y = response[0].Actual_Y
          let actualGraphData: GraphData = new GraphData(this.actual_y, this.actual_x, "Actual Power Production(kWh)")
          this.graphData.push(actualGraphData)
          
          // prepare graph data for predicted plot
          this.pred_x = response[0].Pred_X
          this.pred_y = response[0].Pred_Y
          let predictedGraphData: GraphData = new GraphData(this.pred_y, this.pred_x, "Predicted Power Production(kWh)")
          this.graphData.push(predictedGraphData)

          // draw the graph
          let makeGraph = new GraphService()
          makeGraph.data = this.graphData
          let plot = makeGraph.drawGraph()        
          this.chartData = plot.y
          this.xAxis = plot.x[1] // use the timestamps that includes the prediction

          // get the card data
          this.currentProduction = this.actual_y[this.actual_y.length-1] // last point is the current one
          this.prediction =  this.pred_y[this.pred_y.length-1]
          this.currentTime = (this.xAxis[this.xAxis.length-2]).toString()
          this.predictionTime = (this.xAxis[this.xAxis.length-1]).toString()

          // disable loading since all data has been received for now
          this.isLoading = false
        }
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
        //console.log("House hold data for order page", response)
        //console.log("House hold data", response.Data.Data)
      })
    }



    // get the energy request data on which the bid is being made
    showRequestInfo(){
      // if the page is refreshed, the bid data is lost so give user direction as to what to do   
      if (this.requestForBid.buyerId == "" ){
        Swal.fire({
          icon: 'error',
          title: 'Oops...',
          text: 'Please reselect the bid from the market page',
        }).then((result)=>{
          if (result.isConfirmed){
            //console.log("olo")
            this.router.navigateByUrl('/market')
          }
        })
      }
      else{
        
        
        this.energyAmnt = this.requestForBid.energyAmount
        this.fiatOffer = this.requestForBid.fiatAmount
        this._fiatReceived = this.bidEnergyInput * this.currentAvgPrice
        //let buyerArr = this.requestForBid.buyerId.split('\n')
        //this.buyer = buyerArr[1];
        this.buyer = this.requestForBid.buyerId
        this._buyReqId = this.requestForBid.reqId
        this.totalBids = 0
      }
    }

    // make the sell request bid
    createBid(){




      Swal.fire({
        title: 'Confirm Bid',
        showDenyButton: false,
        showCancelButton: true,
        confirmButtonText: 'Confirm',
        //denyButtonText: denyBtnTxt,
      }).then((result) => {
        /* Read more about isConfirmed, isDenied below */
        if (result.isConfirmed) {
          //console.log("Buy request", this._buyRequest)
          this._fiatReceived = this.bidEnergyInput * this.currentAvgPrice
          //console.log(this._fiatReceived)
          let sellingBid = new SellEnergyRequest(this._sellerId, this.bidEnergyInput, this._fiatReceived, "", this._buyReqId)
          
          //if (this.orderValidation(this._buyRequest)){}
            Swal.fire('Your bid has been placed !!', '', 'success')  
            this._config.makeSellRequest(sellingBid).subscribe(data =>{
              //console.log("The selling bid that is stored in db", data)
            })
          
          // else{
          //   Swal.fire({
          //     icon: 'error',
          //     title: 'Oops...',
          //     text: 'Please enter valid energy amount and ensure you have sufficient fiat balance!',
          //   })
          // }
          
          

        } else if (result.isDismissed) {
          Swal.fire('Request Cancelled!', '', 'info')
        
        }
      })














    }

}
