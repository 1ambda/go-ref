// package: grpc
// file: gateway.proto

import * as jspb from "google-protobuf";

export class WebSocketResponseHeader extends jspb.Message {
  getEventtype(): WebSocketResponseHeader.WebSocketEventType;
  setEventtype(value: WebSocketResponseHeader.WebSocketEventType): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WebSocketResponseHeader.AsObject;
  static toObject(includeInstance: boolean, msg: WebSocketResponseHeader): WebSocketResponseHeader.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WebSocketResponseHeader, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WebSocketResponseHeader;
  static deserializeBinaryFromReader(message: WebSocketResponseHeader, reader: jspb.BinaryReader): WebSocketResponseHeader;
}

export namespace WebSocketResponseHeader {
  export type AsObject = {
    eventtype: WebSocketResponseHeader.WebSocketEventType,
  }

  export enum WebSocketEventType {
    UPDATE_TOTAL_ACCESS_COUNT = 0,
    UPDATE_CURRENT_CONNECTION_COUNT = 1,
    UPDATE_CURRENT_NODE_COUNT = 2,
    UPDATE_CURRENT_MASTER_IDENTIFIER = 3,
  }
}

export class WebSocketRealtimeResponseBody extends jspb.Message {
  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WebSocketRealtimeResponseBody.AsObject;
  static toObject(includeInstance: boolean, msg: WebSocketRealtimeResponseBody): WebSocketRealtimeResponseBody.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WebSocketRealtimeResponseBody, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WebSocketRealtimeResponseBody;
  static deserializeBinaryFromReader(message: WebSocketRealtimeResponseBody, reader: jspb.BinaryReader): WebSocketRealtimeResponseBody;
}

export namespace WebSocketRealtimeResponseBody {
  export type AsObject = {
    value: string,
  }
}

export class WebSocketRealtimeResponse extends jspb.Message {
  hasHeader(): boolean;
  clearHeader(): void;
  getHeader(): WebSocketResponseHeader | undefined;
  setHeader(value?: WebSocketResponseHeader): void;

  hasBody(): boolean;
  clearBody(): void;
  getBody(): WebSocketRealtimeResponseBody | undefined;
  setBody(value?: WebSocketRealtimeResponseBody): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WebSocketRealtimeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: WebSocketRealtimeResponse): WebSocketRealtimeResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WebSocketRealtimeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WebSocketRealtimeResponse;
  static deserializeBinaryFromReader(message: WebSocketRealtimeResponse, reader: jspb.BinaryReader): WebSocketRealtimeResponse;
}

export namespace WebSocketRealtimeResponse {
  export type AsObject = {
    header?: WebSocketResponseHeader.AsObject,
    body?: WebSocketRealtimeResponseBody.AsObject,
  }
}

