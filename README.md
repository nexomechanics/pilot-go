# Pilot Go SDK

```bash
go get github.com/nexomechanics/pilot-go
```

## Usage

```go
import "github.com/nexomechanics/pilot-go"

client := pilot.New("pk_your_api_key")
result, err := client.Send("prod_errors", "Deploy failed on server 3")
if err != nil {
    log.Fatal(err)
}
fmt.Println("remaining:", result.Remaining)
```

## Error handling

```go
result, err := client.Send("prod_errors", "Deploy failed")
if err != nil {
    if e, ok := err.(*pilot.PilotError); ok {
        fmt.Printf("status %d: %s\n", e.StatusCode, e.Message)
    }
}
```
