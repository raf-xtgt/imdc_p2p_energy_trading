import { Component, OnInit, ViewChild } from '@angular/core';
import { ConfigService } from '../config.service';
import { MatTableDataSource } from '@angular/material/table';
import { MatPaginator } from '@angular/material/paginator';
// for the loading
import {ThemePalette} from '@angular/material/core';
import {ProgressSpinnerMode} from '@angular/material/progress-spinner';
import { PotentialClerks } from '../classes';
import Swal from 'sweetalert2'

@Component({
  selector: 'app-manage-clerks',
  templateUrl: './manage-clerks.component.html',
  styleUrls: ['./manage-clerks.component.css']
})
export class ManageClerksComponent implements OnInit {

  // columns to show on the page
  displayedColumns: string[] = ['username', 'userId', 'email', 'smartMeterNo', 'convert']
  
  // the list that will have all our data
  dataSource = new MatTableDataSource<PotentialClerks>(potentialClerks)


  constructor(private _config:ConfigService) { }

    // add the paginator
    @ViewChild(MatPaginator) paginator: MatPaginator | any

  ngOnInit(): void {
    this.getAllUsers()
  }

  // get all the normal users from the database
  getAllUsers(){
    potentialClerks = []
    this._config.getAllUsers().subscribe(data => {
      let response = JSON.parse(JSON.stringify(data))
      console.log("All user data from backend", response)

      for (let i=0; i<response.length; i++){
        let potClerk :PotentialClerks = {
          username: response[i].UserName,
          userId: response[i].UId,
          email: response[i].Email,
          smartMeterNo: response[i].SmartMeterNo,
          button1: "convert"
        }
        potentialClerks.push(potClerk)
      }
      console.log(potentialClerks)
      this.dataSource = new MatTableDataSource<PotentialClerks>(potentialClerks)
      this.dataSource.paginator = this.paginator

    })
  }

  // convert the user into a clerk
  convertToClerk(userId :string){
    console.log("Id of user to be made into clerk", userId)
    Swal.fire({
      title: 'Confirm Clerk Conversion',
      showDenyButton: false,
      showCancelButton: true,
      confirmButtonText: 'Confirm',
      //denyButtonText: denyBtnTxt,
    }).then((result) => {
      /* Read more about isConfirmed, isDenied below */
      if (result.isConfirmed) {

        this._config.convertToClerk(userId).subscribe(data => {
          Swal.fire('User successfully made into a clerk', '', 'success')  
        })
        
      } else if (result.isDismissed) {
        Swal.fire('Request Cancelled!', '', 'info')
      
      }
    })
  }
}

let potentialClerks : PotentialClerks[] = []
