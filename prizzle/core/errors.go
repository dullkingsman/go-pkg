package prizzle

// SQLITE --------------------------------------------------------------------------------------------------------------

type SqliteResultName string

type sqliteResultName struct {
	Ok                     SqliteResultName
	Error                  SqliteResultName
	Internal               SqliteResultName
	Perm                   SqliteResultName
	Abort                  SqliteResultName
	Busy                   SqliteResultName
	Locked                 SqliteResultName
	NoMem                  SqliteResultName
	ReadOnly               SqliteResultName
	Interrupt              SqliteResultName
	IoErr                  SqliteResultName
	Corrupt                SqliteResultName
	NotFound               SqliteResultName
	Full                   SqliteResultName
	CantOpen               SqliteResultName
	Protocol               SqliteResultName
	Empty                  SqliteResultName
	Schema                 SqliteResultName
	TooBig                 SqliteResultName
	Constraint             SqliteResultName
	Mismatch               SqliteResultName
	Misuse                 SqliteResultName
	NolFs                  SqliteResultName
	Auth                   SqliteResultName
	Format                 SqliteResultName
	Range                  SqliteResultName
	NotAdb                 SqliteResultName
	Notice                 SqliteResultName
	Warning                SqliteResultName
	Row                    SqliteResultName
	Done                   SqliteResultName
	OkLoadPermanently      SqliteResultName
	ErrorMissingCollSeq    SqliteResultName
	BusyRecovery           SqliteResultName
	LockedSharedCache      SqliteResultName
	ReadonlyRecovery       SqliteResultName
	IoErrRead              SqliteResultName
	CorruptVab             SqliteResultName
	CantOpenNoteMpDir      SqliteResultName
	ConstraintCheck        SqliteResultName
	AuthUser               SqliteResultName
	NoticeRecoverWal       SqliteResultName
	WarningAutoIndex       SqliteResultName
	ErrorRetry             SqliteResultName
	AbortRollback          SqliteResultName
	BusySnapshot           SqliteResultName
	LockedVtab             SqliteResultName
	ReadonlyCantLock       SqliteResultName
	IoErrShortRead         SqliteResultName
	CorruptSequence        SqliteResultName
	CantOpenIsDir          SqliteResultName
	ConstraintCommitHook   SqliteResultName
	NoticeRecoverRollback  SqliteResultName
	ErrorSnapshot          SqliteResultName
	BusyTimeout            SqliteResultName
	ReadonlyRollback       SqliteResultName
	IoErrWrite             SqliteResultName
	CorruptIndex           SqliteResultName
	CantOpenFullPath       SqliteResultName
	ConstraintForeignKey   SqliteResultName
	ReadOnlyDbMoved        SqliteResultName
	IoErrFSync             SqliteResultName
	CantOpenConvPath       SqliteResultName
	ConstraintFunction     SqliteResultName
	ReadOnlyCantInit       SqliteResultName
	IoErrDirFSync          SqliteResultName
	CantOpenDirtyWal       SqliteResultName
	ConstraintNotNull      SqliteResultName
	ReadOnlyDirectory      SqliteResultName
	IoErrTruncate          SqliteResultName
	CantOpenSymlink        SqliteResultName
	ConstraintPrimaryKey   SqliteResultName
	IoErrFStat             SqliteResultName
	ConstraintTrigger      SqliteResultName
	IoErrUnlock            SqliteResultName
	ConstraintUnique       SqliteResultName
	IoErrRdLock            SqliteResultName
	ConstraintVTab         SqliteResultName
	IoErrDelete            SqliteResultName
	ConstraintRowId        SqliteResultName
	IoErrBlocked           SqliteResultName
	ConstraintPinned       SqliteResultName
	IoErrNoMem             SqliteResultName
	ConstraintDataType     SqliteResultName
	IoErrAccess            SqliteResultName
	IoErrCheckReservedLock SqliteResultName
	IoErrLock              SqliteResultName
	IoErrClose             SqliteResultName
	IoErrDirClose          SqliteResultName
	IoErrShmOpen           SqliteResultName
	IoErrShmSize           SqliteResultName
	IoErrShmLock           SqliteResultName
	IoErrShmMap            SqliteResultName
	IoErrSeek              SqliteResultName
	IoErrDeleteNoEnt       SqliteResultName
	IoErrMMap              SqliteResultName
	IoErrGetTempPath       SqliteResultName
	IoErrConvPath          SqliteResultName
	IoErrVNode             SqliteResultName
	IoErrAuth              SqliteResultName
	IoErrBeginAtomic       SqliteResultName
	IoErrCommitAtomic      SqliteResultName
	IoErrRollbackAtomic    SqliteResultName
	IoErrData              SqliteResultName
	IoErrCorruptFs         SqliteResultName
}

var SqliteResultNameValue = sqliteResultName{
	Ok:                     "SQLITE_OK",
	Error:                  "SQLITE_ERROR",
	Internal:               "SQLITE_INTERNAL",
	Perm:                   "SQLITE_PERM",
	Abort:                  "SQLITE_ABORT",
	Busy:                   "SQLITE_BUSY",
	Locked:                 "SQLITE_LOCKED",
	NoMem:                  "SQLITE_NOMEM",
	ReadOnly:               "SQLITE_READONLY",
	Interrupt:              "SQLITE_INTERRUPT",
	IoErr:                  "SQLITE_IOERR",
	Corrupt:                "SQLITE_CORRUPT",
	NotFound:               "SQLITE_NOTFOUND",
	Full:                   "SQLITE_FULL",
	CantOpen:               "SQLITE_CANTOPEN",
	Protocol:               "SQLITE_PROTOCOL",
	Empty:                  "SQLITE_EMPTY",
	Schema:                 "SQLITE_SCHEMA",
	TooBig:                 "SQLITE_TOOBIG",
	Constraint:             "SQLITE_CONSTRAINT",
	Mismatch:               "SQLITE_MISMATCH",
	Misuse:                 "SQLITE_MISUSE",
	NolFs:                  "SQLITE_NOLFS",
	Auth:                   "SQLITE_AUTH",
	Format:                 "SQLITE_FORMAT",
	Range:                  "SQLITE_RANGE",
	NotAdb:                 "SQLITE_NOTADB",
	Notice:                 "SQLITE_NOTICE",
	Warning:                "SQLITE_WARNING",
	Row:                    "SQLITE_ROW",
	Done:                   "SQLITE_DONE",
	OkLoadPermanently:      "SQLITE_OK_LOAD_PERMANENTLY",
	ErrorMissingCollSeq:    "SQLITE_ERROR_MISSING_COLLSEQ",
	BusyRecovery:           "SQLITE_BUSY_RECOVERY",
	LockedSharedCache:      "SQLITE_LOCKED_SHAREDCACHE",
	ReadonlyRecovery:       "SQLITE_READONLY_RECOVERY",
	IoErrRead:              "SQLITE_IOERR_READ",
	CorruptVab:             "SQLITE_CORRUPT_VTAB",
	CantOpenNoteMpDir:      "SQLITE_CANTOPEN_NOTEMPDIR",
	ConstraintCheck:        "SQLITE_CONSTRAINT_CHECK",
	AuthUser:               "SQLITE_AUTH_USER",
	NoticeRecoverWal:       "SQLITE_NOTICE_RECOVER_WAL",
	WarningAutoIndex:       "SQLITE_WARNING_AUTOINDEX",
	ErrorRetry:             "SQLITE_ERROR_RETRY",
	AbortRollback:          "SQLITE_ABORT_ROLLBACK",
	BusySnapshot:           "SQLITE_BUSY_SNAPSHOT",
	LockedVtab:             "SQLITE_LOCKED_VTAB",
	ReadonlyCantLock:       "SQLITE_READONLY_CANTLOCK",
	IoErrShortRead:         "SQLITE_IOERR_SHORT_READ",
	CorruptSequence:        "SQLITE_CORRUPT_SEQUENCE",
	CantOpenIsDir:          "SQLITE_CANTOPEN_ISDIR",
	ConstraintCommitHook:   "SQLITE_CONSTRAINT_COMMITHOOK",
	NoticeRecoverRollback:  "SQLITE_NOTICE_RECOVER_ROLLBACK",
	ErrorSnapshot:          "SQLITE_ERROR_SNAPSHOT",
	BusyTimeout:            "SQLITE_BUSY_TIMEOUT",
	ReadonlyRollback:       "SQLITE_READONLY_ROLLBACK",
	IoErrWrite:             "SQLITE_IOERR_WRITE",
	CorruptIndex:           "SQLITE_CORRUPT_INDEX",
	CantOpenFullPath:       "SQLITE_CANTOPEN_FULLPATH",
	ConstraintForeignKey:   "SQLITE_CONSTRAINT_FOREIGNKEY",
	ReadOnlyDbMoved:        "SQLITE_READONLY_DBMOVED",
	IoErrFSync:             "SQLITE_IOERR_FSYNC",
	CantOpenConvPath:       "SQLITE_CANTOPEN_CONVPATH",
	ConstraintFunction:     "SQLITE_CONSTRAINT_FUNCTION",
	ReadOnlyCantInit:       "SQLITE_READONLY_CANTINIT",
	IoErrDirFSync:          "SQLITE_IOERR_DIR_FSYNC",
	CantOpenDirtyWal:       "SQLITE_CANTOPEN_DIRTYWAL",
	ConstraintNotNull:      "SQLITE_CONSTRAINT_NOTNULL",
	ReadOnlyDirectory:      "SQLITE_READONLY_DIRECTORY",
	IoErrTruncate:          "SQLITE_IOERR_TRUNCATE",
	CantOpenSymlink:        "SQLITE_CANTOPEN_SYMLINK",
	ConstraintPrimaryKey:   "SQLITE_CONSTRAINT_PRIMARYKEY",
	IoErrFStat:             "SQLITE_IOERR_FSTAT",
	ConstraintTrigger:      "SQLITE_CONSTRAINT_TRIGGER",
	IoErrUnlock:            "SQLITE_IOERR_UNLOCK",
	ConstraintUnique:       "SQLITE_CONSTRAINT_UNIQUE",
	IoErrRdLock:            "SQLITE_IOERR_RDLOCK",
	ConstraintVTab:         "SQLITE_CONSTRAINT_VTAB",
	IoErrDelete:            "SQLITE_IOERR_DELETE",
	ConstraintRowId:        "SQLITE_CONSTRAINT_ROWID",
	IoErrBlocked:           "SQLITE_IOERR_BLOCKED",
	ConstraintPinned:       "SQLITE_CONSTRAINT_PINNED",
	IoErrNoMem:             "SQLITE_IOERR_NOMEM",
	ConstraintDataType:     "SQLITE_CONSTRAINT_DATATYPE",
	IoErrAccess:            "SQLITE_IOERR_ACCESS",
	IoErrCheckReservedLock: "SQLITE_IOERR_CHECKRESERVEDLOCK",
	IoErrLock:              "SQLITE_IOERR_LOCK",
	IoErrClose:             "SQLITE_IOERR_CLOSE",
	IoErrDirClose:          "SQLITE_IOERR_DIR_CLOSE",
	IoErrShmOpen:           "SQLITE_IOERR_SHMOPEN",
	IoErrShmSize:           "SQLITE_IOERR_SHMSIZE",
	IoErrShmLock:           "SQLITE_IOERR_SHMLOCK",
	IoErrShmMap:            "SQLITE_IOERR_SHMMAP",
	IoErrSeek:              "SQLITE_IOERR_SEEK",
	IoErrDeleteNoEnt:       "SQLITE_IOERR_DELETE_NOENT",
	IoErrMMap:              "SQLITE_IOERR_MMAP",
	IoErrGetTempPath:       "SQLITE_IOERR_GETTEMPPATH",
	IoErrConvPath:          "SQLITE_IOERR_CONVPATH",
	IoErrVNode:             "SQLITE_IOERR_VNODE",
	IoErrAuth:              "SQLITE_IOERR_AUTH",
	IoErrBeginAtomic:       "SQLITE_IOERR_BEGIN_ATOMIC",
	IoErrCommitAtomic:      "SQLITE_IOERR_COMMIT_ATOMIC",
	IoErrRollbackAtomic:    "SQLITE_IOERR_ROLLBACK_ATOMIC",
	IoErrData:              "SQLITE_IOERR_DATA",
	IoErrCorruptFs:         "SQLITE_IOERR_CORRUPTFS",
}

var SqliteResultNameValues = []SqliteResultName{
	"SQLITE_OK",
	"SQLITE_ERROR",
	"SQLITE_INTERNAL",
	"SQLITE_PERM",
	"SQLITE_ABORT",
	"SQLITE_BUSY",
	"SQLITE_LOCKED",
	"SQLITE_NOMEM",
	"SQLITE_READONLY",
	"SQLITE_INTERRUPT",
	"SQLITE_IOERR",
	"SQLITE_CORRUPT",
	"SQLITE_NOTFOUND",
	"SQLITE_FULL",
	"SQLITE_CANTOPEN",
	"SQLITE_PROTOCOL",
	"SQLITE_EMPTY",
	"SQLITE_SCHEMA",
	"SQLITE_TOOBIG",
	"SQLITE_CONSTRAINT",
	"SQLITE_MISMATCH",
	"SQLITE_MISUSE",
	"SQLITE_NOLFS",
	"SQLITE_AUTH",
	"SQLITE_FORMAT",
	"SQLITE_RANGE",
	"SQLITE_NOTADB",
	"SQLITE_NOTICE",
	"SQLITE_WARNING",
	"SQLITE_ROW",
	"SQLITE_DONE",
	"SQLITE_OK_LOAD_PERMANENTLY",
	"SQLITE_ERROR_MISSING_COLLSEQ",
	"SQLITE_BUSY_RECOVERY",
	"SQLITE_LOCKED_SHAREDCACHE",
	"SQLITE_READONLY_RECOVERY",
	"SQLITE_IOERR_READ",
	"SQLITE_CORRUPT_VTAB",
	"SQLITE_CANTOPEN_NOTEMPDIR",
	"SQLITE_CONSTRAINT_CHECK",
	"SQLITE_AUTH_USER",
	"SQLITE_NOTICE_RECOVER_WAL",
	"SQLITE_WARNING_AUTOINDEX",
	"SQLITE_ERROR_RETRY",
	"SQLITE_ABORT_ROLLBACK",
	"SQLITE_BUSY_SNAPSHOT",
	"SQLITE_LOCKED_VTAB",
	"SQLITE_READONLY_CANTLOCK",
	"SQLITE_IOERR_SHORT_READ",
	"SQLITE_CORRUPT_SEQUENCE",
	"SQLITE_CANTOPEN_ISDIR",
	"SQLITE_CONSTRAINT_COMMITHOOK",
	"SQLITE_NOTICE_RECOVER_ROLLBACK",
	"SQLITE_ERROR_SNAPSHOT",
	"SQLITE_BUSY_TIMEOUT",
	"SQLITE_READONLY_ROLLBACK",
	"SQLITE_IOERR_WRITE",
	"SQLITE_CORRUPT_INDEX",
	"SQLITE_CANTOPEN_FULLPATH",
	"SQLITE_CONSTRAINT_FOREIGNKEY",
	"SQLITE_READONLY_DBMOVED",
	"SQLITE_IOERR_FSYNC",
	"SQLITE_CANTOPEN_CONVPATH",
	"SQLITE_CONSTRAINT_FUNCTION",
	"SQLITE_READONLY_CANTINIT",
	"SQLITE_IOERR_DIR_FSYNC",
	"SQLITE_CANTOPEN_DIRTYWAL",
	"SQLITE_CONSTRAINT_NOTNULL",
	"SQLITE_READONLY_DIRECTORY",
	"SQLITE_IOERR_TRUNCATE",
	"SQLITE_CANTOPEN_SYMLINK",
	"SQLITE_CONSTRAINT_PRIMARYKEY",
	"SQLITE_IOERR_FSTAT",
	"SQLITE_CONSTRAINT_TRIGGER",
	"SQLITE_IOERR_UNLOCK",
	"SQLITE_CONSTRAINT_UNIQUE",
	"SQLITE_IOERR_RDLOCK",
	"SQLITE_CONSTRAINT_VTAB",
	"SQLITE_IOERR_DELETE",
	"SQLITE_CONSTRAINT_ROWID",
	"SQLITE_IOERR_BLOCKED",
	"SQLITE_CONSTRAINT_PINNED",
	"SQLITE_IOERR_NOMEM",
	"SQLITE_CONSTRAINT_DATATYPE",
	"SQLITE_IOERR_ACCESS",
	"SQLITE_IOERR_CHECKRESERVEDLOCK",
	"SQLITE_IOERR_LOCK",
	"SQLITE_IOERR_CLOSE",
	"SQLITE_IOERR_DIR_CLOSE",
	"SQLITE_IOERR_SHMOPEN",
	"SQLITE_IOERR_SHMSIZE",
	"SQLITE_IOERR_SHMLOCK",
	"SQLITE_IOERR_SHMMAP",
	"SQLITE_IOERR_SEEK",
	"SQLITE_IOERR_DELETE_NOENT",
	"SQLITE_IOERR_MMAP",
	"SQLITE_IOERR_GETTEMPPATH",
	"SQLITE_IOERR_CONVPATH",
	"SQLITE_IOERR_VNODE",
	"SQLITE_IOERR_AUTH",
	"SQLITE_IOERR_BEGIN_ATOMIC",
	"SQLITE_IOERR_COMMIT_ATOMIC",
	"SQLITE_IOERR_ROLLBACK_ATOMIC",
	"SQLITE_IOERR_DATA",
	"SQLITE_IOERR_CORRUPTFS",
}

type SqliteResultCode string

type sqliteResultCode struct {
	Ok                     SqliteResultCode
	Error                  SqliteResultCode
	Internal               SqliteResultCode
	Perm                   SqliteResultCode
	Abort                  SqliteResultCode
	Busy                   SqliteResultCode
	Locked                 SqliteResultCode
	NoMem                  SqliteResultCode
	ReadOnly               SqliteResultCode
	Interrupt              SqliteResultCode
	IoErr                  SqliteResultCode
	Corrupt                SqliteResultCode
	NotFound               SqliteResultCode
	Full                   SqliteResultCode
	CantOpen               SqliteResultCode
	Protocol               SqliteResultCode
	Empty                  SqliteResultCode
	Schema                 SqliteResultCode
	TooBig                 SqliteResultCode
	Constraint             SqliteResultCode
	Mismatch               SqliteResultCode
	Misuse                 SqliteResultCode
	NolFs                  SqliteResultCode
	Auth                   SqliteResultCode
	Format                 SqliteResultCode
	Range                  SqliteResultCode
	NotAdb                 SqliteResultCode
	Notice                 SqliteResultCode
	Warning                SqliteResultCode
	Row                    SqliteResultCode
	Done                   SqliteResultCode
	OkLoadPermanently      SqliteResultCode
	ErrorMissingCollSeq    SqliteResultCode
	BusyRecovery           SqliteResultCode
	LockedSharedCache      SqliteResultCode
	ReadonlyRecovery       SqliteResultCode
	IoErrRead              SqliteResultCode
	CorruptVab             SqliteResultCode
	CantOpenNoteMpDir      SqliteResultCode
	ConstraintCheck        SqliteResultCode
	AuthUser               SqliteResultCode
	NoticeRecoverWal       SqliteResultCode
	WarningAutoIndex       SqliteResultCode
	ErrorRetry             SqliteResultCode
	AbortRollback          SqliteResultCode
	BusySnapshot           SqliteResultCode
	LockedVtab             SqliteResultCode
	ReadonlyCantLock       SqliteResultCode
	IoErrShortRead         SqliteResultCode
	CorruptSequence        SqliteResultCode
	CantOpenIsDir          SqliteResultCode
	ConstraintCommitHook   SqliteResultCode
	NoticeRecoverRollback  SqliteResultCode
	ErrorSnapshot          SqliteResultCode
	BusyTimeout            SqliteResultCode
	ReadonlyRollback       SqliteResultCode
	IoErrWrite             SqliteResultCode
	CorruptIndex           SqliteResultCode
	CantOpenFullPath       SqliteResultCode
	ConstraintForeignKey   SqliteResultCode
	ReadOnlyDbMoved        SqliteResultCode
	IoErrFSync             SqliteResultCode
	CantOpenConvPath       SqliteResultCode
	ConstraintFunction     SqliteResultCode
	ReadOnlyCantInit       SqliteResultCode
	IoErrDirFSync          SqliteResultCode
	CantOpenDirtyWal       SqliteResultCode
	ConstraintNotNull      SqliteResultCode
	ReadOnlyDirectory      SqliteResultCode
	IoErrTruncate          SqliteResultCode
	CantOpenSymlink        SqliteResultCode
	ConstraintPrimaryKey   SqliteResultCode
	IoErrFStat             SqliteResultCode
	ConstraintTrigger      SqliteResultCode
	IoErrUnlock            SqliteResultCode
	ConstraintUnique       SqliteResultCode
	IoErrRdLock            SqliteResultCode
	ConstraintVTab         SqliteResultCode
	IoErrDelete            SqliteResultCode
	ConstraintRowId        SqliteResultCode
	IoErrBlocked           SqliteResultCode
	ConstraintPinned       SqliteResultCode
	IoErrNoMem             SqliteResultCode
	ConstraintDataType     SqliteResultCode
	IoErrAccess            SqliteResultCode
	IoErrCheckReservedLock SqliteResultCode
	IoErrLock              SqliteResultCode
	IoErrClose             SqliteResultCode
	IoErrDirClose          SqliteResultCode
	IoErrShmOpen           SqliteResultCode
	IoErrShmSize           SqliteResultCode
	IoErrShmLock           SqliteResultCode
	IoErrShmMap            SqliteResultCode
	IoErrSeek              SqliteResultCode
	IoErrDeleteNoEnt       SqliteResultCode
	IoErrMMap              SqliteResultCode
	IoErrGetTempPath       SqliteResultCode
	IoErrConvPath          SqliteResultCode
	IoErrVNode             SqliteResultCode
	IoErrAuth              SqliteResultCode
	IoErrBeginAtomic       SqliteResultCode
	IoErrCommitAtomic      SqliteResultCode
	IoErrRollbackAtomic    SqliteResultCode
	IoErrData              SqliteResultCode
	IoErrCorruptFs         SqliteResultCode
}

var SqliteResultCodeValue = sqliteResultCode{
	Ok:                     "0",
	Error:                  "1",
	Internal:               "2",
	Perm:                   "3",
	Abort:                  "4",
	Busy:                   "5",
	Locked:                 "6",
	NoMem:                  "7",
	ReadOnly:               "8",
	Interrupt:              "9",
	IoErr:                  "10",
	Corrupt:                "11",
	NotFound:               "12",
	Full:                   "13",
	CantOpen:               "14",
	Protocol:               "15",
	Empty:                  "16",
	Schema:                 "17",
	TooBig:                 "18",
	Constraint:             "19",
	Mismatch:               "20",
	Misuse:                 "21",
	NolFs:                  "22",
	Auth:                   "23",
	Format:                 "24",
	Range:                  "25",
	NotAdb:                 "26",
	Notice:                 "27",
	Warning:                "28",
	Row:                    "100",
	Done:                   "101",
	OkLoadPermanently:      "256",
	ErrorMissingCollSeq:    "257",
	BusyRecovery:           "261",
	LockedSharedCache:      "262",
	ReadonlyRecovery:       "264",
	IoErrRead:              "266",
	CorruptVab:             "267",
	CantOpenNoteMpDir:      "270",
	ConstraintCheck:        "275",
	AuthUser:               "279",
	NoticeRecoverWal:       "283",
	WarningAutoIndex:       "284",
	ErrorRetry:             "513",
	AbortRollback:          "516",
	BusySnapshot:           "517",
	LockedVtab:             "518",
	ReadonlyCantLock:       "520",
	IoErrShortRead:         "522",
	CorruptSequence:        "523",
	CantOpenIsDir:          "526",
	ConstraintCommitHook:   "531",
	NoticeRecoverRollback:  "539",
	ErrorSnapshot:          "769",
	BusyTimeout:            "773",
	ReadonlyRollback:       "776",
	IoErrWrite:             "778",
	CorruptIndex:           "779",
	CantOpenFullPath:       "782",
	ConstraintForeignKey:   "787",
	ReadOnlyDbMoved:        "1032",
	IoErrFSync:             "1034",
	CantOpenConvPath:       "1038",
	ConstraintFunction:     "1043",
	ReadOnlyCantInit:       "1288",
	IoErrDirFSync:          "1290",
	CantOpenDirtyWal:       "1294",
	ConstraintNotNull:      "1299",
	ReadOnlyDirectory:      "1544",
	IoErrTruncate:          "1546",
	CantOpenSymlink:        "1550",
	ConstraintPrimaryKey:   "1555",
	IoErrFStat:             "1802",
	ConstraintTrigger:      "1811",
	IoErrUnlock:            "2058",
	ConstraintUnique:       "2067",
	IoErrRdLock:            "2314",
	ConstraintVTab:         "2323",
	IoErrDelete:            "2570",
	ConstraintRowId:        "2579",
	IoErrBlocked:           "2826",
	ConstraintPinned:       "2835",
	IoErrNoMem:             "3082",
	ConstraintDataType:     "3091",
	IoErrAccess:            "3338",
	IoErrCheckReservedLock: "3594",
	IoErrLock:              "3850",
	IoErrClose:             "4106",
	IoErrDirClose:          "4362",
	IoErrShmOpen:           "4618",
	IoErrShmSize:           "4874",
	IoErrShmLock:           "5130",
	IoErrShmMap:            "5386",
	IoErrSeek:              "5642",
	IoErrDeleteNoEnt:       "5898",
	IoErrMMap:              "6154",
	IoErrGetTempPath:       "6410",
	IoErrConvPath:          "6666",
	IoErrVNode:             "6922",
	IoErrAuth:              "7178",
	IoErrBeginAtomic:       "7434",
	IoErrCommitAtomic:      "7690",
	IoErrRollbackAtomic:    "7946",
	IoErrData:              "8202",
	IoErrCorruptFs:         "8458",
}

var SqliteResultCodeValues = []SqliteResultCode{
	"0",
	"1",
	"2",
	"3",
	"4",
	"5",
	"6",
	"7",
	"8",
	"9",
	"10",
	"11",
	"12",
	"13",
	"14",
	"15",
	"16",
	"17",
	"18",
	"19",
	"20",
	"21",
	"22",
	"23",
	"24",
	"25",
	"26",
	"27",
	"28",
	"100",
	"101",
	"256",
	"257",
	"261",
	"262",
	"264",
	"266",
	"267",
	"270",
	"275",
	"279",
	"283",
	"284",
	"513",
	"516",
	"517",
	"518",
	"520",
	"522",
	"523",
	"526",
	"531",
	"539",
	"769",
	"773",
	"776",
	"778",
	"779",
	"782",
	"787",
	"1032",
	"1034",
	"1038",
	"1043",
	"1288",
	"1290",
	"1294",
	"1299",
	"1544",
	"1546",
	"1550",
	"1555",
	"1802",
	"1811",
	"2058",
	"2067",
	"2314",
	"2323",
	"2570",
	"2579",
	"2826",
	"2835",
	"3082",
	"3091",
	"3338",
	"3594",
	"3850",
	"4106",
	"4362",
	"4618",
	"4874",
	"5130",
	"5386",
	"5642",
	"5898",
	"6154",
	"6410",
	"6666",
	"6922",
	"7178",
	"7434",
	"7690",
	"7946",
	"8202",
	"8458",
}

var SqliteResultNameToCode = map[SqliteResultName]SqliteResultCode{
	SqliteResultNameValue.Ok:                     "0",
	SqliteResultNameValue.Error:                  "1",
	SqliteResultNameValue.Internal:               "2",
	SqliteResultNameValue.Perm:                   "3",
	SqliteResultNameValue.Abort:                  "4",
	SqliteResultNameValue.Busy:                   "5",
	SqliteResultNameValue.Locked:                 "6",
	SqliteResultNameValue.NoMem:                  "7",
	SqliteResultNameValue.ReadOnly:               "8",
	SqliteResultNameValue.Interrupt:              "9",
	SqliteResultNameValue.IoErr:                  "10",
	SqliteResultNameValue.Corrupt:                "11",
	SqliteResultNameValue.NotFound:               "12",
	SqliteResultNameValue.Full:                   "13",
	SqliteResultNameValue.CantOpen:               "14",
	SqliteResultNameValue.Protocol:               "15",
	SqliteResultNameValue.Empty:                  "16",
	SqliteResultNameValue.Schema:                 "17",
	SqliteResultNameValue.TooBig:                 "18",
	SqliteResultNameValue.Constraint:             "19",
	SqliteResultNameValue.Mismatch:               "20",
	SqliteResultNameValue.Misuse:                 "21",
	SqliteResultNameValue.NolFs:                  "22",
	SqliteResultNameValue.Auth:                   "23",
	SqliteResultNameValue.Format:                 "24",
	SqliteResultNameValue.Range:                  "25",
	SqliteResultNameValue.NotAdb:                 "26",
	SqliteResultNameValue.Notice:                 "27",
	SqliteResultNameValue.Warning:                "28",
	SqliteResultNameValue.Row:                    "100",
	SqliteResultNameValue.Done:                   "101",
	SqliteResultNameValue.OkLoadPermanently:      "256",
	SqliteResultNameValue.ErrorMissingCollSeq:    "257",
	SqliteResultNameValue.BusyRecovery:           "261",
	SqliteResultNameValue.LockedSharedCache:      "262",
	SqliteResultNameValue.ReadonlyRecovery:       "264",
	SqliteResultNameValue.IoErrRead:              "266",
	SqliteResultNameValue.CorruptVab:             "267",
	SqliteResultNameValue.CantOpenNoteMpDir:      "270",
	SqliteResultNameValue.ConstraintCheck:        "275",
	SqliteResultNameValue.AuthUser:               "279",
	SqliteResultNameValue.NoticeRecoverWal:       "283",
	SqliteResultNameValue.WarningAutoIndex:       "284",
	SqliteResultNameValue.ErrorRetry:             "513",
	SqliteResultNameValue.AbortRollback:          "516",
	SqliteResultNameValue.BusySnapshot:           "517",
	SqliteResultNameValue.LockedVtab:             "518",
	SqliteResultNameValue.ReadonlyCantLock:       "520",
	SqliteResultNameValue.IoErrShortRead:         "522",
	SqliteResultNameValue.CorruptSequence:        "523",
	SqliteResultNameValue.CantOpenIsDir:          "526",
	SqliteResultNameValue.ConstraintCommitHook:   "531",
	SqliteResultNameValue.NoticeRecoverRollback:  "539",
	SqliteResultNameValue.ErrorSnapshot:          "769",
	SqliteResultNameValue.BusyTimeout:            "773",
	SqliteResultNameValue.ReadonlyRollback:       "776",
	SqliteResultNameValue.IoErrWrite:             "778",
	SqliteResultNameValue.CorruptIndex:           "779",
	SqliteResultNameValue.CantOpenFullPath:       "782",
	SqliteResultNameValue.ConstraintForeignKey:   "787",
	SqliteResultNameValue.ReadOnlyDbMoved:        "1032",
	SqliteResultNameValue.IoErrFSync:             "1034",
	SqliteResultNameValue.CantOpenConvPath:       "1038",
	SqliteResultNameValue.ConstraintFunction:     "1043",
	SqliteResultNameValue.ReadOnlyCantInit:       "1288",
	SqliteResultNameValue.IoErrDirFSync:          "1290",
	SqliteResultNameValue.CantOpenDirtyWal:       "1294",
	SqliteResultNameValue.ConstraintNotNull:      "1299",
	SqliteResultNameValue.ReadOnlyDirectory:      "1544",
	SqliteResultNameValue.IoErrTruncate:          "1546",
	SqliteResultNameValue.CantOpenSymlink:        "1550",
	SqliteResultNameValue.ConstraintPrimaryKey:   "1555",
	SqliteResultNameValue.IoErrFStat:             "1802",
	SqliteResultNameValue.ConstraintTrigger:      "1811",
	SqliteResultNameValue.IoErrUnlock:            "2058",
	SqliteResultNameValue.ConstraintUnique:       "2067",
	SqliteResultNameValue.IoErrRdLock:            "2314",
	SqliteResultNameValue.ConstraintVTab:         "2323",
	SqliteResultNameValue.IoErrDelete:            "2570",
	SqliteResultNameValue.ConstraintRowId:        "2579",
	SqliteResultNameValue.IoErrBlocked:           "2826",
	SqliteResultNameValue.ConstraintPinned:       "2835",
	SqliteResultNameValue.IoErrNoMem:             "3082",
	SqliteResultNameValue.ConstraintDataType:     "3091",
	SqliteResultNameValue.IoErrAccess:            "3338",
	SqliteResultNameValue.IoErrCheckReservedLock: "3594",
	SqliteResultNameValue.IoErrLock:              "3850",
	SqliteResultNameValue.IoErrClose:             "4106",
	SqliteResultNameValue.IoErrDirClose:          "4362",
	SqliteResultNameValue.IoErrShmOpen:           "4618",
	SqliteResultNameValue.IoErrShmSize:           "4874",
	SqliteResultNameValue.IoErrShmLock:           "5130",
	SqliteResultNameValue.IoErrShmMap:            "5386",
	SqliteResultNameValue.IoErrSeek:              "5642",
	SqliteResultNameValue.IoErrDeleteNoEnt:       "5898",
	SqliteResultNameValue.IoErrMMap:              "6154",
	SqliteResultNameValue.IoErrGetTempPath:       "6410",
	SqliteResultNameValue.IoErrConvPath:          "6666",
	SqliteResultNameValue.IoErrVNode:             "6922",
	SqliteResultNameValue.IoErrAuth:              "7178",
	SqliteResultNameValue.IoErrBeginAtomic:       "7434",
	SqliteResultNameValue.IoErrCommitAtomic:      "7690",
	SqliteResultNameValue.IoErrRollbackAtomic:    "7946",
	SqliteResultNameValue.IoErrData:              "8202",
	SqliteResultNameValue.IoErrCorruptFs:         "8458",
}

var SqliteResultCodeToName = map[SqliteResultCode]SqliteResultName{
	SqliteResultCodeValue.Ok:                     "SQLITE_OK",
	SqliteResultCodeValue.Error:                  "SQLITE_ERROR",
	SqliteResultCodeValue.Internal:               "SQLITE_INTERNAL",
	SqliteResultCodeValue.Perm:                   "SQLITE_PERM",
	SqliteResultCodeValue.Abort:                  "SQLITE_ABORT",
	SqliteResultCodeValue.Busy:                   "SQLITE_BUSY",
	SqliteResultCodeValue.Locked:                 "SQLITE_LOCKED",
	SqliteResultCodeValue.NoMem:                  "SQLITE_NOMEM",
	SqliteResultCodeValue.ReadOnly:               "SQLITE_READONLY",
	SqliteResultCodeValue.Interrupt:              "SQLITE_INTERRUPT",
	SqliteResultCodeValue.IoErr:                  "SQLITE_IOERR",
	SqliteResultCodeValue.Corrupt:                "SQLITE_CORRUPT",
	SqliteResultCodeValue.NotFound:               "SQLITE_NOTFOUND",
	SqliteResultCodeValue.Full:                   "SQLITE_FULL",
	SqliteResultCodeValue.CantOpen:               "SQLITE_CANTOPEN",
	SqliteResultCodeValue.Protocol:               "SQLITE_PROTOCOL",
	SqliteResultCodeValue.Empty:                  "SQLITE_EMPTY",
	SqliteResultCodeValue.Schema:                 "SQLITE_SCHEMA",
	SqliteResultCodeValue.TooBig:                 "SQLITE_TOOBIG",
	SqliteResultCodeValue.Constraint:             "SQLITE_CONSTRAINT",
	SqliteResultCodeValue.Mismatch:               "SQLITE_MISMATCH",
	SqliteResultCodeValue.Misuse:                 "SQLITE_MISUSE",
	SqliteResultCodeValue.NolFs:                  "SQLITE_NOLFS",
	SqliteResultCodeValue.Auth:                   "SQLITE_AUTH",
	SqliteResultCodeValue.Format:                 "SQLITE_FORMAT",
	SqliteResultCodeValue.Range:                  "SQLITE_RANGE",
	SqliteResultCodeValue.NotAdb:                 "SQLITE_NOTADB",
	SqliteResultCodeValue.Notice:                 "SQLITE_NOTICE",
	SqliteResultCodeValue.Warning:                "SQLITE_WARNING",
	SqliteResultCodeValue.Row:                    "SQLITE_ROW",
	SqliteResultCodeValue.Done:                   "SQLITE_DONE",
	SqliteResultCodeValue.OkLoadPermanently:      "SQLITE_OK_LOAD_PERMANENTLY",
	SqliteResultCodeValue.ErrorMissingCollSeq:    "SQLITE_ERROR_MISSING_COLLSEQ",
	SqliteResultCodeValue.BusyRecovery:           "SQLITE_BUSY_RECOVERY",
	SqliteResultCodeValue.LockedSharedCache:      "SQLITE_LOCKED_SHAREDCACHE",
	SqliteResultCodeValue.ReadonlyRecovery:       "SQLITE_READONLY_RECOVERY",
	SqliteResultCodeValue.IoErrRead:              "SQLITE_IOERR_READ",
	SqliteResultCodeValue.CorruptVab:             "SQLITE_CORRUPT_VTAB",
	SqliteResultCodeValue.CantOpenNoteMpDir:      "SQLITE_CANTOPEN_NOTEMPDIR",
	SqliteResultCodeValue.ConstraintCheck:        "SQLITE_CONSTRAINT_CHECK",
	SqliteResultCodeValue.AuthUser:               "SQLITE_AUTH_USER",
	SqliteResultCodeValue.NoticeRecoverWal:       "SQLITE_NOTICE_RECOVER_WAL",
	SqliteResultCodeValue.WarningAutoIndex:       "SQLITE_WARNING_AUTOINDEX",
	SqliteResultCodeValue.ErrorRetry:             "SQLITE_ERROR_RETRY",
	SqliteResultCodeValue.AbortRollback:          "SQLITE_ABORT_ROLLBACK",
	SqliteResultCodeValue.BusySnapshot:           "SQLITE_BUSY_SNAPSHOT",
	SqliteResultCodeValue.LockedVtab:             "SQLITE_LOCKED_VTAB",
	SqliteResultCodeValue.ReadonlyCantLock:       "SQLITE_READONLY_CANTLOCK",
	SqliteResultCodeValue.IoErrShortRead:         "SQLITE_IOERR_SHORT_READ",
	SqliteResultCodeValue.CorruptSequence:        "SQLITE_CORRUPT_SEQUENCE",
	SqliteResultCodeValue.CantOpenIsDir:          "SQLITE_CANTOPEN_ISDIR",
	SqliteResultCodeValue.ConstraintCommitHook:   "SQLITE_CONSTRAINT_COMMITHOOK",
	SqliteResultCodeValue.NoticeRecoverRollback:  "SQLITE_NOTICE_RECOVER_ROLLBACK",
	SqliteResultCodeValue.ErrorSnapshot:          "SQLITE_ERROR_SNAPSHOT",
	SqliteResultCodeValue.BusyTimeout:            "SQLITE_BUSY_TIMEOUT",
	SqliteResultCodeValue.ReadonlyRollback:       "SQLITE_READONLY_ROLLBACK",
	SqliteResultCodeValue.IoErrWrite:             "SQLITE_IOERR_WRITE",
	SqliteResultCodeValue.CorruptIndex:           "SQLITE_CORRUPT_INDEX",
	SqliteResultCodeValue.CantOpenFullPath:       "SQLITE_CANTOPEN_FULLPATH",
	SqliteResultCodeValue.ConstraintForeignKey:   "SQLITE_CONSTRAINT_FOREIGNKEY",
	SqliteResultCodeValue.ReadOnlyDbMoved:        "SQLITE_READONLY_DBMOVED",
	SqliteResultCodeValue.IoErrFSync:             "SQLITE_IOERR_FSYNC",
	SqliteResultCodeValue.CantOpenConvPath:       "SQLITE_CANTOPEN_CONVPATH",
	SqliteResultCodeValue.ConstraintFunction:     "SQLITE_CONSTRAINT_FUNCTION",
	SqliteResultCodeValue.ReadOnlyCantInit:       "SQLITE_READONLY_CANTINIT",
	SqliteResultCodeValue.IoErrDirFSync:          "SQLITE_IOERR_DIR_FSYNC",
	SqliteResultCodeValue.CantOpenDirtyWal:       "SQLITE_CANTOPEN_DIRTYWAL",
	SqliteResultCodeValue.ConstraintNotNull:      "SQLITE_CONSTRAINT_NOTNULL",
	SqliteResultCodeValue.ReadOnlyDirectory:      "SQLITE_READONLY_DIRECTORY",
	SqliteResultCodeValue.IoErrTruncate:          "SQLITE_IOERR_TRUNCATE",
	SqliteResultCodeValue.CantOpenSymlink:        "SQLITE_CANTOPEN_SYMLINK",
	SqliteResultCodeValue.ConstraintPrimaryKey:   "SQLITE_CONSTRAINT_PRIMARYKEY",
	SqliteResultCodeValue.IoErrFStat:             "SQLITE_IOERR_FSTAT",
	SqliteResultCodeValue.ConstraintTrigger:      "SQLITE_CONSTRAINT_TRIGGER",
	SqliteResultCodeValue.IoErrUnlock:            "SQLITE_IOERR_UNLOCK",
	SqliteResultCodeValue.ConstraintUnique:       "SQLITE_CONSTRAINT_UNIQUE",
	SqliteResultCodeValue.IoErrRdLock:            "SQLITE_IOERR_RDLOCK",
	SqliteResultCodeValue.ConstraintVTab:         "SQLITE_CONSTRAINT_VTAB",
	SqliteResultCodeValue.IoErrDelete:            "SQLITE_IOERR_DELETE",
	SqliteResultCodeValue.ConstraintRowId:        "SQLITE_CONSTRAINT_ROWID",
	SqliteResultCodeValue.IoErrBlocked:           "SQLITE_IOERR_BLOCKED",
	SqliteResultCodeValue.ConstraintPinned:       "SQLITE_CONSTRAINT_PINNED",
	SqliteResultCodeValue.IoErrNoMem:             "SQLITE_IOERR_NOMEM",
	SqliteResultCodeValue.ConstraintDataType:     "SQLITE_CONSTRAINT_DATATYPE",
	SqliteResultCodeValue.IoErrAccess:            "SQLITE_IOERR_ACCESS",
	SqliteResultCodeValue.IoErrCheckReservedLock: "SQLITE_IOERR_CHECKRESERVEDLOCK",
	SqliteResultCodeValue.IoErrLock:              "SQLITE_IOERR_LOCK",
	SqliteResultCodeValue.IoErrClose:             "SQLITE_IOERR_CLOSE",
	SqliteResultCodeValue.IoErrDirClose:          "SQLITE_IOERR_DIR_CLOSE",
	SqliteResultCodeValue.IoErrShmOpen:           "SQLITE_IOERR_SHMOPEN",
	SqliteResultCodeValue.IoErrShmSize:           "SQLITE_IOERR_SHMSIZE",
	SqliteResultCodeValue.IoErrShmLock:           "SQLITE_IOERR_SHMLOCK",
	SqliteResultCodeValue.IoErrShmMap:            "SQLITE_IOERR_SHMMAP",
	SqliteResultCodeValue.IoErrSeek:              "SQLITE_IOERR_SEEK",
	SqliteResultCodeValue.IoErrDeleteNoEnt:       "SQLITE_IOERR_DELETE_NOENT",
	SqliteResultCodeValue.IoErrMMap:              "SQLITE_IOERR_MMAP",
	SqliteResultCodeValue.IoErrGetTempPath:       "SQLITE_IOERR_GETTEMPPATH",
	SqliteResultCodeValue.IoErrConvPath:          "SQLITE_IOERR_CONVPATH",
	SqliteResultCodeValue.IoErrVNode:             "SQLITE_IOERR_VNODE",
	SqliteResultCodeValue.IoErrAuth:              "SQLITE_IOERR_AUTH",
	SqliteResultCodeValue.IoErrBeginAtomic:       "SQLITE_IOERR_BEGIN_ATOMIC",
	SqliteResultCodeValue.IoErrCommitAtomic:      "SQLITE_IOERR_COMMIT_ATOMIC",
	SqliteResultCodeValue.IoErrRollbackAtomic:    "SQLITE_IOERR_ROLLBACK_ATOMIC",
	SqliteResultCodeValue.IoErrData:              "SQLITE_IOERR_DATA",
	SqliteResultCodeValue.IoErrCorruptFs:         "SQLITE_IOERR_CORRUPTFS",
}

type SqliteResult struct {
	Code    SqliteResultCode
	Name    SqliteResultName
	Message string
	Meta    string
	Cause   *string
}

var SqliteResults = map[SqliteResultName]SqliteResult{
	SqliteResultNameValue.Ok: {
		Name:    SqliteResultNameValue.Ok,
		Code:    SqliteResultCodeValue.Ok,
		Message: "operation was successful",
		Meta:    "This means that the operation was successful and that there were no errors. Most other result codes indicate an error.",
	},
	SqliteResultNameValue.Error: {
		Name:    SqliteResultNameValue.Error,
		Code:    SqliteResultCodeValue.Error,
		Message: "unknown error",
		Meta:    "A generic error code that is used when no other more specific error code is available.",
	},
	SqliteResultNameValue.Internal: {
		Name:    SqliteResultNameValue.Internal,
		Code:    SqliteResultCodeValue.Internal,
		Message: "internal error",
		Meta:    "Indicates an internal malfunction. In a working version of SQLite, an application should never see this result code. If application does encounter this result code, it shows that there is a bug in the database engine.\n\n\t\tThis result code might be caused by a bug in SQLite. However, application-defined SQL functions or virtual tables, or VFSes, or other extensions can also cause this result code to be returned, so the problem might not be the fault of the core SQLite.",
	},
	SqliteResultNameValue.Perm: {
		Name:    SqliteResultNameValue.Perm,
		Code:    SqliteResultCodeValue.Perm,
		Message: "access mode could not be provided",
		Meta:    "The requested access mode for a newly created database could not be provided.",
	},
	SqliteResultNameValue.Abort: {
		Name:    SqliteResultNameValue.Abort,
		Code:    SqliteResultCodeValue.Abort,
		Message: "operation aborted",
		Meta:    "An operation was aborted prior to completion, usually be application request. See also: SQLITE_INTERRUPT.\n\n\t\tIf the callback function to sqlite3_exec() returns non-zero, then sqlite3_exec() will return SQLITE_ABORT.\n\n\t\tIf a ROLLBACK operation occurs on the same database connection as a pending read or write, then the pending read or write may fail with an SQLITE_ABORT or SQLITE_ABORT_ROLLBACK error.\n\n\t\tIn addition to being a result code, the SQLITE_ABORT value is also used as a conflict resolution mode returned from the sqlite3_vtab_on_conflict() interface.",
	},
	SqliteResultNameValue.Busy: {
		Name:    SqliteResultNameValue.Busy,
		Code:    SqliteResultCodeValue.Busy,
		Message: "database is busy",
		Meta:    "The database file could not be written (or in some cases read) because of concurrent activity by some other database connection, usually a database connection in a separate process.\n\n\t\tFor example, if process A is in the middle of a large write transaction and at the same time process B attempts to start a new write transaction, process B will get back an SQLITE_BUSY result because SQLite only supports one writer at a time. Process B will need to wait for process A to finish its transaction before starting a new transaction. The sqlite3_busy_timeout() and sqlite3_busy_handler() interfaces and the busy_timeout pragma are available to process B to help it deal with SQLITE_BUSY errors.\n\n\t\tAn SQLITE_BUSY error can occur at any point in a transaction: when the transaction is first started, during any write or update operations, or when the transaction commits. To avoid encountering SQLITE_BUSY errors in the middle of a transaction, the application can use BEGIN IMMEDIATE instead of just BEGIN to start a transaction. The BEGIN IMMEDIATE command might itself return SQLITE_BUSY, but if it succeeds, then SQLite guarantees that no subsequent operations on the same database through the next COMMIT will return SQLITE_BUSY.\n\n\t\tSee also: SQLITE_BUSY_RECOVERY and SQLITE_BUSY_SNAPSHOT.\n\n\t\tThe SQLITE_BUSY result code differs from SQLITE_LOCKED in that SQLITE_BUSY indicates a conflict with a separate database connection, probably in a separate process, whereas SQLITE_LOCKED indicates a conflict within the same database connection (or sometimes a database connection with a shared cache).",
	},
	SqliteResultNameValue.Locked: {
		Name:    SqliteResultNameValue.Locked,
		Code:    SqliteResultCodeValue.Locked,
		Message: "lock detected",
		Meta:    "A write operation could not continue because of a conflict within the same database connection or a conflict with a different database connection that uses a shared cache.\n\n\t\tFor example, a DROP TABLE statement cannot be run while another thread is reading from that table on the same database connection because dropping the table would delete the table out from under the concurrent reader.\n\n\t\tThe SQLITE_LOCKED result code differs from SQLITE_BUSY in that SQLITE_LOCKED indicates a conflict on the same database connection (or on a connection with a shared cache) whereas SQLITE_BUSY indicates a conflict with a different database connection, probably in a different process.",
	},
	SqliteResultNameValue.NoMem: {
		Name:    SqliteResultNameValue.NoMem,
		Code:    SqliteResultCodeValue.NoMem,
		Message: "could not allocate memory",
		Meta:    "SQLite was unable to allocate all the memory it needed to complete the operation. In other words, an internal call to sqlite3_malloc() or sqlite3_realloc() has failed in a case where the memory being allocated was required in order to continue the operation.",
	},
	SqliteResultNameValue.ReadOnly: {
		Name:    SqliteResultNameValue.ReadOnly,
		Code:    SqliteResultCodeValue.ReadOnly,
		Message: "write not permitted",
		Meta:    "An attempt is made to alter some data for which the current database connection does not have write permission.",
	},
	SqliteResultNameValue.Locked: {
		Name:    SqliteResultNameValue.Locked,
		Code:    SqliteResultCodeValue.Locked,
		Message: "operation interrupted",
		Meta:    "An operation was interrupted by the sqlite3_interrupt() interface. See also: SQLITE_ABORT",
	},
	SqliteResultNameValue.Locked: {
		Name:    SqliteResultNameValue.Locked,
		Code:    SqliteResultCodeValue.Locked,
		Message: "",
		Meta:    "",
	},
	/*
		(10) SQLITE_IOERR
		The SQLITE_IOERR result code says that the operation could not finish because the operating system reported an I/O error.

		A full disk drive will normally give an SQLITE_FULL error rather than an SQLITE_IOERR error.

		There are many different extended result codes for I/O errors that identify the specific I/O operation that failed.

		(11) SQLITE_CORRUPT
		The SQLITE_CORRUPT result code indicates that the database file has been corrupted. See the How To Corrupt Your Database Files for further discussion on how corruption can occur.

		(12) SQLITE_NOTFOUND
		The SQLITE_NOTFOUND result code is exposed in three ways:

		SQLITE_NOTFOUND can be returned by the sqlite3_file_control() interface to indicate that the file control opcode passed as the third argument was not recognized by the underlying VFS.

		SQLITE_NOTFOUND can also be returned by the xSetSystemCall() method of an sqlite3_vfs object.

		SQLITE_NOTFOUND can be returned by sqlite3_vtab_rhs_value() to indicate that the right-hand operand of a constraint is not available to the xBestIndex method that made the call.

		The SQLITE_NOTFOUND result code is also used internally by the SQLite implementation, but those internal uses are not exposed to the application.

		(13) SQLITE_FULL
		The SQLITE_FULL result code indicates that a write could not complete because the disk is full. Note that this error can occur when trying to write information into the main database file, or it can also occur when writing into temporary disk files.

		Sometimes applications encounter this error even though there is an abundance of primary disk space because the error occurs when writing into temporary disk files on a system where temporary files are stored on a separate partition with much less space that the primary disk.

		(14) SQLITE_CANTOPEN
		The SQLITE_CANTOPEN result code indicates that SQLite was unable to open a file. The file in question might be a primary database file or one of several temporary disk files.

		(15) SQLITE_PROTOCOL
		The SQLITE_PROTOCOL result code indicates a problem with the file locking protocol used by SQLite. The SQLITE_PROTOCOL error is currently only returned when using WAL mode and attempting to start a new transaction. There is a race condition that can occur when two separate database connections both try to start a transaction at the same time in WAL mode. The loser of the race backs off and tries again, after a brief delay. If the same connection loses the locking race dozens of times over a span of multiple seconds, it will eventually give up and return SQLITE_PROTOCOL. The SQLITE_PROTOCOL error should appear in practice very, very rarely, and only when there are many separate processes all competing intensely to write to the same database.

		(16) SQLITE_EMPTY
		The SQLITE_EMPTY result code is not currently used.

		(17) SQLITE_SCHEMA
		The SQLITE_SCHEMA result code indicates that the database schema has changed. This result code can be returned from sqlite3_step() for a prepared statement that was generated using sqlite3_prepare() or sqlite3_prepare16(). If the database schema was changed by some other process in between the time that the statement was prepared and the time the statement was run, this error can result.

		If a prepared statement is generated from sqlite3_prepare_v2() then the statement is automatically re-prepared if the schema changes, up to SQLITE_MAX_SCHEMA_RETRY times (default: 50). The sqlite3_step() interface will only return SQLITE_SCHEMA back to the application if the failure persists after these many retries.

		(18) SQLITE_TOOBIG
		The SQLITE_TOOBIG error code indicates that a string or BLOB was too large. The default maximum length of a string or BLOB in SQLite is 1,000,000,000 bytes. This maximum length can be changed at compile-time using the SQLITE_MAX_LENGTH compile-time option, or at run-time using the sqlite3_limit(db,SQLITE_LIMIT_LENGTH,...) interface. The SQLITE_TOOBIG error results when SQLite encounters a string or BLOB that exceeds the compile-time or run-time limit.

		The SQLITE_TOOBIG error code can also result when an oversized SQL statement is passed into one of the sqlite3_prepare_v2() interfaces. The maximum length of an SQL statement defaults to a much smaller value of 1,000,000,000 bytes. The maximum SQL statement length can be set at compile-time using SQLITE_MAX_SQL_LENGTH or at run-time using sqlite3_limit(db,SQLITE_LIMIT_SQL_LENGTH,...).

		(19) SQLITE_CONSTRAINT
		The SQLITE_CONSTRAINT error code means that an SQL constraint violation occurred while trying to process an SQL statement. Additional information about the failed constraint can be found by consulting the accompanying error message (returned via sqlite3_errmsg() or sqlite3_errmsg16()) or by looking at the extended error code.

		The SQLITE_CONSTRAINT code can also be used as the return value from the xBestIndex() method of a virtual table implementation. When xBestIndex() returns SQLITE_CONSTRAINT, that indicates that the particular combination of inputs submitted to xBestIndex() cannot result in a usable query plan and should not be given further consideration.

		(20) SQLITE_MISMATCH
		The SQLITE_MISMATCH error code indicates a datatype mismatch.

		SQLite is normally very forgiving about mismatches between the type of a value and the declared type of the container in which that value is to be stored. For example, SQLite allows the application to store a large BLOB in a column with a declared type of BOOLEAN. But in a few cases, SQLite is strict about types. The SQLITE_MISMATCH error is returned in those few cases when the types do not match.

		The rowid of a table must be an integer. Attempt to set the rowid to anything other than an integer (or a NULL which will be automatically converted into the next available integer rowid) results in an SQLITE_MISMATCH error.

		(21) SQLITE_MISUSE
		The SQLITE_MISUSE return code might be returned if the application uses any SQLite interface in a way that is undefined or unsupported. For example, using a prepared statement after that prepared statement has been finalized might result in an SQLITE_MISUSE error.

		SQLite tries to detect misuse and report the misuse using this result code. However, there is no guarantee that the detection of misuse will be successful. Misuse detection is probabilistic. Applications should never depend on an SQLITE_MISUSE return value.

		If SQLite ever returns SQLITE_MISUSE from any interface, that means that the application is incorrectly coded and needs to be fixed. Do not ship an application that sometimes returns SQLITE_MISUSE from a standard SQLite interface because that application contains potentially serious bugs.

		(22) SQLITE_NOLFS
		The SQLITE_NOLFS error can be returned on systems that do not support large files when the database grows to be larger than what the filesystem can handle. "NOLFS" stands for "NO Large File Support".

	*/

	SqliteResultNameValue.Locked: {
		Name:    SqliteResultNameValue.Locked,
		Code:    SqliteResultCodeValue.Locked,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.NoMem: {
		Name:    SqliteResultNameValue.NoMem,
		Code:    SqliteResultCodeValue.NoMem,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ReadOnly: {
		Name:    SqliteResultNameValue.ReadOnly,
		Code:    SqliteResultCodeValue.ReadOnly,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Interrupt: {
		Name:    SqliteResultNameValue.Interrupt,
		Code:    SqliteResultCodeValue.Interrupt,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErr: {
		Name:    SqliteResultNameValue.IoErr,
		Code:    SqliteResultCodeValue.IoErr,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Corrupt: {
		Name:    SqliteResultNameValue.Corrupt,
		Code:    SqliteResultCodeValue.Corrupt,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.NotFound: {
		Name:    SqliteResultNameValue.NotFound,
		Code:    SqliteResultCodeValue.NotFound,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Full: {
		Name:    SqliteResultNameValue.Full,
		Code:    SqliteResultCodeValue.Full,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.CantOpen: {
		Name:    SqliteResultNameValue.CantOpen,
		Code:    SqliteResultCodeValue.CantOpen,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Protocol: {
		Name:    SqliteResultNameValue.Protocol,
		Code:    SqliteResultCodeValue.Protocol,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Empty: {
		Name:    SqliteResultNameValue.Empty,
		Code:    SqliteResultCodeValue.Empty,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Schema: {
		Name:    SqliteResultNameValue.Schema,
		Code:    SqliteResultCodeValue.Schema,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.TooBig: {
		Name:    SqliteResultNameValue.TooBig,
		Code:    SqliteResultCodeValue.TooBig,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Constraint: {
		Name:    SqliteResultNameValue.Constraint,
		Code:    SqliteResultCodeValue.Constraint,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Mismatch: {
		Name:    SqliteResultNameValue.Mismatch,
		Code:    SqliteResultCodeValue.Mismatch,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Misuse: {
		Name:    SqliteResultNameValue.Misuse,
		Code:    SqliteResultCodeValue.Misuse,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.NolFs: {
		Name:    SqliteResultNameValue.NolFs,
		Code:    SqliteResultCodeValue.NolFs,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Auth: {
		Name:    SqliteResultNameValue.Auth,
		Code:    SqliteResultCodeValue.Auth,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Format: {
		Name:    SqliteResultNameValue.Format,
		Code:    SqliteResultCodeValue.Format,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Range: {
		Name:    SqliteResultNameValue.Range,
		Code:    SqliteResultCodeValue.Range,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.NotAdb: {
		Name:    SqliteResultNameValue.NotAdb,
		Code:    SqliteResultCodeValue.NotAdb,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Notice: {
		Name:    SqliteResultNameValue.Notice,
		Code:    SqliteResultCodeValue.Notice,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Warning: {
		Name:    SqliteResultNameValue.Warning,
		Code:    SqliteResultCodeValue.Warning,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Row: {
		Name:    SqliteResultNameValue.Row,
		Code:    SqliteResultCodeValue.Row,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.Done: {
		Name:    SqliteResultNameValue.Done,
		Code:    SqliteResultCodeValue.Done,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.OkLoadPermanently: {
		Name:    SqliteResultNameValue.OkLoadPermanently,
		Code:    SqliteResultCodeValue.OkLoadPermanently,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ErrorMissingCollSeq: {
		Name:    SqliteResultNameValue.ErrorMissingCollSeq,
		Code:    SqliteResultCodeValue.ErrorMissingCollSeq,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.BusyRecovery: {
		Name:    SqliteResultNameValue.BusyRecovery,
		Code:    SqliteResultCodeValue.BusyRecovery,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.LockedSharedCache: {
		Name:    SqliteResultNameValue.LockedSharedCache,
		Code:    SqliteResultCodeValue.LockedSharedCache,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ReadonlyRecovery: {
		Name:    SqliteResultNameValue.ReadonlyRecovery,
		Code:    SqliteResultCodeValue.ReadonlyRecovery,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrRead: {
		Name:    SqliteResultNameValue.IoErrRead,
		Code:    SqliteResultCodeValue.IoErrRead,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.CorruptVab: {
		Name:    SqliteResultNameValue.CorruptVab,
		Code:    SqliteResultCodeValue.CorruptVab,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.CantOpenNoteMpDir: {
		Name:    SqliteResultNameValue.CantOpenNoteMpDir,
		Code:    SqliteResultCodeValue.CantOpenNoteMpDir,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ConstraintCheck: {
		Name:    SqliteResultNameValue.ConstraintCheck,
		Code:    SqliteResultCodeValue.ConstraintCheck,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.AuthUser: {
		Name:    SqliteResultNameValue.AuthUser,
		Code:    SqliteResultCodeValue.AuthUser,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.NoticeRecoverWal: {
		Name:    SqliteResultNameValue.NoticeRecoverWal,
		Code:    SqliteResultCodeValue.NoticeRecoverWal,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.WarningAutoIndex: {
		Name:    SqliteResultNameValue.WarningAutoIndex,
		Code:    SqliteResultCodeValue.WarningAutoIndex,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ErrorRetry: {
		Name:    SqliteResultNameValue.ErrorRetry,
		Code:    SqliteResultCodeValue.ErrorRetry,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.AbortRollback: {
		Name:    SqliteResultNameValue.AbortRollback,
		Code:    SqliteResultCodeValue.AbortRollback,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.BusySnapshot: {
		Name:    SqliteResultNameValue.BusySnapshot,
		Code:    SqliteResultCodeValue.BusySnapshot,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.LockedVtab: {
		Name:    SqliteResultNameValue.LockedVtab,
		Code:    SqliteResultCodeValue.LockedVtab,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ReadonlyCantLock: {
		Name:    SqliteResultNameValue.ReadonlyCantLock,
		Code:    SqliteResultCodeValue.ReadonlyCantLock,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrShortRead: {
		Name:    SqliteResultNameValue.IoErrShortRead,
		Code:    SqliteResultCodeValue.IoErrShortRead,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.CorruptSequence: {
		Name:    SqliteResultNameValue.CorruptSequence,
		Code:    SqliteResultCodeValue.CorruptSequence,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.CantOpenIsDir: {
		Name:    SqliteResultNameValue.CantOpenIsDir,
		Code:    SqliteResultCodeValue.CantOpenIsDir,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ConstraintCommitHook: {
		Name:    SqliteResultNameValue.ConstraintCommitHook,
		Code:    SqliteResultCodeValue.ConstraintCommitHook,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.NoticeRecoverRollback: {
		Name:    SqliteResultNameValue.NoticeRecoverRollback,
		Code:    SqliteResultCodeValue.NoticeRecoverRollback,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ErrorSnapshot: {
		Name:    SqliteResultNameValue.ErrorSnapshot,
		Code:    SqliteResultCodeValue.ErrorSnapshot,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.BusyTimeout: {
		Name:    SqliteResultNameValue.BusyTimeout,
		Code:    SqliteResultCodeValue.BusyTimeout,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ReadonlyRollback: {
		Name:    SqliteResultNameValue.ReadonlyRollback,
		Code:    SqliteResultCodeValue.ReadonlyRollback,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrWrite: {
		Name:    SqliteResultNameValue.IoErrWrite,
		Code:    SqliteResultCodeValue.IoErrWrite,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.CorruptIndex: {
		Name:    SqliteResultNameValue.CorruptIndex,
		Code:    SqliteResultCodeValue.CorruptIndex,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.CantOpenFullPath: {
		Name:    SqliteResultNameValue.CantOpenFullPath,
		Code:    SqliteResultCodeValue.CantOpenFullPath,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ConstraintForeignKey: {
		Name:    SqliteResultNameValue.ConstraintForeignKey,
		Code:    SqliteResultCodeValue.ConstraintForeignKey,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ReadOnlyDbMoved: {
		Name:    SqliteResultNameValue.ReadOnlyDbMoved,
		Code:    SqliteResultCodeValue.ReadOnlyDbMoved,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrFSync: {
		Name:    SqliteResultNameValue.IoErrFSync,
		Code:    SqliteResultCodeValue.IoErrFSync,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.CantOpenConvPath: {
		Name:    SqliteResultNameValue.CantOpenConvPath,
		Code:    SqliteResultCodeValue.CantOpenConvPath,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ConstraintFunction: {
		Name:    SqliteResultNameValue.ConstraintFunction,
		Code:    SqliteResultCodeValue.ConstraintFunction,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ReadOnlyCantInit: {
		Name:    SqliteResultNameValue.ReadOnlyCantInit,
		Code:    SqliteResultCodeValue.ReadOnlyCantInit,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrDirFSync: {
		Name:    SqliteResultNameValue.IoErrDirFSync,
		Code:    SqliteResultCodeValue.IoErrDirFSync,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.CantOpenDirtyWal: {
		Name:    SqliteResultNameValue.CantOpenDirtyWal,
		Code:    SqliteResultCodeValue.CantOpenDirtyWal,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ConstraintNotNull: {
		Name:    SqliteResultNameValue.ConstraintNotNull,
		Code:    SqliteResultCodeValue.ConstraintNotNull,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ReadOnlyDirectory: {
		Name:    SqliteResultNameValue.ReadOnlyDirectory,
		Code:    SqliteResultCodeValue.ReadOnlyDirectory,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrTruncate: {
		Name:    SqliteResultNameValue.IoErrTruncate,
		Code:    SqliteResultCodeValue.IoErrTruncate,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.CantOpenSymlink: {
		Name:    SqliteResultNameValue.CantOpenSymlink,
		Code:    SqliteResultCodeValue.CantOpenSymlink,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ConstraintPrimaryKey: {
		Name:    SqliteResultNameValue.ConstraintPrimaryKey,
		Code:    SqliteResultCodeValue.ConstraintPrimaryKey,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrFStat: {
		Name:    SqliteResultNameValue.IoErrFStat,
		Code:    SqliteResultCodeValue.IoErrFStat,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ConstraintTrigger: {
		Name:    SqliteResultNameValue.ConstraintTrigger,
		Code:    SqliteResultCodeValue.ConstraintTrigger,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrUnlock: {
		Name:    SqliteResultNameValue.IoErrUnlock,
		Code:    SqliteResultCodeValue.IoErrUnlock,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ConstraintUnique: {
		Name:    SqliteResultNameValue.ConstraintUnique,
		Code:    SqliteResultCodeValue.ConstraintUnique,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrRdLock: {
		Name:    SqliteResultNameValue.IoErrRdLock,
		Code:    SqliteResultCodeValue.IoErrRdLock,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ConstraintVTab: {
		Name:    SqliteResultNameValue.ConstraintVTab,
		Code:    SqliteResultCodeValue.ConstraintVTab,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrDelete: {
		Name:    SqliteResultNameValue.IoErrDelete,
		Code:    SqliteResultCodeValue.IoErrDelete,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ConstraintRowId: {
		Name:    SqliteResultNameValue.ConstraintRowId,
		Code:    SqliteResultCodeValue.ConstraintRowId,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrBlocked: {
		Name:    SqliteResultNameValue.IoErrBlocked,
		Code:    SqliteResultCodeValue.IoErrBlocked,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ConstraintPinned: {
		Name:    SqliteResultNameValue.ConstraintPinned,
		Code:    SqliteResultCodeValue.ConstraintPinned,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrNoMem: {
		Name:    SqliteResultNameValue.IoErrNoMem,
		Code:    SqliteResultCodeValue.IoErrNoMem,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.ConstraintDataType: {
		Name:    SqliteResultNameValue.ConstraintDataType,
		Code:    SqliteResultCodeValue.ConstraintDataType,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrAccess: {
		Name:    SqliteResultNameValue.IoErrAccess,
		Code:    SqliteResultCodeValue.IoErrAccess,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrCheckReservedLock: {
		Name:    SqliteResultNameValue.IoErrCheckReservedLock,
		Code:    SqliteResultCodeValue.IoErrCheckReservedLock,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrLock: {
		Name:    SqliteResultNameValue.IoErrLock,
		Code:    SqliteResultCodeValue.IoErrLock,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrClose: {
		Name:    SqliteResultNameValue.IoErrClose,
		Code:    SqliteResultCodeValue.IoErrClose,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrDirClose: {
		Name:    SqliteResultNameValue.IoErrDirClose,
		Code:    SqliteResultCodeValue.IoErrDirClose,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrShmOpen: {
		Name:    SqliteResultNameValue.IoErrShmOpen,
		Code:    SqliteResultCodeValue.IoErrShmOpen,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrShmSize: {
		Name:    SqliteResultNameValue.IoErrShmSize,
		Code:    SqliteResultCodeValue.IoErrShmSize,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrShmLock: {
		Name:    SqliteResultNameValue.IoErrShmLock,
		Code:    SqliteResultCodeValue.IoErrShmLock,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrShmMap: {
		Name:    SqliteResultNameValue.IoErrShmMap,
		Code:    SqliteResultCodeValue.IoErrShmMap,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrSeek: {
		Name:    SqliteResultNameValue.IoErrSeek,
		Code:    SqliteResultCodeValue.IoErrSeek,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrDeleteNoEnt: {
		Name:    SqliteResultNameValue.IoErrDeleteNoEnt,
		Code:    SqliteResultCodeValue.IoErrDeleteNoEnt,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrMMap: {
		Name:    SqliteResultNameValue.IoErrMMap,
		Code:    SqliteResultCodeValue.IoErrMMap,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrGetTempPath: {
		Name:    SqliteResultNameValue.IoErrGetTempPath,
		Code:    SqliteResultCodeValue.IoErrGetTempPath,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrConvPath: {
		Name:    SqliteResultNameValue.IoErrConvPath,
		Code:    SqliteResultCodeValue.IoErrConvPath,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrVNode: {
		Name:    SqliteResultNameValue.IoErrVNode,
		Code:    SqliteResultCodeValue.IoErrVNode,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrAuth: {
		Name:    SqliteResultNameValue.IoErrAuth,
		Code:    SqliteResultCodeValue.IoErrAuth,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrBeginAtomic: {
		Name:    SqliteResultNameValue.IoErrBeginAtomic,
		Code:    SqliteResultCodeValue.IoErrBeginAtomic,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrCommitAtomic: {
		Name:    SqliteResultNameValue.IoErrCommitAtomic,
		Code:    SqliteResultCodeValue.IoErrCommitAtomic,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrRollbackAtomic: {
		Name:    SqliteResultNameValue.IoErrRollbackAtomic,
		Code:    SqliteResultCodeValue.IoErrRollbackAtomic,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrData: {
		Name:    SqliteResultNameValue.IoErrData,
		Code:    SqliteResultCodeValue.IoErrData,
		Message: "",
		Meta:    "",
	},
	SqliteResultNameValue.IoErrCorruptFs: {
		Name:    SqliteResultNameValue.IoErrCorruptFs,
		Code:    SqliteResultCodeValue.IoErrCorruptFs,
		Message: "",
		Meta:    "",
	},
}

/*

(23) SQLITE_AUTH
The SQLITE_AUTH error is returned when the authorizer callback indicates that an SQL statement being prepared is not authorized.

(24) SQLITE_FORMAT
The SQLITE_FORMAT error code is not currently used by SQLite.

(25) SQLITE_RANGE
The SQLITE_RANGE error indices that the parameter number argument to one of the sqlite3_bind routines or the column number in one of the sqlite3_column routines is out of range.

(26) SQLITE_NOTADB
When attempting to open a file, the SQLITE_NOTADB error indicates that the file being opened does not appear to be an SQLite database file.

(27) SQLITE_NOTICE
The SQLITE_NOTICE result code is not returned by any C/C++ interface. However, SQLITE_NOTICE (or rather one of its extended error codes) is sometimes used as the first argument in an sqlite3_log() callback to indicate that an unusual operation is taking place.

(28) SQLITE_WARNING
The SQLITE_WARNING result code is not returned by any C/C++ interface. However, SQLITE_WARNING (or rather one of its extended error codes) is sometimes used as the first argument in an sqlite3_log() callback to indicate that an unusual and possibly ill-advised operation is taking place.

(100) SQLITE_ROW
The SQLITE_ROW result code returned by sqlite3_step() indicates that another row of output is available.

(101) SQLITE_DONE
The SQLITE_DONE result code indicates that an operation has completed. The SQLITE_DONE result code is most commonly seen as a return value from sqlite3_step() indicating that the SQL statement has run to completion. But SQLITE_DONE can also be returned by other multi-step interfaces such as sqlite3_backup_step().

(256) SQLITE_OK_LOAD_PERMANENTLY
The sqlite3_load_extension() interface loads an extension into a single database connection. The default behavior is for that extension to be automatically unloaded when the database connection closes. However, if the extension entry point returns SQLITE_OK_LOAD_PERMANENTLY instead of SQLITE_OK, then the extension remains loaded into the process address space after the database connection closes. In other words, the xDlClose methods of the sqlite3_vfs object is not called for the extension when the database connection closes.

The SQLITE_OK_LOAD_PERMANENTLY return code is useful to loadable extensions that register new VFSes, for example.

(257) SQLITE_ERROR_MISSING_COLLSEQ
The SQLITE_ERROR_MISSING_COLLSEQ result code means that an SQL statement could not be prepared because a collating sequence named in that SQL statement could not be located.

Sometimes when this error code is encountered, the sqlite3_prepare_v2() routine will convert the error into SQLITE_ERROR_RETRY and try again to prepare the SQL statement using a different query plan that does not require the use of the unknown collating sequence.

(261) SQLITE_BUSY_RECOVERY
The SQLITE_BUSY_RECOVERY error code is an extended error code for SQLITE_BUSY that indicates that an operation could not continue because another process is busy recovering a WAL mode database file following a crash. The SQLITE_BUSY_RECOVERY error code only occurs on WAL mode databases.

(262) SQLITE_LOCKED_SHAREDCACHE
The SQLITE_LOCKED_SHAREDCACHE result code indicates that access to an SQLite data record is blocked by another database connection that is using the same record in shared cache mode. When two or more database connections share the same cache and one of the connections is in the middle of modifying a record in that cache, then other connections are blocked from accessing that data while the modifications are on-going in order to prevent the readers from seeing a corrupt or partially completed change.

(264) SQLITE_READONLY_RECOVERY
The SQLITE_READONLY_RECOVERY error code is an extended error code for SQLITE_READONLY. The SQLITE_READONLY_RECOVERY error code indicates that a WAL mode database cannot be opened because the database file needs to be recovered and recovery requires write access but only read access is available.

(266) SQLITE_IOERR_READ
The SQLITE_IOERR_READ error code is an extended error code for SQLITE_IOERR indicating an I/O error in the VFS layer while trying to read from a file on disk. This error might result from a hardware malfunction or because a filesystem came unmounted while the file was open.

(267) SQLITE_CORRUPT_VTAB
The SQLITE_CORRUPT_VTAB error code is an extended error code for SQLITE_CORRUPT used by virtual tables. A virtual table might return SQLITE_CORRUPT_VTAB to indicate that content in the virtual table is corrupt.

(270) SQLITE_CANTOPEN_NOTEMPDIR
The SQLITE_CANTOPEN_NOTEMPDIR error code is no longer used.

(275) SQLITE_CONSTRAINT_CHECK
The SQLITE_CONSTRAINT_CHECK error code is an extended error code for SQLITE_CONSTRAINT indicating that a CHECK constraint failed.

(279) SQLITE_AUTH_USER
The SQLITE_AUTH_USER error code is an extended error code for SQLITE_AUTH indicating that an operation was attempted on a database for which the logged in user lacks sufficient authorization.

(283) SQLITE_NOTICE_RECOVER_WAL
The SQLITE_NOTICE_RECOVER_WAL result code is passed to the callback of sqlite3_log() when a WAL mode database file is recovered.

(284) SQLITE_WARNING_AUTOINDEX
The SQLITE_WARNING_AUTOINDEX result code is passed to the callback of sqlite3_log() whenever automatic indexing is used. This can serve as a warning to application designers that the database might benefit from additional indexes.

(513) SQLITE_ERROR_RETRY
The SQLITE_ERROR_RETRY is used internally to provoke sqlite3_prepare_v2() (or one of its sibling routines for creating prepared statements) to try again to prepare a statement that failed with an error on the previous attempt.

(516) SQLITE_ABORT_ROLLBACK
The SQLITE_ABORT_ROLLBACK error code is an extended error code for SQLITE_ABORT indicating that an SQL statement aborted because the transaction that was active when the SQL statement first started was rolled back. Pending write operations always fail with this error when a rollback occurs. A ROLLBACK will cause a pending read operation to fail only if the schema was changed within the transaction being rolled back.

(517) SQLITE_BUSY_SNAPSHOT
The SQLITE_BUSY_SNAPSHOT error code is an extended error code for SQLITE_BUSY that occurs on WAL mode databases when a database connection tries to promote a read transaction into a write transaction but finds that another database connection has already written to the database and thus invalidated prior reads.

The following scenario illustrates how an SQLITE_BUSY_SNAPSHOT error might arise:

Process A starts a read transaction on the database and does one or more SELECT statement. Process A keeps the transaction open.
Process B updates the database, changing values previous read by process A.
Process A now tries to write to the database. But process A's view of the database content is now obsolete because process B has modified the database file after process A read from it. Hence process A gets an SQLITE_BUSY_SNAPSHOT error.
(518) SQLITE_LOCKED_VTAB
The SQLITE_LOCKED_VTAB result code is not used by the SQLite core, but it is available for use by extensions. Virtual table implementations can return this result code to indicate that they cannot complete the current operation because of locks held by other threads or processes.

The R-Tree extension returns this result code when an attempt is made to update the R-Tree while another prepared statement is actively reading the R-Tree. The update cannot proceed because any change to an R-Tree might involve reshuffling and rebalancing of nodes, which would disrupt read cursors, causing some rows to be repeated and other rows to be omitted.

(520) SQLITE_READONLY_CANTLOCK
The SQLITE_READONLY_CANTLOCK error code is an extended error code for SQLITE_READONLY. The SQLITE_READONLY_CANTLOCK error code indicates that SQLite is unable to obtain a read lock on a WAL mode database because the shared-memory file associated with that database is read-only.

(522) SQLITE_IOERR_SHORT_READ
The SQLITE_IOERR_SHORT_READ error code is an extended error code for SQLITE_IOERR indicating that a read attempt in the VFS layer was unable to obtain as many bytes as was requested. This might be due to a truncated file.

(523) SQLITE_CORRUPT_SEQUENCE
The SQLITE_CORRUPT_SEQUENCE result code means that the schema of the sqlite_sequence table is corrupt. The sqlite_sequence table is used to help implement the AUTOINCREMENT feature. The sqlite_sequence table should have the following format:

  CREATE TABLE sqlite_sequence(name,seq);

If SQLite discovers that the sqlite_sequence table has any other format, it returns the SQLITE_CORRUPT_SEQUENCE error.

(526) SQLITE_CANTOPEN_ISDIR
The SQLITE_CANTOPEN_ISDIR error code is an extended error code for SQLITE_CANTOPEN indicating that a file open operation failed because the file is really a directory.

(531) SQLITE_CONSTRAINT_COMMITHOOK
The SQLITE_CONSTRAINT_COMMITHOOK error code is an extended error code for SQLITE_CONSTRAINT indicating that a commit hook callback returned non-zero that thus caused the SQL statement to be rolled back.

(539) SQLITE_NOTICE_RECOVER_ROLLBACK
The SQLITE_NOTICE_RECOVER_ROLLBACK result code is passed to the callback of sqlite3_log() when a hot journal is rolled back.

(769) SQLITE_ERROR_SNAPSHOT
The SQLITE_ERROR_SNAPSHOT result code might be returned when attempting to start a read transaction on an historical version of the database by using the sqlite3_snapshot_open() interface. If the historical snapshot is no longer available, then the read transaction will fail with the SQLITE_ERROR_SNAPSHOT. This error code is only possible if SQLite is compiled with -DSQLITE_ENABLE_SNAPSHOT.

(773) SQLITE_BUSY_TIMEOUT
The SQLITE_BUSY_TIMEOUT error code indicates that a blocking Posix advisory file lock request in the VFS layer failed due to a timeout. Blocking Posix advisory locks are only available as a proprietary SQLite extension and even then are only supported if SQLite is compiled with the SQLITE_EANBLE_SETLK_TIMEOUT compile-time option.

(776) SQLITE_READONLY_ROLLBACK
The SQLITE_READONLY_ROLLBACK error code is an extended error code for SQLITE_READONLY. The SQLITE_READONLY_ROLLBACK error code indicates that a database cannot be opened because it has a hot journal that needs to be rolled back but cannot because the database is readonly.

(778) SQLITE_IOERR_WRITE
The SQLITE_IOERR_WRITE error code is an extended error code for SQLITE_IOERR indicating an I/O error in the VFS layer while trying to write into a file on disk. This error might result from a hardware malfunction or because a filesystem came unmounted while the file was open. This error should not occur if the filesystem is full as there is a separate error code (SQLITE_FULL) for that purpose.

(779) SQLITE_CORRUPT_INDEX
The SQLITE_CORRUPT_INDEX result code means that SQLite detected an entry is or was missing from an index. This is a special case of the SQLITE_CORRUPT error code that suggests that the problem might be resolved by running the REINDEX command, assuming no other problems exist elsewhere in the database file.

(782) SQLITE_CANTOPEN_FULLPATH
The SQLITE_CANTOPEN_FULLPATH error code is an extended error code for SQLITE_CANTOPEN indicating that a file open operation failed because the operating system was unable to convert the filename into a full pathname.

(787) SQLITE_CONSTRAINT_FOREIGNKEY
The SQLITE_CONSTRAINT_FOREIGNKEY error code is an extended error code for SQLITE_CONSTRAINT indicating that a foreign key constraint failed.

(1032) SQLITE_READONLY_DBMOVED
The SQLITE_READONLY_DBMOVED error code is an extended error code for SQLITE_READONLY. The SQLITE_READONLY_DBMOVED error code indicates that a database cannot be modified because the database file has been moved since it was opened, and so any attempt to modify the database might result in database corruption if the processes crashes because the rollback journal would not be correctly named.

(1034) SQLITE_IOERR_FSYNC
The SQLITE_IOERR_FSYNC error code is an extended error code for SQLITE_IOERR indicating an I/O error in the VFS layer while trying to flush previously written content out of OS and/or disk-control buffers and into persistent storage. In other words, this code indicates a problem with the fsync() system call in unix or the FlushFileBuffers() system call in windows.

(1038) SQLITE_CANTOPEN_CONVPATH
The SQLITE_CANTOPEN_CONVPATH error code is an extended error code for SQLITE_CANTOPEN used only by Cygwin VFS and indicating that the cygwin_conv_path() system call failed while trying to open a file. See also: SQLITE_IOERR_CONVPATH

(1043) SQLITE_CONSTRAINT_FUNCTION
The SQLITE_CONSTRAINT_FUNCTION error code is not currently used by the SQLite core. However, this error code is available for use by extension functions.

(1288) SQLITE_READONLY_CANTINIT
The SQLITE_READONLY_CANTINIT result code originates in the xShmMap method of a VFS to indicate that the shared memory region used by WAL mode exists buts its content is unreliable and unusable by the current process since the current process does not have write permission on the shared memory region. (The shared memory region for WAL mode is normally a file with a "-wal" suffix that is mmapped into the process space. If the current process does not have write permission on that file, then it cannot write into shared memory.)

Higher level logic within SQLite will normally intercept the error code and create a temporary in-memory shared memory region so that the current process can at least read the content of the database. This result code should not reach the application interface layer.

(1290) SQLITE_IOERR_DIR_FSYNC
The SQLITE_IOERR_DIR_FSYNC error code is an extended error code for SQLITE_IOERR indicating an I/O error in the VFS layer while trying to invoke fsync() on a directory. The unix VFS attempts to fsync() directories after creating or deleting certain files to ensure that those files will still appear in the filesystem following a power loss or system crash. This error code indicates a problem attempting to perform that fsync().

(1294) SQLITE_CANTOPEN_DIRTYWAL
The SQLITE_CANTOPEN_DIRTYWAL result code is not used at this time.

(1299) SQLITE_CONSTRAINT_NOTNULL
The SQLITE_CONSTRAINT_NOTNULL error code is an extended error code for SQLITE_CONSTRAINT indicating that a NOT NULL constraint failed.

(1544) SQLITE_READONLY_DIRECTORY
The SQLITE_READONLY_DIRECTORY result code indicates that the database is read-only because process does not have permission to create a journal file in the same directory as the database and the creation of a journal file is a prerequisite for writing.

(1546) SQLITE_IOERR_TRUNCATE
The SQLITE_IOERR_TRUNCATE error code is an extended error code for SQLITE_IOERR indicating an I/O error in the VFS layer while trying to truncate a file to a smaller size.

(1550) SQLITE_CANTOPEN_SYMLINK
The SQLITE_CANTOPEN_SYMLINK result code is returned by the sqlite3_open() interface and its siblings when the SQLITE_OPEN_NOFOLLOW flag is used and the database file is a symbolic link.

(1555) SQLITE_CONSTRAINT_PRIMARYKEY
The SQLITE_CONSTRAINT_PRIMARYKEY error code is an extended error code for SQLITE_CONSTRAINT indicating that a PRIMARY KEY constraint failed.

(1802) SQLITE_IOERR_FSTAT
The SQLITE_IOERR_FSTAT error code is an extended error code for SQLITE_IOERR indicating an I/O error in the VFS layer while trying to invoke fstat() (or the equivalent) on a file in order to determine information such as the file size or access permissions.

(1811) SQLITE_CONSTRAINT_TRIGGER
The SQLITE_CONSTRAINT_TRIGGER error code is an extended error code for SQLITE_CONSTRAINT indicating that a RAISE function within a trigger fired, causing the SQL statement to abort.

(2058) SQLITE_IOERR_UNLOCK
The SQLITE_IOERR_UNLOCK error code is an extended error code for SQLITE_IOERR indicating an I/O error within xUnlock method on the sqlite3_io_methods object.

(2067) SQLITE_CONSTRAINT_UNIQUE
The SQLITE_CONSTRAINT_UNIQUE error code is an extended error code for SQLITE_CONSTRAINT indicating that a UNIQUE constraint failed.

(2314) SQLITE_IOERR_RDLOCK
The SQLITE_IOERR_RDLOCK error code is an extended error code for SQLITE_IOERR indicating an I/O error within xLock method on the sqlite3_io_methods object while trying to obtain a read lock.

(2323) SQLITE_CONSTRAINT_VTAB
The SQLITE_CONSTRAINT_VTAB error code is not currently used by the SQLite core. However, this error code is available for use by application-defined virtual tables.

(2570) SQLITE_IOERR_DELETE
The SQLITE_IOERR_DELETE error code is an extended error code for SQLITE_IOERR indicating an I/O error within xDelete method on the sqlite3_vfs object.

(2579) SQLITE_CONSTRAINT_ROWID
The SQLITE_CONSTRAINT_ROWID error code is an extended error code for SQLITE_CONSTRAINT indicating that a rowid is not unique.

(2826) SQLITE_IOERR_BLOCKED
The SQLITE_IOERR_BLOCKED error code is no longer used.

(2835) SQLITE_CONSTRAINT_PINNED
The SQLITE_CONSTRAINT_PINNED error code is an extended error code for SQLITE_CONSTRAINT indicating that an UPDATE trigger attempted do delete the row that was being updated in the middle of the update.

(3082) SQLITE_IOERR_NOMEM
The SQLITE_IOERR_NOMEM error code is sometimes returned by the VFS layer to indicate that an operation could not be completed due to the inability to allocate sufficient memory. This error code is normally converted into SQLITE_NOMEM by the higher layers of SQLite before being returned to the application.

(3091) SQLITE_CONSTRAINT_DATATYPE
The SQLITE_CONSTRAINT_DATATYPE error code is an extended error code for SQLITE_CONSTRAINT indicating that an insert or update attempted to store a value inconsistent with the column's declared type in a table defined as STRICT.

(3338) SQLITE_IOERR_ACCESS
The SQLITE_IOERR_ACCESS error code is an extended error code for SQLITE_IOERR indicating an I/O error within the xAccess method on the sqlite3_vfs object.

(3594) SQLITE_IOERR_CHECKRESERVEDLOCK
The SQLITE_IOERR_CHECKRESERVEDLOCK error code is an extended error code for SQLITE_IOERR indicating an I/O error within the xCheckReservedLock method on the sqlite3_io_methods object.

(3850) SQLITE_IOERR_LOCK
The SQLITE_IOERR_LOCK error code is an extended error code for SQLITE_IOERR indicating an I/O error in the advisory file locking logic. Usually an SQLITE_IOERR_LOCK error indicates a problem obtaining a PENDING lock. However it can also indicate miscellaneous locking errors on some of the specialized VFSes used on Macs.

(4106) SQLITE_IOERR_CLOSE
The SQLITE_IOERR_CLOSE error code is an extended error code for SQLITE_IOERR indicating an I/O error within the xClose method on the sqlite3_io_methods object.

(4362) SQLITE_IOERR_DIR_CLOSE
The SQLITE_IOERR_DIR_CLOSE error code is no longer used.

(4618) SQLITE_IOERR_SHMOPEN
The SQLITE_IOERR_SHMOPEN error code is an extended error code for SQLITE_IOERR indicating an I/O error within the xShmMap method on the sqlite3_io_methods object while trying to open a new shared memory segment.

(4874) SQLITE_IOERR_SHMSIZE
The SQLITE_IOERR_SHMSIZE error code is an extended error code for SQLITE_IOERR indicating an I/O error within the xShmMap method on the sqlite3_io_methods object while trying to enlarge a "shm" file as part of WAL mode transaction processing. This error may indicate that the underlying filesystem volume is out of space.

(5130) SQLITE_IOERR_SHMLOCK
The SQLITE_IOERR_SHMLOCK error code is no longer used.

(5386) SQLITE_IOERR_SHMMAP
The SQLITE_IOERR_SHMMAP error code is an extended error code for SQLITE_IOERR indicating an I/O error within the xShmMap method on the sqlite3_io_methods object while trying to map a shared memory segment into the process address space.

(5642) SQLITE_IOERR_SEEK
The SQLITE_IOERR_SEEK error code is an extended error code for SQLITE_IOERR indicating an I/O error within the xRead or xWrite methods on the sqlite3_io_methods object while trying to seek a file descriptor to the beginning point of the file where the read or write is to occur.

(5898) SQLITE_IOERR_DELETE_NOENT
The SQLITE_IOERR_DELETE_NOENT error code is an extended error code for SQLITE_IOERR indicating that the xDelete method on the sqlite3_vfs object failed because the file being deleted does not exist.

(6154) SQLITE_IOERR_MMAP
The SQLITE_IOERR_MMAP error code is an extended error code for SQLITE_IOERR indicating an I/O error within the xFetch or xUnfetch methods on the sqlite3_io_methods object while trying to map or unmap part of the database file into the process address space.

(6410) SQLITE_IOERR_GETTEMPPATH
The SQLITE_IOERR_GETTEMPPATH error code is an extended error code for SQLITE_IOERR indicating that the VFS is unable to determine a suitable directory in which to place temporary files.

(6666) SQLITE_IOERR_CONVPATH
The SQLITE_IOERR_CONVPATH error code is an extended error code for SQLITE_IOERR used only by Cygwin VFS and indicating that the cygwin_conv_path() system call failed. See also: SQLITE_CANTOPEN_CONVPATH

(6922) SQLITE_IOERR_VNODE
The SQLITE_IOERR_VNODE error code is a code reserved for use by extensions. It is not used by the SQLite core.

(7178) SQLITE_IOERR_AUTH
The SQLITE_IOERR_AUTH error code is a code reserved for use by extensions. It is not used by the SQLite core.

(7434) SQLITE_IOERR_BEGIN_ATOMIC
The SQLITE_IOERR_BEGIN_ATOMIC error code indicates that the underlying operating system reported and error on the SQLITE_FCNTL_BEGIN_ATOMIC_WRITE file-control. This only comes up when SQLITE_ENABLE_ATOMIC_WRITE is enabled and the database is hosted on a filesystem that supports atomic writes.

(7690) SQLITE_IOERR_COMMIT_ATOMIC
The SQLITE_IOERR_COMMIT_ATOMIC error code indicates that the underlying operating system reported and error on the SQLITE_FCNTL_COMMIT_ATOMIC_WRITE file-control. This only comes up when SQLITE_ENABLE_ATOMIC_WRITE is enabled and the database is hosted on a filesystem that supports atomic writes.

(7946) SQLITE_IOERR_ROLLBACK_ATOMIC
The SQLITE_IOERR_ROLLBACK_ATOMIC error code indicates that the underlying operating system reported and error on the SQLITE_FCNTL_ROLLBACK_ATOMIC_WRITE file-control. This only comes up when SQLITE_ENABLE_ATOMIC_WRITE is enabled and the database is hosted on a filesystem that supports atomic writes.

(8202) SQLITE_IOERR_DATA
The SQLITE_IOERR_DATA error code is an extended error code for SQLITE_IOERR used only by checksum VFS shim to indicate that the checksum on a page of the database file is incorrect.

(8458) SQLITE_IOERR_CORRUPTFS
The SQLITE_IOERR_CORRUPTFS error code is an extended error code for SQLITE_IOERR used only by a VFS to indicate that a seek or read failure was due to the request not falling within the file's boundary rather than an ordinary device failure. This often indicates a corrupt filesystem.

*/
