import { Component, OnInit } from '@angular/core';
import { JWTService } from '../userAuth.service';


@Component({
  selector: 'app-order',
  templateUrl: './order.component.html',
  styleUrls: ['./order.component.css']
})
export class OrderComponent implements OnInit {

  constructor(private _jwtServ:JWTService) { }


  username :string =""
  // all of this data needs to come from the backend
  completedTran :number = 10;
  currentFiat :number = 2000
  currentEnergy :number = 12000;
  //userDetails :string = this.username+"\n"+this.currentFiat+"\n"+ this.currentEnergy

  ngOnInit(): void {
    this._jwtServ.verifyToken().subscribe(data => {
      console.log("Verified Token", data)
      let response = JSON.parse(JSON.stringify(data))
      //console.log(response.Username)
      if (data !=null){
        this.username = response.Username
      }
      
    })
  }

  getUserDetails(data: JSON){

  }

}
