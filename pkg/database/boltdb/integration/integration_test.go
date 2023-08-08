package integration_test

import (
	"os"
	"testing"

	"github.com/YuraSahanovskyi/DriveTelegramBot/pkg/database"
	"github.com/YuraSahanovskyi/DriveTelegramBot/pkg/database/boltdb"
	"github.com/spf13/viper"
)

func TestIntegration_BoltDB(t *testing.T) {
	tempDBPath := "test_.db"
	defer os.Remove(tempDBPath)

	// Backup the original config value and replace it with the temp path
	originalDBName := viper.GetString("db_name")
	viper.Set("db_name", tempDBPath)
	defer viper.Set("db_name", originalDBName)

	//Initialize database
	repo, err := boltdb.InitDB()
	if err != nil {
		t.Errorf("cannot initialize database: %v", err)
	}
	bucket := database.Code

	tests := map[int]string{
		0: "0",
		1: "",
		2: "test",
		3: "test3",
	}

	t.Run("Put", func(t *testing.T) {
		for k, v := range tests {
			if err := repo.Put(bucket, int64(k), []byte(v)); err != nil {
				t.Errorf("cannot put into database: %v", err)
			}
		}
	})

	t.Run("Get", func(t *testing.T) {
		for k, v := range tests {
			get, err := repo.Get(bucket, int64(k))
			if err != nil {
				t.Errorf("cannot get from database: %v", err)
			} else if string(get) != v {
				t.Errorf("mismatch in database: want %v, got %v", v, string(get))
			}
		}
	})

	t.Run("Delete", func(t *testing.T) {
		for k := range tests {
			if err := repo.Delete(bucket, int64(k)); err != nil {
				t.Errorf("cannot delete from database: %v", err)
			}
		}
	})
}
