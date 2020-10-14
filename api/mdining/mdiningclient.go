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

	"github.com/MichiganDiningAPI/internal/util/date"
	pb "github.com/anders617/mdining-proto/proto/mdining"
	"github.com/anders617/mdining-proto/proto/mdiningapi"
)

const (
	DiningHallGroupName = "DINING HALLS"
)

// https://prod-dining-services.webplatformsunpublished.umich.edu/
/**
https://prod-dining-services.webplatformsunpublished.umich.edu/dining/menu?key=093665d6ab069c859267fd4001c3c562ba805539ed852978&location=Bursley%20Dining%20Hall&date=11-10-2020&meal=BREAKFAST
https://prod-dining-services.webplatformsunpublished.umich.edu/dining/menu?key=093665d6ab069c859267fd4001c3c562ba805539ed852978&location=Bursley%20Dining%20Hall&date=11-10-2020&meal=BREAKFAST
https://prod-dining-services.webplatformsunpublished.umich.edu/dining/locations?key=093665d6ab069c859267fd4001c3c562ba805539ed852978
https://prod-dining-services.webplatformsunpublished.umich.edu/dining/meal-hours?key=093665d6ab069c859267fd4001c3c562ba805539ed852978&location=Java%20Blu%20At%20Sph&date=11-10-2020
https://prod-dining-services.webplatformsunpublished.umich.edu/dining/capacity?key=093665d6ab069c859267fd4001c3c562ba805539ed852978

*/
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

func (m *MDiningClient) getPB(url string, reply proto.Message, preprocess func(string) string) error {
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
	s = preprocess(s)
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
			DiningHallMeal:   diningHall.Name + m.Name,
			Meal:             m.Name,
			Date:             dateNoTime,
			FormattedDate:    m.FormattedDate,
			RatingCount:      m.RatingCount,
			RatingScore:      m.RatingScore,
			HasCategories:    m.HasCategories,
			Description:      m.Description,
			Category:         m.Category,
			DiningHallName:   diningHall.Name,
			DiningHallCampus: diningHall.Campus}
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
	glog.Infof("GetMenuDetails %s %s", diningHall.Name, url.String())
	preprocess := func(s string) string {
		// Sometimes mdining returns empty string instead of 0
		// This messes with jsonpb unmarshalling since it expects an int
		s = strings.ReplaceAll(s, "portionSize\":\"\"", "portionSize\":0")
		return s
	}
	err := m.getPB(url.String(), &reply, preprocess)
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
	preprocess := func(s string) string { return s }
	err := m.getPB(url.String(), &reply, preprocess)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func (m *MDiningClient) GetDiningHallList() (*map[string]*pb.DiningHalls, error) {
	params := make(url.Values)
	params.Add("_type", "json")
	url := DiningHallListUrl
	url.RawQuery = params.Encode()
	reply := mdiningapi.GetDiningHallsReply{}
	glog.Infof("GetDiningHallList %s", url)
	preprocess := func(s string) string {
		s = strings.ReplaceAll(s, "postalCode\":\"\"", "postalCode\":0")
		// One of the dining hall campuses is named using an int?????????????Wtf
		// Convert the int to a string so protobuf doesn't complain
		s = strings.ReplaceAll(s, "{\"name\":1265", "{\"name\":\"1265\"")
		s = strings.ReplaceAll(s, "\"campus\":1265", "\"campus\":\"1265\"")
		return s
	}
	err := m.getPB(url.String(), &reply, preprocess)
	if err != nil {
		return nil, err
	}
	diningHallsByCampus := make(map[string]*pb.DiningHalls)
	for _, group := range reply.DiningHallGroup {
		diningHalls := &pb.DiningHalls{DiningHalls: []*pb.DiningHall{}}
		for _, diningHall := range group.DiningHall {
			if diningHall.Name != "" {
				diningHalls.DiningHalls = append(diningHalls.DiningHalls, diningHall)
			}
		}
		diningHallsByCampus[group.Name] = diningHalls
	}
	return &diningHallsByCampus, nil
}
