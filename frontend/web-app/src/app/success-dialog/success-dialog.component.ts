import { Component, OnInit } from '@angular/core';
import {MatDialog} from '@angular/material/dialog';

/** WHen user signs up successfully */
@Component({
  selector: 'app-success-dialog',
  templateUrl: './success-dialog.component.html',
  styleUrls: ['./success-dialog.component.css']
})
export class SuccessDialogComponent{

  constructor() { }

  // openDialog() {
  //   this.dialog.open(DialogElementsExampleDialog);
  // }
}


// @Component({
//   selector: 'app-success-dialog',
//   templateUrl: './success-dialog.component.html',
// })
// export class DialogElementsExampleDialog {}