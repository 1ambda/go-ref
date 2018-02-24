import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { environment } from 'environments/environment';

@Component({
  selector: 'app',
  encapsulation: ViewEncapsulation.None,
  styleUrls: [
    './app.component.css'
  ],
  template: `
    <app-navbar [class.mat-elevation-z6]="true"></app-navbar>
    <router-outlet></router-outlet>
  `
})
export class AppComponent implements OnInit {
  public showDevModule: boolean = environment.showDevModule;

  constructor() {}

  public ngOnInit() {
  }

}
