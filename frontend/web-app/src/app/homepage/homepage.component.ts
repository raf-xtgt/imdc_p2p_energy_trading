import { Component, OnInit } from '@angular/core';
import { JWTService } from '../userAuth.service';
// import the custom http service module
import { ConfigService } from '../config.service';
// for routing to a page
import {Router} from '@angular/router';


@Component({
  selector: 'app-homepage',
  templateUrl: './homepage.component.html',
  styleUrls: ['./homepage.component.css']
})
export class HomepageComponent implements OnInit {

  private isRefreshed :boolean = false;
  constructor(private _jwtServ:JWTService, private _config:ConfigService, private router: Router) { }
  private isAuth:boolean = false // to check whether the user is logged in or not
  
   ngOnInit(): void {
      this._jwtServ.verifyToken().subscribe(data => {
        if (data != null){
          console.log("Verified Token", data)
          this.isAuth = true
          this.refreshPage()
        }
      })
  }

  // refresh the page when there is no refreshed key in locals storage
  refreshPage():void {
    if (!localStorage.getItem('refreshed')) { 
      localStorage.setItem('refreshed', 'no reload') 
      location.reload() 
    } else {
      localStorage.removeItem('refreshed') 
    }

  }

  // redirect user to market page
  redirectToMarket():void{
    if (this.isAuth){
      this.router.navigateByUrl('/market');
    }
    else{
      this.router.navigateByUrl('/login')
    }
    
  }

}
