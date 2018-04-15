export * from './browserHistory.service';
import { BrowserHistoryService } from './browserHistory.service';
export * from './geolocation.service';
import { GeolocationService } from './geolocation.service';
export * from './session.service';
import { SessionService } from './session.service';
export const APIS = [BrowserHistoryService, GeolocationService, SessionService];
