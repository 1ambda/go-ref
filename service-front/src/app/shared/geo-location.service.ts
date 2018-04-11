import { Injectable } from '@angular/core'

const geolocator = require('geolocator')

geolocator.config({
  language: "en",
  google: {
    version: "3",
    key: KEY_GOOGLE_MAP_API,
  }
})

@Injectable()
export class GeoLocationService {

  private ip: string
  private latitude: number
  private longitude: number
  private timezone: string
  private formattedAddress: string
  private address: object

  constructor() {
    const options = {
      enableHighAccuracy: true,
      timeout: 5000,
      maximumWait: 10000,     // max wait time for desired accuracy
      maximumAge: 0,          // disable cache
      desiredAccuracy: 30,    // meters
      fallbackToIP: true,     // fallback to IP if Geolocation fails or rejected
      addressLookup: true,    // requires Google API key if true
      timezone: true,         // requires Google API key if true
    }

    console.log("Fetching geolocation information from google")
    geolocator.locate(options, (err, response) => {
      if (err) {
        console.log(err)
        return
      }

      console.log("Fetched geolocation information from google")

      this.ip = response.ip

      if (response.coords) {
        this.latitude = response.coords.latitude
        this.longitude = response.coords.longitude
      }

      if (response.timezone) {
        this.timezone = response.timezone.id
      }

      this.formattedAddress = response.formattedAddress
      this.address = response.address
    })
  }
}

