import { Component, NgModule, OnDestroy, OnInit } from '@angular/core'
import { RouterModule } from '@angular/router'
import { SharedModule } from "../../../shared/shared.module"
import { SessionService } from "../../../shared/session.service"
import { Subscription } from 'rxjs/Subscription'
import { SessionResponse } from "../../../generated/swagger/rest"

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: [ './navbar.component.scss' ]
})
export class Navbar implements OnInit, OnDestroy {
  private subscriptions: Array<Subscription> = []
  private sessionID = null

  constructor(private sessionService: SessionService) {
  }

  ngOnInit() {
    this.subscriptions.push(this.subscribeSession())
  }

  ngOnDestroy() {
    this.subscriptions.forEach(sub => {
      sub.unsubscribe()
    })
  }

  private subscribeSession(): Subscription {
    return this.sessionService.subscribeSession()
      .subscribe((session: SessionResponse) => {
        if (session) {
          this.sessionID = session.sessionID
        } else {
          this.sessionID = null
        }
      })
  }
}

@NgModule({
  imports: [ SharedModule, RouterModule ],
  exports: [ Navbar ],
  declarations: [ Navbar ],
})
export class NavbarModule {
}
