//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package database

import (
	"fmt"
	"testing"
	"time"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

func TestStoragePackage(t *testing.T) { TestingT(t) }

type DataStoreTests struct {
	Store *DefaultDatastore
}

// SetUpSuite is called once before the tests start running.
func (s *DataStoreTests) SetUpSuite(c *C) {
	//s.Store = NewEphemeralDatastore()
}

// SetUpTest is called before each test.
func (s *DataStoreTests) SetUpTest(c *C) {
	// create a new, temporary  datastore
	// s.Store = NewEphemeralDatastore()
	s.Store = NewDefaultDatastore("127.0.0.1", "hsutest")

	// wipe out the collections
	s.Store.DB.C("surveys").RemoveAll(bson.M{})
	s.Store.DB.C("results").RemoveAll(bson.M{})
}

// TearDownSuite is called once after the tests have finished running.
func (s *DataStoreTests) TearDownSuite(c *C) {
}

// TearDownTest is called after each test.
func (s *DataStoreTests) TearDownTest(c *C) {
	s.Store.Close()
}

var _ = Suite(&DataStoreTests{})

// ----------------------------------------------------------------------------
//
func (s *DataStoreTests) TestClose(c *C) {
	s.Store.Close()

	// this call should fail
	_, err := s.Store.GetSurvey("666")
	c.Assert(err, NotNil)
}

// ----------------------------------------------------------------------------
// TestSurveys
//
func (s *DataStoreTests) TestSurveys(c *C) {
	// retrieving a non-existing survey should result in an error
	_, err := s.Store.GetSurvey("unknown")
	c.Assert(err, NotNil)

	s1 := SurveyObj{
		ID:     "testsurvey",
		Name:   "Test Survey",
		Online: true,
	}

	s.Store.AddSurvey(&s1)

	// adding a surevey with the same id twice will result in an error
	err = s.Store.AddSurvey(&s1)
	c.Assert(err, NotNil)

	// there should be one survey in the list
	surveyList, _ := s.Store.ListSurveys()
	c.Assert(len(surveyList), Equals, 1)

	// the returned survey should be equal to local one
	s2, e1 := s.Store.GetSurvey(s1.ID)
	c.Assert(e1, IsNil)

	c.Assert(s2.Name, Equals, s1.Name)

	// we can also retrieve the survey by it's "ReferenceID"
	// s2, _ = s.Store.GetSurvey(s1.ShortID)
	// c.Assert(s2.Name, Equals, s1.Name)

	sUnknown := SurveyObj{
		ID:     "unknown",
		Name:   "Test Survey",
		Online: true,
	}
	// updating an unknown survey should result in an error
	err = s.Store.UpdateSurvey(&sUnknown)
	c.Assert(err, NotNil)

	// updating the name of an existing survey
	s1.Name = "Test Survey (updated)"
	s.Store.UpdateSurvey(&s1)

	// the returned survey should be equal to local one
	s2, _ = s.Store.GetSurvey(s1.ID)
	c.Assert(s2.Name, Equals, s1.Name)
}

// ----------------------------------------------------------------------------
// TestSurveyResults
//
func (s *DataStoreTests) TestSurveyResults(c *C) {

	r1 := ResultObj{
		FBNR:      "buffalo_soldier",
		SurveyID:  "unknown",
		StartTime: time.Now().Local().String(),
		EndTime:   time.Now().Local().String(),
		Datapoints: []DatapointObj{
			DatapointObj{
				ScreenID: "1",
				Position: 1,
				Result:   "L",
				Timing:   666,
			},
		},
	}

	// adding a result to a non-existing survey should result in an error
	err := s.Store.AddSurveyResult(&r1)
	c.Assert(err, NotNil)

	// adding the same result to an exiting survey should be successful
	s1 := SurveyObj{
		ID:     "testsurvey",
		Name:   "Test Survey",
		Online: true,
	}
	s.Store.AddSurvey(&s1)
	c.Assert(s1.ID, NotNil)

	r1.SurveyID = s1.ID
	err = s.Store.AddSurveyResult(&r1)
	c.Assert(err, IsNil)
	c.Assert(r1.ID, NotNil)

	// now list results should yield the result
	results, err := s.Store.ListSurveyResults(s1.ID)
	fmt.Printf("%v", results)
	c.Assert(err, IsNil)
	c.Assert(len(results), Equals, 1)

	// if we add another result, under the same accountid, we should
	// end up with an error
	err = s.Store.AddSurveyResult(&r1)
	c.Assert(err, NotNil)

	// changing the username should allow us to add the result
	r1.FBNR = "dreadlock_rasta"
	err = s.Store.AddSurveyResult(&r1)

	results, err = s.Store.ListSurveyResults(s1.ID)
	c.Assert(err, IsNil)
	c.Assert(len(results), Equals, 2)

}
