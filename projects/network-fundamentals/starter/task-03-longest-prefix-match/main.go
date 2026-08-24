package main

// Route is a single routing table entry.
type Route struct {
	CIDR    string
	NextHop string
}

// Router is a simple longest-prefix-match routing table.
type Router struct {
	routes []Route
}

// NewRouter returns an empty Router.
func NewRouter() *Router {
	return &Router{}
}

// AddRoute appends a route to the table.
func (r *Router) AddRoute(cidr, nextHop string) {
	// TODO: append a Route{CIDR: cidr, NextHop: nextHop} to r.routes.
}

// Lookup finds the next hop for ip using longest-prefix-match.
func (r *Router) Lookup(ip string) (nextHop string, found bool) {
	// TODO: parse ip with net.ParseIP.
	// TODO: for every route, parse its CIDR with net.ParseCIDR and check
	// ipnet.Contains(parsedIP).
	// TODO: among matching routes, keep the one with the largest prefix
	// length (ones, _ := ipnet.Mask.Size()) — that's the most specific match.
	// TODO: return found=false if nothing matched.
	return "", false
}

func main() {}
