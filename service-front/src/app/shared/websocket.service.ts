import {Injectable,} from '@angular/core'
import {Observable} from "rxjs/Observable"
import {ReplaySubject} from "rxjs/ReplaySubject"
import {WebSocketRealtimeResponse} from "../generated/grpc/gateway_pb";

// TODO: DEV, PROD, ...
const WebSocketConnectionUrl = "ws://localhost:50001/ws"

@Injectable()
export class WebsocketService {
  private ws: WebSocket
  private receiveQueue: ReplaySubject<any>
  private sendQueue: ReplaySubject<Object>

  constructor() {
    const ws = new WebSocket(WebSocketConnectionUrl)
    this.ws = ws
    this.sendQueue = new ReplaySubject()
    this.receiveQueue = new ReplaySubject()

    ws.onopen = () => {
      this.sendQueue.subscribe(data => {
        this.ws.send(JSON.stringify(data))
      })
    }

    ws.onmessage = (response) => {
      const parsed = JSON.parse(response.data)
      this.receiveQueue.next(parsed)
    }

    ws.onerror = function (error) {
      console.error("websocket error", error)
    }

    ws.onclose = function () {
      console.error("websocket is disconnected")
    }

    // TODO: reconnect
  }

  public send(data: Object) {
    if (this.ws.readyState !== WebSocket.OPEN) {
      console.warn("websocket is not connected. message will be queued")
    }

    this.sendQueue.next(data)
  }

  public watch(eventType: number): Observable<any> {
    return this.receiveQueue.filter((response: any) => {
      console.log(response)
      return response.header.eventType === eventType
    })
  }
}
