import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

// external components
import { ToolbarComponent } from './toolbar/toolbar.component';
import {HomepageComponent} from './homepage/homepage.component';
import { LoginComponent } from './login/login.component';
import { SignupComponent } from './signup/signup.component';
import { ProfileComponent } from './profile/profile.component';
// import { TradeComponent } from './trade/trade.component';

const routes: Routes = [
  // when url is loaded, the user lands on the homepage
  {path:'', redirectTo:'/homepage', pathMatch:'full'},
  {path:'homepage', component: HomepageComponent},
  {path: 'login', component: LoginComponent}, 
  {path: 'register', component:SignupComponent},
  {path: 'profile', component:ProfileComponent},
  // {path: 'trade', component:TradeComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

export const routingComponents = [HomepageComponent]
