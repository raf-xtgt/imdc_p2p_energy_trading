import { Component, OnInit } from '@angular/core';
import { ConfigService } from '../config.service';

@Component({
  selector: 'app-blockchain-page',
  templateUrl: './blockchain-page.component.html',
  styleUrls: ['./blockchain-page.component.css']
})
export class BlockchainPageComponent implements OnInit {

  constructor(private _config:ConfigService) { }


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
    })
  }

}
