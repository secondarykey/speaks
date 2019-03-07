package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	. "github.com/secondarykey/speaks/config"
)

const SchemaVersion = "1.0.0"

type schemaError struct {
	code    int
	message string
}

type RecordData *string
type Record []RecordData

type FlexRows struct {
	Columns []string
	Records []Record
}

func (s *schemaError) Error() string {
	return fmt.Sprintf("%s[%d]", s.message, s.code)
}

func NewSchemaError(code int, msg string) error {
	return &schemaError{
		code:    code,
		message: msg,
	}
}

type Table interface {
	Create() error
	Drop() error
	Init(tx *sql.Tx) error
	Insert(tx *sql.Tx) (sql.Result, error)
}

var inst *sql.DB

func Evolution1_0_0to0_2_0(newPath, oldPath string) error {
	return fmt.Errorf("Not yet Implemented")
}

func checkSchemaVersion(path, ver string) error {

	//%s be
	pArr := strings.Split(path, "%s")
	if len(pArr) != 2 {
		return NewSchemaError(-1, "Error:database path is '%s' requid["+path+"]")
	}

	rpath := fmt.Sprintf(path, SchemaVersion)
	//exist database file
	_, err := os.Stat(rpath)
	//call version check
	if ver == fmt.Sprintf(`"%s"`, SchemaVersion) || ver == "test" {
		if err == nil {
			return nil
		}
		return NewSchemaError(0, "Create database")
	}

	if err == nil {
		return nil
	}
	//code 0
	return NewSchemaError(0, "Warning:Program schema version,TOML file schema version")
}

func getPath(ver string) string {

	path := Config.Base.Root + "/" + Config.Database.Path
	if ver != "" {
		path = fmt.Sprintf(path, ver)
	}

	return strings.ReplaceAll(path, `"`, "")
}

func Init() error {

	cver := Config.Database.Version
	path := getPath("")

	err := checkSchemaVersion(path, cver)
	if err != nil {

		if se, ok := err.(*schemaError); ok {
			if se.code != 0 {
				log.Println(err.Error() + "[" + path + "]")
				log.Println("Program :" + SchemaVersion)
				log.Println("Config  :" + cver)
				return err
			}
		} else {

			return err
		}
	}

	err = Listen()
	if err != nil {
		return err
	}
	err = dropTables()
	if err != nil {
		return err
	}

	//dbファイルが存在するか？
	err = createTables()
	if err != nil {
		return err
	}
	err = initTables()
	if err != nil {
		return err
	}

	inst.Close()
	return nil
}

func Listen() error {

	var err error
	path := getPath(SchemaVersion)

	log.Println("Database:" + path)

	inst, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	return nil
}

func createTables() error {

	log.Println("Create User Table")
	err := createUserTable()
	if err != nil {
		return err
	}

	log.Println("Create Role Table")
	err = createRoleTable()
	if err != nil {
		return err
	}

	log.Println("Create UserRole Table")
	err = createUserRoleTable()
	if err != nil {
		return err
	}

	log.Println("Create Project Table")
	err = createProjectTable()
	if err != nil {
		return err
	}

	log.Println("Create Member Table")
	err = createMemberTable()
	if err != nil {
		return err
	}

	log.Println("Create Category Table")
	err = createCategoryTable()
	if err != nil {
		return err
	}

	log.Println("Create Message Table")
	err = createMessageTable()
	if err != nil {
		return err
	}

	log.Println("Create Memo Table")
	err = createMemoTable()
	if err != nil {
		return err
	}
	return nil
}

func initTables() error {

	tx, err := inst.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	log.Println("Initialize User Table")
	userId, err := InitUser(tx)
	if err != nil {
		return err
	}

	log.Println("Initialize Role Table")
	err = InitRole(tx)
	if err != nil {
		return err
	}

	log.Println("Initialize UserRole Table")
	err = InitUserRole(tx, userId)
	if err != nil {
		return err
	}

	log.Println("Initialize Project Table")
	err = InitProject(tx)
	if err != nil {
		return err
	}

	log.Println("Initialize Member Table")
	err = InitMember(tx, userId)
	if err != nil {
		return err
	}

	log.Println("Initialize Category Table")
	err = InitCategory(tx)
	if err != nil {
		return err
	}

	log.Println("Commit")
	return tx.Commit()
}

func dropTables() error {

	err := dropUserRoleTable()
	if err != nil {
		return err
	}

	err = dropUserTable()
	if err != nil {
		return err
	}

	err = dropRoleTable()
	if err != nil {
		return err
	}

	err = dropMessageTable()
	if err != nil {
		return err
	}

	err = dropCategoryTable()
	if err != nil {
		return err
	}

	err = dropMemoTable()
	if err != nil {
		return err
	}

	err = dropMemberTable()
	if err != nil {
		return err
	}

	err = dropProjectTable()
	if err != nil {
		return err
	}
	return nil
}

func Exec(sql string) (sql.Result, error) {
	return inst.Exec(sql)
}

func Query(sql string) (*FlexRows, error) {

	rows, err := inst.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	flexRows := &FlexRows{}
	flexRows.Columns, err = rows.Columns()
	if err != nil {
		return nil, err
	}

	records := make([]Record, 0)

	for rows.Next() {
		rec := make([]RecordData, len(flexRows.Columns))
		for i, _ := range rec {
			rec[i] = toPtr("")
		}
		rows.Scan(rec)
		records = append(records, rec)
	}

	flexRows.Records = records
	log.Println(records)

	return flexRows, err
}

func toPtr(s string) *string {
	return &s
}
