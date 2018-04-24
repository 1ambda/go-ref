import { Injectable } from '@angular/core'
import { HttpEvent, HttpEventType, HttpHandler, HttpInterceptor, HttpRequest, HttpResponse, HttpErrorResponse } from '@angular/common/http'
import { Observable } from 'rxjs/Observable'
import { NotificationService } from "./notification.service"


@Injectable()
export class HttpResponseInterceptor implements HttpInterceptor {
  constructor(private notificationService: NotificationService) {
  }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    const started = Date.now()

    return next.handle(req).do((response: any) => {
      if (response instanceof HttpResponse) {
        const elapsed = Date.now() - started
        const message = `[${req.method}] (${elapsed} ms) ${req.url} succeeded`
        console.log(message)
      }

      return response
    }, (response: any) => {
      if (response instanceof HttpErrorResponse) {
        const elapsed = Date.now() - started
        const message = `[${req.method}] (${elapsed} ms) ${req.urlWithParams} failed`
        console.error(message)

        if (response.error && response.error.type && response.error.code) {
          // server-formatted error response
          this.notificationService.displayFormattedError(response.error.code,
            response.error.message, response.error.type)
        } else {
          // unknown error response format
          this.notificationService.displayError("Unknown REST Error", response.message)
        }

      }

      return response
    })
  }
}
