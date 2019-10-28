package dynamoclient

import (
	pb "github.com/MichiganDiningAPI/api/proto"
)

func (d *DynamoClient) QueryMenus(diningHallName *string, date *string, meal *string) (*[]pb.Menu, error) {
	if diningHallName != nil && date != nil && meal != nil {
		menu := pb.Menu{}
		err := d.GetProto(MenuTableName, map[string]string{DiningHallDateMealKey: *diningHallName + *date + *meal}, &menu)
		if err != nil {
			return nil, err
		}
		menus := []pb.Menu{menu}
		return &menus, nil
	}
	//expression.Name("date").Equals("2019-10-25T00:00:00-04:00")
	return nil, nil
}
