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
        <p>Git Commit: <strong>{{projectGitCommit}}</strong></p>
        <p>Git Branch: <strong>{{projectGitBranch}}</strong></p>
        <p>Project Version : <strong>{{projectVersion}}</strong></p>
        <p>Project Build Date: <strong>{{projectBuildDate}}</strong></p>
    </div>
  `
})
export class AboutComponent implements OnInit {
  constructor() {}

  projectGitCommit = PROJECT_GIT_COMMIT
  projectGitBranch = PROJECT_GIT_BRANCH
  projectVersion = PROJECT_VERSION
  projectBuildDate = PROJECT_BUILD_DATE


  public ngOnInit() {
  }
}
