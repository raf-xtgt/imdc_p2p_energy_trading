import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
// The HttpClient service makes use of observables for all transactions. You must import the RxJS observable and operator symbols that appear in the example snippets.
import { Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';

import Swal from 'sweetalert2'

/** This file will handle all the different kinds of modals used by the platform
*/


@Injectable({
    providedIn: 'root'
  })

export class ModalService {
    constuctor(){

    }

    showConfirmationModal(title: string, confirmBtnTxt: string, successMsg: string, cancelMsg:string): boolean{
        let output: boolean = false
        Swal.fire({
            title: title,
            showDenyButton: false,
            showCancelButton: true,
            confirmButtonText: confirmBtnTxt,
            //denyButtonText: denyBtnTxt,
          }).then((result) => {
            /* Read more about isConfirmed, isDenied below */
            if (result.isConfirmed) {
              Swal.fire(successMsg, '', 'success')
              output = true
              
            } else if (result.isDismissed) {
              Swal.fire(cancelMsg, '', 'info')
              output = false
            }
          })
          //console.log(output)
          return output
        
    }
}