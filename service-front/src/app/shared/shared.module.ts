import { CommonModule } from '@angular/common'
import { ModuleWithProviders, NgModule } from '@angular/core'
import { FormsModule } from '@angular/forms'

import { FlexLayoutModule } from '@angular/flex-layout'
import { MaterialModule } from './material.module'
import { CookieService } from 'ngx-cookie-service'
import { SimpleNotificationsModule } from 'angular2-notifications'

import { BrowserHistoryService } from "./browser-history.service"
import { NotificationService } from "./notification.service"
import { WebsocketService } from "./websocket.service"
import { GeoLocationService } from "./geo-location.service"
import { SessionService } from "./session.service"


@NgModule({
  imports: [
    CommonModule,
    MaterialModule,
    FlexLayoutModule,
  ],
  declarations: [],
  exports: [
    CommonModule,
    FormsModule,
    MaterialModule,
    FlexLayoutModule,
  ],
})
export class SharedModule {
  static forRoot(): ModuleWithProviders {
    return {
      ngModule: SharedModule,
      providers: [
        CookieService,

        SessionService,
        WebsocketService,
        GeoLocationService,
        BrowserHistoryService,
        NotificationService,
      ]
    }
  }

}
