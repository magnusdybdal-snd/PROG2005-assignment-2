package utils

import (
	"context"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const testCollectionName = "FuncUtil_test"

func TestMain(m *testing.M) {
	log.Println("Initializing Firestore for utils test...")

	if err := InitFirestore(); err != nil {
		log.Fatalf("FATAL: TestMain failed to initialize Firestore via utils.InitFirestore: %v", err)
	}
	log.Println("Firestore client initialized for utils tests")

	exitCode := m.Run()

	log.Println("Closing Firestore client...")
	CloseFirestore()
	log.Println("Firestore client closed.")
	os.Exit(exitCode)
}

func TestGetFirestoreDoc(t *testing.T) {
	// Start by checking that the firestore client is working
	if FirestoreClient == nil {
		t.Fatal("FATAL: FirestoreClient is nil in test function.")
	}
	ctx := context.Background()

	// === Test Case 1: Success path using dashboardconfig as data example ===
	t.Run("Success getting DashboardConfig document", func(t *testing.T) {

		// Data we expect to get back on successfull call
		expectedData := DashboardConfig{
			Country: "Test Norway",
			IsoCode: "TN",
			Features: struct {
				Temperature      bool     `firestore:"temperature" json:"temperature"`
				Precipitation    bool     `firestore:"precipitation" json:"precipitation"`
				Capital          bool     `firestore:"capital" json:"capital"`
				Coordinates      bool     `firestore:"coordinates" json:"coordinates"`
				Population       bool     `firestore:"population" json:"population"`
				Area             bool     `firestore:"area" json:"area"`
				TargetCurrencies []string `firestore:"targetCurrencies" json:"targetCurrencies"`
			}{
				Temperature:      true,
				Precipitation:    true,
				Capital:          false,
				Coordinates:      true,
				Population:       false,
				Area:             true,
				TargetCurrencies: []string{"EUR", "SEK"},
			},
			LastRetrieval: time.Date(2025, 4, 8, 15, 30, 45, 0, time.UTC), // Fixed time for testing
		}

		// 1. Setup: Add the document to the database for testing
		docID, _, err := FirestoreClient.Collection(testCollectionName).Add(ctx, expectedData)
		if err != nil {
			t.Fatalf("Setup failed: Could not add test document in firestore %s: %v", docID.ID, err)
		}
		// Cleans up using t.Cleanup (runs after all tests are done)
		t.Cleanup(func() {
			_, err2 := docID.Delete(ctx)
			if err2 != nil {
				t.Errorf("Document: %v was not deleted", docID)
			}
		})

		// 2. Run test
		actualData, actualErr := GetFirestoreDocument[DashboardConfig](ctx, docID.ID, testCollectionName)

		// 3. Assertions
		if actualErr != nil {
			t.Fatalf("GetFirestoreDocument[DashboardConfig]() returned an unexpected error: %v", actualErr.Error())
		}

		if !reflect.DeepEqual(actualData, expectedData) {
			t.Errorf("GetFirestoreDocument[DashboardConfig]() returned unexpected data.\nGot:\n%#v\nWant:\n%#v", actualData, expectedData)
		}
	})

	// === Test Case 2: Document not found ===
	t.Run("Document not found", func(t *testing.T) {

		// 1. Setup a fake doc that does not exist
		docID := "this-doc-does-not-exist"

		// 2. Run test
		_, actualErr := GetFirestoreDocument[DashboardConfig](ctx, docID, testCollectionName)

		// 3. Assertions
		if actualErr == nil {
			t.Fatalf("GetFirestoreDocument[DashboardConfig]() expected an error for non-existent document but got nil")
		}
		// 4. We attempt to extract gRPC error codes from firebase (Google Remote Procedure Call) This is a framework for making
		// calls between services and used by google cloud services like firestore
		st, ok := status.FromError(actualErr)
		if !ok { // If status code was not a gRPC code
			t.Errorf("GetFirestoreDocument() returned an error that could not be converted to gRPC: %v", actualErr)
		} else if st.Code() != codes.NotFound {
			t.Errorf("GetFirestoreDocument() expected a Firestore 'Not Found' (code %d), but got: %d: %v", codes.NotFound, st.Code(), actualErr)
		}
	})

	// === Test Case 3: Data Unmarshal Error ===
	t.Run("Error unmarshaling data", func(t *testing.T) {

		// Data in Firestore where 'area' (bool in struct) is stored as a string
		firestoreData := map[string]interface{}{
			"Country": "Mismatch Landia",
			"IsoCode": "XX",
			"Features": map[string]interface{}{
				"Area": "this-is-not-a-bool", // Incompatible type for DashboardConfig.Features.Area
			},
			"LastRetrieval": time.Now(),
		}

		// 1. Setup: Add the document with incompatible data
		docID, _, err := FirestoreClient.Collection(testCollectionName).Add(ctx, firestoreData)
		if err != nil {
			t.Fatalf("Setup failed: Could not add test document in firestore %s: %v", docID.ID, err)
		}
		// Cleans up using t.Cleanup (runs after all tests are done)
		t.Cleanup(func() {
			_, err2 := docID.Delete(ctx)
			if err2 != nil {
				t.Errorf("Document: %v was not deleted", docID)
			}
		})

		// 2. Run test
		_, actualErr := GetFirestoreDocument[DashboardConfig](ctx, docID.ID, testCollectionName)

		// 3. Assertions
		if actualErr == nil {
			t.Fatalf("GetFirestoreDocument[DashboardConfig]() expected an unmarshal error, but got nil")
		}
		expectedErrorMsg := "Failed to unmarshal data for document"
		if !strings.Contains(actualErr.Error(), expectedErrorMsg) {
			t.Errorf("GetFirestoreDocument[DashboardConfig]() error = %q, want error containing %q", actualErr.Error(), expectedErrorMsg)
		}
	})

}
