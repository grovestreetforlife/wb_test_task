// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `openapi: 3.0.0
info:
  title: Order API
  version: "0.1.0"


paths:
  /api/v1/orders/{id}:
    get:
      tags:
        - Orders
      summary: Получить заказ по id
      parameters:
        - in: path
          name: id
          schema:
            type: string
            format: uuid
          required: true
          example: 5d110e48-9e6b-4928-b436-14194b30d54f
          description: Order id

      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/SuccessResponseGetOrder'
        '400':
          description: Bad request
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

        '404':
          description: Not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

        '500':
          description: Interval Server Error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'



  /api/health/live:
    get:
      tags:
        - Health
      summary: Health live
      responses:
        '200':
          description: Healthy
          content:
            application/json:
              schema:
                type: string
                example: "Healthy"
        '503':
          description: Unhealthy
          content:
            application/json:
              schema:
                type: string
                example: "Unhealthy"

  /api/health/readiness:
    get:
      tags:
        - Health
      summary: Health readiness
      responses:
        '200':
          description: Healthy
          content:
            application/json:
              schema:
                type: string
                example: "Healthy"
        '503':
          description: Unhealthy
          content:
            application/json:
              schema:
                type: string
                example: "Unhealthy"

components:
  schemas:
    SuccessResponseGetOrder:
      properties:
        code:
          type: string
          example: OK
        status:
          type: string
          enum: [ok, fail]
        body:
          properties:
            order_uid:
              type: string
              example: 5d110e48-9e6b-4928-b436-14194b30d54f
            track_number:
              type: string
              example: WBILMTESTTRACK3
            entry:
              type: string
              example: WBIL
            delivery:
              properties:
                name:
                  type: string
                  example: Test Testov
                phone:
                  type: string
                  example: '+9720000000'
                zip:
                  type: string
                  example: '2639809'
                city:
                  type: string
                  example: Kiryat Mozkin
                address:
                  type: string
                  example: Ploshad Mira 15
                region:
                  type: string
                  example: Kraiot
                email:
                  type: string
                  format: email
                  example: test@gmail.com
            payment:
              properties:
                transaction: 
                  type: string
                  example: 5d110e48-9e6b-4928-b436-14194b30d54f
                request_id:
                  type: string
                  example: 5d110e48-9e6b-4928-b436-14194b30d54f
                currency:
                  type: string
                  example: USD
                provider:
                  type: string
                  example: wbpay
                amount: 
                  type: integer
                  example: 1817
                payment_dt:
                  type: integer
                  example: 1637907727
                bank:
                  type: string
                  example: alpha
                delivery_cost:
                  type: integer
                  example: 1500
                goods_total:
                  type: integer
                  example: 317
                custom_fee:
                  type: integer
                  example: 0
            items:
              type: array
              items:
                properties:
                  chrt_id: 
                    type: integer
                    example: 9934930
                  track_number:
                    type: string
                    example: WBILMTESTTRACK3
                  price: 
                    type: integer
                    example: 453
                  rid: 
                    type: string
                    example: ab4219087a764ae0btest
                  name:
                    type: string
                    example: Mascaras
                  sale:
                    type: integer
                    example: 30
                  size:
                    type: string
                    example: '0'
                  total_price:
                    type: integer
                    example: 317
                  nm_id:
                    type: integer
                    example: 2389212
                  brand: 
                    type: string
                    example: Vivienne Sabo
                  status: 
                    type: integer
                    example: 202
            locale:
              type: string
              example: en
            internal_signature:
              type: string
              example: ''
            customer_id:
              type: string
              example: test
            delivery_service: 
              type: string
              example: meest
            shard_key:
              type: string
              example: ''
            sm_id:
              type: integer
              example: 99
            date_created:
              type: string
              example: 2021-11-26 06:22:19 +0000 UTC
            oof_shard:
              type: string
              example: '1'
        
        error:
          type: string
          example: ""
          
    ErrorResponse:
      properties:
        code:
          type: string
          example: "ERROR_STRING_CODE"
          description: result string code
        status:
          type: string
          enum: ["ok", "fail"]
        body:
          type: object
          example: null
        error:
          type: string
          example: "ERROR_MESSAGE"
          description: error message
          
          
          
          
          
          
          
          
          
          
          
          
          
          
          
          
          
          
          
          `

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
