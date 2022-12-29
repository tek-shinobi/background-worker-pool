# Simple Worker Queue and Background Worker Pool Template

---

- work is submitted to Work queue using SIGTERM signal. This is just a placeholder. In real, replace this channel with you modus operandi of submitting work
- If you want work task cancellation based on context, add a context arguement to workerProcess like so `workerProcess(ctx context.Context)` and then handle to timeout when doing the work

- for building, running and testing the code,

```
go build main.go
./main
```

(note the two step process. This needs to be compiled. Simply running `go run main.go` will not be able to send the SIGTERM signal to the application process (main), instead you will end up sending the signals to the go run command. If you use some other way to submit work to task queue, this limitation won't be there)

To test, in another terminal, run `kill -s TERM <pid>`

for example, `kill -s TERM 24583` where `24583` is the pid of the main process
