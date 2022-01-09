import { Component, OnInit } from '@angular/core';
import {MatToolbarModule} from '@angular/material/toolbar'

@Component({
  selector: 'app-toolbar',
  templateUrl: './toolbar.component.html',
  styleUrls: ['./toolbar.component.css']
})
export class ToolbarComponent implements OnInit {

  constructor() { }

  ngOnInit(): void {
    
  }

  events: string[] = [];
  opened: boolean = false

  shouldRun = [/(^|\.)plnkr\.co$/, /(^|\.)stackblitz\.io$/].some(h => h.test(window.location.host));

}
