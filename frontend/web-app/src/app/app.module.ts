import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatSliderModule } from '@angular/material/slider';
import { FlexLayoutModule } from '@angular/flex-layout';
import {MatToolbarModule} from '@angular/material/toolbar'
import {MatIconModule} from '@angular/material/icon'; 
import {MatSidenavModule} from '@angular/material/sidenav'; 
import {MatButtonModule} from '@angular/material/button'; 
import {MatGridListModule} from '@angular/material/grid-list';
import {MatTabsModule} from '@angular/material/tabs';
import {MatFormFieldModule} from '@angular/material/form-field';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms'; 
import {MatCardModule} from '@angular/material/card';
import {MatDialogModule} from '@angular/material/dialog';
import {MatListModule} from '@angular/material/list';
//import { NgxChartsModule } from '@swimlane/ngx-charts';
import { ChartsModule } from 'ng2-charts';
// import {Swal} from 'sweetalert2/dist/sweetalert2.js'
import { SweetAlert2Module } from '@sweetalert2/ngx-sweetalert2';
// for the loading screen
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import {MatDividerModule} from '@angular/material/divider';
// for dropdown menu
import {MatMenuModule} from '@angular/material/menu';
// for table
import {MatTableModule} from '@angular/material/table';
// for pagination
import { MatPaginatorModule } from '@angular/material/paginator';




// external components
import { ToolbarComponent } from './toolbar/toolbar.component';
import { HomepageComponent } from './homepage/homepage.component';
import { MarketpageComponent } from './marketpage/marketpage.component';
import { LoginComponent } from './login/login.component';
import { SignupComponent } from './signup/signup.component';
import { SuccessDialogComponent } from './success-dialog/success-dialog.component';
import { FailDialogComponent } from './fail-dialog/fail-dialog.component';
import { ProfileComponent } from './profile/profile.component';
import { OrderComponent } from './order/order.component';
import { LineGraphComponent } from './line-graph/line-graph.component';
import { HouseholdPriceGraphComponent } from './household-price-graph/household-price-graph.component';
import { BidPageComponent } from './bid-page/bid-page.component';
import { TimerComponent } from './timer/timer.component';
import { ManageValidatorsComponent } from './manage-validators/manage-validators.component';
import { BlockchainPageComponent } from './blockchain-page/blockchain-page.component';
import { ManageClerksComponent } from './manage-clerks/manage-clerks.component';
import { ProfilePageComponent } from './profile-page/profile-page.component';
import { NewMarketPageComponent } from './new-market-page/new-market-page.component';
import { OpenRequestsPageComponent } from './open-requests-page/open-requests-page.component';
import { EvPageComponent } from './ev-page/ev-page.component';
import { TransactionComponent } from './transaction/transaction.component';
import { SellOrderComponent } from './sell-order/sell-order.component';

@NgModule({
  declarations: [
    AppComponent,
    ToolbarComponent,
    HomepageComponent,
    MarketpageComponent,
    LoginComponent,
    SignupComponent,
    SuccessDialogComponent,
    FailDialogComponent,
    ProfileComponent,
    OrderComponent,
    LineGraphComponent,
    HouseholdPriceGraphComponent,
    BidPageComponent,
    TimerComponent,
    ManageValidatorsComponent,
    BlockchainPageComponent,
    ManageClerksComponent,
    ProfilePageComponent,
    NewMarketPageComponent,
    OpenRequestsPageComponent,
    EvPageComponent,
    TransactionComponent,
    SellOrderComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    MatSliderModule,
    FlexLayoutModule,
    MatToolbarModule,
    MatIconModule,
    MatSidenavModule,
    MatButtonModule,
    MatGridListModule,
    MatTabsModule,
    MatFormFieldModule,
    HttpClientModule,
    FormsModule,
    MatCardModule,
    MatDialogModule,
    MatListModule,
    ChartsModule,
    SweetAlert2Module.forRoot(),
    MatProgressSpinnerModule,
    MatDividerModule,
    MatMenuModule,
    MatTableModule,
    MatPaginatorModule

  ],
  exports:[
    ChartsModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
