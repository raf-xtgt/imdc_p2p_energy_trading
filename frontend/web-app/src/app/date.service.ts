import { Injectable } from '@angular/core';
// The HttpClient service makes use of observables for all transactions. You must import the RxJS observable and operator symbols that appear in the example snippets.
import { Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';

/** This file will allow the frontend to communicate with the backend
* using Angular's HTTP Client
*/


@Injectable({
  providedIn: 'root'
})

export class DateService{
    constructor() { }
    getCurrentDate(): string {
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