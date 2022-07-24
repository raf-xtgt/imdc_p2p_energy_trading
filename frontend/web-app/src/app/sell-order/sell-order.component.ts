import { Component, OnInit, ViewChild } from '@angular/core';
import { ConfigService } from '../config.service';
import { JWTService } from '../userAuth.service';
import { TransactionInfo } from '../classes';
import { MatTableDataSource } from '@angular/material/table';
import { MatPaginator } from '@angular/material/paginator';

@Component({
  selector: 'app-sell-order',
  templateUrl: './sell-order.component.html',
  styleUrls: ['./sell-order.component.css']
})
export class SellOrderComponent implements OnInit {

  displayedColumns: string[] = ['buyer', 'seller', 'energy', 'fiat', 'transaction']
  dataSource = new MatTableDataSource<TransactionInfo>(trnData)
  
  constructor(private _jwtServ:JWTService, private _config:ConfigService) { }
  private isAuth:boolean = false // to check whether the user is logged in or not
  private _sellerId :string = ""
  // add the paginator
  @ViewChild(MatPaginator) paginator: MatPaginator | any



  ngOnInit(): void {
    this._jwtServ.verifyToken().subscribe(data => {
      let response = JSON.parse(JSON.stringify(data))
      if (data != null){
        //console.log("Verified Token", data)
        this.isAuth = true
        this._sellerId = response.User.UId
        this.buyOrderRequest(this._sellerId)    
      }
    })
  }


  buyOrderRequest(userId: string){
    this._config.getUserSellRequests(this._sellerId).subscribe(data => {

     let response = JSON.parse(JSON.stringify(data))
      console.log("list of transactions", response.Transactions)
      let allUserTransactions = response.Transactions
      trnData = []
      for (let i=0; i<allUserTransactions.length; i++){
        let transaction = allUserTransactions[i]
        let allBidders = transaction.AuctionBids
        for (let j=0; j<allBidders.length; j++){
          let seller = allBidders[j]
          let sellerId = seller.SellerId

          if (sellerId == this._sellerId){
            console.log("Match Found\nUser Id", sellerId , "\nSeller Id", this._sellerId)
            let data:TransactionInfo={
              buyerId: transaction.BuyerId,
              sellerId: this._sellerId,
              energyTraded: (seller.OptEnFromSeller).toFixed(2),
              fiatTraded: (seller.OptSellerReceivable).toFixed(2),
              transactionId: transaction.TId
            }
            trnData.push(data)
          }else{
            console.log("No Match")
          }

          

        }
      }

      // instantiate list
      this.dataSource = new MatTableDataSource<TransactionInfo>(trnData)
      this.dataSource.paginator = this.paginator

    })
  }
}

let trnData : TransactionInfo[] = []