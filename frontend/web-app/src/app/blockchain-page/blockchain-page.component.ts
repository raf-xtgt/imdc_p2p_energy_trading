import { Component, OnInit, ViewChild } from '@angular/core';
import { ConfigService } from '../config.service';
import { JWTService } from '../userAuth.service';
import { Block } from '../classes';
import { MatTableDataSource } from '@angular/material/table';
import { MatPaginator } from '@angular/material/paginator';
// for the loading
import {ThemePalette} from '@angular/material/core';
import {ProgressSpinnerMode} from '@angular/material/progress-spinner';


@Component({
  selector: 'app-blockchain-page',
  templateUrl: './blockchain-page.component.html',
  styleUrls: ['./blockchain-page.component.css']
})
export class BlockchainPageComponent implements OnInit {

  displayedColumns: string[] = ['index', 'hash', 'nonce', 'prevHash']
  dataSource = new MatTableDataSource<Block>(blockData)
  constructor(private _config:ConfigService, private _jwtServ:JWTService) { }
  isValidator:boolean = false;
  // loading before updated blockchain is available
  public isLoading: boolean = true;
  color: ThemePalette = 'primary';
  mode: ProgressSpinnerMode = 'indeterminate';
  value = 100;


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

  getBlockchain(){
    blockData = []
    this._config.getCurrentBlockchain().subscribe(data => {
      
      let response = JSON.parse(JSON.stringify(data))
      console.log("current blockchain", response.Blockchain)
      for (let i=0; i<response.Blockchain.length; i++){
        let block = response.Blockchain[i]
        let data :Block = {
          index: block.Index,
          data:block.Data,
          hash:block.Hash,
          nonce:block.Nonce,
          prevHash:block.PrevHash

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
      //console.log(response.User)

      if (data !=null){

        if (response.User.Type == "validator"){
          this.isValidator = true
          // update the blockchain for validator
          this.updateBlockchain()
        }
        else{
          this.isValidator = false
          // only get the blockchain for normal user
          this.getBlockchain()
        }
      }
    })
  }

}

let blockData : Block[] = []