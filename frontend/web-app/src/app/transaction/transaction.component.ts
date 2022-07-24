import { Component, OnInit, ViewChild } from '@angular/core';
import { ConfigService } from '../config.service';
import { JWTService } from '../userAuth.service';
import { TransactionInfo } from '../classes';
import { MatTableDataSource } from '@angular/material/table';
import { MatPaginator } from '@angular/material/paginator';

@Component({
  selector: 'app-transaction',
  templateUrl: './transaction.component.html',
  styleUrls: ['./transaction.component.css']
})
export class TransactionComponent implements OnInit {

  displayedColumns: string[] = ['buyer', 'seller', 'energy', 'fiat', 'transaction']
  dataSource = new MatTableDataSource<TransactionInfo>(trnData)
  
  constructor(private _jwtServ:JWTService, private _config:ConfigService) { }
  private isAuth:boolean = false // to check whether the user is logged in or not
  private _buyerId :string = ""
  // add the paginator
  @ViewChild(MatPaginator) paginator: MatPaginator | any



  ngOnInit(): void {
    this._jwtServ.verifyToken().subscribe(data => {
      let response = JSON.parse(JSON.stringify(data))
      if (data != null){
        //console.log("Verified Token", data)
        this.isAuth = true
        this._buyerId = response.User.UId
        this.buyOrderRequest(this._buyerId)    
      }
    })
  }

  buyOrderRequest(userId: string){
    this._config.getUserBuyRequests(this._buyerId).subscribe(data => {

     let response = JSON.parse(JSON.stringify(data))
      //console.log("list of transactions", response.Transactions)
      let allUserTransactions = response.Transactions
      trnData = []
      for (let i=0; i<allUserTransactions.length; i++){
        let transaction = allUserTransactions[i]
        let allBidders = transaction.AuctionBids
        for (let j=0; j<allBidders.length; j++){
          let seller = allBidders[j]
          let sellerId = seller.SellerId

          let data:TransactionInfo={
            buyerId: this._buyerId,
            sellerId: sellerId,
            energyTraded: (seller.OptEnFromSeller).toFixed(2),
            fiatTraded: (seller.OptSellerReceivable).toFixed(2),
            transactionId: transaction.TId
          }
          trnData.push(data)

        }
      }

      // instantiate list
      this.dataSource = new MatTableDataSource<TransactionInfo>(trnData)
      this.dataSource.paginator = this.paginator

    })
  }

 
}

let trnData : TransactionInfo[] = []
