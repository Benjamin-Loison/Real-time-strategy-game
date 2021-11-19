package main

func FindMenuByRef(menus []Menu_t, ref int) Menu_t {
	for i := 0 ; i < len(menus); i ++ {
		if menus[i].Ref == ref {
			return menus[i]
		}
	}
	panic("Could not find the required menu")
}

func ActionTypeOfString(s string) ActionType {
	switch (s) {
		case "Building":
			return ActionBuilding
		default:
			panic("Unknown action type.")
	}
}

func MenuElementTypeIfString(s string) MenuElementType {
	switch(s) {
		case "Action":
			return MenuElementAction
		default:
			return MenuElementSubMenu
	}
}

