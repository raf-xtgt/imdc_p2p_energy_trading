import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
// The HttpClient service makes use of observables for all transactions. You must import the RxJS observable and operator symbols that appear in the example snippets.
import { Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';
import { User, Token } from './classes';

/** This file will allow the frontend to communicate with the backend
* using Angular's HTTP Client
*/


@Injectable({
  providedIn: 'root'
})
export class ConfigService {

  private _configUrl: string = "http://localhost:8080/";
  private _registerUrl: string = this._configUrl + "Register";
  private _loginUrl: string = this._configUrl + "Login"
  private _verifyToken: string = this._configUrl+"VerifyToken"
  TOKEN_KEY = 'token';


  //inject the HttpClient service as a dependency 
  constructor(private http: HttpClient) { }

  // add a user to the database
  addNewUser(data: User): Observable<any> {
    const body = JSON.stringify(data)
    //console.log(body)
    //console.log(this._registerUrl)
    return this.http.post<User>(this._registerUrl, body)
  }

  // authenticate a user when they want to login

  authUser(data: User): Observable<any> {
    // const headers = {
    //   headers: new HttpHeaders({ 'Content-Type': 'application/json', 'Cache-Control': 'no-cache' })
    // };
    const body = JSON.stringify(data)
    return this.http.post<User>(this._loginUrl, body)
  }

  // verify the jwt from backend
  verifyToken (){
    const body = JSON.stringify(localStorage.getItem('token'))
    return this.http.post(this._verifyToken, body)

    //return this.http.get(this._verifyToken);

  }
}