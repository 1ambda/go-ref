import { Component, OnDestroy, OnInit } from '@angular/core'

import { WebsocketService } from "../../shared/websocket.service"
import { WebSocketRealtimeResponse, WebSocketResponseHeader } from "../../generated/swagger/websocket"
import { Subscription } from 'rxjs/Subscription'
import { GeoLocationService } from "../../shared/geo-location.service"
import { SessionService } from "../../shared/session.service"
import { BrowserHistoryService } from "../../shared/browser-history.service"

let initialized = false

@Component({
  selector: 'home',
  providers: [],
  styleUrls: [ './home.component.css' ],
  templateUrl: './home.component.html'
})
export class HomeComponent implements OnInit, OnDestroy {

  /**
   * indicates websocket is connected or not
   */
  websocketConnected = false

  /**
   * table related variables
   */
  rows = []
  columns = [
    { name: 'id', prop: 'id' },
    { name: 'browser_name', prop: 'browserName' },
    { name: 'browser_version', prop: 'browserVersion' },
    { name: 'os_name', prop: 'osName' },
    { name: 'os_version', prop: 'osVersion' },
    { name: 'is_mobile', prop: 'isMobile' },
    { name: 'language', prop: 'language' },
    { name: 'timestamp', prop: 'timestamp' },
    { name: 'timezone', prop: 'timezone' },
    { name: 'uuid', prop: 'uuid' },
    { name: 'user_agent', prop: 'userAgent' },
  ]

  /**
   * pagination related variables
   */
  itemCountPerPage = 10
  totalItemCount = 0
  currentPageOffset = 0 // `page_number - 1`
  isTableLoading = true

  /**
   * server-side streamed variables
   */
  currentTotalAccessCount = "0"
  currentUserCount = "0"
  currentNodeCount = "0"
  currentMasterIdentifier = "UNKNOWN"

  private subscriptions: Array<Subscription> = []

  constructor(
    private browserHistoryService: BrowserHistoryService,
    private webSocketService: WebsocketService,
    private geoLocationService: GeoLocationService,
    private sessionService: SessionService) {

  }

  public ngOnInit() {
    // send access record after then get all access records
    if (initialized) {
      this.initialize()
    } else {
      initialized = true
      this.browserHistoryService.addOne().subscribe(_ => {
        this.initialize()
      })
    }
  }

  ngOnDestroy(): void {
    this.subscriptions.forEach(sub => {
      sub.unsubscribe()
    })
  }

  initialize() {
    this.findAllBrowserHistory()
    this.subscriptions.push(this.subscribeWebsocketConnectionCount())
    this.subscriptions.push(this.subscribeBrowserHistoryCount())
    this.subscriptions.push(this.subscribeGatewayLeaderName())
    this.subscriptions.push(this.subscribeGatewayNodeCount())
    this.subscriptions.push(this.subscribeWebsocketConnected())
  }

  findAllBrowserHistory() {
    /** ngx-datatable doesn't display loading-bar when it has no item */
    this.rows = [ {} ]
    this.isTableLoading = true

    this.browserHistoryService.findAll(this.itemCountPerPage, this.currentPageOffset)
      .subscribe(response => {
        this.rows = response.rows

        let page = response.pagination
        this.currentPageOffset = page.currentPageOffset
        this.itemCountPerPage = page.itemCountPerPage
        this.totalItemCount = page.totalItemCount

        this.isTableLoading = false
      })
  }

  subscribeWebsocketConnectionCount(): Subscription {
    const eventType = WebSocketResponseHeader.ResponseTypeEnum.UpdateWebSocketConnectionCount
    return this.webSocketService.watch(eventType)
      .subscribe((response: WebSocketRealtimeResponse) => {
        this.currentUserCount = response.body.value
      })
  }

  subscribeBrowserHistoryCount(): Subscription {
    const eventType = WebSocketResponseHeader.ResponseTypeEnum.UpdateBrowserHistoryCount
    return this.webSocketService.watch(eventType)
      .subscribe((response: WebSocketRealtimeResponse) => {
        this.currentTotalAccessCount = response.body.value
      })
  }

  subscribeGatewayLeaderName(): Subscription {
    const eventType = WebSocketResponseHeader.ResponseTypeEnum.UpdateGatewayLeaderNodeName
    return this.webSocketService.watch(eventType)
      .subscribe((response: WebSocketRealtimeResponse) => {
        this.currentMasterIdentifier = "gateway-" + response.body.value
      })
  }

  subscribeGatewayNodeCount(): Subscription {
    const eventType = WebSocketResponseHeader.ResponseTypeEnum.UpdateGatewayNodeCount
    return this.webSocketService.watch(eventType)
      .subscribe((response: WebSocketRealtimeResponse) => {
        this.currentNodeCount = response.body.value
      })
  }

  subscribeWebsocketConnected(): Subscription {
    return this.webSocketService.watchWebsocketConnected()
      .subscribe(connected => {
        this.websocketConnected = connected
      })
  }

  /**
   * @param event { count, pageSize, limit, offset }
   */
  onPageChange(event) {
    this.currentPageOffset = event.offset
    this.findAllBrowserHistory()
  }
}
