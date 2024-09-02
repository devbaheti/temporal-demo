This sample demonstrates how to implement a DSL workflow. In this sample, we provide 2 sample yaml files each defines a custom workflow that can be processed by this DSL workflow sample code.

Steps to run this sample:
1) Run a [Temporal service](https://github.com/temporalio/samples-go/tree/main/#how-to-use).
2) Run
```
go run main.go
```
to start worker for dsl workflow and workers.
2) You can also write your own json config to play with it.
3) You can replace the dummy activities to your own real activities to build real workflow based on this simple DSL workflow.
