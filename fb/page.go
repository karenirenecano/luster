package fb

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// Page holds access to Facebook page
type Page struct {
	ID      string   // id of page
	Name    string   // name of page
	session *Session // session used to open this page
}

// Fan is user who like of follow page
type Fan struct {
	Profile Profile // user profile
	Time    int32   // time when user liked of started to follow page
	Kind    FanKind // can be "like" or "follow"
}

// FanKind identifies how user become a fan of page
type FanKind int

const (
	// Like represents fan who liked page
	Like FanKind = 1
	// Follow represents fan who follows page
	Follow FanKind = 2
)

// String returns name of FanKind
func (k FanKind) String() string {
	switch k {
	case 1:
		return "Like"
	case 2:
		return "Follow"
	default:
		return "Unknown"
	}
}

// Fans represents collection of fans of some page
type Fans []Fan

// OpenPage opens Facebook page with given name (string identifier)
func (s *Session) OpenPage(name string) (*Page, error) {
	html, err := s.Get(name, nil)
	if err != nil {
		return nil, fmt.Errorf("Cannot open page '%s'", name, err)
	}
	// look for meta tag for mobile applications with page id
	re, _ := regexp.Compile("content=\"fb:\\/\\/page/(\\d+)")
	m := re.FindStringSubmatch(html)
	if len(m) == 0 {
		return nil, fmt.Errorf("Cannot find id for page '%s'", name, err)
	}
	return &Page{ID: m[1], Name: name, session: s}, nil
}

// FetchFans scrapes all users who like or follow this page.
// It looks Facebook limits this information up to 7k users per page.
func (p *Page) FetchFans() (Fans, error) {
	limit := 1000
	var fans []Fan

	// fetch likers
	likers, err := p.fetchFans(false, limit)
	if err != nil {
		return fans, fmt.Errorf("Cannot fetch likers", err)
	}
	fans = append(fans, likers...)

	// fetch followers
	followers, err := p.fetchFans(true, limit)
	if err != nil {
		return fans, fmt.Errorf("Cannot fetch followers", err)
	}
	fans = append(fans, followers...)

	return fans, nil
}

func (p *Page) fetchFans(followers bool, limit int) (Fans, error) {
	var fans Fans
	for offset := 0; true; offset += limit {
		batch, err := p.fetchFansBatch(followers, offset, limit)
		if err != nil {
			return fans, err
		}
		if len(batch) == 0 {
			break // finish when no more fans
		}
		fans = append(fans, batch...)
	}
	return fans, nil
}

func (p *Page) fetchFansBatch(followers bool, offset int, limit int) (Fans, error) {
	// key is input to the endpoint which distinguish between getting likers or followers
	key := "PEOPLE_WHO_LIKE_THIS_PAGE"
	kind := Like
	if followers {
		key = "PEOPLE_WHO_FOLLOW_THIS_PAGE"
		kind = Follow
	}

	// call endpoint with listing people who liked / followe page
	text, err := p.session.Post("pages/admin/people_and_other_pages/entquery/", map[string]string{
		"query_edge_key": key,
		"page_id":        p.ID,
		"offset":         strconv.Itoa(offset),
		"limit":          strconv.Itoa(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot fetch fans batch, page id '%s' limit '%d' offset '%d'", p.ID, offset, limit, err)
	}

	// response is in form 'for (;;);JSON', we need to clean it first
	json := strings.Replace(text, "for (;;);", "", 1)
	json = strings.ReplaceAll(json, "'", "\"")

	// all fans are inside payload.data array
	data := gjson.Get(json, "payload.data").Array()

	// parse each fan which from structure with profile info and timestamp
	fans := make(Fans, len(data))
	for i, d := range data {
		fans[i] = Fan{
			Profile: Profile{
				ID:   gjson.Get(d.Raw, "profile.id").String(),
				Name: gjson.Get(d.Raw, "profile.name").String(),
			},
			Time: int32(gjson.Get(d.Raw, "timestamp").Int()),
			Kind: kind,
		}
	}
	return fans, nil
}
