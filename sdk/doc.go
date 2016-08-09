// Package SDK offers convenience functions for Go code around Hydra's HTTP APIs.
//
//  import "github.com/ory-am/hydra/sdk"
//  import "github.com/ory-am/hydra/client"
//  var hydra, err = sdk.Connect(
// 	sdk.ClientID("client-id"),
// 	sdk.ClientSecret("client-secret"),
//  	sdk.ClustURL("https://localhost:4444"),
//  )
//
//  // Create a new OAuth2 client
//  var newClient, err = hydra.Client.CreateClient(&client.Client{
//  	ID:                "deadbeef",
//  	Secret:            "sup3rs3cret",
//  	RedirectURIs:      []string{"http://yourapp/callback"},
//  	// ...
//  })
//
//  // Retrieve newly created client
//  var gotClient, err = hydra.Client.GetClient(newClient.ID)
//
//  // Remove the newly created client
//  var err = hydra.Client.DeleteClient(newClient.ID)
//
//  // Retrieve list of all clients
//  var clients, err = hydra.Client.GetClients()
package sdk
