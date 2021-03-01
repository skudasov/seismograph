package migrations

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"

	_ "github.com/jinzhu/gorm"

	"github.com/f4hrenh9it/seismograph/back/config"
	"github.com/f4hrenh9it/seismograph/back/testutil"
)

func GormMigrateInit(cfg config.DBConfig) {
	db, err := gorm.Open("postgres", testutil.PgConnectionString(cfg))
	if err != nil {
		log.Fatalf("Error while connecting to database: %s", err.Error())
		return
	}
	defer db.Close()

	db = db.LogMode(true)

	m := gormigrate.New(db, MigrationOptions(), InitMigrations())

	if err = m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
		return
	}
	fmt.Println("migrated successfully!")
}

func InitMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "202005180421",
			Migrate: func(tx *gorm.DB) error {
				// the initial database tables. Do not delete them

				type TestData struct {
					gorm.Model
					TestName     string
					Environment  string
					ProjectID    uint `gorm:"foreignKey:ID;references:Project"`
					DataBlobPath string
				}

				if err := tx.CreateTable(&TestData{}).Error; err != nil {
					return err
				}

				type Project struct {
					gorm.Model
					Name        string
					Description string
					RepoUrl     string
				}
				if err := tx.CreateTable(&Project{}).Error; err != nil {
					return err
				}

				type Instance struct {
					gorm.Model
					ClusterID     uint `gorm:"foreignKey:ID;references:AttackCluster"`
					Region        string
					Name          string
					PublicDNSName string
					PrivateKeyPEM string
					Image         string
					Type          string
				}
				if err := tx.CreateTable(&Instance{}).Error; err != nil {
					return err
				}

				type AttackCluster struct {
					gorm.Model
					Name         string
					ProviderName string
					Region       string
					ProjectID    uint `gorm:"foreignKey:ID;references:Project"`
				}
				if err := tx.CreateTable(&AttackCluster{}).Error; err != nil {
					return err
				}

				type AttackVM struct {
					gorm.Model
					Name            string
					AttackClusterID uint `gorm:"foreignKey:ID;references:AttackCluster"`
				}
				if err := tx.CreateTable(&AttackVM{}).Error; err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTableIfExists("test_data").Error
			},
		},
	}
}

func MigrationOptions() *gormigrate.Options {
	options := gormigrate.DefaultOptions
	options.UseTransaction = true
	options.ValidateUnknownMigrations = true
	return options
}
