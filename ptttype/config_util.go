package ptttype

import (
	"bufio"
	"os"
	"regexp"

	"github.com/Ptt-official-app/go-pttbbs/config_util"
	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

const configPrefix = "go-pttbbs:ptttype"

func InitConfig() (err error) {
	config()

	err = postInitConfig()
	if err != nil {
		return err
	}

	err = checkTypes()
	if err != nil {
		return err
	}

	return initVars()
}

func setStringConfig(idx string, orig string) string {
	return config_util.SetStringConfig(configPrefix, idx, orig)
}

func setBoolConfig(idx string, orig bool) bool {
	return config_util.SetBoolConfig(configPrefix, idx, orig)
}

func setColorConfig(idx string, orig string) string {
	return config_util.SetColorConfig(configPrefix, idx, orig)
}

func setIntConfig(idx string, orig int) int {
	return config_util.SetIntConfig(configPrefix, idx, orig)
}

func setDoubleConfig(idx string, orig float64) float64 {
	return config_util.SetDoubleConfig(configPrefix, idx, orig)
}

func setServiceMode(serviceMode ServiceMode) ServiceMode {
	switch serviceMode {
	case DEV:
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	return serviceMode
}

//SetBBSHOME
//
//This is to safely set BBSHOME
//Public to be used in the tests of other modules.
//
//Params
//	bbshome: new bbshome
//
//Return
//	string: original bbshome
func SetBBSHOME(bbshome string) string {
	origBBSHome := BBSHOME
	log.Debugf("SetBBSHOME: %v", bbshome)

	// config.go
	BBSHOME = bbshome
	BBSPROG = BBSHOME + /* 主程式 */
		string(os.PathSeparator) +
		BBSPROGPOSTFIX

	HAVE_USERAGREEMENT = BBSHOME +
		string(os.PathSeparator) +
		HAVE_USERAGREEMENT_POSTFIX
	HAVE_USERAGREEMENT_VERSION = BBSHOME +
		string(os.PathSeparator) +
		HAVE_USERAGREEMENT_VERSION_POSTFIX
	HAVE_USERAGREEMENT_ACCEPTABLE = BBSHOME +
		string(os.PathSeparator) +
		HAVE_USERAGREEMENT_ACCEPTABLE_POSTFIX

	//common.go
	FN_CONF_BANIP = BBSHOME + // 禁止連線的 IP 列表
		string(os.PathSeparator) +
		FN_CONF_BANIP_POSTFIX
	FN_PASSWD = BBSHOME + /* User records */
		string(os.PathSeparator) +
		FN_PASSWD_POSTFIX
	FN_BOARD = BBSHOME + /* board list */
		string(os.PathSeparator) +
		FN_BOARD_POSTFIX

	//const.go
	FN_FRESH = BBSHOME + string(os.PathSeparator) + FN_FRESH_POSTFIX /* mbbsd/register.c line: 381 */

	FN_ALLOW_EMAIL_LIST = BBSHOME + string(os.PathSeparator) + FN_ALLOW_EMAIL_LIST_POSTFIX
	FN_REJECT_EMAIL_LIST = BBSHOME + string(os.PathSeparator) + FN_REJECT_EMAIL_LIST_POSTFIX

	FN_DEFAULT_FAVS = BBSHOME + string(os.PathSeparator) + FN_DEFAULT_FAVS_POSTFIX

	return origBBSHome
}

//setBBSMName
//
//This is to safely set BBSMNAME
//
//Params
//	bbsmname: new bbsmname
//
//Return
//	string: original bbsmname
func setBBSMName(bbsmname string) string {
	origBBSMName := BBSMNAME
	log.Debugf("SetBBSMNAME: %v", bbsmname)

	BBSMNAME = bbsmname

	// regex-replace

	BN_SECURITY = regexReplace(BN_SECURITY, "BBSMNAME", BBSMNAME)
	BN_NOTE = regexReplace(BN_NOTE, "BBSMNAME", BBSMNAME)
	BN_RECORD = regexReplace(BN_RECORD, "BBSMNAME", BBSMNAME)
	BN_SYSOP = regexReplace(BN_SYSOP, "BBSMNAME", BBSMNAME)
	BN_TEST = regexReplace(BN_SECURITY, "BBSMNAME", BBSMNAME)
	BN_BUGREPORT = regexReplace(BN_BUGREPORT, "BBSMNAME", BBSMNAME)
	BN_LAW = regexReplace(BN_LAW, "BBSMNAME", BBSMNAME)
	BN_NEWBIE = regexReplace(BN_NEWBIE, "BBSMNAME", BBSMNAME)
	BN_ASKBOARD = regexReplace(BN_ASKBOARD, "BBSMNAME", BBSMNAME)
	BN_FOREIGN = regexReplace(BN_FOREIGN, "BBSMNAME", BBSMNAME)

	// config.go
	if IS_BN_FIVECHESS_LOG_INFERRED {
		BN_FIVECHESS_LOG = BBSMNAME + "Five"
	}
	if IS_BN_CCHESS_LOG_INFERRED {
		BN_CCHESS_LOG = BBSMNAME + "CChess"
	}
	if IS_MONEYNAME_INFFERRED {
		MONEYNAME = BBSMNAME + "幣"
	}

	BN_BUGREPORT = BBSMNAME + "Bug"
	BN_LAW = BBSMNAME + "Law"
	BN_NEWBIE = BBSMNAME + "NewHand"
	BN_FOREIGN = BBSMNAME + "Foreign"

	return origBBSMName
}

func regexReplace(str string, substr string, repl string) string {
	theRe := regexp.MustCompile("\\s*" + substr + "\\s*")
	if theRe == nil {
		return str
	}

	return theRe.ReplaceAllString(str, repl)
}

func setCAPTCHAInsertServerAddr(captchaInsertServerAddr string) string {
	origCAPTCHAInsertServerAddr := CAPTCHA_INSERT_SERVER_ADDR

	CAPTCHA_INSERT_SERVER_ADDR = captchaInsertServerAddr

	if IS_CAPTCHA_INSERT_HOST_INFERRED {
		CAPTCHA_INSERT_HOST = CAPTCHA_INSERT_SERVER_ADDR
	}

	return origCAPTCHAInsertServerAddr
}

//setMyHostname
//
//Params
//	myHostName: new my hostname
//
//Return
//	string: orig my hostname
func setMyHostname(myHostname string) string {
	origMyHostname := MYHOSTNAME

	MYHOSTNAME = myHostname

	if IS_AID_HOSTNAME_INFERRED {
		AID_HOSTNAME = MYHOSTNAME
	}

	return origMyHostname

}

//setRecycleBinName
//
//Params
//	recycleBinName: new recycle bin name
//
//Return
//	string: orig recycle bin name
func setRecycleBinName(recycleBinName string) string {
	origRecycleBinName := recycleBinName

	RECYCLE_BIN_NAME = recycleBinName
	RECYCLE_BIN_OWNER = "[" + RECYCLE_BIN_NAME + "]"

	return origRecycleBinName
}

func postInitConfig() error {
	_ = setServiceMode(SERVICE_MODE)
	_ = SetBBSHOME(BBSHOME)
	_ = setBBSMName(BBSMNAME)
	_ = setCAPTCHAInsertServerAddr(CAPTCHA_INSERT_SERVER_ADDR)
	_ = setMyHostname(MYHOSTNAME)
	_ = setRecycleBinName(RECYCLE_BIN_NAME)

	return nil
}

func checkTypes() (err error) {
	if USEREC2_RAW_SZ != DEFAULT_USEREC2_RAW_SZ {
		log.Errorf("userec2 is not aligned: userec2: %v default-userec2: %v", USEREC2_RAW_SZ, DEFAULT_USEREC2_RAW_SZ)

		return ErrInvalidType
	}

	return nil
}

func initVars() (err error) {
	initReservedUserIDs()

	return nil
}

func initReservedUserIDs() {
	ReservedUserIDs = []types.Cstr{}

	filename := SetBBSHomePath(FN_CONF_RESERVED_ID)
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for theBytes, err := types.ReadLine(reader); err == nil; theBytes, err = types.ReadLine(reader) {
		first, _ := types.CstrTokenR(theBytes, BYTES_SPACE)
		if len(first) == 0 {
			continue
		}
		log.Infof("initReservedUserIDs: theBytes: %v first: %v", string(theBytes), string(first))

		firstBytes := types.CstrToBytes(first)
		ReservedUserIDs = append(ReservedUserIDs, firstBytes)
	}

	return
}
