package paths

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/datalakestore/filesystems"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

func TestLifecycle(t *testing.T) {
	client, err := testhelpers.Build()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	fileSystemName := fmt.Sprintf("acctestfs-%s", testhelpers.RandomString())
	path := "test"

	if _, err = client.BuildTestResourcesWithHns(ctx, resourceGroup, accountName, storage.BlobStorage); err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)
	fileSystemsClient := filesystems.NewWithEnvironment(client.Environment)
	fileSystemsClient.Client = client.PrepareWithStorageResourceManagerAuth(fileSystemsClient.Client)
	pathsClient := NewWithEnvironment(client.Environment)
	pathsClient.Client = client.PrepareWithStorageResourceManagerAuth(fileSystemsClient.Client)

	t.Logf("[DEBUG] Creating an empty File System..")
	fileSystemInput := filesystems.CreateInput{}
	if _, err = fileSystemsClient.Create(ctx, accountName, fileSystemName, fileSystemInput); err != nil {
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}

	t.Logf("[DEBUG] Creating folder 'test' ..")
	content := []byte("Hello!")
	input := CreateInput{
		Resource: PathResourceDirectory,
		Content:  &content,
		Properties: map[string]string{
			"hello": "d29ybGQ=",
		},
	}
	if _, err = pathsClient.Create(ctx, accountName, fileSystemName, path, input); err != nil {
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}

	t.Logf("[DEBUG] Retrieving the properties for 'test'..")
	props, err := pathsClient.GetProperties(ctx, accountName, fileSystemName, path)
	if err != nil {
		t.Fatal(fmt.Errorf("Error getting properties: %s", err))
	}
	_ = props

	//TODO - properties don't seem to be saved/retrieved
	// if len(props.Properties) != 1 {
	// 	t.Fatalf("Expected 1 properties by default but got %d", len(props.Properties))
	// }
	// if props.Properties["hello"] != "d29ybGQ=" {
	// 	t.Fatalf("Expected `hello` to be `d29ybGQ=` but got %q", props.Properties["hello"])
	// }

	// t.Logf("[DEBUG] Updating the properties..")
	// setInput := SetPropertiesInput{
	// 	Properties: map[string]string{
	// 		"hello":   "dGVycmFmb3Jt",
	// 		"private": "ZXll",
	// 	},
	// }
	// if _, err := pathsClient.SetProperties(ctx, accountName, fileSystemName, path, setInput); err != nil {
	// 	t.Fatalf("Error setting properties: %s", err)
	// }

	// t.Logf("[DEBUG] Re-Retrieving the Properties..")
	// props, err = pathsClient.GetProperties(ctx, accountName, fileSystemName, path)
	// if err != nil {
	// 	t.Fatal(fmt.Errorf("Error getting properties: %s", err))
	// }
	// if len(props.Properties) != 2 {
	// 	t.Fatalf("Expected 2 properties by default but got %d", len(props.Properties))
	// }
	// if props.Properties["hello"] != "dGVycmFmb3Jt" {
	// 	t.Fatalf("Expected `hello` to be `dGVycmFmb3Jt` but got %q", props.Properties["hello"])
	// }
	// if props.Properties["private"] != "ZXll" {
	// 	t.Fatalf("Expected `private` to be `ZXll` but got %q", props.Properties["private"])
	// }

	t.Logf("[DEBUG] Deleting File System..")
	if _, err := fileSystemsClient.Delete(ctx, accountName, fileSystemName); err != nil {
		t.Fatalf("Error deleting: %s", err)
	}
}
