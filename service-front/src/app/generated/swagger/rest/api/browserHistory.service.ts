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
/* tslint:disable:no-unused-variable member-ordering */

import { Inject, Injectable, Optional }                      from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams,
         HttpResponse, HttpEvent }                           from '@angular/common/http';
import { CustomHttpUrlEncodingCodec }                        from '../encoder';

import { Observable }                                        from 'rxjs/Observable';

import { BrowserHistory } from '../model/browserHistory';
import { BrowserHistoryWithPagination } from '../model/browserHistoryWithPagination';
import { RestError } from '../model/restError';

import { BASE_PATH, COLLECTION_FORMATS }                     from '../variables';
import { Configuration }                                     from '../configuration';


@Injectable()
export class BrowserHistoryService {

    protected basePath = 'http://localhost/api';
    public defaultHeaders = new HttpHeaders();
    public configuration = new Configuration();

    constructor(protected httpClient: HttpClient, @Optional()@Inject(BASE_PATH) basePath: string, @Optional() configuration: Configuration) {
        if (basePath) {
            this.basePath = basePath;
        }
        if (configuration) {
            this.configuration = configuration;
            this.basePath = basePath || configuration.basePath || this.basePath;
        }
    }

    /**
     * @param consumes string[] mime-types
     * @return true: consumes contains 'multipart/form-data', false: otherwise
     */
    private canConsumeForm(consumes: string[]): boolean {
        const form = 'multipart/form-data';
        for (let consume of consumes) {
            if (form === consume) {
                return true;
            }
        }
        return false;
    }


    /**
     * 
     * 
     * @param body 
     * @param observe set whether or not to return the data Observable as the body, response or events. defaults to returning the body.
     * @param reportProgress flag to report request and response progress.
     */
    public addOne(body?: BrowserHistory, observe?: 'body', reportProgress?: boolean): Observable<BrowserHistory>;
    public addOne(body?: BrowserHistory, observe?: 'response', reportProgress?: boolean): Observable<HttpResponse<BrowserHistory>>;
    public addOne(body?: BrowserHistory, observe?: 'events', reportProgress?: boolean): Observable<HttpEvent<BrowserHistory>>;
    public addOne(body?: BrowserHistory, observe: any = 'body', reportProgress: boolean = false ): Observable<any> {

        let headers = this.defaultHeaders;

        // to determine the Accept header
        let httpHeaderAccepts: string[] = [
            'application/json'
        ];
        let httpHeaderAcceptSelected: string | undefined = this.configuration.selectHeaderAccept(httpHeaderAccepts);
        if (httpHeaderAcceptSelected != undefined) {
            headers = headers.set("Accept", httpHeaderAcceptSelected);
        }

        // to determine the Content-Type header
        let consumes: string[] = [
            'application/json'
        ];
        let httpContentTypeSelected:string | undefined = this.configuration.selectHeaderContentType(consumes);
        if (httpContentTypeSelected != undefined) {
            headers = headers.set("Content-Type", httpContentTypeSelected);
        }

        return this.httpClient.post<BrowserHistory>(`${this.basePath}/browser_history`,
            body,
            {
                withCredentials: this.configuration.withCredentials,
                headers: headers,
                observe: observe,
                reportProgress: reportProgress
            }
        );
    }

    /**
     * 
     * 
     * @param itemCountPerPage 
     * @param currentPageOffset 
     * @param filterColumn a column name which will be used for filtering &#x60;BrowserHistoryFilterType&#x60; definition 
     * @param filterValue a column value which will be used for filtering
     * @param sortBy a column name which will be used for sorting
     * @param orderBy &#39;asc&#39; or &#39;desc&#39;
     * @param observe set whether or not to return the data Observable as the body, response or events. defaults to returning the body.
     * @param reportProgress flag to report request and response progress.
     */
    public findAll(itemCountPerPage?: number, currentPageOffset?: number, filterColumn?: string, filterValue?: string, sortBy?: string, orderBy?: string, observe?: 'body', reportProgress?: boolean): Observable<BrowserHistoryWithPagination>;
    public findAll(itemCountPerPage?: number, currentPageOffset?: number, filterColumn?: string, filterValue?: string, sortBy?: string, orderBy?: string, observe?: 'response', reportProgress?: boolean): Observable<HttpResponse<BrowserHistoryWithPagination>>;
    public findAll(itemCountPerPage?: number, currentPageOffset?: number, filterColumn?: string, filterValue?: string, sortBy?: string, orderBy?: string, observe?: 'events', reportProgress?: boolean): Observable<HttpEvent<BrowserHistoryWithPagination>>;
    public findAll(itemCountPerPage?: number, currentPageOffset?: number, filterColumn?: string, filterValue?: string, sortBy?: string, orderBy?: string, observe: any = 'body', reportProgress: boolean = false ): Observable<any> {

        let queryParameters = new HttpParams({encoder: new CustomHttpUrlEncodingCodec()});
        if (itemCountPerPage !== undefined) {
            queryParameters = queryParameters.set('itemCountPerPage', <any>itemCountPerPage);
        }
        if (currentPageOffset !== undefined) {
            queryParameters = queryParameters.set('currentPageOffset', <any>currentPageOffset);
        }
        if (filterColumn !== undefined) {
            queryParameters = queryParameters.set('filterColumn', <any>filterColumn);
        }
        if (filterValue !== undefined) {
            queryParameters = queryParameters.set('filterValue', <any>filterValue);
        }
        if (sortBy !== undefined) {
            queryParameters = queryParameters.set('sortBy', <any>sortBy);
        }
        if (orderBy !== undefined) {
            queryParameters = queryParameters.set('orderBy', <any>orderBy);
        }

        let headers = this.defaultHeaders;

        // to determine the Accept header
        let httpHeaderAccepts: string[] = [
            'application/json'
        ];
        let httpHeaderAcceptSelected: string | undefined = this.configuration.selectHeaderAccept(httpHeaderAccepts);
        if (httpHeaderAcceptSelected != undefined) {
            headers = headers.set("Accept", httpHeaderAcceptSelected);
        }

        // to determine the Content-Type header
        let consumes: string[] = [
            'application/json'
        ];

        return this.httpClient.get<BrowserHistoryWithPagination>(`${this.basePath}/browser_history`,
            {
                params: queryParameters,
                withCredentials: this.configuration.withCredentials,
                headers: headers,
                observe: observe,
                reportProgress: reportProgress
            }
        );
    }

    /**
     * 
     * 
     * @param id 
     * @param body 
     * @param observe set whether or not to return the data Observable as the body, response or events. defaults to returning the body.
     * @param reportProgress flag to report request and response progress.
     */
    public findOne(id: number, body?: BrowserHistory, observe?: 'body', reportProgress?: boolean): Observable<BrowserHistory>;
    public findOne(id: number, body?: BrowserHistory, observe?: 'response', reportProgress?: boolean): Observable<HttpResponse<BrowserHistory>>;
    public findOne(id: number, body?: BrowserHistory, observe?: 'events', reportProgress?: boolean): Observable<HttpEvent<BrowserHistory>>;
    public findOne(id: number, body?: BrowserHistory, observe: any = 'body', reportProgress: boolean = false ): Observable<any> {
        if (id === null || id === undefined) {
            throw new Error('Required parameter id was null or undefined when calling findOne.');
        }

        let headers = this.defaultHeaders;

        // to determine the Accept header
        let httpHeaderAccepts: string[] = [
            'application/json'
        ];
        let httpHeaderAcceptSelected: string | undefined = this.configuration.selectHeaderAccept(httpHeaderAccepts);
        if (httpHeaderAcceptSelected != undefined) {
            headers = headers.set("Accept", httpHeaderAcceptSelected);
        }

        // to determine the Content-Type header
        let consumes: string[] = [
            'application/json'
        ];
        let httpContentTypeSelected:string | undefined = this.configuration.selectHeaderContentType(consumes);
        if (httpContentTypeSelected != undefined) {
            headers = headers.set("Content-Type", httpContentTypeSelected);
        }

        return this.httpClient.get<BrowserHistory>(`${this.basePath}/browser_history/${encodeURIComponent(String(id))}`,
            {
                withCredentials: this.configuration.withCredentials,
                headers: headers,
                observe: observe,
                reportProgress: reportProgress
            }
        );
    }

    /**
     * 
     * 
     * @param id 
     * @param observe set whether or not to return the data Observable as the body, response or events. defaults to returning the body.
     * @param reportProgress flag to report request and response progress.
     */
    public removeOne(id: number, observe?: 'body', reportProgress?: boolean): Observable<any>;
    public removeOne(id: number, observe?: 'response', reportProgress?: boolean): Observable<HttpResponse<any>>;
    public removeOne(id: number, observe?: 'events', reportProgress?: boolean): Observable<HttpEvent<any>>;
    public removeOne(id: number, observe: any = 'body', reportProgress: boolean = false ): Observable<any> {
        if (id === null || id === undefined) {
            throw new Error('Required parameter id was null or undefined when calling removeOne.');
        }

        let headers = this.defaultHeaders;

        // to determine the Accept header
        let httpHeaderAccepts: string[] = [
            'application/json'
        ];
        let httpHeaderAcceptSelected: string | undefined = this.configuration.selectHeaderAccept(httpHeaderAccepts);
        if (httpHeaderAcceptSelected != undefined) {
            headers = headers.set("Accept", httpHeaderAcceptSelected);
        }

        // to determine the Content-Type header
        let consumes: string[] = [
            'application/json'
        ];

        return this.httpClient.delete<any>(`${this.basePath}/browser_history/${encodeURIComponent(String(id))}`,
            {
                withCredentials: this.configuration.withCredentials,
                headers: headers,
                observe: observe,
                reportProgress: reportProgress
            }
        );
    }

}
