import { Injectable } from '@angular/core';
import { ChartDataSets, ChartOptions } from 'chart.js';
import { Color, Label } from 'ng2-charts';
import { color } from 'echarts';
import { GraphData } from './classes';

@Injectable({
    providedIn: 'root'
})

export class GraphService{
    public chartData: ChartDataSets[] = [
        { data: [], label: '' }
      ];
    public xAxis: Label[] = [];
    public lineChartOptions = {
        responsive: true,
    };
    
    public chartColors: any[] = [
    {
        borderColor:"#2793FF",
        backgroundColor: "#B9DCFF",
        fill:true
    }];
    public data: GraphData[] = []
    constructor() { }

    drawGraph(){
        let lineData: ChartDataSets[] = new Array(this.data.length)
        let xAxes = []
        for (let i=0; i<this.data.length; i++){
            lineData[i] = { data: this.data[i].yAxis, label: this.data[i].label }
            xAxes.push(this.data[i].xAxis)
        }
        let output = {
            "y": lineData,
            'x':xAxes
        }
        return output
    }

    // getYAxis(){

    //     let lineChartData: ChartDataSets[] = [
    //       { data: this.data.yAxis, label: 'Household energy price(kWh)' },
    //     ];
    //     this.xAxis = this.data.xAxis
        
    //   }

    // getXAxis(){
    //     this.xAxis = this.data.xAxis
    //     return this.xAxis
    // }
}

