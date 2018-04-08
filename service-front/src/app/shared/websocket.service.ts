import { Injectable, } from '@angular/core'
import { Observable } from "rxjs/Observable"
import { ReplaySubject } from "rxjs/ReplaySubject"

// TODO: DEV, PROD, ...

const ReconnectingWebSocket = require('reconnecting-websocket')

const SERVER_URL = "ws://localhost:50001/endpoint"

@Injectable()
export class WebsocketService {
  private client
  private receiveQueue: ReplaySubject<any>
  private sendQueue: ReplaySubject<Object>

  constructor() {
    this.client = new ReconnectingWebSocket(SERVER_URL)
    this.sendQueue = new ReplaySubject()
    this.receiveQueue = new ReplaySubject()

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
  }

  public send(data: Object) {
    if (this.client.readyState !== this.client.OPEN) {
      console.warn("websocket is not connected. message will be queued")
    }

    this.sendQueue.next(data)
  }


  public watch(responseType: any): Observable<any> {
    return this.receiveQueue.filter((response: any) => {
      console.debug(responseType, response)
      return response.header.responseType === responseType
    })
  }
}
