import { Component, OnInit } from '@angular/core';
import {Validator} from '../classes';
// import the custom http service module
import { ConfigService } from '../config.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-manage-validators',
  templateUrl: './manage-validators.component.html',
  styleUrls: ['./manage-validators.component.css']
})
export class ManageValidatorsComponent implements OnInit {

  constructor(private _config:ConfigService, private router: Router) { }

  ngOnInit(): void {
  }
  
  model = new Validator("", "","","",0, "", "", "", 0)
  submitted = false;


}
