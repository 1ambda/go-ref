/**
 * GatewayWebsocketApi
 * Websocket API Spec for Gateway
 *
 * OpenAPI spec version: 0.0.1
 * 
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 * Do not edit the class manually.
 */


export interface WebSocketResponseHeader {
    responseType: WebSocketResponseHeader.ResponseTypeEnum;
    error?: Error;
}
export namespace WebSocketResponseHeader {
    export type ResponseTypeEnum = 'Error' | 'UpdateBrowserHistoryCount' | 'UpdateWebSocketConnectionCount' | 'UpdateGatewayLeaderNodeName' | 'UpdateGatewayNodeCount';
    export const ResponseTypeEnum = {
        Error: 'Error' as ResponseTypeEnum,
        UpdateBrowserHistoryCount: 'UpdateBrowserHistoryCount' as ResponseTypeEnum,
        UpdateWebSocketConnectionCount: 'UpdateWebSocketConnectionCount' as ResponseTypeEnum,
        UpdateGatewayLeaderNodeName: 'UpdateGatewayLeaderNodeName' as ResponseTypeEnum,
        UpdateGatewayNodeCount: 'UpdateGatewayNodeCount' as ResponseTypeEnum
    }
}
