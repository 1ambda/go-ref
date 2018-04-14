import { Injectable } from '@angular/core'
import { Observable } from "rxjs/Observable"
import { BehaviorSubject } from "rxjs/BehaviorSubject"
import { ReplaySubject } from 'rxjs/ReplaySubject'
import { SessionService } from "./session.service"
import { SessionResponse } from "../generated/swagger/rest"

const ReconnectingWebSocket = require('reconnecting-websocket')

@Injectable()
export class WebsocketService {
  private client
  private receiveQueue: ReplaySubject<any> = new ReplaySubject()
  private sendQueue: ReplaySubject<Object> = new ReplaySubject()

  constructor(sessionService: SessionService) {
    sessionService.subscribeSession().subscribe((session: SessionResponse) => {
      console.info(`Initializing WebsocketService (session: ${session.sessionID})`)

      this.client = new ReconnectingWebSocket(ENDPOINT_SERVICE_GATEWAY_WS)

      this.client.onerror = (error) => {
        console.error('websocket: `onerror`', error)
      }

      this.client.onclose = () => {
        console.warn('websocket: `onclose` (will reconnect)')
      }

      this.client.onopen = () => {
        console.debug("websocket: `onopen`")

        this.sendQueue.subscribe(data => {
          this.client.send(JSON.stringify(data))
        })
      }

      this.client.onmessage = (response) => {
        const parsed = JSON.parse(response.data)
        this.receiveQueue.next(parsed)
      }
    })
  }

  public send(data: Object) {
    if (this.client.readyState !== this.client.OPEN) {
      console.warn("websocket is not connected. message will be queued")
    }

    this.sendQueue.next(data)
  }


  public watch(responseType: any): Observable<any> {
    return this.receiveQueue.filter((response: any) => {
      const isTargetResponseType = response.header.responseType === responseType

      if (isTargetResponseType) {
        console.debug(responseType, response)
      }

      return isTargetResponseType
    })
  }
}
