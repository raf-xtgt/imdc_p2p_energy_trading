 import { Injectable } from '@angular/core';
 import { HttpClient, HttpHeaders } from '@angular/common/http';
 // The HttpClient service makes use of observables for all transactions. You must import the RxJS observable and operator symbols that appear in the example snippets.
 import { Observable, throwError } from 'rxjs';
 import { catchError, retry } from 'rxjs/operators';
 import { User } from './classes';

 /** This file will allow the frontend to communicate with the backend
 * using Angular's HTTP Client
 */


  @Injectable({
    providedIn: 'root'
  })
 export class ConfigService {

  private _configUrl:string = "http://localhost:8080/";
  //private _registerUrl:string = this._configUrl +"Register";


  //inject the HttpClient service as a dependency 
   constructor(private http: HttpClient) { }
    
   // add a user to the database
   addNewUser(data:User): Observable<any>{
    const body = JSON.stringify(data)
    //console.log(body)
    console.log(this._configUrl)
    return this.http.post<User>(this._configUrl, body)
   }

 }