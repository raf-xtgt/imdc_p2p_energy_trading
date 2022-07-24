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
  selector: 'app-transaction',
  templateUrl: './transaction.component.html',
  styleUrls: ['./transaction.component.css']
})
export class TransactionComponent implements OnInit {

  displayedColumns: string[] = ['buyer', 'seller', 'energy', 'fiat']
  dataSource = new MatTableDataSource<Block>(trnData)
  
  constructor() { }

  ngOnInit(): void {
  }

 
}

let trnData : Block[] = []
