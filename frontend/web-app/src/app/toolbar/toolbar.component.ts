import { Component, OnInit } from '@angular/core';
import {MatToolbarModule} from '@angular/material/toolbar'
import { JWTService } from '../userAuth.service';
import { Router } from '@angular/router';


@Component({
  selector: 'app-toolbar',
  templateUrl: './toolbar.component.html',
  styleUrls: ['./toolbar.component.css']
})
export class ToolbarComponent implements OnInit {

  constructor(private _jwtServ:JWTService, private router: Router) { }
  isVerified :boolean = false;
  username :string = ""

  ngOnInit(): void {
    // check if the jwt is stored in local storage or not
      this._jwtServ.verifyToken().subscribe(data => {
        console.log("Verified Token", data)
        let response = JSON.parse(JSON.stringify(data))
        console.log(response.User)
        if (data !=null){
          this.username = response.User.UserName
          this.isLoggedIn()
        }
        
      })
    
    
    
  }


  logout(){
    window.localStorage.removeItem("token")
  }

  isLoggedIn(){
    return this.isVerified = true
  }

  events: string[] = [];
  opened: boolean = false

  shouldRun = [/(^|\.)plnkr\.co$/, /(^|\.)stackblitz\.io$/].some(h => h.test(window.location.host));

}
