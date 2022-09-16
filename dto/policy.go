package dto

import (
	"github.com/casbin/casbin/v2"
)

type Data struct {
	Route  string
	Method string
}

func RoutePolicy(e *casbin.Enforcer, role, route, method string) (bool, error) {

	bol := e.HasNamedPolicy("p", role, route, method)
	if !bol {
		ok, err := e.AddPolicy(role, route, method)
		if !ok || err != nil {
			return false, err
		}
	}

	return true, nil

}

func Plocies() []*Data {
	menu := []*Data{

		//ping
		{
			Route:  "/api/admin/ping",
			Method: "GET",
		},
		//admin
		{
			Route:  "/api/admin/list",
			Method: "GET",
		},
		{
			Route:  "/api/admin/create",
			Method: "POST",
		},
		{
			Route:  "/api/admin/delete",
			Method: "POST",
		},
		{
			Route:  "/api/admin/login",
			Method: "POST",
		},
		{
			Route:  "/api/admin/updatePassword",
			Method: "GET",
		},

		//type1
		{
			Route:  "/api/type1_table/list",
			Method: "GET",
		},
		{
			Route:  "/api/type1_table/create",
			Method: "POST",
		},
		{
			Route:  "/api/type1_table/update",
			Method: "POST",
		},
		{
			Route:  "/api/type1_table/delete",
			Method: "POST",
		},

		//type2
		{
			Route:  "/api/type2_table/list",
			Method: "GET",
		},
		{
			Route:  "/api/type2_table/create",
			Method: "POST",
		},
		{
			Route:  "/api/type2_table/update",
			Method: "POST",
		},
		{
			Route:  "/api/type2_table/delete",
			Method: "POST",
		},

		//type3
		{
			Route:  "/api/type3_table/list",
			Method: "GET",
		},
		{
			Route:  "/api/type3_table/create",
			Method: "POST",
		},
		{
			Route:  "/api/type3_table/update",
			Method: "POST",
		},
		{
			Route:  "/api/type3_table/delete",
			Method: "POST",
		},

		//type4
		{
			Route:  "/api/type4_table/list",
			Method: "GET",
		},
		{
			Route:  "/api/type4_table/create",
			Method: "POST",
		},
		{
			Route:  "/api/type4_table/update",
			Method: "POST",
		},
		{
			Route:  "/api/type4_table/delete",
			Method: "POST",
		},

		//type5
		{
			Route:  "/api/type5_table/list",
			Method: "GET",
		},
		{
			Route:  "/api/type5_table/create",
			Method: "POST",
		},
		{
			Route:  "/api/type5_table/update",
			Method: "POST",
		},
		{
			Route:  "/api/type5_table/delete",
			Method: "POST",
		},

		//type6
		{
			Route:  "/api/type6_table/list",
			Method: "GET",
		},
		{
			Route:  "/api/type6_table/create",
			Method: "POST",
		},
		{
			Route:  "/api/type6_table/update",
			Method: "POST",
		},
		{
			Route:  "/api/type6_table/delete",
			Method: "POST",
		},

		//type7
		{
			Route:  "/api/type7_table/list",
			Method: "GET",
		},
		{
			Route:  "/api/type7_table/create",
			Method: "POST",
		},
		{
			Route:  "/api/type7_table/update",
			Method: "POST",
		},
		{
			Route:  "/api/type7_table/delete",
			Method: "POST",
		},
	}
	return menu
}
