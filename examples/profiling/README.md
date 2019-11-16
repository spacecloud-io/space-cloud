
1. import this inside main.go
	"github.com/pkg/profile"

2. add this first line inside main()
  defer profile.Start().Stop()

./space-cloud run --admin-user="admin" --admin-pass="admin" --admin-secret="topsecret"
go tool pprof -svg ./space-cloud /path/cpu.pprof > out.svg