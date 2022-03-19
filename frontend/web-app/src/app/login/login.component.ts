import { Component, OnInit } from '@angular/core';
import {FormGroup, FormControl, Validators} from '@angular/forms';

import {MatDialog} from '@angular/material/dialog';
import { Router } from '@angular/router';

// import the user class
import {User, Token} from '../classes';

// import the custom http service module
import { ConfigService } from '../config.service';


@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  constructor(private _config:ConfigService, public dialog: MatDialog,  private router: Router) { }

  email = new FormControl('', [Validators.required, Validators.email]);

  // this property is used in the frontend
  model = new User("", "","","",0, "", "")
  submitted = false;
  TOKEN_KEY = 'token';

  login(){
    // only assign the username and password to the model user
    let user = new User (this.model.username, "", this.model.password, "", 0, this.model.uId, "")
    this._config.authUser(user).subscribe(data => {
      console.log("Login response from backend", data)
      localStorage.setItem(this.TOKEN_KEY, data.token);
      console.log("Check local storage")
      this.router.navigateByUrl('/homepage');
    })
  }
  

  redirectToSignUp(){
    this.router.navigateByUrl('/register')
  }


  getErrorMessage() {
    if (this.email.hasError('required')) {
      return 'You must enter a value';
    }

    return this.email.hasError('email') ? 'Not a valid email' : '';
  }


}
