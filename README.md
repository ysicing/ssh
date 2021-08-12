[Deprecated] 请使用 `ergo ops exec`

## go mod ssh

```bash
go get -u github.com/ysicing/ssh
```

## Usage

```go
s := ssh.SSH{
	User:     user,
	Password: pass,
	PkFile:   pkfile,
	}
s.Run(ip, xcmd)
```

## MacOS

```
brew install ysssh
```