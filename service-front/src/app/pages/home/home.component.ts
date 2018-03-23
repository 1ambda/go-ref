import {Component, OnInit} from '@angular/core'

import 'clientjs'
import * as moment from 'moment'

import {Access, AccessService} from '../../generated/swagger/rest'
import {WebsocketService} from "../../shared"
import {WebSocketRealtimeResponse, WebSocketResponseHeader} from "../../generated/swagger/websocket";


@Component({
  selector: 'home', // <home></home>
  providers: [],
  styleUrls: ['./home.component.css'],
  templateUrl: './home.component.html'
})
export class HomeComponent implements OnInit {
  rows = []
  columns = [
    {name: 'id', prop: 'id'},
    {name: 'browser_name', prop: 'browserName'},
    {name: 'browser_version', prop: 'browserVersion'},
    {name: 'os_name', prop: 'osName'},
    {name: 'os_version', prop: 'osVersion'},
    {name: 'is_mobile', prop: 'isMobile'},
    {name: 'language', prop: 'language'},
    {name: 'timestamp', prop: 'timestamp'},
    {name: 'timezone', prop: 'timezone'},
    {name: 'uuid', prop: 'uuid'},
    {name: 'user_agent', prop: 'userAgent'},
  ];

  /**
   * pagination related variables
   */
  itemCountPerPage = 10
  totalItemCount = 0
  currentPageOffset = 0
  /** page_number -1 */
  isTableLoading = true

  /**
   * server-side streamed variables
   */
  totalAccessCount = "0"
  currentUserCount = "0"
  currentNodeCount = "0"
  currentMasterIdentifier = "UNKNOWN"

  constructor(
    private accessService: AccessService,
    private webSocketService: WebsocketService) {
  }

  public ngOnInit() {
    // send access record after then get all access records

    this.sendAccess()
      .subscribe(created => {
        this.fetchAllAccessList()

        // TODO: re-subscribe when websocket is disconnected
        this.subscribeCurrentConnectionCount()
        this.subscribeTotalAccessCount()
        this.subscribeLMasterIdentifier()
        this.subscribeLMasterNodeCount()
      })

  }

  sendAccess() {
    const client = new ClientJS();

    const access: Access = {
      browserName: client.getBrowser(),
      browserVersion: client.getBrowserVersion(),
      osVersion: client.getOS(),
      osName: client.getOSVersion(),
      isMobile: client.isMobile() ? 'Y' : 'N',
      timezone: client.getTimeZone(),
      timestamp: moment.utc().unix().toString(),
      language: client.getLanguage(),
      userAgent: window.navigator.userAgent,
    };


    return this.accessService.addOne(access)
  }

  fetchAllAccessList() {
    /** ngx-datatable doesn't display loading-bar when it has no item */
    this.rows = [{}]
    this.isTableLoading = true

    this.accessService.findAll(this.itemCountPerPage, this.currentPageOffset)
      .subscribe(response => {
        this.rows = response.rows

        let page = response.pagination
        this.currentPageOffset = page.currentPageOffset
        this.itemCountPerPage = page.itemCountPerPage
        this.totalItemCount = page.totalItemCount

        this.isTableLoading = false
      })
  }

  subscribeCurrentConnectionCount() {
    const eventType = WebSocketResponseHeader.ResponseTypeEnum.UpdateConnectionCount
    this.webSocketService.watch(eventType)
      .subscribe((response: WebSocketRealtimeResponse) => {
        this.currentUserCount = response.body.value
      })
  }

  subscribeTotalAccessCount() {
    const eventType = WebSocketResponseHeader.ResponseTypeEnum.UpdateTotalAccessCount
    this.webSocketService.watch(eventType)
      .subscribe((response: WebSocketRealtimeResponse) => {
        this.totalAccessCount = response.body.value
      })
  }

  subscribeLMasterIdentifier() {
    const eventType = WebSocketResponseHeader.ResponseTypeEnum.UpdateMasterIdentifier
    this.webSocketService.watch(eventType)
      .subscribe((response: WebSocketRealtimeResponse) => {
        this.currentMasterIdentifier = response.body.value
      })
  }

  subscribeLMasterNodeCount() {
    const eventType = WebSocketResponseHeader.ResponseTypeEnum.UpdateMasterNodeCount
    this.webSocketService.watch(eventType)
      .subscribe((response: WebSocketRealtimeResponse) => {
        this.currentNodeCount = response.body.value
      })
  }

  /**
   * @param event { count, pageSize, limit, offset }
   */
  onPageChange(event) {
    this.currentPageOffset = event.offset;
    this.fetchAllAccessList()
  }
}
