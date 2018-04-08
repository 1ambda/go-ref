import { Routes } from '@angular/router'
import { HomeComponent } from './pages/home'
import { LocationComponent } from './pages/location'
import { AboutComponent } from './pages/about'
import { NoContentComponent } from './pages/no-content'

export const ROUTES: Routes = [
  { path: '', component: HomeComponent },
  { path: 'home', component: HomeComponent },
  { path: 'location', component: LocationComponent },
  { path: 'about', component: AboutComponent },
  { path: '**', component: NoContentComponent },
]
