import { Injectable } from '@angular/core'
import { SessionService } from "./session.service"
import {
  BrowserHistory,
  BrowserHistoryService as BrowserHistoryApiService,
  BrowserHistoryWithPagination,
  SessionResponse,
} from "../generated/swagger/rest"

import * as moment from 'moment'
import 'moment-timezone'
import 'clientjs'
import { Observable } from 'rxjs/Observable'
import { ReplaySubject } from 'rxjs/ReplaySubject'
import { NotificationService } from "./notification.service"

@Injectable()
export class BrowserHistoryService {

  private browserHistorySendEvent: ReplaySubject<boolean> = new ReplaySubject()

  constructor(private sessionService: SessionService,
              private browserHistoryApiService: BrowserHistoryApiService,
              private notificationService: NotificationService) {
    sessionService.subscribeSession().subscribe((session: SessionResponse) => {
      if (!session) {
        return
      }

      console.info("Initializing BrowserHistoryService")

      // start initialization process for BrowserHistoryServcice
      // after session service is properly setup
      this.addOne().subscribe(_ => {
        this.browserHistorySendEvent.next(true)

        this.notificationService.displayInfo("Browser History", "Persisted")
      })
    })
  }

  private addOne(): Observable<BrowserHistory> {
    const client = new ClientJS()
    const browserHistory: BrowserHistory = {
      browserName: client.getBrowser(),
      browserVersion: client.getBrowserVersion(),
      osVersion: client.getOS(),
      osName: client.getOSVersion(),
      isMobile: client.isMobile(),
      clientTimezone: moment.tz.guess(),
      clientTimestamp: moment().toString(),
      language: client.getLanguage(),
      userAgent: window.navigator.userAgent,
    }

    return this.browserHistoryApiService.addOne(browserHistory)
  }

  findAll(itemCountPerPage: number, currentPageOffset: number,
          filterColumn: string, filterValue: string): Observable<BrowserHistoryWithPagination> {

    return this.browserHistoryApiService.findAll(itemCountPerPage, currentPageOffset,
      filterColumn, filterValue)
  }

  watchBrowserHistorySendEvent(): Observable<boolean> {
    return this.browserHistorySendEvent
  }
}
