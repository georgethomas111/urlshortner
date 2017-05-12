package adapter

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gospackler/caddyshack"
	"github.com/gospackler/couchdb"
)

const (
	defaultBufferSize = 50
)

// View Object is placed over here as a query in couch is possible only with a view
// Helps in forming the javascript that can be used for working on stuff.
type ViewObj struct {
	Name     string
	ViewType reflect.Type
}

// FIXME : Assert Kind of the viewObj passed is a pointer.
func NewViewObj(name string, viewObj caddyshack.StoreObject) *ViewObj {
	// Check if View Exists in the DB.
	// Create View Using the tags if thats the case.
	return &ViewObj{
		Name:     name,
		ViewType: reflect.ValueOf(viewObj).Elem().Type(),
	}
}

// Create the View if it does not exist.
func (v *ViewObj) GetCondition() string {
	// Use the tags to form the condition.
	return "Not implemented"
}

// Initial Version
type CouchQuery struct {
	Condition  string // Code for the view RawJson
	ViewName   string
	Store      *CouchStore
	desDoc     *couchdb.DesignDoc
	Params     string
	Skip       int
	Limit      int
	BufferSize int
}

func NewQuery(line string, viewName string, desDoc string, db *CouchStore) (couchQuery *CouchQuery) {
	// Assuming a design doc is already created.
	line = "\"map\" : \"" + line + "\""
	return createCouchQuery(line, viewName, desDoc, db, "")
}

func NewFilterQuery(line string, viewName string, desDoc string, db *CouchStore, params string) (couchQuery *CouchQuery) {
	// Assuming a design doc is already created.
	line = "\"map\" : \"" + line + "\""
	return createCouchQuery(line, viewName, desDoc, db, params)
}

func NewMRQuery(mrCode string, viewName string, desDoc string, db *CouchStore, params string) (couchQuery *CouchQuery) {
	return createCouchQuery(mrCode, viewName, desDoc, db, params)
}

func createCouchQuery(rawCond string, viewName string, desDoc string, db *CouchStore, params string) (couchQuery *CouchQuery) {
	desDocObj := db.GetDesignDoc(desDoc)

	couchQuery = &CouchQuery{
		desDoc:     desDocObj,
		Condition:  rawCond,
		ViewName:   viewName,
		Store:      db,
		Params:     params,
		BufferSize: defaultBufferSize,
	}

	// Correct the code over here.
	newView := &couchdb.View{Name: viewName}
	newView.RawJson = couchQuery.Condition

	index, status := couchQuery.desDoc.CheckExists(viewName)

	log.Debug("Views -->", couchQuery.desDoc.Views)
	log.Debug("LastView -->", couchQuery.desDoc.LastView)
	if status == false {
		couchQuery.desDoc.AddView(newView)

		err := couchQuery.desDoc.SaveDoc()
		if err != nil {
			panic(err)
		}

	} else {
		log.Debug("Index found at ", index)
		if index < 0 {
			couchQuery.desDoc.LastView = newView
		} else {
			couchQuery.desDoc.Views[index] = newView
		}
	}
	return
}

// Use reflection to create the query from the tags.
func NewObjQuery(obj caddyshack.StoreObject, db *CouchStore) (q *CouchQuery) {
	q = new(CouchQuery)
	prefix := "doc"

	viewName := q.GetViewName(obj)
	view := couchdb.NewView(viewName, prefix, q.getCondition(obj, prefix), q.getEmits(obj))
	keyName := q.getByName(obj, prefix)
	if keyName != "" {
		view.KeyName = keyName
	} else {
		view.KeyName = prefix + "._id"
	}
	//Creates the DesignDoc if it does not exist.
	desDoc := db.GetDesignDoc(db.Res.DesDoc)
	index, status := desDoc.CheckExists(viewName)

	log.Debug("Index, Status", index, status)

	if status == false {
		desDoc.AddView(view)
		log.Debug("Added view", desDoc, view)
		// FIXME: Removing object update for now. Saved by the guard.
		err := desDoc.SaveDoc()
		if err != nil {
			panic(err)
		}
	} else {
		if index < 0 {
			desDoc.LastView = view
		} else {
			desDoc.Views[index] = view
		}
	}

	q.desDoc = desDoc
	q.ViewName = viewName
	q.Store = db
	q.BufferSize = defaultBufferSize

	return
}

func (q *CouchQuery) SetCondition(cond string) {
	q.Condition = cond
}

func (q *CouchQuery) GetCondition() string {
	return q.Condition
}

func (q *CouchQuery) GetViewName(obj caddyshack.StoreObject) string {
	structObj := reflect.ValueOf(obj).Elem()
	typeOfObj := structObj.Type()
	return strings.ToLower(typeOfObj.Name())
}

func (q *CouchQuery) getCondition(obj caddyshack.StoreObject, prefix string) (condStr string) {
	structObj := reflect.ValueOf(obj).Elem()
	typeOfObj := structObj.Type()

	firstCond := true

	if structObj.Kind() == reflect.Struct {
		for index := 0; index < typeOfObj.NumField(); index++ {
			structField := typeOfObj.Field(index)
			fieldCond := structField.Tag.Get("condition")
			if fieldCond != "" {
				if firstCond {
					condStr = condStr + prefix + "." + fieldCond
					firstCond = false
				} else {
					condStr = condStr + " && " + prefix + "." + fieldCond
				}
			}
		}
	} else {
		panic(errors.New("Expected a struct pointer as input."))
	}

	return
}
func (q *CouchQuery) getByName(obj caddyshack.StoreObject, prefix string) (byStr string) {
	structObj := reflect.ValueOf(obj).Elem()
	typeOfObj := structObj.Type()

	firstCond := true

	if structObj.Kind() == reflect.Struct {
		for index := 0; index < typeOfObj.NumField(); index++ {
			structField := typeOfObj.Field(index)
			byName := structField.Tag.Get("by")
			if byName != "" {
				if firstCond {
					byStr = byStr + prefix + "." + byName
					firstCond = false
				} else {
					byStr = byStr + " + \\\",\\\" + " + prefix + "." + byName
				}
			}
		}
	} else {
		panic(errors.New("Expected a struct pointer as input."))
	}
	return
}

func (q *CouchQuery) getEmits(obj caddyshack.StoreObject) (emits string) {
	structObj := reflect.ValueOf(obj).Elem()
	typeOfObj := structObj.Type()

	firstCond := true

	if structObj.Kind() == reflect.Struct {
		for index := 0; index < typeOfObj.NumField(); index++ {
			structField := typeOfObj.Field(index)
			jsonStr := structField.Tag.Get("json")
			if jsonStr != "" {
				if firstCond {
					emits = emits + "\\\"" + jsonStr + "\\\""
					firstCond = false
				} else {
					emits = emits + ", \\\"" + jsonStr + "\\\""
				}
			}
		}
	} else {
		panic(errors.New("Expected a struct pointer as input."))
	}

	return
}

func (q *CouchQuery) Execute() (error, []caddyshack.StoreObject) {
	// Currently O(n) w.r.t to views
	var extraParams string
	if q.Limit != 0 {
		extraParams = fmt.Sprintf("&limit=%d&skip=%d", q.Limit, q.Skip)
	}

	params := q.Params + extraParams
	data, err := q.desDoc.Db.GetView(q.desDoc.Id, q.ViewName, params)
	if err != nil {
		return errors.New("Error retreiving view : " + err.Error()), nil
	} else {
		// Print for now create store Object later.
		// FIXME Handle unmarshalling over here.
		result, err := q.Store.MarshalStoreObjects(data)
		if err != nil {
			return errors.New("Could not Marshal json" + err.Error()), result
		}
		return nil, result
	}

	// Move this section to New
	// The intutive version would be creating an object and then adding methods to it.

	// If it exists get the view back.
	// Otherwise Get Retrieve the Data and Marshal the store Object from the json..
}
