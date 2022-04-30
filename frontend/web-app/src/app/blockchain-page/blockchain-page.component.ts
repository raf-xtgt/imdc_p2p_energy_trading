import { Component, OnInit, ViewChild } from '@angular/core';
import { ConfigService } from '../config.service';
import { JWTService } from '../userAuth.service';
import { Block } from '../classes';
import { MatTableDataSource } from '@angular/material/table';
import { MatPaginator } from '@angular/material/paginator';
// for the loading
import {ThemePalette} from '@angular/material/core';
import {ProgressSpinnerMode} from '@angular/material/progress-spinner';
import {MatDialog} from '@angular/material/dialog';
import Swal from 'sweetalert2'


@Component({
  selector: 'app-blockchain-page',
  templateUrl: './blockchain-page.component.html',
  styleUrls: ['./blockchain-page.component.css']
})
export class BlockchainPageComponent implements OnInit {

  displayedColumns: string[] = ['index', 'hash', 'nonce', 'prevHash']
  dataSource = new MatTableDataSource<Block>(blockData)
  
  constructor(private _config:ConfigService, private _jwtServ:JWTService, public dialog: MatDialog) { }
  isValidator:boolean = false;
  isClerk :boolean = false;
  // loading before updated blockchain is available
  public isLoading: boolean = true;
  color: ThemePalette = 'primary';
  mode: ProgressSpinnerMode = 'indeterminate';
  value = 100;

  
  public blockInfo :any = [] // to hold information for all the blocks

  // add the paginator
  @ViewChild(MatPaginator) paginator: MatPaginator | any

  ngOnInit(): void {

  // uncomment this to re-create the genesis block
    // this.createGenesis();
  this.getUserType()


  }


  createGenesis(){
    this._config.createGenesisBlock().subscribe(data=>{
      console.log("made genesis block")
    })
  }

  updateBlockchain (){
    this._config.updateBlockchain().subscribe(data => {
      console.log("Updated blockchain")
      this.getBlockchain()
    })
  }

  // will invoke integrity check. The backend handles the five new block check
  initINTCheck(){
    this._config.initClerkINTChk().subscribe(data => {
      console.log("Invoke integrity check", data)
      this.getBlockchain()
    })
  }

  getBlockchain(){
    blockData = []
    this._config.getCurrentBlockchain().subscribe(data => {
      
      let response = JSON.parse(JSON.stringify(data))
      console.log("current blockchain", response.Blockchain)
      for (let i=0; i<response.Blockchain.length; i++){
        let block = response.Blockchain[i]
        this.blockInfo.push(block)
        let data :Block = {
          index: block.Index,
          data:block.Data,
          hash:block.Hash,
          nonce:block.Nonce,
          prevHash:block.PrevHash,
        }
        blockData.push(data)
      }
      
      // instantiate list
      this.dataSource = new MatTableDataSource<Block>(blockData)
      this.dataSource.paginator = this.paginator
      
    })
    this.isLoading=false
  }

  getUserType(){
    // check if the jwt is stored in local storage or not
    this._jwtServ.verifyToken().subscribe(data => {
      //console.log("Verified Token", data)
      let response = JSON.parse(JSON.stringify(data))
      console.log(response.User)

      if (data !=null){

        if (response.User.Type == "validator"){
          this.isValidator = true
          // update the blockchain for validator
          this.updateBlockchain()
        }

        else if (response.User.Type == "clerk"){
          // invoke integrity check
          this.initINTCheck()
          //this.getBlockchain()
        }
        else{
          this.isValidator = false
          // only get the blockchain for normal user
          this.getBlockchain()
        }
      }
    })
  }

  // to show the transaction data when the info button is clicked
  showBlockData(blockHash :string){


    //console.log(blockHash)
    let selectedBlock // to hold the data of the selected block
    for (let i=0; i<this.blockInfo.length; i++){
      let block = this.blockInfo[i]
      if (block.Hash == blockHash){
        selectedBlock = block.Data
      }
    }
  
    //console.log(selectedBlock)
    let buyerId, totalEnTraded, totalFiat, tnbIncome, sellers 

    // the string that will be shown on html
    let htmlStr = ""
    for (let j=0; j<selectedBlock.length; j++){
      let info = selectedBlock[j]
      //console.log(info)
      buyerId = info.BuyerId
      totalEnTraded = info.BuyerEnReceivableFromAuction + info.BuyerEnReceivableFromTNB
      totalFiat = info.BuyerPayable // total money in this transaction
      tnbIncome = info.TNBReceivable
      sellers = info.AuctionBids
      
      htmlStr += "<br><b>Buyer:</b> "+buyerId+"\n<br>"
      htmlStr += "<b>Total Energy Traded:</b> "+totalEnTraded.toFixed(2)+" kWH\n<br>"
      htmlStr += "<b>Total Fiat Traded:</b> RM "+totalFiat.toFixed(2)+"\n<br>"
      htmlStr += "<b>TNB Income:</b> RM  "+tnbIncome.toFixed(2)+"\n<br>"
    }



    // the modal
    Swal.fire({
      title: 'Block Info<br>Hash:'+blockHash,
      icon: 'info',
      html: htmlStr,
      showCloseButton: true,
      
    })

  }


}

let blockData : Block[] = []


export class DialogContentExampleDialog {}