# -*- coding: utf-8 -*-
# coding:utf-8
"""
Created on Fri Nov  2 15:25:00 2018

Version 1.0.0

@author: liaojiansheng
@email: liaojiansheng@ztgame.com liasece@gmail.com

***** 注意事项 *****
1   由于JS原生number类型无法表示完整值域的64位整数，最大仅可表示 2^53 次方个整数
    即 int64(−9007199254740992-9007199254740992) uint64(0-9007199254740992)
    如果需求的最大数值超过这些范围，请使用string等其他数据格式

[工具] 将 golang 中由 type TypeName struct{} 中定义的结构转化成为二进制流处理接口
类型定义 例：
// 支持行注释
type Example struct{
    UintType    uint32 // 整形支持 uint8 int8 uint16 int16 
                       // uint32 int32 uint64 int64 int uint
    FloatType   float32         // 浮点类型支持 float32 double(float64)
    StringType  string          // 字符串支持任何字符，包括\0字符，长度限制0-65535
    UniqueType  Example2        // 支持联合嵌套类型，但是该类型必需经过本工具处理
                                // ，以具备读写接口
    UniqueType2 pack.Example2   // golang 端支持引入其他包的结构，但是该类型必需经过本工具处理，以具备读写接口
    SliceType   []AnyType       // 支持 golang slice，值的类型可以为本述的任何类型，暂不支持多维slice
    MapType     map[int]AnyType // 支持 golang map ，键的类型必须为数字或字符串类型，值的类型可以为本述任何类型，暂不支持多维map
    Remark      int    `json:"otherName"` // 支持将 golang 字段名映射为其他 ts 名字，如此，ts端将通过 obj.otherName 访问本字段
}

--初始版本，不保证数据完整性--

更新日志
v0.1.1
    修复一些在消息格式不统一时的溢出BUG
v0.2.0
    消息格式可追加，包括嵌套类型的消息，可从末尾删除
v1.0.0
    完善数据类型支持，支持字符串作为map键值，修复消息前后更新导致的越界undefined错误
v1.1.0
    可以在go结构体声明中添加 // jsonbinary:struct 来告诉脚本不要将该结构翻译为消息，如：
    // jsonbinary:struct
    type ItemMsg struct {
        item int
    }
v1.2.0
    可以在go结构体声明中添加 // jsontype:U 来告诉脚本消息的类型，如：
    // jsonbinary:U
    type ItemMsg struct {
        item int
    }
    则 MsgIdToType(ItemMsgID) == 'U'
"""
import re
import os
import random
import sys
import codecs
import getopt

backstr = []

# 备份代码中的字符串
def backupstr(content):
    def _addtag(matched):
        res = "\'>>>>ja"+str(len(backstr))+"ja<<<<\'"
        backstr.append(matched.group(0))
        return res
    reg = re.compile("\'(.*?)\'")
    content = re.sub(reg, _addtag, content)
    reg = re.compile("\"(.*?)\"")
    content = re.sub(reg, _addtag, content)
    return content
    
# 恢复代码中的备份字符串
def recvstr(content):
    reg = re.compile(r'(\'>>>>ja\d+ja<<<<)\'')
    def _recvtag(matched):
        reg1 = re.compile(r'[>>>>ja](\d+)[ja<<<<]')
        items1 = re.findall(reg1, matched.group(0))
        return backstr[int(items1[0])]
    content = re.sub(reg, _recvtag, content)
    return content

# 加载源代码文件
def LoadFile(filename):
    filecontent = ""
    try :
        file_obj = codecs.open(filename,'r', encoding='utf-8')
        filecontent = file_obj.read()
        filecontent = re.sub("\r\n", "\n", filecontent)
        filecontent = re.sub("\r", "\n", filecontent)
        file_obj.close()
    except IOError as e:
        print("打开文件[%s]失败：%s,%s" %(filename,e.errno, e.strerror))
    return filecontent

# 保存至目标文件
def save_to_file(file_name, contents):
    fh = codecs.open(file_name, 'w', encoding='utf-8')
    fh.write(contents)
    fh.close()
    
# 移除注释代码
def removenotuse(content):
    reg = re.compile("(/\*(\n|.)*?\*/)")
    # 移除行注释时，需要保留行
    def _saveline(matched):
        res = ""
        for i, ch in enumerate(matched.group(0)):
            if ch == '\n':
                res += '\n'
        return res
    content = re.sub(reg, _saveline, content)
    return content

# 移除注释代码
def removenotuseline(content):
    reg = re.compile("(//.*)")
    # 移除行注释时，需要保留行
    def _saveline(matched):
        res = ""
        for i, ch in enumerate(matched.group(0)):
            if ch == '\n':
                res += '\n'
        return res
    content = re.sub(reg, _saveline, content)
    return content

# 格式化代码，去除行首不可见字符
def fmtcode(content) :
    reg = re.compile("\n[\t ]+")
    content = re.sub(reg, "\n", content)
    return content

# 格式化代码，孤立注释（上下皆空行）
def fmtcodesub(content) :
    reg = re.compile("\n[\t ]*\n[\t ]*([\t ]*//.*\n)+[\t ]*\n[\t ]*")
    content = re.sub(reg, "\n", content)

    reg = re.compile("^([\t ]*//.*\n)+[\t ]*\n[\t ]*")
    content = re.sub(reg, "\n", content)

    reg = re.compile("\n[\t ]*\n[\t ]*([\t ]*//.*\n)*([\t ]*//.*\n?)[\t ]*$")
    content = re.sub(reg, "\n", content)
    return content

# 格式化输出代码
def fmtcodeout(content) :
    reg = re.compile(";\s+")
    content = re.sub(reg, ";\n", content)

    reg = re.compile("{\s+")
    content = re.sub(reg, "{\n", content)

    reg = re.compile("}\s*^(else)")
    content = re.sub(reg, "}\n", content)

    reg = re.compile("\s+=\s+")
    content = re.sub(reg, " = ", content)

    reg = re.compile("\)\s*{")
    content = re.sub(reg, " ) { ", content)

    reg = re.compile("\s+{^}")
    content = re.sub(reg, " {", content)

    reg = re.compile("{[ \t\r]+\n")
    content = re.sub(reg, "{\n", content)

    reg = re.compile("\n\s+")
    content = re.sub(reg, "\n", content)
    
    res = ""
    tmplevel = 0
    tmplevel2 = 0
    lastch = ''
    for i,ch in enumerate(content):
        if ch == '(':
            tmplevel2 += 1
        elif ch == ')':
            tmplevel2 -= 1
        if ch == '{':
            tmplevel += 1
        elif ch == '}':
            tmplevel -= 1
        if lastch == '\n':
            for i in range(tmplevel):
                res += '\t'
            for i in range(tmplevel2):
                res += '\t'
        res += ch
        lastch = ch
    return res

# 从tag中获取实际json值
def getjsonname(field,remark):
    reg = re.compile("""['`"]json\s*:\s*['`"]\s*(\w+)\s*['`"]['`"]""")
    ms = re.finditer(reg, remark)
    for i in ms:
        return i.group(1)
    return field

# 判断一个结构代码是否是一个消息,消息内容需要带上消息的前注释
def isMsg(msgcontent):
    reg = re.compile("""//[\t ]*jsonbinary:struct\n[\t ]*""")
    ms = re.finditer(reg, msgcontent)
    for i in ms:
        return False
    return True

def getMsgType(name,msgcontent):
    reg = re.compile("""//[\t ]*jsontype:(\w)\n[\t ]*""")
    ms = re.finditer(reg, msgcontent)
    for i in ms:
        return i.group(1)
    return name[0]

def getgoint64(jsonname):
    read = jsonname+' = readBinaryInt64(data[offset:offset+8])\n\
            offset+=8\n'
    send = 'writeBinaryInt64(data[offset:offset+8], obj.'+jsonname+')\n\
            offset+=8\n'
    size = '8'
    return read,send,size
def getgoint32(jsonname):
    read = jsonname+' = readBinaryInt32(data[offset:offset+4])\n\
            offset+=4\n'
    send = 'writeBinaryInt32(data[offset:offset+4], obj.'+jsonname+')\n\
            offset+=4\n'
    size = '4'
    return read,send,size
def getgoint(jsonname):
    read = jsonname+' = readBinaryInt(data[offset:offset+4])\n\
            offset+=4\n'
    send = 'writeBinaryInt(data[offset:offset+4], obj.'+jsonname+')\n\
            offset+=4\n'
    size = '4'
    return read,send,size
def getgoint16(jsonname):
    read = jsonname+' = readBinaryInt16(data[offset:offset+2])\n\
            offset+=2\n'
    send = 'writeBinaryInt16(data[offset:offset+2], obj.'+jsonname+')\n\
            offset+=2\n'
    size = '2'
    return read,send,size
def getgoint8(jsonname):
    read = jsonname+' = readBinaryInt8(data[offset:offset+1])\n\
            offset+=1\n'
    send = 'writeBinaryInt8(data[offset:offset+1],obj.'+jsonname+')\n\
            offset+=1\n'
    size = '1'
    return read,send,size
def getgouint64(jsonname):
    read = jsonname+' = binary.BigEndian.Uint64(data[offset:offset+8])\n\
            offset+=8\n'
    send = 'binary.BigEndian.PutUint64(data[offset:offset+8], obj.'+jsonname+')\n\
            offset+=8\n'
    size = '8'
    return read,send,size
def getgouint32(jsonname):
    read = jsonname+' = binary.BigEndian.Uint32(data[offset:offset+4])\n\
            offset+=4\n'
    send = 'binary.BigEndian.PutUint32(data[offset:offset+4], obj.'+jsonname+')\n\
            offset+=4\n'
    size = '4'
    return read,send,size
def getgouint16(jsonname):
    read = jsonname+' = binary.BigEndian.Uint16(data[offset:offset+2])\n\
            offset+=2\n'
    send = 'binary.BigEndian.PutUint16(data[offset:offset+2], obj.'+jsonname+')\n\
            offset+=2\n'
    size = '2'
    return read,send,size
def getgouint8(jsonname):
    read = jsonname+' = readBinaryUint8(data[offset:offset+1])\n\
            offset+=1\n'
    send = 'writeBinaryUint8(data[offset:offset+1],obj.'+jsonname+')\n\
            offset+=1\n'
    size = '1'
    return read,send,size

def getgofloat(jsonname):
    read = jsonname+' = readBinaryFloat32(data[offset:offset+4])\n\
            offset+=4\n'
    send = 'writeBinaryFloat32(data[offset:offset+4], obj.'+jsonname+')\n\
            offset+=4\n'
    size = '4'
    return read,send,size
def getgodouble(jsonname):
    read = jsonname+' = readBinaryFloat64(data[offset:offset+8])\n\
            offset+=8\n'
    send = 'writeBinaryFloat64(data[offset:offset+8], obj.'+jsonname+')\n\
            offset+=8\n'
    size = '8'
    return read,send,size

def getgostringfunc():
    jsonname = "strfunc"
    read = 'func readBinaryString(data []byte) string {\
    '+jsonname+'len := binary.BigEndian.Uint32(data[:4])\n\
    if int('+jsonname+'len) + 4 > len(data) {\n\
        return ""\n\
    }\n\
    return string(data[4:4+'+jsonname+'len])\n\
}\n'
    send = 'func writeBinaryString(data []byte,obj string) int {\
    objlen := len(obj)\n\
    binary.BigEndian.PutUint32(data[:4],uint32(objlen))\n\
    copy(data[4:4+objlen], obj)\n\
    return 4+objlen\n\
}\n'
    return read,send

def getgobasefunc():
    jsonname = "strfunc"
    code = 'func bool2int(value bool) int {\
    if value {\n\
        return 1\n\
    }\n\
    return 0\n\
}\n\
func readBinaryInt64(data []byte) int64 {\n\
    // 大端模式\n\
    num := int64(0)\n\
    num |= int64(data[7]) << 0\n\
    num |= int64(data[6]) << 8\n\
    num |= int64(data[5]) << 16\n\
    num |= int64(data[4]) << 24\n\
    num |= int64(data[3]) << 32\n\
    num |= int64(data[2]) << 40\n\
    num |= int64(data[1]) << 48\n\
    num |= int64(data[0]) << 56\n\
    return num\n\
}\n\
func writeBinaryInt64(data []byte, num int64) {\n\
    // 大端模式\n\
    data[7] = byte((num >> 0) & 0xff)\n\
    data[6] = byte((num >> 8) & 0xff)\n\
    data[5] = byte((num >> 16) & 0xff)\n\
    data[4] = byte((num >> 24) & 0xff)\n\
    data[3] = byte((num >> 32) & 0xff)\n\
    data[2] = byte((num >> 40) & 0xff)\n\
    data[1] = byte((num >> 48) & 0xff)\n\
    data[0] = byte((num >> 56) & 0xff)\n\
}\n\
func readBinaryInt32(data []byte) int32 {\n\
    // 大端模式\n\
    num := int32(0)\n\
    num |= int32(data[3]) << 0\n\
    num |= int32(data[2]) << 8\n\
    num |= int32(data[1]) << 16\n\
    num |= int32(data[0]) << 24\n\
    return num\n\
}\n\
func writeBinaryInt32(data []byte, num int32) {\n\
    // 大端模式\n\
    data[3] = byte((num >> 0) & 0xff)\n\
    data[2] = byte((num >> 8) & 0xff)\n\
    data[1] = byte((num >> 16) & 0xff)\n\
    data[0] = byte((num >> 24) & 0xff)\n\
}\n\
func readBinaryInt(data []byte) int {\n\
    return int(readBinaryInt32(data))\n\
}\n\
func writeBinaryInt(data []byte, num int) {\n\
    writeBinaryInt32(data,int32(num))\n\
}\n\
func readBinaryInt16(data []byte) int16 {\n\
    // 大端模式\n\
    num := int16(0)\n\
    num |= int16(data[1]) << 0\n\
    num |= int16(data[0]) << 8\n\
    return num\n\
}\n\
func writeBinaryInt16(data []byte, num int16) {\n\
    // 大端模式\n\
    data[1] = byte((num >> 0) & 0xff)\n\
    data[0] = byte((num >> 8) & 0xff)\n\
}\n\
func readBinaryInt8(data []byte) int8 {\n\
    // 大端模式\n\
    num := int8(0)\n\
    num |= int8(data[0]) << 0\n\
    return num\n\
}\n\
func writeBinaryInt8(data []byte, num int8) {\n\
    // 大端模式\n\
    data[0] = byte(num)\n\
}\n\
func readBinaryBool(data []byte) bool {\n\
    // 大端模式\n\
    num := int8(0)\n\
    num |= int8(data[0]) << 0\n\
    return num>0\n\
}\n\
func writeBinaryBool(data []byte, num bool) {\n\
    // 大端模式\n\
    if num == true {\n\
    data[0] = byte(1)\n\
    } else {\n\
    data[0] = byte(0)\n\
    }\n\
}\n\
func readBinaryUint8(data []byte) uint8 {\n\
return uint8(data[0])\n\
}\n\
func writeBinaryUint8(data []byte, num uint8) {\n\
data[0] = byte(num)\n\
}\n\
func writeBinaryFloat32(data []byte, num float32) {\n\
bits := math.Float32bits(num)\n\
binary.BigEndian.PutUint32(data,bits)\n\
}\n\
func readBinaryFloat32(data []byte) float32 {\n\
bits := binary.BigEndian.Uint32(data)\n\
return math.Float32frombits(bits)\n\
}\n\
func writeBinaryFloat64(data []byte, num float64) {\n\
bits := math.Float64bits(num)\n\
binary.BigEndian.PutUint64(data,bits)\n\
}\n\
func readBinaryFloat64(data []byte) float64 {\n\
bits := binary.BigEndian.Uint64(data)\n\
return math.Float64frombits(bits)\n\
}\n'
    return code

def getgostring(jsonname):
    read = jsonname+' = readBinaryString(data[offset:])\n\
    '+'offset += 4 + len(obj.'+jsonname+')\n'
    send = 'writeBinaryString(data[offset:],obj.'+jsonname+')\n\
    '+'offset += 4 + len(obj.'+jsonname+')\n'
    size = '4 + len(obj.'+jsonname+')'
    return read,send,size

def getgobool(jsonname):
    read = jsonname+' = uint8(data[offset]) != 0\n\
            offset += 1\n'
    send = 'data[offset] = uint8(bool2int(obj.'+jsonname+'))\n\
            offset += 1\n'
    size = '1'
    return read,send,size

needpack = ''

def getgobybasetype(typestr,jsonname,ispoint,fieldnum):
    global needpack
    if ispoint:
        typestr = typestr[1:]
    if typestr == 'uint64' :
        ret,se,size = getgouint64(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'uint32' or typestr == 'uint':
        ret,se,size =  getgouint32(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'uint16':
        ret,se,size =  getgouint16(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'uint8' or typestr == 'byte':
        ret,se,size =  getgouint8(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'int64' :
        ret,se,size =  getgoint64(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'int32' :
        ret,se,size =  getgoint32(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'int':
        ret,se,size =  getgoint(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'int16':
        ret,se,size =  getgoint16(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'int8':
        ret,se,size =  getgoint8(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'bool':
        ret,se,size =  getgobool(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'string':
        ret,se,size =  getgostring(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'float32':
        ret,se,size =  getgofloat(jsonname)
        return 'obj.'+ret,se,size
    elif typestr == 'float64':
        ret,se,size =  getgodouble(jsonname)
        return 'obj.'+ret,se,size
    else :
        reg = re.compile("""[a-zA-Z0-9_]+""")
        ty = re.fullmatch(reg, typestr)
        if ty :
            # print(typestr)
            read = ""
            if ispoint:
                read += "offset += ReadMsg"+typestr+"ByBytes(data[offset:], obj."+jsonname+")\n"
            else :
                read += "offset += ReadMsg"+typestr+"ByBytes(data[offset:], &obj."+jsonname+")\n"
            send = ''
            if ispoint:
                send = "offset += WriteMsg"+typestr+"ByObj(data[offset:], obj."+jsonname+")\n"
            else :
                send = "offset += WriteMsg"+typestr+"ByObj(data[offset:], &obj."+jsonname+")\n"
            size = 'obj.'+jsonname+'.GetSize()'
            return read,send,size
        else :
            reg = re.compile("""([a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)""")
            ty = re.fullmatch(reg, typestr)
            if ty :
                needpack = '"'+ty.group(1)+'"'
                # print(typestr)
                read = ""
                if ispoint:
                    read += "offset += "+ty.group(1)+".ReadMsg"+ty.group(2)+'ByBytes(data[offset:], obj.'+jsonname+')\n'
                else :
                    read += "offset += "+ty.group(1)+".ReadMsg"+ty.group(2)+'ByBytes(data[offset:], &obj.'+jsonname+')\n'
                send = ''
                if ispoint:
                    send = "offset += "+ty.group(1)+".WriteMsg"+ty.group(2)+"ByObj(data[offset:], obj."+jsonname+")\n"
                else :
                    send = "offset += "+ty.group(1)+".WriteMsg"+ty.group(2)+"ByObj(data[offset:], &obj."+jsonname+")\n"
                size = 'obj.'+jsonname+'.GetSize()'
                return read,send,size
    # else :
        # print("Error Unknow typename:"+typestr)
    return "","",""
def isbasetype(typestr):
    if typestr == 'uint64' :
        return True
    elif typestr == 'uint32' or typestr == 'uint':
        return True
    elif typestr == 'uint16':
        return True
    elif typestr == 'uint8' or typestr == 'byte':
        return True
    elif typestr == 'int64' :
        return True
    elif typestr == 'int32' or typestr == 'int':
        return True
    elif typestr == 'int16':
        return True
    elif typestr == 'int8':
        return True
    elif typestr == 'bool':
        return True
    elif typestr == 'string':
        return True
    elif typestr == 'float32' or typestr == 'float64':
        return True
    return False
def getgobybasetypesub(typestr):
    global needpack
    if typestr == 'int64' :
        return 'readBinaryInt64(data[offset:offset+8])','writeBinaryInt64',8
    elif typestr == 'int32' :
        return 'readBinaryInt32(data[offset:offset+4])','writeBinaryInt32',4
    elif typestr == 'int':
        return 'readBinaryInt(data[offset:offset+4])','writeBinaryInt',4
    elif typestr == 'int16':
        return 'readBinaryInt16(data[offset:offset+2])','writeBinaryInt16',2
    elif typestr == 'int8':
        return 'readBinaryInt8(data[offset:offset+1])','writeBinaryInt8',1
    elif typestr == 'uint64' :
        return 'binary.BigEndian.Uint64(data[offset:offset+8])','binary.BigEndian.PutUint64',8
    elif typestr == 'uint32' or typestr == 'uint':
        return 'binary.BigEndian.Uint32(data[offset:offset+4])','binary.BigEndian.PutUint32',4
    elif typestr == 'uint16':
        return 'binary.BigEndian.Uint16(data[offset:offset+2])','binary.BigEndian.PutUint16',2
    elif typestr == 'uint8' or typestr == 'byte':
        return 'readBinaryUint8(data[offset:offset+1])','writeBinaryUint8',1
    elif typestr == 'bool':
        return 'readBinaryBool(data[offset:offset+1])','writeBinaryBool',1
    elif typestr == 'string':
        return 'readBinaryString(data[offset:])','writeBinaryString',0
    elif typestr == 'float32':
        return 'readBinaryFloat32(data[offset:offset+4])','writeBinaryFloat32',4
    elif typestr == 'float64':
        return 'readBinaryFloat64(data[offset:offset+8])','writeBinaryFloat64',8
    else :
        reg = re.compile("""[a-zA-Z0-9_]+""")
        ty = re.fullmatch(reg, typestr)
        if ty :
            # print(typestr)
            return "ReadMsg"+typestr+"ByBytes","WriteMsg"+typestr+"ByObj",0
        else :
            reg = re.compile("""([a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)""")
            ty = re.fullmatch(reg, typestr)
            if ty :
                needpack = '"'+ty.group(1)+'"'
                return ty.group(1)+".ReadMsg"+ty.group(2)+"ByBytes",ty.group(1)+".WriteMsg"+ty.group(2)+"ByObj"
    # else :
        # print("Error Unknow typename:"+typestr)
    return "","",0

# 数组转换代码
def getgoslice(typestr,jsonname,fieldnum):
    reg = re.compile("""\[\]([\w*]+)""")
    ty = re.fullmatch(reg, typestr)
    if ty:
        subtype = ty.group(1)
        getvaluest = ''
        getpointerst = ''
        ispoint = True
        sizerely = ''
        if subtype[0] != '*' and subtype != 'string':
            ispoint = False
            getvaluest = '*'
            getpointerst = '&'
        if subtype[0] == '*' :
            subtype = subtype[1:]
        subtypecode,subtypecodesend,subleng = getgobybasetypesub(subtype)
        readint,sendint,size = getgouint32(jsonname+"_slent")
        if subtypecode == "":
            reg = re.compile("""\[\]([\w*_]+)""")
            ty = re.fullmatch(reg, typestr)
            if ty :
                subtypecode,subtypecodesend,size,sizerely = getgoslice(subtype,ty.group(1),fieldnum+100000)
        if subtypecode == "":
            print("Error Unknow subtype in go:"+subtype)
            return "","","",""
        read = jsonname+"_slent := uint32(0)\n"
        # 越界判断
        read += 'if offset + '+size+' > data__len{\n\
                    return endpos\n\
                    }\n'
        read += jsonname+"_slen := 0\n" + readint + jsonname+"_slen = int("+jsonname+"_slent)\n"
        send = 'binary.BigEndian.PutUint32(data[offset:offset+4],uint32(len(obj.'+jsonname+')))\n\
                offset += 4\n'
        if isbasetype(subtype) :
            read += '\
                obj.'+jsonname+' = make('+typestr+','+jsonname+'_slen)\n\
                i'+str(fieldnum)+'i := 0\n\
                for '+jsonname+'_slen > i'+str(fieldnum)+'i {\n\
                if offset + '+str(subleng)+' > data__len{\n\
                    return endpos\n\
                }\n\
                tmp'+jsonname+'value := '+subtypecode+'\n\
                obj.'+jsonname+'[i'+str(fieldnum)+'i] = tmp'+jsonname+'value\n\
                offset += '+str(subleng)+'\n\
                i'+str(fieldnum)+'i++\n\
                }\n\
                '
            send += 'i'+str(fieldnum)+'i := 0\n\
                '+jsonname+'_slen := len(obj.'+jsonname+')\n\
                for '+jsonname+'_slen > i'+str(fieldnum)+'i {\n\
                '+subtypecodesend+'(data[offset:offset+'+str(subleng)+'],obj.'+jsonname+'[i'+str(fieldnum)+'i])\n\
                offset += '+str(subleng)+'\n\
                i'+str(fieldnum)+'i++\n\
                }\n\
                '
            if subtype == 'string' :
                sizerely += 'sizerely'+subtype+''+str(fieldnum)+' := func()int{\n\
                        resnum := 0\n\
                        i'+str(fieldnum)+'i := 0\n\
                        '+jsonname+'_slen := len(obj.'+jsonname+')\n\
                        for '+jsonname+'_slen > i'+str(fieldnum)+'i {\n\
                            resnum += len(obj.'+jsonname+'[i'+str(fieldnum)+'i]) + 2\n\
                            i'+str(fieldnum)+'i++\n\
                        }\n\
                        return resnum\n\
                    }\n'
                size += ' + sizerely'+subtype+''+str(fieldnum)+'()'
            else :
                size += ' + len(obj.'+jsonname+') * '+str(subleng)
        else :
            read += '\
                obj.'+jsonname+' = make('+typestr+','+jsonname+'_slen)\n\
                i'+str(fieldnum)+'i := 0\n\
                for '+jsonname+'_slen > i'+str(fieldnum)+'i {\n'
            if subtype != 'string':
                if ispoint :
                    read += 'tmpvalue'+subtype+' := &'+subtype+'{}\n\
                            offset += '+subtypecode+'(data[offset:],tmpvalue'+subtype+')\n\
                            obj.'+jsonname+'[i'+str(fieldnum)+'i] = tmpvalue'+subtype+'\n'
                else :
                    read += 'offset += '+subtypecode+'(data[offset:],&obj.'+jsonname+'[i'+str(fieldnum)+'i])\n'
                read += 'i'+str(fieldnum)+'i++\n\
                    }\n\
                    '
            else :
                read += 'obj.'+jsonname+'[i'+str(fieldnum)+'i] += '+subtypecode+'\n\
                        if offset + 2 > data__len{\n\
                            return endpos\n\
                        }\n\
                        offset += 2 + len(obj.'+jsonname+'[i'+str(fieldnum)+'i])\n\
                        i'+str(fieldnum)+'i++\n\
                    }\n\
                    '
            send += 'i'+str(fieldnum)+'i := 0\n\
                '+jsonname+'_slen := len(obj.'+jsonname+')\n\
                for '+jsonname+'_slen > i'+str(fieldnum)+'i {\n\
                offset += '+subtypecodesend+'(data[offset:],'+getpointerst+'obj.'+jsonname+'[i'+str(fieldnum)+'i])\n\
                i'+str(fieldnum)+'i++\n\
                }\n\
                '
            sizerely += 'sizerely'+subtype+''+str(fieldnum)+' := func()int{\n\
                    resnum := 0\n\
                    i'+str(fieldnum)+'i := 0\n\
                    '+jsonname+'_slen := len(obj.'+jsonname+')\n\
                    for '+jsonname+'_slen > i'+str(fieldnum)+'i {\n\
                        resnum += obj.'+jsonname+'[i'+str(fieldnum)+'i].GetSize()\n\
                        i'+str(fieldnum)+'i++\n\
                    }\n\
                    return resnum\n\
                }\n'
            size += ' + sizerely'+subtype+''+str(fieldnum)+'()'
        return read,send,size,sizerely
    # print("not's slice in go:"+typestr)
    return "","","",""

def getgomap(typestr,jsonname,fieldnum):
    reg = re.compile("""map\s*\[(\s*[a-zA-Z0-9_]+\s*)\]\s*([a-zA-Z0-9*_]+)\s*""")
    ty = re.fullmatch(reg, typestr)
    if ty:
        keytype = ty.group(1)
        valuetype = ty.group(2)
        getvaluest = ''
        getpointerst = ''
        getpointerfa = ''
        if valuetype[0] == '*':
            valuetype = valuetype[1:]
            getvaluest = '*'
            getpointerst = '&'
        else :
            getpointerfa = '&'
        keytypecode,keytypecodesend,keyleng = getgobybasetypesub(keytype)
        valuetypecode,valuetypecodesend,valueleng = getgobybasetypesub(valuetype)
        if keytypecode == "" or valuetypecode == "":
            print("Error Unknow in go keytypecode:"+keytypecode+";valuetypecode:"+valuetypecode)
            return "","","",""
        readint,sendint,size = getgouint32(jsonname+"_slent")
        read = jsonname+"_slent := uint32(0)\n"
        read += 'if offset + '+size+' > data__len{\n\
                    return endpos\n\
                }\n'
        read +=  readint + "\n"
        send = 'binary.BigEndian.PutUint32(data[offset:offset+2],uint32(len(obj.'+jsonname+')))\n\
                offset += 2\n'
        sizerely = ''

        catkeyv = ''
        catkeyvread = ''
        keyoffset = 'offset += '+str(keyleng)
        keyoffsetread = 'offset += '+str(keyleng)
        if keytype == 'string' :
            catkeyv = ''+jsonname+'_kcatlen := '
            catkeyvread = ''+jsonname+'_kcatlen := len(key'+jsonname+')'
            keyoffset = 'offset += '+jsonname+'_kcatlen'
            keyoffsetread = 'offset += '+jsonname+'_kcatlen + 2'
        catvaluev = ''
        catvaluevread = ''
        valueoffset = 'offset += '+str(valueleng)
        valueoffsetread = 'offset += '+str(valueleng)
        if valuetype == 'string' :
            valueleng = 2
            catvaluev = ''+jsonname+'_vcatlen := '
            catvaluevread = ''+jsonname+'_vcatlen := len(value'+jsonname+')'
            valueoffset = 'offset += '+jsonname+'_vcatlen'
            valueoffsetread = 'offset += '+jsonname+'_vcatlen + 2'

        if isbasetype(valuetype) :
            read += '\
                obj.'+jsonname+' = make('+typestr+')\n\
                i'+str(fieldnum)+'i := uint32(0)\n\
                for '+jsonname+'_slent > i'+str(fieldnum)+'i {\n\
                if offset + '+str(keyleng)+' > data__len{\n\
                    return endpos\n\
                }\n\
                key'+jsonname+' := '+keytypecode+'\n\
                '+catkeyvread+'\n\
                '+keyoffsetread+'\n\
                if offset + '+str(valueleng)+' > data__len{\n\
                    return endpos\n\
                }\n\
                value'+jsonname+' := '+valuetypecode+'\n\
                '+catvaluevread+'\n\
                '+valueoffsetread+'\n\
                obj.'+jsonname+'[key'+jsonname+'] = '+getpointerst+'value'+jsonname+'\n\
                i'+str(fieldnum)+'i++\n\
                }\n\
                '
            send += 'for '+jsonname+'key,'+jsonname+'value := range obj.'+jsonname+' {\n\
                '+catkeyv+keytypecodesend+'(data[offset:],'+jsonname+'key)\n\
                '+keyoffset+'\n\
                '+catvaluev+valuetypecodesend+'(data[offset:],'+jsonname+'value);\
                '+valueoffset+'\n\
                }\n\
                '
            if valuetype == 'string' and keytype != 'string' :
                # 值是string 键 是基础类型
                sizerely += 'sizerely'+valuetype+''+str(fieldnum)+' := func()int{\n\
                        resnum := 0\n\
                        for _,'+jsonname+'value := range obj.'+jsonname+' {\n\
                            resnum += len('+jsonname+'value) + 2\n\
                        }\n\
                        resnum += len(obj.'+jsonname+') * ('+str(keyleng)+')\n\
                        return resnum\n\
                    }\n'
            elif valuetype != 'string' and keytype == 'string' :
                # 值是基础类型 键是string
                sizerely += 'sizerely'+valuetype+''+str(fieldnum)+' := func()int{\n\
                        resnum := 0\n\
                        for '+jsonname+'key,_ := range obj.'+jsonname+' {\n\
                            resnum += len('+jsonname+'key) + 2\n\
                        }\n\
                        resnum += len(obj.'+jsonname+') * ('+str(valueleng)+')\n\
                        return resnum\n\
                    }\n'
            elif valuetype == 'string' and keytype == 'string' :
                # 值是string 键是string
                sizerely += 'sizerely'+valuetype+''+str(fieldnum)+' := func()int{\n\
                        resnum := 0\n\
                        for '+jsonname+'value,'+jsonname+'key := range obj.'+jsonname+' {\n\
                            resnum += len('+jsonname+'value) + 2\n\
                            resnum += len('+jsonname+'key) + 2\n\
                        }\n\
                        return resnum\n\
                    }\n'
        else :
            # 值 不是基础类型
            read += '\
                obj.'+jsonname+' = make('+typestr+')\n\
                i'+str(fieldnum)+'i := uint32(0)\n\
                for '+jsonname+'_slent > i'+str(fieldnum)+'i {\n\
                if offset + '+str(keyleng)+' > data__len{\n\
                    return endpos\n\
                }\n\
                key'+jsonname+' := '+keytypecode+'\n\
                '+catkeyvread+'\n\
                '+keyoffsetread+'\n\
                tmpvalue'+valuetype+' := '+valuetype+'{}\n\
                leng := '+valuetypecode+'(data[offset:],&tmpvalue'+valuetype+')\n\
                obj.'+jsonname+'[key'+jsonname+'] = '+getpointerst+'tmpvalue'+valuetype+'\n\
                offset += leng\n\
                i'+str(fieldnum)+'i++\n\
                }\n\
                '
            send += 'for '+jsonname+'key,'+jsonname+'value := range obj.'+jsonname+' {\n\
                '+catkeyv+keytypecodesend+'(data[offset:],'+jsonname+'key)\n\
                '+keyoffset+'\n\
                offset += '+valuetypecodesend+'(data[offset:],'+getpointerfa+''+jsonname+'value);\
                }\n\
                '
            if keytype == 'string' :
                # 键是string
                sizerely += 'sizerely'+valuetype+''+str(fieldnum)+' := func()int{\n\
                        resnum := 0\n\
                        for '+jsonname+'key,'+jsonname+'value := range obj.'+jsonname+' {\n\
                            resnum += '+jsonname+'value.GetSize()\n\
                            resnum += len('+jsonname+'key) + 2\n\
                        }\n\
                        return resnum\n\
                    }\n'
            else :
                # 键是基础类型
                sizerely += 'sizerely'+valuetype+''+str(fieldnum)+' := func()int{\n\
                        resnum := 0\n\
                        for _,'+jsonname+'value := range obj.'+jsonname+' {\n\
                            resnum += '+jsonname+'value.GetSize()\n\
                        }\n\
                        resnum += len(obj.'+jsonname+') * '+str(keyleng)+'\n\
                        return resnum\n\
                    }\n'
        if isbasetype(valuetype) and isbasetype(keytype) and valuetype!='string' and keytype != 'string' :
            # 键值都是可定量大小的类型
            size += ' + len(obj.'+jsonname+') * ('+str(keyleng)+' + '+str(valueleng)+')'
        else :
            # 需要依赖
            size += ' + sizerely'+valuetype+''+str(fieldnum)+'()'
        return read,send,size,sizerely
    # print("not's map in go:"+typestr)
    return "","","",""

def getgofield(field,fieldnum):
    name = recvstr(field.group(1))
    typestr = recvstr(field.group(2))
    # remark = recvstr(field.group(3))
    jsonname = name
    
    codestr = ""
    codestrsend = ""
    ispoint = False
    if typestr[0] == '*':
        ispoint = True
    codestr,codestrsend,size = getgobybasetype(typestr,jsonname,ispoint,fieldnum)
    if codestr != "":
        codestr = 'if offset + '+size+' > data__len{\n\
                    return endpos\n\
                    }\n' + codestr
        return jsonname,codestr,codestrsend,size,""
    codestr,codestrsend,size,sizerely = getgoslice(typestr,jsonname,fieldnum)
    if codestr != "":
        return jsonname,codestr,codestrsend,size,sizerely
    codestr,codestrsend,size,sizerely = getgomap(typestr,jsonname,fieldnum)
    if codestr != "":
        return jsonname,codestr,codestrsend,size,sizerely
    print("Error Unknow typename in GO:"+typestr)
    return "","","","",""
    # print("++"+jsonname+"++")

    # print("-"+typestr+"-")
    # print("-"+remark+"-")
    # print("done")

def getgomsgcode(content):
    # print(content)
    content = removenotuseline(content)
    res = ""
    ressend = ""
    size = ''
    sizerely = ''
    reg = re.compile("""[\t ]*(\w+)[\t ]+([^\s]+)[\t ]*([\w'`\": \t]+)?((\s*//.*\n?)*)""")
    gomsgs = re.finditer(reg, content)
    fieldnum = 0 
    for i in gomsgs:
        # print("    "+i.group(1))
        fieldnum+=1
        name,code,codesend,tmpsize,tmpsizerely = getgofield(i,fieldnum)
        if fieldnum > 1:
            if fieldnum % 5 == 0:
                size += ' + \n'
            else :
                size += ' + '
        size += tmpsize
        sizerely += tmpsizerely
        res += code
        ressend += codesend
    return res,ressend,size,sizerely
        # print(gettsfield(i))


def getgomsg(msgdef):
    name = recvstr(msgdef.group(3))
    content = recvstr(msgdef.group(4))
    # print(""+name)
    # print("-"+name+"-")
    # print("+"+content+"+")
    fieldcode,fieldcodesend,size,sizerely = getgomsgcode(content)
    getdatalencode = ''
    if fieldcode != '' :
        getdatalencode = 'data := indata[offset:offset+objsize]\n\
        offset = 0\n\
        data__len := len(data)'
    if size == '' :
        size = '0'
    res = 'func ReadMsg'+name+'ByBytes(indata []byte, obj *'+name+') int {\n\
        offset := 0\n\
        if len(indata) < 4 {\n\
        return 0\n\
        }\n\
        objsize := int(binary.BigEndian.Uint32(indata[offset:offset+4]))\n\
        offset += 4\n\
        if objsize == 0 {\n\
        return 4\n\
        }\n\
        if offset + objsize > len(indata){\n\
        return offset\n\
        }\n\
        endpos := offset+objsize\n\
        '+getdatalencode+'\n\
        '+fieldcode+'\nreturn endpos\n\
        }'
    ressend = 'func WriteMsg'+name+'ByObj(data []byte, obj *'+name+') int {\n\
        if obj == nil {\n\
            binary.BigEndian.PutUint32(data[0:4],0)\n\
            return 4\n\
        }\n\
        objsize := obj.GetSize() - 4\n\
        offset := 0\n\
        binary.BigEndian.PutUint32(data[offset:offset+4],uint32(objsize))\n\
        offset += 4\n\
        '+fieldcodesend+'\nreturn offset\n\
        }'
    ressize = 'func GetSize'+name+'(obj *'+name+') int {\n\
        if obj == nil {\n\
            return 4\n\
        }\n\
        '+sizerely+'\n\
        return 4 + '+size+'\n\
        }\n'
    return name,res,ressend,ressize

# go 解析函数
def getgoheadfunc(names,types,packname,proto):
    ressend = "func WriteBinary(msgname string, obj interface{}, data []byte) ([]byte,int) {\nswitch(msgname){\n"
    resinterface = ''
    resinterfaceread = ''
    resinterfacegetid = ''
    resinterfacegetname = ''
    resinterfacegetsize = ''
    resinterfacegetjsonstring = ''
    constmsgid = 'const (\n'
    constmsgname = 'const (\n'
    MsgIdToString = 'func MsgIdToString(id uint16) string {\nswitch(id){\n'
    StringToMsgId = 'func StringToMsgId(msgname string) uint16 {\nswitch(msgname){\n'
    MsgIdToType = 'func MsgIdToType(id uint16) rune {\nswitch(id){\n'
    # 消息号
    re,se,leng = getgobybasetypesub("uint16")
    times = 0
    for i in names:
        ressend += ' case '+i+'Name: {\n\
        '+se+'(data[:'+str(leng)+'], '+i+'ID)\n\
        offset := '+str(leng)+'\n\
        offset += WriteMsg'+ i +'ByObj(data[offset:],obj.(*'+i+'))\n\
        return data,offset\n}\n'
        MsgIdToString += ' case '+i+'ID: \nreturn '+i+'Name\n'
        StringToMsgId += 'case '+i+'Name: \nreturn '+i+'ID\n'
        MsgIdToType += ' case '+i+'ID: \nreturn rune(\''+types[times]+'\')\n'
        if proto == 'protobuf':
            resinterface += 'func (this *'+i+') WriteBinary(data []byte) int {\n\
    this.MarshalTo(data)\nreturn this.ProtoSize()\n}\n\n'
            resinterfaceread += 'func (this *'+i+') ReadBinary(data []byte) int {\n\
offset := len(data)\n\
this.Unmarshal(data)\n\
return offset\n}\n\n'
            resinterfacegetsize += 'func (this *'+i+') GetSize() int {\n\
            return this.ProtoSize()\n\
        }\n'
        else :
            resinterface += 'func (this *'+i+') WriteBinary(data []byte) int {\n\
    return WriteMsg'+ i +'ByObj(data,this)\n}\n\n'
            resinterfaceread += 'func (this *'+i+') ReadBinary(data []byte) int {\n\
return ReadMsg'+ i +'ByBytes(data, this)\n}\n\n'
            resinterfacegetsize += 'func (this *'+i+') GetSize() int {\n\
            return GetSize'+i+'(this)\n\
        }\n'
        resinterfacegetid += 'func (this *'+i+') GetMsgId() uint16 {\n\
            return '+i+'ID\n\
        }\n'
        resinterfacegetname += 'func (this *'+i+') GetMsgName() string {\n\
            return '+i+'Name\n\
        }\n'
        resinterfacegetjsonstring += 'func (this *'+i+') GetJson() string {\n\
            json,_ := json.Marshal(this)\n\
            return string(json)\n\
        }\n'

        constmsgid += ''+i+'ID = '+str(times+36)+'\n'
        constmsgname += ''+i+'Name = "'+packname+'.'+i+'"\n'
        times += 1
    ressend += 'default:\nlog.Error("未知的消息名称："+msgname)\n}\n}\n'
    ressend += 'default:\nreturn data,0\n}\n}\n'
    MsgIdToString += 'default:\nreturn ""\n}\n}\n'
    StringToMsgId += 'default:\nreturn 0\n}\n}\n'
    MsgIdToType += 'default:\nreturn rune(0)\n}\n}\n'
    constmsgid += ')\n'
    constmsgname += ')\n'
    
    return "" + constmsgid + constmsgname  + resinterface+ resinterfaceread + MsgIdToString + StringToMsgId + MsgIdToType + resinterfacegetid + resinterfacegetname + resinterfacegetsize + resinterfacegetjsonstring

def getgocode(packname,gomsgs,proto):
    msgNameList = []
    msgTypeList = []
    readstr,sendstr = getgostringfunc()
    basefunccode = getgobasefunc()
    res = readstr+sendstr + basefunccode
    for i in gomsgs:
        ismsg = isMsg(i.group(0))
        # print(i+"\n\n")
        name,code,codesend,codesize = getgomsg(i)
        if ismsg:
            msgNameList.append(name)
            msgTypeList.append(getMsgType(name,i.group(0)))
        res += code + "\n" + codesend + "\n" + codesize + "\n"
    # 根据是否使用其他proto决定是否需要编解码接口
    rescontent = getgoheadfunc(msgNameList,msgTypeList,packname,proto)
    if proto == 'protobuf':
        rescontent += ''
    else :
        rescontent += res
    # 返回最终结果
    return rescontent,packname
    
def procfile(filename,gofile,proto,showdetail):
    global backstr
    backstr = []
    
    if showdetail:
        print("gen: "+filename,end='')
    content = LoadFile(filename)
    # print(classs)
    # 保存代碼中的字符串
    content = backupstr(content)
    # 去除代码注释
    content = removenotuse(content)
    # 简单格式化代码
    content = fmtcode(content)
    content = fmtcodesub(content)

    reg = re.compile(r'package\s+(\w+)\n')
    pack = re.finditer(reg, content)
    packname = ""
    for i in pack:
        packname = i.group(1)
        break
    reg = re.compile(r'((//.*\n?)*)\n[\t ]*type\s+(\w+)\s+struct\s*{\s*([\s\S]*?)\s*\n?[\t ]*}')
    gomsgs = re.finditer(reg, content)

    # GOLang
    if gofile != "" :
        # go代码生成
        gomsgs = re.finditer(reg, content)
        gocontent,packname = getgocode(packname,gomsgs,proto)
        # 插入头
        totalcontent = 'package '+packname+'\n\n\
import (\n'
        # 根据是否使用其他协议区分引用的包
        if proto == 'protobuf':
            totalcontent += ''
        else :
            totalcontent += '\
\t"encoding/binary"\n\
\t"math"\n'
        # 包尾部
        totalcontent += '\
\t"encoding/json"\n\
\t'+needpack+'\n\
)\n\n'  
        # 实现体
        totalcontent += gocontent
        # 简单格式化代码
        totalcontent = fmtcodeout(totalcontent)
        # 保存go代码
        save_to_file(gofile,totalcontent)
        if showdetail:
            print(" Done")

def main():
    pwd = os.getcwd()
    goname = ''
    gofile = ''
    outtype = []
    proto = ''
    showdetail = False

    # 尝试解析命令行参数
    try:
        opts, args = getopt.getopt(sys.argv[1:],"hi:o:t:p:d",["help","infile=","outtype=","tsarg=","proto=","detail"])
    except getopt.GetoptError as e:
        print('\033[1;31mError:'+str(e)+'\nExample: go2ts.py -i <inputfile> -o <outtype> -t <tsarg> -p <proto> -d <detail>\033[0m')
        sys.exit(2)
    for opt, arg in opts:
        if opt in ("-h", "--help"):
            print('Usage: go2ts.py -i <inputfile> -o <outtype> -t <tsarg> -p <proto>')
            print('  -h, --help       Display this help and exit.')
            print('  -i, --infile     Message srtuct definition, must be golang source file.')
            print('  -o, --outtype    Output code type, e.g., go, ts.')
            print('  -p, --proto      Create Marshal/Unmarshal use other proto.')
            print('  -d, --detail     Show detail infomation.')
            sys.exit()
        elif opt in ("-i", "--infile"):
            index = arg.rfind('.')
            goname = arg[:index]
        elif opt in ("-o", "--outtype"):
            outtype.append(arg)
        elif opt in ("-p", "--proto"):
            proto = arg
        elif opt in ("-d", "--detail"):
            showdetail = True

    if 'go' in outtype:
        gofile = os.path.join(pwd,goname + "_binary.go")
    sourcegofile = os.path.join(pwd,goname+".go")

    procfile(sourcegofile,gofile,proto,showdetail)

if __name__ == '__main__':
    main()
