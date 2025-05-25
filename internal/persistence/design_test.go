package persistence

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/arvaliullin/wapa/internal/domain"
	_ "github.com/lib/pq"
)

const testDbConn = "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable"

func cleanDB(t *testing.T) {
	err := WithConnection(testDbConn, func(conn *sql.DB) error {
		_, err := conn.Exec("TRUNCATE TABLE composer.design RESTART IDENTITY CASCADE")
		return err
	})
	if err != nil {
		t.Fatalf("failed to clean database: %v", err)
	}
}

func TestCreateDesign(t *testing.T) {
	cleanDB(t)
	repo := &DesignRepository{DbConnection: testDbConn}

	design := domain.Design{
		Name: "TestDesign",
		Lang: "go",
		JS:   "test.js",
		Wasm: "test.wasm",
		Functions: []domain.Function{
			{Function: "testFunction", Args: []float64{1, 2, 3}},
		},
	}

	id, err := repo.Create(design)
	if err != nil {
		t.Fatalf("failed to create design: %v", err)
	}
	if id == "" {
		t.Fatalf("expected non-empty id, got empty")
	}

	err = WithConnection(testDbConn, func(conn *sql.DB) error {
		var count int
		query := `SELECT COUNT(*) FROM composer.design WHERE id = $1`
		err := conn.QueryRow(query, id).Scan(&count)
		if err != nil {
			return err
		}
		if count != 1 {
			t.Fatalf("expected 1 record with id %s, but found %d", id, count)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("database validation failed: %v", err)
	}
}

func TestGetByID(t *testing.T) {
	cleanDB(t)
	repo := &DesignRepository{DbConnection: testDbConn}

	design := domain.Design{
		Name: "TestGetByID",
		Lang: "go",
		JS:   "get_by_id.js",
		Wasm: "get_by_id.wasm",
		Functions: []domain.Function{
			{Function: "exampleFunction", Args: []float64{1.0, 2.0}},
		},
	}

	id, err := repo.Create(design)
	if err != nil {
		t.Fatalf("failed to create design: %v", err)
	}

	foundDesign, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("failed to get design by id: %v", err)
	}

	if foundDesign.ID != id {
		t.Errorf("expected id %s, got %s", id, foundDesign.ID)
	}
	if !reflect.DeepEqual(foundDesign.Functions, design.Functions) {
		t.Errorf("expected functions %v, got %v", design.Functions, foundDesign.Functions)
	}
}

func TestUpdateDesign(t *testing.T) {
	cleanDB(t)
	repo := &DesignRepository{DbConnection: testDbConn}

	design := domain.Design{
		Name: "BeforeUpdate",
		Lang: "go",
		JS:   "before_update.js",
		Wasm: "before_update.wasm",
		Functions: []domain.Function{
			{Function: "oldFunction", Args: []float64{1.0}},
		},
	}

	id, err := repo.Create(design)
	if err != nil {
		t.Fatalf("failed to create design: %v", err)
	}

	updatedDesign := domain.Design{
		ID:   id,
		Name: "AfterUpdate",
		Lang: "cpp",
		JS:   "after_update.js",
		Wasm: "after_update.wasm",
		Functions: []domain.Function{
			{Function: "newFunction", Args: []float64{2.0}},
		},
	}

	err = repo.Update(updatedDesign)
	if err != nil {
		t.Fatalf("failed to update design: %v", err)
	}

	foundDesign, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("failed to get design after update: %v", err)
	}

	if foundDesign.Name != updatedDesign.Name {
		t.Errorf("expected name %s, got %s", updatedDesign.Name, foundDesign.Name)
	}
	if foundDesign.Lang != updatedDesign.Lang {
		t.Errorf("expected lang %s, got %s", updatedDesign.Lang, foundDesign.Lang)
	}
	if !reflect.DeepEqual(foundDesign.Functions, updatedDesign.Functions) {
		t.Errorf("expected functions %v, got %v", updatedDesign.Functions, foundDesign.Functions)
	}
}

func TestDeleteDesign(t *testing.T) {
	cleanDB(t)
	repo := &DesignRepository{DbConnection: testDbConn}

	design := domain.Design{
		Name: "TestDelete",
		Lang: "go",
		JS:   "delete.js",
		Wasm: "delete.wasm",
		Functions: []domain.Function{
			{Function: "deleteFunction", Args: []float64{1.0}},
		},
	}

	id, err := repo.Create(design)
	if err != nil {
		t.Fatalf("failed to create design: %v", err)
	}

	err = repo.Delete(id)
	if err != nil {
		t.Fatalf("failed to delete design: %v", err)
	}
}

func TestGetAllDesigns(t *testing.T) {
	cleanDB(t)
	repo := &DesignRepository{DbConnection: testDbConn}

	initialDesigns := []domain.Design{
		{
			Name: "Design1",
			Lang: "go",
			JS:   "design1.js",
			Wasm: "design1.wasm",
			Functions: []domain.Function{
				{Function: "func1", Args: []float64{1, 2}},
			},
		},
		{
			Name: "Design2",
			Lang: "rust",
			JS:   "design2.js",
			Wasm: "design2.wasm",
			Functions: []domain.Function{
				{Function: "func2", Args: []float64{3, 4}},
			},
		},
	}

	for _, d := range initialDesigns {
		_, err := repo.Create(d)
		if err != nil {
			t.Fatalf("failed to create design during setup: %v", err)
		}
	}

	designs, err := repo.GetAll()
	if err != nil {
		t.Fatalf("failed to get all designs: %v", err)
	}

	if len(designs) != len(initialDesigns) {
		t.Errorf("expected %d designs, got %d", len(initialDesigns), len(designs))
	}

	for i, d := range designs {
		expected := initialDesigns[i]
		if d.Name != expected.Name || d.Lang != expected.Lang || d.JS != expected.JS {
			t.Errorf("got unexpected design at index %d: %+v", i, d)
		}
		expectedFunctions, _ := json.Marshal(expected.Functions)
		actualFunctions, _ := json.Marshal(d.Functions)
		if string(expectedFunctions) != string(actualFunctions) {
			t.Errorf("expected functions %s, got %s", string(expectedFunctions), string(actualFunctions))
		}
	}
}
