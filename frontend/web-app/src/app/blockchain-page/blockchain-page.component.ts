import { Component, OnInit, ViewChild } from '@angular/core';
import { ConfigService } from '../config.service';
import { Block } from '../classes';
import { MatTableDataSource } from '@angular/material/table';
import { MatPaginator } from '@angular/material/paginator';
@Component({
  selector: 'app-blockchain-page',
  templateUrl: './blockchain-page.component.html',
  styleUrls: ['./blockchain-page.component.css']
})
export class BlockchainPageComponent implements OnInit {

  displayedColumns: string[] = ['index', 'hash', 'nonce', 'prevHash']
  dataSource = new MatTableDataSource<Block>(blockData)
  constructor(private _config:ConfigService) { }

  // add the paginator
  @ViewChild(MatPaginator) paginator: MatPaginator | any

  ngOnInit(): void {

  // uncomment this to re-create the genesis block
    // this.createGenesis();

  this.updateBlockchain()

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
      this.dataSource = new MatTableDataSource<Block>(blockData)
      this.dataSource.paginator = this.paginator
    })
  }

}

const blockData : Block[] = []