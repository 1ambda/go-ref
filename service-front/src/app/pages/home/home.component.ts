import {
  Component,
  OnInit
} from '@angular/core'

import 'clientjs'
import * as moment from 'moment'

import { AccessService, Access } from '../../generated/swagger';

@Component({
  selector: 'home',  // <home></home>
  providers: [],
  styleUrls: [ './home.component.css' ],
  templateUrl: './home.component.html'
})
export class HomeComponent implements OnInit {
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
  ];

  /**
   * pagination related variables
   */
  itemCountPerPage = 10
  totalItemCount = 0
  currentPageOffset = 0 /** page_number -1 */
  isTableLoading = true

  constructor(
    private accessService: AccessService
  ) {}

  public ngOnInit() {
    // send access record after then get all access records

    this.sendAccess()
      .subscribe(created => {
        this.fetchAllAccessList()
      })
  }

  sendAccess() {
    const client = new ClientJS();

    const access: Access = {
      browserName: client.getBrowser(),
      browserVersion: client.getBrowserVersion(),
      osVersion: client.getOS(),
      osName: client.getOSVersion(),
      isMobile: client.isMobile() ? 'Y' : 'N' ,
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


  /**
   * @param event { count, pageSize, limit, offset }
   */
  onPageChange(event){
    console.log(event)
    this.currentPageOffset = event.offset;
    this.fetchAllAccessList()
  }
}
