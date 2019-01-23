package main

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
)

// コンテナ外部のcronからファイル指定実行する
// docker-compose exec web go run app/cron/notifyreminder.go
func main()  {
	services.InitDB()

}
