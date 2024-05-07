package main

import (
	"context"
	// "encoding/json" 
    "fmt"
	"os"
	"github.com/joho/godotenv"
	plaid "github.com/plaid/plaid-go/v21/plaid"
)

var CLIENT_ID string
var CLIENT_SECRET string

var transactionData []Trans

var cursor *string
var added []Trans
var modified []Trans
var removed []string

type Trans struct {
    Name string                 `json: "name"`
    Logo string                 `json: "logo"`
    Cost float64                `json: "cost"`
    DateTime string             `json: "dateTime"`
    // Loc plaid.Location          `json: "location"`
}

func loadEnv() {
    err := godotenv.Load(".env")
    if err != nil {
        fmt.Println(err)
    }
    CLIENT_ID = os.Getenv("CLIENT_ID")
    CLIENT_SECRET = os.Getenv("CLIENT_SECRET")
}

func configClient() *plaid.APIClient {
    configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", CLIENT_ID)
	configuration.AddDefaultHeader("PLAID-SECRET", CLIENT_SECRET)
	configuration.UseEnvironment(plaid.Sandbox)
    return plaid.NewAPIClient(configuration)
}

func addTransactionWorker(accessToken string, client *plaid.APIClient) {
    envAccessToken := os.Getenv(accessToken)
    // TODO: find a way to store cursor for this item id 

    ctx := context.Background()
	hasMore := true

	// Iterate through each page of new transaction updates for item
	for hasMore {
		request := plaid.NewTransactionsSyncRequest(envAccessToken)
		if cursor != nil {
			request.SetCursor(*cursor)
        }

		resp, _, _:= client.PlaidApi.TransactionsSync(
			ctx,
		).TransactionsSyncRequest(*request).Execute()

		// Add this page of results
        appendTransactions(&added, resp.GetAdded()...)
        appendTransactions(&modified, resp.GetModified()...)
        for i := 0; i < len(resp.GetRemoved()); i++ {
            remTrans := resp.GetRemoved()[i]
            removed = append(removed, remTrans.GetTransactionId())
        }

		hasMore = resp.GetHasMore()

		// Update cursor to the next cursor
		nextCursor := resp.GetNextCursor()
		cursor = &nextCursor
	}
}

// Helper function for converting the plaid.Transaction data into data that I can turn to JSON
func appendTransactions(transArr *[]Trans, newTransactions ...plaid.Transaction) {
    for i := 0; i < len(newTransactions); i++ {
        t := newTransactions[i] 

        var logo string;
        if t.HasLogoUrl() {
            logo = *t.LogoUrl.Get()
        } else {
            logo = ""
        }

        var transTime string
        if t.Datetime.IsSet() {
            tempDate := *t.Datetime.Get()
            transTime = tempDate.String()
        } else {
            transTime = t.Date
        }

        newTrans := Trans{t.Name, logo, t.Amount, transTime}
        *transArr = append(*transArr, newTrans)
    }
}

func main() {
    loadEnv()
    client := configClient()
    addTransactionWorker("AMEX_ACCESS_TOKEN", client)

    // res, err := json.Marshal(added)
    // if err != nil {
    //     fmt.Println(err)
    //     return
    // }
    // fmt.Println(string(res))
}
