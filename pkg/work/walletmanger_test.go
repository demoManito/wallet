package work

import (
	"fmt"
	"testing"

	"wallet/models"
)

func TestLargeVolume(t *testing.T) {
	for _, id := range LargeVolume {
		query := models.FormatWalletManage(id)

		fmt.Println(query)
	}
}
