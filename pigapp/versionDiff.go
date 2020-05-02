package main

import (
	"fmt"
	"net/http"
	"sort"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"
)

// VersionDiffFacade is the Facade for the Diff page.
type VersionDiffFacade struct {
	PageMeta                           PageMeta
	UserProfile                        UserProfile
	IdiomLeft, IdiomRight              *IdiomHistory
	ImplIDs                            []int
	ImplLeft, ImplRight                map[int]Impl
	CreationImplIDs                    map[int]bool
	DeletionImplIDs                    map[int]bool
	PreviousChangePath, NextChangePath string
}

func versionDiff(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	ctx := r.Context()

	idiomIDStr := vars["idiomId"]
	idiomID := String2Int(idiomIDStr)
	v1Str := vars["v1"]
	v1 := String2Int(v1Str)
	v2Str := vars["v2"]
	v2 := String2Int(v2Str)
	if v2 < v1 {
		return PiError{fmt.Sprintf("Won't compare v%v with older v%v", v1, v2), http.StatusBadRequest}
	}
	if v2 == v1 {
		return PiError{fmt.Sprintf("Won't compare v%v with itself", v1), http.StatusBadRequest}
	}

	var err error
	var left *IdiomHistory
	if v1 == 0 {
		// Dummy empty object, to show an idiom creation
		left = &IdiomHistory{}
		left.Idiom.Id = idiomID
	} else {
		_, left, err = dao.getIdiomHistory(ctx, idiomID, v1)
		if err != nil {
			return PiError{err.Error(), http.StatusNotFound}
		}
	}
	_, right, err := dao.getIdiomHistory(ctx, idiomID, v2)
	if err != nil {
		return PiError{err.Error(), http.StatusNotFound}
	}

	removeUntouchedImpl(left, right)

	implIDs := make([]int, 0, len(right.Implementations)+1)
	implLeft := map[int]Impl{}
	implRight := map[int]Impl{}
	creationImplIDs := map[int]bool{}
	deletionImplIDs := map[int]bool{}
	for _, impl := range left.Implementations {
		implIDs = append(implIDs, impl.Id)
		implLeft[impl.Id] = impl
	}
	for _, impl := range right.Implementations {
		if _, ok := implLeft[impl.Id]; !ok {
			implIDs = append(implIDs, impl.Id)
			creationImplIDs[impl.Id] = true
		}
		implRight[impl.Id] = impl
	}
	for _, impl := range left.Implementations {
		if _, ok := implRight[impl.Id]; !ok {
			deletionImplIDs[impl.Id] = true
		}
	}
	// Recently created... first?
	sort.Sort(sort.Reverse(sort.IntSlice(implIDs)))

	userProfile := readUserProfile(r)
	myToggles := copyToggles(toggles)
	myToggles["writable"] = false
	myToggles["actionEditIdiom"] = false
	myToggles["actionIdiomHistory"] = false
	myToggles["actionAddImpl"] = false
	// Note that custom toggles don't work
	// in deeply nested templates...
	data := &VersionDiffFacade{
		PageMeta: PageMeta{
			PageTitle:             right.Title,
			Toggles:               myToggles,
			PreventIndexingRobots: true,
		},
		UserProfile:     userProfile,
		IdiomLeft:       left,
		IdiomRight:      right,
		ImplIDs:         implIDs,
		ImplLeft:        implLeft,
		ImplRight:       implRight,
		CreationImplIDs: creationImplIDs,
		DeletionImplIDs: deletionImplIDs,
	}
	// Note: the Prev/Next links wouldn't work in a case where version numbers
	// wouldn't be perfectly sequential.
	if left.Version >= 2 {
		data.PreviousChangePath = fmt.Sprintf("/idiom/%d/diff/%d/%d", left.Id, left.Version-1, left.Version)
	}
	_, _, errNext := dao.getIdiomHistory(ctx, right.Id, right.Version+1)
	if errNext == nil {
		data.NextChangePath = fmt.Sprintf("/idiom/%d/diff/%d/%d", right.Id, right.Version, right.Version+1)
	}
	return templates.ExecuteTemplate(w, "page-idiom-version-diff", data)
}

// removeUntouchedImpl strips all non-relevant implementations from diff operands
func removeUntouchedImpl(a, b *IdiomHistory) {
	// two maps ImplID -> version
	mapa := make(map[int]int, len(a.Implementations))
	for _, impl := range a.Implementations {
		mapa[impl.Id] = impl.Version
	}
	mapb := make(map[int]int, len(b.Implementations))
	for _, impl := range b.Implementations {
		mapb[impl.Id] = impl.Version
	}
	// if same version, remove from both sides
	touchedA := make([]Impl, 0, len(a.Implementations))
	for _, impl := range a.Implementations {
		// Keep if only in a, or in b with different version
		if mapa[impl.Id] != mapb[impl.Id] {
			touchedA = append(touchedA, impl)
		}
	}
	a.Implementations = touchedA
	touchedB := make([]Impl, 0, len(b.Implementations))
	for _, impl := range b.Implementations {
		// Keep if only in b, or in a with different version
		if mapa[impl.Id] != mapb[impl.Id] {
			touchedB = append(touchedB, impl)
		}
	}
	b.Implementations = touchedB

	// also, the two sides should have same impl order (except from creation/deletion)
}
