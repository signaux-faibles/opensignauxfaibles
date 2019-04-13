package engine

import (
	"fmt"

	daclient "github.com/signaux-faibles/datapi/client"
)

// ExportToDatapi exports Public to Datapi
func ExportToDatapi(url string, user string, password string) error {
	client := daclient.DatapiServer{
		URL: url,
	}

	err := client.Connect(user, password)

	query := daclient.QueryParams{
		Key: map[string]string{},
	}
	objects, err := client.Get("system", query)

	fmt.Println(objects)
	return err
}
