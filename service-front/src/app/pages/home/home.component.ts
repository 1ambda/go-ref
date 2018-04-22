import { Component, OnDestroy, OnInit } from '@angular/core'

import { WebsocketService } from "../../shared/websocket.service"
import { WebSocketRealtimeResponse, WebSocketResponseHeader } from "../../generated/swagger/websocket"
import { Subscription } from 'rxjs/Subscription'
import { GeoLocationService } from "../../shared/geo-location.service"
import { SessionService } from "../../shared/session.service"
import { BrowserHistoryService } from "../../shared/browser-history.service"
import { NotificationService } from 'app/shared/notification.service'
import { FormControl, FormGroup } from '@angular/forms'

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
    { name: 'ID', prop: 'recordId' },
    { name: 'Session', prop: 'sessionId', width: 320, },
    { name: 'Browser Name', prop: 'browserName' },
    { name: 'Browser Version', prop: 'browserVersion' },
    { name: 'OS Name', prop: 'osName' },
    { name: 'OS Version', prop: 'osVersion' },
    { name: 'Mobile', prop: 'isMobile' },
    { name: 'Language', prop: 'language' },
    { name: 'Client Timestamp', prop: 'clientTimestamp', width: 300, },
    { name: 'Client Timezone', prop: 'clientTimezone' },
    { name: 'User Agent', prop: 'userAgent', width: 1000, },
  ]

  /**
   * pagination related variables
   */
  itemCountPerPage = 10
  totalItemCount = 0
  currentPageOffset = 0 // `page_number - 1`
  isTableLoading = true

  /**
   * filter related variables
   */
  defaultFilterColumn = 'SessionID'
  filterColumn = this.defaultFilterColumn
  filterValue = ''

  /**
   * server-side streamed variables
   */
  currentTotalAccessCount = '0'
  currentUserCount = '0'
  currentNodeCount = '0'
  currentMasterIdentifier = 'UNKNOWN'

  private subscriptions: Array<Subscription> = []

  constructor(
    private browserHistoryService: BrowserHistoryService,
    private webSocketService: WebsocketService,
    private geoLocationService: GeoLocationService,
    private sessionService: SessionService,
    private notificationService: NotificationService) {
  }

  public ngOnInit() {
    // start initialization process after browser history is sent
    this.subscriptions.push(this.subscribeBrowserHistorySendEvent())
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

    let filterColumn = this.filterColumn
    const filterValue = this.filterValue
    if (!filterValue || filterValue === "") {
      filterColumn = ""
    }

    this.browserHistoryService.findAll(this.itemCountPerPage, this.currentPageOffset,
      filterColumn, filterValue)
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

  onFilterValueChange(filterValue) {
    this.filterValue = filterValue.trim()
  }

  onFilterColumnChange(filterColumn) {
    this.filterColumn = filterColumn.trim()
  }

  onFilterSearchClick() {
    this.currentPageOffset = 0
    this.findAllBrowserHistory()
  }

  private subscribeBrowserHistorySendEvent(): Subscription {
    return this.browserHistoryService.watchBrowserHistorySendEvent()
      .subscribe(_ => {
        this.initialize()
      })
  }
}
