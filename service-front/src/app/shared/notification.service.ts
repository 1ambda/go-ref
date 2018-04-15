import { Injectable } from '@angular/core'
import { NotificationsService as NotificationLibrary } from 'angular2-notifications'

// https://github.com/flauc/angular2-notifications
const NOTIFICATION_OPTION = {
  timeOut: 5000,
  showProgressBar: true,
  pauseOnHover: true,
  clickToClose: true,
  animate: 'fromRight'
}

@Injectable()
export class NotificationService {
  constructor(private notificationLibrary: NotificationLibrary) {
  }

  displayError(title: string, message: string) {
    this.notificationLibrary.error(title, message, NOTIFICATION_OPTION)
  }

  displayWarn(title: string, message: string) {
    this.notificationLibrary.warn(title, message, NOTIFICATION_OPTION)
  }
}
