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
      this.getUserIncomeData()
    })
  }

  getUserIncomeData(){
    this._config.getUserIncome(this._userId).subscribe(data => {
      
      let response = JSON.parse(JSON.stringify(data))
      console.log("Response from backend", response)
      //let incomeInfo = response.Receivable
   
      //this.hash = response.BlockHashes
      for (let x=0; x<response.BlockHashes.length; x++){
        let str = "Block"+(x+1)
        this.hash.push(str)
      }
      this.fiatReceived = response.Receivable
      this.energySold = response.EnergySold
      console.log(this.energySold)

      // income graph data
      let incomeGraphData: GraphData = new GraphData(this.fiatReceived, this.hash, "Income")
      this.graphData.push(incomeGraphData)
      
      // energy graph data
      let energyGraphData: GraphData = new GraphData(this.energySold, this.hash, "Energy Amount Traded")
      this.graphData.push(energyGraphData)

      // draw the graph
      let makeGraph = new GraphService()
      makeGraph.data = this.graphData
      let plot = makeGraph.drawGraph()        
      this.chartData = plot.y
      this.xAxis = plot.x[0] // use the timestamps that includes the prediction
      console.log(plot)

      // calculate totals
      this.totalIncome = this.summation(this.fiatReceived)
      this.totalEnSold = this.summation(this.energySold)


    })
  }

  summation(arr: number[]): number {
    let summation :number = 0
    for (let i=0 ; i<arr.length; i++){
      summation += arr[i]
    }
    return summation

  }

}
