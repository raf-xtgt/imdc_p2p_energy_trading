import { dashCaseToCamelCase } from '@angular/compiler/src/util';
import { Component, OnInit } from '@angular/core';
import {FormControl, Validators} from '@angular/forms';
 
// import the user class
import {User} from '../classes';

// import the custom http service module
import { ConfigService } from '../config.service';

@Component({
  selector: 'app-signup',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.css']
})
export class SignupComponent implements OnInit {

  constructor(private _config:ConfigService) { }

  ngOnInit(): void {
  }
  
  email = new FormControl('', [Validators.required, Validators.email]);


  //select = ['He is fat', 'Mikasa']
  model = new User("", "")
  submitted = false;

  addUser(){
    this.submitted = true;
    //console.log(this.model.fullname)
    //console.log(this.model.email)
    let newUser = new User(this.model.fullname, this.model.email)
    this._config.addUser(newUser)
  }

}
