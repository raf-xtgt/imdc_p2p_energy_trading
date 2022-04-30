import { Component, OnInit } from '@angular/core';
// import the custom http service module
import { ConfigService } from '../config.service';
import {GraphData} from '../classes';
import { ChartDataSets } from 'chart.js';
import { Label } from 'ng2-charts';
import { GraphService } from '../graph.service';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})

export class ProfileComponent implements OnInit {
  
  constructor(private _config:ConfigService) { }

  TOKEN_KEY = 'token';
  public username :string = "";
  public email :string = "";
  public address :string = "";
  public smartMetreNo :number = 0;
  public _userId :string = "";


  // income graph
  public hash :string[] = []
  public dates :string[] = []
  public fiatReceived :number[] = []
  public totalIncome :number = 0
  public totalEnSold :number = 0
  public energySold :number[] = []

  public graphData :GraphData[] = []
  // y and x axis
  public chartData: ChartDataSets[] = [{data:[], label:'Income chart'}];
  public xAxis: Label[] = []; 
  
  public barChartOptions = {
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

  async ngOnInit(): Promise<void> {
    await this.verifyUserJWT()
    
  }

  verifyUserJWT(){
    //let token : string = localStorage.getItem("token")
    this._config.verifyToken().subscribe(data => {
      console.log("Verified Token", data)
      let response = JSON.parse(JSON.stringify(data))
      console.log(response)
      this.username = response.User.UserName
      this.email = response.User.Email
      this.address = response.User.Address
      this.smartMetreNo = response.User.SmartMeterNo
      this._userId = response.User.UId
      console.log(this._userId)
      if (this.username == 'tnb'){
        console.log("show income for tnb")
        this.getTNBIncomeData()
      }
      else{
        this.getUserIncomeData()
      }
      
    })
  }

  // to get the income for normal users
  getUserIncomeData(){
    this._config.getUserIncome(this._userId).subscribe(data => {
      
      let response = JSON.parse(JSON.stringify(data))
      console.log("Response from backend", response)
      //let incomeInfo = response.Receivable


      let allDates: string[] = response.Dates
      let fiatRec: string[] = response.Dates
      let seen :string[] = []
      
      let newFiatReceivables :number[]= []
      let newEnTrade :number[] = []
      for (let i=0; i<response.Dates.length; i++){
        let current = response.Dates[i]
        let newFiatVal = response.Receivable[i]
        let newEnVal = response.EnergySold[i]
        if ( (this.seenDate(current, seen)) == false ){
          seen.push(current)
          
          for (let j = i+1; j<response.Dates.length; j++){
            let next = response.Dates[j]
            if (current == next){
              newFiatVal +=  response.Receivable[j]
              newEnVal += response.EnergySold[j]
            }

            if (j==response.Dates.length-1){
              newFiatReceivables.push(newFiatVal)
              newEnTrade.push(newEnVal)
            } 

          }
          
        }
      }
      console.log(newFiatReceivables)
      console.log(newEnTrade)
      console.log(seen)

      // income graph data
      let incomeGraphData: GraphData = new GraphData(newFiatReceivables, seen, "Income")
      this.graphData.push(incomeGraphData)
      
      // energy graph data
      let energyGraphData: GraphData = new GraphData(newEnTrade, seen, "Energy Amount Traded")
      this.graphData.push(energyGraphData)

      // draw the graph
      let makeGraph = new GraphService()
      makeGraph.data = this.graphData
      let plot = makeGraph.drawGraph()        
      this.chartData = plot.y
      this.xAxis = plot.x[0] // use the timestamps that includes the prediction
      //console.log(plot)

      // calculate totals
      this.totalIncome = this.summation(newFiatReceivables)
      this.totalEnSold = this.summation(newEnTrade)


    })
  }


  // to get the income for tnb
  getTNBIncomeData(){
    this._config.getTNBIncome().subscribe(data => {
      
      let response = JSON.parse(JSON.stringify(data))
      console.log("Response from backend", response)
      //let incomeInfo = response.Receivable

      // income graph data
      let incomeGraphData: GraphData = new GraphData(response.Receivable, response.Dates, "Income")
      this.graphData.push(incomeGraphData)
      
      // energy graph data
      let energyGraphData: GraphData = new GraphData(response.EnergySold, response.Dates, "Energy Sold")
      this.graphData.push(energyGraphData)

      // draw the graph
      let makeGraph = new GraphService()
      makeGraph.data = this.graphData
      let plot = makeGraph.drawGraph()        
      this.chartData = plot.y
      this.xAxis = plot.x[0] // use the timestamps that includes the prediction
      //console.log(plot)

      // calculate totals
      this.totalIncome = this.summation(response.Receivable)
      this.totalEnSold = this.summation(response.EnergySold)


    })
  }

  summation(arr: number[]): number {
    let summation :number = 0
    for (let i=0 ; i<arr.length; i++){
      summation += arr[i]
    }
    return summation

  }

  // to check whether a date exists in a list or not
  seenDate(date:string, allDates:string[]): boolean{
    for (let k=0; k<allDates.length; k++){
      if (allDates[k] == date){
        return true
        break
        }
    }
    return false
  }

}
