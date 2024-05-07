package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"github.com/joho/godotenv"
	plaid "github.com/plaid/plaid-go/v21/plaid"
)

var CLIENT_ID string
var CLIENT_SECRET string

var transactionData []Trans

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
    ctx := context.Background()

    const iso8601TimeFormat = "2006-01-02"
    startDate := time.Now().Add(-365 * 24 * time.Hour).Format(iso8601TimeFormat)
    endDate := time.Now().Format(iso8601TimeFormat)

    request := plaid.NewTransactionsGetRequest(
      envAccessToken,
      startDate,
      endDate,
    )

    transactionsResp, _, err := client.PlaidApi.TransactionsGet(ctx).TransactionsGetRequest(*request).Execute()
    if err != nil {
        fmt.Println(err)
        return
    }

    t := transactionsResp.Transactions[15]

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

    // var location plaid.Location
    // var newTrans Trans
    // if t.LogoUrl.IsSet() {
    //     newTrans = Trans{t.Name, *t.LogoUrl.Get(), t.Amount, t.Datetime, t.Location}
    // } else {
    //     newTrans = Trans{t.Name, "", t.Amount, t.Datetime, t.Location}
    // }

    newTrans := Trans{t.Name, logo, t.Amount, transTime}
    transactionData = append(transactionData, newTrans)
}

func main() {
    loadEnv()
    client := configClient()
    addTransactionWorker("AMEX_ACCESS_TOKEN", client)

    res, err := json.Marshal(transactionData)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(res))
}
