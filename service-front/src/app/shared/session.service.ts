import { Injectable } from '@angular/core'

import { SessionRequest, SessionResponse, SessionService as SessionApiService } from "../generated/swagger/rest"
import { CookieService } from 'ngx-cookie-service'

const SESSION_KEY = "sessionID"

@Injectable()
export class SessionService {

  private session: SessionResponse

  constructor(private sessionApiService: SessionApiService,
              private cookieService: CookieService) {

    const sessionID = cookieService.get(SESSION_KEY)
    const emptySession = this.isEmptySession(sessionID)

    const request: SessionRequest = { sessionID: sessionID, }
    sessionApiService.validateOrGenerate(request)
      .subscribe(response => {
        this.session = response
        const sessionID = response.sessionID

        if (emptySession) {
          console.log(`New session is issued: ${sessionID}`)
        } else if (!emptySession && response.refreshed) {
          console.log(`Existing session is used or refreshed: ${sessionID}`)
        } else if (emptySession && !response.refreshed) {
          console.log("Can't find session. Issued new session", response)
        } else {
          console.warn("Unknown session handshake case", response)
        }

        this.cookieService.set(SESSION_KEY, this.session.sessionID)
      })
  }

  isEmptySession(sessionId: string): boolean {
    return sessionId === undefined || sessionId === null || sessionId == ""
  }

}
