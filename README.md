# go_simple_backdoor
Simple bind shells written in Go

## Server
```bash
go build main.go
chmod +x ./main
./main <PORT_NUMBER>
```

## Client

- Using **Netcat**:
```bash
netcat <IP_SERVER> <PORT_NUMBER>
or
nc <IP_SERVER> <PORT_NUMBER>
```

- Using **Telnet**:
```bash
telnet <IP_SERVER> <PORT_NUMBER>
```
