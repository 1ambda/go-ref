import { Component, NgModule, OnDestroy, OnInit } from '@angular/core'
import { RouterModule } from '@angular/router'
import { SharedModule } from "../../../shared/shared.module"

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
  imports: [ SharedModule, RouterModule ],
  exports: [ Navbar ],
  declarations: [ Navbar ],
})
export class NavbarModule {
}
