package y_crdt

import "testing"

func TestDecodeRelativePositionTags(t *testing.T) {
	itemID := GenID(1, 2)
	rp0 := &RelativePosition{Item: &itemID, Assoc: 0}
	dec0 := DecodeRelativePosition(EncodeRelativePosition(rp0))
	if dec0.Item == nil || dec0.Item.Client != 1 || dec0.Item.Clock != 2 {
		t.Fatalf("case 0: got item %#v", dec0.Item)
	}

	rp1 := &RelativePosition{Tname: "XmlText", Assoc: -1}
	dec1 := DecodeRelativePosition(EncodeRelativePosition(rp1))
	if dec1.Tname != "XmlText" || dec1.Assoc != -1 {
		t.Fatalf("case 1: got tname=%q assoc=%v", dec1.Tname, dec1.Assoc)
	}

	typeID := GenID(3, 4)
	rp2 := &RelativePosition{Type: &typeID, Assoc: 1}
	dec2 := DecodeRelativePosition(EncodeRelativePosition(rp2))
	if dec2.Type == nil || dec2.Type.Client != 3 || dec2.Type.Clock != 4 {
		t.Fatalf("case 2: got type %#v", dec2.Type)
	}
}
