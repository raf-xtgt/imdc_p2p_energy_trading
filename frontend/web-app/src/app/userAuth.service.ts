import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

/** This file will allow the frontend to communicate with the backend
* using Angular's HTTP Client to verify the token stored in the local storage
*/


@Injectable({
  providedIn: 'root'
})
export class JWTService {

  private _configUrl: string = "http://localhost:8080/";
  private _verifyToken: string = this._configUrl+"VerifyToken"
  private _getUserDetailsUrl: string = this._configUrl + 'GetUserDetails'
  TOKEN_KEY = 'token';


  //inject the HttpClient service as a dependency 
  constructor(private http: HttpClient) { }

  // verify the jwt from backend
  verifyToken (){
    // get the token from local storage
    const body = JSON.stringify(localStorage.getItem('token'))
    return this.http.post(this._verifyToken, body)
  }

  // check if the JWT is stored in local storage or not
  checkToken(){
    if (localStorage.getItem("token")=== null){
      return false
    }
    else{
      return true
    }
  }

  // given userId, get the username of the user
  gerUsername(userId: string){
    const body = JSON.stringify(userId)
    return this.http.post(this._getUserDetailsUrl, body)
  }
}