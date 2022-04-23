import { Component, OnInit } from '@angular/core';
// import the custom http service module
import { ConfigService } from '../config.service';

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
      console.log("Response from backend", data)
    })
  }

}
