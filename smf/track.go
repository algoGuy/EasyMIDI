package smf

import (
	"container/list"
)

//Track midi track struct
type Track struct {

	//EventsList linked list of midi events
	eventsList list.List
}

//AddEvent adds events to EventsList
func (t *Track) AddEvent(event Event) error {

	//check nil
	if event == nil {
		return &MidiError{"Try to add nil event"}
	}

	//add event
	t.eventsList.PushBack(event)

	return nil
}

//RemoveAt removes event at position
func (t *Track) RemoveAt(position uint32) error {

	//check position
	if position >= uint32(t.eventsList.Len()) {
		return &MidiError{"Event number out of range"}
	}

	//remove value
	for i, e := uint32(0), t.eventsList.Front(); e != nil; e = e.Next() {
		if i == position {
			t.eventsList.Remove(e)
		}
	}

	return nil
}

//Len return number of events
func (t *Track) Len() int {

	return t.eventsList.Len()
}

//GetAllEvents return array of all events in track
func (t *Track) GetAllEvents() []Event {

	//create result
	result := make([]Event, t.eventsList.Len())

	//add values
	for i, e := 0, t.eventsList.Front(); e != nil; e = e.Next() {
		result[i] = e.Value.(Event)
		i++
	}

	return result
}

//GetIterator return iterator for track
func (t *Track) GetIterator() *TrackIterator {

	return newTrackIterator(t)
}

//TrackFromArray create track from events array
func TrackFromArray(events []Event) (*Track, error) {

	//check nil
	if events == nil {
		return nil, &MidiError{"nil events array"}
	}

	//create track
	track := &Track{}

	//add events
	for _, val := range events {

		err := track.AddEvent(val)
		if err != nil {
			return nil, err
		}
	}

	return track, nil
}
