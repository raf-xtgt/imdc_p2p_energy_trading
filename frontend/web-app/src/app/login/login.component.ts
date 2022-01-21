import { Component, OnInit } from '@angular/core';
import {FormGroup, FormControl, Validators} from '@angular/forms';

import {MatDialog} from '@angular/material/dialog';
 

// import the user class
import {User} from '../classes';

// import the custom http service module
import { ConfigService } from '../config.service';


@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  constructor(private _config:ConfigService, public dialog: MatDialog) { }

  

  email = new FormControl('', [Validators.required, Validators.email]);

  // this property is used in the frontend
  model = new User("", "","","",0)
  submitted = false;


  login(){
    // only assign the username and password to the model user
    let user = new User (this.model.username, "", this.model.password, "", 0)
    this._config.authUser(user).subscribe(data => {
      console.log("Login response from backend", data)
      if (data!= null){
        console.log("User successfully logged in")
      }else{
        console.log("Credentials do not match")
      }
    })
  }



  getErrorMessage() {
    if (this.email.hasError('required')) {
      return 'You must enter a value';
    }

    return this.email.hasError('email') ? 'Not a valid email' : '';
  }


}
