package msg

import (
	"math"
	"reflect"
	"testing"

	"github.com/grafana/metrictank/schema"
)

func TestWriteReadPointMsg(t *testing.T) {
	mp := schema.MetricPoint{
		MKey: schema.MKey{
			Org: 123,
		},
		Time:  math.MaxUint32,
		Value: 123.45,
	}
	buf := make([]byte, 0, 33)
	out, err := WritePointMsg(mp, buf, FormatMetricPoint)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	_, ok := IsPointMsg(out)
	if !ok {
		t.Fatal("IsPointMsg: exp true, got false")
	}

	leftover, outPoint, err := ReadPointMsg(out, 6)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if len(leftover) > 0 {
		t.Fatalf("expected no leftover. got %v", leftover)
	}

	if !reflect.DeepEqual(mp, outPoint) {
		t.Fatalf("expected point %v, got %v", mp, outPoint)
	}
}

func TestWriteReadPointMsgWithoutOrg(t *testing.T) {
	mp := schema.MetricPoint{
		MKey: schema.MKey{
			Org: 123,
		},
		Time:  math.MaxUint32,
		Value: 123.45,
	}
	buf := make([]byte, 0, 29)
	out, err := WritePointMsg(mp, buf, FormatMetricPointWithoutOrg)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	_, ok := IsPointMsg(out)
	if !ok {
		t.Fatal("IsPointMsg: exp true, got false")
	}

	exp := mp
	exp.MKey.Org = 6 // ReadPointMsg will have to set the default org and we want to check it
	leftover, outPoint, err := ReadPointMsg(out, 6)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if len(leftover) > 0 {
		t.Fatalf("expected no leftover. got %v", leftover)
	}

	if !reflect.DeepEqual(exp, outPoint) {
		t.Fatalf("expected point %v, got %v", exp, outPoint)
	}
}

func TestWriteReadIndexControlMsg(t *testing.T) {
	mp := schema.ControlMsg{
		Op:   schema.OpArchive,
		Defs: make([]schema.MetricDefinition, 1),
	}
	out, err := WriteIndexControlMsg(&mp)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	ok := IsIndexControlMsg(out)
	if !ok {
		t.Fatal("IsPointMsg: exp true, got false")
	}

	outMsg, err := ReadIndexControlMsg(out)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if !reflect.DeepEqual(mp, outMsg) {
		t.Fatalf("expected point %v, got %v", mp, outMsg)
	}
}
