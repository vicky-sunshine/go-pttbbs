package cache

import (
	"reflect"
	"sync"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestGetBCache(t *testing.T) {
	setupTest()
	defer teardownTest()

	boards := [3]ptttype.BoardHeaderRaw{testBoardHeader0, testBoardHeader1, testBoardHeader2}
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.BCache),
		unsafe.Sizeof(boards),
		unsafe.Pointer(&boards),
	)

	type args struct {
		bidInCache ptttype.Bid
	}
	tests := []struct {
		name          string
		args          args
		expectedBoard *ptttype.BoardHeaderRaw
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			args:          args{1},
			expectedBoard: &testBoardHeader0,
		},
		{
			args:          args{2},
			expectedBoard: &testBoardHeader1,
		},
		{
			args:          args{3},
			expectedBoard: &testBoardHeader2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBoard, err := GetBCache(tt.args.bidInCache)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBoard, tt.expectedBoard) {
				t.Errorf("GetBCache() = %v, want %v", gotBoard, tt.expectedBoard)
			}
		})
	}
}

func TestIsHiddenBoardFriend(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = LoadUHash()

	ReloadBCache()

	type args struct {
		bidInCache ptttype.BidInStore
		uidInCache ptttype.UidInStore
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			args:     args{0, 0}, //board: SYSOP user: SYSOP
			expected: true,
		},
		{
			args:     args{0, 1}, //board: SYSOP user: CodingMan
			expected: false,
		},
		{
			args:     args{0, 2}, //board: SYSOP user: pichu
			expected: false,
		},
		{
			args:     args{0, 3}, //board: SYSOP user: Kahou
			expected: true,
		},
		{
			args:     args{0, 4}, //board: SYSOP user: Kahou2
			expected: false,
		},
		{
			args:     args{0, 5}, //board: SYSOP user: (non-exist)
			expected: false,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := IsHiddenBoardFriend(tt.args.bidInCache, tt.args.uidInCache); got != tt.expected {
				t.Errorf("IsHiddenBoardFriend() = %v, want %v", got, tt.expected)
			}
		})
		wg.Wait()
	}
}

func TestNumBoards(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	tests := []struct {
		name     string
		expected int32
	}{
		// TODO: Add test cases.
		{
			expected: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumBoards(); got != tt.expected {
				t.Errorf("NumBoards() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestReloadBCache(t *testing.T) {
	setupTest()
	defer teardownTest()

	tests := []struct {
		name                  string
		expectedNBoard        int32
		expectedSortedByName  []ptttype.BidInStore
		expectedSortedByClass []ptttype.BidInStore
		expectedBCacheName    []ptttype.BoardID_t
		expectedBCacheTitle   []ptttype.BoardTitle_t
	}{
		// TODO: Add test cases.
		{
			expectedNBoard:        int32(12),
			expectedSortedByName:  testSortedByName,
			expectedSortedByClass: testSortedByClass,
			expectedBCacheName:    testBCacheName,
			expectedBCacheTitle:   testBCacheTitle,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReloadBCache()

			nBoard := int32(0)
			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.BNumber),
				types.INT32_SZ,
				unsafe.Pointer(&nBoard),
			)

			if !reflect.DeepEqual(nBoard, tt.expectedNBoard) {
				t.Errorf("ReloadBCache() = %v, want %v", nBoard, tt.expectedNBoard)
			}

			bsorted := [ptttype.BSORT_BY_MAX][ptttype.MAX_BOARD]ptttype.BidInStore{}
			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.BSorted),
				unsafe.Sizeof(bsorted),
				unsafe.Pointer(&bsorted),
			)

			for idx := int32(0); idx < nBoard; idx++ {
				board, _ := GetBCache(ptttype.Bid(idx + 1))
				if types.Cstrcmp(board.Brdname[:], tt.expectedBCacheName[idx][:]) != 0 {
					t.Errorf("bcacheName (%v/%v) = %v, want %v", idx, nBoard, string(board.Brdname[:]), string(tt.expectedBCacheName[idx][:]))
				}
				if types.Cstrcmp(board.Title[:], tt.expectedBCacheTitle[idx][:]) != 0 {
					t.Errorf("bcacheTitle (%v/%v) = %v, want %v", idx, nBoard, string(board.Title[:]), string(tt.expectedBCacheTitle[idx][:]))
				}
			}

			bsortedByName := bsorted[ptttype.BSORT_BY_NAME][:nBoard]
			bsortedByClass := bsorted[ptttype.BSORT_BY_CLASS][:nBoard]
			if !reflect.DeepEqual(bsortedByName, tt.expectedSortedByName) {
				t.Errorf("bsorted-by-name = %v, want %v", bsortedByName, tt.expectedSortedByName)
			}
			if !reflect.DeepEqual(bsortedByClass, tt.expectedSortedByClass) {
				t.Errorf("bsorted-by-class = %v, want %v", bsortedByClass, tt.expectedSortedByClass)
			}

		})
	}
}

func Test_reloadCacheLoadBottom(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	tests := []struct {
		name     string
		expected uint8
	}{
		// TODO: Add test cases.
		{
			expected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reloadCacheLoadBottom()
			nBottom := uint8(0)
			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.NBottom)+unsafe.Sizeof(nBottom)*9,
				unsafe.Sizeof(nBottom),
				unsafe.Pointer(&nBottom),
			)

			if nBottom != tt.expected {
				t.Errorf("nBottom: %v want: %v", nBottom, tt.expected)
			}
		})
	}
}

func TestGetBTotal(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	bid1 := ptttype.Bid(3)
	bid1InCache := bid1.ToBidInStore()
	total1 := int32(5)

	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.Total)+types.INT32_SZ*uintptr(bid1InCache),
		types.INT32_SZ,
		unsafe.Pointer(&total1),
	)

	type args struct {
		bid ptttype.Bid
	}
	tests := []struct {
		name          string
		args          args
		expectedTotal int32
	}{
		// TODO: Add test cases.
		{
			args:          args{bid1},
			expectedTotal: total1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTotal := GetBTotal(tt.args.bid); gotTotal != tt.expectedTotal {
				t.Errorf("GetBTotal() = %v, want %v", gotTotal, tt.expectedTotal)
			}
		})
	}
}

func TestSetBTotal(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	type args struct {
		bid ptttype.Bid
	}
	tests := []struct {
		name                 string
		args                 args
		wantErr              bool
		expectedTotal        int32
		expectedLastPostTime types.Time4
	}{
		// TODO: Add test cases.
		{
			args:                 args{10},
			expectedTotal:        2,
			expectedLastPostTime: 1607203395,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetBTotal(tt.args.bid); (err != nil) != tt.wantErr {
				t.Errorf("SetBTotal() error = %v, wantErr %v", err, tt.wantErr)
			}

			total := GetBTotal(tt.args.bid)
			if total != tt.expectedTotal {
				t.Errorf("SetBTotal: total: %v want: %v", total, tt.expectedTotal)
			}

			bidInCache := tt.args.bid.ToBidInStore()
			lastPostTime := types.Time4(0)
			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.LastPostTime)+types.TIME4_SZ*uintptr(bidInCache),
				types.TIME4_SZ,
				unsafe.Pointer(&lastPostTime),
			)
			if lastPostTime != tt.expectedLastPostTime {
				t.Errorf("SetBTotal: lastPostTime: %v want: %v", lastPostTime, tt.expectedLastPostTime)
			}
		})
	}
}

func TestSetBottomTotal(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	type args struct {
		bid ptttype.Bid
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		expectedTotal uint8
	}{
		// TODO: Add test cases.
		{
			args:          args{10},
			expectedTotal: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetBottomTotal(tt.args.bid); (err != nil) != tt.wantErr {
				t.Errorf("SetBottomTotal() error = %v, wantErr %v", err, tt.wantErr)
			}

			bidInCache := tt.args.bid.ToBidInStore()
			total := uint8(0)
			const uint8sz = unsafe.Sizeof(total)

			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.NBottom)+uint8sz*uintptr(bidInCache),
				uint8sz,
				unsafe.Pointer(&total),
			)
			if total != tt.expectedTotal {
				t.Errorf("SetBottomTotal: total: %v want: %v", total, tt.expectedTotal)
			}

		})
	}
}

func TestGetBid(t *testing.T) {
	setupTest()
	defer teardownTest()

	ReloadBCache()

	boardID0 := &ptttype.BoardID_t{}
	copy(boardID0[:], []byte("WhoAmI"))

	boardID1 := &ptttype.BoardID_t{}
	copy(boardID1[:], []byte("SYSOP"))

	type args struct {
		boardID *ptttype.BoardID_t
	}
	tests := []struct {
		name        string
		args        args
		expectedBid ptttype.Bid
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			args:        args{boardID: boardID0},
			expectedBid: 10,
		},
		{
			args:        args{boardID: boardID1},
			expectedBid: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBid, err := GetBid(tt.args.boardID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBid, tt.expectedBid) {
				t.Errorf("GetBid() = %v, want %v", gotBid, tt.expectedBid)
			}
		})
	}
}
