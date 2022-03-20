import { Component, OnInit } from '@angular/core';
import {Validator} from '../classes';
// import the custom http service module
import { ConfigService } from '../config.service';
import { Router } from '@angular/router';
import Swal from 'sweetalert2'

@Component({
  selector: 'app-manage-validators',
  templateUrl: './manage-validators.component.html',
  styleUrls: ['./manage-validators.component.css']
})
export class ManageValidatorsComponent implements OnInit {

  constructor(private _config:ConfigService, private router: Router) { }

  ngOnInit(): void {
  }
  
  model = new Validator("", "","","","", "", "", 0)
  submitted = false;


    /** Method to add user to the database */
    addValidator(){
      this.submitted = true;
      //console.log(this.model.fullname)
      //console.log(this.model.email)
      let newValidator = new Validator(this.model.username, this.model.email, this.model.password, this.model.address, this.model.uId, "validator", this.model.fullName, this.model.ICNum)
      this._config.addNewValidator(newValidator).subscribe(data => {
        console.log("Data from backend for validator addition", data)
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
        title: 'Validator Addition Successful',
        text: 'Click OK and login to the site',
      }).then((result)=>{
        if (result.isConfirmed){
          this.router.navigateByUrl('/manageValidators')
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
          this.router.navigateByUrl('/manageValidators')
        }
      })
    }

}
