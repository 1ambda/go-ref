import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http'
import { RouterModule, PreloadAllModules } from '@angular/router'
import { BrowserAnimationsModule } from '@angular/platform-browser/animations'

import { environment } from 'environments/environment'
import { ROUTES } from './app.routes'
import { AppComponent } from './app.component'
import { APP_RESOLVER_PROVIDERS } from './app.resolver'

import { WebsocketService } from "./shared"

import { HomeComponent } from './pages/home'
import { AboutComponent } from './pages/about'
import { NoContentComponent } from './pages/no-content'
import { DevModuleModule } from './pages/+dev-module'

/**
 * rxjs global import
 */

import 'rxjs/add/operator/filter'

/**
 * swagger generated clients
 *
 * See: https://github.com/swagger-api/swagger-codegen/tree/master/samples/client/petstore/typescript-angular-v4/npm
 */
import { ApiModule, Configuration, ConfigurationParameters } from './generated/swagger'

/**
 * angular material
 *
 * See: https://material.angular.io/guide/getting-started
 */
import 'hammerjs';

/**
 * angular flex layout
 *
 * See: https://github.com/angular/flex-layout/wiki/Webpack-Configuration
 */
import { FlexLayoutModule } from '@angular/flex-layout';

/**
 * ngx-datatable
 */
import { NgxDatatableModule } from '@swimlane/ngx-datatable';


/**
 * swagger client
 */
export function apiConfigFactory (): Configuration {
  const params: ConfigurationParameters = {
    // set configuration parameters here.
  }

  if (!environment.production) {
    params.basePath = "http://localhost:50002/api"
  }

  // TODO: dev, prod

  return new Configuration(params);
}


import '../styles/styles.scss';

// Application wide providers
const APP_PROVIDERS = [
  ...APP_RESOLVER_PROVIDERS,
];

/**
 * App Modules
 */
import { NavbarModule } from './pages/common/navbar'
import { MatCardModule } from '@angular/material'

/**
 * `AppModule` is the main entry point.
 */
@NgModule({
  bootstrap: [ AppComponent ],
  declarations: [
    AppComponent,
    AboutComponent,
    HomeComponent,
    NoContentComponent,
  ],
  /**
   * Import Angular's modules.
   */
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    FlexLayoutModule,
    FormsModule,
    NgxDatatableModule,
    ApiModule.forRoot(apiConfigFactory),
    HttpClientModule,
    RouterModule.forRoot(ROUTES, {
      useHash: Boolean(history.pushState) === false,
      preloadingStrategy: PreloadAllModules
    }),

    MatCardModule,
    NavbarModule,

    /**
     * This section will import the `DevModuleModule` only in certain build types.
     * When the module is not imported it will get tree shaked.
     * This is a simple example, a big app should probably implement some logic
     */
    ...environment.showDevModule ? [ DevModuleModule ] : [],
  ],
  /**
   * Expose our Services and Providers into Angular's dependency injection.
   */
  providers: [
    WebsocketService,
    environment.ENV_PROVIDERS,
    APP_PROVIDERS
  ]
})
export class AppModule {}
