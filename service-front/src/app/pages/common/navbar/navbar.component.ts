import { Component, NgModule, OnDestroy, OnInit } from '@angular/core'
import { MatButtonModule, MatMenuModule, MatToolbarModule } from '@angular/material'
import { RouterModule } from '@angular/router'

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.html',
  styleUrls: [ './navbar.scss' ]
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
  imports: [ MatButtonModule, MatToolbarModule, MatMenuModule, RouterModule ],
  exports: [ Navbar ],
  declarations: [ Navbar ],
})
export class NavbarModule {
}
