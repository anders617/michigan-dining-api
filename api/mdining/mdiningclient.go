package mdiningclient

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	pb "github.com/MichiganDiningAPI/api/proto"
	"github.com/MichiganDiningAPI/api/proto/mdiningapi"
	"github.com/MichiganDiningAPI/internal/util/date"
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
	b, err1 := ioutil.ReadAll(res.Body)
	if err1 != nil {
		return err1
	}
	s := string(b)
	um := jsonpb.Unmarshaler{AllowUnknownFields: true}
	// um.Unmarshal(res.Body, reply)
	// Sometimes mdining returns empty string instead of 0
	// This messes with jsonpb unmarshalling since it expects an int
	s = strings.ReplaceAll(s, "portionSize\":\"\"", "portionSize\":0")
	s = strings.ReplaceAll(s, "postalCode\":\"\"", "postalCode\":0")
	err = um.Unmarshal(strings.NewReader(s), reply)
	if err != nil {
		glog.Errorf("Error unmarshalling json: %s", err)
		return err
	}
	return nil
}

func (m *MDiningClient) GetAllMenus(diningHalls *pb.DiningHalls) (*[]*pb.Menu, error) {
	var wg sync.WaitGroup
	diningHallMenus := make([]*[]*pb.Menu, len(diningHalls.DiningHalls))
	for idx, diningHall := range diningHalls.DiningHalls {
		wg.Add(1)
		go func(idx int, diningHall *pb.DiningHall) {
			defer wg.Done()
			menu, err := m.GetMenus(diningHall)
			if err != nil {
				glog.Warningf("Error getting %s menus %s", diningHall.Name, err)
				diningHallMenus[idx] = nil
				return
			}
			diningHallMenus[idx] = menu
		}(idx, diningHall)
	}
	wg.Wait()
	menus := make([]*pb.Menu, 0)
	for _, menu := range diningHallMenus {
		if menu == nil {
			continue
		}
		menus = append(menus, *menu...)
	}
	return &menus, nil
}

func (m *MDiningClient) GetMenus(diningHall *pb.DiningHall) (*[]*pb.Menu, error) {
	reply, err := m.GetMenuDetails(diningHall)
	if err != nil {
		return nil, err
	}
	menus := make([]*pb.Menu, 0)
	glog.Infof("Parsing menus for %s", diningHall.Name)
	for _, m := range reply.Menu {
		if m == nil {
			// TODO: Why nil?
			continue
		}
		dateTime, err := date.Parse(&m.Date)
		if err != nil {
			glog.Warningf("Could not parse date %s", m.Date)
			continue
		}
		dateNoTime := date.FormatNoTime(dateTime)
		menu := pb.Menu{
			DiningHallMeal: diningHall.Name + m.Name,
			Meal:           m.Name,
			Date:           dateNoTime,
			FormattedDate:  m.FormattedDate,
			RatingCount:    m.RatingCount,
			RatingScore:    m.RatingScore,
			HasCategories:  m.HasCategories,
			Description:    m.Description,
			Category:       m.Category,
			DiningHallName: diningHall.Name}
		menus = append(menus, &menu)
	}
	return &menus, nil
}

func (m *MDiningClient) GetMenuDetails(diningHall *pb.DiningHall) (*mdiningapi.GetMenuDetailsReply, error) {
	params := make(url.Values)
	params.Add("_type", "json")
	params.Add("diningHall", diningHall.Name)
	url := *DiningHallMenuDetailsBaseUrl
	url.RawQuery = params.Encode()
	reply := mdiningapi.GetMenuDetailsReply{}
	glog.Infof("GetMenuDetails %s %s", diningHall.Name, url)
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
	glog.Infof("GetMenuBase %s %s", diningHall.Name, url)
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
	glog.Infof("GetDiningHallList %s", url)
	err := m.getPB(url.String(), &reply)
	if err != nil {
		// Don't return err here. There are multiple "diningHallGroup"
		// objects that have different structures than the one we want.
		// This causes the pb unmarshaller to return an error even
		// when the target diningHallGroup is processed
		// Should probably fix the parsing
		// return nil, err
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
