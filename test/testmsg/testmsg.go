package testmsg

type SModuleInfo struct {
	ModuleID   string
	ModuleAddr string
	// 服务器序号 重复不影响正常运行
	// 但是其改动会影响 配置读取/ModuleName/Log文件名
	ModuleNumber uint32
	// 服务器数字版本
	// 命名规则为： YYYYMMDDhhmm (年月日时分)
	Version uint64
}

type TestMyProto struct {
	Int                        int
	Byte_nil_1                 byte
	Int8_nil_1                 int8
	Int16_nil_1                int16
	Int32_nil_1                int32
	Int64_nil_1                int64
	Uint_nil_1                 uint
	Uint8_nil_1                uint8
	Uint16_nil_1               uint16
	Uint32_nil_1               uint32
	Uint64_nil_1               uint64
	String_nil_1               string
	MyType_nil_1               SModuleInfo
	MyTypeP_nil_1              *SModuleInfo
	MyTypeP1_nil_1             *SModuleInfo
	MyType_1_nil_1             SModuleInfo
	MyTypeP_1_nil_1            *SModuleInfo
	MyTypeP1_1_nil_1           *SModuleInfo
	SliceInt_nil_1             []int
	SliceByte_nil_1            []byte
	SliceInt8_nil_1            []int8
	SliceInt16_nil_1           []int16
	SliceInt32_nil_1           []int32
	SliceInt64_nil_1           []int64
	SliceUint_nil_1            []uint
	SliceUint8_nil_1           []uint8
	SliceUint16_nil_1          []uint16
	SliceUint32_nil_1          []uint32
	SliceUint64_nil_1          []uint64
	SliceMyType_nil_1          []SModuleInfo
	SliceMyTypeP_nil_1         []*SModuleInfo
	SliceMyType1_nil_1         []SModuleInfo
	SliceMyTypeP1_nil_1        []*SModuleInfo
	Int_nil_nil_1              int
	Byte_nil_nil_1             byte
	Int8_nil_nil_1             int8
	Int16_nil_nil_1            int16
	Int32_nil_nil_1            int32
	Int64_nil_nil_1            int64
	Uint_nil_nil_1             uint
	Uint8_nil_nil_1            uint8
	Uint16_nil_nil_1           uint16
	Uint32_nil_nil_1           uint32
	Uint64_nil_nil_1           uint64
	String_nil_nil_1           string
	MyType_nil_nil_1           SModuleInfo
	MyTypeP_nil_nil_1          *SModuleInfo
	MyTypeP1_nil_nil_1         *SModuleInfo
	MyType_1_nil_nil_1         SModuleInfo
	MyTypeP_1_nil_nil_1        *SModuleInfo
	MyTypeP1_1_nil_nil_1       *SModuleInfo
	SliceInt_nil_nil_1         []int
	SliceByte_nil_nil_1        []byte
	SliceInt8_nil_nil_1        []int8
	SliceInt16_nil_nil_1       []int16
	SliceInt32_nil_nil_1       []int32
	SliceInt64_nil_nil_1       []int64
	SliceUint_nil_nil_1        []uint
	SliceUint8_nil_nil_1       []uint8
	SliceUint16_nil_nil_1      []uint16
	SliceUint32_nil_nil_1      []uint32
	SliceUint64_nil_nil_1      []uint64
	SliceMyType_nil_nil_1      []SModuleInfo
	SliceMyTypeP_nil_nil_1     []*SModuleInfo
	SliceMyType1_nil_nil_1     []SModuleInfo
	SliceMyTypeP1_nil_nil_1    []*SModuleInfo
	SliceMyType_1_nil_nil_1    []SModuleInfo
	SliceMyTypeP_1_nil_nil_1   []*SModuleInfo
	SliceMyType1_1_nil_nil_1   []SModuleInfo
	SliceMyTypeP1_1_nil_nil_1  []*SModuleInfo
	MapIntInt_nil_nil_1        map[int]int
	MapInt8Int_nil_nil_1       map[int8]int
	MapInt16Int_nil_nil_1      map[int16]int
	MapIntByteInt_nil_nil_1    map[byte]int
	MapIntInt32Int_nil_nil_1   map[int32]int
	MapIntInt64Int_nil_nil_1   map[int64]int
	MapIntByte_nil_nil_1       map[int]byte
	MapInt8Byte_nil_nil_1      map[int8]byte
	MapInt16Byte_nil_nil_1     map[int16]byte
	MapIntByteByte_nil_nil_1   map[byte]byte
	MapIntInt32Byte_nil_nil_1  map[int32]byte
	MapIntInt64Byte_nil_nil_1  map[int64]byte
	MapIntInt8_nil_nil_1       map[int]int8
	MapInt8Int8_nil_nil_1      map[int8]int8
	MapInt16Int8_nil_nil_1     map[int16]int8
	MapIntByteInt8_nil_nil_1   map[byte]int8
	MapIntInt32Int8_nil_nil_1  map[int32]int8
	MapIntInt64Int8_nil_nil_1  map[int64]int8
	MapIntInt16_nil_nil_1      map[int]int16
	MapInt8Int16_nil_nil_1     map[int8]int16
	MapInt16Int16_nil_nil_1    map[int16]int16
	MapIntByteInt16_nil_nil_1  map[byte]int16
	MapIntInt32Int16_nil_nil_1 map[int32]int16
	MapIntInt64Int16_nil_nil_1 map[int64]int16
	MapIntInt32_nil_nil_1      map[int]int32
	Byte                       byte
	Int8                       int8
	Int16                      int16
	Int32                      int32
	Int64                      int64
	Uint                       uint
	Uint8                      uint8
	Uint16                     uint16
	Uint32                     uint32
	Uint64                     uint64
	String                     string
	MyType                     SModuleInfo
	MyTypeP                    *SModuleInfo
	MyTypeP1                   *SModuleInfo
	MyType_1                   SModuleInfo
	MyTypeP_1                  *SModuleInfo
	MyTypeP1_1                 *SModuleInfo
	SliceInt                   []int
	SliceByte                  []byte
	SliceInt8                  []int8
	SliceInt16                 []int16
	SliceInt32                 []int32
	SliceInt64                 []int64
	SliceUint                  []uint
	SliceUint8                 []uint8
	SliceUint16                []uint16
	SliceUint32                []uint32
	SliceUint64                []uint64
	SliceMyType                []SModuleInfo
	SliceMyTypeP               []*SModuleInfo
	SliceMyType1               []SModuleInfo
	SliceMyTypeP1              []*SModuleInfo
	Int_nil                    int
	Byte_nil                   byte
	Int8_nil                   int8
	Int16_nil                  int16
	Int32_nil                  int32
	Int64_nil                  int64
	Uint_nil                   uint
	Uint8_nil                  uint8
	Uint16_nil                 uint16
	Uint32_nil                 uint32
	Uint64_nil                 uint64
	String_nil                 string
	MyType_nil                 SModuleInfo
	MyTypeP_nil                *SModuleInfo
	MyTypeP1_nil               *SModuleInfo
	MyType_1_nil               SModuleInfo
	MyTypeP_1_nil              *SModuleInfo
	MyTypeP1_1_nil             *SModuleInfo
	SliceInt_nil               []int
	SliceByte_nil              []byte
	SliceInt8_nil              []int8
	SliceInt16_nil             []int16
	SliceInt32_nil             []int32
	SliceInt64_nil             []int64
	SliceUint_nil              []uint
	SliceUint8_nil             []uint8
	SliceUint16_nil            []uint16
	SliceUint32_nil            []uint32
	SliceUint64_nil            []uint64
	SliceMyType_nil            []SModuleInfo
	SliceMyTypeP_nil           []*SModuleInfo
	SliceMyType1_nil           []SModuleInfo
	SliceMyTypeP1_nil          []*SModuleInfo
	SliceMyType_1_nil          []SModuleInfo
	SliceMyTypeP_1_nil         []*SModuleInfo
	SliceMyType1_1_nil         []SModuleInfo
	SliceMyTypeP1_1_nil        []*SModuleInfo
	MapIntInt_nil              map[int]int
	MapInt8Int_nil             map[int8]int
	MapInt16Int_nil            map[int16]int
	MapIntByteInt_nil          map[byte]int
	MapIntInt32Int_nil         map[int32]int
	MapIntInt64Int_nil         map[int64]int
	MapIntByte_nil             map[int]byte
	MapInt8Byte_nil            map[int8]byte
	MapInt16Byte_nil           map[int16]byte
	MapIntByteByte_nil         map[byte]byte
	MapIntInt32Byte_nil        map[int32]byte
	MapIntInt64Byte_nil        map[int64]byte
	MapIntInt8_nil             map[int]int8
	MapInt8Int8_nil            map[int8]int8
	MapInt16Int8_nil           map[int16]int8
	MapIntByteInt8_nil         map[byte]int8
	MapIntInt32Int8_nil        map[int32]int8
	MapIntInt64Int8_nil        map[int64]int8
	MapIntInt16_nil            map[int]int16
	MapInt8Int16_nil           map[int8]int16
	MapInt16Int16_nil          map[int16]int16
	MapIntByteInt16_nil        map[byte]int16
	MapIntInt32Int16_nil       map[int32]int16
	MapIntInt64Int16_nil       map[int64]int16
	MapIntInt32_nil            map[int]int32
	MapInt8Int32_nil           map[int8]int32
	MapInt16Int32_nil          map[int16]int32
	MapIntByteInt32_nil        map[byte]int32
	MapIntInt32Int32_nil       map[int32]int32
	MapIntInt64Int32_nil       map[int64]int32
	MapIntInt64_nil            map[int]int64
	MapInt8Int64_nil           map[int8]int64
	MapInt16Int64_nil          map[int16]int64
	MapIntByteInt64_nil        map[byte]int64
	MapIntInt32Int64_nil       map[int32]int64
	MapIntInt64Int64_nil       map[int64]int64
	MapIntUint_nil             map[int]uint
	MapInt8Uint_nil            map[int8]uint
	MapInt16Uint_nil           map[int16]uint
	MapIntByteUint_nil         map[byte]uint
	MapIntInt32Uint_nil        map[int32]uint
	MapIntInt64Uint_nil        map[int64]uint
	MapIntUint8_nil            map[int]uint8
	MapInt8Uint8_nil           map[int8]uint8
	MapInt16Uint8_nil          map[int16]uint8
	MapIntByteUint8_nil        map[byte]uint8
	MapIntInt32Uint8_nil       map[int32]uint8
	MapIntInt64Uint8_nil       map[int64]uint8
	MapIntUint16_nil           map[int]uint16
	MapInt8Uint16_nil          map[int8]uint16
	MapInt16Uint16_nil         map[int16]uint16
	MapIntByteUint16_nil       map[byte]uint16
	MapIntInt32Uint16_nil      map[int32]uint16
	MapIntInt64Uint16_nil      map[int64]uint16
	MapIntUint32_nil           map[int]uint32
	MapInt8Uint32_nil          map[int8]uint32
	MapInt16Uint32_nil         map[int16]uint32
	MapIntByteUint32_nil       map[byte]uint32
	MapIntInt32Uint32_nil      map[int32]uint32
	MapIntInt64Uint32_nil      map[int64]uint32
	MapIntUint64_nil           map[int]uint64
	MapInt8Uint64_nil          map[int8]uint64
	MapInt16Uint64_nil         map[int16]uint64
	MapIntByteUint64_nil       map[byte]uint64
	MapIntInt32Uint64_nil      map[int32]uint64
	MapIntInt64Uint64_nil      map[int64]uint64
	MapIntMyType_nil           map[int]SModuleInfo
	MapInt8MyType_nil          map[int8]SModuleInfo
	MapInt16MyType_nil         map[int16]SModuleInfo
	MapIntByteMyType_nil       map[byte]SModuleInfo
	MapIntInt32MyType_nil      map[int32]SModuleInfo
	MapIntInt64MyType_nil      map[int64]SModuleInfo
	MapIntMyTypeP_nil          map[int]*SModuleInfo
	MapInt8MyTypeP_nil         map[int8]*SModuleInfo
	MapInt16MyTypeP_nil        map[int16]*SModuleInfo
	MapIntByteMyTypeP_nil      map[byte]*SModuleInfo
	MapIntInt32MyTypeP_nil     map[int32]*SModuleInfo
	MapIntInt64MyTypeP_nil     map[int64]*SModuleInfo
	Mapstring_nil              map[string]string
	Mapstring1_nil             map[string]string
	MapIntInt32Uint32_1_nil    map[int32]uint32
	MapIntInt64Uint32_1_nil    map[int64]uint32
	MapIntUint64_1_nil         map[int]uint64
	MapInt8Uint64_1_nil        map[int8]uint64
	MapInt16Uint64_1_nil       map[int16]uint64
	MapIntByteUint64_1_nil     map[byte]uint64
	MapIntInt32Uint64_1_nil    map[int32]uint64
	MapIntInt64Uint64_1_nil    map[int64]uint64
	SliceMyType_1              []SModuleInfo
	SliceMyTypeP_1             []*SModuleInfo
	SliceMyType1_1             []SModuleInfo
	SliceMyTypeP1_1            []*SModuleInfo
	MapIntInt                  map[int]int
	MapInt8Int                 map[int8]int
	MapInt16Int                map[int16]int
	MapIntByteInt              map[byte]int
	MapIntInt32Int             map[int32]int
	MapIntInt64Int             map[int64]int
	MapIntByte                 map[int]byte
	MapInt8Byte                map[int8]byte
	MapInt16Byte               map[int16]byte
	MapIntByteByte             map[byte]byte
	MapIntInt32Byte            map[int32]byte
	MapIntInt64Byte            map[int64]byte
	MapIntInt8                 map[int]int8
	MapInt8Int8                map[int8]int8
	MapInt16Int8               map[int16]int8
	MapIntByteInt8             map[byte]int8
	MapIntInt32Int8            map[int32]int8
	MapIntInt64Int8            map[int64]int8
	MapIntInt16                map[int]int16
	MapInt8Int16               map[int8]int16
	MapInt16Int16              map[int16]int16
	MapIntByteInt16            map[byte]int16
	MapIntInt32Int16           map[int32]int16
	MapIntInt64Int16           map[int64]int16
	MapIntInt32                map[int]int32
	MapInt8Int32               map[int8]int32
	MapInt16Int32              map[int16]int32
	MapIntByteInt32            map[byte]int32
	MapIntInt32Int32           map[int32]int32
	MapIntInt64Int32           map[int64]int32
	MapIntInt64                map[int]int64
	MapInt8Int64               map[int8]int64
	MapInt16Int64              map[int16]int64
	MapIntByteInt64            map[byte]int64
	MapIntInt32Int64           map[int32]int64
	MapIntInt64Int64           map[int64]int64
	MapIntUint                 map[int]uint
	MapInt8Uint                map[int8]uint
	MapInt16Uint               map[int16]uint
	MapIntByteUint             map[byte]uint
	MapIntInt32Uint            map[int32]uint
	MapIntInt64Uint            map[int64]uint
	MapIntUint8                map[int]uint8
	MapInt8Uint8               map[int8]uint8
	MapInt16Uint8              map[int16]uint8
	MapIntByteUint8            map[byte]uint8
	MapIntInt32Uint8           map[int32]uint8
	MapIntInt64Uint8           map[int64]uint8
	MapIntUint16               map[int]uint16
	MapInt8Uint16              map[int8]uint16
	MapInt16Uint16             map[int16]uint16
	MapIntByteUint16           map[byte]uint16
	MapIntInt32Uint16          map[int32]uint16
	MapIntInt64Uint16          map[int64]uint16
	MapIntUint32               map[int]uint32
	MapInt8Uint32              map[int8]uint32
	MapInt16Uint32             map[int16]uint32
	MapIntByteUint32           map[byte]uint32
	MapIntInt32Uint32          map[int32]uint32
	MapIntInt64Uint32          map[int64]uint32
	MapIntUint64               map[int]uint64
	MapInt8Uint64              map[int8]uint64
	MapInt16Uint64             map[int16]uint64
	MapIntByteUint64           map[byte]uint64
	MapIntInt32Uint64          map[int32]uint64
	MapIntInt64Uint64          map[int64]uint64
	MapIntMyType               map[int]SModuleInfo
	MapInt8MyType              map[int8]SModuleInfo
	MapInt16MyType             map[int16]SModuleInfo
	MapIntByteMyType           map[byte]SModuleInfo
	MapIntInt32MyType          map[int32]SModuleInfo
	MapIntInt64MyType          map[int64]SModuleInfo
	MapIntMyTypeP              map[int]*SModuleInfo
	MapInt8MyTypeP             map[int8]*SModuleInfo
	MapInt16MyTypeP            map[int16]*SModuleInfo
	MapIntByteMyTypeP          map[byte]*SModuleInfo
	MapIntInt32MyTypeP         map[int32]*SModuleInfo
	MapIntInt64MyTypeP         map[int64]*SModuleInfo
	Mapstring                  map[string]string
	Mapstring1                 map[string]string
	MapIntInt32Uint32_1        map[int32]uint32
	MapIntInt64Uint32_1        map[int64]uint32
	MapIntUint64_1             map[int]uint64
	MapInt8Uint64_1            map[int8]uint64
	MapInt16Uint64_1           map[int16]uint64
	MapIntByteUint64_1         map[byte]uint64
	MapIntInt32Uint64_1        map[int32]uint64
	MapIntInt64Uint64_1        map[int64]uint64
}

var testModuleInfo SModuleInfo = SModuleInfo{
	ModuleID:     "123",
	ModuleAddr:   "123",
	ModuleNumber: 123,
	Version:      123,
}

var TetstValue TestMyProto = TestMyProto{
	Int:          115,
	Byte:         115,
	Int8:         115,
	Int16:        115,
	Int32:        115,
	Int64:        115,
	Uint:         115,
	Uint8:        115,
	Uint16:       115,
	Uint32:       115,
	Uint64:       115,
	String:       "115",
	MyType:       testModuleInfo,
	MyTypeP:      &testModuleInfo,
	MyTypeP1:     nil,
	SliceInt:     []int{115},
	SliceByte:    []byte{115},
	SliceInt8:    []int8{115},
	SliceInt16:   []int16{115},
	SliceInt64:   []int64{115},
	SliceUint:    []uint{115},
	SliceUint8:   []uint8{115},
	SliceUint16:  []uint16{115},
	SliceUint32:  []uint32{115},
	SliceUint64:  []uint64{115},
	SliceMyType:  []SModuleInfo{testModuleInfo},
	SliceMyTypeP: []*SModuleInfo{&testModuleInfo, nil, &testModuleInfo},
	// SliceMyType1:       []SModuleInfo{testModuleInfo},
	// SliceMyTypeP1:      []*SModuleInfo{&testModuleInfo, nil, &testModuleInfo},
	MapIntInt:         map[int]int{1: 115, 2: 116, 3: 117},
	MapInt8Int:        map[int8]int{1: 115, 2: 116, 3: 117},
	MapInt16Int:       map[int16]int{1: 115, 2: 116, 3: 117},
	MapIntByteInt:     map[byte]int{1: 115, 2: 116, 3: 117},
	MapIntInt32Int:    map[int32]int{1: 115, 2: 116, 3: 117},
	MapIntByte:        map[int]byte{1: 115, 2: 116, 3: 117},
	MapInt8Byte:       map[int8]byte{1: 115, 2: 116, 3: 117},
	MapInt16Byte:      map[int16]byte{1: 115, 2: 116, 3: 117},
	MapIntByteByte:    map[byte]byte{1: 115, 2: 116, 3: 117},
	MapIntInt32Byte:   map[int32]byte{1: 115, 2: 116, 3: 117},
	MapIntInt64Byte:   map[int64]byte{1: 115, 2: 116, 3: 117},
	MapIntInt8:        map[int]int8{1: 115, 2: 116, 3: 117},
	MapInt8Int8:       map[int8]int8{1: 115, 2: 116, 3: 117},
	MapInt16Int8:      map[int16]int8{1: 115, 2: 116, 3: 117},
	MapIntByteInt8:    map[byte]int8{1: 115, 2: 116, 3: 117},
	MapIntInt32Int8:   map[int32]int8{1: 115, 2: 116, 3: 117},
	MapIntInt64Int8:   map[int64]int8{1: 115, 2: 116, 3: 117},
	MapIntInt16:       map[int]int16{1: 115, 2: 116, 3: 117},
	MapInt8Int16:      map[int8]int16{1: 115, 2: 116, 3: 117},
	MapInt16Int16:     map[int16]int16{1: 115, 2: 116, 3: 117},
	MapIntByteInt16:   map[byte]int16{1: 115, 2: 116, 3: 117},
	MapIntInt32Int16:  map[int32]int16{1: 115, 2: 116, 3: 117},
	MapIntInt64Int16:  map[int64]int16{1: 115, 2: 116, 3: 117},
	MapIntInt32:       map[int]int32{1: 115, 2: 116, 3: 117},
	MapInt8Int32:      map[int8]int32{1: 115, 2: 116, 3: 117},
	MapInt16Int32:     map[int16]int32{1: 115, 2: 116, 3: 117},
	MapIntByteInt32:   map[byte]int32{1: 115, 2: 116, 3: 117},
	MapIntInt32Int32:  map[int32]int32{1: 115, 2: 116, 3: 117},
	MapIntInt64Int32:  map[int64]int32{1: 115, 2: 116, 3: 117},
	MapIntInt64:       map[int]int64{1: 115, 2: 116, 3: 117},
	MapInt8Int64:      map[int8]int64{1: 115, 2: 116, 3: 117},
	MapInt16Int64:     map[int16]int64{1: 115, 2: 116, 3: 117},
	MapIntByteInt64:   map[byte]int64{1: 115, 2: 116, 3: 117},
	MapIntInt32Int64:  map[int32]int64{1: 115, 2: 116, 3: 117},
	MapIntInt64Int64:  map[int64]int64{1: 115, 2: 116, 3: 117},
	MapIntUint:        map[int]uint{1: 115, 2: 116, 3: 117},
	MapInt8Uint:       map[int8]uint{1: 115, 2: 116, 3: 117},
	MapInt16Uint:      map[int16]uint{1: 115, 2: 116, 3: 117},
	MapIntByteUint:    map[byte]uint{1: 115, 2: 116, 3: 117},
	MapIntInt32Uint:   map[int32]uint{1: 115, 2: 116, 3: 117},
	MapIntInt64Uint:   map[int64]uint{1: 115, 2: 116, 3: 117},
	MapIntUint8:       map[int]uint8{1: 115, 2: 116, 3: 117},
	MapInt8Uint8:      map[int8]uint8{1: 115, 2: 116, 3: 117},
	MapInt16Uint8:     map[int16]uint8{1: 115, 2: 116, 3: 117},
	MapIntByteUint8:   map[byte]uint8{1: 115, 2: 116, 3: 117},
	MapIntInt32Uint8:  map[int32]uint8{1: 115, 2: 116, 3: 117},
	MapIntInt64Uint8:  map[int64]uint8{1: 115, 2: 116, 3: 117},
	MapIntUint16:      map[int]uint16{1: 115, 2: 116, 3: 117},
	MapInt8Uint16:     map[int8]uint16{1: 115, 2: 116, 3: 117},
	MapInt16Uint16:    map[int16]uint16{1: 115, 2: 116, 3: 117},
	MapIntByteUint16:  map[byte]uint16{1: 115, 2: 116, 3: 117},
	MapIntInt32Uint16: map[int32]uint16{1: 115, 2: 116, 3: 117},
	MapIntInt64Uint16: map[int64]uint16{1: 115, 2: 116, 3: 117},
	MapIntUint32:      map[int]uint32{1: 115, 2: 116, 3: 117},
	MapInt8Uint32:     map[int8]uint32{1: 115, 2: 116, 3: 117},
	MapInt16Uint32:    map[int16]uint32{1: 115, 2: 116, 3: 117},
	MapIntByteUint32:  map[byte]uint32{1: 115, 2: 116, 3: 117},
	MapIntInt32Uint32: map[int32]uint32{1: 115, 2: 116, 3: 117},
	MapIntInt64Uint32: map[int64]uint32{1: 115, 2: 116, 3: 117},
	MapIntUint64:      map[int]uint64{1: 115, 2: 116, 3: 117},
	MapInt8Uint64:     map[int8]uint64{1: 115, 2: 116, 3: 117},
	MapInt16Uint64:    map[int16]uint64{1: 115, 2: 116, 3: 117},
	MapIntByteUint64:  map[byte]uint64{1: 115, 2: 116, 3: 117},
	MapIntInt32Uint64: map[int32]uint64{1: 115, 2: 116, 3: 117},
	MapIntInt64Uint64: map[int64]uint64{1: 115, 2: 116, 3: 117},
	MapIntMyType:      map[int]SModuleInfo{1: testModuleInfo, 2: testModuleInfo, 3: testModuleInfo, 4: testModuleInfo, 5: testModuleInfo, 6: testModuleInfo, 7: testModuleInfo, 8: testModuleInfo, 9: testModuleInfo},
	MapInt8MyType:     map[int8]SModuleInfo{1: testModuleInfo, 2: testModuleInfo, 3: testModuleInfo, 4: testModuleInfo, 5: testModuleInfo, 6: testModuleInfo, 7: testModuleInfo, 8: testModuleInfo, 9: testModuleInfo},
	MapInt16MyType:    map[int16]SModuleInfo{1: testModuleInfo, 2: testModuleInfo, 3: testModuleInfo, 4: testModuleInfo, 5: testModuleInfo, 6: testModuleInfo, 7: testModuleInfo, 8: testModuleInfo, 9: testModuleInfo},
	MapIntByteMyType:  map[byte]SModuleInfo{1: testModuleInfo, 2: testModuleInfo, 3: testModuleInfo, 4: testModuleInfo, 5: testModuleInfo, 6: testModuleInfo, 7: testModuleInfo, 8: testModuleInfo, 9: testModuleInfo},
	MapIntInt32MyType: map[int32]SModuleInfo{1: testModuleInfo, 2: testModuleInfo, 3: testModuleInfo, 4: testModuleInfo, 5: testModuleInfo, 6: testModuleInfo, 7: testModuleInfo, 8: testModuleInfo, 9: testModuleInfo},
	// MapIntInt64MyType: map[int64]SModuleInfo{1: testModuleInfo, 2: testModuleInfo, 3: testModuleInfo, 4: testModuleInfo, 5: testModuleInfo, 6: testModuleInfo, 7: testModuleInfo, 8: testModuleInfo, 9: testModuleInfo},
	MapIntMyTypeP:   map[int]*SModuleInfo{1: &testModuleInfo, 2: nil, 3: nil, 4: nil, 5: &testModuleInfo, 6: &testModuleInfo, 7: &testModuleInfo, 8: &testModuleInfo, 9: &testModuleInfo},
	MapInt8MyTypeP:  map[int8]*SModuleInfo{1: &testModuleInfo, 2: nil, 3: &testModuleInfo, 4: &testModuleInfo, 5: &testModuleInfo, 6: &testModuleInfo, 7: &testModuleInfo, 8: &testModuleInfo, 9: &testModuleInfo},
	MapInt16MyTypeP: map[int16]*SModuleInfo{1: &testModuleInfo, 2: nil, 3: &testModuleInfo, 4: &testModuleInfo, 5: &testModuleInfo, 6: &testModuleInfo, 7: &testModuleInfo, 8: &testModuleInfo, 9: &testModuleInfo},
	// MapIntByteMyTypeP:  map[byte]*SModuleInfo{1: &testModuleInfo, 2: nil, 3: &testModuleInfo, 4: &testModuleInfo, 5: &testModuleInfo, 6: &testModuleInfo, 7: &testModuleInfo, 8: &testModuleInfo, 9: &testModuleInfo},
	MapIntInt32MyTypeP: map[int32]*SModuleInfo{1: &testModuleInfo, 2: nil, 3: &testModuleInfo, 4: &testModuleInfo, 5: &testModuleInfo, 6: &testModuleInfo, 7: &testModuleInfo, 8: &testModuleInfo, 9: &testModuleInfo},
	MapIntInt64MyTypeP: map[int64]*SModuleInfo{1: &testModuleInfo, 2: nil, 3: &testModuleInfo, 4: &testModuleInfo, 5: &testModuleInfo, 6: &testModuleInfo, 7: &testModuleInfo, 8: &testModuleInfo, 9: &testModuleInfo},
	// Mapstring:          map[string]string{"1": "115"},
	Mapstring1: map[string]string{"1": "115"},
}
