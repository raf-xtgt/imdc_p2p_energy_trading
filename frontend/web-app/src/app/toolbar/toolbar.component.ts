import { Component, OnInit } from '@angular/core';
import {MatToolbarModule} from '@angular/material/toolbar'
import { JWTService } from '../userAuth.service';

@Component({
  selector: 'app-toolbar',
  templateUrl: './toolbar.component.html',
  styleUrls: ['./toolbar.component.css']
})
export class ToolbarComponent implements OnInit {

  constructor(private _jwtServ:JWTService) { }
  isVerified :boolean = false;
  username :string = ""

  ngOnInit(): void {
    this._jwtServ.verifyToken().subscribe(data => {
      console.log("Verified Token", data)
      let response = JSON.parse(JSON.stringify(data))
      //console.log(response.Username)
      if (data !=null){
        this.username = response.Username
        this.isLoggedIn()
      }
      
    })
    
    
    
  }

  isLoggedIn(){
    return this.isVerified = true
  }

  events: string[] = [];
  opened: boolean = false

  shouldRun = [/(^|\.)plnkr\.co$/, /(^|\.)stackblitz\.io$/].some(h => h.test(window.location.host));

}
