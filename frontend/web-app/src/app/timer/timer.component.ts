import { Component} from '@angular/core';

@Component({
  selector: 'app-timer',
  templateUrl: './timer.component.html',
  styleUrls: ['./timer.component.css']
})
export class TimerComponent {

  public dateNow = new Date();
  milliSecondsInASecond = 1000;
  hoursInADay = 24;
  minutesInAnHour = 60;
  SecondsInAMinute  = 60;

  public timeDifference: number=0;
  public secondsToDday: number=0;
  public minutesToDday: number=0;


  constructor() { }

private allocateTimeUnits (timeDifference: number):string {
    this.secondsToDday = Math.floor((timeDifference) / (this.milliSecondsInASecond) % this.SecondsInAMinute);
    this.minutesToDday = Math.floor((timeDifference) / (this.milliSecondsInASecond * this.minutesInAnHour) % this.SecondsInAMinute);
    //this.hoursToDday = Math.floor((timeDifference) / (this.milliSecondsInASecond * this.minutesInAnHour * this.SecondsInAMinute) % this.hoursInADay);
    //this.daysToDday = Math.floor((timeDifference) / (this.milliSecondsInASecond * this.minutesInAnHour * this.SecondsInAMinute * this.hoursInADay));
    console.log("Time remaining ", this.minutesToDday," Min ", this.secondsToDday, " S")
    let output:string = this.minutesToDday+ " Min "+ this.secondsToDday+ " S"
    return  output
  }

public getTimeDiff(request: any): string{
  let remTimeStr:string = ""
  if (!request.RequestClosed) {
    let remTime: number = 0
    let reqTime: string = request.ReqTime
    let reqDate: Date = new Date(reqTime)
    let reqMs: number = reqDate.getTime()
    remTime = this.dateNow.getTime() - reqMs
    remTimeStr = this.allocateTimeUnits(remTime)
  }
  return remTimeStr
}


}
