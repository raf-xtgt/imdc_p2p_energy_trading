import { Component, OnInit } from '@angular/core';
import { ChartDataSets, ChartOptions } from 'chart.js';
import { Color, Label } from 'ng2-charts';

// import the custom http service module
import { ConfigService } from '../config.service';
// import class
import {HouseholdEnergyData} from '../classes';
import { color } from 'echarts';

@Component({
  selector: 'app-household-price-graph',
  templateUrl: './household-price-graph.component.html',
  styleUrls: ['./household-price-graph.component.css']
})
export class HouseholdPriceGraphComponent implements OnInit {

  constructor(private _config:ConfigService) { }
  private model = new HouseholdEnergyData("", 0, [0,0],  "", 0)
  public chartData: ChartDataSets[] = [
    { data: this.model.data, label: 'Electricity Price per Household(kWh)' }
  ];
  public timeStamps: Label[] = ['12AM', '1AM', '2AM', '3AM', '4AM', '5AM', '6AM', '7AM', '8AM', '9AM', '10AM', '11AM', '12PM', '1PM', '2PM', '3PM','4PM','5PM','6PM', '7PM','8PM','9PM','10PM','11PM'];
  public lineChartOptions = {
    responsive: true,
  };

  public chartColors: any[] = [
    {
      borderColor:"#2793FF",
      backgroundColor: "#B9DCFF",
      fill:true
    }];
    ngOnInit(): void {
        
    }

    ngAfterContentInit(): void {
    let dateToday = this.getHouseholdMarketData()
    this.model.dateStr = dateToday
    try{
      this._config.getHouseholdData(this.model).subscribe(data => {
        let response = JSON.parse(JSON.stringify(data))
        if (this.model.data != undefined){
          this.model.data = response.Data.Data
        this.model.day = response.Data.Day
        this.model.dateStr = response.Data.DateStr
        this.model.average = response.Data.Average
        this.chartData = this.drawGraph()  
        }
        else{
          console.log("lol")
        }
        
        console.log("House hold data", response)
        //console.log("House hold data", response.Data.Data)
      })
    }catch (err){
      console.log(err)
      window.location.reload()
    }
    
  }

  drawGraph():ChartDataSets[]{
    let lineChartData: ChartDataSets[] = [
      { data: this.model.data, label: 'Household energy price(kWh)' },
    ];
    return lineChartData
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
