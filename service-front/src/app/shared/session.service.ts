import { Injectable } from '@angular/core'

import { SessionRequest, SessionResponse, SessionService as SessionApiService } from "../generated/swagger/rest"
import { CookieService } from 'ngx-cookie-service'

const SESSION_KEY = "sessionID"

@Injectable()
export class SessionService {

  private session: SessionResponse

  constructor(private sessionApiService: SessionApiService,
              private cookieService: CookieService) {

    const sessionID = cookieService.get("SESSION_KEY")

    if (this.isEmptySession(sessionID)) {
      console.log("Session is empty")
    }

    const request: SessionRequest = { sessionID: sessionID, }
    sessionApiService.validateOrGenerate(request)
      .subscribe(response => {
        this.session = response

        console.log(response.refreshed)
        console.log(response.refreshCount)
      })
  }

  isEmptySession(sessionId: string): boolean {
    return sessionId === undefined || sessionId === null || sessionId == ""
  }

}
