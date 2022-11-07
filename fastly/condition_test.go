package fastly

import (
	"testing"
)

func TestClient_Conditions(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "conditions/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var condition *Condition
	record(t, "conditions/create", func(c *Client) {
		condition, err = c.CreateCondition(&CreateConditionInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           String("test/condition"),
			Statement:      String("req.url~+\"index.html\""),
			Type:           String("REQUEST"),
			Priority:       Int(1),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// // Ensure deleted
	defer func() {
		record(t, "conditions/cleanup", func(c *Client) {
			_ = c.DeleteCondition(&DeleteConditionInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test/condition",
			})
		})
	}()

	if condition.Name != "test/condition" {
		t.Errorf("bad name: %q", condition.Name)
	}
	if condition.Statement != "req.url~+\"index.html\"" {
		t.Errorf("bad statement: %q", condition.Statement)
	}
	if condition.Type != "REQUEST" {
		t.Errorf("bad type: %s", condition.Type)
	}
	if condition.Priority != 1 {
		t.Errorf("bad priority: %d", condition.Priority)
	}

	// List
	var conditions []*Condition
	record(t, "conditions/list", func(c *Client) {
		conditions, err = c.ListConditions(&ListConditionsInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(conditions) < 1 {
		t.Errorf("bad conditions: %v", conditions)
	}

	// Get
	var newCondition *Condition
	record(t, "conditions/get", func(c *Client) {
		newCondition, err = c.GetCondition(&GetConditionInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test/condition",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if condition.Name != newCondition.Name {
		t.Errorf("bad name: %q (%q)", condition.Name, newCondition.Name)
	}
	if condition.Statement != "req.url~+\"index.html\"" {
		t.Errorf("bad statement: %q", condition.Statement)
	}
	if condition.Type != "REQUEST" {
		t.Errorf("bad type: %s", condition.Type)
	}
	if condition.Priority != 1 {
		t.Errorf("bad priority: %d", condition.Priority)
	}

	// Update
	var updatedCondition *Condition
	record(t, "conditions/update", func(c *Client) {
		updatedCondition, err = c.UpdateCondition(&UpdateConditionInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test/condition",
			Statement:      String("req.url~+\"updated.html\""),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if updatedCondition.Statement != "req.url~+\"updated.html\"" {
		t.Errorf("bad statement: %q", updatedCondition.Statement)
	}

	// Delete
	record(t, "conditions/delete", func(c *Client) {
		err = c.DeleteCondition(&DeleteConditionInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test/condition",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListConditions_validation(t *testing.T) {
	var err error
	_, err = testClient.ListConditions(&ListConditionsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListConditions(&ListConditionsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateCondition_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateCondition(&CreateConditionInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateCondition(&CreateConditionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetCondition_validation(t *testing.T) {
	var err error

	_, err = testClient.GetCondition(&GetConditionInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetCondition(&GetConditionInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetCondition(&GetConditionInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateCondition_validation(t *testing.T) {
	var err error

	_, err = testClient.UpdateCondition(&UpdateConditionInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateCondition(&UpdateConditionInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateCondition(&UpdateConditionInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteCondition_validation(t *testing.T) {
	var err error

	err = testClient.DeleteCondition(&DeleteConditionInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteCondition(&DeleteConditionInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteCondition(&DeleteConditionInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}
