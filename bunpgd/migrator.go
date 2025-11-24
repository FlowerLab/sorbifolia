package bunpgd

import (
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"github.com/uptrace/bun/migrate/sqlschema"
)

var _ sqlschema.MigratorDialect = (*Dialect)(nil)

type migratorFunc func(b []byte, op any) ([]byte, error)

func (f migratorFunc) AppendSQL(b []byte, op any) ([]byte, error) { return f(b, op) }

func (d *Dialect) NewMigrator(db *bun.DB, schemaName string) sqlschema.Migrator {
	return migratorFunc(func(b []byte, op any) (_ []byte, err error) {
		gen := db.QueryGen()

		switch op := op.(type) {
		case *migrate.CreateTableOp:
			return sqlschema.NewBaseMigrator(db).AppendCreateTable(b, op)

		case *migrate.DropTableOp:
			return sqlschema.NewBaseMigrator(db).AppendDropTable(b, schemaName, op.TableName)

		case *migrate.RenameTableOp:
			b = append(b, "ALTER TABLE "...)
			b = gen.AppendQuery(b, "?.?", bun.Ident(schemaName), bun.Ident(op.TableName))
			b = append(b, " RENAME TO "...)
			b = gen.AppendName(b, op.NewName)

		case *migrate.RenameColumnOp:
			b = append(b, "ALTER TABLE "...)
			b = gen.AppendQuery(b, "?.?", bun.Ident(schemaName), bun.Ident(op.TableName))
			b = append(b, " RENAME COLUMN "...)
			b = gen.AppendName(b, op.OldName)

			b = append(b, " TO "...)
			b = gen.AppendName(b, op.NewName)

		case *migrate.AddColumnOp:
			b = append(b, "ALTER TABLE "...)
			b = gen.AppendQuery(b, "?.?", bun.Ident(schemaName), bun.Ident(op.TableName))
			b = append(b, " ADD COLUMN "...)
			b = gen.AppendName(b, op.ColumnName)
			b = append(b, " "...)

			b, err = op.Column.AppendQuery(gen, b)
			if err != nil {
				return nil, err
			}

			if op.Column.GetDefaultValue() != "" {
				b = append(b, " DEFAULT "...)
				b = append(b, op.Column.GetDefaultValue()...)
				b = append(b, " "...)
			}

			if op.Column.GetIsIdentity() {
				b = d.AppendSequence(b, nil, nil)
			}

		case *migrate.DropColumnOp:
			b = append(b, "ALTER TABLE "...)
			b = gen.AppendQuery(b, "?.?", bun.Ident(schemaName), bun.Ident(op.TableName))
			b = append(b, " DROP COLUMN "...)
			b = gen.AppendName(b, op.ColumnName)

		case *migrate.AddPrimaryKeyOp:
			b = append(b, "ALTER TABLE "...)
			b = gen.AppendQuery(b, "?.?", bun.Ident(schemaName), bun.Ident(op.TableName))

			b = append(b, " ADD PRIMARY KEY ("...)
			b, _ = op.PrimaryKey.Columns.AppendQuery(gen, b)
			b = append(b, ")"...)

		case *migrate.ChangePrimaryKeyOp:
			b = append(b, "ALTER TABLE "...)
			b = gen.AppendQuery(b, "?.?", bun.Ident(schemaName), bun.Ident(op.TableName))

			b = append(b, " DROP CONSTRAINT "...)
			b = gen.AppendName(b, op.Old.Name)

			b = append(b, ",  ADD PRIMARY KEY ("...)
			b, _ = op.New.Columns.AppendQuery(gen, b)
			b = append(b, ")"...)

		case *migrate.DropPrimaryKeyOp:
			b = append(b, "ALTER TABLE "...)
			b = gen.AppendQuery(b, "?.?", bun.Ident(schemaName), bun.Ident(op.TableName))

			b = append(b, " DROP CONSTRAINT "...)
			b = gen.AppendName(b, op.PrimaryKey.Name)

		case *migrate.AddUniqueConstraintOp:
			b = append(b, "ALTER TABLE "...)
			b = gen.AppendQuery(b, "?.?", bun.Ident(schemaName), bun.Ident(op.TableName))

			b = append(b, " ADD CONSTRAINT "...)
			if op.Unique.Name != "" {
				b = gen.AppendName(b, op.Unique.Name)
			} else {
				b = gen.AppendName(b, fmt.Sprintf("%s_%s_key", op.TableName, op.Unique.Columns))
			}

			b = append(b, " UNIQUE ("...)
			b, _ = op.Unique.Columns.AppendQuery(gen, b)
			b = append(b, ")"...)

		case *migrate.DropUniqueConstraintOp:
			b = append(b, "ALTER TABLE "...)
			b = gen.AppendQuery(b, "?.?", bun.Ident(schemaName), bun.Ident(op.TableName))

			b = append(b, " DROP CONSTRAINT "...)
			b = gen.AppendName(b, op.Unique.Name)

		case *migrate.DropForeignKeyOp:
			b = append(b, "ALTER TABLE "...)
			b = gen.AppendQuery(b, "?.?", bun.Ident(schemaName), bun.Ident(op.TableName()))

			b = append(b, " DROP CONSTRAINT "...)
			b = gen.AppendName(b, op.ConstraintName)

		// case *migrate.ChangeColumnTypeOp:
		// case *migrate.AddForeignKeyOp:
		default:
			return nil, fmt.Errorf("migrator: unsupported op type %T", op)
		}

		return b, err
	})
}
