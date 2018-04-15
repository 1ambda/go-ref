import { Component, OnInit, ViewEncapsulation } from '@angular/core'
import { environment } from 'environments/environment'
import { WebsocketService } from 'app/shared/websocket.service'
import { Subscription } from 'rxjs/Subscription'

@Component({
  selector: 'app',
  encapsulation: ViewEncapsulation.None,
  styleUrls: [ './app.component.css' ],
  templateUrl: './app.component.html'
})
export class AppComponent implements OnInit {
  public showDevModule: boolean = environment.showDevModule

  /**
   * indicates websocket is connected or not
   */
  private websocketConnected = false
  private subscriptions: Array<Subscription> = []

  constructor(private webSocketService: WebsocketService) {
  }

  public ngOnInit() {
    this.subscriptions.push(
      this.webSocketService.watchWebsocketConnected()
        .subscribe(connected => {
          if (connected != this.websocketConnected) {
            this.websocketConnected = connected
          }
        })
    )
  }

  ngOnDestroy(): void {
    this.subscriptions.forEach(sub => {
      sub.unsubscribe()
    })
  }

}
