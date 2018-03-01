import {Component, NgModule, OnInit, OnDestroy} from '@angular/core'
import {MatButtonModule, MatToolbarModule, MatMenuModule} from '@angular/material'
import {RouterModule} from '@angular/router'

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.html',
  styleUrls: ['./navbar.scss']
})
export class Navbar implements OnInit, OnDestroy {
  constructor() {
  }

  ngOnInit() {
  }

  ngOnDestroy() {
  }
}

@NgModule({
  imports: [MatButtonModule, MatToolbarModule, MatMenuModule, RouterModule],
  exports: [Navbar],
  declarations: [Navbar],
})
export class NavbarModule {}
