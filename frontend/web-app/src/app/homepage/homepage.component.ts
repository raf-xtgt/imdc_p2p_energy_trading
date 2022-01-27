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

  constructor(private _jwtServ:JWTService, private _config:ConfigService) { }
  
   ngOnInit(): void {
      this._jwtServ.verifyToken().subscribe(data => {
        console.log("Verified Token", data)
      })
    
  }

}
