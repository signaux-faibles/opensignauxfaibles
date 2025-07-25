package engine

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"time"

	"github.com/globalsign/mgo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"github.com/spf13/viper"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// MigrationFS are postgres migration files
var MigrationFS, _ = fs.Sub(migrationFS, "migrations")

// VersionTable is the name of the table to track migrations
const VersionTable = "sfdata_migrations"

// DB type centralisant les accès à une base de données
type DB struct {
	DB         *mgo.Database
	DBStatus   *mgo.Database
	PostgresDB *pgxpool.Pool
}

// InitDB Initialisation de la connexion MongoDB et de la base de données
// PostgreSQL.
// Cette fonction réalise les migrations - le cas échéant - de la base
// PostgreSQL.
func InitDB() (DB, error) {
	dbDial := viper.GetString("DB_DIAL")
	dbDatabase := viper.GetString("DB")

	// définition de 2 connexions pour isoler les requêtes (TODO: utile ?)
	mongostatus, err := mgo.Dial(dbDial)
	if err != nil {
		return DB{}, fmt.Errorf("erreur de connexion (status) à MongoDB : %w", err)
	}
	mongostatus.SetSocketTimeout(72000 * time.Second)

	mongodb, err := mgo.Dial(dbDial)
	if err != nil {
		return DB{}, fmt.Errorf("erreur de connexion (data) à MongoDB : %w", err)
	}
	mongodb.SetSocketTimeout(72000 * time.Second)
	dbstatus := mongostatus.DB(dbDatabase)
	db := mongodb.DB(dbDatabase)

	// Création d'index sur la collection Admin, pour selection et tri de GetBatches()
	_ = db.C("Admin").EnsureIndex(mgo.Index{
		Key: []string{"_id.type", "_id.key"},
	})

	conn, err := pgxpool.New(context.Background(), viper.GetString("POSTGRES_DB_URL"))

	if err == nil {
		// Test connectivity with postgreSQL database
		err = conn.Ping(context.Background())
	}

	if err != nil {
		return DB{}, fmt.Errorf("erreur de connexion à PostgreSQL : %w", err)
	}

	// Run database migrations
	if err := runMigrations(conn); err != nil {
		return DB{}, fmt.Errorf("erreur lors de l'exécution des migrations : %w", err)
	}

	return DB{
		DB:         db,
		DBStatus:   dbstatus,
		PostgresDB: conn,
	}, nil
}

// runMigrations executes database migrations using Tern
func runMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("impossible d'acquérir une connexion du pool : %w", err)
	}
	defer conn.Release()

	migrator, err := migrate.NewMigrator(ctx, conn.Conn(), "sfdata_migrations")
	if err != nil {
		return fmt.Errorf("impossible d'initialiser les migrations : %w", err)
	}

	if err := migrator.LoadMigrations(MigrationFS); err != nil {
		return fmt.Errorf("impossible de charger les migrations : %w", err)
	}

	currentVersion, errCV := migrator.GetCurrentVersion(ctx)

	if err := migrator.Migrate(ctx); err != nil {
		return fmt.Errorf("échec de la migration : %w", err)
	}

	newVersion, errNV := migrator.GetCurrentVersion(ctx)

	if newVersion != currentVersion && errCV == nil && errNV == nil {
		slog.Info("Migrations exécutées", "version initiale", currentVersion, "nouvelle version", newVersion)
	}

	return nil
}
