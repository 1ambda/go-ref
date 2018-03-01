// package: grpc
// file: gateway.proto

import * as gateway_pb from "./gateway_pb";
export class Gateway {
  static serviceName = "grpc.Gateway";
}
export namespace Gateway {
  export class SubscribeTotalAccessCount {
    static readonly methodName = "SubscribeTotalAccessCount";
    static readonly service = Gateway;
    static readonly requestStream = false;
    static readonly responseStream = true;
    static readonly requestType = gateway_pb.EmptyRequest;
    static readonly responseType = gateway_pb.CountResponse;
  }
  export class SubscribeCurrentUserCount {
    static readonly methodName = "SubscribeCurrentUserCount";
    static readonly service = Gateway;
    static readonly requestStream = false;
    static readonly responseStream = true;
    static readonly requestType = gateway_pb.EmptyRequest;
    static readonly responseType = gateway_pb.CountResponse;
  }
}
