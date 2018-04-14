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

@Injectable()
export class BrowserHistoryService {

  private browserHistorySendEvent: ReplaySubject<boolean> = new ReplaySubject()

  constructor(private sessionService: SessionService,
              private browserHistoryApiService: BrowserHistoryApiService) {
    sessionService.subscribeSession().subscribe((session: SessionResponse) => {
      console.info(`Initializing BrowserHistoryService (session: ${session.sessionID})`)

      // start initialization process for BrowserHistoryServcice
      // after session service is properly setup
      this.addOne().subscribe(_ => {
        this.browserHistorySendEvent.next(true)
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

  findAll(itemCountPerPage: number,
          currentPageOffset: number): Observable<BrowserHistoryWithPagination> {

    return this.browserHistoryApiService.findAll(itemCountPerPage, currentPageOffset)
  }

  watchBrowserHistorySendEvent(): Observable<boolean> {
    return this.browserHistorySendEvent
  }
}
