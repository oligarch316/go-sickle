let types = ./types.dhall

let ObservLogConfig =
      { Type = types.ObservLogConfig
      , default =
        { encoding = "console"
        , level = "info"
        , enableCaller = False
        , enableStacktrace = False
        }
      }

let ObservConfig =
      { Type = types.ObservConfig, default.log = ObservLogConfig.default }

let PluginConfig =
      { Type = types.PluginConfig
      , default =
        { files = [] : List Text
        , directories = [] : List Text
        , trees = [] : List Text
        }
      }

let ConsumerPluginsConfig =
      { Type = types.ConsumerPluginsConfig
      , default =
        { any = [] : List Text
        , parsed = [] : List Text
        , classified = [] : List Text
        , collection = [] : List Text
        , media = [] : List Text
        }
      }

let ConsumerConfig =
      { Type = types.ConsumerConfig
      , default.plugins = ConsumerPluginsConfig.default
      }

let TransformerConfig =
      { Type = types.TransformerConfig, default.plugins = [] : List Text }

let Config =
      { Type = types.Config
      , default =
        { observ = ObservConfig.default
        , plugin = PluginConfig.default
        , consumer = ConsumerConfig.default
        , transformer = TransformerConfig.default
        }
      }

in  { ObservLogConfig
    , ObservConfig
    , PluginConfig
    , ConsumerPluginsConfig
    , ConsumerConfig
    , TransformerConfig
    , Config
    }
