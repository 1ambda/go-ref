import { Injectable } from '@angular/core'
import { SessionService } from "./session.service"
import { Geolocation, GeolocationService as GeolocationApiService, SessionResponse } from "../generated/swagger/rest"
import { NotificationService } from "./notification.service"
import { ReplaySubject } from 'rxjs/ReplaySubject'
import { Observable } from 'rxjs/Observable'

const geolocator = require('geolocator')

const GEOLOCATOR_CONFIG = {
  language: "en",
  google: {
    version: "3",
    key: KEY_GOOGLE_MAP_API,
  }
}

geolocator.config(GEOLOCATOR_CONFIG)

const GEOLOCATOR_OPTIONS = {
  enableHighAccuracy: true,
  timeout: 5000,
  maximumWait: 10000,     // max wait time for desired accuracy
  maximumAge: 0,          // disable cache
  desiredAccuracy: 30,    // meters
  fallbackToIP: true,     // fallback to IP if Geolocation fails or rejected
  addressLookup: true,    // requires Google API key if true
  timezone: true,         // requires Google API key if true
}

@Injectable()
export class GeoLocationService {

  private geolocation: Geolocation
  private geolocationReplay: ReplaySubject<Geolocation> = new ReplaySubject<Geolocation>()

  constructor(private sessionService: SessionService,
              private notificationService: NotificationService,
              private geolocationApiService: GeolocationApiService) {

    this.sessionService.subscribeSession().subscribe((session: SessionResponse) => {
      if (!session) {
        return
      }

      console.info("Initializing GeoLocationService")

      geolocator.locate(GEOLOCATOR_OPTIONS, (err, response) => {
        if (err) {
          console.log(err)
          this.notificationService.displayError("Geolocation", "Failed to fetch geolocation")
          return
        }

        const request: Geolocation = {
          apiProvider: "Google",
          apiLanguage: GEOLOCATOR_CONFIG.language,
          apiVersion: GEOLOCATOR_CONFIG.google.version,
          apiDesiredAccuracy: GEOLOCATOR_OPTIONS.desiredAccuracy,

          provider: response.provider,
          timezone: response.timezone ? response.timezone.id : null,
          ip: response.ip,
          googlePlaceID: response.placeId,

          latitude: response.coords ? response.coords.latitude : null,
          longitude: response.coords ? response.coords.longitude : null,

          formattedAddress: response.address ? response.address.formattedAddress : null,
          commonName: response.address ? response.address.commonName : null,
          streetNumber: response.address ? response.address.streetNumber : null,
          street: response.address ? response.address.street : null,

          route: response.address ? response.address.route : null,
          neighborhood: response.address ? response.address.neighborhood : null,
          town: response.address ? response.address.town : null,
          city: response.address ? response.address.city : null,
          region: response.address ? response.address.region : null,
          postalCode: response.address ? response.address.postalCode : null,

          state: response.address ? response.address.state : null,
          stateCode: response.address ? response.address.stateCode : null,
          country: response.address ? response.address.country : null,
          countryCode: response.address ? response.address.countryCode : null,
        }

        this.geolocationApiService.add(request).subscribe(response => {
          this.notificationService.displayInfo("Geolocation", "Analyzed")

          this.geolocation = response
          this.geolocationReplay.next(response)
        })
      })
    })
  }

  watchGeolocation(): Observable<Geolocation> {
    return this.geolocationReplay
  }
}

