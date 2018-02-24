import {
  Component,
  OnInit
} from '@angular/core'

@Component({
  selector: 'about',
  styles: [`
  `],
  template: `
    <h1>About</h1>
    <div>
        Hello :)
    </div>
  `
})
export class AboutComponent implements OnInit {
  constructor() {}

  public ngOnInit() {
  }
}
