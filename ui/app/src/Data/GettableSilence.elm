{-
   Alertmanager API
   API of the Prometheus Alertmanager (https://github.com/prometheus/alertmanager)

   OpenAPI spec version: 0.0.1

   NOTE: This file is auto generated by the openapi-generator.
   https://github.com/openapitools/openapi-generator.git
   Do not edit this file manually.
-}


module Data.GettableSilence exposing (GettableSilence, decoder, encoder)

import Data.Matchers as Matchers exposing (Matchers)
import Data.SilenceStatus as SilenceStatus exposing (SilenceStatus)
import DateTime exposing (DateTime)
import Dict exposing (Dict)
import Json.Decode as Decode exposing (Decoder)
import Json.Decode.Pipeline exposing (optional, required)
import Json.Encode as Encode


type alias GettableSilence =
    { matchers : Matchers
    , startsAt : DateTime
    , endsAt : DateTime
    , createdBy : String
    , comment : String
    , id : String
    , status : SilenceStatus
    , updatedAt : DateTime
    }


decoder : Decoder GettableSilence
decoder =
    Decode.succeed GettableSilence
        |> required "matchers" Matchers.decoder
        |> required "startsAt" DateTime.decoder
        |> required "endsAt" DateTime.decoder
        |> required "createdBy" Decode.string
        |> required "comment" Decode.string
        |> required "id" Decode.string
        |> required "status" SilenceStatus.decoder
        |> required "updatedAt" DateTime.decoder


encoder : GettableSilence -> Encode.Value
encoder model =
    Encode.object
        [ ( "matchers", Matchers.encoder model.matchers )
        , ( "startsAt", DateTime.encoder model.startsAt )
        , ( "endsAt", DateTime.encoder model.endsAt )
        , ( "createdBy", Encode.string model.createdBy )
        , ( "comment", Encode.string model.comment )
        , ( "id", Encode.string model.id )
        , ( "status", SilenceStatus.encoder model.status )
        , ( "updatedAt", DateTime.encoder model.updatedAt )
        ]
