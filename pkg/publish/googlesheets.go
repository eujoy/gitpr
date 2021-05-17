package publish

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eujoy/gitpr/internal/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const (
	isLaterAddedTrue  = "Yes"
	isLaterAddedFalse = "No"
)

// GoogleSheetsService describes the google sheets wrapper service.
type GoogleSheetsService struct {
	*sheets.Service
}

// NewGoogleSheetsService creates and returns a new google sheets wrapper service.
//nolint:staticcheck
func NewGoogleSheetsService() (*GoogleSheetsService, error) {
	var b []byte
	var err error

	credentials := os.Getenv("GCP_CREDENTIALS")
	if credentials != "" {
		b = []byte(credentials)
	} else {

		b, err = ioutil.ReadFile("credentials.json")
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
			return nil, err
		}
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
		return nil, err
	}

	return &GoogleSheetsService{srv}, nil
}

// WriteTicketInformationToSpreadsheet write all the ticket information alongside the respective details for the report
// in the sprint specific sheet alongside the respective information to the overall sheet.
// func (s *Service) WriteTicketInformationToSpreadsheet(
// 	spreadsheetID string,
// 	sheetName string,
// 	overallSheetName string,
// 	ticketDetails map[string][]*domain.Ticket,
// 	ticketDetailsMapKeys []string,
// 	reportData *domain.Report,
// ) error {
// 	err := s.CreateAndCleanupSheet(spreadsheetID, sheetName)
// 	if err != nil {
// 		fmt.Printf("Failed to create spreadsheet if it does not exist and clean it up wit error : %v\n", err)
// 		return err
// 	}
//
// 	startingPoint := 1
// 	for _, state := range ticketDetailsMapKeys {
// 		err := s.WriteTicketList(spreadsheetID, sheetName, fmt.Sprintf("A%d", startingPoint), state, ticketDetails[state])
// 		if err != nil {
// 			fmt.Printf("Error writing ticket information in spreadsheet: %v\n", err)
// 			return err
// 		}
//
// 		startingPoint += len(ticketDetails[state]) + 4
// 	}
//
// 	err = s.WriteReportData(spreadsheetID, sheetName, "M1", reportData)
// 	if err != nil {
// 		fmt.Printf("Error writing report data in spreadsheet: %v\n", err)
// 		return err
// 	}
//
// 	err = s.WriteReportDataInOverall(spreadsheetID, overallSheetName, fmt.Sprintf("A%d", reportData.SprintNumber+1), reportData)
// 	if err != nil {
// 		fmt.Printf("Error writing report data in overall data sheet in spreadsheet: %v\n", err)
// 		return err
// 	}
//
// 	return nil
// }

// WriteOverallSheetHeader in the provided overall data sheet. This function shall be used only in case the overall data sheet does not exist.
func (s *GoogleSheetsService) WriteOverallSheetHeader(spreadsheetID string, sheetName string) error {
	var vr sheets.ValueRange

	rangeInSheet := fmt.Sprintf("%v!A1", sheetName)

	listOfValues := [][]interface{}{
		{
			"#",
			"Sprint Name",
			"Start Date",
			"End Date",
			"Comments Total",
			"Comments Avg",
			"Review Comments Total",
			"Review Comments Avg",
			"Commits Total",
			"Commits Avg",
			"Additions Total",
			"Additions Avg",
			"Deletions Total",
			"Deletions Avg",
			"Modified Lines Total",
			"Modified Lines Avg",
			"Changed Files Total",
			"Changed Files Avg",
			"Lead Time Total",
			"Lead Time Avg",
			"Time to Merge Total",
			"Time to Merge Avg",
			"Flow Rate Created",
			"Flow Rate Merged",
			"Flow Ratio",
		},
	}

	vr.Values = listOfValues

	_, err := s.Spreadsheets.Values.Update(spreadsheetID, rangeInSheet, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		fmt.Printf("Failed to write data in spreadsheet with error: %v", err)
		return err
	}

	return nil
}

// WriteReportData to the provided spreadsheet.
func (s *GoogleSheetsService) WritePullRequestReportData(spreadsheetID string, sheetName string, cellRange string, sprint *domain.SprintSummary, prMetrics *domain.PullRequestMetrics, prFlowRatio *domain.PullRequestFlowRatio) error {
	var vr sheets.ValueRange

	rangeInSheet := fmt.Sprintf("%v!%v", sheetName, cellRange)

	listOfValues := [][]interface{}{
		{
			sprint.Number,
			sprint.Name,
			sprint.StartDate.Format("2006-01-02"),
			sprint.EndDate.Format("2006-01-02"),
			prMetrics.Total.Comments,
			convertToFloatWithTwoDecimals(prMetrics.Average.Comments),
			prMetrics.Total.ReviewComments,
			convertToFloatWithTwoDecimals(prMetrics.Average.ReviewComments),
			prMetrics.Total.Commits,
			convertToFloatWithTwoDecimals(prMetrics.Average.Commits),
			prMetrics.Total.Additions,
			convertToFloatWithTwoDecimals(prMetrics.Average.Additions),
			prMetrics.Total.Deletions,
			convertToFloatWithTwoDecimals(prMetrics.Average.Deletions),
			prMetrics.Total.Additions + prMetrics.Total.Deletions,
			convertToFloatWithTwoDecimals(prMetrics.Average.Additions + prMetrics.Average.Deletions),
			prMetrics.Total.ChangedFiles,
			convertToFloatWithTwoDecimals(prMetrics.Average.ChangedFiles),
			convertDurationToHourDecimal(prMetrics.Total.LeadTime),
			convertDurationToHourDecimal(prMetrics.Average.LeadTime),
			convertDurationToHourDecimal(prMetrics.Total.TimeToMerge),
			convertDurationToHourDecimal(prMetrics.Average.TimeToMerge),
			prFlowRatio.Created,
			prFlowRatio.Merged,
			prFlowRatio.Ratio,
		},
	}

	vr.Values = listOfValues

	_, err := s.Spreadsheets.Values.Update(spreadsheetID, rangeInSheet, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		fmt.Printf("Failed to write data in spreadsheet with error: %v", err)
		return err
	}

	return nil
}

// CreateAndCleanupOverallSheet checks if the overall sheet exists, create it if it doesn't exist and add the respective header line.
func (s *GoogleSheetsService) CreateAndCleanupOverallSheet(spreadsheetID string, sheetName string) error {
	sheetExists, err := s.CheckIfSheetExists(spreadsheetID, sheetName)
	if err != nil {
		fmt.Printf("Error checking if sheet exists in spreadsheet: %v\n", err)
		return err
	}

	if sheetExists {
		return nil
	}

	err = s.CreateSheet(spreadsheetID, sheetName)
	if err != nil {
		fmt.Printf("Error creating the sheet %q in spreadsheet: %v\n", sheetName, err)
		return err
	}

	err = s.WriteOverallSheetHeader(spreadsheetID, sheetName)
	if err != nil {
		fmt.Printf("Error to write the header line in overall data sheet %q in spreadsheet: %v\n", sheetName, err)
		return err
	}

	return nil
}

// CreateSheet in an existing spreadsheet with a given name.
func (s *GoogleSheetsService) CreateSheet(spreadsheetID string, sheetName string) error {
	req := sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: sheetName,
			},
		},
	}

	rbb := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{&req},
	}

	_, err := s.Spreadsheets.BatchUpdate(spreadsheetID, rbb).Do()
	if err != nil {
		fmt.Printf("Failed to create new sheet : %v", err)
		return err
	}

	return nil
}

// CheckIfSheetExists in a provided spreadsheet.
func (s *GoogleSheetsService) CheckIfSheetExists(spreadsheetID string, sheetName string) (bool, error) {
	_, err := s.Spreadsheets.Values.Get(spreadsheetID, fmt.Sprintf("%v!A1", sheetName)).Do()
	if err != nil {
		if strings.Contains(err.Error(), "Unable to parse range") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}

	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	tok := &oauth2.Token{}
	var err error

	token := os.Getenv("GCP_AUTH_TOKEN")
	if token != "" {
		err = json.Unmarshal([]byte(token), tok)
	} else {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}

		defer func() {
			err := f.Close()
			if err != nil {
				log.Fatalf("Failed to close the file with error : %v", err)
			}
		}()

		err = json.NewDecoder(f).Decode(tok)
	}

	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatalf("Failed to close the file with error : %v", err)
		}
	}()

	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Fatalf("Failed to encode token with error : %v", err)
	}
}

// -- Helper Functions --

func convertToFloatWithTwoDecimals(inputVal float64) float64 {
	outputVal, err := strconv.ParseFloat(fmt.Sprintf("%.2f", inputVal), 64)
	if err != nil {
		return 0.0
	}

	return outputVal
}

func convertDurationToHourDecimal(dur time.Duration) float64 {
	if dur == time.Duration(0) {
		return 0.0
	}

	days := int64(dur.Hours() / 24)
	hours := int64(math.Mod(dur.Hours(), 24))
	minutes := int64(math.Mod(dur.Minutes(), 60))

	formattedDuration := fmt.Sprintf("%02d.%02d", days*24 + hours, minutes)

	floatDuration, err := strconv.ParseFloat(formattedDuration, 64)
	if err != nil {
		return 0.0
	}

	return floatDuration
}
