import { Injectable } from '@angular/core'

import { SessionRequest, SessionResponse, SessionService as SessionApiService } from "../generated/swagger/rest"
import { CookieService } from 'ngx-cookie-service'
import { ReplaySubject } from 'rxjs/ReplaySubject'
import { Observable } from 'rxjs/Observable'

const SESSION_KEY = "sessionID"

@Injectable()
export class SessionService {

  private session: SessionResponse = null

  private sessionReplay: ReplaySubject<SessionResponse> = new ReplaySubject()

  constructor(private sessionApiService: SessionApiService,
              private cookieService: CookieService) {

    console.log("Initializing SessionService")

    const sessionID = cookieService.get(SESSION_KEY)
    const emptySession = this.isEmptySession(sessionID)

    const request: SessionRequest = { sessionID: sessionID, }
    sessionApiService.validateOrGenerate(request)
      .subscribe(response => {
        const sessionID = response.sessionID

        this.session = response
        this.cookieService.set(SESSION_KEY, sessionID)

        if (emptySession) {
          console.log(`New session is issued: ${sessionID}`)
        } else if (!emptySession && response.refreshed) {
          console.log(`Existing session is used or refreshed: ${sessionID}`)
        } else if (emptySession && !response.refreshed) {
          console.log("Can't find session. Issued new session", response)
        } else {
          console.warn("Unknown session handshake case", response)
        }

        this.sessionReplay.next(this.session)
      })
  }

  isEmptySession(sessionId: string): boolean {
    return sessionId === undefined || sessionId === null || sessionId == ""
  }

  subscribeSession(): Observable<SessionResponse> {
    return this.sessionReplay
  }

  clear() {
    this.cookieService.delete(SESSION_KEY)
    this.sessionReplay.next(null)
  }
}
