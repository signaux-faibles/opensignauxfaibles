package engine

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"

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
	PostgresDB *pgxpool.Pool
}

// InitDB Initialisation de la la base de données PostgreSQL.
// Cette fonction réalise les migrations - le cas échéant - de la base
// PostgreSQL.
func InitDB() (DB, error) {
<<<<<<< HEAD
	conn, err := pgxpool.New(context.Background(), viper.GetString("POSTGRES_DB_URL"))
=======
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

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, viper.GetString("POSTGRES_DB_URL"))

	poolConn, _ := conn.Acquire(ctx)
	rows, _ := poolConn.Conn().Query(ctx, `
  SELECT table_name
  FROM information_schema.tables
  WHERE table_schema='public'`)
	for rows.Next() {
		var name string
		rows.Scan(&name)
	}
>>>>>>> fe5bd47 (Refine migrations + frontend_ap views)

	if err == nil {
		// Test connectivity with postgreSQL database
		err = conn.Ping(ctx)
	}

	if err != nil {
		return DB{}, fmt.Errorf("erreur de connexion à PostgreSQL : %w", err)
	}

	// Run database migrations
	if err := runMigrations(ctx, conn); err != nil {
		return DB{}, fmt.Errorf("erreur lors de l'exécution des migrations : %w", err)
	}

	return DB{
		PostgresDB: conn,
	}, nil
}

// runMigrations executes database migrations using Tern
func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("impossible d'acquérir une connexion du pool : %w", err)
	}
	defer conn.Release()

	migrator, err := migrate.NewMigrator(ctx, conn.Conn(), VersionTable)
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
