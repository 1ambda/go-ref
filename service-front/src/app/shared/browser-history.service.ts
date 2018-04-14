import { Injectable } from '@angular/core'
import { SessionService } from "./session.service"
import {
  BrowserHistory,
  BrowserHistoryService as BrowserHistoryApiService,
  BrowserHistoryWithPagination,
  SessionResponse,
} from "../generated/swagger/rest"

import * as moment from 'moment'
import 'clientjs'
import { Observable } from 'rxjs/Observable'

@Injectable()
export class BrowserHistoryService {

  private initialized = false

  constructor(private sessionService: SessionService,
              private browserHistoryApiService: BrowserHistoryApiService) {
    sessionService.subscribeSession().subscribe((session: SessionResponse) => {
      console.info(`Initializing BrowserHistoryService (session: ${session.sessionID})`)
    })
  }

  addOne(): Observable<BrowserHistory> {
    const client = new ClientJS()
    const browserHistory: BrowserHistory = {
      browserName: client.getBrowser(),
      browserVersion: client.getBrowserVersion(),
      osVersion: client.getOS(),
      osName: client.getOSVersion(),
      isMobile: client.isMobile() ? 'Y' : 'N',
      timezone: client.getTimeZone(),
      timestamp: moment.utc().unix().toString(),
      language: client.getLanguage(),
      userAgent: window.navigator.userAgent,
    }

    this.initialized = true

    return this.browserHistoryApiService.addOne(browserHistory)
  }

  findAll(itemCountPerPage: number,
          currentPageOffset: number): Observable<BrowserHistoryWithPagination> {


    return this.browserHistoryApiService.findAll(itemCountPerPage, currentPageOffset)
  }

}
