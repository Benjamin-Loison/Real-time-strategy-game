package main

func FindMenuByRef(menus []Menu_t, ref int) Menu_t {
	/*
	for i := 0 ; i < len(menus); i ++ {
		if menus[i].Ref == ref {
			return menus[i]
		}
	}*/
	// The new menu labelling allows us to use the ref as array index.
	if ref < len(menus) {
		return menus[ref]
	} else {
		panic("Could not find the required menu")
	}
}


func FindActionByRef(actions []Action_t, ref int) Action_t {
	for i := 0 ; i < len(actions); i ++ {
		if actions[i].Ref == ref {
			return actions[i]
		}
	}
	panic("Could not find the required action")
}

func ActionTypeOfString(s string) ActionType {
	switch (s) {
		case "Building":
			return ActionBuilding
		case "Quit":
			return ActionQuitGame
		case "MappedKeys":
			return ActionMappedKeys
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

