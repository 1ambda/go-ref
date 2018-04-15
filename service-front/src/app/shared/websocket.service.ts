import { Injectable } from '@angular/core'
import { Observable } from "rxjs/Observable"
import { BehaviorSubject } from "rxjs/BehaviorSubject"
import { ReplaySubject } from 'rxjs/ReplaySubject'
import { SessionService } from "./session.service"
import { SessionResponse } from "../generated/swagger/rest"
import { NotificationService } from "./notification.service"
import { WebSocketError, WebSocketResponseHeader } from 'app/generated/swagger/websocket'

const ReconnectingWebSocket = require('reconnecting-websocket')

@Injectable()
export class WebsocketService {
  private client
  private receiveQueue: ReplaySubject<any> = new ReplaySubject()
  private sendQueue: ReplaySubject<Object> = new ReplaySubject()
  private websocketConnectedEvent: BehaviorSubject<boolean> = new BehaviorSubject(false)

  constructor(private sessionService: SessionService,
              private notificationService: NotificationService) {
    this.sessionService.subscribeSession().subscribe((session: SessionResponse) => {
      console.info(`Initializing WebsocketService (session: ${session.sessionID})`)

      if (this.client) {
        return
      }

      this.client = new ReconnectingWebSocket(ENDPOINT_SERVICE_GATEWAY_WS)

      this.client.onerror = (error) => {
        let message = error.message
        if (error.message === "" || error.message === undefined || error.message === null) {
          message = "Connection refused"
        }

        console.error('websocket: `onerror`', message)
        // this.notificationService.displayError("Error (WS)", message)
      }

      this.client.onclose = () => {
        console.warn('websocket: `onclose` (will reconnect)')

        this.websocketConnectedEvent.next(false)
        this.notificationService.displayWarn("Disconnected (WS)", "will reconnect")
      }

      this.client.onopen = () => {
        console.debug("websocket: `onopen`")

        this.websocketConnectedEvent.next(true)
        this.notificationService.displaySuccess("Connected (WS)", "")

        this.sendQueue.subscribe(data => {
          this.client.send(JSON.stringify(data))
        })
      }

      this.client.onmessage = (response) => {
        const parsed = JSON.parse(response.data)
        this.receiveQueue.next(parsed)


        if (parsed.header.responseType === WebSocketResponseHeader.ResponseTypeEnum.Error) {
          const errorResponse: WebSocketError = parsed.header.error
          console.log(errorResponse)

          const message = `${errorResponse.type} (${errorResponse.code})`
          this.notificationService.displayError("Error (WS Server)", message)

          if (errorResponse.type === WebSocketError.TypeEnum.InvalidSession) {
            this.sessionService.clear()
          }
        }
      }
    })
  }

  public send(data: Object) {
    if (this.client.readyState !== this.client.OPEN) {
      console.warn("websocket is not connected. message will be queued")
    }

    this.sendQueue.next(data)
  }


  public watch(targetResponseType: any): Observable<any> {
    return this.receiveQueue.filter((response: any) => {
      const eventType = response.header.responseType

      if (eventType === targetResponseType) {
        console.debug(targetResponseType, response)
      }

      return (eventType === targetResponseType)
    })
  }

  public watchWebsocketConnected(): Observable<boolean> {
    return this.websocketConnectedEvent
  }
}
