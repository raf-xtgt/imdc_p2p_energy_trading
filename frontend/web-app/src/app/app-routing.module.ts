import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

// external components
import { ToolbarComponent } from './toolbar/toolbar.component';
import {HomepageComponent} from './homepage/homepage.component';
import { LoginComponent } from './login/login.component';
import { SignupComponent } from './signup/signup.component';
import { ProfileComponent } from './profile/profile.component';
import { OrderComponent } from './order/order.component';
import { MarketpageComponent } from './marketpage/marketpage.component';
import { BidPageComponent } from './bid-page/bid-page.component';
import { ManageValidatorsComponent } from './manage-validators/manage-validators.component';
import { BlockchainPageComponent } from './blockchain-page/blockchain-page.component';
import { ManageClerksComponent } from './manage-clerks/manage-clerks.component';
import { NewMarketPageComponent } from './new-market-page/new-market-page.component';
import { EvPageComponent } from './ev-page/ev-page.component';
import { TransactionComponent } from './transaction/transaction.component';

const routes: Routes = [
  // when url is loaded, the user lands on the homepage
  {path:'', redirectTo:'/homepage', pathMatch:'full'},
  {path:'homepage', component: HomepageComponent},
  {path: 'login', component: LoginComponent}, 
  {path: 'register', component:SignupComponent},
  {path: 'profile', component:ProfileComponent},
  {path: 'order', component:OrderComponent},
  {path: 'market', component:NewMarketPageComponent},
  {path: 'bid', component:BidPageComponent},
  {path: 'manageValidators', component:ManageValidatorsComponent},
  {path: 'blockchain', component:BlockchainPageComponent},
  {path: 'manageClerks', component:ManageClerksComponent},
  {path: 'evPage', component:EvPageComponent},
  {path: 'transactions', component:TransactionComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

export const routingComponents = [HomepageComponent]
