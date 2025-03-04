package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"
)

type Establishment struct {
	Name           string `xml:"BusinessName"`
	RatingValue    string `xml:"RatingValue"`
	Address        string `xml:"AddressLine1"`
	LocalAuthority string `xml:"LocalAuthorityName"`
}

type Establishments struct {
	Establishments []Establishment `xml:"EstablishmentCollection>EstablishmentDetail"`
}

// aica food safety app
//
// This program prompts the user to enter a postcode and then fetches food safety ratings
// for cafes in that postcode area from the Food Standards Agency API. The results are
// displayed in a tabular format including the company name, rating value, address, and
// local authority.
//
// The program performs the following steps:
//  1. Prompts the user to enter a postcode.
//  2. Constructs the API URL using the entered postcode.
//  3. Sends an HTTP GET request to the API.
//  4. Reads and unmarshals the XML response into a Go struct.
//  5. Prints the number of responses received.
//  6. Displays the results in a tabular format with columns for company name, rating value,
//     address, and local authority.
func main() {
	var postcode string
	fmt.Print("Enter postcode: ")
	fmt.Scanln(&postcode)

	url := fmt.Sprintf("https://api1-ratings.food.gov.uk/search/cafe/%s/xml", postcode)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	var establishments Establishments
	err = xml.Unmarshal(body, &establishments)
	if err != nil {
		fmt.Println("Error unmarshalling XML:", err)
		os.Exit(1)
	}

	fmt.Printf("Number of responses: %d\n", len(establishments.Establishments))

	// Create a new tab writer
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(writer, "Company\tRating\tAddress\tLocal Authority")
	for _, establishment := range establishments.Establishments {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", establishment.Name, establishment.RatingValue, establishment.Address, establishment.LocalAuthority)
	}
	writer.Flush()
}
