import { Component, OnInit } from "@angular/core"
import { GeoLocationService } from "../../shared/geo-location.service"
import { Subscription } from 'rxjs/Subscription'
import { Geolocation } from "../../generated/swagger/rest"


@Component({
  selector: 'location',
  providers: [],
  styleUrls: [ './location.component.css' ],
  templateUrl: './location.component.html'
})
export class LocationComponent implements OnInit {

  private title: string = 'My first AGM project'
  private lat: number = 51.678418
  private lng: number = 7.809007

  private subscriptions: Array<Subscription> = []

  constructor(private geoLocationService: GeoLocationService) {
  }

  ngOnInit(): void {
    this.subscriptions.push(this.subscribeGeolocation())
  }

  ngOnDestroy(): void {
    this.subscriptions.forEach(sub => {
      sub.unsubscribe()
    })
  }

  private subscribeGeolocation(): Subscription {
    return this.geoLocationService.watchGeolocation()
      .subscribe((geolocation: Geolocation) => {
        this.lat = geolocation.latitude
        this.lng = geolocation.longitude
      })
  }

}
