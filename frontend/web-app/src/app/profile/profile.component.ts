import { Component, OnInit } from '@angular/core';
// import the custom http service module
import { ConfigService } from '../config.service';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})

export class ProfileComponent implements OnInit {
  TOKEN_KEY = 'token';
  jwt = ''
  constructor(private _config:ConfigService) { }

  ngOnInit(): void {
    this.verifyUserJWT()
  }

  verifyUserJWT(){
    //let token : string = localStorage.getItem("token")
    this._config.verifyToken().subscribe(data => {
      console.log("Verified Token", data)
    })
  }

}
