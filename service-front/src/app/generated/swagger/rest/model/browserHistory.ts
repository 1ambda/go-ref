/**
 * GatewayRestApi
 * REST API Spec for Gateway
 *
 * OpenAPI spec version: 0.0.1
 * 
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 * Do not edit the class manually.
 */


export interface BrowserHistory {
    browserName: string;
    browserVersion: string;
    osName: string;
    osVersion: string;
    isMobile: boolean;
    clientTimestamp: string;
    clientTimezone: string;
    language: string;
    userAgent: string;
    recordId?: number;
    sessionId?: string;
}
