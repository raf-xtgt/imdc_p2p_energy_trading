import { Component, OnInit } from '@angular/core';
import { JWTService } from '../userAuth.service';
// import the custom http service module
import { ConfigService } from '../config.service';

@Component({
  selector: 'app-homepage',
  templateUrl: './homepage.component.html',
  styleUrls: ['./homepage.component.css']
})
export class HomepageComponent implements OnInit {

  private isRefreshed :boolean = false;
  constructor(private _jwtServ:JWTService, private _config:ConfigService) { }
  
   ngOnInit(): void {
      this._jwtServ.verifyToken().subscribe(data => {
        if (data != null){
          console.log("Verified Token", data)
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

}
