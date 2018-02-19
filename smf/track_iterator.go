package smf

import (
	"container/list"
)

//TrackIterator used for iteration over track
type TrackIterator struct {
	trackRef   *Track
	elementRef *list.Element
}

//MoveNext iterate to next element return true if there next element
func (iterator *TrackIterator) MoveNext() bool {

	if iterator.elementRef == nil {

		//check nil front
		if iterator.trackRef.eventsList.Front() == nil {
			return false
		}

		iterator.elementRef = iterator.trackRef.eventsList.Front()
		return true
	}

	//iterate next
	iterator.elementRef = iterator.elementRef.Next()
	return !(iterator.elementRef == nil)
}

//GetValue current value of iterator
func (iterator *TrackIterator) GetValue() Event {

	//check no elements
	if iterator.elementRef == nil {
		return nil
	}

	return iterator.elementRef.Value.(Event)
}

//newTrackIterator return new track iterator
func newTrackIterator(track *Track) *TrackIterator {

	return &TrackIterator{trackRef: track, elementRef: nil}
}
