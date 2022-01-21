import { dashCaseToCamelCase } from '@angular/compiler/src/util';
import { Component, OnInit } from '@angular/core';
import {FormControl, Validators} from '@angular/forms';
import {MatDialog} from '@angular/material/dialog';
 
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

  constructor(private _config:ConfigService, public dialog: MatDialog) { }

  ngOnInit(): void {
  }
  

  //select = ['He is fat', 'Mikasa']
  model = new User("", "","","",0)
  submitted = false;

  /** Method to check whether the username and email is unique or not.
   * If an account with the same username and email, exist, then
   * ask for using different credentials
   */
  checkUsernameAndEmail () {
    
  }

  /** Method to add user to the database */
  addUser(){
    this.submitted = true;
    //console.log(this.model.fullname)
    //console.log(this.model.email)
    let newUser = new User(this.model.username, this.model.email, this.model.password, this.model.address, this.model.smartMeterNo)
    this._config.addNewUser(newUser).subscribe(data => {
      console.log(data)
      if (data.Res){
        this.openSuccessDialog()
      }else{
        this.openFailDialog()
      }
    })
  }

  openSuccessDialog() {
    this.dialog.open(DialogElementsSuccess);
  }

  openFailDialog() {
    this.dialog.open(DialogElementsFail);
  }


}


@Component({
  selector: 'app-success-dialog',
  templateUrl: '../success-dialog/./success-dialog.component.html',
})
export class DialogElementsSuccess {}


@Component({
  selector: 'app-fail-dialog',
  templateUrl: '../fail-dialog/./fail-dialog.component.html',
})
export class DialogElementsFail {}