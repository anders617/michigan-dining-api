package mdiningserver

import (
	"context"
	"sync"
	"time"

	"github.com/MichiganDiningAPI/db/dynamoclient"
	"github.com/MichiganDiningAPI/internal/processing/mdiningprocessing"
	"github.com/MichiganDiningAPI/internal/util/date"
	pb "github.com/anders617/mdining-proto/proto/mdining"
	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	dc                *dynamoclient.DynamoClient
	diningHalls       *pb.DiningHalls
	items             *pb.Items
	filterableEntries *pb.FilterableEntries
	foodStats         *[]*pb.FoodStat
	lastFetch         time.Time
	mu                sync.RWMutex
}

func New() *Server {
	s := Server{dc: dynamoclient.New()}
	s.diningHalls = nil
	s.items = nil
	s.filterableEntries = nil
	s.foodStats = nil
	s.fetchData()
	return &s
}

func (s *Server) fetchData() {
	go func() {
		for {
			s.lastFetch = date.Now()
			wg := &sync.WaitGroup{}
			wg.Add(3)
			go s.fetchDiningHalls(wg)
			go s.fetchItemsAndFilterableEntries(wg)
			go s.fetchFoodStats(wg)
			wg.Wait()
			timeToNextFetch := date.NextFetchTime().Sub(date.Now())
			// Reload 30 mins after time that fetch is scheduler for to give it some time
			timeToNextFetch = timeToNextFetch + time.Minute*30
			glog.Infof("Scheduling fetch in %v", timeToNextFetch)
			<-time.After(timeToNextFetch)
		}
	}()
}

func (s *Server) fetchDiningHalls(wg *sync.WaitGroup) {
	defer wg.Done()
	tmp, err := s.dc.QueryDiningHalls()
	if err != nil {
		glog.Fatalf("QueryDiningHalls err %s", err)
	}
	s.mu.Lock()
	s.diningHalls = tmp
	s.mu.Unlock()
	glog.Infof("QueryDiningHalls Success")
}

func (s *Server) fetchItemsAndFilterableEntries(wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	var foods *[]*pb.Food
	// Get all foods after today
	startDate := date.FormatNoTime(date.Now())
	foods, err = s.dc.QueryFoodsDateRange(nil, &startDate, nil)
	if err != nil {
		glog.Fatalf("QueryFoodsDateRange err %s", err)
	}
	glog.Infof("QueryFoodsDateRange Success")
	tmpItems := mdiningprocessing.FoodsToItems(foods)
	s.mu.Lock()
	s.items = tmpItems
	s.mu.Unlock()
	tmpFE := mdiningprocessing.ItemsToFilterableEntries(s.items)
	s.mu.Lock()
	s.filterableEntries = tmpFE
	s.mu.Unlock()
}

func (s *Server) fetchFoodStats(wg *sync.WaitGroup) {
	defer wg.Done()
	tmp, err := s.dc.QueryFoodStats()
	if err != nil {
		glog.Fatalf("QueryFoodStats err %s", err)
	}
	s.mu.Lock()
	s.foodStats = tmp
	s.mu.Unlock()
	glog.Infof("QueryFoodStats Success")
}

//
// Handler for GetDiningHalls request
//
func (s *Server) GetDiningHalls(ctx context.Context, req *pb.DiningHallsRequest) (*pb.DiningHallsReply, error) {
	glog.Infof("GetDiningHalls req{%v}", req)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.diningHalls == nil {
		return nil, status.Error(codes.Unavailable, "Fetching data...")
	}
	// Currently just returns static dining halls data that's checked into git
	return &pb.DiningHallsReply{DiningHalls: s.diningHalls.DiningHalls}, nil
}

//
// Handler for GetItems request
//
func (s *Server) GetItems(ctx context.Context, req *pb.ItemsRequest) (*pb.ItemsReply, error) {
	glog.Infof("GetItems req{%v}", req)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.items == nil {
		return nil, status.Error(codes.Unavailable, "Fetching data...")
	}
	return &pb.ItemsReply{Items: s.items.Items}, nil
}

func (s *Server) GetFilterableEntries(ctx context.Context, req *pb.FilterableEntriesRequest) (*pb.FilterableEntriesReply, error) {
	glog.Infof("GetFilterableEntries req{%v}", req)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.filterableEntries == nil {
		return nil, status.Error(codes.Unavailable, "Fetching data...")
	}
	return &pb.FilterableEntriesReply{FilterableEntries: s.filterableEntries.FilterableEntries}, nil
}

func (s *Server) GetAll(ctx context.Context, req *pb.AllRequest) (*pb.AllReply, error) {
	glog.Infof("GetAll req{%v}", req)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.filterableEntries == nil || s.items == nil || s.diningHalls == nil {
		return nil, status.Error(codes.Unavailable, "Fetching data...")
	}
	reply := pb.AllReply{DiningHalls: s.diningHalls.DiningHalls, Items: s.items.Items, FilterableEntries: s.filterableEntries.FilterableEntries}
	defer glog.Infof("GetAll res{Items: %d, FilterableEntries: %d, DiningHalls: %d}", len(reply.Items), len(reply.FilterableEntries), len(reply.DiningHalls))
	return &reply, nil
}

func (s *Server) GetMenu(ctx context.Context, req *pb.MenuRequest) (*pb.MenuReply, error) {
	glog.Infof("GetMenu req{%v}", req)
	diningHall, date, meal := &req.DiningHall, &req.Date, &req.Meal
	if *diningHall == "" {
		diningHall = nil
	}
	if *date == "" {
		date = nil
	}
	if *meal == "" {
		meal = nil
	}
	menus, err := s.dc.QueryMenus(diningHall, date, meal)
	if err != nil {
		glog.Infof("GetMenu Error %s", err)
		return nil, err
	}
	glog.Infof("GetMenu res{%d menus}", len(*menus))
	return &pb.MenuReply{Menus: *menus}, nil
}

func (s *Server) GetFood(ctx context.Context, req *pb.FoodRequest) (*pb.FoodReply, error) {
	glog.Infof("GetFood req{%v}", req)
	name, date, startDate, endDate := &req.Name, &req.Date, &req.StartDate, &req.EndDate
	if *name == "" {
		name = nil
	}
	if *date == "" {
		date = nil
	}
	if *startDate == "" {
		startDate = nil
	}
	if *endDate == "" {
		endDate = nil
	}
	var foods *[]*pb.Food
	var err error
	if startDate != nil || endDate != nil {
		foods, err = s.dc.QueryFoodsDateRange(name, startDate, endDate)
	} else {
		foods, err = s.dc.QueryFoods(name, date)
	}
	if err != nil {
		glog.Infof("GetFood Error %s", err)
		return nil, err
	}
	glog.Infof("GetFood res{%d foods}", len(*foods))
	return &pb.FoodReply{Foods: *foods}, nil
}

func (s *Server) GetFoodStats(ctx context.Context, req *pb.FoodStatsRequest) (*pb.FoodStatsReply, error) {
	glog.Infof("GetFoodStats req{%v}", req)
	s.mu.RLock()
	if s.foodStats == nil {
		return nil, status.Error(codes.Unavailable, "Fetching data...")
	}
	defer s.mu.RUnlock()
	return &pb.FoodStatsReply{FoodStats: *s.foodStats}, nil
}
