package analyticsclient

import (
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/google/uuid"
)

const analyticsURL = "https://www.google-analytics.com/collect"
const sendHitContentType = "application/x-www-form-urlencoded"

const analyticsVersion = "1"
const analyticsTrackingID = "UA-85646494-2"
const analyticsDataSource = "server"
const analyticsHitType = "pageview"

//
// AnalyticsClient - A type for sending page hits to Google analytics
//
type AnalyticsClient struct {
	// HTTP client used to send requests
	client *http.Client
	// Map from remote IP to UUID to do some very basic user tracking (no persistent accross runs)
	userIDs map[string]string
	// Mutex for syncing userIDs read/writes
	mu sync.Mutex
}

//
// New - Create a new AnalyticsClient object
//
func New() *AnalyticsClient {
	ac := new(AnalyticsClient)
	ac.client = new(http.Client)
	ac.userIDs = make(map[string]string)
	return ac
}

//
// SendHit - Send a hit to Google Analytics
//
// See https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters for more info
//
func (ac *AnalyticsClient) SendHit(r *http.Request) {
	ip := getIPFromRemoteAddr(r.RemoteAddr)
	params := url.Values{}
	params.Set("v", analyticsVersion)
	params.Set("tid", analyticsTrackingID)
	params.Set("ds", analyticsDataSource)
	params.Set("cid", ac.getUserID(ip))
	params.Set("uip", ip)
	params.Set("ua", r.Header.Get("User-Agent"))
	params.Set("dp", r.URL.Path)
	params.Set("t", analyticsHitType)
	payload := params.Encode()
	_, err := ac.client.Post(analyticsURL, sendHitContentType, strings.NewReader(payload))
	if err != nil {
		glog.Warningf("Analytics SendHit Failed: %s", err)
	}
}

func getIPFromRemoteAddr(remoteAddr string) string {
	comps := strings.Split(remoteAddr, ":")
	if len(comps) != 2 {
		// If the address is not parsable, just return localhost
		return "localhost"
	}
	return comps[0]
}

func (ac *AnalyticsClient) getUserID(ip string) string {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	uid, exists := ac.userIDs[ip]
	if !exists {
		uid = uuid.New().String()
		ac.userIDs[ip] = uid
	}
	return uid
}
