package wechat

import "testing"

func TestAddCardTicket(t *testing.T) {

}

func TestListdCardTicket(t *testing.T) {
	err := testWx.CardBatchget(nil, 0, 100)
	if err != nil {
		t.Error(err)
		return
	}

}
