package ptt

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

func Test_registerCountEmail(t *testing.T) {
	type args struct {
		user  *ptttype.UserecRaw
		email *ptttype.Email_t
	}
	tests := []struct {
		name          string
		args          args
		expectedCount int
		wantErr       bool
	}{
		// TODO: Add test cases.
		{},
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			defer teardownTest()

			gotCount, err := registerCountEmail(tt.args.user, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerCountEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.expectedCount {
				t.Errorf("registerCountEmail() = %v, expected %v", gotCount, tt.expectedCount)
			}
		})
	}
}

func Test_getSystemUaVersion(t *testing.T) {
	tests := []struct {
		name     string
		expected uint8
	}{
		// TODO: Add test cases.
		{expected: 128},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			defer teardownTest()

			if got := getSystemUaVersion(); got != tt.expected {
				t.Errorf("getSystemUaVersion() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestSetupNewUser(t *testing.T) {
	//setup/teardown move to for-loop

	type args struct {
		user *ptttype.UserecRaw
	}
	tests := []struct {
		name     string
		args     args
		expected *ptttype.UserecRaw
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{testSetupNewUser1},
			expected: testSetupNewUser1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			defer teardownTest()

			if err := SetupNewUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("SetupNewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			_, got, err := InitCurrentUser(&tt.args.user.UserID)
			if err != nil {
				t.Errorf("SetupNewUser (InitCurrentUser): err: %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SetupNewUser (InitCurrentUser): got: %v expected: %v", got, tt.expected)
			}
		})
	}
}

func Test_isToCleanUser(t *testing.T) {
	setupTest()
	defer teardownTest()

	file, err := os.OpenFile(ptttype.FN_FRESH, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Warnf("unable to open-file: e: %v", err)
	}
	_, _ = file.Write([]byte("temp"))
	file.Close()

	newTime1 := time.Now().Add(-3700 * types.TS_TO_NANO_TS)
	newTime2 := time.Now().Add(-2700 * types.TS_TO_NANO_TS)

	tests := []struct {
		name     string
		newTime  time.Time
		isDelete bool
		expected bool
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name:     "old-file",
			newTime:  newTime1,
			expected: true,
		},
		{
			name:     "new-file",
			newTime:  newTime2,
			expected: false,
		},
		{
			name:     "no file",
			isDelete: true,
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isDelete {
				os.Remove(ptttype.FN_FRESH)
			} else {
				err = os.Chtimes(ptttype.FN_FRESH, tt.newTime, tt.newTime)
				if err != nil {
					t.Errorf("unable to Chtimes e: %v", err)
				}
			}

			got, err := isToCleanUser()
			if (err != nil) != tt.wantErr {
				t.Errorf("isToCleanUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("isToCleanUser() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func Test_touchFresh(t *testing.T) {
	setupTest()
	defer teardownTest()

	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := touchFresh(); (err != nil) != tt.wantErr {
				t.Errorf("touchFresh() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkAndExpireAccount(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		uid         ptttype.Uid
		user        *ptttype.UserecRaw
		expireRange int
	}
	tests := []struct {
		name     string
		args     args
		expected int
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{uid: 2, user: testUserecRaw2},
			expected: -87067,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkAndExpireAccount(tt.args.uid, tt.args.user, tt.args.expireRange)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkAndExpireAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got > tt.expected {
				t.Errorf("checkAndExpireAccount() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func Test_computeUserExpireValue(t *testing.T) {
	setupTest()
	defer teardownTest()

	user1 := &ptttype.UserecRaw{}
	user1.UserLevel |= ptttype.PERM_XEMPT

	user2 := &ptttype.UserecRaw{}
	copy(user2.UserID[:], ptttype.USER_ID_REGNEW[:])
	user2.LastLogin = types.NowTS() - 10*60

	user3 := &ptttype.UserecRaw{}
	copy(user3.UserID[:], ptttype.USER_ID_REGNEW[:])
	user3.LastLogin = types.NowTS() - 400*60

	type args struct {
		user *ptttype.UserecRaw
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			name:     "PERM_XEMPT",
			args:     args{user1},
			expected: true,
		},
		{
			name:     "new (valid)",
			args:     args{user2},
			expected: true,
		},
		{
			name:     "new (invalid)",
			args:     args{user3},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeUserExpireValue(tt.args.user)
			if got < 0 && tt.expected {
				t.Errorf("computeUserExpireValue() = %v, expected %v", got, tt.expected)
			} else if got >= 0 && !tt.expected {
				t.Errorf("computeUserExpireValue() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestNewRegister(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		userID          *ptttype.UserID_t
		passwd          []byte
		fromHost        *ptttype.IPv4_t
		email           *ptttype.Email_t
		isEmailVerified bool
		isAdbannerUSong bool
		nickname        *ptttype.Nickname_t
		realname        *ptttype.RealName_t
		career          *ptttype.Career_t
		address         *ptttype.Address_t
		over18          bool
	}
	tests := []struct {
		name        string
		args        args
		expected    *ptttype.UserecRaw
		expectedUID ptttype.Uid
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				userID:   &testNewRegister1.UserID,
				passwd:   testNewRegister1Passwd,
				fromHost: &testNewRegister1.LastHost,
				email:    &testNewRegister1.Email,
				nickname: &testNewRegister1.Nickname,
				realname: &testNewRegister1.RealName,
				career:   &testNewRegister1.Career,
				address:  &testNewRegister1.Address,
				over18:   testNewRegister1.Over18,
			},
			expected:    testNewRegister1,
			expectedUID: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUid, got, err := NewRegister(tt.args.userID, tt.args.passwd, tt.args.fromHost, tt.args.email, tt.args.isEmailVerified, tt.args.isAdbannerUSong, tt.args.nickname, tt.args.realname, tt.args.career, tt.args.address, tt.args.over18)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRegister() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			isGood, err := cmbbs.CheckPasswd(got.PasswdHash[:], testNewRegister1Passwd)
			if err != nil || !isGood {
				t.Errorf("NewRegister() unable to checkpasswd: passwd: %v", string(testNewRegister1Passwd))
			}
			copy(testNewRegister1.PasswdHash[:], got.PasswdHash[:])
			testNewRegister1.LastLogin = got.LastLogin
			testNewRegister1.FirstLogin = got.FirstLogin
			testNewRegister1.LastSeen = got.LastSeen

			if gotUid != tt.expectedUID {
				t.Errorf("NewRegister: uid: %v expected: %v", gotUid, tt.expectedUID)
			}

			testutil.TDeepEqual(t, "userec", got, tt.expected)
		})
	}
}

func Test_ensureErasingOldUser(t *testing.T) {
	setupTest()
	defer teardownTest()

	userID2 := &ptttype.UserID_t{}
	copy(userID2[:], []byte("CodingMan"))

	type args struct {
		uid    ptttype.Uid
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{2, userID2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ensureErasingOldUser(tt.args.uid, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("ensureErasingOldUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tryCleanUser(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = os.Remove(ptttype.FN_FRESH)

	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tryCleanUser(); (err != nil) != tt.wantErr {
				t.Errorf("tryCleanUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_registerCheckAndUpdateEmaildb(t *testing.T) {
	email2 := &ptttype.Email_t{}
	copy(email2[:], []byte("test@test.test"))

	type args struct {
		user  *ptttype.UserecRaw
		email *ptttype.Email_t
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{testUserecRaw2, email2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registerCheckAndUpdateEmaildb(tt.args.user, tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("registerCheckAndUpdateEmaildb() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	setupTest()
	defer teardownTest()

	testNewUserID := &ptttype.UserID_t{}
	copy(testNewUserID[:], []byte(ptttype.STR_REGNEW))

	testGuestUserID := &ptttype.UserID_t{}
	copy(testGuestUserID[:], []byte(ptttype.STR_GUEST))

	testReserve0UserID := &ptttype.UserID_t{}
	copy(testReserve0UserID[:], []byte("reserve0"))

	testReserve1UserID := &ptttype.UserID_t{}
	copy(testReserve1UserID[:], []byte("reserve1"))

	type args struct {
		userID          *ptttype.UserID_t
		passwd          []byte
		fromHost        *ptttype.IPv4_t
		email           *ptttype.Email_t
		isEmailVerified bool
		isAdbannerUSong bool
		nickname        *ptttype.Nickname_t
		realname        *ptttype.RealName_t
		career          *ptttype.Career_t
		address         *ptttype.Address_t
		over18          bool
	}
	tests := []struct {
		name         string
		args         args
		expectedUid  ptttype.Uid
		expectedUser *ptttype.UserecRaw
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				userID:   &testNewRegister1.UserID,
				passwd:   testNewRegister1Passwd,
				fromHost: &testNewRegister1.LastHost,
				email:    &testNewRegister1.Email,
				nickname: &testNewRegister1.Nickname,
				realname: &testNewRegister1.RealName,
				career:   &testNewRegister1.Career,
				address:  &testNewRegister1.Address,
				over18:   testNewRegister1.Over18,
			},
			expectedUser: testNewRegister1,
			expectedUid:  6,
		},
		{
			args: args{
				userID: testNewUserID,
				passwd: testNewRegister1Passwd,
				over18: true,
			},
			wantErr: true,
		},
		{
			args: args{
				userID: testGuestUserID,
				passwd: testNewRegister1Passwd,
				over18: true,
			},
			wantErr: true,
		},
		{
			name: "reserve0",
			args: args{
				userID: testReserve0UserID,
				passwd: testNewRegister1Passwd,
				over18: true,
			},
			wantErr: true,
		},
		{
			args: args{
				userID: testReserve1UserID,
				passwd: testNewRegister1Passwd,
				over18: true,
			},
			wantErr: true,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotUid, gotUser, err := Register(tt.args.userID, tt.args.passwd, tt.args.fromHost, tt.args.email, tt.args.isEmailVerified, tt.args.isAdbannerUSong, tt.args.nickname, tt.args.realname, tt.args.career, tt.args.address, tt.args.over18)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}
			if !reflect.DeepEqual(gotUid, tt.expectedUid) {
				t.Errorf("Register() gotUid = %v, want %v", gotUid, tt.expectedUid)
			}

			copy(gotUser.PasswdHash[:], testNewRegister1.PasswdHash[:])
			gotUser.LastLogin = testNewRegister1.LastLogin
			gotUser.FirstLogin = testNewRegister1.FirstLogin
			gotUser.LastSeen = testNewRegister1.LastSeen

			testutil.TDeepEqual(t, "user", gotUser, tt.expectedUser)
		})
		wg.Wait()
	}
}

func TestCheckEmailAllowRejectLists(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "A-pass",
			args: args{email: "abc@gmail.com"},
		},
		{
			name:    "A-fail",
			args:    args{email: "abd@gmail.com"},
			wantErr: true,
		},
		{
			name: "D-pass",
			args: args{email: "test@ptt.test"},
		},
		{
			name: "D-pass",
			args: args{email: "test@sub.ptt.test"},
		},
		{
			name: "P-pass",
			args: args{email: "test@cs.nthu.edu.tw"},
		},
		{
			name: "P-pass",
			args: args{email: "test@cs.nthu.edu.cn"},
		},
		{
			name: "S-pass",
			args: args{email: "test3@ntu.edu.tw"},
		},
		{
			name:    "S-fail",
			args:    args{email: "test@ntu.edu.sg"},
			wantErr: true,
		},
		{
			name:    "S-fail",
			args:    args{email: "test@csie.ntu.edu.tw"},
			wantErr: true,
		},
		{
			name:    "ban-P",
			args:    args{email: "test2@ptt.test"},
			wantErr: true,
		},
		{
			name:    "ban-A",
			args:    args{email: "test@ntu.edu.tw"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckEmailAllowRejectLists(tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("CheckEmailAllowRejectLists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
