import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-bid-page',
  templateUrl: './bid-page.component.html',
  styleUrls: ['./bid-page.component.css']
})
export class BidPageComponent implements OnInit {

  public bidEnergyInput: number = 0;
  constructor() { }

  ngOnInit(): void {
  }

  createBid(){}

}
