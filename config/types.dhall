#!/usr/bin/env -S dhall --file

let ObservLogConfig =
      { encoding : Text
      , level : Text
      , enableCaller : Bool
      , enableStacktrace : Bool
      }

let ObservConfig = { log : ObservLogConfig }

let PluginConfig =
      { files : List Text
      , directories : List Text
      , trees : List Text
      }

let ConsumerPluginsConfig =
      { any : List Text
      , parsed : List Text
      , classified : List Text
      , collection : List Text
      , media : List Text
      }

let ConsumerConfig = { plugins : ConsumerPluginsConfig }

let TransformerConfig = { plugins : List Text }

let Config =
      { observ : ObservConfig
      , plugin : PluginConfig
      , consumer : ConsumerConfig
      , transformer : TransformerConfig
      }

in  { ObservLogConfig = ObservLogConfig
    , ObservConfig = ObservConfig
    , PluginConfig = PluginConfig
    , ConsumerPluginsConfig = ConsumerPluginsConfig
    , ConsumerConfig = ConsumerConfig
    , TransformerConfig = TransformerConfig
    , Config = Config
    }
