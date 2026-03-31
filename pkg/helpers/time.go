package helper

import (
	"fmt"
	"time"
)

func GetOrderDate() string {
	y, m, d := time.Now().Date()
	return fmt.Sprintf("%d%d%d", y, int(m), d)
}
