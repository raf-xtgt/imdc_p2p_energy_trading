import { Component, OnInit } from '@angular/core';
import { JWTService } from '../userAuth.service';
@Component({
  selector: 'app-homepage',
  templateUrl: './homepage.component.html',
  styleUrls: ['./homepage.component.css']
})
export class HomepageComponent implements OnInit {

  constructor(private _jwtServ:JWTService) { }

  ngOnInit(): void {
    this._jwtServ.verifyToken().subscribe(data => {
      console.log("Verified Token", data)
    })
  }
  

}
