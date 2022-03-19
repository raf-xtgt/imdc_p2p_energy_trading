import { dashCaseToCamelCase } from '@angular/compiler/src/util';
import { Component, OnInit } from '@angular/core';
import {FormControl, Validators} from '@angular/forms';
import Swal from 'sweetalert2'
import { Router } from '@angular/router';

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

  constructor(private _config:ConfigService, private router: Router) { }

  ngOnInit(): void {
  }
  
  model = new User("", "","","",0, "", "")
  submitted = false;


  /** Method to add user to the database */
  addUser(){
    this.submitted = true;
    //console.log(this.model.fullname)
    //console.log(this.model.email)
    let newUser = new User(this.model.username, this.model.email, this.model.password, this.model.address, this.model.smartMeterNo, this.model.uId, "normal")
    this._config.addNewUser(newUser).subscribe(data => {
      console.log("Data from backend", data)
      if (data.Res){
        this.openSuccessDialog()
      }else{
        this.openFailDialog()
      }
    })
  }

  openSuccessDialog() {
    Swal.fire({
      icon: 'success',
      title: 'Sign Up Successful',
      text: 'Click OK and login to the site',
    }).then((result)=>{
      if (result.isConfirmed){
        this.router.navigateByUrl('/login')
      }
    })
  }

  openFailDialog() {
    Swal.fire({
      icon: 'error',
      title: 'Oops...',
      text: 'Username and/or email already being used. Please enter a different username and/or email',
    }).then((result)=>{
      if (result.isConfirmed){
        this.router.navigateByUrl('/register')
      }
    })
  }


}
