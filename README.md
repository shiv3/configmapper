# configmapper

Golang Config Parser for Kubernetes Configmap with Toml


## Usage

Defines a config struct
```go
type Config struct {
    EnvConfig EnvConfig `mapstructure:"env"`
}
type EnvConfig struct {
    Env         string `validate:"required"`
    ServiceName string `mapstructure:"service_name" validate:"required"`
    ProjectID   string `mapstructure:"project_id" validate:"required"`
    LogLevel    string `mapstructure:"log_level" validate:"required"`
}
```

Create a configmap file
`manifests/configmap.yaml`
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-configmap
  labels:
    app: test
data:
  config.toml: |-
    [env]
        env = "1"
        service_name = "2"
        project_id = "3"
        log_level = "4"
```

Initialize config from struct
`main.go`
```go
c, err := configmapper.Initialize(Config{})
if err != nil {
    panic(err) // error handling
}
config := c.(Config)
// use config
...

```

if run local, run golang with LOCAL=1
```sh
LOCAL=1 go run main.go
```
