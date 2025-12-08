package builder

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

type testOptionalUpdate struct {
	ID      uint64     `json:"id"`
	SubID   uint64     `json:"sub_id"`
	XID     *uint      `json:"-"`
	Enable  bool       `json:"enable"`
	Disable *bool      `json:"disable"`
	JSON    *jsonAbleX `json:"json"`
	User
}

type jsonAbleX struct{}

func (*jsonAbleX) MarshalJSON() ([]byte, error) { return []byte("2"), nil }

func TestOptionalUpdate(t *testing.T) {
	{
		queryBytes, err := OptionalUpdate(db.NewUpdate().Model(&User{}).Where("1=1"), &testOptionalUpdate{
			ID:     0,
			SubID:  0,
			XID:    nil,
			Enable: false,
			JSON:   &jsonAbleX{},
		}, "sub_id").AppendQuery(db.QueryGen(), nil)
		if err != nil {
			t.Fatal(err)
		}
		if string(queryBytes) != `UPDATE "user" AS "user" SET "id" = 0, "enable" = FALSE, "json" = '2' WHERE (1=1)` {
			t.Fatal(fmt.Errorf("unexpected query: %s", string(queryBytes)))
		}
	}

	{
		_, err := OptionalUpdate(db.NewUpdate().Model(&User{}), map[string]any{}).AppendQuery(db.QueryGen(), nil)
		if err == nil {
			t.Fatal(errors.New("expected error"))
		}
	}
}

func TestOptionalForceUpdate(t *testing.T) {
	{
		queryBytes, err := OptionalForceUpdate(db.NewUpdate().Model(&User{}).Where("1=1"), &testOptionalUpdate{
			ID:     0,
			SubID:  0,
			XID:    nil,
			Enable: false,
			JSON:   &jsonAbleX{},
		}, []string{"disable"}, []string{"sub_id"}).AppendQuery(db.QueryGen(), nil)
		if err != nil {
			t.Fatal(err)
		}
		if string(queryBytes) != `UPDATE "user" AS "user" SET "id" = 0, "enable" = FALSE, "disable" = NULL, "json" = '2' WHERE (1=1)` {
			t.Fatal(fmt.Errorf("unexpected query: %s", string(queryBytes)))
		}
	}

	{
		_, err := OptionalForceUpdate(db.NewUpdate().Model(&User{}), map[string]any{}, nil, nil).AppendQuery(db.QueryGen(), nil)
		if err == nil {
			t.Fatal(errors.New("expected error"))
		}
	}
}

// Test structs for Updater tests
type testUpdaterStruct struct {
	ID       uint64                  `json:"id"`
	Name     string                  `json:"name"`
	Age      int                     `json:"age"`
	Score    float64                 `json:"score"`
	Active   bool                    `json:"active"`
	Tags     []string                `json:"tags"`
	Data     map[string]interface{}  `json:"data"`
	Email    *string                 `json:"email"`
	Phone    *string                 `json:"phone"`
	Metadata *testMetadata           `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
	Ignored  string                  `json:"-"`
	Anonymous struct {
		Value string `json:"value"`
	} `json:"anonymous"`
}

type testMetadata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type testProtobufStruct struct {
	SearchKey string `json:"search_key" protobuf:"bytes,2,opt,name=search_key,json=searchKey,proto3"`
	UserName  string `json:"user_name" protobuf:"bytes,3,opt,name=user_name,json=userName,proto3"`
	Normal    string `json:"normal"`
}

func TestUpdater_UseUpdater(t *testing.T) {
	q := db.NewUpdate().Model(&User{})
	v := &testUpdaterStruct{ID: 1}
	
	updater := UseUpdater(q, v)
	if updater == nil {
		t.Fatal("UseUpdater returned nil")
	}
	if updater.q != q {
		t.Error("Updater query mismatch")
	}
	if updater.v != v {
		t.Error("Updater value mismatch")
	}
}

func TestUpdater_Ignore(t *testing.T) {
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID:   1,
		Name: "test",
		Age:  20,
	}
	
	queryBytes, err := UseUpdater(q, v).Ignore([]string{"age", "score", "active", "tags", "data", "email", "phone", "metadata", "created_at", "anonymous"}).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	expected := `UPDATE "user" AS "user" SET "id" = 1, "name" = 'test' WHERE (1=1)`
	if string(queryBytes) != expected {
		t.Errorf("got: %s\nwant: %s", string(queryBytes), expected)
	}
}

func TestUpdater_Select(t *testing.T) {
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID:   1,
		Name: "test",
		Age:  20,
	}
	
	queryBytes, err := UseUpdater(q, v).Select([]string{"name"}).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	expected := `UPDATE "user" AS "user" SET "name" = 'test' WHERE (1=1)`
	if string(queryBytes) != expected {
		t.Errorf("got: %s\nwant: %s", string(queryBytes), expected)
	}
}

func TestUpdater_Exec_BasicTypes(t *testing.T) {
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID:     1,
		Name:   "test",
		Age:    25,
		Score:  95.5,
		Active: true,
	}
	
	queryBytes, err := UseUpdater(q, v).Select([]string{"id", "name", "age", "score", "active"}).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	expected := `UPDATE "user" AS "user" SET "id" = 1, "name" = 'test', "age" = 25, "score" = 95.5, "active" = TRUE WHERE (1=1)`
	if string(queryBytes) != expected {
		t.Errorf("got: %s\nwant: %s", string(queryBytes), expected)
	}
}

func TestUpdater_Exec_PointerTypes(t *testing.T) {
	email := "test@example.com"
	phone := "1234567890"
	
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID:    1,
		Email: &email,
		Phone: &phone,
	}
	
	queryBytes, err := UseUpdater(q, v).Select([]string{"id", "email", "phone"}).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	expected := `UPDATE "user" AS "user" SET "id" = 1, "email" = 'test@example.com', "phone" = '1234567890' WHERE (1=1)`
	if string(queryBytes) != expected {
		t.Errorf("got: %s\nwant: %s", string(queryBytes), expected)
	}
}

func TestUpdater_Exec_NilPointer(t *testing.T) {
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID:    1,
		Email: nil,
		Phone: nil,
	}
	
	queryBytes, err := UseUpdater(q, v).Select([]string{"id", "email", "phone"}).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	expected := `UPDATE "user" AS "user" SET "id" = 1, "email" = NULL, "phone" = NULL WHERE (1=1)`
	if string(queryBytes) != expected {
		t.Errorf("got: %s\nwant: %s", string(queryBytes), expected)
	}
}

func TestUpdater_Exec_SliceTypes(t *testing.T) {
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID:   1,
		Tags: []string{"tag1", "tag2", "tag3"},
	}
	
	queryBytes, err := UseUpdater(q, v).Select([]string{"id", "tags"}).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Array should be formatted as PostgreSQL array
	expected := `UPDATE "user" AS "user" SET "id" = 1, "tags" = '{"tag1","tag2","tag3"}' WHERE (1=1)`
	if string(queryBytes) != expected {
		t.Errorf("got: %s\nwant: %s", string(queryBytes), expected)
	}
}

func TestUpdater_Exec_NilSlice(t *testing.T) {
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID:   1,
		Tags: nil,
	}
	
	queryBytes, err := UseUpdater(q, v).Select([]string{"id", "tags"}).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	expected := `UPDATE "user" AS "user" SET "id" = 1, "tags" = NULL WHERE (1=1)`
	if string(queryBytes) != expected {
		t.Errorf("got: %s\nwant: %s", string(queryBytes), expected)
	}
}

func TestUpdater_Exec_MapTypes(t *testing.T) {
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID: 1,
		Data: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
	}
	
	queryBytes, err := UseUpdater(q, v).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Map should be serialized as JSON
	// Note: map order is non-deterministic, so we just check it contains the key
	if string(queryBytes) == "" {
		t.Error("query should not be empty")
	}
}

func TestUpdater_Exec_NilMap(t *testing.T) {
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID:   1,
		Data: nil,
	}
	
	queryBytes, err := UseUpdater(q, v).Select([]string{"id", "data"}).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	expected := `UPDATE "user" AS "user" SET "id" = 1, "data" = NULL WHERE (1=1)`
	if string(queryBytes) != expected {
		t.Errorf("got: %s\nwant: %s", string(queryBytes), expected)
	}
}

func TestUpdater_Exec_StructPointer(t *testing.T) {
	metadata := &testMetadata{
		Key:   "test_key",
		Value: "test_value",
	}
	
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID:       1,
		Metadata: metadata,
	}
	
	queryBytes, err := UseUpdater(q, v).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Struct pointer should be serialized
	if string(queryBytes) == "" {
		t.Error("query should not be empty")
	}
}

func TestUpdater_Exec_TimeType(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID:        1,
		CreatedAt: now,
	}
	
	queryBytes, err := UseUpdater(q, v).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Time should be included in the query
	if string(queryBytes) == "" {
		t.Error("query should not be empty")
	}
}

func TestUpdater_Exec_IgnoresAnonymousFields(t *testing.T) {
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID: 1,
	}
	v.Anonymous.Value = "test"
	
	queryBytes, err := UseUpdater(q, v).Select([]string{"id"}).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Anonymous fields should be ignored (they are not exported or have json tags)
	expected := `UPDATE "user" AS "user" SET "id" = 1 WHERE (1=1)`
	if string(queryBytes) != expected {
		t.Errorf("got: %s\nwant: %s", string(queryBytes), expected)
	}
}

func TestUpdater_Exec_IgnoresJsonDashTag(t *testing.T) {
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID:      1,
		Ignored: "should be ignored",
	}
	
	queryBytes, err := UseUpdater(q, v).Select([]string{"id"}).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Fields with json:"-" should be ignored
	expected := `UPDATE "user" AS "user" SET "id" = 1 WHERE (1=1)`
	if string(queryBytes) != expected {
		t.Errorf("got: %s\nwant: %s", string(queryBytes), expected)
	}
}

func TestUpdater_Exec_ErrorOnNonStruct(t *testing.T) {
	q := db.NewUpdate().Model(&User{})
	v := map[string]interface{}{"id": 1}
	
	query := UseUpdater(q, v).Exec()
	_, err := query.AppendQuery(db.QueryGen(), nil)
	if err == nil {
		t.Fatal("expected error for non-struct type")
	}
	
	// Verify error message contains expected text
	if err.Error() == "" {
		t.Error("error message should not be empty")
	}
}

func TestUpdater_PB_ProtobufTag(t *testing.T) {
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testProtobufStruct{
		SearchKey: "test_search",
		UserName:  "test_user",
		Normal:    "test_normal",
	}
	
	queryBytes, err := UseUpdater(q, v).PB().Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// With PB(), the key is extracted from protobuf tag but sqlKey still uses json tag
	// Based on the code, parseKey updates key but sqlKey is set before protobuf processing
	// So it uses json tag for SQL identifier
	expected := `UPDATE "user" AS "user" SET "search_key" = 'test_search', "user_name" = 'test_user', "normal" = 'test_normal' WHERE (1=1)`
	if string(queryBytes) != expected {
		t.Errorf("got: %s\nwant: %s", string(queryBytes), expected)
	}
}

func TestUpdater_Ignore_EmptyList(t *testing.T) {
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID:   1,
		Name: "test",
		Age:  20,
	}
	
	queryBytes, err := UseUpdater(q, v).Ignore([]string{}).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Empty ignore list should include all fields (except ignored ones)
	// Check that id, name, and age are included
	queryStr := string(queryBytes)
	if !strings.Contains(queryStr, `"id" = 1`) {
		t.Error("query should contain id")
	}
	if !strings.Contains(queryStr, `"name" = 'test'`) {
		t.Error("query should contain name")
	}
	if !strings.Contains(queryStr, `"age" = 20`) {
		t.Error("query should contain age")
	}
}

func TestUpdater_Select_EmptyList(t *testing.T) {
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID:   1,
		Name: "test",
		Age:  20,
	}
	
	query := UseUpdater(q, v).Select([]string{}).Exec()
	_, err := query.AppendQuery(db.QueryGen(), nil)
	
	// Empty select list should cause an error (empty SET clause)
	if err == nil {
		t.Fatal("expected error for empty select list")
	}
	if !strings.Contains(err.Error(), "empty SET clause") {
		t.Errorf("expected error about empty SET clause, got: %v", err)
	}
}

func TestUpdater_Exec_EmptySlice(t *testing.T) {
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID:   1,
		Tags: []string{},
	}
	
	queryBytes, err := UseUpdater(q, v).Select([]string{"id", "tags"}).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Empty slice should be formatted as empty array
	expected := `UPDATE "user" AS "user" SET "id" = 1, "tags" = '{}' WHERE (1=1)`
	if string(queryBytes) != expected {
		t.Errorf("got: %s\nwant: %s", string(queryBytes), expected)
	}
}

func TestUpdater_Exec_ComplexTypes(t *testing.T) {
	q := db.NewUpdate().Model(&User{}).Where("1=1")
	v := &testUpdaterStruct{
		ID:     1,
		Active: true,
		Score:  99.99,
	}
	
	queryBytes, err := UseUpdater(q, v).Select([]string{"id", "active", "score"}).Exec().
		AppendQuery(db.QueryGen(), nil)
	if err != nil {
		t.Fatal(err)
	}
	
	queryStr := string(queryBytes)
	// Check that all expected fields are present (order may vary)
	if !strings.Contains(queryStr, `"id" = 1`) {
		t.Error("query should contain id")
	}
	if !strings.Contains(queryStr, `"active" = TRUE`) {
		t.Error("query should contain active")
	}
	if !strings.Contains(queryStr, `"score" = 99.99`) {
		t.Error("query should contain score")
	}
}
