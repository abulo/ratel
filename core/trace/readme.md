## 配置文件

```yaml
[jaeger]
EnableRPCMetrics= true
[jaeger.Reporter]
LocalAgentHostPort = "127.0.0.1:6831"
LogSpans = true
[jaeger.Sampler]
Param = 0.0001
```