package mdiningclient2

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	"github.com/MichiganDiningAPI/internal/util/date"
	pb "github.com/anders617/mdining-proto/proto/mdining"
	"github.com/anders617/mdining-proto/proto/mdiningapi2"
)

const (
	DiningHallGroupName = "DINING HALLS"
)

var (
	GetLocationsUrl = urlFrom("https://prod-dining-services.webplatformsunpublished.umich.edu/dining/locations")
	GetMenuUrl      = urlFrom("https://prod-dining-services.webplatformsunpublished.umich.edu/dining/menu")
	GetMealHoursUrl = urlFrom("https://prod-dining-services.webplatformsunpublished.umich.edu/dining/meal-hours")
)

// Construct URL without having second return value
func urlFrom(s string) *url.URL {
	u, _ := url.Parse(s)
	return u
}

type MDiningClient2 struct {
	client *http.Client
	apiKey string
}

func New(apiKey string) *MDiningClient2 {
	mc := new(MDiningClient2)
	mc.client = new(http.Client)
	mc.apiKey = apiKey
	return mc
}

func (m *MDiningClient2) getPB(url string, reply proto.Message, preprocess func(string) string) error {
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

func (m *MDiningClient2) GetAllMenus(partialMenus *[]*pb.Menu) (*[]*pb.Menu, error) {
	var wg sync.WaitGroup
	for _, partialMenu := range *partialMenus {
		wg.Add(1)
		go func(partialMenu *pb.Menu) {
			defer wg.Done()
			if !partialMenu.HasCategories {
				return
			}
			dateTime, err := date.ParseNoTime(&partialMenu.Date)
			if err != nil {
				glog.Errorf("Error parsing date (%s)", partialMenu.Date)
				return
			}
			reply, err := m.GetMenu(partialMenu.DiningHallName, partialMenu.Meal, dateTime)
			if err != nil {
				glog.Errorf("Error retrieving menu (%s, %s, %s)", partialMenu.DiningHallName, partialMenu.Meal, partialMenu.Date)
				return
			}
			if reply.Menu == nil {
				glog.Errorf("No menu in reply (%s, %s, %s)", partialMenu.DiningHallName, partialMenu.Meal, partialMenu.Date)
				return
			}
			partialMenu.Category = []*pb.Category{}
			for _, category := range reply.Menu.Category {
				newCategory := pb.Category{Name: category.Name, MenuItem: []*pb.MenuItem{}}
				for _, menuItem := range category.MenuItem {
					newMenuItem := pb.MenuItem{
						Name:      menuItem.Name,
						Attribute: menuItem.Attribute,
						Allergens: menuItem.Allergens,
						ItemSizes: []*pb.ItemSizes{},
					}
					itemSize := menuItem.ItemSizes
					portionSize, err := strconv.Atoi(itemSize.GetPortionSize())
					if err != nil {
						portionSize = 0
					}
					newItemSize := pb.ItemSizes{
						PortionSize:     int32(portionSize),
						ServingSize:     itemSize.ServingSize,
						NutritionalInfo: []*pb.NutritionalInfo{},
					}
					for _, nutrition := range itemSize.Nutrition {
						value := 0
						unit := ""
						for i := 1; i < len(nutrition.Value); i++ {
							v, err := strconv.Atoi(nutrition.Value[0:i])
							if err != nil {
								unit = nutrition.Value[i:]
								break
							}
							value = v
						}
						newItemSize.NutritionalInfo = append(newItemSize.NutritionalInfo, &pb.NutritionalInfo{
							Name:              nutrition.Name,
							Value:             int32(value),
							Units:             unit,
							PercentDailyValue: nutrition.PercentDailyValue,
						})
					}
					newMenuItem.ItemSizes = append(newMenuItem.ItemSizes, &newItemSize)
					newCategory.MenuItem = append(newCategory.MenuItem, &newMenuItem)
				}
				partialMenu.Category = append(partialMenu.Category, &newCategory)
			}
		}(partialMenu)
	}
	wg.Wait()
	return partialMenus, nil
}

func (m *MDiningClient2) GetMenu(location string, meal string, d time.Time) (*mdiningapi2.GetMenuReply, error) {
	params := make(url.Values)
	params.Add("key", m.apiKey)
	params.Add("location", location)
	params.Add("meal", meal)
	params.Add("date", date.FormatMDiningAPINoTime(d))
	url := *GetMenuUrl
	url.RawQuery = params.Encode()
	reply := mdiningapi2.GetMenuReply{}
	glog.Infof("GetMenu %s", url.String())
	preprocess := func(s string) string { return s }
	err := m.getPB(url.String(), &reply, preprocess)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

func (m *MDiningClient2) GetMealHours(location string, d time.Time) (*mdiningapi2.GetMealHoursReply, error) {
	params := make(url.Values)
	params.Add("key", m.apiKey)
	params.Add("location", location)
	params.Add("date", date.FormatMDiningAPINoTime(d))
	url := *GetMealHoursUrl
	url.RawQuery = params.Encode()
	glog.Infof("GetMealHours %s", url.String())
	preprocess := func(s string) string { return s }
	reply := &mdiningapi2.GetMealHoursReply{}
	err := m.getPB(url.String(), reply, preprocess)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (m *MDiningClient2) GetDiningHallList(dates []time.Time) (*map[string]*pb.DiningHalls, *[]*pb.Menu, error) {
	params := make(url.Values)
	params.Add("key", m.apiKey)
	url := *GetLocationsUrl
	url.RawQuery = params.Encode()
	reply := mdiningapi2.GetLocationsReply{}
	glog.Infof("API KEY: \"%s\"", m.apiKey)
	glog.Infof("GetDiningHallList %s", url.String())
	// Wrap in object so that it is convertible to PB
	preprocess := func(s string) string { return strings.Join([]string{"{\"location\":", s, "}"}, "") }
	err := m.getPB(url.String(), &reply, preprocess)
	if err != nil {
		return nil, nil, err
	}
	diningHallsByCampus := make(map[string]*pb.DiningHalls)
	diningHalls := make([]*pb.DiningHall, 0)
	for _, location := range reply.Location {
		if _, exists := diningHallsByCampus[location.Campus]; !exists {
			diningHallsByCampus[location.Campus] = &pb.DiningHalls{DiningHalls: []*pb.DiningHall{}}
		}
		dh := pb.DiningHall{
			Name:   location.Name,
			Campus: location.Campus,
			Building: &pb.DiningHall_Building{
				Name:    location.Buildingpreferredname,
				Address: location.Address,
			},
			Type:         location.Type,
			SortPosition: 0,
		}
		diningHalls = append(diningHalls, &dh)
	}
	menus := make([][]*pb.Menu, len(diningHalls))
	var wg sync.WaitGroup
	for lidx, location := range diningHalls {
		wg.Add(1)
		go func(diningHallIdx int, locationName string, locationCampus string) {
			defer wg.Done()
			for _, d := range dates {
				reply, err := m.GetMealHours(locationName, d)
				if err != nil {
					continue
				}
				dayEvents := &pb.DiningHall_DayEvent{
					Key:           date.Format(d),
					CalendarEvent: []*pb.DiningHall_DayEvent_CalendarEvent{},
				}
				for _, hour := range reply.Hours {
					dayEvents.CalendarEvent = append(dayEvents.CalendarEvent, &pb.DiningHall_DayEvent_CalendarEvent{
						EventDayEnd:    hour.EventDayEnd,
						EventDayStart:  hour.EventDayStart,
						EventTimeStart: hour.EventTimeStart,
						EventTimeEnd:   hour.EventTimeEnd,
						EventTitle:     hour.EventTitle,
					})
				}
				for _, meal := range reply.Meal {
					menus[diningHallIdx] = append(menus[diningHallIdx], &pb.Menu{
						Meal:             meal.Name,
						HasCategories:    meal.HasMenu,
						Description:      meal.Description,
						DiningHallName:   locationName,
						DiningHallCampus: locationCampus,
						DiningHallMeal:   locationName + meal.Name,
						FormattedDate:    date.FormatNoTime(d),
						Date:             date.FormatNoTime(d),
					})
				}
				diningHalls[diningHallIdx].DayEvents = append(diningHalls[diningHallIdx].DayEvents, dayEvents)
			}
		}(lidx, location.Name, location.Campus)
	}
	wg.Wait()

	for _, diningHall := range diningHalls {
		diningHallsByCampus[diningHall.Campus].DiningHalls = append(diningHallsByCampus[diningHall.Campus].DiningHalls, diningHall)
	}

	menuList := []*pb.Menu{}
	for _, dhMenus := range menus {
		menuList = append(menuList, dhMenus...)
	}

	return &diningHallsByCampus, &menuList, nil
}
