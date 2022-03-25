#!/usr/bin/env -S dhall --file

let types = ./types.dhall let schemas = ./schemas.dhall in types // schemas
