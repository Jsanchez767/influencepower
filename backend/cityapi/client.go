package cityapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	BaseURL = "https://webapi.legistar.com/v1/chicago"
)

// Client for Chicago City Clerk ELMS API
type Client struct {
	HTTPClient *http.Client
	BaseURL    string
}

// NewClient creates a new City API client
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		BaseURL: BaseURL,
	}
}

// Matter represents legislation from the City Clerk API
type Matter struct {
	MatterID            int                    `json:"MatterId"`
	MatterGUID          string                 `json:"MatterGuid"`
	MatterFile          string                 `json:"MatterFile"`
	MatterName          string                 `json:"MatterName"`
	MatterTitle         string                 `json:"MatterTitle"`
	MatterTypeID        int                    `json:"MatterTypeId"`
	MatterTypeName      string                 `json:"MatterTypeName"`
	MatterStatusID      int                    `json:"MatterStatusId"`
	MatterStatusName    string                 `json:"MatterStatusName"`
	MatterIntroDate     string                 `json:"MatterIntroDate"`
	MatterAgendaDate    string                 `json:"MatterAgendaDate"`
	MatterPassedDate    string                 `json:"MatterPassedDate"`
	MatterEnactmentDate string                 `json:"MatterEnactmentDate"`
	MatterEnactmentNumber string               `json:"MatterEnactmentNumber"`
	MatterRequester     string                 `json:"MatterRequester"`
	MatterSponsors      []MatterSponsor        `json:"MatterSponsors"`
	MatterAttachments   []MatterAttachment     `json:"MatterAttachments"`
	MatterText          string                 `json:"MatterText"`
	MatterVersion       string                 `json:"MatterVersion"` // API returns string, not int
}

type MatterSponsor struct {
	MatterSponsorID   int    `json:"MatterSponsorId"`
	MatterSponsorGUID string `json:"MatterSponsorGuid"`
	MatterSponsorName string `json:"MatterSponsorName"`
}

type MatterAttachment struct {
	MatterAttachmentID   int    `json:"MatterAttachmentId"`
	MatterAttachmentGUID string `json:"MatterAttachmentGuid"`
	MatterAttachmentName string `json:"MatterAttachmentName"`
	MatterAttachmentHyperlink string `json:"MatterAttachmentHyperlink"`
}

// Event represents council meetings/events
type Event struct {
	EventID          int          `json:"EventId"`
	EventGUID        string       `json:"EventGuid"`
	EventBodyID      int          `json:"EventBodyId"`
	EventBodyName    string       `json:"EventBodyName"`
	EventDate        string       `json:"EventDate"`
	EventTime        string       `json:"EventTime"`
	EventLocation    string       `json:"EventLocation"`
	EventAgendaFile  string       `json:"EventAgendaFile"`
	EventMinutesFile string       `json:"EventMinutesFile"`
	EventVideoURL    string       `json:"EventVideoUrl"`
	EventItems       []EventItem  `json:"EventItems"`
}

type EventItem struct {
	EventItemID          int    `json:"EventItemId"`
	EventItemGUID        string `json:"EventItemGuid"`
	EventItemMatterID    int    `json:"EventItemMatterId"`
	EventItemAgendaSequence int `json:"EventItemAgendaSequence"`
	EventItemAgendaNumber string `json:"EventItemAgendaNumber"`
	EventItemAction      string `json:"EventItemAction"`
	EventItemActionText  string `json:"EventItemActionText"`
}

// Vote represents a vote on legislation
type Vote struct {
	VoteID        int    `json:"VoteId"`
	VoteGUID      string `json:"VoteGuid"`
	VotePersonID  int    `json:"VotePersonId"`
	VotePersonName string `json:"VotePersonName"`
	VoteValue     string `json:"VoteValueName"` // "Yea", "Nay", "Abstain", etc.
	VoteEventID   int    `json:"VoteEventId"`
	VoteMatterID  int    `json:"VoteMatterId"`
	VoteDate      string `json:"VoteDate"`
}

// Body represents committee or council body
type Body struct {
	BodyID        int    `json:"BodyId"`
	BodyGUID      string `json:"BodyGuid"`
	BodyName      string `json:"BodyName"`
	BodyTypeID    int    `json:"BodyTypeId"`
	BodyTypeName  string `json:"BodyTypeName"`
	BodyMeetFlag  int    `json:"BodyMeetFlag"`
}

// Person represents an official from the API
type Person struct {
	PersonID        int    `json:"PersonId"`
	PersonGUID      string `json:"PersonGuid"`
	PersonFirstName string `json:"PersonFirstName"`
	PersonLastName  string `json:"PersonLastName"`
	PersonFullName  string `json:"PersonFullName"`
	PersonEmail     string `json:"PersonEmail"`
	PersonWWW       string `json:"PersonWWW"`
}

// GetMatters fetches legislation/matters from the API
func (c *Client) GetMatters(params map[string]string) ([]Matter, error) {
	endpoint := fmt.Sprintf("%s/matters", c.BaseURL)
	
	// Build query parameters
	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Add(key, value)
	}
	
	if len(queryParams) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())
	}
	
	var matters []Matter
	err := c.doRequest(endpoint, &matters)
	return matters, err
}

// GetMatterByID fetches a specific matter by ID
func (c *Client) GetMatterByID(matterID int) (*Matter, error) {
	endpoint := fmt.Sprintf("%s/matters/%d", c.BaseURL, matterID)
	
	var matter Matter
	err := c.doRequest(endpoint, &matter)
	return &matter, err
}

// GetEvents fetches council events/meetings
func (c *Client) GetEvents(params map[string]string) ([]Event, error) {
	endpoint := fmt.Sprintf("%s/events", c.BaseURL)
	
	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Add(key, value)
	}
	
	if len(queryParams) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())
	}
	
	var events []Event
	err := c.doRequest(endpoint, &events)
	return events, err
}

// GetEventByID fetches a specific event by ID
func (c *Client) GetEventByID(eventID int) (*Event, error) {
	endpoint := fmt.Sprintf("%s/events/%d", c.BaseURL, eventID)
	
	var event Event
	err := c.doRequest(endpoint, &event)
	return &event, err
}

// GetVotes fetches votes for a matter
func (c *Client) GetVotes(matterID int) ([]Vote, error) {
	endpoint := fmt.Sprintf("%s/matters/%d/votes", c.BaseURL, matterID)
	
	var votes []Vote
	err := c.doRequest(endpoint, &votes)
	return votes, err
}

// GetBodies fetches all bodies (committees, council)
func (c *Client) GetBodies() ([]Body, error) {
	endpoint := fmt.Sprintf("%s/bodies", c.BaseURL)
	
	var bodies []Body
	err := c.doRequest(endpoint, &bodies)
	return bodies, err
}

// GetPersons fetches all persons (officials)
func (c *Client) GetPersons() ([]Person, error) {
	endpoint := fmt.Sprintf("%s/persons", c.BaseURL)
	
	var persons []Person
	err := c.doRequest(endpoint, &persons)
	return persons, err
}

// doRequest performs the HTTP request and unmarshals the response
func (c *Client) doRequest(endpoint string, result interface{}) error {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Accept", "application/json")
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	err = json.Unmarshal(body, result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return nil
}
