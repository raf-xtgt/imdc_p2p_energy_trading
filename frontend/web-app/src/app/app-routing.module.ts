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

const routes: Routes = [
  // when url is loaded, the user lands on the homepage
  {path:'', redirectTo:'/homepage', pathMatch:'full'},
  {path:'homepage', component: HomepageComponent},
  {path: 'login', component: LoginComponent}, 
  {path: 'register', component:SignupComponent},
  {path: 'profile', component:ProfileComponent},
  {path: 'order', component:OrderComponent},
  {path: 'market', component:MarketpageComponent},
  {path: 'bid', component:BidPageComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

export const routingComponents = [HomepageComponent]
