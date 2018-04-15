// Code generated by go-swagger; DO NOT EDIT.

package rest_server

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "REST API Spec for Gateway",
    "title": "GatewayRestApi",
    "version": "0.0.1"
  },
  "basePath": "/api",
  "paths": {
    "/browser_history": {
      "get": {
        "tags": [
          "BrowserHistory"
        ],
        "operationId": "findAll",
        "parameters": [
          {
            "type": "integer",
            "format": "int64",
            "default": 10,
            "name": "itemCountPerPage",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int32",
            "default": 0,
            "name": "currentPageOffset",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "BrowserHistory records with pagination info",
            "schema": {
              "$ref": "#/definitions/BrowserHistoryWithPagination"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/RestError"
            }
          }
        }
      },
      "post": {
        "tags": [
          "BrowserHistory"
        ],
        "operationId": "addOne",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/BrowserHistory"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Created",
            "schema": {
              "$ref": "#/definitions/BrowserHistory"
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/RestError"
            }
          }
        }
      }
    },
    "/browser_history/{id}": {
      "get": {
        "tags": [
          "BrowserHistory"
        ],
        "operationId": "findOne",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/BrowserHistory"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/BrowserHistory"
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/RestError"
            }
          }
        }
      },
      "delete": {
        "tags": [
          "BrowserHistory"
        ],
        "operationId": "removeOne",
        "responses": {
          "204": {
            "description": "Deleted"
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/RestError"
            }
          }
        }
      },
      "parameters": [
        {
          "type": "integer",
          "format": "int64",
          "name": "id",
          "in": "path",
          "required": true
        }
      ]
    },
    "/geolocation": {
      "post": {
        "tags": [
          "geolocation"
        ],
        "operationId": "add",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Geolocation"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Created",
            "schema": {
              "$ref": "#/definitions/Geolocation"
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/RestError"
            }
          }
        }
      }
    },
    "/session": {
      "post": {
        "tags": [
          "session"
        ],
        "operationId": "validateOrGenerate",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/SessionRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/SessionResponse"
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/RestError"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "BrowserHistory": {
      "type": "object",
      "required": [
        "browserName",
        "browserVersion",
        "osName",
        "osVersion",
        "isMobile",
        "language",
        "clientTimezone",
        "clientTimestamp",
        "userAgent"
      ],
      "properties": {
        "browserName": {
          "type": "string"
        },
        "browserVersion": {
          "type": "string"
        },
        "clientTimestamp": {
          "type": "string"
        },
        "clientTimezone": {
          "type": "string"
        },
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "isMobile": {
          "type": "boolean"
        },
        "language": {
          "type": "string"
        },
        "osName": {
          "type": "string"
        },
        "osVersion": {
          "type": "string"
        },
        "userAgent": {
          "type": "string"
        },
        "uuid": {
          "type": "string"
        }
      }
    },
    "BrowserHistoryWithPagination": {
      "type": "object",
      "required": [
        "pagination",
        "rows"
      ],
      "properties": {
        "pagination": {
          "$ref": "#/definitions/Pagination"
        },
        "rows": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/BrowserHistory"
          }
        }
      }
    },
    "Geolocation": {
      "type": "object",
      "properties": {
        "city": {
          "type": "string"
        },
        "common_name": {
          "type": "string"
        },
        "country": {
          "type": "string"
        },
        "country_code": {
          "type": "string"
        },
        "formatted_address": {
          "type": "string"
        },
        "googlePlaceID": {
          "type": "string"
        },
        "ip": {
          "type": "string"
        },
        "latitude": {
          "type": "number",
          "format": "float"
        },
        "longitude": {
          "type": "number",
          "format": "float"
        },
        "neighborhood": {
          "type": "string"
        },
        "postal_code": {
          "type": "string"
        },
        "provider": {
          "type": "string"
        },
        "region": {
          "type": "string"
        },
        "route": {
          "type": "string"
        },
        "state": {
          "type": "string"
        },
        "state_code": {
          "type": "string"
        },
        "street": {
          "type": "string"
        },
        "street_number": {
          "type": "string"
        },
        "timezone": {
          "type": "string"
        },
        "town": {
          "type": "string"
        }
      }
    },
    "Pagination": {
      "type": "object",
      "required": [
        "itemCountPerPage",
        "currentPageOffset",
        "totalItemCount"
      ],
      "properties": {
        "currentPageOffset": {
          "type": "integer",
          "format": "int32"
        },
        "itemCountPerPage": {
          "type": "integer",
          "format": "int64"
        },
        "totalItemCount": {
          "type": "integer",
          "format": "int64"
        }
      }
    },
    "RestError": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer",
          "format": "int64"
        },
        "message": {
          "type": "string"
        },
        "timestamp": {
          "type": "string"
        },
        "type": {
          "type": "string",
          "enum": [
            "InvalidSession",
            "InternalServer",
            "RecordDoesNotxist"
          ]
        }
      }
    },
    "SessionRequest": {
      "type": "object",
      "required": [
        "sessionID"
      ],
      "properties": {
        "sessionID": {
          "type": "string"
        }
      }
    },
    "SessionResponse": {
      "type": "object",
      "required": [
        "sessionID",
        "createdAt",
        "updatedAt",
        "expiredAt",
        "refreshed",
        "refreshCount"
      ],
      "properties": {
        "createdAt": {
          "type": "integer",
          "format": "int64"
        },
        "expiredAt": {
          "type": "integer",
          "format": "int64"
        },
        "refreshCount": {
          "type": "integer",
          "format": "int64"
        },
        "refreshed": {
          "type": "boolean"
        },
        "sessionID": {
          "type": "string"
        },
        "updatedAt": {
          "type": "integer",
          "format": "int64"
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "REST API Spec for Gateway",
    "title": "GatewayRestApi",
    "version": "0.0.1"
  },
  "basePath": "/api",
  "paths": {
    "/browser_history": {
      "get": {
        "tags": [
          "BrowserHistory"
        ],
        "operationId": "findAll",
        "parameters": [
          {
            "type": "integer",
            "format": "int64",
            "default": 10,
            "name": "itemCountPerPage",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int32",
            "default": 0,
            "name": "currentPageOffset",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "BrowserHistory records with pagination info",
            "schema": {
              "$ref": "#/definitions/BrowserHistoryWithPagination"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/RestError"
            }
          }
        }
      },
      "post": {
        "tags": [
          "BrowserHistory"
        ],
        "operationId": "addOne",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/BrowserHistory"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Created",
            "schema": {
              "$ref": "#/definitions/BrowserHistory"
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/RestError"
            }
          }
        }
      }
    },
    "/browser_history/{id}": {
      "get": {
        "tags": [
          "BrowserHistory"
        ],
        "operationId": "findOne",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/BrowserHistory"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/BrowserHistory"
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/RestError"
            }
          }
        }
      },
      "delete": {
        "tags": [
          "BrowserHistory"
        ],
        "operationId": "removeOne",
        "responses": {
          "204": {
            "description": "Deleted"
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/RestError"
            }
          }
        }
      },
      "parameters": [
        {
          "type": "integer",
          "format": "int64",
          "name": "id",
          "in": "path",
          "required": true
        }
      ]
    },
    "/geolocation": {
      "post": {
        "tags": [
          "geolocation"
        ],
        "operationId": "add",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Geolocation"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Created",
            "schema": {
              "$ref": "#/definitions/Geolocation"
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/RestError"
            }
          }
        }
      }
    },
    "/session": {
      "post": {
        "tags": [
          "session"
        ],
        "operationId": "validateOrGenerate",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/SessionRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/SessionResponse"
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/RestError"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "BrowserHistory": {
      "type": "object",
      "required": [
        "browserName",
        "browserVersion",
        "osName",
        "osVersion",
        "isMobile",
        "language",
        "clientTimezone",
        "clientTimestamp",
        "userAgent"
      ],
      "properties": {
        "browserName": {
          "type": "string"
        },
        "browserVersion": {
          "type": "string"
        },
        "clientTimestamp": {
          "type": "string"
        },
        "clientTimezone": {
          "type": "string"
        },
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "isMobile": {
          "type": "boolean"
        },
        "language": {
          "type": "string"
        },
        "osName": {
          "type": "string"
        },
        "osVersion": {
          "type": "string"
        },
        "userAgent": {
          "type": "string"
        },
        "uuid": {
          "type": "string"
        }
      }
    },
    "BrowserHistoryWithPagination": {
      "type": "object",
      "required": [
        "pagination",
        "rows"
      ],
      "properties": {
        "pagination": {
          "$ref": "#/definitions/Pagination"
        },
        "rows": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/BrowserHistory"
          }
        }
      }
    },
    "Geolocation": {
      "type": "object",
      "properties": {
        "city": {
          "type": "string"
        },
        "common_name": {
          "type": "string"
        },
        "country": {
          "type": "string"
        },
        "country_code": {
          "type": "string"
        },
        "formatted_address": {
          "type": "string"
        },
        "googlePlaceID": {
          "type": "string"
        },
        "ip": {
          "type": "string"
        },
        "latitude": {
          "type": "number",
          "format": "float"
        },
        "longitude": {
          "type": "number",
          "format": "float"
        },
        "neighborhood": {
          "type": "string"
        },
        "postal_code": {
          "type": "string"
        },
        "provider": {
          "type": "string"
        },
        "region": {
          "type": "string"
        },
        "route": {
          "type": "string"
        },
        "state": {
          "type": "string"
        },
        "state_code": {
          "type": "string"
        },
        "street": {
          "type": "string"
        },
        "street_number": {
          "type": "string"
        },
        "timezone": {
          "type": "string"
        },
        "town": {
          "type": "string"
        }
      }
    },
    "Pagination": {
      "type": "object",
      "required": [
        "itemCountPerPage",
        "currentPageOffset",
        "totalItemCount"
      ],
      "properties": {
        "currentPageOffset": {
          "type": "integer",
          "format": "int32"
        },
        "itemCountPerPage": {
          "type": "integer",
          "format": "int64"
        },
        "totalItemCount": {
          "type": "integer",
          "format": "int64"
        }
      }
    },
    "RestError": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer",
          "format": "int64"
        },
        "message": {
          "type": "string"
        },
        "timestamp": {
          "type": "string"
        },
        "type": {
          "type": "string",
          "enum": [
            "InvalidSession",
            "InternalServer",
            "RecordDoesNotxist"
          ]
        }
      }
    },
    "SessionRequest": {
      "type": "object",
      "required": [
        "sessionID"
      ],
      "properties": {
        "sessionID": {
          "type": "string"
        }
      }
    },
    "SessionResponse": {
      "type": "object",
      "required": [
        "sessionID",
        "createdAt",
        "updatedAt",
        "expiredAt",
        "refreshed",
        "refreshCount"
      ],
      "properties": {
        "createdAt": {
          "type": "integer",
          "format": "int64"
        },
        "expiredAt": {
          "type": "integer",
          "format": "int64"
        },
        "refreshCount": {
          "type": "integer",
          "format": "int64"
        },
        "refreshed": {
          "type": "boolean"
        },
        "sessionID": {
          "type": "string"
        },
        "updatedAt": {
          "type": "integer",
          "format": "int64"
        }
      }
    }
  }
}`))
}
