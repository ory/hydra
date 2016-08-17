// Package SDK offers convenience functions for Go code around Hydra's HTTP APIs.
//
//  import "github.com/ory-am/hydra/sdk"
//  import "github.com/ory-am/hydra/client"
//
//  var hydra, err = sdk.Connect(
// 	sdk.ClientID("client-id"),
// 	sdk.ClientSecret("client-secret"),
//  	sdk.ClusterURL("https://localhost:4444"),
//  )
//
//  // You now have access to the various API endpoints of hydra, for example the oauth2 client endpoint:
//  var newClient, err = hydra.Client.CreateClient(&client.Client{
//  	ID:                "deadbeef",
//  	Secret:            "sup3rs3cret",
//  	RedirectURIs:      []string{"http://yourapp/callback"},
//  	// ...
//  })
//
//  // Retrieve newly created client
//  var gotClient, err = hydra.Client.GetClient(newClient.ID)
package sdk
