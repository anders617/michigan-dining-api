package mdiningclient

import (
	"net/http"
	"net/url"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	pb "github.com/MichiganDiningAPI/api/proto"
	"github.com/MichiganDiningAPI/api/proto/mdiningapi"
)

const (
	DiningHallGroupName = "DINING HALLS"
)

var (
	DiningHallListUrl            = urlFrom("https://mobile.its.umich.edu/michigan/services/dining/shallowDiningHallGroups")
	DiningHallMenuBaseUrl        = urlFrom("https://mobile.its.umich.edu/michigan/services/dining/shallowMenusByDiningHall")
	DiningHallMenuDetailsBaseUrl = urlFrom("https://mobile.its.umich.edu/michigan/services/dining/menusByDiningHall")
)

// Construct URL without having second return value
func urlFrom(s string) *url.URL {
	u, _ := url.Parse(s)
	return u
}

type MDiningClient struct {
	client *http.Client
}

func New() *MDiningClient {
	mc := new(MDiningClient)
	mc.client = new(http.Client)
	return mc
}

func (m *MDiningClient) getPB(url string, reply proto.Message) error {
	res, err := m.client.Get(url)
	if err != nil {
		glog.Error("Network error: %s", err)
		return err
	}
	defer res.Body.Close()
	um := jsonpb.Unmarshaler{AllowUnknownFields: true}
	um.Unmarshal(res.Body, reply)
	if err != nil {
		glog.Error("Error unmarshalling json: %s", err)
		return err
	}
	return nil
}

func (m *MDiningClient) GetMenuDetails(diningHall *pb.DiningHall, base *mdiningapi.MenuBase) (*mdiningapi.GetMenuDetailsReply, error) {
	date, _ := time.Parse(time.RFC3339, base.Date)
	params := make(url.Values)
	params.Add("_type", "json")
	params.Add("diningHall", diningHall.Name)
	params.Add("menu", base.Name)
	params.Add("date", date.Format("02-01-06"))
	url := DiningHallMenuDetailsBaseUrl
	url.RawQuery = params.Encode()
	reply := mdiningapi.GetMenuDetailsReply{}
	err := m.getPB(url.String(), &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func (m *MDiningClient) GetMenuBase(diningHall *pb.DiningHall) (*mdiningapi.GetMenuBaseReply, error) {
	params := make(url.Values)
	params.Add("_type", "json")
	params.Add("diningHall", diningHall.Name)
	url := DiningHallMenuBaseUrl
	url.RawQuery = params.Encode()
	reply := mdiningapi.GetMenuBaseReply{}
	err := m.getPB(url.String(), &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func (m *MDiningClient) GetDiningHallList() (*pb.DiningHalls, error) {
	params := make(url.Values)
	params.Add("_type", "json")
	url := DiningHallListUrl
	url.RawQuery = params.Encode()
	reply := mdiningapi.GetDiningHallsReply{}
	err := m.getPB(url.String(), &reply)
	if err != nil {
		return nil, err
	}
	diningHalls := pb.DiningHalls{}
	for _, group := range reply.DiningHallGroup {
		if group.Name == DiningHallGroupName {
			diningHalls.DiningHalls = group.DiningHall
			break
		}
	}
	return &diningHalls, nil
}
