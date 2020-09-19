package gitwatch

import (
	"reflect"
	"time"

	"github.com/xaque208/znet/pkg/events"
)

// EventNames are the available named events that this package exports.
var EventNames []string

func init() {
	EventNames = []string{
		"NewCommit",
		"NewTag",
	}
}

// NewCommit is the event when a new commit is detected.
type NewCommit struct {
	Time       *time.Time `json:"time"`
	Name       string     `json:"name"`
	URL        string     `json:"url"`
	Hash       string     `json:"hash"`
	Branch     string     `json:"branch"`
	Collection string     `json:"collection"`
}

func (n *NewCommit) PassFilter(filter interface{}) bool {
	return passGitFilter(filter.(GitFilter), n)
}

// NewTag is the event when a new tag is detected.
type NewTag struct {
	Time       *time.Time `json:"time"`
	Name       string     `json:"name"`
	URL        string     `json:"url"`
	Tag        string     `json:"tag"`
	Collection string     `json:"collection"`
}

func (n *NewTag) PassFilter(filter interface{}) bool {
	return passGitFilter(filter.(GitFilter), n)
}

// GitFilter is a way of reducing when the executions fire, in the case of many
// repos.  This does mean the events are handed to all scubscribers, and the
// subscriber is responsible for reducing the executions.
type GitFilter struct {
	Name       []string
	URL        []string
	Branch     []string
	Collection []string
}

func (f *GitFilter) Filter(ev events.Event) bool {
	// return f.passGitFilter(ev)

	return false
}

func passGitFilter(filter GitFilter, x interface{}) bool {
	var xName string
	var xURL string
	var xBranch string
	var xCollection string

	t := reflect.TypeOf(x).String()

	switch t {
	case "NewTag":
		xName = x.(NewTag).Name
		xURL = x.(NewTag).URL
		xCollection = x.(NewTag).Collection
	case "NewCommit":
		xName = x.(NewCommit).Name
		xURL = x.(NewCommit).URL
		xBranch = x.(NewCommit).Branch
		xCollection = x.(NewCommit).Collection
	default:
		return false
	}

	passName := func() bool {
		if len(filter.Name) == 0 {
			return true
		}

		for _, name := range filter.Name {
			if name == xName {
				return true
			}
		}

		return false
	}

	passURL := func() bool {
		if len(filter.URL) == 0 {
			return true
		}

		for _, url := range filter.URL {
			if url == xURL {
				return true
			}
		}
		return false
	}

	passBranch := func() bool {
		if len(filter.Branch) == 0 {
			return true
		}

		for _, branch := range filter.Branch {
			if branch == xBranch {
				return true
			}
		}
		return false
	}

	passCollection := func() bool {
		if len(filter.Collection) == 0 {
			return true
		}

		for _, collection := range filter.Collection {
			if collection == xCollection {
				return true
			}
		}
		return false
	}

	return passName() && passURL() && passBranch() && passCollection()
}
